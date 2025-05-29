package ssml

import (
	"fmt"
	"time"
)

// AudioData 表示音频数据
type AudioData struct {
	SampleRate int           // 采样率
	Channels   int           // 声道数
	Duration   time.Duration // 持续时间
	Data       []float64     // 音频数据（归一化的浮点数）
}

// AudioPostProcessor 音频后处理器
type AudioPostProcessor struct {
	sampleRate int
}

// NewAudioPostProcessor 创建音频后处理器
func NewAudioPostProcessor(sampleRate int) *AudioPostProcessor {
	return &AudioPostProcessor{
		sampleRate: sampleRate,
	}
}

// ProcessTTSAudio 处理 TTS 生成的音频
func (app *AudioPostProcessor) ProcessTTSAudio(ttsAudio *AudioData, instructions []AudioInstruction) (*AudioData, error) {
	if len(instructions) == 0 {
		return ttsAudio, nil
	}

	// 创建结果音频数据
	result := &AudioData{
		SampleRate: ttsAudio.SampleRate,
		Channels:   ttsAudio.Channels,
		Data:       make([]float64, 0),
	}

	// 按位置排序指令
	sortedInstructions := app.sortInstructionsByPosition(instructions)

	currentSample := 0

	for _, instruction := range sortedInstructions {
		switch instruction.Type {
		case "break", "auto_break":
			// 复制到停顿位置的音频
			breakPosition := app.timeToSamples(instruction.Duration, ttsAudio.SampleRate)
			if currentSample < len(ttsAudio.Data) {
				endSample := currentSample + breakPosition
				if endSample > len(ttsAudio.Data) {
					endSample = len(ttsAudio.Data)
				}
				result.Data = append(result.Data, ttsAudio.Data[currentSample:endSample]...)
				currentSample = endSample
			}

			// 插入静音
			silenceSamples := app.timeToSamples(instruction.Duration, ttsAudio.SampleRate)
			silence := make([]float64, silenceSamples)
			result.Data = append(result.Data, silence...)

		case "audio":
			// 这里可以插入外部音频文件
			// 目前只是添加一个占位符
			fmt.Printf("需要插入音频文件: %s\n", instruction.AudioFile)
		}
	}

	// 添加剩余的音频数据
	if currentSample < len(ttsAudio.Data) {
		result.Data = append(result.Data, ttsAudio.Data[currentSample:]...)
	}

	result.Duration = time.Duration(len(result.Data)) * time.Second / time.Duration(result.SampleRate)

	return result, nil
}

// InsertSilence 在指定位置插入静音
func (app *AudioPostProcessor) InsertSilence(audio *AudioData, position time.Duration, duration time.Duration) *AudioData {
	positionSamples := app.timeToSamples(position, audio.SampleRate)
	durationSamples := app.timeToSamples(duration, audio.SampleRate)

	// 创建新的音频数据
	newData := make([]float64, len(audio.Data)+durationSamples)

	// 复制插入点之前的数据
	copy(newData[:positionSamples], audio.Data[:positionSamples])

	// 插入静音（零值）
	for i := positionSamples; i < positionSamples+durationSamples; i++ {
		newData[i] = 0.0
	}

	// 复制插入点之后的数据
	copy(newData[positionSamples+durationSamples:], audio.Data[positionSamples:])

	return &AudioData{
		SampleRate: audio.SampleRate,
		Channels:   audio.Channels,
		Duration:   time.Duration(len(newData)) * time.Second / time.Duration(audio.SampleRate),
		Data:       newData,
	}
}

// InsertAudioFile 在指定位置插入音频文件（接口定义）
func (app *AudioPostProcessor) InsertAudioFile(audio *AudioData, position time.Duration, audioFile string) (*AudioData, error) {
	// 这里需要实际的音频文件加载逻辑
	// 目前返回原音频作为占位符
	fmt.Printf("TODO: 插入音频文件 %s 到位置 %v\n", audioFile, position)
	return audio, nil
}

