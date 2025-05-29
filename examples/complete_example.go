package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"time"

	"ssml-parser/ssml"
)

func main() {
	fmt.Println("SSML 到音频文件完整示例")
	fmt.Println("=========================")

	// 步骤 1: 准备 SSML 内容
	ssmlContent := `<?xml version="1.0" encoding="UTF-8"?>
<speak version="1.0" xml:lang="zh-CN">
    <voice name="xiaoxiao" gender="female">
        欢迎使用 SSML 音频处理系统！
        <break time="1s"/>
        
        <prosody rate="fast" pitch="high" volume="loud">
            这是快速、高音调、大音量的语音。
        </prosody>
        
        <break time="500ms"/>
        
        <emphasis level="strong">这里是重点强调的内容。</emphasis>
        
        <break time="800ms"/>
        
        <sub alias="人工智能">AI</sub>技术正在快速发展。
        
        <break time="200ms"/>
        
        感谢您的使用！
    </voice>
</speak>`

	fmt.Println("原始 SSML 内容:")
	fmt.Println(ssmlContent)
	fmt.Println()

	// 步骤 2: 解析 SSML
	fmt.Println("步骤 1: 解析 SSML...")
	parser := ssml.NewParser(nil)
	parseResult, err := parser.Parse(ssmlContent)
	if err != nil {
		log.Fatalf("SSML 解析失败: %v", err)
	}
	fmt.Printf("✓ SSML 解析成功，版本: %s, 语言: %s\n",
		parseResult.Root.Version, parseResult.Root.Lang)

	// 步骤 3: 提取文本和处理指令
	fmt.Println("\n步骤 2: 提取文本和音频处理指令...")
	audioProcessor := ssml.NewAudioProcessor()
	result, err := audioProcessor.ProcessSSML(parseResult.Root)
	if err != nil {
		log.Fatalf("音频处理失败: %v", err)
	}

	fmt.Printf("✓ 提取的纯文本: %s\n", result.GetTextForTTS())
	fmt.Printf("✓ 音频片段数: %d\n", len(result.Segments))
	fmt.Printf("✓ 处理指令数: %d\n", len(result.Instructions))

	// 显示提取的音频片段
	fmt.Println("\n音频片段详情:")
	for i, segment := range result.Segments {
		fmt.Printf("  %d. [%v-%v] \"%s\" (声音=%s, 语速=%s, 音调=%s, 音量=%s)\n",
			i+1,
			segment.StartTime,
			segment.StartTime+segment.Duration,
			segment.Text,
			segment.Properties.Voice,
			segment.Properties.Rate,
			segment.Properties.Pitch,
			segment.Properties.Volume)
	}

	// 显示处理指令
	fmt.Println("\n处理指令详情:")
	for i, instruction := range result.Instructions {
		switch instruction.Type {
		case "break":
			fmt.Printf("  %d. 停顿 %v (位置: %d)\n", i+1, instruction.Duration, instruction.Position)
		case "auto_break":
			fmt.Printf("  %d. 自动停顿 %v (位置: %d)\n", i+1, instruction.Duration, instruction.Position)
		case "audio":
			fmt.Printf("  %d. 插入音频: %s (位置: %d)\n", i+1, instruction.AudioFile, instruction.Position)
		}
	}

	// 步骤 4: 使用增强的 Mock TTS 生成基础音频
	fmt.Println("\n步骤 3: 使用 TTS 生成基础音频...")
	mockTTS := NewEnhancedMockTTS(44100)

	// 使用第一个音频片段的属性作为默认属性
	var properties *ssml.AudioProperties
	if len(result.Segments) > 0 {
		properties = result.Segments[0].Properties
	} else {
		properties = &ssml.AudioProperties{
			Rate:     "medium",
			Pitch:    "medium",
			Volume:   "medium",
			Voice:    "default",
			Gender:   "neutral",
			Language: "zh-CN",
			Emphasis: "none",
		}
	}

	ttsAudio, err := mockTTS.GenerateAudio(result.GetTextForTTS(), properties)
	if err != nil {
		log.Fatalf("TTS 生成失败: %v", err)
	}

	fmt.Printf("✓ TTS 音频生成成功\n")
	fmt.Printf("  - 采样率: %d Hz\n", ttsAudio.SampleRate)
	fmt.Printf("  - 声道数: %d\n", ttsAudio.Channels)
	fmt.Printf("  - 原始时长: %v\n", ttsAudio.Duration)
	fmt.Printf("  - 数据点数: %d\n", len(ttsAudio.Data))

	// 步骤 5: 应用音频后处理
	fmt.Println("\n步骤 4: 应用音频后处理...")
	postProcessor := ssml.NewAudioPostProcessor(ttsAudio.SampleRate)
	finalAudio, err := postProcessor.ProcessTTSAudio(ttsAudio, result.Instructions)
	if err != nil {
		log.Fatalf("音频后处理失败: %v", err)
	}

	fmt.Printf("✓ 音频后处理完成\n")
	fmt.Printf("  - 最终时长: %v\n", finalAudio.Duration)
	fmt.Printf("  - 最终数据点数: %d\n", len(finalAudio.Data))

	// 显示处理效果对比
	fmt.Printf("  - 时长变化: %v → %v (增加了 %v)\n",
		ttsAudio.Duration,
		finalAudio.Duration,
		finalAudio.Duration-ttsAudio.Duration)

	// 步骤 6: 保存为 WAV 文件
	fmt.Println("\n步骤 5: 保存音频文件...")
	filename := "output_ssml_audio.wav"
	err = SaveToWAV(finalAudio, filename)
	if err != nil {
		log.Fatalf("保存音频文件失败: %v", err)
	}

	fmt.Printf("✓ 音频文件保存成功: %s\n", filename)

	// 步骤 7: 验证 WAV 文件
	fmt.Println("\n步骤 6: 验证音频文件...")
	err = ValidateWAVFile(filename)
	if err != nil {
		log.Printf("WAV 文件验证失败: %v", err)
	}

	// 步骤 8: 显示文件信息
	fmt.Println("\n步骤 7: 显示文件信息...")
	ShowFileInfo(filename)

	// 步骤 9: 显示处理报告
	fmt.Println("\n=== 完整处理报告 ===")
	fmt.Println(result.FormatProcessingReport())

	fmt.Println("=== 处理完成! ===")
	fmt.Printf("输出文件: %s\n", filename)
	fmt.Println("您可以使用音频播放软件播放生成的 WAV 文件。")
	fmt.Println()
	fmt.Println("## 使用说明")
	fmt.Println("1. 生成的是 16 位 PCM WAV 文件")
	fmt.Println("2. 可以用 Windows Media Player, VLC, Audacity 等播放")
	fmt.Println("3. 文件包含了 SSML 中定义的停顿效果")
	fmt.Println("4. 这是模拟音频，实际应用中可替换为真实 TTS 引擎")
}

