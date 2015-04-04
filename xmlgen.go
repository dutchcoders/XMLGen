// xmlGen - A tool for generating native Golang types from xml objects.
// Copyright (C) 2015 Remco Verhoef
// Based on https://github.com/bemasher/JSONGen
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"encoding/xml"
	"errors"
	"flag"
	"fmt"
	"go/format"
	"io"

	"code.google.com/p/go-charset/charset"
	_ "code.google.com/p/go-charset/data"

	"log"
	"os"
	"strings"
	"unicode"
)

var config Config

type Config struct {
	dumpFilename string

	dumpFile  *os.File
	inputFile *os.File

	titleCase bool
	normalize bool
}

func (c *Config) Parse() (err error) {
	flag.StringVar(&config.dumpFilename, "dump", os.DevNull, "Dump tree structure to file.")
	flag.BoolVar(&config.normalize, "normalize", true, "Squash arrays of struct and determine primitive array type.")
	flag.BoolVar(&config.titleCase, "title", true, "Convert identifiers to title case, treating '_' and '-' as word boundaries.")

	flag.Parse()

	if flag.NArg() == 0 {
		config.inputFile = os.Stdin
	} else {
		config.inputFile, err = os.Open(flag.Arg(0))
		if err != nil {
			return
		}
	}

	c.dumpFile, err = os.Create(c.dumpFilename)
	if err != nil {
		return
	}

	return
}

func (c Config) Close() {
	c.dumpFile.Close()
	c.inputFile.Close()
}

// Field name sanitizer.
type Ident string

// Golang identifiers must begin with a letter and may contain letters, digits
// and _'s. If config.titleCase is true, -, _ and spaces are treated as word
// boundaries, otherwise only spaces are treated as word boundaries.
func (id Ident) String() (s string) {
	// Trim non-letter characters from the left of the identifier.
	s = strings.TrimLeftFunc(string(id), func(r rune) bool {
		return !unicode.IsLetter(r)
	})

	// Remove any invalid characters in the identifier.
	s = strings.Map(func(r rune) rune {
		if r == ' ' {
			return ' '
		}

		// Convert -'s to _'s or spaces depending on configuration.
		if r == '-' || r == '_' {
			if config.titleCase {
				return ' '
			}
			return '_'
		}

		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			return r
		}

		return -1
	}, s)

	// Perform title casing.
	s = strings.Title(s)
	// Remove spaces from the identifier.
	s = strings.Map(func(r rune) rune {
		if r == ' ' {
			return -1
		}
		return r
	}, s)

	// If the identifier is empty, output an _.
	if len(s) == 0 {
		s = "_"
	}

	return
}

// Returns a field tag for the original field name.
func (tree Tree) Tag() string {
	if tree.Attr {
		return fmt.Sprintf("`"+`xml:"%s,attr"`+"`", string(tree.Name))
	} else {
		return fmt.Sprintf("`"+`xml:"%s"`+"`", string(tree.Name))
	}
}

// xml values are translated to go types as follows:
// null   -> interface{}
// bool   -> bool
// int    -> int64
// float  -> float64
// string -> string
// object -> struct
type Type int

const (
	Interface Type = iota + 1
	Bool
	Int
	Float
	String
	Struct
)

func (t Type) String() string {
	switch t {
	case Interface:
		return "interface{}"
	case Bool:
		return "bool"
	case Int:
		return "int64"
	case Float:
		return "float64"
	case String:
		return "string"
	case Struct:
		return "struct"
	}
	return "unset"
}

// Necessary for dumping the tree for debugging.
func (t Type) MarshalText() (text []byte, err error) {
	return []byte(t.String()), nil
}

// A type tree describes parsed xml input. Elements have a name, type and
// children, list specifies if the type is a list.
type Tree struct {
	parent   *Tree
	Name     Ident `xml:",omitempty"`
	List     bool  `xml:",omitempty"`
	Attr     bool
	Type     Type
	Children []*Tree `xml:",omitempty"`
}

// A tree implements the sort interface on it's children's sanitized names.
func (t Tree) Len() int {
	return len(t.Children)
}

func (t Tree) Less(i, j int) bool {
	return t.Children[i].Name.String() < t.Children[j].Name.String()
}

func (t Tree) Swap(i, j int) {
	t.Children[i], t.Children[j] = t.Children[j], t.Children[i]
}