// ApplyEmphasis 应用强调效果（通过调整音量）
func (app *AudioPostProcessor) ApplyEmphasis(audio *AudioData, startTime, duration time.Duration, level string) *AudioData {
	startSample := app.timeToSamples(startTime, audio.SampleRate)
	durationSamples := app.timeToSamples(duration, audio.SampleRate)
	endSample := startSample + durationSamples

	if endSample > len(audio.Data) {
		endSample = len(audio.Data)
	}

	// 根据强调级别调整音量
	var multiplier float64
	switch level {
	case "strong":
		multiplier = 1.3
	case "moderate":
		multiplier = 1.15
	case "reduced":
		multiplier = 0.85
	default:
		multiplier = 1.0
	}

	// 创建新的音频数据
	newData := make([]float64, len(audio.Data))
	copy(newData, audio.Data)

	// 应用强调
	for i := startSample; i < endSample && i < len(newData); i++ {
		newData[i] *= multiplier
		// 确保不超出范围
		if newData[i] > 1.0 {
			newData[i] = 1.0
		} else if newData[i] < -1.0 {
			newData[i] = -1.0
		}
	}

	return &AudioData{
		SampleRate: audio.SampleRate,
		Channels:   audio.Channels,
		Duration:   audio.Duration,
		Data:       newData,
	}
}

// ChangeSpeed 改变音频速度（简单的重采样）
func (app *AudioPostProcessor) ChangeSpeed(audio *AudioData, speedMultiplier float64) *AudioData {
	if speedMultiplier == 1.0 {
		return audio
	}

	newLength := int(float64(len(audio.Data)) / speedMultiplier)
	newData := make([]float64, newLength)

	// 简单的线性插值重采样
	for i := 0; i < newLength; i++ {
		sourceIndex := float64(i) * speedMultiplier
		leftIndex := int(sourceIndex)
		rightIndex := leftIndex + 1

		if rightIndex >= len(audio.Data) {
			rightIndex = len(audio.Data) - 1
		}
		if leftIndex >= len(audio.Data) {
			leftIndex = len(audio.Data) - 1
		}

		// 线性插值
		fraction := sourceIndex - float64(leftIndex)
		newData[i] = audio.Data[leftIndex]*(1-fraction) + audio.Data[rightIndex]*fraction
	}

	return &AudioData{
		SampleRate: audio.SampleRate,
		Channels:   audio.Channels,
		Duration:   time.Duration(newLength) * time.Second / time.Duration(audio.SampleRate),
		Data:       newData,
	}
}

// 辅助方法

// timeToSamples 将时间转换为采样点数
func (app *AudioPostProcessor) timeToSamples(duration time.Duration, sampleRate int) int {
	return int(duration.Seconds() * float64(sampleRate))
}

// samplesToTime 将采样点数转换为时间
func (app *AudioPostProcessor) samplesToTime(samples int, sampleRate int) time.Duration {
	return time.Duration(samples) * time.Second / time.Duration(sampleRate)
}

// sortInstructionsByPosition 按位置排序指令
func (app *AudioPostProcessor) sortInstructionsByPosition(instructions []AudioInstruction) []AudioInstruction {
	result := make([]AudioInstruction, len(instructions))
	copy(result, instructions)

	// 简单的冒泡排序
	for i := 0; i < len(result)-1; i++ {
		for j := 0; j < len(result)-1-i; j++ {
			if result[j].Position > result[j+1].Position {
				result[j], result[j+1] = result[j+1], result[j]
			}
		}
	}

	return result
}

// TTSAdapter TTS 适配器接口
type TTSAdapter interface {
	// GenerateAudio 生成音频
	GenerateAudio(text string, properties *AudioProperties) (*AudioData, error)
}

