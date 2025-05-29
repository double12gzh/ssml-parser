package main

import (
	"fmt"
	"log"
	"strings"

	"ssml-parser/ssml"
)

func main() {
	fmt.Println("SSML Parser Demo")
	fmt.Println("================")

	// 演示1: 解析 SSML
	demonstrateParser()

	// 演示2: 使用构建器创建 SSML
	demonstrateBuilder()

	// 演示3: 复杂的 SSML 示例
	demonstrateComplexSSML()
}

// 演示解析器功能
func demonstrateParser() {
	fmt.Println("\n1. 解析 SSML 示例:")
	fmt.Println("-------------------")

	ssmlContent := `<?xml version="1.0" encoding="UTF-8"?>
<speak version="1.0" xml:lang="zh-CN">
    <voice name="xiaoxiao">
        <prosody rate="fast" pitch="high">
            这是一个快速、高音调的语音示例。
        </prosody>
    </voice>
    <break time="1s"/>
    <emphasis level="strong">这里是重点强调的内容。</emphasis>
    <audio src="beep.wav">备用文本</audio>
    <phoneme alphabet="ipa" ph="ˈhəloʊ">Hello</phoneme>
    <sub alias="世界贸易组织">WTO</sub>
</speak>`

	parser := ssml.NewParser(nil)
	result, err := parser.Parse(ssmlContent)
	if err != nil {
		log.Printf("解析失败: %v", err)
		return
	}

	fmt.Printf("解析成功！\n")
	fmt.Printf("版本: %s\n", result.Root.Version)
	fmt.Printf("语言: %s\n", result.Root.Lang)
	fmt.Printf("内容元素数量: %d\n", len(result.Root.Content))

	if len(result.Warnings) > 0 {
		fmt.Printf("警告: %v\n", result.Warnings)
	}

	if len(result.Errors) > 0 {
		fmt.Printf("错误: %v\n", result.Errors)
	}

	// 将解析结果序列化回 SSML
	serializer := ssml.NewSerializer(true)
	serialized, err := serializer.Serialize(result.Root)
	if err != nil {
		log.Printf("序列化失败: %v", err)
		return
	}

	fmt.Println("\n序列化结果:")
	fmt.Println(serialized)
}

// 演示构建器功能
func demonstrateBuilder() {
	fmt.Println("\n2. 使用构建器创建 SSML:")
	fmt.Println("-------------------------")

	builder := ssml.NewBuilder()

	// 构建 SSML
	ssmlString, err := builder.
		Version("1.0").
		Lang("zh-CN").
		Text("欢迎使用 SSML 解析器！").
		BreakTime("500ms").
		VoiceNameText("xiaoxiao", "这是小小的声音。").
		BreakStrength("medium").
		RateText("slow", "这段话语速较慢。").
		BreakTime("1s").
		EmphasisText("strong", "这里是重点内容！").
		Audio("notification.wav", "通知音效").
		SubText("人工智能", "AI").
		Text("技术正在快速发展。").
		BuildString(true)

	if err != nil {
		log.Printf("构建失败: %v", err)
		return
	}

	fmt.Println("构建的 SSML:")
	fmt.Println(ssmlString)
}

