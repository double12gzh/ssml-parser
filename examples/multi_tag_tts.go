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
	fmt.Println("🚀 多标签 SSML 并发 TTS 处理示例")
	fmt.Println("================================")

	// 包含多种标签的复杂 SSML
	ssmlContent := `<?xml version="1.0" encoding="UTF-8"?>
<speak version="1.0" xml:lang="zh-CN">
    <voice name="xiaoxiao" gender="female">
        欢迎来到我们的智能客服系统！
        <break time="800ms"/>
        
        <prosody rate="slow" pitch="low" volume="loud">
            请注意，我们的服务时间是工作日上午九点到下午六点。
        </prosody>
        <break time="500ms"/>
        
        <emphasis level="strong">
            重要提醒：
        </emphasis>
        请保持您的账户信息安全，不要向他人透露密码。
        <break time="1s"/>
        
        <voice name="xiaogang" gender="male">
            <prosody rate="fast" pitch="high">
                如果您需要技术支持，请按1。
            </prosody>
            <break time="300ms"/>
            
            <emphasis level="moderate">
                如果您需要账户服务，请按2。
            </emphasis>
        </voice>
        <break time="600ms"/>
        
        <sub alias="人工智能">AI</sub>助手将为您提供24小时服务。
        <break time="400ms"/>
        
        <prosody rate="x-slow" volume="soft">
            感谢您的耐心等待，祝您使用愉快！
        </prosody>
    </voice>
</speak>`

	fmt.Println("📝 复杂 SSML 内容:")
	fmt.Println(ssmlContent)
	fmt.Println()

	// 运行多标签处理示例
	RunMultiTagTTSExample(ssmlContent)
}

// MultiTagSegment 多标签文本片段
type MultiTagSegment struct {
	Text          string                // 文本内容
	OriginalStart int                   // 原始位置
	OriginalEnd   int                   // 结束位置
	Properties    *ssml.AudioProperties // 音频属性
	SentenceIndex int                   // 句子索引
	TagStack      []TagInfo             // 标签栈信息
}

// TagInfo 标签信息
type TagInfo struct {
	Type       string            // 标签类型
	Attributes map[string]string // 标签属性
	StartPos   int               // 开始位置
	EndPos     int               // 结束位置
}

// MultiTagTTSResult 多标签TTS结果
type MultiTagTTSResult struct {
	Index       int
	Text        string
	AudioData   *ssml.AudioData
	Properties  *ssml.AudioProperties
	AppliedTags []TagInfo
	Error       error
	ProcessTime time.Duration
}

// ComplexTagInsertPoint 复杂标签插入点
type ComplexTagInsertPoint struct {
	Type          string
	AfterSentence int
	Duration      time.Duration
	Properties    map[string]string
	Priority      int // 插入优先级
}

