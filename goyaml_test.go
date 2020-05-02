package main

import (
	"testing"
	"github.com/goccy/go-yaml"
	"io/ioutil"
	"os"
)

func TestUnmarshal(t *testing.T) {
	yml := `
%YAML 1.2
---
a: 1
b: c
`
	var v struct {
		A int
		B string
	}
	if err := yaml.Unmarshal([]byte(yml), &v); err != nil {
		t.Errorf("unmarshall error %s", err)
	}
	if (v.A != 1) {
		t.Errorf("Unmarshal v.A = %d; want 1", v.A)
	}
	if (v.B != "c") {
		t.Errorf("Unmarshal v.B = %s; want 'c'", v.B)
	}
}

func TestReadUnmarshal(t *testing.T) {
	yml, err := ioutil.ReadFile("./test/test1.yaml")
	if (err != nil) {
		t.Errorf("Readfile error : %s", err)
	} 
	var v struct {
		A int
		B string
	}
	if err := yaml.Unmarshal([]byte(yml), &v); err != nil {
		t.Errorf("Unmarshal = %s", err)
	}
	if (v.A != 1) {
		t.Errorf("Unmarshal v.A = %d; want 1", v.A)
	}
	if (v.B != "c") {
		t.Errorf("Unmarshal v.B = %s; want 'c'", v.B)
	}
}
func TestDecoder(t *testing.T) {
	yamlFile, err := os.Open("./test/test1.yaml")
	if (err != nil) {
		t.Errorf("Readfile error : %s", err)
	} 
	var v struct {
		A int
		B string
	}
	dec := yaml.NewDecoder(yamlFile)
	if err := dec.Decode(&v); err != nil {
		t.Errorf("Unmarshal = %s", err)
	}
	if (v.A != 1) {
		t.Errorf("Unmarshal v.A = %d; want 1", v.A)
	}
	if (v.B != "c") {
		t.Errorf("Unmarshal v.B = %s; want 'c'", v.B)
	}
}


func TestDecoder_InvalidCases(t *testing.T) {
	const src = `---
a:
- b
  c: d
`
	var v struct {
		A []string
	}
	err := yaml.Unmarshal([]byte(src), &v)
	if err == nil {
		t.Fatalf("expected error")
	}

	if err.Error() != yaml.FormatError(err, false, true) {
		t.Logf("err.Error() = %s", err.Error())
		t.Logf("yaml.FormatError(err, false, true) = %s", yaml.FormatError(err, false, true))
		t.Fatal(`err.Error() should match yaml.FormatError(err, false, true)`)
	}
	const ref =
`[3:3] unexpected key name
   1 | ---
   2 | a:
>  3 | - b
   4 |   c: d
        ^
`		
	t.Logf("%s", err)
	if err.Error() != ref {
		t.Errorf("Expecting \n#%s# having \n#%s#", ref, err)
	}
}

func TestUnmarshalNested(t *testing.T) {
	yml := `
%YAML 1.2
---
a: 
  aa:
    bb:
      - true
      - 0
      - 1.1
      - !!int 2
      - 3s
      - 2015-01-01
      - !!timestamp "2015-01-01"
      - ""
      - plip
      - !!str "plop"
      - !!string "plop"
      - ~
      - null
      - !expr "*.*"
b: c
`
	var v interface{}
	err := yaml.Unmarshal([]byte(yml), &v)
	t.Logf("%#v", v)

	if err != nil {
		t.Errorf("unmarshall error %s", err)
	}
}


