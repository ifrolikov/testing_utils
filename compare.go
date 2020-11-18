package testing_utils

import (
	"fmt"
	"reflect"
)

func Equal(x1 interface{}, x2 interface{}) (bool, string) {
	x1CanBeNil, x2CanBeNil := false, false
	switch reflect.TypeOf(x1).Kind() {
		case reflect.Chan, reflect.Func, reflect.Map,
		reflect.Ptr, reflect.UnsafePointer, reflect.Interface, reflect.Slice:
			x1CanBeNil = true
	}
	switch reflect.TypeOf(x2).Kind() {
	case reflect.Chan, reflect.Func, reflect.Map,
		reflect.Ptr, reflect.UnsafePointer, reflect.Interface, reflect.Slice:
		x2CanBeNil = true
	}
	if (x1CanBeNil && !x2CanBeNil) || (!x1CanBeNil && x2CanBeNil) {
		return false, fmt.Sprintf("one of value can be nil\nx1: %+v\nx2: %+v", x1, x2)
	}

	if x1CanBeNil && x2CanBeNil {
		if reflect.ValueOf(x1).IsNil() && reflect.ValueOf(x2).IsNil() {
			return true, ""
		} else {
			if reflect.ValueOf(x1).IsNil() || reflect.ValueOf(x2).IsNil() {
				return false, fmt.Sprintf("one of value is nil\nx1: %+v\nx2: %+v", x1, x2)
			}
		}
	}

	if (reflect.TypeOf(x1).Kind() == reflect.Ptr &&
		reflect.TypeOf(x2).Kind() != reflect.Ptr) ||
		(reflect.TypeOf(x1).Kind() != reflect.Ptr &&
			reflect.TypeOf(x2).Kind() == reflect.Ptr) {
		return false, fmt.Sprintf("pointers not equals:\nx1: %+v\nx2: %+v", x1, x2)
	}

	x1Ptr := getPtr(x1)
	x2Ptr := getPtr(x2)

	x1Kind := reflect.TypeOf(
		reflect.ValueOf(x1Ptr).Elem().Interface(),
	).Kind()
	x2Kind := reflect.TypeOf(
		reflect.ValueOf(x2Ptr).Elem().Interface(),
	).Kind()
	if x1Kind != x2Kind {
		return false, fmt.Sprintf(
			"kind of value not equals:\nx1: %s\nx2: %s\nx1: %+v\nx2: %+v",
			x1Kind,
			x2Kind,
			x1,
			x2,
		)
	}

	switch x1Kind {
	case reflect.Struct:
		var ok bool
		var desc string
		if ok, desc = compareStruct(x1Ptr, x2Ptr); ok {
			ok, desc = compareStruct(x2Ptr, x1Ptr)
		}
		return ok, desc
	case reflect.Map:
		return compareMap(x1Ptr, x2Ptr)
	case reflect.Slice:
		fallthrough
	case reflect.Array:
		return compareSlice(x1Ptr, x2Ptr)
	case reflect.Interface:
		return false, fmt.Sprintf("unsupported kind %s", x1Kind)
	default:
		return compareScalar(x1Ptr, x2Ptr)
	}
}

func compareStruct(x1Ptr interface{}, x2Ptr interface{}) (bool, string) {
	elem1 := reflect.ValueOf(reflect.ValueOf(x1Ptr).Elem().Interface())
	elem2 := reflect.ValueOf(reflect.ValueOf(x2Ptr).Elem().Interface())

	for i := 0; i < elem1.NumField(); i++ {
		if field1 := elem1.Field(i); !field1.CanInterface() {
			continue
		}
		field1Name := elem1.Type().Field(i).Name
		equalElem2Index := -1
		for j := 0; j < elem2.NumField(); j++ {
			if field2 := elem1.Field(i); !field2.CanInterface() {
				continue
			}
			field2Name := elem2.Type().Field(j).Name
			if field2Name == field1Name {
				equalElem2Index = j
				break
			}
		}
		if equalElem2Index == -1 {
			return false, fmt.Sprintf(
				"field %s not found in x2:\nx1: %+v\nx2: %+v",
				field1Name,
				x1Ptr,
				x2Ptr,
			)
		}
		field1 := elem1.Field(i).Interface()
		field2 := elem2.Field(equalElem2Index).Interface()
		if ok, desc := Equal(field1, field2); !ok {
			return false, fmt.Sprintf("field %s is not equal with some element in x2 because %s\nx1: %+v\nx2: %+v\n",
				field1Name, desc, x1Ptr, x2Ptr)
		}
	}
	return true, ""
}

func compareMap(x1Ptr interface{}, x2Ptr interface{}) (bool, string) {
	elem1 := reflect.ValueOf(x1Ptr).Elem()
	elem2 := reflect.ValueOf(x2Ptr).Elem()
	if elem1.Len() != elem2.Len() {
		return false, fmt.Sprintf(
			"lenghts of slices not equals:\nx1: %+v\nx2: %+v",
			x1Ptr,
			x2Ptr,
		)
	}
	for _, x1Key := range elem1.MapKeys() {
		var x2Key *reflect.Value
		for _, key2 := range elem2.MapKeys() {
			if ok, _ := Equal(key2.Interface(), x1Key.Interface()); ok {
				x2Key = &key2
				break
			}
		}
		if x2Key == nil {
			return false, fmt.Sprintf("x1Key %+v not found in x2\nx1: %+v\nx2: %+v\n",
				x1Key.Interface(), x1Ptr, x2Ptr)
		}
		x1Val := elem1.MapIndex(x1Key).Interface()
		x2Val := elem2.MapIndex(*x2Key).Interface()
		if ok, desc := Equal(x1Val, x2Val); !ok {
			return false, fmt.Sprintf("x1Val %+v not found in x2 because %s\nx1: %+v\nx2: %+v\n",
				desc, x1Key.Interface(), x1Ptr, x2Ptr)
		}
	}
	return true, ""
}

func compareSlice(x1Ptr interface{}, x2Ptr interface{}) (bool, string) {
	elem1 := reflect.ValueOf(reflect.ValueOf(x1Ptr).Elem().Interface())
	elem2 := reflect.ValueOf(reflect.ValueOf(x2Ptr).Elem().Interface())

	if elem2.Len() != elem2.Len() {
		return false, fmt.Sprintf(
			"lenghts of slices not equals:\nx1: %+v\nx2: %+v",
			x1Ptr,
			x2Ptr,
		)
	}

	for i := 0; i < elem1.Len(); i++ {
		val1 := elem1.Index(i)
		isVal1FoundInX2 := false
		for j := 0; j < elem2.Len(); j++ {
			val2 := elem2.Index(j)
			if isEqual, _ := Equal(val1.Interface(), val2.Interface()); isEqual {
				isVal1FoundInX2 = true
				break
			}
		}
		if isVal1FoundInX2 == false {
			return false, fmt.Sprintf("val1 '%+v' not found in x2\nx1: %+v\nx2: %+v\n",
				val1, elem1.Interface(), elem2.Interface())
		}
	}
	return true, ""
}

func compareScalar(x1Ptr interface{}, x2Ptr interface{}) (bool, string) {
	val1 := reflect.ValueOf(x1Ptr).Elem().Interface()
	val2 := reflect.ValueOf(x2Ptr).Elem().Interface()
	if val1 != val2 {
		return false, fmt.Sprintf(
			"value not equals:\nx1: %+v\nx2: %+v",
			val1,
			val2,
		)
	}
	return true, ""
}

func getPtr(x interface{}) interface{} {
	if reflect.TypeOf(x).Kind() == reflect.Ptr {
		return x
	} else {
		return &x
	}
}
