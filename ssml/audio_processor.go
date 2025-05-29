package ssml

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// AudioSegment 表示一个音频片段
type AudioSegment struct {
	Text       string           // 要合成的文本
	StartTime  time.Duration    // 在最终音频中的开始时间
	Duration   time.Duration    // 预期持续时间
	Properties *AudioProperties // 音频属性
}

// AudioProperties 包含音频处理的属性
type AudioProperties struct {
	Rate     string // 语速：slow, medium, fast 或具体数值
	Pitch    string // 音调：low, medium, high 或具体数值
	Volume   string // 音量：silent, soft, medium, loud 或具体数值
	Voice    string // 声音名称
	Gender   string // 性别：male, female, neutral
	Language string // 语言
	Emphasis string // 强调级别：none, reduced, moderate, strong
}

// AudioInstruction 表示音频处理指令
type AudioInstruction struct {
	Type       string           // 指令类型：break, audio, silence, emphasis
	Position   int              // 在文本中的位置（字符索引）
	Duration   time.Duration    // 持续时间（用于 break）
	AudioFile  string           // 音频文件路径（用于 audio）
	Properties *AudioProperties // 音频属性变更
}

// AudioProcessingResult 包含文本提取和音频处理指令
type AudioProcessingResult struct {
	PlainText     string             // 提取的纯文本
	Segments      []AudioSegment     // 音频片段
	Instructions  []AudioInstruction // 音频处理指令
	TotalDuration time.Duration      // 预计总时长
}

// AudioProcessor 音频处理器
type AudioProcessor struct {
	defaultProperties *AudioProperties
	charToTimeRatio   time.Duration // 每字符的预计发音时间
}

// NewAudioProcessor 创建新的音频处理器
func NewAudioProcessor() *AudioProcessor {
	return &AudioProcessor{
		defaultProperties: &AudioProperties{
			Rate:     "medium",
			Pitch:    "medium",
			Volume:   "medium",
			Voice:    "default",
			Gender:   "neutral",
			Language: "zh-CN",
			Emphasis: "none",
		},
		charToTimeRatio: 150 * time.Millisecond, // 默认每字符150ms
	}
}

// ProcessSSML 处理 SSML 并生成音频处理结果
func (ap *AudioProcessor) ProcessSSML(speak *Speak) (*AudioProcessingResult, error) {
	result := &AudioProcessingResult{
		Segments:     make([]AudioSegment, 0),
		Instructions: make([]AudioInstruction, 0),
	}

	// 创建处理上下文
	ctx := &processingContext{
		processor:        ap,
		result:           result,
		currentTime:      0,
		textPosition:     0,
		propertyStack:    []*AudioProperties{ap.copyProperties(ap.defaultProperties)},
		plainTextBuilder: &strings.Builder{},
	}

	// 处理所有内容
	for _, content := range speak.Content {
		ctx.processContent(content)
	}

	result.PlainText = ctx.plainTextBuilder.String()
	result.TotalDuration = ctx.currentTime

	return result, nil
}

// processingContext 处理上下文
type processingContext struct {
	processor        *AudioProcessor
	result           *AudioProcessingResult
	currentTime      time.Duration
	textPosition     int
	propertyStack    []*AudioProperties
	plainTextBuilder *strings.Builder
}

// processElement 处理单个元素
func (ctx *processingContext) processElement(element SSMLElement) {
	switch elem := element.(type) {
	case *Text:
		ctx.processText(elem)
	case *Break:
		ctx.processBreak(elem)
	case *Audio:
		ctx.processAudio(elem)
	case *Voice:
		ctx.processVoice(elem)
	case *Prosody:
		ctx.processProsody(elem)
	case *Emphasis:
		ctx.processEmphasis(elem)
	case *Phoneme:
		ctx.processPhoneme(elem)
	case *Sub:
		ctx.processSub(elem)
	case *Paragraph, *Sentence, *W:
		ctx.processContainer(elem)
	}
}

