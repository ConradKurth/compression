package main

import (
	"bytes"
	"encoding/gob"
	"errors"
	"os"
	"sort"
)

var NonInit = errors.New("NOT_INIT")

type Node struct {
	Char  rune
	Count int
	Code  []byte
	Left  *Node
	Right *Node
}

type Sorted []Node

func (s Sorted) Len() int {
	return len(s)
}

func (s Sorted) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s Sorted) Less(i, j int) bool {
	return s[i].Count > s[j].Count
}

type saveObject struct {
	Root     *Node
	Encoding []byte
}

type C struct {
	Root        *Node
	Encoding    []byte
	Text        string
	initialized bool
}

// create an empty compression object
func New() *C {
	return &C{}
}

// getNodeByRune return the node specified by the rune
func getNodeByRune(root *Node, r rune) *Node {
	if root == nil {
		return nil
	}
	if root.Char == r {
		return root
	}
	n := getNodeByRune(root.Left, r)
	if n != nil {
		return n
	}
	n = getNodeByRune(root.Right, r)
	return n
}

// getNodeByByte will return the node by checking each byte in the
// passed in byte array
func getNodeByByte(root *Node, encoding *[]byte) *Node {
	if root == nil {
		return nil
	}

	if len(*encoding) == 0 {
		return root
	}

	n := root
	b := (*encoding)[0]
	*encoding = append((*encoding)[:0], (*encoding)[1:]...)
	nextNode := root.Right
	if b == 0 {
		nextNode = root.Left
	}

	t := getNodeByByte(nextNode, encoding)
	if t != nil {
		n = t
	} else {
		*encoding = append([]byte{b}, *encoding...)
	}
	return n
}

// decode is a function that will decode the compressed file
func decode(encoding []byte, root *Node) string {

	temp := make([]byte, len(encoding))
	copy(temp, encoding)

	var decode bytes.Buffer
	for len(temp) > 1 {
		n := getNodeByByte(root, &temp)
		decode.WriteRune(n.Char)
	}
	return decode.String()
}

// SaveEncoding will save the struct to the specificed file
func (c *C) SaveEncoding(file string) error {
	if !c.initialized {
		return NonInit
	}
	f, err := os.Create(file)
	if err != nil {
		return err
	}

	defer f.Close()

	enc := gob.NewEncoder(f)
	s := saveObject{Root: c.Root, Encoding: c.Encoding}
	if err := enc.Encode(s); err != nil {
		return err
	}
	return nil
}

// GetEncoding will return the byte array for the encoded text
// this struct needs to be initialized first
func (c *C) GetEncoding() ([]byte, error) {
	if !c.initialized {
		return nil, NonInit
	}

	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	s := saveObject{Root: c.Root, Encoding: c.Encoding}
	if err := enc.Encode(s); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

// LoadEncoding will load an encoding saved by a file and then
// decode the string in the function
func (c *C) LoadEncoding(file string) (string, error) {
	f, err := os.Open(file)
	if err != nil {
		return "", err
	}

	dec := gob.NewDecoder(f)
	s := saveObject{}

	if err := dec.Decode(s); err != nil {
		return "", err
	}
	c.initialized = true
	c.Root = s.Root
	c.Encoding = s.Encoding
	c.Text = decode(s.Encoding, s.Root)
	return c.Text, nil
}

// countRunes is a helper function that will count each letter in a string
// then return an ordered list of nodes that contain the count and rune
func countRunes(text string) *[]Node {
	m := make(map[rune]Node)
	for _, l := range text {
		n, ok := m[l]
		if !ok {
			n = Node{Char: l}
		}
		n.Count++
		m[l] = n
	}
	list := make([]Node, 0, len(m))
	for _, v := range m {
		list = append(list, v)
	}

	sort.Sort(Sorted(list))
	return &list
}

// Compress is a function that will compress the text given.
// this is a main point of entry
func (c *C) Compress(text string) {
	root := constructTree(text)

	c.Encoding = encodeString(text, root)
	c.initialized = true
	c.Root = root
}

// encodeString will take a text string and root mapping that will be used
// to encode the string with the proper encoding
func encodeString(text string, root *Node) []byte {
	resp := make([]byte, 0, len(text))
	for _, l := range text {
		e := getNodeByRune(root, l)
		resp = append(resp, e.Code...)
	}
	return resp
}

// popNode is a helper function that will pop off a node and
// return that popped off node
func popNode(nodes *[]Node) *Node {
	if len(*nodes) == 0 {
		return nil
	}

	n := (*nodes)[len(*nodes)-1]
	*nodes = (*nodes)[:len(*nodes)-1]
	return &n
}

// assignCodes is a recursive function that take the constructed
// tree in construct tree and gives each node that is not a grouping node
// and assigned them a code of 0 or 1 depending the route taken
func assignCodes(node *Node, code []byte) {
	if node == nil {
		return
	}
	if node.Char != -1 {
		tmp := make([]byte, len(code))
		copy(tmp, code)
		(*node).Code = tmp
	}
	assignCodes(node.Left, append(code, 0))
	assignCodes(node.Right, append(code, 1))
}

// constructTree is a function that will take text and contruct
// and binary tree that will be the mapping for our encoding
func constructTree(text string) *Node {

	nodes := countRunes(text)
	for len(*nodes) > 1 {
		n1 := popNode(nodes)
		n2 := popNode(nodes)
		p := Node{Count: n2.Count + n1.Count,
			Char:  -1,
			Right: n1,
			Left:  n2}
		*nodes = append(*nodes, p)
		sort.Sort(Sorted(*nodes))
	}
	top := &(*nodes)[0]
	assignCodes(top, nil)
	return top
}
