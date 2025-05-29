package ssml

// Builder SSML 构建器
type Builder struct {
	speak *Speak
}

// NewBuilder 创建新的构建器
func NewBuilder() *Builder {
	return &Builder{
		speak: &Speak{
			Content: make([]interface{}, 0),
		},
	}
}

// Version 设置 SSML 版本
func (b *Builder) Version(version string) *Builder {
	b.speak.Version = version
	return b
}

// Lang 设置语言
func (b *Builder) Lang(lang string) *Builder {
	b.speak.Lang = lang
	return b
}

// Text 添加文本内容
func (b *Builder) Text(text string) *Builder {
	b.speak.Content = append(b.speak.Content, Text{Content: text})
	return b
}

// Audio 添加音频元素
func (b *Builder) Audio(src string, fallbackText string) *Builder {
	audio := &Audio{Src: src}
	if fallbackText != "" {
		audio.Content = []interface{}{Text{Content: fallbackText}}
	}
	b.speak.Content = append(b.speak.Content, audio)
	return b
}

// Break 添加停顿元素
func (b *Builder) Break(time string, strength string) *Builder {
	breakElem := &Break{
		Time:     time,
		Strength: strength,
	}
	b.speak.Content = append(b.speak.Content, breakElem)
	return b
}

// BreakTime 添加时间停顿
func (b *Builder) BreakTime(time string) *Builder {
	return b.Break(time, "")
}

// BreakStrength 添加强度停顿
func (b *Builder) BreakStrength(strength string) *Builder {
	return b.Break("", strength)
}

// Emphasis 添加强调元素
func (b *Builder) Emphasis(level string, builderFunc func(*ElementBuilder)) *Builder {
	emphasis := &Emphasis{Level: level}

	if builderFunc != nil {
		elementBuilder := &ElementBuilder{}
		builderFunc(elementBuilder)
		emphasis.Content = elementBuilder.content
	}

	b.speak.Content = append(b.speak.Content, emphasis)
	return b
}

// EmphasisText 添加强调文本
func (b *Builder) EmphasisText(level string, text string) *Builder {
	return b.Emphasis(level, func(eb *ElementBuilder) {
		eb.Text(text)
	})
}

// Paragraph 添加段落元素
func (b *Builder) Paragraph(builderFunc func(*ElementBuilder)) *Builder {
	paragraph := &Paragraph{}

	if builderFunc != nil {
		elementBuilder := &ElementBuilder{}
		builderFunc(elementBuilder)
		paragraph.Content = elementBuilder.content
	}

	b.speak.Content = append(b.speak.Content, paragraph)
	return b
}

// ParagraphText 添加段落文本
func (b *Builder) ParagraphText(text string) *Builder {
	return b.Paragraph(func(eb *ElementBuilder) {
		eb.Text(text)
	})
}

// Phoneme 添加发音元素
func (b *Builder) Phoneme(alphabet, ph string, builderFunc func(*ElementBuilder)) *Builder {
	phoneme := &Phoneme{
		Alphabet: alphabet,
		Ph:       ph,
	}

	if builderFunc != nil {
		elementBuilder := &ElementBuilder{}
		builderFunc(elementBuilder)
		phoneme.Content = elementBuilder.content
	}

	b.speak.Content = append(b.speak.Content, phoneme)
	return b
}

// PhonemeText 添加发音文本
func (b *Builder) PhonemeText(alphabet, ph, text string) *Builder {
	return b.Phoneme(alphabet, ph, func(eb *ElementBuilder) {
		eb.Text(text)
	})
}

// Prosody 添加韵律元素
func (b *Builder) Prosody(rate, pitch, range_, volume string, builderFunc func(*ElementBuilder)) *Builder {
	prosody := &Prosody{
		Rate:   rate,
		Pitch:  pitch,
		Range:  range_,
		Volume: volume,
	}

	if builderFunc != nil {
		elementBuilder := &ElementBuilder{}
		builderFunc(elementBuilder)
		prosody.Content = elementBuilder.content
	}

	b.speak.Content = append(b.speak.Content, prosody)
	return b
}

// ProsodyText 添加韵律文本
func (b *Builder) ProsodyText(rate, pitch, range_, volume, text string) *Builder {
	return b.Prosody(rate, pitch, range_, volume, func(eb *ElementBuilder) {
		eb.Text(text)
	})
}

// Rate 添加速度控制的韵律
func (b *Builder) Rate(rate string, builderFunc func(*ElementBuilder)) *Builder {
	return b.Prosody(rate, "", "", "", builderFunc)
}

