package jsonvx

import (
	"errors"
	"reflect"
	"testing"
)

type SplitTokens struct {
	msg            string
	input          []byte
	cfg            *Config
	expectedChunks [][2]int
	expectedErr    error
}

func RunSpliceTokens(t *testing.T, tests []SplitTokens) {
	for _, test := range tests {
		t.Run(test.msg, func(t *testing.T) {
			l := NewLexer(test.input, test.cfg)
			tokens := l.Tokens()
			chunks, err := tokens.Split()

			if !reflect.DeepEqual(chunks, test.expectedChunks) || !errors.Is(err, test.expectedErr) {
				t.Errorf("got (%v, %v), expected (%v, %v)", chunks, err, test.expectedChunks, test.expectedErr)
			}
		})
	}
}

func TestJSONSplit(t *testing.T) {
	var tests = []SplitTokens{
		// {msg: "Splice nothing", input: []byte(``), expectedChunks: [][2]int{{0, 0}}, expectedErr: nil, cfg: nil},
		{msg: "Splice single null", input: []byte(`null`), expectedChunks: [][2]int{{0, 0}}, expectedErr: nil, cfg: nil},
	}

	RunSpliceTokens(t, tests)
}
