package main

import (
	"reflect"
	"testing"
)

type TreeTest struct {
	Text       string
	SeedValues []Node
	Root       *Node
	Encoded    []byte
}

var tests = []TreeTest{
	TreeTest{
		Text: "cbccaccaccbcacacf",
		SeedValues: ([]Node{
			Node{Char: rune('c'), Count: 10},
			Node{Char: rune('a'), Count: 4},
			Node{Char: rune('b'), Count: 2},
			Node{Char: rune('f'), Count: 1},
		}),
		Root: &Node{Char: -1, Count: 17,
			Left: &Node{Char: rune('c'), Code: []byte{0}, Count: 10},
			Right: &Node{Char: -1, Count: 7,
				Left: &Node{Char: rune('a'), Code: []byte{1, 0}, Count: 4},
				Right: &Node{Char: -1, Count: 3,
					Left:  &Node{Char: rune('b'), Code: []byte{1, 1, 0}, Count: 2},
					Right: &Node{Char: rune('f'), Code: []byte{1, 1, 1}, Count: 1},
				},
			},
		},
		Encoded: []byte{0, 1, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 1, 1},
	},
}

func TestCountRunes(t *testing.T) {
	for _, v := range tests {
		n := countRunes(v.Text)
		if !reflect.DeepEqual(v.SeedValues, *n) {
			t.Error("Did not equal values for string", v.SeedValues, *n)
		}
	}
}

var falseTest = map[string]*[]Node{
	"aaabbw": &([]Node{
		Node{Char: rune('b'), Count: 4},
		Node{Char: rune('a'), Count: 3},
	}),
	"qww": &([]Node{
		Node{Char: rune('w'), Count: 1},
		Node{Char: rune('q'), Count: 1},
	}),
}

func TestCountRunesFalse(t *testing.T) {
	t.Parallel()
	for k, v := range falseTest {
		n := countRunes(k)
		if reflect.DeepEqual(v, n) {
			t.Error("Did equal values for string", k, *v, *n)
		}
	}
}

func TestTreeCreation(t *testing.T) {
	for _, v := range tests {
		r := constructTree(v.Text)
		if !reflect.DeepEqual(v.Root, r) {
			t.Error("Not equal roots", *r, *v.Root)
		}
	}
}

func TestEncoding(t *testing.T) {
	for _, v := range tests {
		e := encodeString(v.Text, v.Root)
		if !reflect.DeepEqual(v.Encoded, e) {
			t.Error("Error encoding", v.Encoded, e)
		}
	}
}

func TestDecoding(t *testing.T) {
	for _, v := range tests {
		text := decode(v.Encoded, v.Root)
		if text != v.Text {
			t.Error("Text did not equal", text, v.Text)
		}
	}
}
