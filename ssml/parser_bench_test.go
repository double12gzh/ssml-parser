package ssml

import (
	"fmt"
	"strings"
	"testing"
)

// BenchmarkParseSimpleSSML 测试解析简单 SSML 的性能
func BenchmarkParseSimpleSSML(b *testing.B) {
	ssmlContent := `<?xml version="1.0" encoding="UTF-8"?>
<speak version="1.0" xml:lang="zh-CN">
    Hello World
</speak>`

	parser := NewParser(nil)
	b.SetBytes(int64(len(ssmlContent)))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := parser.Parse(ssmlContent)
		if err != nil {
			b.Fatalf("解析失败: %v", err)
		}
	}
}

// BenchmarkParseComplexSSML 测试解析复杂 SSML 的性能
func BenchmarkParseComplexSSML(b *testing.B) {
	ssmlContent := `<?xml version="1.0" encoding="UTF-8"?>
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

	parser := NewParser(nil)
	b.SetBytes(int64(len(ssmlContent)))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := parser.Parse(ssmlContent)
		if err != nil {
			b.Fatalf("解析失败: %v", err)
		}
	}
}

// BenchmarkParseLargeSSML 测试解析大型 SSML 文档的性能
func BenchmarkParseLargeSSML(b *testing.B) {
	// 生成大型 SSML 文档
	var builder strings.Builder
	builder.WriteString(`<?xml version="1.0" encoding="UTF-8"?>`)
	builder.WriteString(`<speak version="1.0" xml:lang="zh-CN">`)

	// 添加 1000 个段落
	for i := 0; i < 1000; i++ {
		builder.WriteString(fmt.Sprintf(`
    <p>
        <voice name="voice%d">
            这是第 %d 个段落。
            <prosody rate="medium" pitch="medium">
                包含韵律控制的文本内容。
                <emphasis level="moderate">重要信息 %d</emphasis>
            </prosody>
            <break time="100ms"/>
        </voice>
    </p>`, i%10, i+1, i+1))
	}

	builder.WriteString(`</speak>`)
	ssmlContent := builder.String()

	parser := NewParser(nil)
	b.SetBytes(int64(len(ssmlContent)))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := parser.Parse(ssmlContent)
		if err != nil {
			b.Fatalf("解析失败: %v", err)
		}
	}
}

// BenchmarkParseDeepNesting 测试解析深度嵌套 SSML 的性能
func BenchmarkParseDeepNesting(b *testing.B) {
	var builder strings.Builder
	builder.WriteString(`<?xml version="1.0" encoding="UTF-8"?>`)
	builder.WriteString(`<speak version="1.0" xml:lang="zh-CN">`)

	// 创建 50 层嵌套
	nestingDepth := 50
	for i := 0; i < nestingDepth; i++ {
		if i%2 == 0 {
			builder.WriteString(`<emphasis level="moderate">`)
		} else {
			builder.WriteString(`<prosody rate="medium">`)
		}
	}

	builder.WriteString(`深度嵌套的内容`)

	for i := nestingDepth - 1; i >= 0; i-- {
		if i%2 == 0 {
			builder.WriteString(`</emphasis>`)
		} else {
			builder.WriteString(`</prosody>`)
		}
	}

	builder.WriteString(`</speak>`)
	ssmlContent := builder.String()

	// 使用更大的嵌套深度限制
	config := &ValidationConfig{
		MaxNestingDepth:      100,
		AllowUnknownElements: true,
	}

	parser := NewParser(config)
	b.SetBytes(int64(len(ssmlContent)))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := parser.Parse(ssmlContent)
		if err != nil {
			b.Fatalf("解析失败: %v", err)
		}
	}
}

// BenchmarkParseWithStrictValidation 测试严格验证模式的性能
func BenchmarkParseWithStrictValidation(b *testing.B) {
	ssmlContent := `<?xml version="1.0" encoding="UTF-8"?>
<speak version="1.0" xml:lang="zh-CN">
    <voice name="xiaoxiao" gender="female">
        <prosody rate="fast" pitch="high">
            严格验证模式下的解析测试。
            <emphasis level="strong">重要内容</emphasis>
        </prosody>
        <break time="500ms"/>
        <audio src="test.wav">音频内容</audio>
    </voice>
</speak>`

	strictConfig := &ValidationConfig{
		StrictMode:           true,
		AllowUnknownElements: false,
		MaxNestingDepth:      10,
	}

	parser := NewParser(strictConfig)
	b.SetBytes(int64(len(ssmlContent)))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := parser.Parse(ssmlContent)
		if err != nil {
			b.Fatalf("解析失败: %v", err)
		}
	}
}

// BenchmarkParseWithBasicValidation 测试基本验证模式的性能
func BenchmarkParseWithBasicValidation(b *testing.B) {
	ssmlContent := `<?xml version="1.0" encoding="UTF-8"?>
<speak version="1.0" xml:lang="zh-CN">
    <voice name="xiaoxiao" gender="female">
        <prosody rate="fast" pitch="high">
            基本验证模式下的解析测试。
            <emphasis level="strong">重要内容</emphasis>
        </prosody>
        <break time="500ms"/>
        <audio src="test.wav">音频内容</audio>
    </voice>
</speak>`

	basicConfig := &ValidationConfig{
		StrictMode:           false,
		AllowUnknownElements: true,
		MaxNestingDepth:      10,
	}

	parser := NewParser(basicConfig)
	b.SetBytes(int64(len(ssmlContent)))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := parser.Parse(ssmlContent)
		if err != nil {
			b.Fatalf("解析失败: %v", err)
		}
	}
}

// BenchmarkParseMultipleVoices 测试解析多声音 SSML 的性能
func BenchmarkParseMultipleVoices(b *testing.B) {
	var builder strings.Builder
	builder.WriteString(`<?xml version="1.0" encoding="UTF-8"?>`)
	builder.WriteString(`<speak version="1.0" xml:lang="zh-CN">`)

	voices := []string{"xiaoxiao", "xiaoyu", "xiaoyun", "xiaomeng", "ruoxi"}
	genders := []string{"female", "male"}
	rates := []string{"slow", "medium", "fast"}

	// 生成 100 个不同的声音片段
	for i := 0; i < 100; i++ {
		voice := voices[i%len(voices)]
		gender := genders[i%len(genders)]
		rate := rates[i%len(rates)]

		builder.WriteString(fmt.Sprintf(`
    <voice name="%s" gender="%s">
        <prosody rate="%s">
            这是第 %d 段声音内容。
            <emphasis level="moderate">重要信息</emphasis>
            <break time="200ms"/>
        </prosody>
    </voice>`, voice, gender, rate, i+1))
	}

	builder.WriteString(`</speak>`)
	ssmlContent := builder.String()

	parser := NewParser(nil)
	b.SetBytes(int64(len(ssmlContent)))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := parser.Parse(ssmlContent)
		if err != nil {
			b.Fatalf("解析失败: %v", err)
		}
	}
}

// BenchmarkParseAllElementTypes 测试解析包含所有元素类型的 SSML 性能
func BenchmarkParseAllElementTypes(b *testing.B) {
	ssmlContent := `<?xml version="1.0" encoding="UTF-8"?>
<speak version="1.0" xml:lang="zh-CN">
    <p>段落开始
        <s>第一个句子
            <voice name="xiaoxiao" gender="female" age="young">
                <prosody rate="fast" pitch="high" volume="loud" range="wide">
                    快速高音调大音量语音
                    <emphasis level="strong">
                        强调内容
                        <w role="verb">运行</w>
                    </emphasis>
                </prosody>
                <break time="500ms"/>
                <phoneme alphabet="ipa" ph="ˈhəloʊ">Hello</phoneme>
                <sub alias="人工智能">AI</sub>
                <audio src="beep.wav">音频备用文本</audio>
            </voice>
        </s>
        <s>第二个句子包含更多元素</s>
    </p>
    <break time="1s"/>
    <p>另一个段落</p>
</speak>`

	parser := NewParser(nil)
	b.SetBytes(int64(len(ssmlContent)))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := parser.Parse(ssmlContent)
		if err != nil {
			b.Fatalf("解析失败: %v", err)
		}
	}
}

// BenchmarkParseRealWorldSSML 测试解析真实世界 SSML 示例的性能
func BenchmarkParseRealWorldSSML(b *testing.B) {
	ssmlContent := `<?xml version="1.0" encoding="UTF-8"?>
<speak version="1.0" xml:lang="zh-CN">
    <p>欢迎使用语音合成服务。</p>
    <break time="500ms"/>
    
    <voice name="narrator" gender="neutral">
        <p>今天的天气预报：</p>
        <prosody rate="medium" pitch="medium">
            <s>上午多云，气温 18 到 22 摄氏度。</s>
            <break time="300ms"/>
            <s>下午转晴，最高温度 <emphasis level="moderate">25 摄氏度</emphasis>。</s>
            <break time="300ms"/>
            <s>建议外出时携带薄外套。</s>
        </prosody>
    </voice>
    
    <break time="1s"/>
    
    <p>新闻播报：</p>
    <voice name="news" gender="male" age="adult">
        <prosody rate="slow" pitch="low">
            <s>科技新闻：<sub alias="人工智能">AI</sub> 技术取得重大突破。</s>
            <break time="400ms"/>
            <s>研究人员开发出新的 <phoneme alphabet="ipa" ph="ˈælɡəˌrɪðəm">algorithm</phoneme>。</s>
        </prosody>
    </voice>
    
    <break time="800ms"/>
    <audio src="end-chime.wav">结束音效</audio>
    
    <p>感谢收听，再见！</p>
</speak>`

	parser := NewParser(nil)
	b.SetBytes(int64(len(ssmlContent)))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := parser.Parse(ssmlContent)
		if err != nil {
			b.Fatalf("解析失败: %v", err)
		}
	}
}

// BenchmarkParseFromReader 测试从 Reader 解析的性能
func BenchmarkParseFromReader(b *testing.B) {
	ssmlContent := `<?xml version="1.0" encoding="UTF-8"?>
<speak version="1.0" xml:lang="zh-CN">
    <voice name="xiaoxiao">
        <prosody rate="fast" pitch="high">
            从 Reader 解析的性能测试。
        </prosody>
        <break time="500ms"/>
        <emphasis level="strong">重要内容</emphasis>
    </voice>
</speak>`

	parser := NewParser(nil)
	b.SetBytes(int64(len(ssmlContent)))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		reader := strings.NewReader(ssmlContent)
		_, err := parser.ParseReader(reader)
		if err != nil {
			b.Fatalf("解析失败: %v", err)
		}
	}
}
