# SSML 到音频文件 - 完整示例

本文档展示了如何使用 SSML 解析器从 SSML 文本生成最终的音频文件。

## 🚀 快速开始

### 最简单的使用方法

```bash
# 运行快速开始示例
go run examples/quick_start.go examples/wav_writer.go
```

**输出结果**：
```
🎵 SSML 快速开始示例
==================
📝 SSML 内容:
<?xml version="1.0" encoding="UTF-8"?>
<speak version="1.0" xml:lang="zh-CN">
    你好，世界！
    <break time="1s"/>
    这是一个简单的 SSML 示例。
    <emphasis level="strong">重点内容！</emphasis>
</speak>

🔄 正在处理...
✅ 处理完成！
📝 提取文本: 你好，世界！这是一个简单的 SSML 示例。重点内容！
⏱️  音频时长: 3.7s
🎵 采样率: 44100 Hz
📊 数据点数: 163,170
💾 音频已保存: quick_start_output.wav
🎧 可以用音频播放器播放该文件
```

### 完整功能演示

```bash
# 运行完整功能示例
go run examples/complete_example.go examples/wav_writer.go
```

**输出结果**：
```
SSML 到音频文件完整示例
=========================

步骤 1: 解析 SSML...
✓ SSML 解析成功，版本: 1.0, 语言: zh-CN

步骤 2: 提取文本和音频处理指令...
✓ 提取的纯文本: 欢迎使用 SSML 音频处理系统！这是快速、高音调、大音量的语音。这里是重点强调的内容。人工智能技术正在快速发展。感谢您的使用！
✓ 音频片段数: 6
✓ 处理指令数: 4

步骤 3: 使用 TTS 生成基础音频...
✓ TTS 音频生成成功
  - 采样率: 44100 Hz
  - 原始时长: 12.8s

步骤 4: 应用音频后处理...
✓ 音频后处理完成
  - 最终时长: 15.3s
  - 时长变化: 12.8s → 15.3s (增加了 2.5s)

步骤 5: 保存音频文件...
✓ 音频文件保存成功: output_ssml_audio.wav

步骤 6: 验证音频文件...
✓ WAV 文件格式验证通过: output_ssml_audio.wav

文件信息:
  - 文件名: output_ssml_audio.wav
  - 文件大小: 1349504 字节 (1317.88 KB)
```

## 📊 核心功能演示

### 1. SSML 解析
- ✅ 支持所有主要 SSML 元素
- ✅ 提取纯文本用于 TTS
- ✅ 生成音频处理指令

### 2. 音频处理
- ✅ Mock TTS 生成模拟音频
- ✅ 根据 SSML 插入停顿（静音）
- ✅ 支持不同语音属性

### 3. 文件输出
- ✅ 生成标准 WAV 文件
- ✅ 16 位 PCM，44.1kHz 采样率
- ✅ 文件格式验证

## 🎯 支持的 SSML 元素

| 元素 | 功能 | 示例效果 |
|------|------|----------|
| `<break time="1s"/>` | 插入 1 秒静音 | 实际停顿 1 秒 |
| `<voice name="xiaoxiao">` | 改变声音 | 使用指定声音模型 |
| `<prosody rate="fast">` | 语速控制 | 影响 TTS 生成 |
| `<prosody pitch="high">` | 音调控制 | 调整音频频率 |
| `<prosody volume="loud">` | 音量控制 | 调整音频幅度 |
| `<emphasis level="strong">` | 强调效果 | 音量增强 30% |
| `<sub alias="人工智能">AI</sub>` | 文本替换 | "AI" → "人工智能" |

## 🔄 完整处理流程

```
1. SSML 文本输入
   ↓
2. XML 解析 (Parser)
   ↓
3. 文本提取 + 指令生成 (AudioProcessor)
   ↓
4. TTS 音频生成 (MockTTSAdapter)
   ↓
5. 音频后处理 (AudioPostProcessor)
   ↓
6. WAV 文件保存 (WAVWriter)
   ↓
7. 本地音频文件
```

## 💻 核心代码

### 最简单的使用方式

```go
// 一行代码完成 SSML 到音频转换
audioData, audioResult, err := ProcessSSMLToAudio(ssmlContent)

// 保存到文件
err = SaveAudioToFile(audioData, "output.wav")
```

### 完整的流程控制

```go
// 1. 解析 SSML
parser := ssml.NewParser(nil)
parseResult, err := parser.Parse(ssmlContent)

// 2. 提取文本和指令
audioProcessor := ssml.NewAudioProcessor()
result, err := audioProcessor.ProcessSSML(parseResult.Root)

// 3. TTS 生成音频
mockTTS := NewEnhancedMockTTS(44100)
ttsAudio, err := mockTTS.GenerateAudio(result.GetTextForTTS(), properties)

// 4. 音频后处理
postProcessor := ssml.NewAudioPostProcessor(ttsAudio.SampleRate)
finalAudio, err := postProcessor.ProcessTTSAudio(ttsAudio, result.Instructions)

// 5. 保存文件
err = SaveToWAV(finalAudio, "output.wav")
```

## 📁 生成的文件

运行示例后会生成以下文件：

1. **quick_start_output.wav** - 快速示例生成的音频文件
2. **output_ssml_audio.wav** - 完整示例生成的音频文件

### 文件特性
- **格式**: 16 位 PCM WAV
- **采样率**: 44,100 Hz (CD 质量)
- **兼容性**: 可用 Windows Media Player、VLC、Audacity 播放

## 🔧 实际应用

### 替换为真实 TTS 引擎

```go
type YourTTSAdapter struct {
    apiKey string
    baseURL string
}

func (adapter *YourTTSAdapter) GenerateAudio(text string, properties *ssml.AudioProperties) (*ssml.AudioData, error) {
    // 调用真实的 TTS API
    response := callTTSAPI(adapter.apiKey, text, properties)
    
    // 转换为标准格式
    return convertToAudioData(response), nil
}
```

### 集成到应用中

```go
// Web 服务示例
func handleTTSRequest(w http.ResponseWriter, r *http.Request) {
    ssmlContent := r.FormValue("ssml")
    
    // 处理 SSML
    audioData, _, err := ProcessSSMLToAudio(ssmlContent)
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
    
    // 返回音频文件
    w.Header().Set("Content-Type", "audio/wav")
    ServeAudio(w, audioData)
}
```

## 📚 相关文档

- [USAGE_GUIDE.md](USAGE_GUIDE.md) - 详细使用指南
- [AUDIO_PROCESSING.md](AUDIO_PROCESSING.md) - 音频处理文档
- [README.md](README.md) - 项目主文档

## 🎯 总结

这个完整的示例展示了：

✅ **SSML 解析**: 完整支持 SSML 标准  
✅ **文本提取**: 智能提取 TTS 所需文本  
✅ **音频生成**: 模拟 TTS 引擎生成音频  
✅ **音频处理**: 根据 SSML 指令处理音频  
✅ **文件保存**: 生成标准 WAV 音频文件  
✅ **易于集成**: 可轻松替换为真实 TTS 引擎  

完整的从 SSML 到音频文件的处理系统！🎵 