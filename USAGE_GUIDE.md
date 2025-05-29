# SSML 到音频文件 - 完整使用指南

本指南将演示如何使用 SSML 解析器将 SSML 文本转换为最终的音频文件。

## 🎯 完整流程概览

```
SSML 文本 → 解析 → 提取文本和指令 → TTS 生成音频 → 音频后处理 → WAV 文件
```

## 📋 步骤详解

### 步骤 1: 准备 SSML 内容

```xml
<?xml version="1.0" encoding="UTF-8"?>
<speak version="1.0" xml:lang="zh-CN">
    <voice name="xiaoxiao" gender="female">
        欢迎使用 SSML 音频处理系统！
        <break time="1s"/>
        
        <prosody rate="fast" pitch="high" volume="loud">
            这是快速、高音调、大音量的语音。
        </prosody>
        
        <break time="500ms"/>
        
        <emphasis level="strong">这里是重点强调的内容。</emphasis>
        
        <break time="800ms"/>
        
        <sub alias="人工智能">AI</sub>技术正在快速发展。
        
        <break time="200ms"/>
        
        感谢您的使用！
    </voice>
</speak>
```

### 步骤 2: 解析 SSML

```go
// 创建解析器
parser := ssml.NewParser(nil)

// 解析 SSML 字符串
parseResult, err := parser.Parse(ssmlContent)
if err != nil {
    log.Fatalf("SSML 解析失败: %v", err)
}
```

**输出**：
- ✅ SSML 解析成功，版本: 1.0, 语言: zh-CN

### 步骤 3: 提取文本和处理指令

```go
// 创建音频处理器
audioProcessor := ssml.NewAudioProcessor()

// 处理 SSML，提取文本和指令
result, err := audioProcessor.ProcessSSML(parseResult.Root)
if err != nil {
    log.Fatalf("音频处理失败: %v", err)
}

// 获取提取的文本
text := result.GetTextForTTS()
instructions := result.Instructions
```

**输出**：
- 📝 **提取的纯文本**: "欢迎使用 SSML 音频处理系统！这是快速、高音调、大音量的语音。这里是重点强调的内容。人工智能技术正在快速发展。感谢您的使用！"
- 🎵 **音频片段数**: 6
- 📋 **处理指令数**: 4

**音频片段详情**：
```
1. [0s-2.55s] "欢迎使用 SSML 音频处理系统！" (声音=xiaoxiao, 语速=medium, 音调=medium, 音量=medium)
2. [3.55s-5.35s] "这是快速、高音调、大音量的语音。" (声音=xiaoxiao, 语速=fast, 音调=high, 音量=loud)
3. [5.85s-7.5s] "这里是重点强调的内容。" (声音=xiaoxiao, 语速=medium, 音调=medium, 音量=medium)
4. [8.3s-8.9s] "人工智能" (声音=xiaoxiao, 语速=medium, 音调=medium, 音量=medium)
5. [8.9s-10.25s] "技术正在快速发展。" (声音=xiaoxiao, 语速=medium, 音调=medium, 音量=medium)
6. [10.45s-11.5s] "感谢您的使用！" (声音=xiaoxiao, 语速=medium, 音调=medium, 音量=medium)
```

**处理指令详情**：
```
1. 停顿 1s (位置: 39)
2. 停顿 500ms (位置: 87)
3. 停顿 800ms (位置: 120)
4. 停顿 200ms (位置: 159)
```

### 步骤 4: 使用 TTS 生成基础音频

```go
// 创建 TTS 适配器（这里使用模拟 TTS）
mockTTS := NewEnhancedMockTTS(44100)

// 获取音频属性
var properties *ssml.AudioProperties
if len(result.Segments) > 0 {
    properties = result.Segments[0].Properties
}

// 生成音频
ttsAudio, err := mockTTS.GenerateAudio(result.GetTextForTTS(), properties)
if err != nil {
    log.Fatalf("TTS 生成失败: %v", err)
}
```

**输出**：
- ✅ TTS 音频生成成功
- 🎵 **采样率**: 44100 Hz
- 📻 **声道数**: 1
- ⏱️ **原始时长**: 12.8s
- 📊 **数据点数**: 564,480

### 步骤 5: 应用音频后处理

```go
// 创建音频后处理器
postProcessor := ssml.NewAudioPostProcessor(ttsAudio.SampleRate)

// 应用处理指令（插入静音等）
finalAudio, err := postProcessor.ProcessTTSAudio(ttsAudio, result.Instructions)
if err != nil {
    log.Fatalf("音频后处理失败: %v", err)
}
```

