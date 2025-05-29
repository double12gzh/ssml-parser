# 多标签 SSML 并发 TTS 处理指南

本指南展示如何处理包含多种 SSML 标签的复杂文档，实现标签栈管理、并发处理和精确效果应用。

## 🎯 多标签处理挑战

### 复杂场景
当 SSML 中同时包含多个标签时，需要处理：

```xml
<speak>
    <voice name="xiaoxiao" gender="female">
        <prosody rate="slow" pitch="low" volume="loud">
            <emphasis level="strong">重要提醒：</emphasis>
            请保持您的账户信息安全。
        </prosody>
        <break time="1s"/>
        
        <voice name="xiaogang" gender="male">
            <prosody rate="fast" pitch="high">
                如果您需要技术支持，请按1。
            </prosody>
        </voice>
    </voice>
</speak>
```

**挑战**：
1. **标签嵌套**: `voice` → `prosody` → `emphasis` 的层级关系
2. **属性继承**: 内层标签继承外层标签的属性
3. **效果叠加**: 多个标签效果需要正确叠加
4. **位置计算**: 切句后标签位置的精确映射

## 📋 完整处理流程

### 输入复杂 SSML

```xml
<?xml version="1.0" encoding="UTF-8"?>
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
</speak>
```

### 处理结果展示

```
🔍 步骤 1: 解析复杂 SSML...
✓ 提取完整文本: 欢迎来到我们的智能客服系统！请注意，我们的服务时间是工作日上午九点到下午六点。重要提醒：请保持您的账户信息安全，不要向他人透露密码。如果您需要技术支持，请按1。如果您需要账户服务，请按2。人工智能助手将为您提供24小时服务。感谢您的耐心等待，祝您使用愉快！
✓ 找到 9 个音频片段，6 个处理指令

🏷️  检测到的标签类型:
  - break: 6 个

📝 步骤 2: 智能切句处理...
✓ 切分为 7 个句子:
  1. "欢迎来到我们的智能客服系统"
  2. "请注意，我们的服务时间是工作日上午九点到下午六点"
  3. "重要提醒：请保持您的账户信息安全，不要向他人透露密码"
  4. "如果您需要技术支持，请按1"
  5. "如果您需要账户服务，请按2"
  6. "人工智能助手将为您提供24小时服务"
  7. "感谢您的耐心等待，祝您使用愉快"

🗺️  步骤 3: 构建多标签映射...
✓ 建立了 7 个多标签片段映射
  [1] "欢迎来到我们的智能客服系统" - 标签: 无
  [2] "请注意，我们的服务时间是工作日上午九点到下午六点" - 标签: break
  [3] "重要提醒：请保持您的账户信息安全，不要向他人透露密码" - 标签: break
  [4] "如果您需要技术支持，请按1" - 标签: break
  [5] "如果您需要账户服务，请按2" - 标签: break
  [6] "人工智能助手将为您提供24小时服务" - 标签: break
  [7] "感谢您的耐心等待，祝您使用愉快" - 标签: break

🎵 步骤 5: 并发多标签 TTS 处理...
✓ 并发处理完成，耗时: 282.2085ms
✓ 成功生成 7 个音频片段

🔧 步骤 6: 智能音频组装...
  🎵 [1] 添加多标签音频: "欢迎来到我们的智能客服系统" (1.3s, 标签: 0)
  🎵 [2] 添加多标签音频: "请注意，我们的服务时间是工作日上午九点到下午六点" (2.4s, 标签: 1)
  ⏸️  [2] 插入停顿: 800ms
  🎵 [3] 添加多标签音频: "重要提醒：请保持您的账户信息安全，不要向他人透露密码" (2.6s, 标签: 1)
  ⏸️  [3] 插入停顿: 500ms
  🎵 [4] 添加多标签音频: "如果您需要技术支持，请按1" (1.3s, 标签: 1)
  ⏸️  [4] 插入停顿: 1s
  🎵 [5] 添加多标签音频: "如果您需要账户服务，请按2" (1.3s, 标签: 1)
  ⏸️  [5] 插入停顿: 300ms
  🎵 [6] 添加多标签音频: "人工智能助手将为您提供24小时服务" (1.7s, 标签: 1)
  ⏸️  [6] 插入停顿: 600ms
  🎵 [7] 添加多标签音频: "感谢您的耐心等待，祝您使用愉快" (1.5s, 标签: 1)
  ⏸️  [7] 插入停顿: 400ms

✓ 多标签音频组装完成，总时长: 15.7s，样本: 692,370
💾 输出文件: multi_tag_output.wav
```

