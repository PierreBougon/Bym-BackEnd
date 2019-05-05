package models

type Table interface {
	Validate() (map[string]interface{}, bool)
}