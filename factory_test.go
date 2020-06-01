package main

import (
	"testing"
	"log"
	"reflect"
)


//func TestFactory(t *testing.T) {
//	
//	dut := struct{
//		OR []string
//	}{[]string{"astring"}}
//
//	expected := &OR{[]node{("astring")}}
//
//	node := Choose("OR").New(dut.OR)
//	if node != expected {
//		t.Errorf("Want coherency in %#v : %#v instead of %#v",dut, node, expected)
//	}
//	log.Print(node)
//}

func TestBigUglySwitch(t *testing.T) {
	
	dut := "astring"

	expected := &leaf{reflect.ValueOf("astring")}

	node := BigUglySwitch(dut)
	if nil != node.IsCoherentWith(expected) {
		t.Errorf("Want coherency in %#v :\n%#v instead of\n %#v",dut, node, expected)
	}
	if nil != expected.IsCoherentWith(node) {
		t.Errorf("Want coherency in %#v :\n%#v instead of\n %#v",dut, node, expected)
	}
	log.Print(node)
}

func TestBigUglySwitchTable(t *testing.T) {
	testMap :=make (map[string]interface{})
	testMap["Coherent"] = []int{2,2}
	testMap["OR"] = []int{1,2}
	testMap["Not"] = 2
	testMapNode1 := &nArray{[]node{&leaf{reflect.ValueOf(2)},&leaf{reflect.ValueOf(1)}}}
	testMapCoherent  := &Coherent{testMapNode1}
	testMapOR  := &OR{testMapNode1}
	testMapNot  := &Not{&leaf{reflect.ValueOf(1)}}
	testMapNode := &nArray{[]node{testMapCoherent,testMapOR,testMapNot}}

	testStruct := struct {
		Coherent []int
		OR []int
		Not int
	} { []int{2,2}, []int{1,2}, 2,}
	
	tables := []struct{ dut interface{}; expected node;}{
		{"astring",&leaf{reflect.ValueOf("astring")}},
		{2,&leaf{reflect.ValueOf(2)}},
		{2.0,&leaf{reflect.ValueOf(2.0)}},
		{float32(2.0),&leaf{reflect.ValueOf(float32(2.0))}},
		{[]int{2},&nArray{[]node{&leaf{reflect.ValueOf(2)}}}},
		{[]int{2,3},&nArray{[]node{&leaf{reflect.ValueOf(2)},&leaf{reflect.ValueOf(3)}}}},
		{[]int{2,3},&nArray{[]node{&leaf{reflect.ValueOf(3)},&leaf{reflect.ValueOf(2)}}}},
		{[]int{-1,3},&nArray{[]node{&leaf{reflect.ValueOf(3)},&leaf{reflect.ValueOf(2)},&leaf{reflect.ValueOf(4)}}}},
		{[][]int{{2,3}},&nArray{[]node{&nArray{[]node{
			&leaf{reflect.ValueOf(2)},&leaf{reflect.ValueOf(3)}}}}}},
		{testMap, testMapNode},
		{testStruct, testMapNode},
	}

	for _, line := range tables {
		node := BigUglySwitch(line.dut)
		if nil != node.IsCoherentWith(line.expected) {
			t.Errorf("Want coherency in %#v :\n%v instead of\n %v",line.dut, node, line.expected)
		}
		if nil != line.expected.IsCoherentWith(node) {
			t.Errorf("Want coherency in %#v :\n%#v instead of\n %#v",line.dut, node, line.expected)
		}
	}
}