**输出**：
- ✅ 音频后处理完成
- ⏱️ **最终时长**: 15.3s
- 📊 **最终数据点数**: 674,730
- 📈 **时长变化**: 12.8s → 15.3s (增加了 2.5s 的停顿)

### 步骤 6: 保存为 WAV 文件

```go
// 保存为 WAV 文件
filename := "output_ssml_audio.wav"
err = SaveToWAV(finalAudio, filename)
if err != nil {
    log.Fatalf("保存音频文件失败: %v", err)
}
```

**输出**：
- ✅ WAV 文件保存成功，包含 674,730 个样本
- 📁 **文件大小**: 1,349,504 字节 (1.32 MB)
- 🎧 **格式**: 16 位 PCM WAV 文件

## 🎧 音频文件特性

生成的 WAV 文件具有以下特性：

- **格式**: 16 位 PCM WAV
- **采样率**: 44,100 Hz（CD 质量）
- **声道**: 单声道
- **时长**: 15.3 秒
- **包含效果**:
  - 1 秒停顿（在 "欢迎使用 SSML 音频处理系统！" 后）
  - 500 毫秒停顿
  - 800 毫秒停顿  
  - 200 毫秒停顿
  - 不同的语音属性（音调、语速、音量）

## 🚀 运行示例

```bash
# 运行完整示例
go run examples/complete_example.go examples/wav_writer.go

# 输出文件
output_ssml_audio.wav
```

## 🎵 播放音频

您可以使用以下软件播放生成的 WAV 文件：

- **Windows**: Windows Media Player
- **跨平台**: VLC Media Player
- **专业音频**: Audacity
- **命令行**: `ffplay output_ssml_audio.wav`

## 🔧 实际应用适配

### 集成真实 TTS 引擎

要使用真实的 TTS 引擎，只需实现 `TTSAdapter` 接口：

```go
type YourTTSAdapter struct {
    // 您的 TTS 引擎配置
}

func (adapter *YourTTSAdapter) GenerateAudio(text string, properties *ssml.AudioProperties) (*ssml.AudioData, error) {
    // 调用您的 TTS 引擎
    // 返回 AudioData 格式的音频数据
    
    // 示例：调用百度、阿里云、腾讯云等 TTS API
    rawAudio := callYourTTSAPI(text, properties)
    
    return &ssml.AudioData{
        SampleRate: 16000,
        Channels:   1,
        Duration:   calculateDuration(rawAudio),
        Data:       convertToFloat64(rawAudio),
    }, nil
}
```

### 支持的 SSML 属性映射

| SSML 属性 | 处理方式 | 示例 |
|-----------|----------|------|
| `<break time="1s"/>` | 插入静音 | 1 秒无声音频 |
| `<voice name="xiaoxiao">` | 改变 TTS 声音 | 使用指定声音模型 |
| `<prosody rate="fast">` | 调整语速 | TTS 参数调整 |
| `<prosody pitch="high">` | 调整音调 | TTS 参数调整 |
| `<prosody volume="loud">` | 调整音量 | 音频增益处理 |
| `<emphasis level="strong">` | 强调效果 | 音量增强 +30% |
| `<sub alias="人工智能">AI</sub>` | 文本替换 | "AI" → "人工智能" |

## 📊 性能数据

基于示例运行的性能数据：

- **文本长度**: 64 个字符
- **处理时间**: < 1 秒
- **内存使用**: 约 2 MB
- **输出文件**: 1.32 MB
- **音频质量**: CD 级别 (44.1kHz/16bit)

## 🎯 应用场景

这个完整的 SSML 到音频转换系统可以用于：

1. **语音合成应用** - 智能客服、语音助手
2. **教育软件** - 语音教学内容生成
3. **播客制作** - 自动化音频内容生成
4. **无障碍应用** - 文本转语音服务
5. **游戏开发** - 动态语音内容生成
6. **智能硬件** - IoT 设备语音反馈

## 🔧 进一步定制

您可以根据需要进一步定制：

- **音频格式**: 支持 MP3、OGG 等格式输出
- **多声道**: 支持立体声、5.1 环绕声
- **音频效果**: 添加混响、回声、均衡器
- **实时处理**: 支持流式音频生成
- **批量处理**: 支持多文件批量转换

## 📚 相关文档

- [AUDIO_PROCESSING.md](AUDIO_PROCESSING.md) - 详细的音频处理文档
- [BENCHMARK.md](BENCHMARK.md) - 性能测试文档  
- [README.md](README.md) - 项目主文档 