package main

import (
	"reflect"
	"fmt"
	//"log"
)

// le pattern d'abstract factory n'aide pas vraiment ici, car C'est vraiment le contenu exact de l'entré qui compte.
// Je n'ai pas trouvé de moyen d'utiliser les interfaces pour éviter ça. Le gros switch est inévitable,
// pour créer les nodes.
// 
//map et struct sont gérer de la même façon et doivent l'être, elle gène un array pour gérer les clefs particulières
// comme "OR" directement
//Kind() pas encore utilisé : 
        //Uintptr
	//Complex64
	//Complex128
	//Chan
	//Func
	//Interface
	//Ptr
	//UnsafePointer	
func  BigUglySwitch(in interface{}) node {
	v := reflect.ValueOf(in)
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
			//	typeOfT.Field(i).Name, f.Type(), f.Interface())
		}
		return &nArray{a}
	}
	case reflect.Map: {
		mapNode := make (map[node]node)
		var returnArrayNode []node
		
		iter := v.MapRange()
		for iter.Next() {
			k := iter.Key()
			v := iter.Value()
			if k.Kind() == reflect.String {
				returnArrayNode = appendArrayWithKey(returnArrayNode, mapNode, k.String(), v.Interface())
			}			
		}
		if len (mapNode) != 0 {
			returnArrayNode = append(returnArrayNode,&nStruct{mapNode})
		}
		return &nArray{returnArrayNode}
	}
	case reflect.Struct:
		mapNode := make(map[node]node)
		var returnArrayNode []node
		for i := 0; i < v.NumField(); i++ {
			if v.Field(i).CanInterface() {			
				fieldName := v.Type().Field(i).Name
				//log.Printf("^: %v %s\n",v,fieldName)
				fieldValue := v.Field(i).Interface()
				returnArrayNode = appendArrayWithKey(returnArrayNode, mapNode, fieldName, fieldValue)
			}
		}
		if len (mapNode) != 0 {
			returnArrayNode = append(returnArrayNode,&nStruct{mapNode})
		}
		return &nArray{returnArrayNode}
	default:
		fmt.Printf("\n%v\n",v)
	}
	return nil
}


func appendArrayWithKey(returnArrayNode []node, mapNode map[node]node, key string, value interface{}) []node {
	valueNode := BigUglySwitch(value) 
	switch key {
	case "OR":
		returnArrayNode = append(returnArrayNode,&OR{valueNode})
	case "Coherent":
		returnArrayNode = append(returnArrayNode,&Coherent{valueNode})
	case "Not":
		returnArrayNode = append(returnArrayNode,&Not{valueNode})
	default: {
			keyNode := Str(key)
			mapNode[keyNode] = valueNode
	}
	}
	return (returnArrayNode)
}