// Returns canonical golang of the type structure.
func (t *Tree) Format() (formatted []byte, err error) {
	// Store the raw source for debugging.
	unformatted := []byte("type " + t.formatHelper(0))

	// Attempt to format the source.
	formatted, err = format.Source(unformatted)

	// If formatting failed, return the unformatted source and the error.
	if err != nil {
		formatted = unformatted
	}
	return
}

func (t *Tree) formatHelper(depth int) (r string) {
	indent := strings.Repeat("\t", depth)

	// Print the name of the current element.
	r += indent + t.Name.String() + " "

	// On return append a tag if the field name differs from the parsed name.
	defer func() {
		if depth != 0 {
			r += " " + t.Tag()
		}
		r += "\n"
	}()

	// Prefix the type with [] if list is true.
	if t.List {
		r += "[]"
	}

	// Print type
	r += t.Type.String()

	// If the type is a struct, print struct and enclosing curly braces.
	// fmt.Println("formatHelper", t.Name, t.Type)
	if t.Type == Struct {
		r += " {\n"
		defer func() {
			r += indent + "}"
		}()

		// Recurse for each child of the struct.
		for _, child := range t.Children {
			r += child.formatHelper(depth + 1)
		}
	} else {
	}

	return
}

// Given a value which xml has been parsed into, populates the tree.
func (tree *Tree) Populate(decoder *xml.Decoder) {
	current := tree

	for {
		// Read tokens from the XML document in a stream.
		t, _ := decoder.Token()
		if t == nil {
			break
		}

		// Inspect the type of the token just read.
		switch se := t.(type) {
		case xml.StartElement:
			child := &Tree{parent: current, Attr: false, List: false}
			child.Name = Ident(se.Name.Local)
			current.Children = append(current.Children, child)
			current = child

			// att attributes
			for _, attr := range se.Attr {
				if attr.Name.Local == "xmlns" {
					continue
				}

				child := &Tree{parent: current, Name: Ident(attr.Name.Local), Attr: true, List: false, Type: String}
				current.Children = append(current.Children, child)
			}

		case xml.EndElement:
			if len(current.Children) > 1 {
				current.Type = Struct
			} else {
				current.Type = String
			}

			current = current.parent

		}
	}

}

var ErrNotSupported = errors.New("Not supported")

// merge will merge the Children of src into dst
func merge(dst, src *Tree) error {
	if dst.Type != Struct {
		return ErrNotSupported
	}

	if src.Type != Struct {
		return ErrNotSupported
	}

	return deepMerge(dst, src)
}

func deepMerge(dst *Tree, src *Tree) error {
outer:
	for _, sChild := range src.Children {
		for _, dChild := range dst.Children {

			if dChild.Name == sChild.Name {
				merge(sChild, dChild)
				continue outer
			}
		}

		dst.Children = append(dst.Children, sChild)
	}

	return nil
}

// Flattens homogeneous lists of primitive types and squashes lists of struct
// into one struct. If fields have conflicting types while squashing a
// list of struct, the offending field is converted to the empty interface.
func (t *Tree) Normalize() {
	// Normalize from the bottom up so use depth first iteration.
	for idx := range t.Children {
		t.Children[idx].Normalize()
	}

	var prevChild *Tree

	newChildren := []*Tree{}
	for _, child := range t.Children {
		if prevChild != nil && prevChild.Name == child.Name {
			// merge new types for missing fields
			merge(prevChild, child)
			prevChild.List = true
			continue
		}

		newChildren = append(newChildren, child)
		prevChild = child
	}

	t.Children = newChildren
}

func init() {
	log.SetFlags(log.Lshortfile)

	if err := config.Parse(); err != nil {
		log.Fatal("Error parsing flags:", err)
	}
}

func main() {
	defer config.Close()

	xmlDecoder := xml.NewDecoder(config.inputFile)
	xmlDecoder.CharsetReader = charset.NewReader

	root := Tree{parent: nil, Type: Struct, List: false}
	root.Populate(xmlDecoder)

	var data interface{}
	err := xmlDecoder.Decode(&data)

	if err == io.EOF {
	} else if err != nil {
		log.Fatal("Error decoding input: ", err)
	}

	if config.normalize {
		root.Normalize()
	}

	indented, err := xml.MarshalIndent(root, "", "\t")
	if err != nil {
		log.Fatal("Error encoding root:", err)
	}

	_, err = config.dumpFile.Write(indented)
	if err != nil {
		log.Fatal("Error dumping root:", err)
	}

	source, err := root.Format()
	fmt.Println(string(source))
	if err != nil {
		log.Fatal("Error formatting source:", err)
	}
}
