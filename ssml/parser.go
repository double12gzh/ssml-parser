package ssml

import (
	"encoding/xml"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Parser SSML 解析器
type Parser struct {
	config *ValidationConfig
}

// NewParser 创建新的解析器
func NewParser(config *ValidationConfig) *Parser {
	if config == nil {
		config = DefaultValidationConfig()
	}
	return &Parser{config: config}
}

// Parse 解析 SSML 字符串
func (p *Parser) Parse(ssmlContent string) (*ParseResult, error) {
	return p.ParseReader(strings.NewReader(ssmlContent))
}

// ParseReader 从 Reader 解析 SSML
func (p *Parser) ParseReader(reader io.Reader) (*ParseResult, error) {
	result := &ParseResult{
		Warnings: []string{},
		Errors:   []string{},
	}

	decoder := xml.NewDecoder(reader)
	var root *Speak

	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("XML parsing error: %v", err))
			return result, err
		}

		switch se := token.(type) {
		case xml.StartElement:
			if se.Name.Local == "speak" {
				speak := &Speak{}
				if err := p.parseSpeak(decoder, se, speak); err != nil {
					result.Errors = append(result.Errors, fmt.Sprintf("Error parsing speak element: %v", err))
					return result, err
				}
				root = speak
			} else {
				result.Warnings = append(result.Warnings, fmt.Sprintf("Root element should be 'speak', found '%s'", se.Name.Local))
			}
		}
	}

	if root == nil {
		result.Errors = append(result.Errors, "No valid SSML root element found")
		return result, fmt.Errorf("no valid SSML root element found")
	}

	result.Root = root

	// 验证解析结果
	if err := p.validate(root, result); err != nil {
		return result, err
	}

	return result, nil
}

// parseSpeak 解析 speak 元素
func (p *Parser) parseSpeak(decoder *xml.Decoder, start xml.StartElement, speak *Speak) error {
	speak.XMLName = start.Name

	// 解析属性
	for _, attr := range start.Attr {
		switch attr.Name.Local {
		case "version":
			speak.Version = attr.Value
		case "lang":
			speak.Lang = attr.Value
		}
	}

	// 解析内容
	content, err := p.parseContent(decoder, "speak")
	if err != nil {
		return err
	}
	speak.Content = content

	return nil
}

// parseContent 解析元素内容
func (p *Parser) parseContent(decoder *xml.Decoder, parentTag string) ([]interface{}, error) {
	var content []interface{}

	for {
		token, err := decoder.Token()
		if err != nil {
			return nil, err
		}

		switch se := token.(type) {
		case xml.StartElement:
			element, err := p.parseElement(decoder, se)
			if err != nil {
				return nil, err
			}
			if element != nil {
				content = append(content, element)
			}

		case xml.CharData:
			text := strings.TrimSpace(string(se))
			if text != "" {
				content = append(content, Text{Content: text})
			}

		case xml.EndElement:
			if se.Name.Local == parentTag {
				return content, nil
			}
		}
	}
}

// parseElement 解析单个元素
func (p *Parser) parseElement(decoder *xml.Decoder, start xml.StartElement) (interface{}, error) {
	switch start.Name.Local {
	case "audio":
		return p.parseAudio(decoder, start)
	case "break":
		return p.parseBreak(decoder, start)
	case "emphasis":
		return p.parseEmphasis(decoder, start)
	case "p":
		return p.parseParagraph(decoder, start)
	case "phoneme":
		return p.parsePhoneme(decoder, start)
	case "prosody":
		return p.parseProsody(decoder, start)
	case "s":
		return p.parseSentence(decoder, start)
	case "sub":
		return p.parseSub(decoder, start)
	case "voice":
		return p.parseVoice(decoder, start)
	case "w":
		return p.parseW(decoder, start)
	default:
		if !p.config.AllowUnknownElements {
			return nil, fmt.Errorf("unknown element: %s", start.Name.Local)
		}
		// 跳过未知元素
		return p.skipElement(decoder, start.Name.Local)
	}
}

