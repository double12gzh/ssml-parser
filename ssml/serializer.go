package ssml

import (
	"fmt"
	"html"
	"strings"
)

// Serializer SSML 序列化器
type Serializer struct {
	Pretty bool
	Indent string
}

// NewSerializer 创建新的序列化器
func NewSerializer(pretty bool) *Serializer {
	indent := ""
	if pretty {
		indent = "  "
	}
	return &Serializer{
		Pretty: pretty,
		Indent: indent,
	}
}

// escapeString 转义字符串中的 XML 特殊字符
func (s *Serializer) escapeString(str string) string {
	return html.EscapeString(str)
}

// Serialize 将 Speak 结构体序列化为 SSML 字符串
func (s *Serializer) Serialize(speak *Speak) (string, error) {
	if speak == nil {
		return "", fmt.Errorf("speak is nil")
	}

	var builder strings.Builder
	
	// 写入 XML 声明
	builder.WriteString(`<?xml version="1.0" encoding="UTF-8"?>`)
	if s.Pretty {
		builder.WriteString("\n")
	}

	// 序列化 speak 元素
	if err := s.serializeSpeak(&builder, speak, 0); err != nil {
		return "", err
	}

	return builder.String(), nil
}

// serializeSpeak 序列化 speak 元素
func (s *Serializer) serializeSpeak(builder *strings.Builder, speak *Speak, depth int) error {
	s.writeIndent(builder, depth)
	builder.WriteString("<speak")

	// 写入属性
	if speak.Version != "" {
		builder.WriteString(fmt.Sprintf(` version="%s"`, s.escapeString(speak.Version)))
	}
	if speak.Lang != "" {
		builder.WriteString(fmt.Sprintf(` xml:lang="%s"`, s.escapeString(speak.Lang)))
	}

	if len(speak.Content) == 0 {
		builder.WriteString("/>")
		if s.Pretty {
			builder.WriteString("\n")
		}
		return nil
	}

	builder.WriteString(">")
	if s.Pretty {
		builder.WriteString("\n")
	}

	// 序列化内容
	if err := s.serializeContent(builder, speak.Content, depth+1); err != nil {
		return err
	}

	s.writeIndent(builder, depth)
	builder.WriteString("</speak>")
	if s.Pretty {
		builder.WriteString("\n")
	}

	return nil
}

// serializeContent 序列化内容
func (s *Serializer) serializeContent(builder *strings.Builder, content []interface{}, depth int) error {
	for _, item := range content {
		switch v := item.(type) {
		case Text:
			s.writeIndent(builder, depth)
			builder.WriteString(s.escapeString(v.Content))
			if s.Pretty {
				builder.WriteString("\n")
			}

		case *Audio:
			if err := s.serializeAudio(builder, v, depth); err != nil {
				return err
			}

		case *Break:
			if err := s.serializeBreak(builder, v, depth); err != nil {
				return err
			}

		case *Emphasis:
			if err := s.serializeEmphasis(builder, v, depth); err != nil {
				return err
			}

		case *Paragraph:
			if err := s.serializeParagraph(builder, v, depth); err != nil {
				return err
			}

		case *Phoneme:
			if err := s.serializePhoneme(builder, v, depth); err != nil {
				return err
			}

		case *Prosody:
			if err := s.serializeProsody(builder, v, depth); err != nil {
				return err
			}

		case *Sentence:
			if err := s.serializeSentence(builder, v, depth); err != nil {
				return err
			}

		case *Sub:
			if err := s.serializeSub(builder, v, depth); err != nil {
				return err
			}

		case *Voice:
			if err := s.serializeVoice(builder, v, depth); err != nil {
				return err
			}

		case *W:
			if err := s.serializeW(builder, v, depth); err != nil {
				return err
			}

		default:
			return fmt.Errorf("unknown content type: %T", v)
		}
	}
	return nil
}

// serializeAudio 序列化 audio 元素
func (s *Serializer) serializeAudio(builder *strings.Builder, audio *Audio, depth int) error {
	s.writeIndent(builder, depth)
	builder.WriteString("<audio")

	if audio.Src != "" {
		builder.WriteString(fmt.Sprintf(` src="%s"`, s.escapeString(audio.Src)))
	}

	if len(audio.Content) == 0 {
		builder.WriteString("/>")
		if s.Pretty {
			builder.WriteString("\n")
		}
		return nil
	}

	builder.WriteString(">")
	if s.Pretty {
		builder.WriteString("\n")
	}

	if err := s.serializeContent(builder, audio.Content, depth+1); err != nil {
		return err
	}

	s.writeIndent(builder, depth)
	builder.WriteString("</audio>")
	if s.Pretty {
		builder.WriteString("\n")
	}

	return nil
}

