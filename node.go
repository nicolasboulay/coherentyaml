package main

import (
	"errors"
	"log"
	"fmt"
	"reflect"
)

type node interface {
	GetChild() []node
	IsCoherent() error
	IsCoherentWith(n node) error
	String() string
	//New(interface{}) node
}

type OR struct {
	child node
}

func (or *OR) GetChild() []node {
	a, ok := or.child.(*nArray);
	if ok {
		return a.value
	}	
	return []node{or.child} 
}

func (or *OR) IsCoherent() error {
	children := or.GetChild()
	for _, child := range children {
		err := child.IsCoherent() 
		if (err != nil) {
			return err
		}
	}
	return nil
}

func (or *OR) IsCoherentWith(n node) error {
	children := or.GetChild()
	var err error 
	for _, child := range children {
		err = child.IsCoherentWith(n) 
		if (err == nil) {
			return nil
		}
	}
	for _, child := range children {
		err = n.IsCoherentWith(child) 
		if (err == nil) {
			return nil
		}
	}
	return err
}

func (o *OR) String() string {
	return fmt.Sprintf("OR{%v}", o.GetChild())
}


type Coherent struct {
	child node
}

func (c *Coherent) GetChild() []node {
	a, ok := c.child.(*nArray);
	if ok {
		return a.value
	}	
	return []node{c.child} 
}

func (c *Coherent) IsCoherent() error {
	children := c.GetChild()
	for _, child := range children {
		for _, child2 := range children {
			err := child.IsCoherentWith(child2) 
			if (err != nil) {
				return child2.IsCoherentWith(child)
			}
		}
	}
	return nil
}

func (c *Coherent) IsCoherentWith(n node) error {
	//log.Printf("Coherent %v %v ?",c,n)
	children := c.GetChild()
	var err error 
	for _, child := range children {
		err = child.IsCoherentWith(n) 
		if (err != nil) {
			return err
		}
	}
	return nil
}

func (c *Coherent) String() string {
	return fmt.Sprintf("Coherent{%v}", c.GetChild())
}

type Not struct {
	child node
}

func (n *Not) GetChild() []node {
	a, ok := n.child.(*nArray);
	if ok {
		return a.value
	}	
	return []node{n.child} 
}

func (n *Not) IsCoherent() error {
	children := n.GetChild()
	if (len(children)) != 1 {
		return errors.New("'Not' could have only one child")
	}
	
	err := children[0].IsCoherent() 
	if (err != nil) {
		return err
	}
	
	return nil
}

func (n *Not) IsCoherentWith(o node) error {
	
	if (len(n.GetChild())) != 1 {
		return errors.New("'Not' could have only one child")
	}
	var err error
	//log.Printf("%v",n.child[0])
	//log.Printf("%v",o)
	err = n.GetChild()[0].IsCoherentWith(o)
	if (err != nil) {
		//log.Print(nil)
		return nil
	}
	//log.Print("not : Both node should be different")
	return errors.New("Both node should be different") //TODO: printing and referencing the yaml text
}

func (n *Not) String() string {
	return fmt.Sprintf("Not{%v}", n.GetChild())
}

type Str string

var StrZero Str = Str("")

func (s *Str) GetChild() []node {
	return []node{}
} 

func (s *Str) IsCoherent() error {
	return nil
}

func (s *Str) IsCoherentWith(n node) error {
	s2, ok := n.(*Str);
	if (!ok) {
		//case with OR/Not/Coherency/.. in between
		return n.IsCoherentWith(s)
	}
	if (s2 == s) {
		return nil
	} else if *s == StrZero || *s2 == StrZero {
		return nil
	}
	
	return errors.New("String shall be coherent")
}


func (s *Str) String() string {
	return fmt.Sprintf("str{%s}",*s)
}

type leaf struct {
	value reflect.Value
}

func (l * leaf) GetChild() []node {
	return []node{}
}

func (l *leaf) IsCoherent() error {
	return nil
}

