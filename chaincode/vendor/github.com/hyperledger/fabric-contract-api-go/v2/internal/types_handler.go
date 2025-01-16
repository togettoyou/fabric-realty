// Copyright the Hyperledger Fabric contributors. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"fmt"
	"reflect"
	"sort"
	"unicode"

	"github.com/hyperledger/fabric-contract-api-go/v2/internal/types"
	"github.com/hyperledger/fabric-contract-api-go/v2/internal/utils"
)

func basicTypesAsSlice() []string {
	typesArr := []string{}

	for el := range types.BasicTypes {
		typesArr = append(typesArr, el.String())
	}
	sort.Strings(typesArr)

	return typesArr
}

func listBasicTypes() string {
	return utils.SliceAsCommaSentence(basicTypesAsSlice())
}

func arrayOfValidType(array reflect.Value, additionalTypes []reflect.Type) error {
	if array.Len() < 1 {
		return fmt.Errorf("arrays must have length greater than 0")
	}

	return typeIsValid(array.Index(0).Type(), additionalTypes, false)
}

func structOfValidType(obj reflect.Type, additionalTypes []reflect.Type) error {
	if obj.Kind() == reflect.Ptr {
		obj = obj.Elem()
	}

	for i := 0; i < obj.NumField(); i++ {
		field := obj.Field(i)

		if runes := []rune(field.Name); len(runes) > 0 && !unicode.IsUpper(runes[0]) && field.Tag.Get("metadata") == "" {
			// Skip validation for private fields, except those tagged as metadata
			continue
		}

		err := typeIsValid(field.Type, additionalTypes, false)

		if err != nil {
			return err
		}
	}

	return nil
}

func typeInSlice(a reflect.Type, list []reflect.Type) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func typeIsValid(t reflect.Type, additionalTypes []reflect.Type, allowError bool) error {
	kind := t.Kind()
	if kind == reflect.Array {
		array := reflect.New(t).Elem()
		return arrayOfValidType(array, additionalTypes)
	} else if kind == reflect.Slice {
		slice := reflect.MakeSlice(t, 1, 1)
		return typeIsValid(slice.Index(0).Type(), additionalTypes, false)
	} else if kind == reflect.Map {
		if t.Key().Kind() != reflect.String {
			return fmt.Errorf("map key type %s is not valid. Expected string", t.Key().String())
		}

		return typeIsValid(t.Elem(), additionalTypes, false)
	} else if !typeInSlice(t, additionalTypes) {
		if kind == reflect.Struct {
			additionalTypes = append(additionalTypes, t)
			additionalTypes = append(additionalTypes, reflect.PointerTo(t))
			// add self for cyclic
			return structOfValidType(t, additionalTypes)
		} else if kind == reflect.Ptr && t.Elem().Kind() == reflect.Struct {
			additionalTypes = append(additionalTypes, t)
			additionalTypes = append(additionalTypes, t.Elem())
			// add self for cyclic
			return structOfValidType(t, additionalTypes)
		} else if _, ok := types.BasicTypes[t.Kind()]; !ok || (!allowError && t == types.ErrorType) || (t.Kind() == reflect.Interface && t.String() != "interface {}" && t.String() != "error") {
			errStr := ""

			if allowError {
				errStr = " error,"
			}

			return fmt.Errorf("type %s is not valid. Expected a struct or one of the basic types%s %s or an array/slice of these", t.String(), errStr, listBasicTypes())
		}
	}

	return nil
}

func typeMatchesInterface(toMatch reflect.Type, iface reflect.Type) error {
	if iface.Kind() != reflect.Interface {
		return fmt.Errorf("type passed for interface is not an interface")
	}

	for i := 0; i < iface.NumMethod(); i++ {
		ifaceMethod := iface.Method(i)
		matchMethod, exists := toMatch.MethodByName(ifaceMethod.Name)

		if !exists {
			return fmt.Errorf("missing function %s", ifaceMethod.Name)
		}

		ifaceNumIn := ifaceMethod.Type.NumIn()
		matchNumIn := matchMethod.Type.NumIn() - 1 // skip over which the function is acting on

		if ifaceNumIn != matchNumIn {
			return fmt.Errorf("parameter mismatch in method %s. Expected %d, got %d", ifaceMethod.Name, ifaceNumIn, matchNumIn)
		}

		for j := 0; j < ifaceNumIn; j++ {
			ifaceIn := ifaceMethod.Type.In(j)
			matchIn := matchMethod.Type.In(j + 1)

			if ifaceIn.Kind() != matchIn.Kind() {
				return fmt.Errorf("parameter mismatch in method %s at parameter %d. Expected %s, got %s", ifaceMethod.Name, j, ifaceIn.Name(), matchIn.Name())
			}
		}

		ifaceNumOut := ifaceMethod.Type.NumOut()
		matchNumOut := matchMethod.Type.NumOut()
		if ifaceNumOut != matchNumOut {
			return fmt.Errorf("return mismatch in method %s. Expected %d, got %d", ifaceMethod.Name, ifaceNumOut, matchNumOut)
		}

		for j := 0; j < ifaceNumOut; j++ {
			ifaceOut := ifaceMethod.Type.Out(j)
			matchOut := matchMethod.Type.Out(j)

			if ifaceOut.Kind() != matchOut.Kind() {
				return fmt.Errorf("return mismatch in method %s at return %d. Expected %s, got %s", ifaceMethod.Name, j, ifaceOut.Name(), matchOut.Name())
			}
		}
	}

	return nil
}