// serializeBreak 序列化 break 元素
func (s *Serializer) serializeBreak(builder *strings.Builder, breakElem *Break, depth int) error {
	s.writeIndent(builder, depth)
	builder.WriteString("<break")

	if breakElem.Time != "" {
		builder.WriteString(fmt.Sprintf(` time="%s"`, s.escapeString(breakElem.Time)))
	}
	if breakElem.Strength != "" {
		builder.WriteString(fmt.Sprintf(` strength="%s"`, s.escapeString(breakElem.Strength)))
	}

	builder.WriteString("/>")
	if s.Pretty {
		builder.WriteString("\n")
	}

	return nil
}

// serializeEmphasis 序列化 emphasis 元素
func (s *Serializer) serializeEmphasis(builder *strings.Builder, emphasis *Emphasis, depth int) error {
	s.writeIndent(builder, depth)
	builder.WriteString("<emphasis")

	if emphasis.Level != "" {
		builder.WriteString(fmt.Sprintf(` level="%s"`, s.escapeString(emphasis.Level)))
	}

	if len(emphasis.Content) == 0 {
		builder.WriteString("/>")
		if s.Pretty {
			builder.WriteString("\n")
		}
		return nil
	}

	builder.WriteString(">")
	if s.Pretty {
		builder.WriteString("\n")
	}

	if err := s.serializeContent(builder, emphasis.Content, depth+1); err != nil {
		return err
	}

	s.writeIndent(builder, depth)
	builder.WriteString("</emphasis>")
	if s.Pretty {
		builder.WriteString("\n")
	}

	return nil
}

// serializeParagraph 序列化 p 元素
func (s *Serializer) serializeParagraph(builder *strings.Builder, paragraph *Paragraph, depth int) error {
	s.writeIndent(builder, depth)
	builder.WriteString("<p>")
	if s.Pretty {
		builder.WriteString("\n")
	}

	if err := s.serializeContent(builder, paragraph.Content, depth+1); err != nil {
		return err
	}

	s.writeIndent(builder, depth)
	builder.WriteString("</p>")
	if s.Pretty {
		builder.WriteString("\n")
	}

	return nil
}

// serializePhoneme 序列化 phoneme 元素
func (s *Serializer) serializePhoneme(builder *strings.Builder, phoneme *Phoneme, depth int) error {
	s.writeIndent(builder, depth)
	builder.WriteString("<phoneme")

	if phoneme.Alphabet != "" {
		builder.WriteString(fmt.Sprintf(` alphabet="%s"`, s.escapeString(phoneme.Alphabet)))
	}
	if phoneme.Ph != "" {
		builder.WriteString(fmt.Sprintf(` ph="%s"`, s.escapeString(phoneme.Ph)))
	}

	if len(phoneme.Content) == 0 {
		builder.WriteString("/>")
		if s.Pretty {
			builder.WriteString("\n")
		}
		return nil
	}

	builder.WriteString(">")
	if s.Pretty {
		builder.WriteString("\n")
	}

	if err := s.serializeContent(builder, phoneme.Content, depth+1); err != nil {
		return err
	}

	s.writeIndent(builder, depth)
	builder.WriteString("</phoneme>")
	if s.Pretty {
		builder.WriteString("\n")
	}

	return nil
}

// serializeProsody 序列化 prosody 元素
func (s *Serializer) serializeProsody(builder *strings.Builder, prosody *Prosody, depth int) error {
	s.writeIndent(builder, depth)
	builder.WriteString("<prosody")

	if prosody.Rate != "" {
		builder.WriteString(fmt.Sprintf(` rate="%s"`, s.escapeString(prosody.Rate)))
	}
	if prosody.Pitch != "" {
		builder.WriteString(fmt.Sprintf(` pitch="%s"`, s.escapeString(prosody.Pitch)))
	}
	if prosody.Range != "" {
		builder.WriteString(fmt.Sprintf(` range="%s"`, s.escapeString(prosody.Range)))
	}
	if prosody.Volume != "" {
		builder.WriteString(fmt.Sprintf(` volume="%s"`, s.escapeString(prosody.Volume)))
	}

	if len(prosody.Content) == 0 {
		builder.WriteString("/>")
		if s.Pretty {
			builder.WriteString("\n")
		}
		return nil
	}

	builder.WriteString(">")
	if s.Pretty {
		builder.WriteString("\n")
	}

	if err := s.serializeContent(builder, prosody.Content, depth+1); err != nil {
		return err
	}

	s.writeIndent(builder, depth)
	builder.WriteString("</prosody>")
	if s.Pretty {
		builder.WriteString("\n")
	}

	return nil
}

