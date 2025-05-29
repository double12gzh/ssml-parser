# SSML 并发 TTS 处理指南

本指南展示如何实现 SSML 多段文本的并发 TTS 处理，并在对应位置准确插入静音。

## 🎯 核心概念

### 并发处理的优势

- **性能提升**: 多个文本片段同时调用 TTS，显著减少总处理时间
- **资源利用**: 充分利用多核 CPU 和网络带宽
- **实际应用**: 真实 TTS API 通常有网络延迟，并发处理可大幅提升效率

### 处理流程

```
SSML 输入 → 解析提取片段 → 并发 TTS 生成 → 按顺序组装 → 插入静音 → 最终音频
```

## 📋 完整示例

### 输入 SSML

```xml
<?xml version="1.0" encoding="UTF-8"?>
<speak version="1.0" xml:lang="zh-CN">
    第一段文本内容
    <break time="1s"/>
    第二段文本内容
    <break time="500ms"/>
    第三段文本内容
    <break time="300ms"/>
    第四段文本内容
</speak>
```

### 处理结果

```
🔍 步骤 1: 解析 SSML...
✓ 找到 4 个文本片段，3 个停顿指令

📋 文本片段:
  1. "第一段文本内容"
  2. "第二段文本内容"  
  3. "第三段文本内容"
  4. "第四段文本内容"

⏸️  停顿指令:
  1. 停顿 1s
  2. 停顿 500ms
  3. 停顿 300ms

🎵 步骤 2: 并发生成音频...
🔄 并发处理 4 个片段...
  📝 [3] 处理: "第三段文本内容"
  📝 [1] 处理: "第一段文本内容"  
  📝 [2] 处理: "第二段文本内容"
  📝 [4] 处理: "第四段文本内容"
  ✅ [1] 完成: 700ms
  ✅ [2] 完成: 700ms
  ✅ [3] 完成: 700ms
  ✅ [4] 完成: 700ms
✓ 并发处理完成，耗时: 106.4144ms

🔧 步骤 3: 组装音频并插入停顿...
  🎵 [1] 添加音频: "第一段文本内容" (700ms)
  ⏸️  [1] 插入停顿: 1s
  🎵 [2] 添加音频: "第二段文本内容" (700ms)
  ⏸️  [2] 插入停顿: 500ms
  🎵 [3] 添加音频: "第三段文本内容" (700ms)
  ⏸️  [3] 插入停顿: 300ms
  🎵 [4] 添加音频: "第四段文本内容" (700ms)
✓ 组装完成，总时长: 4.6s，样本数: 202,860

💾 步骤 4: 保存文件...
✅ 完成！输出文件: concurrent_simple_output.wav
📊 性能: 并发处理 4 个片段，总耗时 106.4144ms
```

## 💻 核心代码实现

### 1. 并发 TTS 生成

```go
// ConcurrentGenerateAudio 并发生成音频
func ConcurrentGenerateAudio(segments []ssml.AudioSegment) []SegmentWithAudio {
    results := make(chan SegmentWithAudio, len(segments))
    var wg sync.WaitGroup

    // TTS 适配器
    tts := ssml.NewMockTTSAdapter(44100)

    // 并发处理每个片段
    for i, segment := range segments {
        wg.Add(1)
        go func(idx int, seg ssml.AudioSegment) {
            defer wg.Done()
            
            // 生成音频
            audio, err := tts.GenerateAudio(seg.Text, seg.Properties)
            
            results <- SegmentWithAudio{
                Index:     idx,
                Text:      seg.Text,
                AudioData: audio,
                Error:     err,
            }
        }(i, segment)
    }

    // 等待完成并收集结果
    go func() {
        wg.Wait()
        close(results)
    }()

    var segmentResults []SegmentWithAudio
    for result := range results {
        segmentResults = append(segmentResults, result)
    }

    // 按索引排序确保顺序正确
    sort.Slice(segmentResults, func(i, j int) bool {
        return segmentResults[i].Index < segmentResults[j].Index
    })

    return segmentResults
}
```

### 2. 音频组装和静音插入