// EnhancedMockTTS 增强的模拟 TTS，生成更真实的音频
type EnhancedMockTTS struct {
	sampleRate int
}

// NewEnhancedMockTTS 创建增强的模拟 TTS
func NewEnhancedMockTTS(sampleRate int) *EnhancedMockTTS {
	return &EnhancedMockTTS{sampleRate: sampleRate}
}

// GenerateAudio 生成更真实的模拟音频
func (tts *EnhancedMockTTS) GenerateAudio(text string, properties *ssml.AudioProperties) (*ssml.AudioData, error) {
	// 根据文本长度计算时长（每字符 200ms）
	textLength := len([]rune(text))
	durationMs := textLength * 200

	sampleCount := (durationMs * tts.sampleRate) / 1000
	data := make([]float64, sampleCount)

	// 生成更复杂的音频波形模拟语音
	baseFreq := 200.0 // 基础频率

	// 根据音调调整频率
	switch properties.Pitch {
	case "x-low":
		baseFreq *= 0.7
	case "low":
		baseFreq *= 0.85
	case "medium":
		baseFreq *= 1.0
	case "high":
		baseFreq *= 1.2
	case "x-high":
		baseFreq *= 1.4
	}

	// 根据音量调整幅度
	amplitude := 0.3
	switch properties.Volume {
	case "silent":
		amplitude = 0.0
	case "x-soft":
		amplitude = 0.1
	case "soft":
		amplitude = 0.2
	case "medium":
		amplitude = 0.3
	case "loud":
		amplitude = 0.5
	case "x-loud":
		amplitude = 0.7
	}

	// 生成模拟的语音波形
	for i := range data {
		t := float64(i) / float64(tts.sampleRate)

		// 主要频率分量
		wave1 := math.Sin(2 * math.Pi * baseFreq * t)

		// 添加谐波分量使声音更丰富
		wave2 := 0.3 * math.Sin(2*math.Pi*baseFreq*2*t)
		wave3 := 0.2 * math.Sin(2*math.Pi*baseFreq*3*t)

		// 添加一些随机噪声模拟语音的自然性
		noise := 0.05 * (2*math.Sin(float64(i)*0.01) - 1)

		// 添加幅度调制模拟语音的节奏感
		envelope := 0.5 + 0.5*math.Sin(2*math.Pi*t*2) // 2Hz 的调制

		// 组合所有分量
		sample := amplitude * envelope * (wave1 + wave2 + wave3 + noise)

		// 限制幅度范围
		if sample > 1.0 {
			sample = 1.0
		} else if sample < -1.0 {
			sample = -1.0
		}

		data[i] = sample
	}

	return &ssml.AudioData{
		SampleRate: tts.sampleRate,
		Channels:   1,
		Duration:   time.Duration(durationMs) * time.Millisecond,
		Data:       data,
	}, nil
}

// SaveToWAV 将音频数据保存为 WAV 文件
func SaveToWAV(audio *ssml.AudioData, filename string) error {
	fmt.Printf("正在保存 WAV 文件: %s...\n", filename)

	// 创建 WAV 写入器
	writer, err := NewWAVWriter(filename, audio.SampleRate, audio.Channels)
	if err != nil {
		return fmt.Errorf("创建 WAV 写入器失败: %w", err)
	}
	defer writer.Close()

	// 将浮点数音频数据转换为 16 位整数
	samples := make([]int16, len(audio.Data))
	for i, sample := range audio.Data {
		// 将 [-1.0, 1.0] 范围的浮点数转换为 [-32768, 32767] 范围的整数
		if sample > 1.0 {
			sample = 1.0
		} else if sample < -1.0 {
			sample = -1.0
		}
		samples[i] = int16(sample * 32767)
	}

	// 写入音频数据
	err = writer.WriteSamples(samples)
	if err != nil {
		return fmt.Errorf("写入音频数据失败: %w", err)
	}

	fmt.Printf("✓ WAV 文件保存成功，包含 %d 个样本\n", len(samples))
	return nil
}

// ShowFileInfo 显示文件信息
func ShowFileInfo(filename string) {
	info, err := os.Stat(filename)
	if err != nil {
		fmt.Printf("获取文件信息失败: %v\n", err)
		return
	}

	fmt.Printf("文件信息:\n")
	fmt.Printf("  - 文件名: %s\n", info.Name())
	fmt.Printf("  - 文件大小: %d 字节 (%.2f KB)\n", info.Size(), float64(info.Size())/1024)
	fmt.Printf("  - 修改时间: %s\n", info.ModTime().Format("2006-01-02 15:04:05"))
}