// serializeSentence 序列化 s 元素
func (s *Serializer) serializeSentence(builder *strings.Builder, sentence *Sentence, depth int) error {
	s.writeIndent(builder, depth)
	builder.WriteString("<s>")
	if s.Pretty {
		builder.WriteString("\n")
	}

	if err := s.serializeContent(builder, sentence.Content, depth+1); err != nil {
		return err
	}

	s.writeIndent(builder, depth)
	builder.WriteString("</s>")
	if s.Pretty {
		builder.WriteString("\n")
	}

	return nil
}

// serializeSub 序列化 sub 元素
func (s *Serializer) serializeSub(builder *strings.Builder, sub *Sub, depth int) error {
	s.writeIndent(builder, depth)
	builder.WriteString("<sub")

	if sub.Alias != "" {
		builder.WriteString(fmt.Sprintf(` alias="%s"`, s.escapeString(sub.Alias)))
	}

	if len(sub.Content) == 0 {
		builder.WriteString("/>")
		if s.Pretty {
			builder.WriteString("\n")
		}
		return nil
	}

	builder.WriteString(">")
	if s.Pretty {
		builder.WriteString("\n")
	}

	if err := s.serializeContent(builder, sub.Content, depth+1); err != nil {
		return err
	}

	s.writeIndent(builder, depth)
	builder.WriteString("</sub>")
	if s.Pretty {
		builder.WriteString("\n")
	}

	return nil
}

// serializeVoice 序列化 voice 元素
func (s *Serializer) serializeVoice(builder *strings.Builder, voice *Voice, depth int) error {
	s.writeIndent(builder, depth)
	builder.WriteString("<voice")

	if voice.Gender != "" {
		builder.WriteString(fmt.Sprintf(` gender="%s"`, s.escapeString(voice.Gender)))
	}
	if voice.Age != "" {
		builder.WriteString(fmt.Sprintf(` age="%s"`, s.escapeString(voice.Age)))
	}
	if voice.Variant != "" {
		builder.WriteString(fmt.Sprintf(` variant="%s"`, s.escapeString(voice.Variant)))
	}
	if voice.Name != "" {
		builder.WriteString(fmt.Sprintf(` name="%s"`, s.escapeString(voice.Name)))
	}
	if voice.Languages != "" {
		builder.WriteString(fmt.Sprintf(` xml:lang="%s"`, s.escapeString(voice.Languages)))
	}

	if len(voice.Content) == 0 {
		builder.WriteString("/>")
		if s.Pretty {
			builder.WriteString("\n")
		}
		return nil
	}

	builder.WriteString(">")
	if s.Pretty {
		builder.WriteString("\n")
	}

	if err := s.serializeContent(builder, voice.Content, depth+1); err != nil {
		return err
	}

	s.writeIndent(builder, depth)
	builder.WriteString("</voice>")
	if s.Pretty {
		builder.WriteString("\n")
	}

	return nil
}

// serializeW 序列化 w 元素
func (s *Serializer) serializeW(builder *strings.Builder, w *W, depth int) error {
	s.writeIndent(builder, depth)
	builder.WriteString("<w")

	if w.Role != "" {
		builder.WriteString(fmt.Sprintf(` role="%s"`, s.escapeString(w.Role)))
	}

	if len(w.Content) == 0 {
		builder.WriteString("/>")
		if s.Pretty {
			builder.WriteString("\n")
		}
		return nil
	}

	builder.WriteString(">")
	if s.Pretty {
		builder.WriteString("\n")
	}

	if err := s.serializeContent(builder, w.Content, depth+1); err != nil {
		return err
	}

	s.writeIndent(builder, depth)
	builder.WriteString("</w>")
	if s.Pretty {
		builder.WriteString("\n")
	}

	return nil
}

// writeIndent 写入缩进
func (s *Serializer) writeIndent(builder *strings.Builder, depth int) {
	if s.Pretty {
		for i := 0; i < depth; i++ {
			builder.WriteString(s.Indent)
		}
	}
}