// processText 处理文本
func (ctx *processingContext) processText(text *Text) {
	content := strings.TrimSpace(text.Content)
	if content == "" {
		return
	}

	// 添加到纯文本
	ctx.plainTextBuilder.WriteString(content)

	// 计算预期持续时间
	duration := ctx.calculateTextDuration(content)

	// 创建音频片段
	segment := AudioSegment{
		Text:       content,
		StartTime:  ctx.currentTime,
		Duration:   duration,
		Properties: ctx.getCurrentProperties(),
	}

	ctx.result.Segments = append(ctx.result.Segments, segment)
	ctx.currentTime += duration
	ctx.textPosition += len(content)
}

// processBreak 处理停顿
func (ctx *processingContext) processBreak(br *Break) {
	var duration time.Duration

	if br.Time != "" {
		// 解析时间格式（如 "1s", "500ms"）
		if d, err := ctx.parseDuration(br.Time); err == nil {
			duration = d
		}
	} else if br.Strength != "" {
		// 根据强度设置默认时间
		switch br.Strength {
		case "none":
			duration = 0
		case "x-weak":
			duration = 100 * time.Millisecond
		case "weak":
			duration = 250 * time.Millisecond
		case "medium":
			duration = 500 * time.Millisecond
		case "strong":
			duration = 1 * time.Second
		case "x-strong":
			duration = 2 * time.Second
		default:
			duration = 500 * time.Millisecond
		}
	} else {
		duration = 500 * time.Millisecond // 默认停顿
	}

	// 添加停顿指令
	instruction := AudioInstruction{
		Type:     "break",
		Position: ctx.textPosition,
		Duration: duration,
	}

	ctx.result.Instructions = append(ctx.result.Instructions, instruction)
	ctx.currentTime += duration
}

// processAudio 处理音频插入
func (ctx *processingContext) processAudio(audio *Audio) {
	// 添加音频指令
	instruction := AudioInstruction{
		Type:      "audio",
		Position:  ctx.textPosition,
		AudioFile: audio.Src,
	}

	ctx.result.Instructions = append(ctx.result.Instructions, instruction)

	// 如果有备用文本，处理它
	for _, content := range audio.Content {
		ctx.processContent(content)
	}
}

// processVoice 处理声音变化
func (ctx *processingContext) processVoice(voice *Voice) {
	// 创建新的属性
	newProps := ctx.copyCurrentProperties()
	if voice.Name != "" {
		newProps.Voice = voice.Name
	}
	if voice.Gender != "" {
		newProps.Gender = voice.Gender
	}
	if voice.Languages != "" {
		newProps.Language = voice.Languages
	}

	ctx.pushProperties(newProps)

	// 处理子元素
	for _, content := range voice.Content {
		ctx.processContent(content)
	}

	ctx.popProperties()
}

// processProsody 处理韵律
func (ctx *processingContext) processProsody(prosody *Prosody) {
	newProps := ctx.copyCurrentProperties()

	if prosody.Rate != "" {
		newProps.Rate = prosody.Rate
	}
	if prosody.Pitch != "" {
		newProps.Pitch = prosody.Pitch
	}
	if prosody.Volume != "" {
		newProps.Volume = prosody.Volume
	}

	ctx.pushProperties(newProps)

	// 处理子元素
	for _, content := range prosody.Content {
		ctx.processContent(content)
	}

	ctx.popProperties()
}

// processEmphasis 处理强调
func (ctx *processingContext) processEmphasis(emphasis *Emphasis) {
	newProps := ctx.copyCurrentProperties()
	newProps.Emphasis = emphasis.Level

	ctx.pushProperties(newProps)

	// 处理子元素
	for _, content := range emphasis.Content {
		ctx.processContent(content)
	}

	ctx.popProperties()
}

// processPhoneme 处理发音
func (ctx *processingContext) processPhoneme(phoneme *Phoneme) {
	// 对于发音，我们使用原始文本，但可以添加发音指导信息
	for _, content := range phoneme.Content {
		ctx.processContent(content)
	}
}

