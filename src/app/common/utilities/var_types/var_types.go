package var_types

// GetIfInterfaceIsIntOrString checks if an interface is either a string or an int.
// Lol, do I I have to explain that?
func GetIfInterfaceIsIntOrString(i interface{}) string {
	var interfaceType string

	switch i.(type) {
	case int:
		interfaceType = "int"
	case string:
		interfaceType = "string"
	case float64:
		interfaceType = "float64"
	default:
		interfaceType = "unknown"
	}

	return interfaceType
}

// GetInterfaceAsIntOrString -
func GetInterfaceAsIntOrString() {

}