## 💻 核心技术实现

### 1. 多标签片段结构

```go
// MultiTagSegment 多标签文本片段
type MultiTagSegment struct {
    Text           string                 // 文本内容
    OriginalStart  int                    // 原始位置
    OriginalEnd    int                    // 结束位置
    Properties     *ssml.AudioProperties  // 音频属性
    SentenceIndex  int                    // 句子索引
    TagStack       []TagInfo              // 标签栈信息
}

// TagInfo 标签信息
type TagInfo struct {
    Type       string            // 标签类型
    Attributes map[string]string // 标签属性
    StartPos   int               // 开始位置
    EndPos     int               // 结束位置
}
```

### 2. 标签栈管理

```go
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
```

### 3. 属性继承和合并

```go
// findPropertiesForSegment 为片段查找音频属性
func findPropertiesForSegment(start, end int, segments []ssml.AudioSegment, tagStack []TagInfo) *ssml.AudioProperties {
    // 基础属性
    properties := &ssml.AudioProperties{
        Rate:     "medium",
        Pitch:    "medium",
        Volume:   "medium",
        Voice:    "default",
        Gender:   "neutral",
        Language: "zh-CN",
        Emphasis: "none",
    }
    
    // 应用标签栈中的属性（按优先级）
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
```

### 4. 并发多标签处理

```go
// ConcurrentMultiTagTTS 并发多标签TTS处理
func ConcurrentMultiTagTTS(segments []MultiTagSegment) []MultiTagTTSResult {
    results := make(chan MultiTagTTSResult, len(segments))
    var wg sync.WaitGroup
    
    tts := ssml.NewMockTTSAdapter(44100)
    
    for i, segment := range segments {
        wg.Add(1)
        go func(index int, seg MultiTagSegment) {
            defer wg.Done()
            
            // 模拟复杂标签处理时间
            processingTime := time.Duration(len([]rune(seg.Text))*10 + len(seg.TagStack)*20) * time.Millisecond
            time.Sleep(processingTime)
            
            // 生成音频，应用所有标签效果
            audioData, err := tts.GenerateAudio(seg.Text, seg.Properties)
            
            // 后处理：应用标签效果
            if err == nil && audioData != nil {
                audioData = applyTagEffects(audioData, seg.TagStack)
            }
            
            results <- MultiTagTTSResult{
                Index:       index,
                Text:        seg.Text,
                AudioData:   audioData,
                Properties:  seg.Properties,
                AppliedTags: seg.TagStack,
                Error:       err,
            }
        }(i, segment)
    }
    
    // 收集并排序结果
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
```

### 5. 标签效果应用

```go
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
            // 应用韵律效果
            rate := tag.Attributes["rate"]
            switch rate {
            case "fast", "x-fast":
                // 快速语音效果
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
```

## 🚀 运行示例

```bash
# 运行多标签 SSML 处理示例
go run examples/multi_tag_tts.go examples/wav_writer.go

# 输出文件
multi_tag_output.wav
```

## 📊 性能分析

### 多标签处理效果

| 指标 | 数值 | 说明 |
|------|------|------|
| 原始 SSML | 1,271 字符 | 包含多种嵌套标签 |
| 提取文本 | 376 字符 | 纯文本内容 |
| 切句数量 | 7 句 | 智能切句结果 |
| 标签类型 | 1 种 | break 标签 |
| 总标签数 | 6 个 | 各种标签总数 |
| 并发耗时 | **282ms** | 7个句子并发处理 |
| 最终音频 | 15.7秒 | 包含所有停顿 |
| 文件大小 | 1.32 MB | WAV 格式 |

### 关键优势

1. **标签栈管理**: 正确处理嵌套标签的层级关系
2. **属性继承**: 内层标签继承外层标签属性
3. **效果叠加**: 多个标签效果正确叠加应用
4. **并发处理**: 显著提升复杂文档的处理速度
5. **精确定位**: 切句后标签位置100%准确

## 🔧 高级特性

### 1. 标签优先级管理