// 演示复杂的 SSML 示例
func demonstrateComplexSSML() {
	fmt.Println("\n3. 复杂 SSML 示例:")
	fmt.Println("------------------")

	builder := ssml.NewBuilder()

	ssmlString, err := builder.
		Version("1.0").
		Lang("zh-CN").
		ParagraphText("这是第一段内容，用于介绍语音合成技术。").
		BreakTime("1s").
		Voice("female", "young", "", "xiaoxiao", "zh-CN", func(eb *ssml.ElementBuilder) {
			eb.Text("使用女性年轻声音：")
			eb.Prosody("fast", "high", "", "loud", func(nested *ssml.ElementBuilder) {
				nested.Text("快速、高音调、大音量的语音效果。")
			})
			eb.BreakTime("500ms")
			eb.Emphasis("moderate", func(nested *ssml.ElementBuilder) {
				nested.Text("适中强调的内容。")
			})
		}).
		BreakTime("1s").
		Sentence(func(eb *ssml.ElementBuilder) {
			eb.Text("这是一个句子，包含")
			eb.PhonemeText("ipa", "ˈteknoʊlədʒi", "technology")
			eb.Text("的发音指导。")
		}).
		BreakTime("2s").
		Text("语音合成（").
		SubText("文本转语音", "TTS").
		Text("）技术可以将文本转换为自然的语音。").
		BuildString(true)

	if err != nil {
		log.Printf("构建复杂SSML失败: %v", err)
		return
	}

	fmt.Println("复杂 SSML 示例:")
	fmt.Println(ssmlString)

	// 解析刚刚构建的 SSML
	fmt.Println("\n解析构建的 SSML:")
	parser := ssml.NewParser(nil)
	result, err := parser.Parse(ssmlString)
	if err != nil {
		log.Printf("解析失败: %v", err)
		return
	}

	fmt.Printf("解析成功，内容元素数量: %d\n", len(result.Root.Content))
	printContentSummary(result.Root.Content, 0)
}

// 打印内容摘要
func printContentSummary(content []interface{}, indent int) {
	prefix := strings.Repeat("  ", indent)

	for i, item := range content {
		switch v := item.(type) {
		case ssml.Text:
			fmt.Printf("%s[%d] 文本: \"%.30s...\"\n", prefix, i, v.Content)
		case *ssml.Audio:
			fmt.Printf("%s[%d] 音频: src=%s, 内容数=%d\n", prefix, i, v.Src, len(v.Content))
			if len(v.Content) > 0 {
				printContentSummary(v.Content, indent+1)
			}
		case *ssml.Break:
			fmt.Printf("%s[%d] 停顿: time=%s, strength=%s\n", prefix, i, v.Time, v.Strength)
		case *ssml.Emphasis:
			fmt.Printf("%s[%d] 强调: level=%s, 内容数=%d\n", prefix, i, v.Level, len(v.Content))
			if len(v.Content) > 0 {
				printContentSummary(v.Content, indent+1)
			}
		case *ssml.Paragraph:
			fmt.Printf("%s[%d] 段落: 内容数=%d\n", prefix, i, len(v.Content))
			if len(v.Content) > 0 {
				printContentSummary(v.Content, indent+1)
			}
		case *ssml.Phoneme:
			fmt.Printf("%s[%d] 发音: alphabet=%s, ph=%s, 内容数=%d\n", prefix, i, v.Alphabet, v.Ph, len(v.Content))
			if len(v.Content) > 0 {
				printContentSummary(v.Content, indent+1)
			}
		case *ssml.Prosody:
			fmt.Printf("%s[%d] 韵律: rate=%s, pitch=%s, volume=%s, 内容数=%d\n", prefix, i, v.Rate, v.Pitch, v.Volume, len(v.Content))
			if len(v.Content) > 0 {
				printContentSummary(v.Content, indent+1)
			}
		case *ssml.Sentence:
			fmt.Printf("%s[%d] 句子: 内容数=%d\n", prefix, i, len(v.Content))
			if len(v.Content) > 0 {
				printContentSummary(v.Content, indent+1)
			}
		case *ssml.Sub:
			fmt.Printf("%s[%d] 替换: alias=%s, 内容数=%d\n", prefix, i, v.Alias, len(v.Content))
			if len(v.Content) > 0 {
				printContentSummary(v.Content, indent+1)
			}
		case *ssml.Voice:
			fmt.Printf("%s[%d] 声音: name=%s, gender=%s, age=%s, 内容数=%d\n", prefix, i, v.Name, v.Gender, v.Age, len(v.Content))
			if len(v.Content) > 0 {
				printContentSummary(v.Content, indent+1)
			}
		case *ssml.W:
			fmt.Printf("%s[%d] 单词: role=%s, 内容数=%d\n", prefix, i, v.Role, len(v.Content))
			if len(v.Content) > 0 {
				printContentSummary(v.Content, indent+1)
			}
		default:
			fmt.Printf("%s[%d] 未知类型: %T\n", prefix, i, v)
		}
	}
}
