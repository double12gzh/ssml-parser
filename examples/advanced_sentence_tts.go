package main

import (
	"fmt"
	"log"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"ssml-parser/ssml"
)

func main() {
	fmt.Println("🚀 高级切句并发 TTS 示例")
	fmt.Println("=========================")

	// 包含长文本和多个标签的复杂 SSML
	ssmlContent := `<?xml version="1.0" encoding="UTF-8"?>
<speak version="1.0" xml:lang="zh-CN">
    欢迎使用我们的智能语音系统。这是一个非常先进的系统，可以处理复杂的文本。
    <break time="1s"/>
    
    <prosody rate="fast">
        我们的系统支持多种语音效果。您可以调整语速、音调和音量。这些功能都很实用。
    </prosody>
    <break time="500ms"/>
    
    <emphasis level="strong">
        重要提醒：请仔细阅读使用说明。正确使用可以获得最佳体验！
    </emphasis>
    <break time="800ms"/>
    
    感谢您的使用。如果您有任何问题，请联系我们的客服团队。
</speak>`

	fmt.Println("📝 原始 SSML:")
	fmt.Println(ssmlContent)
	fmt.Println()

	DemoAdvancedSentenceTTS(ssmlContent)
}

// SentenceTTSResult 切句TTS结果
type SentenceTTSResult struct {
	Index       int
	Text        string
	AudioData   *ssml.AudioData
	Error       error
	ProcessTime time.Duration
}

// TagInsertPoint 标签插入点
type TagInsertPoint struct {
	Type          string
	AfterSentence int
	Duration      time.Duration
}

// DemoAdvancedSentenceTTS 演示高级切句TTS处理
func DemoAdvancedSentenceTTS(ssmlContent string) {
	// 1. 解析SSML
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

	fullText := result.GetTextForTTS()
	fmt.Printf("✓ 提取文本: %s\n", fullText)

	// 2. 切句服务
	fmt.Println("\n📝 步骤 2: 调用切句服务...")
	sentences := SmartSentenceSegmentation(fullText)
	fmt.Printf("✓ 切分为 %d 个句子:\n", len(sentences))
	for i, sentence := range sentences {
		fmt.Printf("  %d. \"%s\"\n", i+1, sentence)
	}

	// 3. 计算标签位置
	fmt.Println("\n📍 步骤 3: 计算标签插入位置...")
	insertPoints := CalculateBreakPositions(fullText, sentences, result.Instructions)
	fmt.Printf("✓ 计算了 %d 个插入点:\n", len(insertPoints))
	for i, point := range insertPoints {
		fmt.Printf("  %d. %s 在句子 %d 后，停顿 %v\n",
			i+1, point.Type, point.AfterSentence+1, point.Duration)
	}

	// 4. 并发TTS
	fmt.Println("\n🎵 步骤 4: 并发 TTS 生成...")
	start := time.Now()
	ttsResults := ParallelTTSGeneration(sentences)
	ttsTime := time.Since(start)

	fmt.Printf("✓ 并发处理完成，耗时: %v\n", ttsTime)

	// 5. 组装音频
	fmt.Println("\n🔧 步骤 5: 组装音频...")
	finalAudio := AssembleAudioWithBreaks(ttsResults, insertPoints)

	fmt.Printf("✓ 最终音频时长: %v\n", finalAudio.Duration)

	// 6. 保存文件
	fmt.Println("\n💾 步骤 6: 保存文件...")
	filename := "advanced_sentence_output.wav"
	err = SaveSentenceAudio(finalAudio, filename)
	if err != nil {
		log.Fatalf("保存失败: %v", err)
	}

	fmt.Printf("✅ 完成！输出文件: %s\n", filename)
	fmt.Printf("📊 统计: %d 句子，%v 处理时间\n", len(sentences), ttsTime)
}

// SmartSentenceSegmentation 智能切句
func SmartSentenceSegmentation(text string) []string {
	fmt.Println("🔪 智能切句处理...")

	// 模拟切句服务延迟
	time.Sleep(30 * time.Millisecond)

	// 多级切句策略
	var sentences []string

	// 1. 按句号、感叹号、问号切分
	mainParts := regexp.MustCompile(`[。！？]+`).Split(text, -1)

	for _, part := range mainParts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		// 2. 如果片段过长，按逗号切分
		if len([]rune(part)) > 30 {
			subParts := regexp.MustCompile(`[，；]+`).Split(part, -1)
			for _, subPart := range subParts {
				subPart = strings.TrimSpace(subPart)
				if subPart != "" {
					sentences = append(sentences, subPart)
				}
			}
		} else {
			sentences = append(sentences, part)
		}
	}

	fmt.Printf("  ✓ 切句结果: %d 句\n", len(sentences))
	return sentences
}

