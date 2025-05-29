package main

import (
	"fmt"
	"log"
	"time"

	"ssml-parser/ssml"
)

func main() {
	fmt.Println("SSML 音频处理演示")
	fmt.Println("=================")

	// 1. 基本音频处理演示
	demonstrateBasicAudioProcessing()

	// 2. 复杂 SSML 音频处理演示
	demonstrateComplexAudioProcessing()

	// 3. 完整 SSML 到音频流程演示
	demonstrateCompleteSSMLToAudio()
}

// 基本音频处理演示
func demonstrateBasicAudioProcessing() {
	fmt.Println("\n1. 基本音频处理演示")
	fmt.Println("-------------------")

	// 简单的 SSML 内容
	ssmlContent := `<?xml version="1.0" encoding="UTF-8"?>
<speak version="1.0" xml:lang="zh-CN">
    欢迎使用语音合成。
    <break time="1s"/>
    这里是一个停顿。
    <break time="500ms"/>
    短暂停顿后继续。
</speak>`

	// 解析 SSML
	parser := ssml.NewParser(nil)
	parseResult, err := parser.Parse(ssmlContent)
	if err != nil {
		log.Fatalf("SSML 解析失败: %v", err)
	}

	// 音频处理
	audioProcessor := ssml.NewAudioProcessor()
	result, err := audioProcessor.ProcessSSML(parseResult.Root)
	if err != nil {
		log.Fatalf("音频处理失败: %v", err)
	}

	// 显示处理结果
	fmt.Printf("提取的文本: %s\n", result.GetTextForTTS())
	fmt.Printf("预计总时长: %v\n", result.TotalDuration)
	fmt.Printf("音频片段数: %d\n", len(result.Segments))
	fmt.Printf("处理指令数: %d\n", len(result.Instructions))

	// 显示停顿指令
	breaks := result.GetBreakInstructions()
	fmt.Printf("\n停顿指令:\n")
	for i, br := range breaks {
		fmt.Printf("  %d. 位置 %d: %v\n", i+1, br.Position, br.Duration)
	}
}

// 复杂 SSML 音频处理演示
func demonstrateComplexAudioProcessing() {
	fmt.Println("\n2. 复杂 SSML 音频处理演示")
	fmt.Println("-------------------------")

	// 复杂的 SSML 内容
	ssmlContent := `<?xml version="1.0" encoding="UTF-8"?>
<speak version="1.0" xml:lang="zh-CN">
    <voice name="xiaoxiao" gender="female">
        <prosody rate="fast" pitch="high" volume="loud">
            这是快速、高音调、大音量的语音。
            <emphasis level="strong">这里是重点强调</emphasis>
        </prosody>
        <break time="800ms"/>
        <sub alias="人工智能">AI</sub>技术发展迅速。
    </voice>
    <break time="1s"/>
    <audio src="notification.wav">通知音效</audio>
    <p>
        <s>这是第一个句子。</s>
        <s>这是第二个句子。</s>
    </p>
</speak>`

	// 解析和处理
	parser := ssml.NewParser(nil)
	parseResult, err := parser.Parse(ssmlContent)
	if err != nil {
		log.Fatalf("SSML 解析失败: %v", err)
	}

	audioProcessor := ssml.NewAudioProcessor()
	result, err := audioProcessor.ProcessSSML(parseResult.Root)
	if err != nil {
		log.Fatalf("音频处理失败: %v", err)
	}

	// 显示详细的处理报告
	fmt.Println(result.FormatProcessingReport())
}

