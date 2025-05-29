package main

import (
	"fmt"
	"log"

	"ssml-parser/ssml"
)

func main() {
	// 示例1: 基本解析
	basicParsing()

	// 示例2: 基本构建
	basicBuilding()

	// 示例3: 语音控制
	voiceControl()

	// 示例4: 韵律控制
	prosodyControl()
}

// 基本解析示例
func basicParsing() {
	fmt.Println("=== 基本解析示例 ===")

	ssmlContent := `<?xml version="1.0" encoding="UTF-8"?>
	<speak version="1.0" xml:lang="zh-CN">
		欢迎使用 SSML 解析器！
		<break time="500ms"/>
		这是一个简单的示例。
	</speak>`

	parser := ssml.NewParser(nil)
	result, err := parser.Parse(ssmlContent)
	if err != nil {
		log.Printf("解析失败: %v", err)
		return
	}

	fmt.Printf("解析成功！版本: %s, 语言: %s\n", result.Root.Version, result.Root.Lang)
	fmt.Printf("内容元素数量: %d\n", len(result.Root.Content))

	// 序列化回 SSML
	serializer := ssml.NewSerializer(true)
	serialized, _ := serializer.Serialize(result.Root)
	fmt.Println("序列化结果:")
	fmt.Println(serialized)
	fmt.Println()
}

// 基本构建示例
func basicBuilding() {
	fmt.Println("=== 基本构建示例 ===")

	builder := ssml.NewBuilder()
	ssmlString, err := builder.
		Version("1.0").
		Lang("zh-CN").
		Text("这是用构建器创建的 SSML。").
		BreakTime("500ms").
		Text("非常简单易用！").
		BuildString(true)

	if err != nil {
		log.Printf("构建失败: %v", err)
		return
	}

	fmt.Println("构建的 SSML:")
	fmt.Println(ssmlString)
	fmt.Println()
}

// 语音控制示例
func voiceControl() {
	fmt.Println("=== 语音控制示例 ===")

	builder := ssml.NewBuilder()
	ssmlString, err := builder.
		Version("1.0").
		Lang("zh-CN").
		Text("默认声音说话。").
		BreakTime("500ms").
		VoiceNameText("xiaoxiao", "这是小小的声音。").
		BreakTime("500ms").
		Voice("male", "adult", "", "", "zh-CN", func(eb *ssml.ElementBuilder) {
			eb.Text("这是成年男性的声音。")
		}).
		BuildString(true)

	if err != nil {
		log.Printf("构建失败: %v", err)
		return
	}

	fmt.Println("语音控制 SSML:")
	fmt.Println(ssmlString)
	fmt.Println()
}

// 韵律控制示例
func prosodyControl() {
	fmt.Println("=== 韵律控制示例 ===")

	builder := ssml.NewBuilder()
	ssmlString, err := builder.
		Version("1.0").
		Lang("zh-CN").
		Text("正常语速和音调。").
		BreakTime("500ms").
		RateText("slow", "这段话语速较慢。").
		BreakTime("500ms").
		RateText("fast", "这段话语速较快。").
		BreakTime("500ms").
		PitchText("low", "这段话音调较低。").
		BreakTime("500ms").
		PitchText("high", "这段话音调较高。").
		BreakTime("500ms").
		VolumeText("soft", "这段话音量较小。").
		BreakTime("500ms").
		VolumeText("loud", "这段话音量较大。").
		BuildString(true)

	if err != nil {
		log.Printf("构建失败: %v", err)
		return
	}

	fmt.Println("韵律控制 SSML:")
	fmt.Println(ssmlString)
	fmt.Println()
}