// RateText 添加速度控制的韵律文本
func (b *Builder) RateText(rate, text string) *Builder {
	return b.Rate(rate, func(eb *ElementBuilder) {
		eb.Text(text)
	})
}

// Pitch 添加音调控制的韵律
func (b *Builder) Pitch(pitch string, builderFunc func(*ElementBuilder)) *Builder {
	return b.Prosody("", pitch, "", "", builderFunc)
}

// PitchText 添加音调控制的韵律文本
func (b *Builder) PitchText(pitch, text string) *Builder {
	return b.Pitch(pitch, func(eb *ElementBuilder) {
		eb.Text(text)
	})
}

// Volume 添加音量控制的韵律
func (b *Builder) Volume(volume string, builderFunc func(*ElementBuilder)) *Builder {
	return b.Prosody("", "", "", volume, builderFunc)
}

// VolumeText 添加音量控制的韵律文本
func (b *Builder) VolumeText(volume, text string) *Builder {
	return b.Volume(volume, func(eb *ElementBuilder) {
		eb.Text(text)
	})
}

// Sentence 添加句子元素
func (b *Builder) Sentence(builderFunc func(*ElementBuilder)) *Builder {
	sentence := &Sentence{}

	if builderFunc != nil {
		elementBuilder := &ElementBuilder{}
		builderFunc(elementBuilder)
		sentence.Content = elementBuilder.content
	}

	b.speak.Content = append(b.speak.Content, sentence)
	return b
}

// SentenceText 添加句子文本
func (b *Builder) SentenceText(text string) *Builder {
	return b.Sentence(func(eb *ElementBuilder) {
		eb.Text(text)
	})
}

// Sub 添加替换元素
func (b *Builder) Sub(alias string, builderFunc func(*ElementBuilder)) *Builder {
	sub := &Sub{Alias: alias}

	if builderFunc != nil {
		elementBuilder := &ElementBuilder{}
		builderFunc(elementBuilder)
		sub.Content = elementBuilder.content
	}

	b.speak.Content = append(b.speak.Content, sub)
	return b
}

// SubText 添加替换文本
func (b *Builder) SubText(alias, text string) *Builder {
	return b.Sub(alias, func(eb *ElementBuilder) {
		eb.Text(text)
	})
}

// Voice 添加声音元素
func (b *Builder) Voice(gender, age, variant, name, lang string, builderFunc func(*ElementBuilder)) *Builder {
	voice := &Voice{
		Gender:    gender,
		Age:       age,
		Variant:   variant,
		Name:      name,
		Languages: lang,
	}

	if builderFunc != nil {
		elementBuilder := &ElementBuilder{}
		builderFunc(elementBuilder)
		voice.Content = elementBuilder.content
	}

	b.speak.Content = append(b.speak.Content, voice)
	return b
}

// VoiceText 添加声音文本
func (b *Builder) VoiceText(gender, age, variant, name, lang, text string) *Builder {
	return b.Voice(gender, age, variant, name, lang, func(eb *ElementBuilder) {
		eb.Text(text)
	})
}

// VoiceName 根据名称添加声音
func (b *Builder) VoiceName(name string, builderFunc func(*ElementBuilder)) *Builder {
	return b.Voice("", "", "", name, "", builderFunc)
}

// VoiceNameText 根据名称添加声音文本
func (b *Builder) VoiceNameText(name, text string) *Builder {
	return b.VoiceName(name, func(eb *ElementBuilder) {
		eb.Text(text)
	})
}

// W 添加单词元素
func (b *Builder) W(role string, builderFunc func(*ElementBuilder)) *Builder {
	w := &W{Role: role}

	if builderFunc != nil {
		elementBuilder := &ElementBuilder{}
		builderFunc(elementBuilder)
		w.Content = elementBuilder.content
	}

	b.speak.Content = append(b.speak.Content, w)
	return b
}

// WText 添加单词文本
func (b *Builder) WText(role, text string) *Builder {
	return b.W(role, func(eb *ElementBuilder) {
		eb.Text(text)
	})
}

// Build 构建 SSML
func (b *Builder) Build() *Speak {
	return b.speak
}

// BuildString 构建 SSML 字符串
func (b *Builder) BuildString(pretty bool) (string, error) {
	serializer := NewSerializer(pretty)
	return serializer.Serialize(b.speak)
}

// ElementBuilder 元素构建器，用于构建嵌套内容
type ElementBuilder struct {
	content []interface{}
}

// Text 添加文本
func (eb *ElementBuilder) Text(text string) *ElementBuilder {
	eb.content = append(eb.content, Text{Content: text})
	return eb
}

