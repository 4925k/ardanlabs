package nlp_test

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/4925k/practical_go_foundations/nlp"
)

func ExampleTokenize() {
	text := "hello world"
	tokens := nlp.Tokenize(text)
	fmt.Println(tokens)

	// Output: [hello world]
}

func TestTokenize(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "hello world ",
			args: args{
				text: "hello world",
			},
			want: []string{"hello", "world"},
		},
		{
			name: "empty",
			args: args{
				text: "",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// if tt.args.want != tt.want {} // cant compare slices with ==
			if got := nlp.Tokenize(tt.args.text); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Tokenize() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func FuzzToken(f *testing.F) {
	f.Fuzz(func(t *testing.T, text string) {
		tokens := nlp.Tokenize(text)
		lText := strings.ToLower(text)
		for _, tok := range tokens {
			if !strings.Contains(lText, tok) {
				t.Fatal(tok)
			}
		}
	})

}
