package internal

import "fmt"

var (
	errEmployeeNotFound error = fmt.Errorf("employee has not been found")
	errNoCommonManager  error = fmt.Errorf("no common manager")
)
