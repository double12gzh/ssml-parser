package ssml

import (
	"strings"
	"testing"
)

// TestParseBasicSSML 测试基本 SSML 解析
func TestParseBasicSSML(t *testing.T) {
	ssmlContent := `<?xml version="1.0" encoding="UTF-8"?>
<speak version="1.0" xml:lang="zh-CN">
    Hello World
</speak>`

	parser := NewParser(nil)
	result, err := parser.Parse(ssmlContent)
	if err != nil {
		t.Fatalf("解析失败: %v", err)
	}

	if result.Root == nil {
		t.Fatal("根元素为空")
	}

	if result.Root.Version != "1.0" {
		t.Errorf("期望版本 '1.0'，得到 '%s'", result.Root.Version)
	}

	if result.Root.Lang != "zh-CN" {
		t.Errorf("期望语言 'zh-CN'，得到 '%s'", result.Root.Lang)
	}

	if len(result.Root.Content) != 1 {
		t.Errorf("期望 1 个内容元素，得到 %d", len(result.Root.Content))
	}
}

// TestParseComplexSSML 测试复杂 SSML 解析
func TestParseComplexSSML(t *testing.T) {
	ssmlContent := `<?xml version="1.0" encoding="UTF-8"?>
<speak version="1.0" xml:lang="zh-CN">
    <voice name="xiaoxiao" gender="female">
        <prosody rate="fast" pitch="high" volume="loud">
            这是一个快速、高音调、大音量的语音示例。
        </prosody>
        <break time="1s"/>
        <emphasis level="strong">这里是重点强调的内容。</emphasis>
    </voice>
    <audio src="beep.wav">备用文本</audio>
    <phoneme alphabet="ipa" ph="ˈhəloʊ">Hello</phoneme>
    <sub alias="世界贸易组织">WTO</sub>
</speak>`

	parser := NewParser(nil)
	result, err := parser.Parse(ssmlContent)
	if err != nil {
		t.Fatalf("解析失败: %v", err)
	}

	if result.Root == nil {
		t.Fatal("根元素为空")
	}

	// 验证内容数量
	if len(result.Root.Content) < 4 {
		t.Errorf("期望至少 4 个内容元素，得到 %d", len(result.Root.Content))
	}

	// 验证 voice 元素
	var voiceFound bool
	for _, item := range result.Root.Content {
		if voice, ok := item.(*Voice); ok {
			voiceFound = true
			if voice.Name != "xiaoxiao" {
				t.Errorf("期望声音名称 'xiaoxiao'，得到 '%s'", voice.Name)
			}
			if voice.Gender != "female" {
				t.Errorf("期望性别 'female'，得到 '%s'", voice.Gender)
			}
		}
	}
	if !voiceFound {
		t.Error("未找到 voice 元素")
	}
}

// TestParseInvalidSSML 测试无效 SSML
func TestParseInvalidSSML(t *testing.T) {
	testCases := []struct {
		name    string
		content string
	}{
		{
			name:    "无根元素",
			content: `<?xml version="1.0"?><invalid>content</invalid>`,
		},
		{
			name:    "格式错误的XML",
			content: `<?xml version="1.0"?><speak><unclosed>`,
		},
		{
			name:    "空内容",
			content: ``,
		},
	}

	parser := NewParser(nil)
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := parser.Parse(tc.content)
			if err == nil {
				t.Error("期望解析失败，但成功了")
			}
		})
	}
}

// TestParseStrictMode 测试严格模式
func TestParseStrictMode(t *testing.T) {
	ssmlContent := `<?xml version="1.0" encoding="UTF-8"?>
<speak xml:lang="zh-CN">
    Hello World
</speak>`

	// 严格模式应该要求版本号
	strictConfig := &ValidationConfig{
		StrictMode:           true,
		AllowUnknownElements: false,
		MaxNestingDepth:      10,
	}

	parser := NewParser(strictConfig)
	_, err := parser.Parse(ssmlContent)
	if err == nil {
		t.Error("严格模式下期望解析失败（缺少版本号），但成功了")
	}
}

// TestParseUnknownElements 测试未知元素处理
func TestParseUnknownElements(t *testing.T) {
	ssmlContent := `<?xml version="1.0" encoding="UTF-8"?>
<speak version="1.0" xml:lang="zh-CN">
    <unknown>content</unknown>
    Hello World
</speak>`

	// 允许未知元素
	configAllow := &ValidationConfig{
		AllowUnknownElements: true,
		MaxNestingDepth:      10,
	}

	parser := NewParser(configAllow)
	result, err := parser.Parse(ssmlContent)
	if err != nil {
		t.Fatalf("允许未知元素时解析失败: %v", err)
	}

	if len(result.Root.Content) == 0 {
		t.Error("期望有内容元素")
	}

	// 不允许未知元素
	configDisallow := &ValidationConfig{
		AllowUnknownElements: false,
		MaxNestingDepth:      10,
	}

	parser2 := NewParser(configDisallow)
	_, err = parser2.Parse(ssmlContent)
	if err == nil {
		t.Error("不允许未知元素时期望解析失败，但成功了")
	}
}

