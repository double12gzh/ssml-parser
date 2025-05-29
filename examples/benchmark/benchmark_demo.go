package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"ssml-parser/ssml"
)

func main() {
	fmt.Println("SSML Parser Benchmark Demo")
	fmt.Println("==========================")

	// 演示不同场景的性能对比
	demonstrateSimpleVsComplex()
	demonstrateValidationModes()
	demonstrateLargeDocumentPerformance()
}

// 演示简单 vs 复杂 SSML 的性能差异
func demonstrateSimpleVsComplex() {
	fmt.Println("\n1. 简单 vs 复杂 SSML 性能对比:")
	fmt.Println("--------------------------------")

	simpleSSML := `<?xml version="1.0" encoding="UTF-8"?>
<speak version="1.0" xml:lang="zh-CN">
    Hello World
</speak>`

	complexSSML := `<?xml version="1.0" encoding="UTF-8"?>
<speak version="1.0" xml:lang="zh-CN">
    <voice name="xiaoxiao" gender="female" age="young">
        <prosody rate="fast" pitch="high" volume="loud" range="wide">
            这是一个包含多种韵律控制的复杂语音示例。
            <emphasis level="strong">这里是重点强调的内容。</emphasis>
            <break time="500ms"/>
            语音合成技术可以将文本转换为自然的语音。
        </prosody>
        <break time="1s"/>
        <sub alias="人工智能">AI</sub>技术正在快速发展。
    </voice>
    <audio src="notification.wav">通知音效</audio>
    <phoneme alphabet="ipa" ph="ˈteknoʊlədʒi">technology</phoneme>
    <p>
        <s>这是第一个句子。</s>
        <s>这是第二个句子。</s>
    </p>
    <w role="verb">运行</w>在现代设备上。
</speak>`

	parser := ssml.NewParser(nil)

	// 测试简单 SSML
	start := time.Now()
	iterations := 10000
	for i := 0; i < iterations; i++ {
		_, err := parser.Parse(simpleSSML)
		if err != nil {
			log.Fatal(err)
		}
	}
	simpleTime := time.Since(start)

	// 测试复杂 SSML
	start = time.Now()
	for i := 0; i < iterations; i++ {
		_, err := parser.Parse(complexSSML)
		if err != nil {
			log.Fatal(err)
		}
	}
	complexTime := time.Since(start)

	fmt.Printf("简单 SSML (%d 次解析): %v (平均 %v/次)\n",
		iterations, simpleTime, simpleTime/time.Duration(iterations))
	fmt.Printf("复杂 SSML (%d 次解析): %v (平均 %v/次)\n",
		iterations, complexTime, complexTime/time.Duration(iterations))
	fmt.Printf("复杂度增加系数: %.2fx\n", float64(complexTime)/float64(simpleTime))
}

// 演示不同验证模式的性能差异
func demonstrateValidationModes() {
	fmt.Println("\n2. 验证模式性能对比:")
	fmt.Println("--------------------")

	ssmlContent := `<?xml version="1.0" encoding="UTF-8"?>
<speak version="1.0" xml:lang="zh-CN">
    <voice name="xiaoxiao" gender="female">
        <prosody rate="fast" pitch="high">
            验证模式性能测试。
            <emphasis level="strong">重要内容</emphasis>
        </prosody>
        <break time="500ms"/>
        <audio src="test.wav">音频内容</audio>
    </voice>
</speak>`

	// 基本验证模式
	basicConfig := &ssml.ValidationConfig{
		StrictMode:           false,
		AllowUnknownElements: true,
		MaxNestingDepth:      10,
	}
	basicParser := ssml.NewParser(basicConfig)

	// 严格验证模式
	strictConfig := &ssml.ValidationConfig{
		StrictMode:           true,
		AllowUnknownElements: false,
		MaxNestingDepth:      10,
	}
	strictParser := ssml.NewParser(strictConfig)

	iterations := 5000

	// 测试基本验证
	start := time.Now()
	for i := 0; i < iterations; i++ {
		_, err := basicParser.Parse(ssmlContent)
		if err != nil {
			log.Fatal(err)
		}
	}
	basicTime := time.Since(start)

	// 测试严格验证
	start = time.Now()
	for i := 0; i < iterations; i++ {
		_, err := strictParser.Parse(ssmlContent)
		if err != nil {
			log.Fatal(err)
		}
	}
	strictTime := time.Since(start)

	fmt.Printf("基本验证模式 (%d 次): %v (平均 %v/次)\n",
		iterations, basicTime, basicTime/time.Duration(iterations))
	fmt.Printf("严格验证模式 (%d 次): %v (平均 %v/次)\n",
		iterations, strictTime, strictTime/time.Duration(iterations))
	fmt.Printf("严格模式开销: %.2fx\n", float64(strictTime)/float64(basicTime))
}

// 演示大型文档解析性能
func demonstrateLargeDocumentPerformance() {
	fmt.Println("\n3. 大型文档解析性能:")
	fmt.Println("--------------------")

	parser := ssml.NewParser(nil)

	// 测试不同大小的文档
	sizes := []int{10, 100, 500, 1000}

	for _, size := range sizes {
		// 生成指定大小的 SSML 文档
		ssmlContent := generateLargeSSML(size)

		// 测试解析时间
		start := time.Now()
		result, err := parser.Parse(ssmlContent)
		if err != nil {
			log.Fatal(err)
		}
		duration := time.Since(start)

		fmt.Printf("文档大小: %d 元素, 解析时间: %v, 吞吐量: %.2f MB/s\n",
			size, duration, float64(len(ssmlContent))/1024/1024/duration.Seconds())

		// 显示解析结果统计
		fmt.Printf("  - 文档大小: %d 字节\n", len(ssmlContent))
		fmt.Printf("  - 内容元素: %d 个\n", len(result.Root.Content))
		fmt.Printf("  - 平均每元素: %v\n", duration/time.Duration(size))
		fmt.Println()
	}
}

// 生成指定大小的 SSML 文档
func generateLargeSSML(elementCount int) string {
	var builder strings.Builder
	builder.WriteString(`<?xml version="1.0" encoding="UTF-8"?>`)
	builder.WriteString(`<speak version="1.0" xml:lang="zh-CN">`)

	for i := 0; i < elementCount; i++ {
		switch i % 6 {
		case 0:
			builder.WriteString(fmt.Sprintf(`
    <voice name="voice%d">
        这是第 %d 个声音片段。
    </voice>`, i%5, i+1))
		case 1:
			builder.WriteString(fmt.Sprintf(`
    <prosody rate="medium" pitch="medium">
        这是第 %d 个韵律片段。
    </prosody>`, i+1))
		case 2:
			builder.WriteString(fmt.Sprintf(`
    <emphasis level="moderate">
        这是第 %d 个强调片段。
    </emphasis>`, i+1))
		case 3:
			builder.WriteString(`<break time="100ms"/>`)
		case 4:
			builder.WriteString(fmt.Sprintf(`
    <sub alias="替换文本%d">原始文本%d</sub>`, i+1, i+1))
		case 5:
			builder.WriteString(fmt.Sprintf(`
    <p>这是第 %d 个段落。</p>`, i+1))
		}
	}

	builder.WriteString(`</speak>`)
	return builder.String()
}