// CalculateBreakPositions 计算break标签的精确插入位置
func CalculateBreakPositions(fullText string, sentences []string, instructions []ssml.AudioInstruction) []TagInsertPoint {
	var points []TagInsertPoint

	// 建立句子位置映射
	sentencePositions := make([]int, len(sentences))
	currentPos := 0

	for i, sentence := range sentences {
		pos := strings.Index(fullText[currentPos:], sentence)
		if pos != -1 {
			sentencePositions[i] = currentPos + pos
			currentPos = currentPos + pos + len(sentence)
		}
	}

	// 为每个break指令找到对应的插入位置
	for _, instruction := range instructions {
		if instruction.Type != "break" {
			continue
		}

		// 找到break应该插入到哪个句子之后
		afterSentence := -1
		for i, sentPos := range sentencePositions {
			sentEnd := sentPos + len(sentences[i])
			if instruction.Position <= sentEnd {
				afterSentence = i
				break
			}
		}

		if afterSentence == -1 {
			afterSentence = len(sentences) - 1
		}

		point := TagInsertPoint{
			Type:          instruction.Type,
			AfterSentence: afterSentence,
			Duration:      instruction.Duration,
		}

		points = append(points, point)
	}

	return points
}

// ParallelTTSGeneration 并行TTS生成
func ParallelTTSGeneration(sentences []string) []SentenceTTSResult {
	fmt.Printf("🔄 并行处理 %d 个句子...\n", len(sentences))

	results := make(chan SentenceTTSResult, len(sentences))
	var wg sync.WaitGroup

	tts := ssml.NewMockTTSAdapter(44100)

	for i, sentence := range sentences {
		wg.Add(1)
		go func(idx int, text string) {
			defer wg.Done()

			fmt.Printf("  📝 [%d] 处理: \"%s\"\n", idx+1, text)

			start := time.Now()
			time.Sleep(time.Duration(len(text)*6) * time.Millisecond)

			// 使用默认属性生成音频
			properties := &ssml.AudioProperties{
				Rate:     "medium",
				Pitch:    "medium",
				Volume:   "medium",
				Voice:    "default",
				Gender:   "neutral",
				Language: "zh-CN",
				Emphasis: "none",
			}

			audio, err := tts.GenerateAudio(text, properties)
			duration := time.Since(start)

			if err != nil {
				fmt.Printf("  ❌ [%d] 失败: %v\n", idx+1, err)
			} else {
				fmt.Printf("  ✅ [%d] 完成: %v (耗时: %v)\n", idx+1, audio.Duration, duration)
			}

			results <- SentenceTTSResult{
				Index:       idx,
				Text:        text,
				AudioData:   audio,
				Error:       err,
				ProcessTime: duration,
			}
		}(i, sentence)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	var ttsResults []SentenceTTSResult
	for result := range results {
		ttsResults = append(ttsResults, result)
	}

	sort.Slice(ttsResults, func(i, j int) bool {
		return ttsResults[i].Index < ttsResults[j].Index
	})

	return ttsResults
}

// AssembleAudioWithBreaks 组装音频并插入break
func AssembleAudioWithBreaks(ttsResults []SentenceTTSResult, insertPoints []TagInsertPoint) *ssml.AudioData {
	fmt.Println("🔧 组装音频并插入停顿...")

	if len(ttsResults) == 0 {
		return &ssml.AudioData{
			SampleRate: 44100,
			Channels:   1,
			Duration:   0,
			Data:       []float64{},
		}
	}

	sampleRate := ttsResults[0].AudioData.SampleRate
	var finalData []float64

	for i, result := range ttsResults {
		if result.Error != nil {
			continue
		}

		fmt.Printf("  🎵 [%d] 添加: \"%s\" (%v)\n",
			i+1, result.Text, result.AudioData.Duration)

		finalData = append(finalData, result.AudioData.Data...)

		// 检查是否需要在此句子后插入停顿
		for _, point := range insertPoints {
			if point.AfterSentence == result.Index {
				fmt.Printf("  ⏸️  [%d] 插入停顿: %v\n", i+1, point.Duration)

				silenceSamples := int(point.Duration.Seconds() * float64(sampleRate))
				silence := make([]float64, silenceSamples)
				finalData = append(finalData, silence...)
			}
		}
	}

	totalDuration := time.Duration(len(finalData)) * time.Second / time.Duration(sampleRate)
	fmt.Printf("✓ 组装完成，总时长: %v，样本: %d\n", totalDuration, len(finalData))

	return &ssml.AudioData{
		SampleRate: sampleRate,
		Channels:   1,
		Duration:   totalDuration,
		Data:       finalData,
	}
}

// SaveSentenceAudio 保存音频文件
func SaveSentenceAudio(audio *ssml.AudioData, filename string) error {
	writer, err := NewWAVWriter(filename, audio.SampleRate, audio.Channels)
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
