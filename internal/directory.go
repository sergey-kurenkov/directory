package internal

type CommonManager struct {
	Employee1 string
	Employee2 string
	Manager   string
}

type Directory interface {
	FindClosestCommonManager(employee1, employee2 string) []CommonManager
}

func NewDirectory(top *OrgUnit) Directory {
	return &directoryImpl{
		top: top,
	}
}

type Employee struct {
	Name string `json:"name"`
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
