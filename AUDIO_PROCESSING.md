# SSML 音频处理功能

本文档介绍如何使用 SSML 解析器的音频处理功能，实现从 SSML 文本到音频的完整处理流程。

## 功能概述

SSML 音频处理模块提供以下核心功能：

1. **文本提取** - 从 SSML 中提取纯文本用于 TTS
2. **音频指令生成** - 根据 SSML 元素生成音频处理指令
3. **音频后处理** - 对 TTS 生成的音频进行二次加工
4. **完整流程** - 一站式 SSML 到音频的处理

## 核心组件

### 1. AudioProcessor - 音频处理器

负责解析 SSML 并提取文本和处理指令：

```go
processor := ssml.NewAudioProcessor()
result, err := processor.ProcessSSML(speak)
```

### 2. AudioPostProcessor - 音频后处理器

对 TTS 生成的音频进行后处理：

```go
postProcessor := ssml.NewAudioPostProcessor(44100)
finalAudio, err := postProcessor.ProcessTTSAudio(ttsAudio, instructions)
```

### 3. CompleteSSMLToAudioProcessor - 完整处理器

提供一站式的 SSML 到音频处理：

```go
processor := ssml.NewCompleteSSMLToAudioProcessor(ttsAdapter, 44100)
audioData, audioResult, err := processor.ProcessSSMLToAudio(ssmlContent)
```

## 支持的 SSML 元素

### 停顿处理 (`<break>`)

```xml
<break time="1s"/>        <!-- 1秒停顿 -->
<break time="500ms"/>     <!-- 500毫秒停顿 -->
<break strength="strong"/> <!-- 强停顿 -->
```

**处理方式**: 在 TTS 音频中插入相应时长的静音

### 音频插入 (`<audio>`)

```xml
<audio src="notification.wav">备用文本</audio>
```

**处理方式**: 在指定位置插入外部音频文件

### 声音变化 (`<voice>`)

```xml
<voice name="xiaoxiao" gender="female">
    文本内容
</voice>
```

**处理方式**: 改变 TTS 的声音属性

### 韵律控制 (`<prosody>`)

```xml
<prosody rate="fast" pitch="high" volume="loud">
    快速、高音调、大音量的语音
</prosody>
```

**处理方式**: 调整语速、音调、音量等参数

### 强调 (`<emphasis>`)

```xml
<emphasis level="strong">重点强调的内容</emphasis>
```

**处理方式**: 通过音量调整实现强调效果

### 文本替换 (`<sub>`)

```xml
<sub alias="人工智能">AI</sub>
```

**处理方式**: 用别名文本替换原文本进行 TTS

## 使用示例

### 基本使用

```go
package main

import (
    "fmt"
    "log"
    "ssml-parser/ssml"
)

func main() {
    // SSML 内容
    ssmlContent := `<?xml version="1.0" encoding="UTF-8"?>
    <speak version="1.0" xml:lang="zh-CN">
        欢迎使用语音合成。
        <break time="1s"/>
        这里是一个停顿。
    </speak>`

    // 创建处理器
    parser := ssml.NewParser(nil)
    parseResult, err := parser.Parse(ssmlContent)
    if err != nil {
        log.Fatal(err)
    }

    // 音频处理
    audioProcessor := ssml.NewAudioProcessor()
    result, err := audioProcessor.ProcessSSML(parseResult.Root)
    if err != nil {
        log.Fatal(err)
    }

    // 获取结果
    fmt.Printf("提取的文本: %s\n", result.GetTextForTTS())
    fmt.Printf("处理指令数: %d\n", len(result.Instructions))
}
```

### 完整流程

```go
func processSSMLToAudio(ssmlContent string) {
    // 创建 TTS 适配器（需要您实现）
    ttsAdapter := ssml.NewMockTTSAdapter(44100)
    
    // 创建完整处理器
    processor := ssml.NewCompleteSSMLToAudioProcessor(ttsAdapter, 44100)
    
    // 处理 SSML
    audioData, audioResult, err := processor.ProcessSSMLToAudio(ssmlContent)
    if err != nil {
        log.Fatal(err)
    }
    
    // 使用音频数据
    fmt.Printf("音频时长: %v\n", audioData.Duration)
    fmt.Printf("采样率: %d Hz\n", audioData.SampleRate)
    
    // 保存或播放音频...
}
```

## TTS 适配器接口

要集成真实的 TTS 引擎，需要实现 `TTSAdapter` 接口：

```go
type TTSAdapter interface {
    GenerateAudio(text string, properties *AudioProperties) (*AudioData, error)
}

type YourTTSAdapter struct {
    // 您的 TTS 引擎配置
}

func (adapter *YourTTSAdapter) GenerateAudio(text string, properties *AudioProperties) (*AudioData, error) {
    // 调用您的 TTS 引擎
    // 返回 AudioData 格式的音频数据
}
```

## 音频后处理功能

### 插入静音

```go
postProcessor := ssml.NewAudioPostProcessor(44100)
audioWithSilence := postProcessor.InsertSilence(audio, 2*time.Second, 1*time.Second)
```

### 改变速度

```go
fasterAudio := postProcessor.ChangeSpeed(audio, 1.5)  // 1.5倍速
slowerAudio := postProcessor.ChangeSpeed(audio, 0.75) // 0.75倍速
```

### 应用强调效果

```go
emphasizedAudio := postProcessor.ApplyEmphasis(audio, startTime, duration, "strong")
```

## 数据结构

### AudioData - 音频数据

```go
type AudioData struct {
    SampleRate int           // 采样率
    Channels   int           // 声道数
    Duration   time.Duration // 持续时间
    Data       []float64     // 音频数据（归一化浮点数）
}
```

### AudioInstruction - 音频指令

```go
type AudioInstruction struct {
    Type       string        // 指令类型：break, audio, emphasis
    Position   int           // 在文本中的位置
    Duration   time.Duration // 持续时间（用于 break）
    AudioFile  string        // 音频文件路径（用于 audio）
    Properties *AudioProperties // 音频属性变更
}
```

### AudioSegment - 音频片段

```go
type AudioSegment struct {
    Text       string        // 要合成的文本
    StartTime  time.Duration // 在最终音频中的开始时间
    Duration   time.Duration // 预期持续时间
    Properties *AudioProperties // 音频属性
}
```

## 运行示例

```bash
# 运行音频处理演示
go run examples/audio_processing_demo.go

# 运行基本使用示例
go run examples/basic_usage.go
```

## 实际应用建议

1. **集成真实 TTS 引擎** - 替换 MockTTSAdapter 为真实的 TTS 实现
2. **音频文件支持** - 实现真实的音频文件加载和插入功能
3. **格式转换** - 添加音频格式转换功能（WAV、MP3等）
4. **性能优化** - 对于大型文档，考虑流式处理
5. **错误处理** - 添加更完善的错误处理和恢复机制

## 扩展功能

- 支持更多 SSML 元素（如 `<mark>`, `<say-as>` 等）
- 添加音频效果（回声、混响等）
- 支持多声道音频处理
- 实现音频可视化功能
- 添加音频质量分析工具 