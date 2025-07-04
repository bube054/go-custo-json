package jsonvx

type Json struct {
}

func (j *Json) IsRawJson() bool {
	return true
}

func (j *Json) Parse() any {
	return true
}

func (j *Json) RawJson() any {
	return true
}

func (j *Json) Stringify() any {
	return 1
}

func NewJson() any {
	return 0
}