// RunMultiTagTTSExample 运行多标签TTS示例
func RunMultiTagTTSExample(ssmlContent string) {
	// 步骤 1: 解析复杂 SSML
	fmt.Println("🔍 步骤 1: 解析复杂 SSML...")

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

	fullText := result.GetTextForTTS()
	fmt.Printf("✓ 提取完整文本: %s\n", fullText)
	fmt.Printf("✓ 找到 %d 个音频片段，%d 个处理指令\n",
		len(result.Segments), len(result.Instructions))

	// 显示所有标签信息
	fmt.Println("\n🏷️  检测到的标签类型:")
	tagTypes := make(map[string]int)
	for _, instruction := range result.Instructions {
		tagTypes[instruction.Type]++
	}
	for tagType, count := range tagTypes {
		fmt.Printf("  - %s: %d 个\n", tagType, count)
	}

	// 步骤 2: 智能切句
	fmt.Println("\n📝 步骤 2: 智能切句处理...")
	sentences := AdvancedSentenceSegmentation(fullText)
	fmt.Printf("✓ 切分为 %d 个句子:\n", len(sentences))
	for i, sentence := range sentences {
		fmt.Printf("  %d. \"%s\"\n", i+1, sentence)
	}

	// 步骤 3: 构建多标签映射
	fmt.Println("\n🗺️  步骤 3: 构建多标签映射...")
	multiTagSegments := BuildMultiTagMapping(fullText, sentences, result.Segments, result.Instructions)
	fmt.Printf("✓ 建立了 %d 个多标签片段映射\n", len(multiTagSegments))

	// 显示每个片段的标签信息
	for i, segment := range multiTagSegments {
		fmt.Printf("  [%d] \"%s\" - 标签: ", i+1, segment.Text)
		if len(segment.TagStack) == 0 {
			fmt.Printf("无")
		} else {
			for j, tag := range segment.TagStack {
				if j > 0 {
					fmt.Printf(", ")
				}
				fmt.Printf("%s", tag.Type)
			}
		}
		fmt.Println()
	}

	// 步骤 4: 计算复杂标签插入点
	fmt.Println("\n📍 步骤 4: 计算复杂标签插入点...")
	insertPoints := CalculateComplexTagPositions(fullText, sentences, result.Instructions)
	fmt.Printf("✓ 计算了 %d 个插入点:\n", len(insertPoints))
	for i, point := range insertPoints {
		fmt.Printf("  %d. %s 在句子 %d 后", i+1, point.Type, point.AfterSentence+1)
		if point.Duration > 0 {
			fmt.Printf("，停顿 %v", point.Duration)
		}
		if len(point.Properties) > 0 {
			fmt.Printf("，属性: %v", point.Properties)
		}
		fmt.Println()
	}

	// 步骤 5: 并发多标签TTS处理
	fmt.Println("\n🎵 步骤 5: 并发多标签 TTS 处理...")
	start := time.Now()
	ttsResults := ConcurrentMultiTagTTS(multiTagSegments)
	ttsTime := time.Since(start)

	fmt.Printf("✓ 并发处理完成，耗时: %v\n", ttsTime)
	fmt.Printf("✓ 成功生成 %d 个音频片段\n", len(ttsResults))

	// 显示每个片段的处理结果
	for i, result := range ttsResults {
		fmt.Printf("  [%d] \"%s\" - 时长: %v，标签: %d 个\n",
			i+1, result.Text, result.AudioData.Duration, len(result.AppliedTags))
	}

	// 步骤 6: 智能音频组装
	fmt.Println("\n🔧 步骤 6: 智能音频组装...")
	finalAudio := AssembleMultiTagAudio(ttsResults, insertPoints)

	fmt.Printf("✓ 音频组装完成\n")
	fmt.Printf("  - 最终时长: %v\n", finalAudio.Duration)
	fmt.Printf("  - 数据点数: %d\n", len(finalAudio.Data))

	// 步骤 7: 保存文件
	fmt.Println("\n💾 步骤 7: 保存多标签音频文件...")
	filename := "multi_tag_output.wav"
	err = SaveMultiTagAudio(finalAudio, filename)
	if err != nil {
		log.Fatalf("保存失败: %v", err)
	}

	fmt.Printf("✓ 音频文件保存成功: %s\n", filename)

	// 步骤 8: 详细统计
	fmt.Println("\n📊 多标签处理统计:")
	fmt.Printf("  - 原始 SSML 长度: %d 字符\n", len(ssmlContent))
	fmt.Printf("  - 提取文本长度: %d 字符\n", len(fullText))
	fmt.Printf("  - 切句数量: %d 句\n", len(sentences))
	fmt.Printf("  - 标签类型数: %d 种\n", len(tagTypes))
	fmt.Printf("  - 总标签数量: %d 个\n", len(result.Instructions))
	fmt.Printf("  - 并发处理时间: %v\n", ttsTime)
	fmt.Printf("  - 平均每句耗时: %v\n", ttsTime/time.Duration(len(sentences)))
	fmt.Printf("  - 最终音频大小: %.2f MB\n", float64(len(finalAudio.Data)*2)/1024/1024)

	fmt.Println("\n🎉 多标签 SSML 处理完成！")
}

