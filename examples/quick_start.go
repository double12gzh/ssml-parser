package main

import (
	"fmt"
	"log"
	"ssml-parser/ssml"
)

func main() {
	fmt.Println("🎵 SSML 快速开始示例")
	fmt.Println("==================")

	// 1. 简单的 SSML 内容
	ssmlContent := `<?xml version="1.0" encoding="UTF-8"?>
<speak version="1.0" xml:lang="zh-CN">
    你好，世界！
    <break time="1s"/>
    这是一个简单的 SSML 示例。
    <emphasis level="strong">重点内容！</emphasis>
</speak>`

	fmt.Println("📝 SSML 内容:")
	fmt.Println(ssmlContent)
	fmt.Println()

	// 2. 一行代码完成 SSML 到音频转换
	fmt.Println("🔄 正在处理...")
	audioData, audioResult, err := ProcessSSMLToAudio(ssmlContent)
	if err != nil {
		log.Fatalf("处理失败: %v", err)
	}

	// 3. 显示结果
	fmt.Printf("✅ 处理完成！\n")
	fmt.Printf("📝 提取文本: %s\n", audioResult.GetTextForTTS())
	fmt.Printf("⏱️  音频时长: %v\n", audioData.Duration)
	fmt.Printf("🎵 采样率: %d Hz\n", audioData.SampleRate)
	fmt.Printf("📊 数据点数: %d\n", len(audioData.Data))

	// 4. 保存音频文件
	filename := "quick_start_output.wav"
	err = SaveAudioToFile(audioData, filename)
	if err != nil {
		log.Fatalf("保存文件失败: %v", err)
	}

	fmt.Printf("💾 音频已保存: %s\n", filename)
	fmt.Println("🎧 可以用音频播放器播放该文件")
}

// ProcessSSMLToAudio 一行代码处理 SSML 到音频的完整流程
func ProcessSSMLToAudio(ssmlContent string) (*ssml.AudioData, *ssml.AudioProcessingResult, error) {
	// 创建 TTS 适配器
	ttsAdapter := ssml.NewMockTTSAdapter(44100)

	// 创建完整处理器
	processor := ssml.NewCompleteSSMLToAudioProcessor(ttsAdapter, 44100)

	// 处理 SSML
	return processor.ProcessSSMLToAudio(ssmlContent)
}

// SaveAudioToFile 保存音频到文件
func SaveAudioToFile(audio *ssml.AudioData, filename string) error {
	// 创建 WAV 写入器
	writer, err := NewWAVWriter(filename, audio.SampleRate, audio.Channels)
	if err != nil {
		return err
	}
	defer writer.Close()

	// 转换并写入数据
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