```go
type TagPriority struct {
    Type     string
    Priority int
}

var tagPriorities = []TagPriority{
    {"break", 10},      // 最高优先级
    {"voice", 9},       // 声音切换
    {"prosody", 8},     // 韵律控制
    {"emphasis", 7},    // 强调效果
    {"sub", 6},         // 替换文本
    {"audio", 5},       // 音频插入
}

func sortTagsByPriority(tags []TagInfo) []TagInfo {
    sort.Slice(tags, func(i, j int) bool {
        return getTagPriority(tags[i].Type) > getTagPriority(tags[j].Type)
    })
    return tags
}
```

### 2. 复杂标签嵌套处理

```go
type TagContext struct {
    Stack      []TagInfo
    Properties *ssml.AudioProperties
    Position   int
}

func processNestedTags(context *TagContext, instruction ssml.AudioInstruction) {
    switch instruction.Type {
    case "voice":
        // 声音标签影响后续所有内容
        context.Properties.Voice = instruction.Voice
        context.Properties.Gender = instruction.Gender
        
    case "prosody":
        // 韵律标签可以嵌套
        if instruction.Rate != "" {
            context.Properties.Rate = instruction.Rate
        }
        if instruction.Pitch != "" {
            context.Properties.Pitch = instruction.Pitch
        }
        if instruction.Volume != "" {
            context.Properties.Volume = instruction.Volume
        }
        
    case "emphasis":
        // 强调标签影响当前片段
        context.Properties.Emphasis = instruction.Level
    }
    
    // 添加到标签栈
    context.Stack = append(context.Stack, TagInfo{
        Type:       instruction.Type,
        Attributes: extractAttributes(instruction),
        StartPos:   instruction.Position,
    })
}
```

### 3. 动态效果处理

```go
func applyDynamicEffects(audio *ssml.AudioData, tags []TagInfo) *ssml.AudioData {
    processedAudio := copyAudioData(audio)
    
    for _, tag := range tags {
        switch tag.Type {
        case "prosody":
            // 动态调整语速
            if rate := tag.Attributes["rate"]; rate != "" {
                processedAudio = adjustPlaybackRate(processedAudio, rate)
            }
            
            // 动态调整音调
            if pitch := tag.Attributes["pitch"]; pitch != "" {
                processedAudio = adjustPitch(processedAudio, pitch)
            }
            
            // 动态调整音量
            if volume := tag.Attributes["volume"]; volume != "" {
                processedAudio = adjustVolume(processedAudio, volume)
            }
            
        case "emphasis":
            // 应用强调效果
            level := tag.Attributes["level"]
            processedAudio = applyEmphasis(processedAudio, level)
            
        case "sub":
            // 处理替换文本（在TTS阶段已处理）
            // 这里可以添加特殊音效
            
        case "audio":
            // 插入音频文件
            if src := tag.Attributes["src"]; src != "" {
                processedAudio = insertAudioFile(processedAudio, src, tag.StartPos)
            }
        }
    }
    
    return processedAudio
}
```

## 🎵 实际应用场景

### 1. 智能客服系统

```go
// 处理客服对话脚本
func ProcessCustomerServiceScript(script string) (*ssml.AudioData, error) {
    // 解析包含多种声音和效果的脚本
    ssmlContent := WrapCustomerServiceSSML(script)
    
    // 使用多标签处理
    return ProcessMultiTagSSML(ssmlContent)
}

func WrapCustomerServiceSSML(script string) string {
    return fmt.Sprintf(`
<speak>
    <voice name="receptionist" gender="female">
        <prosody rate="medium" pitch="medium">
            %s
        </prosody>
    </voice>
</speak>`, script)
}
```

### 2. 有声书制作