// 完整 SSML 到音频流程演示
func demonstrateCompleteSSMLToAudio() {
	fmt.Println("\n3. 完整 SSML 到音频流程演示")
	fmt.Println("-----------------------------")

	// 新闻播报的 SSML 内容
	ssmlContent := `<?xml version="1.0" encoding="UTF-8"?>
<speak version="1.0" xml:lang="zh-CN">
    <p>欢迎收听今日新闻。</p>
    <break time="500ms"/>
    
    <voice name="news" gender="male">
        <prosody rate="medium" pitch="low">
            <s>今日头条：科技新闻。</s>
            <break time="300ms"/>
            <s><sub alias="人工智能">AI</sub>技术取得重大突破。</s>
            <break time="300ms"/>
            <s><emphasis level="moderate">这项技术</emphasis>将改变我们的生活。</s>
        </prosody>
    </voice>
    
    <break time="1s"/>
    <audio src="transition.wav">过渡音效</audio>
    
    <voice name="weather" gender="female">
        <prosody rate="fast" pitch="medium">
            <p>天气预报：</p>
            <s>今天晴朗，温度适宜。</s>
            <break time="200ms"/>
            <s>适合外出活动。</s>
        </prosody>
    </voice>
    
    <break time="500ms"/>
    <p>感谢收听，再见！</p>
</speak>`

	// 创建完整的处理器（使用模拟 TTS）
	mockTTS := ssml.NewMockTTSAdapter(44100) // 44.1kHz 采样率
	processor := ssml.NewCompleteSSMLToAudioProcessor(mockTTS, 44100)

	// 处理 SSML 到音频
	fmt.Println("正在处理 SSML...")
	audioData, audioResult, err := processor.ProcessSSMLToAudio(ssmlContent)
	if err != nil {
		log.Fatalf("处理失败: %v", err)
	}

	// 显示结果
	fmt.Printf("处理完成！\n")
	fmt.Printf("最终音频信息:\n")
	fmt.Printf("  - 采样率: %d Hz\n", audioData.SampleRate)
	fmt.Printf("  - 声道数: %d\n", audioData.Channels)
	fmt.Printf("  - 持续时间: %v\n", audioData.Duration)
	fmt.Printf("  - 音频数据点数: %d\n", len(audioData.Data))

	fmt.Printf("\n文本处理结果:\n")
	fmt.Printf("  - 纯文本: %s\n", audioResult.GetTextForTTS())
	fmt.Printf("  - 音频片段数: %d\n", len(audioResult.Segments))
	fmt.Printf("  - 处理指令数: %d\n", len(audioResult.Instructions))

	// 显示主要处理指令
	fmt.Printf("\n主要处理指令:\n")
	for i, instruction := range audioResult.Instructions {
		switch instruction.Type {
		case "break":
			fmt.Printf("  %d. 停顿 %v (位置: %d)\n", i+1, instruction.Duration, instruction.Position)
		case "audio":
			fmt.Printf("  %d. 插入音频: %s (位置: %d)\n", i+1, instruction.AudioFile, instruction.Position)
		case "auto_break":
			fmt.Printf("  %d. 自动停顿 %v (位置: %d)\n", i+1, instruction.Duration, instruction.Position)
		}
	}

	// 演示音频后处理功能
	demonstrateAudioPostProcessing(audioData)
}

// 音频后处理功能演示
func demonstrateAudioPostProcessing(originalAudio *ssml.AudioData) {
	fmt.Println("\n4. 音频后处理功能演示")
	fmt.Println("---------------------")

	postProcessor := ssml.NewAudioPostProcessor(originalAudio.SampleRate)

	// 1. 插入静音
	fmt.Println("演示静音插入:")
	audioWithSilence := postProcessor.InsertSilence(originalAudio, 2*time.Second, 1*time.Second)
	fmt.Printf("  原始音频: %v\n", originalAudio.Duration)
	fmt.Printf("  插入1秒静音后: %v\n", audioWithSilence.Duration)

	// 2. 改变速度
	fmt.Println("\n演示速度调整:")
	fasterAudio := postProcessor.ChangeSpeed(originalAudio, 1.5)  // 快1.5倍
	slowerAudio := postProcessor.ChangeSpeed(originalAudio, 0.75) // 慢0.75倍
	fmt.Printf("  原始音频: %v\n", originalAudio.Duration)
	fmt.Printf("  1.5倍速: %v\n", fasterAudio.Duration)
	fmt.Printf("  0.75倍速: %v\n", slowerAudio.Duration)

	// 3. 应用强调效果
	fmt.Println("\n演示强调效果:")
	emphasizedAudio := postProcessor.ApplyEmphasis(originalAudio, 1*time.Second, 2*time.Second, "strong")
	fmt.Printf("  在1-3秒位置应用强调效果\n")
	fmt.Printf("  处理后音频时长: %v\n", emphasizedAudio.Duration)

	fmt.Println("\n音频处理功能演示完成！")
	fmt.Println("实际应用中，您可以:")
	fmt.Println("  1. 集成真实的 TTS 引擎")
	fmt.Println("  2. 加载和插入真实的音频文件")
	fmt.Println("  3. 应用更复杂的音频效果")
	fmt.Println("  4. 输出到音频文件（WAV、MP3等）")
}

// 使用示例函数
func ExampleUsage() {
	// 这是一个示例，展示如何在实际项目中使用

	// 1. 创建 TTS 适配器（这里需要您实现真实的 TTS 接口）
	// ttsAdapter := &YourRealTTSAdapter{...}
	mockTTS := ssml.NewMockTTSAdapter(22050)

	// 2. 创建完整处理器
	processor := ssml.NewCompleteSSMLToAudioProcessor(mockTTS, 22050)

	// 3. 处理 SSML
	ssml := `<speak version="1.0" xml:lang="zh-CN">
		<voice name="narrator">
			欢迎使用我们的服务。
			<break time="500ms"/>
			请选择您需要的功能。
		</voice>
	</speak>`

	audioData, _, err := processor.ProcessSSMLToAudio(ssml)
	if err != nil {
		log.Fatal(err)
	}

	// 4. 使用生成的音频数据
	_ = audioData // 这里可以保存到文件或播放
}
