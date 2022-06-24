package utils

import (
	"bytes"
	"encoding/gob"
	"reflect"
)

func IsEthereumAddress(address string) bool {
	if len(address) != 42 {
		return false
	}
	if address[0] != '0' || address[1] != 'x' {
		return false
	}
	return true
}

func DeepCopy(src, dist interface{}) (err error) {
	buf := bytes.Buffer{}
	if err = gob.NewEncoder(&buf).Encode(src); err != nil {
		return
	}
	return gob.NewDecoder(&buf).Decode(dist)
}

func GetUnderlyingAsValue(data interface{}) reflect.Value {
	return reflect.ValueOf(data)
}