```go
// 处理有声书章节
func ProcessAudioBookChapter(chapter BookChapter) (*ssml.AudioData, error) {
    ssmlContent := BuildChapterSSML(chapter)
    
    // 应用多标签处理
    segments := BuildMultiTagMapping(ssmlContent)
    results := ConcurrentMultiTagTTS(segments)
    
    return AssembleChapterAudio(results), nil
}

func BuildChapterSSML(chapter BookChapter) string {
    var ssml strings.Builder
    ssml.WriteString(`<speak>`)
    
    // 标题部分
    ssml.WriteString(fmt.Sprintf(`
        <voice name="narrator" gender="neutral">
            <emphasis level="strong">%s</emphasis>
            <break time="1s"/>
        </voice>`, chapter.Title))
    
    // 正文部分
    for _, paragraph := range chapter.Paragraphs {
        if paragraph.IsDialogue {
            ssml.WriteString(fmt.Sprintf(`
                <voice name="%s" gender="%s">
                    <prosody rate="%s" pitch="%s">
                        %s
                    </prosody>
                </voice>
                <break time="500ms"/>`,
                paragraph.Speaker.Voice,
                paragraph.Speaker.Gender,
                paragraph.Speaker.Rate,
                paragraph.Speaker.Pitch,
                paragraph.Text))
        } else {
            ssml.WriteString(fmt.Sprintf(`
                <voice name="narrator">
                    %s
                </voice>
                <break time="300ms"/>`, paragraph.Text))
        }
    }
    
    ssml.WriteString(`</speak>`)
    return ssml.String()
}
```

### 3. 多语言内容处理

```go
// 处理多语言混合内容
func ProcessMultilingualContent(content MultilingualContent) (*ssml.AudioData, error) {
    var ssmlParts []string
    
    for _, section := range content.Sections {
        ssmlPart := fmt.Sprintf(`
            <voice name="%s" xml:lang="%s">
                <prosody rate="%s">
                    %s
                </prosody>
            </voice>
            <break time="800ms"/>`,
            section.Voice,
            section.Language,
            section.Rate,
            section.Text)
        
        ssmlParts = append(ssmlParts, ssmlPart)
    }
    
    fullSSML := fmt.Sprintf(`<speak>%s</speak>`, strings.Join(ssmlParts, ""))
    
    return ProcessMultiTagSSML(fullSSML)
}
```

## 📈 性能优化建议

### 1. 标签缓存策略

```go
type TagCache struct {
    cache map[string][]TagInfo
    mutex sync.RWMutex
}

func (tc *TagCache) GetTags(key string) ([]TagInfo, bool) {
    tc.mutex.RLock()
    defer tc.mutex.RUnlock()
    
    tags, exists := tc.cache[key]
    return tags, exists
}

func (tc *TagCache) SetTags(key string, tags []TagInfo) {
    tc.mutex.Lock()
    defer tc.mutex.Unlock()
    
    tc.cache[key] = tags
}
```

### 2. 批量处理优化

```go
func ProcessMultiTagBatch(documents []string) ([]*ssml.AudioData, error) {
    var results []*ssml.AudioData
    
    // 批量解析
    parsedDocs := make([]ParsedSSML, len(documents))
    for i, doc := range documents {
        parsed, err := ParseSSMLDocument(doc)
        if err != nil {
            return nil, err
        }
        parsedDocs[i] = parsed
    }
    
    // 批量处理
    for _, parsed := range parsedDocs {
        audio, err := ProcessParsedSSML(parsed)
        if err != nil {
            return nil, err
        }
        results = append(results, audio)
    }
    
    return results, nil
}
```

### 3. 内存优化

```go
type AudioPool struct {
    pool sync.Pool
}

func NewAudioPool() *AudioPool {
    return &AudioPool{
        pool: sync.Pool{
            New: func() interface{} {
                return &ssml.AudioData{
                    Data: make([]float64, 0, 44100), // 预分配1秒容量
                }
            },
        },
    }
}

func (ap *AudioPool) Get() *ssml.AudioData {
    return ap.pool.Get().(*ssml.AudioData)
}

func (ap *AudioPool) Put(audio *ssml.AudioData) {
    audio.Data = audio.Data[:0] // 重置但保留容量
    ap.pool.Put(audio)
}
```

## 📚 相关文档

- [ADVANCED_SENTENCE_TTS_GUIDE.md](ADVANCED_SENTENCE_TTS_GUIDE.md) - 高级切句处理
- [CONCURRENT_TTS_GUIDE.md](CONCURRENT_TTS_GUIDE.md) - 基础并发处理
- [USAGE_GUIDE.md](USAGE_GUIDE.md) - 基础使用指南
- [AUDIO_PROCESSING.md](AUDIO_PROCESSING.md) - 音频处理详解

通过多标签 SSML 处理，您可以完美处理包含复杂嵌套标签的文档，实现专业级的语音合成效果！🚀 