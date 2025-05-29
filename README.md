# SSML Parser - Go 语言 SSML 解析器

一个用 Go 语言实现的完整 SSML（Speech Synthesis Markup Language）解析器，支持解析、构建和序列化 SSML 文档。

## 功能特性

- ✅ **完整的 SSML 支持**: 支持所有主要的 SSML 元素
- ✅ **解析器**: 将 SSML 字符串解析为结构化数据
- ✅ **序列化器**: 将结构化数据转换回 SSML 字符串
- ✅ **构建器**: 提供链式 API 轻松构建 SSML
- ✅ **验证**: 支持基本和严格验证模式
- ✅ **错误处理**: 详细的错误信息和警告
- ✅ **类型安全**: 完全类型安全的 Go 实现
- ✅ **音频处理**: 从 SSML 到音频的完整处理流程
- ✅ **TTS 集成**: 支持 TTS 引擎集成和音频后处理

## 支持的 SSML 元素

| 元素 | 描述 | 示例 |
|------|------|------|
| `<speak>` | 根元素 | `<speak version="1.0" xml:lang="zh-CN">` |
| `<voice>` | 声音控制 | `<voice name="xiaoxiao" gender="female">` |
| `<prosody>` | 韵律控制（速度、音调、音量） | `<prosody rate="fast" pitch="high">` |
| `<emphasis>` | 强调 | `<emphasis level="strong">` |
| `<break>` | 停顿 | `<break time="1s"/>` |
| `<audio>` | 音频插入 | `<audio src="beep.wav">` |
| `<phoneme>` | 发音指导 | `<phoneme alphabet="ipa" ph="həˈloʊ">` |
| `<sub>` | 文本替换 | `<sub alias="世界贸易组织">WTO</sub>` |
| `<p>` | 段落 | `<p>这是一个段落</p>` |
| `<s>` | 句子 | `<s>这是一个句子</s>` |
| `<w>` | 单词 | `<w role="verb">run</w>` |

## 安装

```bash
go mod init your-project
go get ssml-parser
```

## 快速开始

### 1. 解析 SSML

```go
package main

import (
    "fmt"
    "log"
    "ssml-parser/ssml"
)

func main() {
    ssmlContent := `<?xml version="1.0" encoding="UTF-8"?>
    <speak version="1.0" xml:lang="zh-CN">
        <voice name="xiaoxiao">
            <prosody rate="fast" pitch="high">
                这是一个快速、高音调的语音示例。
            </prosody>
        </voice>
        <break time="1s"/>
        <emphasis level="strong">这里是重点强调的内容。</emphasis>
    </speak>`

    parser := ssml.NewParser(nil)
    result, err := parser.Parse(ssmlContent)
    if err != nil {
        log.Fatalf("解析失败: %v", err)
    }

    fmt.Printf("解析成功！版本: %s, 语言: %s\n", 
        result.Root.Version, result.Root.Lang)
}
```

### 2. 使用构建器创建 SSML

```go
package main

import (
    "fmt"
    "log"
    "ssml-parser/ssml"
)

func main() {
    builder := ssml.NewBuilder()
    
    ssmlString, err := builder.
        Version("1.0").
        Lang("zh-CN").
        Text("欢迎使用 SSML 解析器！").
        BreakTime("500ms").
        VoiceNameText("xiaoxiao", "这是小小的声音。").
        EmphasisText("strong", "这里是重点内容！").
        BuildString(true)

    if err != nil {
        log.Fatalf("构建失败: %v", err)
    }

    fmt.Println(ssmlString)
}
```

### 3. 序列化 SSML

```go
// 将解析结果序列化回 SSML 字符串
serializer := ssml.NewSerializer(true) // true 表示格式化输出
serialized, err := serializer.Serialize(result.Root)
if err != nil {
    log.Fatalf("序列化失败: %v", err)
}

fmt.Println(serialized)
```

## API 文档

### Parser（解析器）

```go
// 创建解析器
parser := ssml.NewParser(config)

// 解析 SSML 字符串
result, err := parser.Parse(ssmlContent)

// 从 io.Reader 解析
result, err := parser.ParseReader(reader)
```

### Builder（构建器）

```go
builder := ssml.NewBuilder()

// 设置基本属性
builder.Version("1.0").Lang("zh-CN")

// 添加文本
builder.Text("普通文本")

// 添加语音控制
builder.VoiceNameText("xiaoxiao", "指定声音的文本")

// 添加韵律控制
builder.RateText("fast", "快速语音")
builder.PitchText("high", "高音调语音")
builder.VolumeText("loud", "大音量语音")

// 添加停顿
builder.BreakTime("1s")        // 时间停顿
builder.BreakStrength("medium") // 强度停顿

// 添加强调
builder.EmphasisText("strong", "强调文本")

// 添加音频
builder.Audio("sound.wav", "备用文本")

// 添加替换
builder.SubText("人工智能", "AI")

// 构建结果
speak := builder.Build()           // 返回 *Speak
ssmlString, err := builder.BuildString(true) // 返回格式化字符串
```

### Serializer（序列化器）

```go
serializer := ssml.NewSerializer(pretty) // pretty: 是否格式化输出
ssmlString, err := serializer.Serialize(speak)
```