// AdvancedSentenceSegmentation 高级切句
func AdvancedSentenceSegmentation(text string) []string {
	fmt.Println("🔪 高级切句处理...")

	time.Sleep(40 * time.Millisecond)

	var sentences []string

	// 1. 按主要标点切分
	mainParts := regexp.MustCompile(`[。！？]+`).Split(text, -1)

	for _, part := range mainParts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		// 2. 处理长句子
		if len([]rune(part)) > 40 {
			// 按逗号、分号切分
			subParts := regexp.MustCompile(`[，；：]+`).Split(part, -1)
			for _, subPart := range subParts {
				subPart = strings.TrimSpace(subPart)
				if subPart != "" {
					// 3. 进一步处理超长片段
					if len([]rune(subPart)) > 25 {
						// 按顿号、括号切分
						finalParts := regexp.MustCompile(`[、（）【】]+`).Split(subPart, -1)
						for _, finalPart := range finalParts {
							finalPart = strings.TrimSpace(finalPart)
							if finalPart != "" {
								sentences = append(sentences, finalPart)
							}
						}
					} else {
						sentences = append(sentences, subPart)
					}
				}
			}
		} else {
			sentences = append(sentences, part)
		}
	}

	fmt.Printf("  ✓ 高级切句完成: %d 句\n", len(sentences))
	return sentences
}

// BuildMultiTagMapping 构建多标签映射
func BuildMultiTagMapping(fullText string, sentences []string, originalSegments []ssml.AudioSegment, instructions []ssml.AudioInstruction) []MultiTagSegment {
	var segments []MultiTagSegment

	currentPos := 0

	for i, sentence := range sentences {
		// 查找句子在原文中的位置
		startPos := strings.Index(fullText[currentPos:], sentence)
		if startPos == -1 {
			startPos = findFuzzyMatch(fullText[currentPos:], sentence)
		}

		if startPos != -1 {
			actualStart := currentPos + startPos
			actualEnd := actualStart + len(sentence)

			// 查找应用到这个片段的所有标签
			tagStack := findTagsForPosition(actualStart, actualEnd, instructions)

			// 查找对应的音频属性
			properties := findPropertiesForSegment(actualStart, actualEnd, originalSegments, tagStack)

			segment := MultiTagSegment{
				Text:          sentence,
				OriginalStart: actualStart,
				OriginalEnd:   actualEnd,
				Properties:    properties,
				SentenceIndex: i,
				TagStack:      tagStack,
			}

			segments = append(segments, segment)
			currentPos = actualEnd
		}
	}

	return segments
}

// findFuzzyMatch 模糊匹配
func findFuzzyMatch(text string, target string) int {
	words := strings.Fields(target)
	if len(words) == 0 {
		return -1
	}

	// 尝试匹配前几个字符
	if len(target) > 3 {
		return strings.Index(text, target[:3])
	}

	return strings.Index(text, words[0])
}

// findTagsForPosition 查找位置对应的标签
func findTagsForPosition(start, end int, instructions []ssml.AudioInstruction) []TagInfo {
	var tags []TagInfo

	for _, instruction := range instructions {
		// 检查标签是否影响这个位置
		if instruction.Position >= start && instruction.Position <= end {
			tag := TagInfo{
				Type:       instruction.Type,
				Attributes: make(map[string]string),
				StartPos:   instruction.Position,
				EndPos:     instruction.Position,
			}

			// 根据标签类型设置属性
			switch instruction.Type {
			case "break":
				tag.Attributes["time"] = instruction.Duration.String()
			case "prosody":
				// 这里需要从原始指令中提取prosody属性
				tag.Attributes["rate"] = "medium"
				tag.Attributes["pitch"] = "medium"
				tag.Attributes["volume"] = "medium"
			case "emphasis":
				tag.Attributes["level"] = "moderate"
			case "voice":
				tag.Attributes["name"] = "default"
				tag.Attributes["gender"] = "neutral"
			}

			tags = append(tags, tag)
		}
	}

	return tags
}

