package interface_utils

import (
	"errors"
	"reflect"
	"strconv"
)


// GetJSONValueIfInt returns an interface if it is an integer
func GetJSONValueIfInt(i interface{}) (int, error) {
	// Get value as string first
	strVar, err := GetJSONValueIfString(i)
	if err != nil {
		return 0, err
	}
	// Then convert to int
	intVal, err := strconv.Atoi(strVar)
	if err != nil {
		return 0, err
	}
	return intVal, nil
}

// GetJSONValueIfString returns the string value of an interface if
// a valid string
func GetValueIfInt(i interface{}) (string, error) {
	var interfaceVal string
	var err error
	if reflect.TypeOf(i).String() != "string" {
		err = errors.New("interface is not a string")
	} else {
		interfaceVal = i.(string)
	}
	return interfaceVal, err
}

// GetJSONValueIfString returns the string value of an interface if
// a valid string
func GetJSONValueIfString(i interface{}) (string, error) {
	var interfaceVal string
	var err error
	if reflect.TypeOf(i).String() != "string" {
		err = errors.New("interface is not a string")
	} else {
		interfaceVal = i.(string)
	}
	return interfaceVal, err
}