package models

type ParagraphElement struct {
	ID        int
	Index     int
	TextRun   *TextRun
	PageBreak *PageBreak
}

func (p ParagraphElement) GetContent() string {
	//TODO handle element type
	return p.TextRun.Content
}

func (p *ParagraphElement) SetContent(s string) {
	//TODO handle element type
	p.TextRun.Content = s
}