// findPropertiesForSegment 为片段查找音频属性
func findPropertiesForSegment(start, end int, segments []ssml.AudioSegment, tagStack []TagInfo) *ssml.AudioProperties {
	// 从原始片段中查找属性
	for _, segment := range segments {
		if segment.Properties != nil {
			return segment.Properties
		}
	}

	// 根据标签栈构建属性
	properties := &ssml.AudioProperties{
		Rate:     "medium",
		Pitch:    "medium",
		Volume:   "medium",
		Voice:    "default",
		Gender:   "neutral",
		Language: "zh-CN",
		Emphasis: "none",
	}

	// 应用标签栈中的属性
	for _, tag := range tagStack {
		switch tag.Type {
		case "prosody":
			if rate, exists := tag.Attributes["rate"]; exists {
				properties.Rate = rate
			}
			if pitch, exists := tag.Attributes["pitch"]; exists {
				properties.Pitch = pitch
			}
			if volume, exists := tag.Attributes["volume"]; exists {
				properties.Volume = volume
			}
		case "emphasis":
			if level, exists := tag.Attributes["level"]; exists {
				properties.Emphasis = level
			}
		case "voice":
			if name, exists := tag.Attributes["name"]; exists {
				properties.Voice = name
			}
			if gender, exists := tag.Attributes["gender"]; exists {
				properties.Gender = gender
			}
		}
	}

	return properties
}

// CalculateComplexTagPositions 计算复杂标签位置
func CalculateComplexTagPositions(fullText string, sentences []string, instructions []ssml.AudioInstruction) []ComplexTagInsertPoint {
	var points []ComplexTagInsertPoint

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

	// 处理每个指令
	for _, instruction := range instructions {
		if instruction.Type == "break" {
			// 找到break应该插入的位置
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

			point := ComplexTagInsertPoint{
				Type:          instruction.Type,
				AfterSentence: afterSentence,
				Duration:      instruction.Duration,
				Properties:    make(map[string]string),
				Priority:      1, // break有最高优先级
			}

			points = append(points, point)
		}
		// 可以添加其他类型的标签处理
	}

	// 按优先级和位置排序
	sort.Slice(points, func(i, j int) bool {
		if points[i].AfterSentence == points[j].AfterSentence {
			return points[i].Priority > points[j].Priority
		}
		return points[i].AfterSentence < points[j].AfterSentence
	})

	return points
}