// parseAudio 解析 audio 元素
func (p *Parser) parseAudio(decoder *xml.Decoder, start xml.StartElement) (*Audio, error) {
	audio := &Audio{XMLName: start.Name}

	for _, attr := range start.Attr {
		if attr.Name.Local == "src" {
			audio.Src = attr.Value
		}
	}

	content, err := p.parseContent(decoder, "audio")
	if err != nil {
		return nil, err
	}
	audio.Content = content

	return audio, nil
}

// parseBreak 解析 break 元素
func (p *Parser) parseBreak(decoder *xml.Decoder, start xml.StartElement) (*Break, error) {
	breakElem := &Break{XMLName: start.Name}

	for _, attr := range start.Attr {
		switch attr.Name.Local {
		case "time":
			breakElem.Time = attr.Value
		case "strength":
			breakElem.Strength = attr.Value
		}
	}

	// break 是自闭合元素，跳过到结束标签
	for {
		token, err := decoder.Token()
		if err != nil {
			return nil, err
		}
		if endElement, ok := token.(xml.EndElement); ok && endElement.Name.Local == "break" {
			break
		}
	}

	return breakElem, nil
}

// parseEmphasis 解析 emphasis 元素
func (p *Parser) parseEmphasis(decoder *xml.Decoder, start xml.StartElement) (*Emphasis, error) {
	emphasis := &Emphasis{XMLName: start.Name}

	for _, attr := range start.Attr {
		if attr.Name.Local == "level" {
			emphasis.Level = attr.Value
		}
	}

	content, err := p.parseContent(decoder, "emphasis")
	if err != nil {
		return nil, err
	}
	emphasis.Content = content

	return emphasis, nil
}

// parseParagraph 解析 p 元素
func (p *Parser) parseParagraph(decoder *xml.Decoder, start xml.StartElement) (*Paragraph, error) {
	paragraph := &Paragraph{XMLName: start.Name}

	content, err := p.parseContent(decoder, "p")
	if err != nil {
		return nil, err
	}
	paragraph.Content = content

	return paragraph, nil
}

// parsePhoneme 解析 phoneme 元素
func (p *Parser) parsePhoneme(decoder *xml.Decoder, start xml.StartElement) (*Phoneme, error) {
	phoneme := &Phoneme{XMLName: start.Name}

	for _, attr := range start.Attr {
		switch attr.Name.Local {
		case "alphabet":
			phoneme.Alphabet = attr.Value
		case "ph":
			phoneme.Ph = attr.Value
		}
	}

	content, err := p.parseContent(decoder, "phoneme")
	if err != nil {
		return nil, err
	}
	phoneme.Content = content

	return phoneme, nil
}

// parseProsody 解析 prosody 元素
func (p *Parser) parseProsody(decoder *xml.Decoder, start xml.StartElement) (*Prosody, error) {
	prosody := &Prosody{XMLName: start.Name}

	for _, attr := range start.Attr {
		switch attr.Name.Local {
		case "rate":
			prosody.Rate = attr.Value
		case "pitch":
			prosody.Pitch = attr.Value
		case "range":
			prosody.Range = attr.Value
		case "volume":
			prosody.Volume = attr.Value
		}
	}

	content, err := p.parseContent(decoder, "prosody")
	if err != nil {
		return nil, err
	}
	prosody.Content = content

	return prosody, nil
}

// parseSentence 解析 s 元素
func (p *Parser) parseSentence(decoder *xml.Decoder, start xml.StartElement) (*Sentence, error) {
	sentence := &Sentence{XMLName: start.Name}

	content, err := p.parseContent(decoder, "s")
	if err != nil {
		return nil, err
	}
	sentence.Content = content

	return sentence, nil
}

// parseSub 解析 sub 元素
func (p *Parser) parseSub(decoder *xml.Decoder, start xml.StartElement) (*Sub, error) {
	sub := &Sub{XMLName: start.Name}

	for _, attr := range start.Attr {
		if attr.Name.Local == "alias" {
			sub.Alias = attr.Value
		}
	}

	content, err := p.parseContent(decoder, "sub")
	if err != nil {
		return nil, err
	}
	sub.Content = content

	return sub, nil
}

