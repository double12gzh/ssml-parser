package ssml

import (
	"encoding/xml"
	"time"
)

// SSML 根元素
type Speak struct {
	XMLName xml.Name `xml:"speak"`
	Version string   `xml:"version,attr,omitempty"`
	Lang    string   `xml:"xml:lang,attr,omitempty"`
	Content []interface{}
}

// 文本内容
type Text struct {
	Content string
}

// 音频元素
type Audio struct {
	XMLName xml.Name `xml:"audio"`
	Src     string   `xml:"src,attr"`
	Content []interface{}
}

// 停顿元素
type Break struct {
	XMLName  xml.Name `xml:"break"`
	Time     string   `xml:"time,attr,omitempty"`
	Strength string   `xml:"strength,attr,omitempty"`
}

// 强调元素
type Emphasis struct {
	XMLName xml.Name `xml:"emphasis"`
	Level   string   `xml:"level,attr,omitempty"`
	Content []interface{}
}

// 段落元素
type Paragraph struct {
	XMLName xml.Name `xml:"p"`
	Content []interface{}
}

// 发音指导元素
type Phoneme struct {
	XMLName  xml.Name `xml:"phoneme"`
	Alphabet string   `xml:"alphabet,attr,omitempty"`
	Ph       string   `xml:"ph,attr"`
	Content  []interface{}
}

// 韵律元素（音调、速度、音量等）
type Prosody struct {
	XMLName xml.Name `xml:"prosody"`
	Rate    string   `xml:"rate,attr,omitempty"`
	Pitch   string   `xml:"pitch,attr,omitempty"`
	Range   string   `xml:"range,attr,omitempty"`
	Volume  string   `xml:"volume,attr,omitempty"`
	Content []interface{}
}

// 句子元素
type Sentence struct {
	XMLName xml.Name `xml:"s"`
	Content []interface{}
}

// 替换元素
type Sub struct {
	XMLName xml.Name `xml:"sub"`
	Alias   string   `xml:"alias,attr"`
	Content []interface{}
}

// 声音元素
type Voice struct {
	XMLName   xml.Name `xml:"voice"`
	Gender    string   `xml:"gender,attr,omitempty"`
	Age       string   `xml:"age,attr,omitempty"`
	Variant   string   `xml:"variant,attr,omitempty"`
	Name      string   `xml:"name,attr,omitempty"`
	Languages string   `xml:"xml:lang,attr,omitempty"`
	Content   []interface{}
}

// 单词元素
type W struct {
	XMLName xml.Name `xml:"w"`
	Role    string   `xml:"role,attr,omitempty"`
	Content []interface{}
}

// 解析结果结构
type ParseResult struct {
	Root     *Speak
	Warnings []string
	Errors   []string
}

// SSML 元素接口
type SSMLElement interface {
	GetContent() []interface{}
	SetContent([]interface{})
}

// 实现接口方法
func (s *Speak) GetContent() []interface{}      { return s.Content }
func (s *Speak) SetContent(c []interface{})     { s.Content = c }
func (t *Text) GetContent() []interface{}       { return nil }
func (t *Text) SetContent(c []interface{})      { /* Text 没有子元素 */ }
func (b *Break) GetContent() []interface{}      { return nil }
func (b *Break) SetContent(c []interface{})     { /* Break 没有子元素 */ }
func (a *Audio) GetContent() []interface{}      { return a.Content }
func (a *Audio) SetContent(c []interface{})     { a.Content = c }
func (e *Emphasis) GetContent() []interface{}   { return e.Content }
func (e *Emphasis) SetContent(c []interface{})  { e.Content = c }
func (p *Paragraph) GetContent() []interface{}  { return p.Content }
func (p *Paragraph) SetContent(c []interface{}) { p.Content = c }
func (p *Phoneme) GetContent() []interface{}    { return p.Content }
func (p *Phoneme) SetContent(c []interface{})   { p.Content = c }
func (p *Prosody) GetContent() []interface{}    { return p.Content }
func (p *Prosody) SetContent(c []interface{})   { p.Content = c }
func (s *Sentence) GetContent() []interface{}   { return s.Content }
func (s *Sentence) SetContent(c []interface{})  { s.Content = c }
func (s *Sub) GetContent() []interface{}        { return s.Content }
func (s *Sub) SetContent(c []interface{})       { s.Content = c }
func (v *Voice) GetContent() []interface{}      { return v.Content }
func (v *Voice) SetContent(c []interface{})     { v.Content = c }
func (w *W) GetContent() []interface{}          { return w.Content }
func (w *W) SetContent(c []interface{})         { w.Content = c }

// 验证器配置
type ValidationConfig struct {
	StrictMode           bool
	AllowUnknownElements bool
	MaxNestingDepth      int
	MaxDuration          time.Duration
}

// 默认验证配置
func DefaultValidationConfig() *ValidationConfig {
	return &ValidationConfig{
		StrictMode:           false,
		AllowUnknownElements: true,
		MaxNestingDepth:      10,
		MaxDuration:          time.Hour,
	}
}
