package models

type Alignment string

const (
	ALIGNMENT_UNSPECIFIED = Alignment("UNSPECIFIED")
	ALIGNMENT_START       = Alignment("START")
	ALIGNMENT_END         = Alignment("END")
	ALIGNMENT_CENTER      = Alignment("CENTER")
	ALIGNMENT_JUSTIFIED   = Alignment("JUSTIFIED")
)
