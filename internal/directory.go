package internal

type Directory interface {
	FindClosestCommonManager(employee1, employee2 string) (*Manager, error)
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
	Manager  Manager     `json:"manager"`
	Reports  []*Employee `json:"reports"`
	OrgUnits []*OrgUnit  `json:"org-units"`
}