### ValidationConfig（验证配置）

```go
config := &ssml.ValidationConfig{
    StrictMode:           false,  // 严格模式
    AllowUnknownElements: true,   // 允许未知元素
    MaxNestingDepth:      10,     // 最大嵌套深度
    MaxDuration:          time.Hour, // 最大时长
}

parser := ssml.NewParser(config)
```

## 复杂示例

### 嵌套结构

```go
builder := ssml.NewBuilder()

ssml, err := builder.
    Version("1.0").
    Lang("zh-CN").
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
    BuildString(true)
```

### 发音指导

```go
builder.Sentence(func(eb *ssml.ElementBuilder) {
    eb.Text("这个单词")
    eb.PhonemeText("ipa", "ˈteknoʊlədʒi", "technology")
    eb.Text("的发音是这样的。")
})
```

## 错误处理

解析器会返回详细的错误信息和警告：

```go
result, err := parser.Parse(ssmlContent)
if err != nil {
    log.Printf("解析失败: %v", err)
    return
}

// 检查警告
if len(result.Warnings) > 0 {
    for _, warning := range result.Warnings {
        log.Printf("警告: %s", warning)
    }
}

// 检查错误
if len(result.Errors) > 0 {
    for _, error := range result.Errors {
        log.Printf("错误: %s", error)
    }
}
```

## 性能特性

- **内存效率**: 使用结构化数据避免重复解析
- **流式处理**: 支持 `io.Reader` 接口
- **并发安全**: 解析器和序列化器都是并发安全的
- **快速验证**: 可配置的验证级别

## 测试和性能

### 运行测试

```bash
# 运行所有单元测试
go test ./ssml -v

# 运行性能测试
go test ./ssml -bench=Benchmark -benchmem

# 运行特定的性能测试
go test ./ssml -bench=BenchmarkParseSimpleSSML -benchmem
```

### 性能基准

基于 Intel i7-8700K @ 3.70GHz 的测试结果：

| 场景 | 平均延迟 | 吞吐量 | 内存/操作 | 推荐用途 |
|------|----------|---------|-----------|----------|
| 简单 SSML | ~2.7µs | 37.8 MB/s | 1.2KB | 实时处理 |
| 复杂 SSML | ~24µs | 36.4 MB/s | 8KB | 常规应用 |
| 大型文档 | ~9ms | 37.4 MB/s | 2.5MB | 批量处理 |

### 性能演示

```bash
# 运行完整的音频处理演示
go run examples/audio/audio_processing_demo.go

# 运行基本使用示例
go run examples/basic/basic_usage.go

# 运行性能基准演示
go run examples/benchmark/benchmark_demo.go
```

详细的性能测试文档请参考 [BENCHMARK.md](BENCHMARK.md)。

## 音频处理功能

SSML 解析器提供完整的音频处理功能，支持从 SSML 到音频的端到端处理流程。

### 核心功能

- **文本提取**: 从 SSML 中提取纯文本用于 TTS
- **音频指令生成**: 根据 SSML 元素生成音频处理指令
- **音频后处理**: 对 TTS 生成的音频进行二次加工
- **TTS 集成**: 支持自定义 TTS 引擎集成

### 基本使用

```go
// 1. 创建音频处理器
audioProcessor := ssml.NewAudioProcessor()

// 2. 处理 SSML
result, err := audioProcessor.ProcessSSML(parseResult.Root)

// 3. 获取提取的文本和处理指令
text := result.GetTextForTTS()
instructions := result.Instructions
```

### 完整音频处理流程

```go
// 创建 TTS 适配器
ttsAdapter := ssml.NewMockTTSAdapter(44100)

// 创建完整处理器
processor := ssml.NewCompleteSSMLToAudioProcessor(ttsAdapter, 44100)

// 处理 SSML 到音频
audioData, audioResult, err := processor.ProcessSSMLToAudio(ssmlContent)
```

### 支持的音频处理

- **停顿处理**: `<break>` 元素转换为静音插入
- **音频插入**: `<audio>` 元素支持外部音频文件
- **声音变化**: `<voice>` 元素改变 TTS 属性
- **韵律控制**: `<prosody>` 元素调整语速、音调、音量
- **强调效果**: `<emphasis>` 元素通过音量调整实现

### 音频后处理功能

```go
postProcessor := ssml.NewAudioPostProcessor(44100)

// 插入静音
audioWithSilence := postProcessor.InsertSilence(audio, position, duration)

// 改变速度
fasterAudio := postProcessor.ChangeSpeed(audio, 1.5)

// 应用强调效果
emphasizedAudio := postProcessor.ApplyEmphasis(audio, startTime, duration, "strong")
```

### 运行音频处理演示

```bash
# 运行完整的音频处理演示
go run examples/audio/audio_processing_demo.go
```

详细的音频处理文档请参考 [AUDIO_PROCESSING.md](AUDIO_PROCESSING.md)。

## 贡献

欢迎提交 Issues 和 Pull Requests！

## 许可证

MIT License

## 更新日志

### v1.0.0
- 初始版本
- 完整的 SSML 解析支持
- 构建器模式 API
- 序列化功能
- 验证功能 