// ConcurrentMultiTagTTS 并发多标签TTS处理
func ConcurrentMultiTagTTS(segments []MultiTagSegment) []MultiTagTTSResult {
	fmt.Printf("🔄 并发处理 %d 个多标签片段...\n", len(segments))

	results := make(chan MultiTagTTSResult, len(segments))
	var wg sync.WaitGroup

	// 创建增强的TTS适配器
	tts := ssml.NewMockTTSAdapter(44100)

	for i, segment := range segments {
		wg.Add(1)
		go func(index int, seg MultiTagSegment) {
			defer wg.Done()

			fmt.Printf("  📝 [%d] 处理多标签: \"%s\" (标签: %d 个)\n",
				index+1, seg.Text, len(seg.TagStack))

			start := time.Now()

			// 模拟复杂标签处理时间
			processingTime := time.Duration(len([]rune(seg.Text))*10+len(seg.TagStack)*20) * time.Millisecond
			time.Sleep(processingTime)

			// 生成音频，应用所有标签效果
			audioData, err := tts.GenerateAudio(seg.Text, seg.Properties)

			// 后处理：应用标签效果
			if err == nil && audioData != nil {
				audioData = applyTagEffects(audioData, seg.TagStack)
			}

			duration := time.Since(start)

			if err != nil {
				fmt.Printf("  ❌ [%d] 处理失败: %v\n", index+1, err)
			} else {
				fmt.Printf("  ✅ [%d] 完成: %v (耗时: %v, 标签: %d)\n",
					index+1, audioData.Duration, duration, len(seg.TagStack))
			}

			results <- MultiTagTTSResult{
				Index:       index,
				Text:        seg.Text,
				AudioData:   audioData,
				Properties:  seg.Properties,
				AppliedTags: seg.TagStack,
				Error:       err,
				ProcessTime: duration,
			}
		}(i, segment)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	var ttsResults []MultiTagTTSResult
	for result := range results {
		ttsResults = append(ttsResults, result)
	}

	sort.Slice(ttsResults, func(i, j int) bool {
		return ttsResults[i].Index < ttsResults[j].Index
	})

	return ttsResults
}

// applyTagEffects 应用标签效果
func applyTagEffects(audio *ssml.AudioData, tags []TagInfo) *ssml.AudioData {
	// 复制音频数据
	processedAudio := &ssml.AudioData{
		SampleRate: audio.SampleRate,
		Channels:   audio.Channels,
		Duration:   audio.Duration,
		Data:       make([]float64, len(audio.Data)),
	}
	copy(processedAudio.Data, audio.Data)

	// 应用每个标签的效果
	for _, tag := range tags {
		switch tag.Type {
		case "emphasis":
			// 应用强调效果（增加音量）
			level := tag.Attributes["level"]
			multiplier := 1.0
			switch level {
			case "strong":
				multiplier = 1.5
			case "moderate":
				multiplier = 1.2
			case "reduced":
				multiplier = 0.8
			}

			for i := range processedAudio.Data {
				processedAudio.Data[i] *= multiplier
				// 防止溢出
				if processedAudio.Data[i] > 1.0 {
					processedAudio.Data[i] = 1.0
				} else if processedAudio.Data[i] < -1.0 {
					processedAudio.Data[i] = -1.0
				}
			}

		case "prosody":
			// 应用韵律效果（这里简化处理）
			rate := tag.Attributes["rate"]
			switch rate {
			case "fast", "x-fast":
				// 快速语音效果（这里简化为音量调整）
				for i := range processedAudio.Data {
					processedAudio.Data[i] *= 1.1
				}
			case "slow", "x-slow":
				// 慢速语音效果
				for i := range processedAudio.Data {
					processedAudio.Data[i] *= 0.9
				}
			}
		}
	}

	return processedAudio
}

// AssembleMultiTagAudio 组装多标签音频
func AssembleMultiTagAudio(ttsResults []MultiTagTTSResult, insertPoints []ComplexTagInsertPoint) *ssml.AudioData {
	fmt.Println("🔧 组装多标签音频...")

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

		fmt.Printf("  🎵 [%d] 添加多标签音频: \"%s\" (%v, 标签: %d)\n",
			i+1, result.Text, result.AudioData.Duration, len(result.AppliedTags))

		// 添加音频数据
		finalData = append(finalData, result.AudioData.Data...)

		// 检查是否需要插入标签效果
		for _, point := range insertPoints {
			if point.AfterSentence == result.Index {
				switch point.Type {
				case "break":
					fmt.Printf("  ⏸️  [%d] 插入停顿: %v\n", i+1, point.Duration)

					silenceSamples := int(point.Duration.Seconds() * float64(sampleRate))
					silence := make([]float64, silenceSamples)
					finalData = append(finalData, silence...)

				default:
					fmt.Printf("  🏷️  [%d] 应用 %s 效果\n", i+1, point.Type)
				}
			}
		}
	}

	totalDuration := time.Duration(len(finalData)) * time.Second / time.Duration(sampleRate)
	fmt.Printf("✓ 多标签音频组装完成，总时长: %v，样本: %d\n", totalDuration, len(finalData))

	return &ssml.AudioData{
		SampleRate: sampleRate,
		Channels:   1,
		Duration:   totalDuration,
		Data:       finalData,
	}
}

// SaveMultiTagAudio 保存多标签音频
func SaveMultiTagAudio(audio *ssml.AudioData, filename string) error {
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
