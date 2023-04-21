package models

type TextStyle struct {
	ID              int
	FontFamily      string
	FontWeight      int
	Bold            bool
	Underline       bool
	Italic          bool
	BackgroundColor Color
	ForegroundColor Color
}
