package validation_utils

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/dembygenesis/quiz_maker_auth/src/app/common/utilities/date"
	"reflect"
	"strings"
)

var validate *validator.Validate

func init() {
	configValidate()
}

// configValidate adds custom validator rules, and custom tag name returns
func configValidate() {
	validate = validator.New()
	// ======================================================
	// Add custom validators
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	_ = validate.RegisterValidation("date_format", func(fl validator.FieldLevel) bool {

		// Ignore
		if requiredExists(&fl) {
			// Ensure "date_format" has a valid date format
			return date.ValidDate(fl.Field().String()) == true
		} else {
			// Don't bother enforcing "date_format" if the field is not required
			return false
		}
	})
}

// requiredExists checks if the field has the option "required" specified
// (e.g) date_format is required to be followed
func requiredExists(fl *validator.FieldLevel) bool {
	p := (*fl).Parent()
	sf := (*fl).StructFieldName()
	sit := reflect.Indirect(reflect.ValueOf(p.Interface())).Type()
	st, found := sit.FieldByName(sf)
	if found == false {
		return false
	}
	tags := st.Tag.Get("validate")
	if strings.Contains(tags, "required") {
		return true
	}
	return false
}

// ValidateStructParams2 validates the struct validation rules
// provided in the tag using the "validator" library
func ValidateStructParams2(p interface{}) ([]string, error) {
	var missingParams []string
	if reflect.ValueOf(p).Kind().String() == "ptr" {
		structType := reflect.Indirect(reflect.ValueOf(p)).Kind().String()
		if structType == "slice" {
			return missingParams, nil
		}
		if structType != "struct" {
			return missingParams, errors.New("parameter passed is not a struct: " + structType)
		}
	} else if reflect.ValueOf(p).Kind() != reflect.Struct {
		return missingParams, errors.New("parameter passed is not a struct")
	}

	err := validate.Struct(p)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		for _, err := range validationErrors {
			strToTrim := `Key: '` + err.Namespace() + `' Error:Field validation for '` + err.Field() + "'"
			trimmedMsg := strings.Replace(err.Error(), strToTrim, err.Field(), -1)
			missingParams = append(missingParams, trimmedMsg)
		}
	}
	return missingParams, nil
}

// ValidateStructParams - scans a struct and returns errors if required fields are empty
func ValidateStructParams(p interface{}) []string {
	var missingParameters []string
	s := p
	v := reflect.ValueOf(s)
	numberOfFields := reflect.Indirect(v).NumField()
	v2 := reflect.TypeOf(s)
	typeOfS := v.Type()
	for i := 0; i < numberOfFields; i++ {
		propertyType := typeOfS.Field(i).Type.String()
		propertyValue := v.Field(i).Interface()
		propertyName := v2.Field(i).Tag.Get("json")
		required := v2.Field(i).Tag.Get("required")
		if required != "false" {
			if propertyType == "string" {
				if propertyValue == "" {
					missingParameters = append(missingParameters, ``+propertyName+` empty`)
				}
			} else if propertyType == "float64" {
				if propertyValue == 0 {
					missingParameters = append(missingParameters, ``+propertyName+` empty`)
				}
			}
		}
	}
	return missingParameters
}