// parseVoice 解析 voice 元素
func (p *Parser) parseVoice(decoder *xml.Decoder, start xml.StartElement) (*Voice, error) {
	voice := &Voice{XMLName: start.Name}

	for _, attr := range start.Attr {
		switch attr.Name.Local {
		case "gender":
			voice.Gender = attr.Value
		case "age":
			voice.Age = attr.Value
		case "variant":
			voice.Variant = attr.Value
		case "name":
			voice.Name = attr.Value
		case "lang":
			voice.Languages = attr.Value
		}
	}

	content, err := p.parseContent(decoder, "voice")
	if err != nil {
		return nil, err
	}
	voice.Content = content

	return voice, nil
}

// parseW 解析 w 元素
func (p *Parser) parseW(decoder *xml.Decoder, start xml.StartElement) (*W, error) {
	w := &W{XMLName: start.Name}

	for _, attr := range start.Attr {
		if attr.Name.Local == "role" {
			w.Role = attr.Value
		}
	}

	content, err := p.parseContent(decoder, "w")
	if err != nil {
		return nil, err
	}
	w.Content = content

	return w, nil
}

// skipElement 跳过未知元素
func (p *Parser) skipElement(decoder *xml.Decoder, tagName string) (interface{}, error) {
	for {
		token, err := decoder.Token()
		if err != nil {
			return nil, err
		}
		if endElement, ok := token.(xml.EndElement); ok && endElement.Name.Local == tagName {
			break
		}
	}
	return nil, nil
}

// validate 验证解析结果
func (p *Parser) validate(speak *Speak, result *ParseResult) error {
	if p.config.StrictMode {
		return p.strictValidation(speak, result)
	}
	return p.basicValidation(speak, result)
}

// basicValidation 基本验证
func (p *Parser) basicValidation(speak *Speak, result *ParseResult) error {
	if speak.Version == "" {
		result.Warnings = append(result.Warnings, "SSML version not specified")
	}

	if speak.Lang == "" {
		result.Warnings = append(result.Warnings, "Language not specified")
	}

	return p.validateNestingDepth(speak.Content, 0, result)
}

// strictValidation 严格验证
func (p *Parser) strictValidation(speak *Speak, result *ParseResult) error {
	if speak.Version == "" {
		result.Errors = append(result.Errors, "SSML version is required in strict mode")
		return fmt.Errorf("SSML version is required")
	}

	if speak.Lang == "" {
		result.Errors = append(result.Errors, "Language is required in strict mode")
		return fmt.Errorf("language is required")
	}

	return p.validateNestingDepth(speak.Content, 0, result)
}

// validateNestingDepth 验证嵌套深度
func (p *Parser) validateNestingDepth(content []interface{}, depth int, result *ParseResult) error {
	if depth > p.config.MaxNestingDepth {
		err := fmt.Errorf("maximum nesting depth exceeded: %d", p.config.MaxNestingDepth)
		result.Errors = append(result.Errors, err.Error())
		return err
	}

	for _, item := range content {
		if element, ok := item.(SSMLElement); ok {
			if err := p.validateNestingDepth(element.GetContent(), depth+1, result); err != nil {
				return err
			}
		}
	}

	return nil
}

// validateTimeAttribute 验证时间属性
func (p *Parser) validateTimeAttribute(timeStr string) error {
	if timeStr == "" {
		return nil
	}

	// 支持的时间格式：1s, 500ms, 1.5s 等
	if strings.HasSuffix(timeStr, "ms") {
		_, err := strconv.ParseFloat(strings.TrimSuffix(timeStr, "ms"), 64)
		return err
	}

	if strings.HasSuffix(timeStr, "s") {
		_, err := strconv.ParseFloat(strings.TrimSuffix(timeStr, "s"), 64)
		return err
	}

	return fmt.Errorf("invalid time format: %s", timeStr)
}
