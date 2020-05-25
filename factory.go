package main

import (
	"reflect"
	"fmt"
)

//type nodeFactory interface {
//	New(interface{}) node
//}
//
//type ORFactory struct {
//}
//
//func (*ORFactory) New(in interface{}) node {
//	return &OR{&StrZero}
//}
//
//type CoherentFactory struct {
//}
//
//func (*CoherentFactory) New(in interface{}) node {
//	return &Coherent{node{}}
//}
//
//type LeafFactory struct {
//}
//
//func (*LeafFactory) New(in interface{}) node {
//	return &leaf{}
//}

//func Choose(key interface{}) nodeFactory {
//	v := reflect.ValueOf(key)
//	
//	switch v.Kind() {
////	case reflect.Bool:
////		fmt.Printf("bool: %v\n", v.Bool())
////	case reflect.Int, reflect.Int8, reflect.Int32, reflect.Int64:
////		fmt.Printf("int: %v\n", v.Int())
////	case reflect.Uint, reflect.Uint8, reflect.Uint32, reflect.Uint64:
////		fmt.Printf("int: %v\n", v.Uint())
////	case reflect.Float32, reflect.Float64:
////		fmt.Printf("float: %v\n", v.Float())
//	case reflect.String:
//		if v.String() == "OR" {
//			return &ORFactory{}
//		}
////	case reflect.Slice:
////		fmt.Printf("slice: len=%d, %v\n", v.Len(), v.Interface())
////	case reflect.Map:
////		fmt.Printf("map: \n");
////		iter := reflect.ValueOf(yaml).MapRange()
////		for iter.Next() {
////			k := iter.Key()
////			fmt.Printf("[%v] ", k);
////			v := iter.Value()
////			yamlToNode(v.Interface())
////		}
////	case reflect.Chan:
////		fmt.Printf("chan %v\n", v.Interface())
//	default:
//		fmt.Printf("\n%v\n",v)
//	}
//	return nil
//}

func  BigUglySwitch(in interface{}) node {
	v := reflect.ValueOf(in)
	
	//Uintptr
	//Complex64
	//Complex128
	//Chan
	//Func
	//Interface
	//Ptr
	//Struct
	//UnsafePointer	
	switch v.Kind() {
	case reflect.Bool,
	     reflect.Int, reflect.Int8, reflect.Int32, reflect.Int64,
	     reflect.Uint, reflect.Uint8, reflect.Uint32, reflect.Uint64,
	     reflect.Float32, reflect.Float64,
	     reflect.String:
		return &leaf{v}
	case reflect.Slice, reflect.Array: {
		var a []node
		for i := 0; i < v.Len(); i++ {
			elem := v.Index(i)
			n := BigUglySwitch(elem.Interface())		
			a = append(a, n)		
			//fmt.Printf("%d: %s %s = %v\n", i,
			//	typeOfT.Field(i).Name, f.Type(), f.Interface())
		}
		return &nArray{a}
		//fmt.Printf("slice: len=%d, %v\n", v.Len(), v.Interface())
	}
	case reflect.Map:
		fmt.Printf("map: \n");
		iter := v.MapRange()
		for iter.Next() {
			k := iter.Key()
			v := iter.Value()
			fmt.Printf("[%v] ", k);
			if k.Kind() == reflect.String {
				switch k.String() {
				case "OR":
					return &OR{BigUglySwitch(v.Interface())}
				case "Coherent":
					return &Coherent{BigUglySwitch(v.Interface())}
				}
			}
			//v := iter.Value()
			//yamlToNode(v.Interface())
			
		}
	default:
		fmt.Printf("\n%v\n",v)
	}
	return nil
}