// processSub 处理文本替换
func (ctx *processingContext) processSub(sub *Sub) {
	// 创建替换文本
	aliasText := &Text{Content: sub.Alias}
	ctx.processText(aliasText)
}

// processContainer 处理容器元素
func (ctx *processingContext) processContainer(element SSMLElement) {
	content := element.GetContent()

	switch element.(type) {
	case *Paragraph:
		// 处理段落内容
		for _, child := range content {
			ctx.processContent(child)
		}
		// 段落后添加短暂停顿
		ctx.addAutomaticBreak(300 * time.Millisecond)
	case *Sentence:
		// 处理句子内容
		for _, child := range content {
			ctx.processContent(child)
		}
		// 句子后添加短暂停顿
		ctx.addAutomaticBreak(200 * time.Millisecond)
	case *W:
		// 处理单词内容
		for _, child := range content {
			ctx.processContent(child)
		}
	}
}

// 辅助方法

// getCurrentProperties 获取当前属性
func (ctx *processingContext) getCurrentProperties() *AudioProperties {
	return ctx.copyProperties(ctx.propertyStack[len(ctx.propertyStack)-1])
}

// copyCurrentProperties 复制当前属性
func (ctx *processingContext) copyCurrentProperties() *AudioProperties {
	return ctx.copyProperties(ctx.propertyStack[len(ctx.propertyStack)-1])
}

// copyProperties 复制属性
func (ap *AudioProcessor) copyProperties(props *AudioProperties) *AudioProperties {
	return &AudioProperties{
		Rate:     props.Rate,
		Pitch:    props.Pitch,
		Volume:   props.Volume,
		Voice:    props.Voice,
		Gender:   props.Gender,
		Language: props.Language,
		Emphasis: props.Emphasis,
	}
}

// copyProperties (context method)
func (ctx *processingContext) copyProperties(props *AudioProperties) *AudioProperties {
	return ctx.processor.copyProperties(props)
}

// pushProperties 压入属性栈
func (ctx *processingContext) pushProperties(props *AudioProperties) {
	ctx.propertyStack = append(ctx.propertyStack, props)
}

// popProperties 弹出属性栈
func (ctx *processingContext) popProperties() {
	if len(ctx.propertyStack) > 1 {
		ctx.propertyStack = ctx.propertyStack[:len(ctx.propertyStack)-1]
	}
}

// calculateTextDuration 计算文本预期持续时间
func (ctx *processingContext) calculateTextDuration(text string) time.Duration {
	// 基础时间
	baseTime := time.Duration(len([]rune(text))) * ctx.processor.charToTimeRatio

	// 根据语速调整
	props := ctx.getCurrentProperties()
	multiplier := ctx.getRateMultiplier(props.Rate)

	return time.Duration(float64(baseTime) * multiplier)
}

// getRateMultiplier 获取语速倍数
func (ctx *processingContext) getRateMultiplier(rate string) float64 {
	switch rate {
	case "x-slow":
		return 2.0
	case "slow":
		return 1.5
	case "medium":
		return 1.0
	case "fast":
		return 0.75
	case "x-fast":
		return 0.5
	default:
		// 尝试解析百分比或数值
		if strings.HasSuffix(rate, "%") {
			if val, err := strconv.ParseFloat(strings.TrimSuffix(rate, "%"), 64); err == nil {
				return val / 100.0
			}
		}
		return 1.0
	}
}

// parseDuration 解析时间字符串
func (ctx *processingContext) parseDuration(timeStr string) (time.Duration, error) {
	timeStr = strings.TrimSpace(timeStr)

	if strings.HasSuffix(timeStr, "ms") {
		val := strings.TrimSuffix(timeStr, "ms")
		if ms, err := strconv.ParseFloat(val, 64); err == nil {
			return time.Duration(ms) * time.Millisecond, nil
		}
	} else if strings.HasSuffix(timeStr, "s") {
		val := strings.TrimSuffix(timeStr, "s")
		if s, err := strconv.ParseFloat(val, 64); err == nil {
			return time.Duration(s * float64(time.Second)), nil
		}
	}

	return time.ParseDuration(timeStr)
}