// TestParseNestingDepth 测试嵌套深度限制
func TestParseNestingDepth(t *testing.T) {
	// 创建深度嵌套的 SSML
	var builder strings.Builder
	builder.WriteString(`<?xml version="1.0" encoding="UTF-8"?>`)
	builder.WriteString(`<speak version="1.0" xml:lang="zh-CN">`)
	
	// 创建 15 层嵌套
	for i := 0; i < 15; i++ {
		builder.WriteString(`<emphasis level="moderate">`)
	}
	builder.WriteString(`深度嵌套内容`)
	for i := 0; i < 15; i++ {
		builder.WriteString(`</emphasis>`)
	}
	builder.WriteString(`</speak>`)

	config := &ValidationConfig{
		MaxNestingDepth: 10,
	}

	parser := NewParser(config)
	_, err := parser.Parse(builder.String())
	if err == nil {
		t.Error("期望因嵌套深度超限而解析失败，但成功了")
	}
}

// TestParseAllElements 测试所有支持的元素
func TestParseAllElements(t *testing.T) {
	ssmlContent := `<?xml version="1.0" encoding="UTF-8"?>
<speak version="1.0" xml:lang="zh-CN">
    <p>段落内容</p>
    <s>句子内容</s>
    <voice name="xiaoxiao">声音内容</voice>
    <prosody rate="fast">韵律内容</prosody>
    <emphasis level="strong">强调内容</emphasis>
    <break time="1s"/>
    <audio src="test.wav">音频备用文本</audio>
    <phoneme alphabet="ipa" ph="test">发音内容</phoneme>
    <sub alias="替换">原文</sub>
    <w role="verb">单词</w>
</speak>`

	parser := NewParser(nil)
	result, err := parser.Parse(ssmlContent)
	if err != nil {
		t.Fatalf("解析所有元素失败: %v", err)
	}

	if len(result.Root.Content) != 10 {
		t.Errorf("期望 10 个元素，得到 %d", len(result.Root.Content))
	}

	// 验证各种元素类型
	elementTypes := make(map[string]bool)
	for _, item := range result.Root.Content {
		switch item.(type) {
		case *Paragraph:
			elementTypes["paragraph"] = true
		case *Sentence:
			elementTypes["sentence"] = true
		case *Voice:
			elementTypes["voice"] = true
		case *Prosody:
			elementTypes["prosody"] = true
		case *Emphasis:
			elementTypes["emphasis"] = true
		case *Break:
			elementTypes["break"] = true
		case *Audio:
			elementTypes["audio"] = true
		case *Phoneme:
			elementTypes["phoneme"] = true
		case *Sub:
			elementTypes["sub"] = true
		case *W:
			elementTypes["w"] = true
		}
	}

	expectedTypes := []string{"paragraph", "sentence", "voice", "prosody", "emphasis", "break", "audio", "phoneme", "sub", "w"}
	for _, expectedType := range expectedTypes {
		if !elementTypes[expectedType] {
			t.Errorf("未找到 %s 元素", expectedType)
		}
	}
}

// TestParseSerializeRoundtrip 测试解析-序列化往返
func TestParseSerializeRoundtrip(t *testing.T) {
	originalSSML := `<?xml version="1.0" encoding="UTF-8"?>
<speak version="1.0" xml:lang="zh-CN">
  <voice name="xiaoxiao">
    <prosody rate="fast" pitch="high">
      这是测试内容。
    </prosody>
    <break time="500ms"/>
    <emphasis level="strong">重要内容</emphasis>
  </voice>
</speak>`

	parser := NewParser(nil)
	result, err := parser.Parse(originalSSML)
	if err != nil {
		t.Fatalf("解析失败: %v", err)
	}

	serializer := NewSerializer(true)
	serializedSSML, err := serializer.Serialize(result.Root)
	if err != nil {
		t.Fatalf("序列化失败: %v", err)
	}

	// 再次解析序列化后的结果
	result2, err := parser.Parse(serializedSSML)
	if err != nil {
		t.Fatalf("再次解析失败: %v", err)
	}

	// 验证基本属性一致
	if result.Root.Version != result2.Root.Version {
		t.Errorf("版本不一致: %s vs %s", result.Root.Version, result2.Root.Version)
	}

	if result.Root.Lang != result2.Root.Lang {
		t.Errorf("语言不一致: %s vs %s", result.Root.Lang, result2.Root.Lang)
	}

	if len(result.Root.Content) != len(result2.Root.Content) {
		t.Errorf("内容元素数量不一致: %d vs %d", len(result.Root.Content), len(result2.Root.Content))
	}
} 