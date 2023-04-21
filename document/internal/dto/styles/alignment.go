package dto

type Alignment int

const (
	AlignmentUnknown Alignment = iota
	AlignmentStart
	AlignmentEnd
	AlignmentCenter
	AlignmentJustified
)
