package mapper

type ErrorMappingNotFound struct {
	iface_name string
}

func (e *ErrorMappingNotFound) Error() string {
	return "ERROR: mapping not found"
}