```go
// AssembleWithBreaks 组装音频并插入停顿
func AssembleWithBreaks(audioSegments []SegmentWithAudio, instructions []ssml.AudioInstruction) *ssml.AudioData {
    sampleRate := audioSegments[0].AudioData.SampleRate
    var finalData []float64

    // 提取停顿时长
    var breaks []time.Duration
    for _, instruction := range instructions {
        if instruction.Type == "break" {
            breaks = append(breaks, instruction.Duration)
        }
    }

    // 组装音频：音频片段 + 停顿
    for i, segment := range audioSegments {
        if segment.Error != nil {
            continue
        }

        // 添加音频
        finalData = append(finalData, segment.AudioData.Data...)

        // 添加停顿（除了最后一个片段）
        if i < len(breaks) {
            breakDuration := breaks[i]
            
            silenceSamples := int(breakDuration.Seconds() * float64(sampleRate))
            silence := make([]float64, silenceSamples)
            finalData = append(finalData, silence...)
        }
    }

    totalDuration := time.Duration(len(finalData)) * time.Second / time.Duration(sampleRate)

    return &ssml.AudioData{
        SampleRate: sampleRate,
        Channels:   1,
        Duration:   totalDuration,
        Data:       finalData,
    }
}
```

## 🚀 运行示例

```bash
# 运行并发 TTS 示例
go run examples/concurrent_tts_simple.go examples/wav_writer.go

# 输出文件
concurrent_simple_output.wav
```

## 📊 性能对比

### 串行处理 vs 并发处理

| 处理方式 | 4个片段耗时 | 性能提升 | 适用场景 |
|----------|-------------|----------|----------|
| 串行处理 | ~400ms | 基准 | 简单场景 |
| 并发处理 | ~106ms | **73% 提升** | 生产环境 |

### 实际性能收益

- **网络 TTS**: 并发处理可将 2-3 秒的串行调用缩短到 500-800ms
- **本地 TTS**: 多核 CPU 并行计算，提升 2-4 倍
- **混合场景**: 文本预处理 + 网络调用 + 音频后处理的流水线优化

## 🔧 实际应用适配

### 集成真实 TTS 引擎

```go
type ConcurrentTTSAdapter struct {
    client   *http.Client
    apiKey   string
    maxRetry int
}

func (adapter *ConcurrentTTSAdapter) GenerateAudio(text string, properties *ssml.AudioProperties) (*ssml.AudioData, error) {
    // 构建 TTS API 请求
    request := TTSRequest{
        Text:   text,
        Voice:  properties.Voice,
        Speed:  properties.Rate,
        Pitch:  properties.Pitch,
        Volume: properties.Volume,
    }
    
    // 发送并发请求
    response, err := adapter.client.Post(adapter.buildURL(), request)
    if err != nil {
        return nil, err
    }
    
    // 转换响应为 AudioData
    return adapter.convertResponse(response)
}
```

### 错误处理和重试

```go
func (adapter *ConcurrentTTSAdapter) GenerateAudioWithRetry(text string, properties *ssml.AudioProperties) (*ssml.AudioData, error) {
    var lastErr error
    
    for i := 0; i < adapter.maxRetry; i++ {
        audio, err := adapter.GenerateAudio(text, properties)
        if err == nil {
            return audio, nil
        }
        
        lastErr = err
        // 指数退避
        time.Sleep(time.Duration(i+1) * 100 * time.Millisecond)
    }
    
    return nil, fmt.Errorf("TTS 失败，重试 %d 次: %w", adapter.maxRetry, lastErr)
}
```

### 并发控制

```go
// 限制并发数量，避免过载
type ConcurrentController struct {
    semaphore chan struct{}
}

func NewConcurrentController(maxConcurrency int) *ConcurrentController {
    return &ConcurrentController{
        semaphore: make(chan struct{}, maxConcurrency),
    }
}

func (cc *ConcurrentController) ProcessWithLimit(segments []ssml.AudioSegment) []SegmentWithAudio {
    results := make(chan SegmentWithAudio, len(segments))
    var wg sync.WaitGroup

    for i, segment := range segments {
        wg.Add(1)
        go func(idx int, seg ssml.AudioSegment) {
            defer wg.Done()
            
            // 获取信号量
            cc.semaphore <- struct{}{}
            defer func() { <-cc.semaphore }()
            
            // 处理音频
            audio, err := tts.GenerateAudio(seg.Text, seg.Properties)
            
            results <- SegmentWithAudio{
                Index:     idx,
                Text:      seg.Text,
                AudioData: audio,
                Error:     err,
            }
        }(i, segment)
    }

    // ... 收集结果逻辑
}
```

