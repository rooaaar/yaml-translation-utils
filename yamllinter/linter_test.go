package main

import (
	"reflect"
	"testing"

	"gopkg.in/yaml.v2"
)

func TestAlphabeticSort(t *testing.T) {
	interfaces := []struct {
		yamlStr string
		errs    []yamlError
	}{
		{
			"",
			nil,
		},
		{
			`
a: a
b: b
c: c
`,
			nil,
		},
		{
			`
a:
  aa: aa
  bb: bb
  dd:
    ddaa: ddaa
    ddbb: ddbb
b: b
c: c
`,
			nil,
		},
		{
			`
a:
  aa: aa
  ab: ab
  ac: ac
b:
  ba: ba
  bb: bb
  bc: bc
c: c
`,
			nil,
		},
		{
			`
b: b
a: a
c: c
`,
			[]yamlError{
				yamlError{notSorted, []string{"a"}},
			},
		},
		{
			`
b: b
c: c
a: a
`,
			[]yamlError{
				yamlError{notSorted, []string{"a"}},
			},
		},
		{
			`
x: x
a: a
b: b
c: c
`,
			[]yamlError{
				yamlError{notSorted, []string{"a"}},
				yamlError{notSorted, []string{"b"}},
				yamlError{notSorted, []string{"c"}},
			},
		},
		{
			`
a:
  aa: aa
  ac: ac
  ab: ab
b: b
c: c
`,
			[]yamlError{
				yamlError{notSorted, []string{"a", "ab"}},
			},
		},
		{
			`
a:
  ax: ax
  aa: aa
  ab: ab
  ac: ac
b: b
c: c
`,
			[]yamlError{
				yamlError{notSorted, []string{"a", "aa"}},
				yamlError{notSorted, []string{"a", "ab"}},
				yamlError{notSorted, []string{"a", "ac"}},
			},
		},
		{
			`
a:
  ax: ax
  aa: aa
  ab: ab
  ac: ac
b: b
c: c
`,
			[]yamlError{
				yamlError{notSorted, []string{"a", "aa"}},
				yamlError{notSorted, []string{"a", "ab"}},
				yamlError{notSorted, []string{"a", "ac"}},
			},
		},
		{
			`
a:
  aa: aa
  ab: ab
  ac: ac
b:
  bx: bx
  ba: ba
  bb: bb
  bc: bc
c: c
`,
			[]yamlError{
				yamlError{notSorted, []string{"b", "ba"}},
				yamlError{notSorted, []string{"b", "bb"}},
				yamlError{notSorted, []string{"b", "bc"}},
			},
		},
		{
			`
a:
  aa: aa
  ab: ab
  ac: ac
b:
  ba: ba
  bb:
    bbx: bbx
    bba: bba
    bbb: bbb
    bbc: bbc
  bc: bc
c: c
`,
			[]yamlError{
				yamlError{notSorted, []string{"b", "bb", "bba"}},
				yamlError{notSorted, []string{"b", "bb", "bbb"}},
				yamlError{notSorted, []string{"b", "bb", "bbc"}},
			},
		},
		{
			`
a:
  aa: aa
  ab: ab
  ac: ac
b:
  ba: ba
  bb:
    bbb: bbb
    bba: bba
    bbc: bbc
  bc: bc
c: c
`,
			[]yamlError{
				yamlError{notSorted, []string{"b", "bb", "bba"}},
			},
		},
	}

	for _, in := range interfaces {
		node := unmarshalYaml(in.yamlStr)
		errs := lint(unmarshalYaml(in.yamlStr))
		if !reflect.DeepEqual(errs, in.errs) {
			t.Errorf("lint(%v) != %v, got: %v", node, in.errs, errs)
		}
	}
}

func TestKeyErrors(t *testing.T) {
	interfaces := []struct {
		yamlStr string
		errs    []yamlError
	}{
		{
			"a.b: a.b",
			[]yamlError{
				yamlError{notAlphaNumericDashUnderline, []string{"a.b"}},
			},
		},
		{
			`
a1: a1
b1: b1
XX: XX
`,
			[]yamlError{
				yamlError{notSorted, []string{"XX"}},
				yamlError{notLower, []string{"XX"}},
			},
		},
		{
			`
a:
  a.1: a.1
  a.a: a.a
`,
			[]yamlError{
				yamlError{notAlphaNumericDashUnderline, []string{"a", "a.1"}},
				yamlError{notAlphaNumericDashUnderline, []string{"a", "a.a"}},
			},
		},
		{
			`
a:
  aa:
    aaa:
      aaaa:
        AAAAA: AAAAA
`,
			[]yamlError{
				yamlError{notLower, []string{"a", "aa", "aaa", "aaaa", "AAAAA"}},
			},
		},
	}

	for _, in := range interfaces {
		node := unmarshalYaml(in.yamlStr)
		errs := lint(unmarshalYaml(in.yamlStr))
		if !reflect.DeepEqual(errs, in.errs) {
			t.Errorf("lint(%v) != %v, got: %v", node, in.errs, errs)
		}
	}
}

func unmarshalYaml(str string) yaml.MapSlice {
	var data yaml.MapSlice
	if err := yaml.Unmarshal([]byte(str), &data); err != nil {
		panic(err)
	}

	return data
}