// Audio 添加音频
func (eb *ElementBuilder) Audio(src string, fallbackText string) *ElementBuilder {
	audio := &Audio{Src: src}
	if fallbackText != "" {
		audio.Content = []interface{}{Text{Content: fallbackText}}
	}
	eb.content = append(eb.content, audio)
	return eb
}

// Break 添加停顿
func (eb *ElementBuilder) Break(time string, strength string) *ElementBuilder {
	breakElem := &Break{
		Time:     time,
		Strength: strength,
	}
	eb.content = append(eb.content, breakElem)
	return eb
}

// BreakTime 添加时间停顿
func (eb *ElementBuilder) BreakTime(time string) *ElementBuilder {
	return eb.Break(time, "")
}

// BreakStrength 添加强度停顿
func (eb *ElementBuilder) BreakStrength(strength string) *ElementBuilder {
	return eb.Break("", strength)
}

// Emphasis 添加强调
func (eb *ElementBuilder) Emphasis(level string, builderFunc func(*ElementBuilder)) *ElementBuilder {
	emphasis := &Emphasis{Level: level}

	if builderFunc != nil {
		nestedBuilder := &ElementBuilder{}
		builderFunc(nestedBuilder)
		emphasis.Content = nestedBuilder.content
	}

	eb.content = append(eb.content, emphasis)
	return eb
}

// EmphasisText 添加强调文本
func (eb *ElementBuilder) EmphasisText(level string, text string) *ElementBuilder {
	return eb.Emphasis(level, func(nested *ElementBuilder) {
		nested.Text(text)
	})
}

// Phoneme 添加发音
func (eb *ElementBuilder) Phoneme(alphabet, ph string, builderFunc func(*ElementBuilder)) *ElementBuilder {
	phoneme := &Phoneme{
		Alphabet: alphabet,
		Ph:       ph,
	}

	if builderFunc != nil {
		nestedBuilder := &ElementBuilder{}
		builderFunc(nestedBuilder)
		phoneme.Content = nestedBuilder.content
	}

	eb.content = append(eb.content, phoneme)
	return eb
}

// PhonemeText 添加发音文本
func (eb *ElementBuilder) PhonemeText(alphabet, ph, text string) *ElementBuilder {
	return eb.Phoneme(alphabet, ph, func(nested *ElementBuilder) {
		nested.Text(text)
	})
}

// Prosody 添加韵律
func (eb *ElementBuilder) Prosody(rate, pitch, range_, volume string, builderFunc func(*ElementBuilder)) *ElementBuilder {
	prosody := &Prosody{
		Rate:   rate,
		Pitch:  pitch,
		Range:  range_,
		Volume: volume,
	}

	if builderFunc != nil {
		nestedBuilder := &ElementBuilder{}
		builderFunc(nestedBuilder)
		prosody.Content = nestedBuilder.content
	}

	eb.content = append(eb.content, prosody)
	return eb
}

// Sub 添加替换
func (eb *ElementBuilder) Sub(alias string, builderFunc func(*ElementBuilder)) *ElementBuilder {
	sub := &Sub{Alias: alias}

	if builderFunc != nil {
		nestedBuilder := &ElementBuilder{}
		builderFunc(nestedBuilder)
		sub.Content = nestedBuilder.content
	}

	eb.content = append(eb.content, sub)
	return eb
}

// SubText 添加替换文本
func (eb *ElementBuilder) SubText(alias, text string) *ElementBuilder {
	return eb.Sub(alias, func(nested *ElementBuilder) {
		nested.Text(text)
	})
}

// Voice 添加声音
func (eb *ElementBuilder) Voice(gender, age, variant, name, lang string, builderFunc func(*ElementBuilder)) *ElementBuilder {
	voice := &Voice{
		Gender:    gender,
		Age:       age,
		Variant:   variant,
		Name:      name,
		Languages: lang,
	}

	if builderFunc != nil {
		nestedBuilder := &ElementBuilder{}
		builderFunc(nestedBuilder)
		voice.Content = nestedBuilder.content
	}

	eb.content = append(eb.content, voice)
	return eb
}

// W 添加单词
func (eb *ElementBuilder) W(role string, builderFunc func(*ElementBuilder)) *ElementBuilder {
	w := &W{Role: role}

	if builderFunc != nil {
		nestedBuilder := &ElementBuilder{}
		builderFunc(nestedBuilder)
		w.Content = nestedBuilder.content
	}

	eb.content = append(eb.content, w)
	return eb
}
