// JSONGen - A tool for generating native Golang types from JSON objects.
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
	"bytes"
	"encoding/xml"
	"reflect"
	"testing"
)

func Parse(raw string) (tree Tree, err error) {
	buf := bytes.NewBufferString(raw)
	xmlDecoder := xml.NewDecoder(buf)

	tree.Populate(xmlDecoder)
	tree.Normalize()

	return
}

type TreeTestCase struct {
	Source string
	Tree   Tree
}

func (tc TreeTestCase) TestTree(t *testing.T) {
	tree, err := Parse(tc.Source)

	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(tree, tc.Tree) {
		t.Errorf("Expected: %+v Got: %#v", tc.Tree, tree)
	}

	formatted, err := tree.Format()
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%s\n", tc.Source)
	t.Logf("%s\n", string(formatted))
}

func TestString(t *testing.T) {
	testCases := []TreeTestCase{}

	for _, testCase := range testCases {
		testCase.TestTree(t)
	}
}

func TestStringList(t *testing.T) {
	testCases := []TreeTestCase{}

	for _, testCase := range testCases {
		testCase.TestTree(t)
	}
}

func TestStruct(t *testing.T) {
	testCases := []TreeTestCase{}

	for _, testCase := range testCases {
		testCase.TestTree(t)
	}
}

func (tc TreeTestCase) TestFormat(t *testing.T) {
	source, err := tc.Tree.Format()

	if err != nil {
		t.Fatal(err)
	}

	if string(source) != tc.Source {
		t.Errorf("Expected: %q Got: %q", tc.Source, source)
	}
}

type SanitizerTestCase struct {
	Source, Sanitized string
	TitleCase         bool
}

func TestSanitizier(t *testing.T) {
	testCases := []SanitizerTestCase{
		{"Sanitary", "Sanitary", true},
		{"sanitary", "Sanitary", true},
		{"_Sanitary", "Sanitary", true},
		{"_sanitary", "Sanitary", true},
		{"Sanitary", "Sanitary", false},
		{"sanitary", "Sanitary", false},
		{"_Sanitary", "Sanitary", false},
		{"_sanitary", "Sanitary", false},

		{"titlecase", "Titlecase", true},
		{"title-case", "TitleCase", true},
		{"title_case", "TitleCase", true},
		{"title case", "TitleCase", true},
		{"titlecase", "Titlecase", false},
		{"title-case", "Title_case", false},
		{"title_case", "Title_case", false},
		{"title case", "TitleCase", false},

		{"123", "_", true},
		{"123.foo", "Foo", true},
		{".foo123", "Foo123", true},
		{".foo.123", "Foo123", true},
	}

	for _, testCase := range testCases {
		sanitized := Ident(testCase.Source)
		config.titleCase = testCase.TitleCase
		if testCase.Sanitized != sanitized.String() {
			t.Fatalf("Source: %q Expected: %q Got: %q\n", testCase.Source, testCase.Sanitized, sanitized.String())
		}
	}
}