// addAutomaticBreak 添加自动停顿（避免重复添加）
func (ctx *processingContext) addAutomaticBreak(duration time.Duration) {
	// 检查最后一个指令是否已经是停顿，避免重复
	if len(ctx.result.Instructions) > 0 {
		lastInstruction := ctx.result.Instructions[len(ctx.result.Instructions)-1]
		if lastInstruction.Type == "auto_break" && lastInstruction.Position == ctx.textPosition {
			// 如果最后一个指令是同位置的自动停顿，就不再添加
			return
		}
	}

	instruction := AudioInstruction{
		Type:     "auto_break",
		Position: ctx.textPosition,
		Duration: duration,
	}

	ctx.result.Instructions = append(ctx.result.Instructions, instruction)
	ctx.currentTime += duration
}

// GetTextForTTS 获取用于 TTS 的纯文本
func (result *AudioProcessingResult) GetTextForTTS() string {
	return result.PlainText
}

// GetBreakInstructions 获取所有停顿指令
func (result *AudioProcessingResult) GetBreakInstructions() []AudioInstruction {
	var breaks []AudioInstruction
	for _, instruction := range result.Instructions {
		if instruction.Type == "break" || instruction.Type == "auto_break" {
			breaks = append(breaks, instruction)
		}
	}
	return breaks
}

// GetAudioInstructions 获取所有音频插入指令
func (result *AudioProcessingResult) GetAudioInstructions() []AudioInstruction {
	var audios []AudioInstruction
	for _, instruction := range result.Instructions {
		if instruction.Type == "audio" {
			audios = append(audios, instruction)
		}
	}
	return audios
}

// FormatProcessingReport 生成处理报告
func (result *AudioProcessingResult) FormatProcessingReport() string {
	var report strings.Builder

	report.WriteString("=== SSML 音频处理报告 ===\n\n")
	report.WriteString(fmt.Sprintf("纯文本: %s\n", result.PlainText))
	report.WriteString(fmt.Sprintf("预计总时长: %v\n", result.TotalDuration))
	report.WriteString(fmt.Sprintf("音频片段数: %d\n", len(result.Segments)))
	report.WriteString(fmt.Sprintf("处理指令数: %d\n\n", len(result.Instructions)))

	if len(result.Segments) > 0 {
		report.WriteString("=== 音频片段 ===\n")
		for i, segment := range result.Segments {
			report.WriteString(fmt.Sprintf("%d. [%v - %v] \"%s\"\n",
				i+1,
				segment.StartTime,
				segment.StartTime+segment.Duration,
				segment.Text))
			report.WriteString(fmt.Sprintf("   属性: 声音=%s, 语速=%s, 音调=%s, 音量=%s\n",
				segment.Properties.Voice,
				segment.Properties.Rate,
				segment.Properties.Pitch,
				segment.Properties.Volume))
		}
		report.WriteString("\n")
	}

	if len(result.Instructions) > 0 {
		report.WriteString("=== 处理指令 ===\n")
		for i, instruction := range result.Instructions {
			report.WriteString(fmt.Sprintf("%d. 位置 %d: %s",
				i+1, instruction.Position, instruction.Type))

			switch instruction.Type {
			case "break", "auto_break":
				report.WriteString(fmt.Sprintf(" (时长: %v)", instruction.Duration))
			case "audio":
				report.WriteString(fmt.Sprintf(" (文件: %s)", instruction.AudioFile))
			}
			report.WriteString("\n")
		}
	}

	return report.String()
}

// processContent 处理内容项（可能是 Text 或 SSMLElement）
func (ctx *processingContext) processContent(content interface{}) {
	if text, ok := content.(Text); ok {
		ctx.processText(&text)
	} else if text, ok := content.(*Text); ok {
		ctx.processText(text)
	} else if elem, ok := content.(SSMLElement); ok {
		ctx.processElement(elem)
	}
}