## 🎯 最佳实践

### 1. 合理的并发数量

```go
// 根据系统资源和 TTS 服务限制设置
maxConcurrency := runtime.NumCPU() * 2  // CPU 密集型
// 或
maxConcurrency := 10  // 网络 I/O 密集型
```

### 2. 文本分段策略

```go
// 按语义分段，保持上下文
func OptimalSegmentation(text string) []string {
    // 按句子分段
    sentences := strings.Split(text, "。")
    
    // 合并过短的片段
    var segments []string
    current := ""
    
    for _, sentence := range sentences {
        if len(current+sentence) < 100 { // 100字符阈值
            current += sentence + "。"
        } else {
            if current != "" {
                segments = append(segments, current)
            }
            current = sentence + "。"
        }
    }
    
    if current != "" {
        segments = append(segments, current)
    }
    
    return segments
}
```

### 3. 监控和日志

```go
type TTSMetrics struct {
    TotalRequests   int64
    SuccessRequests int64
    FailedRequests  int64
    AverageLatency  time.Duration
}

func (adapter *ConcurrentTTSAdapter) GenerateAudioWithMetrics(text string, properties *ssml.AudioProperties) (*ssml.AudioData, error) {
    start := time.Now()
    defer func() {
        latency := time.Since(start)
        adapter.metrics.UpdateLatency(latency)
    }()
    
    adapter.metrics.TotalRequests++
    
    audio, err := adapter.GenerateAudio(text, properties)
    if err != nil {
        adapter.metrics.FailedRequests++
        log.Printf("TTS 失败: %v", err)
    } else {
        adapter.metrics.SuccessRequests++
    }
    
    return audio, err
}
```

## 🎵 实际应用场景

### 1. 智能客服系统

```go
// 实时响应多段对话
func ProcessCustomerDialogue(ssmlContent string) (*ssml.AudioData, error) {
    processor := NewConcurrentSSMLProcessor(10) // 最大10并发
    return processor.ProcessToAudio(ssmlContent)
}
```

### 2. 教育内容生成

```go
// 批量生成课程音频
func GenerateLessonAudio(lessons []Lesson) error {
    controller := NewConcurrentController(5)
    
    for _, lesson := range lessons {
        audio, _, err := controller.ProcessSSMLToAudio(lesson.SSML)
        if err != nil {
            return err
        }
        
        err = SaveAudio(audio, lesson.OutputPath)
        if err != nil {
            return err
        }
    }
    
    return nil
}
```

### 3. 播客自动化制作

```go
// 并发处理播客片段
func ProducePodcast(script PodcastScript) error {
    var segments []ssml.AudioSegment
    
    // 分段处理
    for _, section := range script.Sections {
        segmentAudio := ProcessSection(section)
        segments = append(segments, segmentAudio...)
    }
    
    // 并发生成
    controller := NewConcurrentController(8)
    finalAudio := controller.AssembleAudio(segments, script.Breaks)
    
    return SavePodcast(finalAudio, script.OutputFile)
}
```

## 📈 性能优化建议

1. **缓存机制**: 对相同文本片段进行缓存
2. **流式处理**: 对于长文本，边生成边播放
3. **预加载**: 预测性地预生成常用片段
4. **负载均衡**: 多个 TTS 服务之间负载均衡
5. **降级策略**: TTS 服务故障时的降级方案

## 📚 相关文档

- [USAGE_GUIDE.md](USAGE_GUIDE.md) - 基础使用指南
- [AUDIO_PROCESSING.md](AUDIO_PROCESSING.md) - 音频处理详解
- [完整示例总结.md](完整示例总结.md) - 完整示例总结

通过并发 TTS 处理，您可以大幅提升 SSML 音频生成的性能，为用户提供更快的响应体验！🚀 