func (l *leaf) IsCoherentWith(n node) error {
	l2, ok := n.(*leaf);
	if (!ok) {
		//case with OR/Not/Coherency/.. in between
		return n.IsCoherentWith(l)
	}

	if (l2.value.Interface() == l.value.Interface()) {
		return nil
	}
	equalKind := l2.value.Kind() == l.value.Kind() 
	if (equalKind) {
		if(l.isNeutral() || l2.isNeutral()) {
			return nil
		}
	}
	
	return fmt.Errorf("Incoherent leaf %v vs %v (%s vs %s)", l, l2, l.value.Kind(), l2.value.Kind())
}

func (l *leaf) isNeutral() bool {
	i := l.value
	switch i.Kind() {
	case reflect.Bool:
		return false
	case reflect.Int, reflect.Int8, reflect.Int32, reflect.Int64:
		return i.Int() == -1 
	case reflect.Uint, reflect.Uint8, reflect.Uint32, reflect.Uint64:
		return i.Uint() == 1 
	case reflect.Float32, reflect.Float64:
		return i.Float() == 1.0 
	case reflect.String:
		return i.String() == "" 
	default: {
		fmt.Printf("isNeutral unknwon: %v\n",i)
		return false
	}
	}
	return false
}

func (l *leaf) String() string {
	return fmt.Sprintf("%v", l.value)
}

type nStruct struct {
	value map[node]node
}

func (n *nStruct) GetChild() []node {
	var v []node
	for k := range n.value {
		v = append(v, k)
	}
	return v
}

func (n *nStruct) IsCoherent() error {
	c := n.GetChild()
	for _,node := range c {
		err := node.IsCoherent()
		if (err != nil) {
			return err
		}
	}
	return nil
}

func (n *nStruct) IsCoherentWith(n2 node) error {//TODO ! c'est complètement faux
	c := n.GetChild()
	for _,k := range c {
		err := k.IsCoherentWith(n2)
		if (err != nil) {
			return err
		}
	}
	return nil
}

func (n *nStruct) String() string {
	return fmt.Sprintf("%v", n.value)
}

type nArray struct {
	value []node
}

func (a *nArray) GetChild() []node {
	return a.value
}

func (a *nArray) IsCoherent() error {
	c := a.GetChild()
	for _,node := range c {
		err := node.IsCoherent()
		if (err != nil) {
			return err
		}
	}
	return nil
}

// array is coherent with an other array
// Array are not ordered : so an element must be coherent with an other element in the other array, symmetricaly
// multiplicity are not defined, 
// 
func (a *nArray) IsCoherentWith(n2 node) error {
	a2, ok := n2.(*nArray)
	if !ok {
		return fmt.Errorf("Array needed : %v vs %v", a, n2)
	}
	c  :=  a.GetChild()
	c2 := a2.GetChild()
	for _,k := range c {
		ok := false
		for _,k2 := range c2 {
			err := k2.IsCoherentWith(k)
			log.Print(err)
			if (err == nil) {
				ok = true
				break
			}
		}
		if (!ok) {
			return errors.New("'Array' value should match without order")
		}

	}
	return nil
}

func (n *nArray) String() string {
	return fmt.Sprintf("%v", n.value)
}

//func yamlToNode(yaml interface{}) node{
//	v := reflect.ValueOf(yaml)
//	
//	switch v.Kind() {
//	case reflect.Bool:
//		fmt.Printf("bool: %v\n", v.Bool())
//	case reflect.Int, reflect.Int8, reflect.Int32, reflect.Int64:
//		fmt.Printf("int: %v\n", v.Int())
//	case reflect.Uint, reflect.Uint8, reflect.Uint32, reflect.Uint64:
//		fmt.Printf("int: %v\n", v.Uint())
//	case reflect.Float32, reflect.Float64:
//		fmt.Printf("float: %v\n", v.Float())
//	case reflect.String:
//		fmt.Printf("string: %v\n", v.String())
//	case reflect.Slice:
//		fmt.Printf("slice: len=%d, %v\n", v.Len(), v.Interface())
//	case reflect.Map:
//		fmt.Printf("map: \n");
//		iter := reflect.ValueOf(yaml).MapRange()
//		for iter.Next() {
//			k := iter.Key()
//			fmt.Printf("[%v] ", k);
//			v := iter.Value()
//			yamlToNode(v.Interface())
//		}
//	case reflect.Chan:
//		fmt.Printf("chan %v\n", v.Interface())
//	default:
//		fmt.Printf("\n%v\n",v)
//	}
//	return nil
//}

