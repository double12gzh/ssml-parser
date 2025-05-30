package main

import (
	"fmt"
	"log"
	"sort"
	"sync"
	"time"

	"ssml-parser/examples/wav"
	"ssml-parser/ssml"
)

func main() {
	fmt.Println("🚀 SSML 并发 TTS 简化示例")
	fmt.Println("==========================")

	// 包含多段文本和停顿的 SSML
	ssmlContent := `<?xml version="1.0" encoding="UTF-8"?>
<speak version="1.0" xml:lang="zh-CN">
    第一段文本内容
    <break time="1s"/>
    第二段文本内容  
    <break time="500ms"/>
    第三段文本内容
    <break time="300ms"/>
    第四段文本内容
</speak>`

	fmt.Println("📝 SSML 内容:")
	fmt.Println(ssmlContent)
	fmt.Println()

	// 解析并处理
	DemoConcurrentTTS(ssmlContent)
}

// SegmentWithAudio 音频片段结果
type SegmentWithAudio struct {
	Index     int
	Text      string
	AudioData *ssml.AudioData
	Error     error
}

// DemoConcurrentTTS 演示并发 TTS 处理
func DemoConcurrentTTS(ssmlContent string) {
	// 1. 解析 SSML
	fmt.Println("🔍 步骤 1: 解析 SSML...")
	parser := ssml.NewParser(nil)
	parseResult, err := parser.Parse(ssmlContent)
	if err != nil {
		log.Fatalf("解析失败: %v", err)
	}

	audioProcessor := ssml.NewAudioProcessor()
	result, err := audioProcessor.ProcessSSML(parseResult.Root)
	if err != nil {
		log.Fatalf("处理失败: %v", err)
	}

	fmt.Printf("✓ 找到 %d 个文本片段，%d 个停顿指令\n",
		len(result.Segments), len(result.Instructions))

	// 显示片段
	fmt.Println("\n📋 文本片段:")
	for i, segment := range result.Segments {
		fmt.Printf("  %d. \"%s\"\n", i+1, segment.Text)
	}

	// 显示停顿
	fmt.Println("\n⏸️  停顿指令:")
	for i, instruction := range result.Instructions {
		if instruction.Type == "break" {
			fmt.Printf("  %d. 停顿 %v\n", i+1, instruction.Duration)
		}
	}

	// 2. 并发生成音频
	fmt.Println("\n🎵 步骤 2: 并发生成音频...")
	start := time.Now()
	audioSegments := ConcurrentGenerateAudio(result.Segments)
	ttsTime := time.Since(start)

	fmt.Printf("✓ 并发处理完成，耗时: %v\n", ttsTime)

	// 3. 组装最终音频
	fmt.Println("\n🔧 步骤 3: 组装音频并插入停顿...")
	finalAudio := AssembleWithBreaks(audioSegments, result.Instructions)

	fmt.Printf("✓ 最终音频时长: %v\n", finalAudio.Duration)

	// 4. 保存文件
	fmt.Println("\n💾 步骤 4: 保存文件...")
	filename := "concurrent_simple_output.wav"
	err = SaveSimpleAudio(finalAudio, filename)
	if err != nil {
		log.Fatalf("保存失败: %v", err)
	}

	fmt.Printf("✅ 完成！输出文件: %s\n", filename)
	fmt.Printf("📊 性能: 并发处理 %d 个片段，总耗时 %v\n", len(result.Segments), ttsTime)
}

// ConcurrentGenerateAudio 并发生成音频
func ConcurrentGenerateAudio(segments []ssml.AudioSegment) []SegmentWithAudio {
	fmt.Printf("🔄 并发处理 %d 个片段...\n", len(segments))

	results := make(chan SegmentWithAudio, len(segments))
	var wg sync.WaitGroup

	// TTS 适配器
	tts := ssml.NewMockTTSAdapter(44100)

	// 并发处理每个片段
	for i, segment := range segments {
		wg.Add(1)
		go func(idx int, seg ssml.AudioSegment) {
			defer wg.Done()

			fmt.Printf("  📝 [%d] 处理: \"%s\"\n", idx+1, seg.Text)

			// 模拟处理时间
			time.Sleep(time.Duration(len(seg.Text)*5) * time.Millisecond)

			// 生成音频
			audio, err := tts.GenerateAudio(seg.Text, seg.Properties)
			if err != nil {
				fmt.Printf("  ❌ [%d] 失败: %v\n", idx+1, err)
			} else {
				fmt.Printf("  ✅ [%d] 完成: %v\n", idx+1, audio.Duration)
			}

			results <- SegmentWithAudio{
				Index:     idx,
				Text:      seg.Text,
				AudioData: audio,
				Error:     err,
			}
		}(i, segment)
	}

	// 等待完成
	go func() {
		wg.Wait()
		close(results)
	}()

	// 收集结果
	var segmentResults []SegmentWithAudio
	for result := range results {
		segmentResults = append(segmentResults, result)
	}

	// 按索引排序
	sort.Slice(segmentResults, func(i, j int) bool {
		return segmentResults[i].Index < segmentResults[j].Index
	})

	return segmentResults
}

// AssembleWithBreaks 组装音频并插入停顿
func AssembleWithBreaks(audioSegments []SegmentWithAudio, instructions []ssml.AudioInstruction) *ssml.AudioData {
	fmt.Println("🔧 组装音频...")

	if len(audioSegments) == 0 {
		return &ssml.AudioData{
			SampleRate: 44100,
			Channels:   1,
			Duration:   0,
			Data:       []float64{},
		}
	}

	sampleRate := audioSegments[0].AudioData.SampleRate
	var finalData []float64

	// 提取停顿时长
	var breaks []time.Duration
	for _, instruction := range instructions {
		if instruction.Type == "break" {
			breaks = append(breaks, instruction.Duration)
		}
	}

	// 组装音频：音频片段 + 停顿
	for i, segment := range audioSegments {
		if segment.Error != nil {
			continue
		}

		fmt.Printf("  🎵 [%d] 添加音频: \"%s\" (%v)\n",
			i+1, segment.Text, segment.AudioData.Duration)

		// 添加音频
		finalData = append(finalData, segment.AudioData.Data...)

		// 添加停顿（除了最后一个片段）
		if i < len(breaks) {
			breakDuration := breaks[i]
			fmt.Printf("  ⏸️  [%d] 插入停顿: %v\n", i+1, breakDuration)

			silenceSamples := int(breakDuration.Seconds() * float64(sampleRate))
			silence := make([]float64, silenceSamples)
			finalData = append(finalData, silence...)
		}
	}

	totalDuration := time.Duration(len(finalData)) * time.Second / time.Duration(sampleRate)
	fmt.Printf("✓ 组装完成，总时长: %v，样本数: %d\n", totalDuration, len(finalData))

	return &ssml.AudioData{
		SampleRate: sampleRate,
		Channels:   1,
		Duration:   totalDuration,
		Data:       finalData,
	}
}

// SaveSimpleAudio 保存音频文件
func SaveSimpleAudio(audio *ssml.AudioData, filename string) error {
	writer, err := wav.NewWAVWriter(filename, audio.SampleRate, audio.Channels)
	if err != nil {
		return err
	}
	defer writer.Close()

	samples := make([]int16, len(audio.Data))
	for i, sample := range audio.Data {
		if sample > 1.0 {
			sample = 1.0
		} else if sample < -1.0 {
			sample = -1.0
		}
		samples[i] = int16(sample * 32767)
	}

	return writer.WriteSamples(samples)
}
