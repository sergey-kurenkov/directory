package internal

import "github.com/satori/go.uuid"

type CommonManager struct {
	Employee1 string
	Employee2 string
	Manager   string
}

type Directory interface {
	FindClosestCommonManager(employee1, employee2 string) []CommonManager
}

func NewDirectory(top *OrgUnit) Directory {
	return newDirectory(top)
}

type Employee struct {
	Name string `json:"name"`
	UUID uuid.UUID `json:"-"`
}

type Manager struct {
	Employee
}

type OrgUnit struct {
	Name     string      `json:"org-unit-name"`
	Manager  Manager     `json:"manager"`
	Reports  []*Employee `json:"reports"`
	OrgUnits []*OrgUnit  `json:"org-units"`
}
