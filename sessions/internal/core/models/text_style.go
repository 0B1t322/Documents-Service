package models

type TextStyle struct {
	ID              int
	FontFamily      string
	FontWeight      int
	FontSize        Dimension
	Bold            bool
	Underline       bool
	Italic          bool
	BackgroundColor Color
	ForegroundColor Color
}