// MockTTSAdapter 模拟的 TTS 适配器（用于测试）
type MockTTSAdapter struct {
	sampleRate int
}

// NewMockTTSAdapter 创建模拟 TTS 适配器
func NewMockTTSAdapter(sampleRate int) *MockTTSAdapter {
	return &MockTTSAdapter{sampleRate: sampleRate}
}

// GenerateAudio 生成模拟音频数据
func (mock *MockTTSAdapter) GenerateAudio(text string, properties *AudioProperties) (*AudioData, error) {
	// 为每个字符生成 100ms 的模拟音频
	textLength := len([]rune(text))
	durationMs := textLength * 100 // 每字符 100ms

	sampleCount := (durationMs * mock.sampleRate) / 1000
	data := make([]float64, sampleCount)

	// 生成简单的模拟音频模式
	for i := range data {
		data[i] = 0.1 * (float64(i%100) / 100.0) // 简单的音频模式
	}

	return &AudioData{
		SampleRate: mock.sampleRate,
		Channels:   1,
		Duration:   time.Duration(durationMs) * time.Millisecond,
		Data:       data,
	}, nil
}

// CompleteSSMLToAudioProcessor 完整的 SSML 到音频处理器
type CompleteSSMLToAudioProcessor struct {
	ssmlParser     *Parser
	audioProcessor *AudioProcessor
	postProcessor  *AudioPostProcessor
	ttsAdapter     TTSAdapter
}

// NewCompleteSSMLToAudioProcessor 创建完整的处理器
func NewCompleteSSMLToAudioProcessor(ttsAdapter TTSAdapter, sampleRate int) *CompleteSSMLToAudioProcessor {
	return &CompleteSSMLToAudioProcessor{
		ssmlParser:     NewParser(nil),
		audioProcessor: NewAudioProcessor(),
		postProcessor:  NewAudioPostProcessor(sampleRate),
		ttsAdapter:     ttsAdapter,
	}
}

// ProcessSSMLToAudio 将 SSML 转换为处理后的音频
func (processor *CompleteSSMLToAudioProcessor) ProcessSSMLToAudio(ssmlContent string) (*AudioData, *AudioProcessingResult, error) {
	// 1. 解析 SSML
	parseResult, err := processor.ssmlParser.Parse(ssmlContent)
	if err != nil {
		return nil, nil, fmt.Errorf("SSML 解析失败: %w", err)
	}

	// 2. 处理 SSML，提取文本和指令
	audioResult, err := processor.audioProcessor.ProcessSSML(parseResult.Root)
	if err != nil {
		return nil, nil, fmt.Errorf("音频处理失败: %w", err)
	}

	// 检查是否有文本内容
	if audioResult.PlainText == "" {
		// 如果没有文本，创建一个空的音频数据
		return &AudioData{
			SampleRate: 44100,
			Channels:   1,
			Duration:   0,
			Data:       []float64{},
		}, audioResult, nil
	}

	// 3. 使用 TTS 生成基础音频
	var properties *AudioProperties
	if len(audioResult.Segments) > 0 {
		properties = audioResult.Segments[0].Properties
	} else {
		// 使用默认属性
		properties = &AudioProperties{
			Rate:     "medium",
			Pitch:    "medium",
			Volume:   "medium",
			Voice:    "default",
			Gender:   "neutral",
			Language: "zh-CN",
			Emphasis: "none",
		}
	}

	ttsAudio, err := processor.ttsAdapter.GenerateAudio(audioResult.GetTextForTTS(), properties)
	if err != nil {
		return nil, nil, fmt.Errorf("TTS 生成失败: %w", err)
	}

	// 4. 应用后处理（插入静音等）
	finalAudio, err := processor.postProcessor.ProcessTTSAudio(ttsAudio, audioResult.Instructions)
	if err != nil {
		return nil, nil, fmt.Errorf("音频后处理失败: %w", err)
	}

	return finalAudio, audioResult, nil
}
