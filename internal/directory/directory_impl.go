package directory

import (
	"fmt"
	"strings"

	uuid "github.com/satori/go.uuid"
)

type directoryImpl struct {
	top *OrgUnit
}

func newDirectory(top *OrgUnit) Directory {
	items := []*OrgUnit{top}

	for len(items) > 0 {
		currentItem := items[len(items)-1]
		items = items[:len(items)-1]

		currentItem.Manager.UUID = uuid.NewV4()
		for _, employee := range currentItem.Reports {
			employee.UUID = uuid.NewV4()
		}

		items = append(items, currentItem.OrgUnits...)
	}

	return &directoryImpl{top: top}
}

type duplicateKeys = map[string]struct{}

func (this *directoryImpl) FindClosestCommonManager(employeeName1, employeeName2 string) ([]CommonManager, error) {
	allEmployees1, err := this.findEmployee(employeeName1)
	if err != nil {
		return nil, err
	}

	allEmployees2, err := this.findEmployee(employeeName2)
	if err != nil {
		return nil, nil
	}

	commonManagers := []CommonManager{}

	duplKeys := duplicateKeys{}

	for _, e1 := range allEmployees1 {
		for _, e2 := range allEmployees2 {
			cm := this.findCommonManager(e1, e2, duplKeys)
			if cm != nil {
				commonManagers = append(commonManagers, *cm)
			}
		}
	}

	return commonManagers, nil
}

type orgUnits []*OrgUnit

type orgUnit2Traverse struct {
	currentOrgUnit *OrgUnit
	parentOrgUnits orgUnits
}

func (this *directoryImpl) findEmployee(employeeName string) ([]*foundEmployee, error) {
	result := []*foundEmployee{}
	items := []*orgUnit2Traverse{{this.top, orgUnits{}}}

	for len(items) > 0 {
		currentItem := items[len(items)-1]
		items = items[:len(items)-1]

		if currentItem.currentOrgUnit.Manager.Name == employeeName {
			e := foundEmployee{
				managersOrgUnits: orgUnits{},
				ownUnit:          orgUnits{},
				employeeName:     currentItem.currentOrgUnit.Manager.Name,
				employeeUUID:     currentItem.currentOrgUnit.Manager.UUID,
			}
			e.managersOrgUnits = append(e.managersOrgUnits, currentItem.parentOrgUnits...)
			e.ownUnit = append(e.ownUnit, currentItem.parentOrgUnits...)
			e.ownUnit = append(e.ownUnit, currentItem.currentOrgUnit)
			result = append(result, &e)
		}

		currentItem.parentOrgUnits = append(currentItem.parentOrgUnits, currentItem.currentOrgUnit)
		for _, employee := range currentItem.currentOrgUnit.Reports {
			if employee.Name == employeeName {
				e := foundEmployee{
					managersOrgUnits: orgUnits{},
					ownUnit:          orgUnits{},
					employeeName:     employeeName,
					employeeUUID:     employee.UUID,
				}

				e.managersOrgUnits = append(e.managersOrgUnits, currentItem.parentOrgUnits...)
				e.ownUnit = append(e.ownUnit, currentItem.parentOrgUnits...)
				result = append(result, &e)
			}
		}

		for _, orgUnit := range currentItem.currentOrgUnit.OrgUnits {
			items = append(items, &orgUnit2Traverse{orgUnit, currentItem.parentOrgUnits})
		}
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("%w: %v", errNoEmployee, employeeName)
	}

	return result, nil
}

func (this *directoryImpl) findCommonManager(e1, e2 *foundEmployee, duplicateKeys duplicateKeys) *CommonManager {
	maxLen := len(e1.managersOrgUnits)

	if len(e2.managersOrgUnits) < maxLen {
		maxLen = len(e2.managersOrgUnits)
	}

	if maxLen == 0 {
		return nil
	}

	if e1.employeeUUID == e2.employeeUUID {
		return nil
	}

	i := 0
	for ; i < maxLen; i++ {
		if e1.managersOrgUnits[i] != e2.managersOrgUnits[i] {
			if i == 0 {
				return nil
			}

			break
		}
	}

	commonManager := &CommonManager{
		Employee1: e1.makeFullEmployeeName(),
		Employee2: e2.makeFullEmployeeName(),
		Manager:   e1.makeManagerName(i - 1),
	}

	const pattern4Key = "%s:%s"

	if e1.employeeName == e2.employeeName {
		if _, ok := duplicateKeys[fmt.Sprintf(pattern4Key, e1.employeeUUID, e2.employeeUUID)]; ok {
			return nil
		}

		duplicateKeys[fmt.Sprintf(pattern4Key, e1.employeeUUID, e2.employeeUUID)] = struct{}{}
		duplicateKeys[fmt.Sprintf(pattern4Key, e2.employeeUUID, e1.employeeUUID)] = struct{}{}
	}

	return commonManager
}

type foundEmployee struct {
	managersOrgUnits orgUnits
	ownUnit          orgUnits
	employeeName     string
	employeeUUID     uuid.UUID
}

func (this *foundEmployee) makeFullEmployeeName() string {
	var fullName strings.Builder

	for _, orgUnit := range this.ownUnit {
		fullName.WriteString("/")
		fullName.WriteString(orgUnit.Name)
	}

	fullName.WriteString("/")
	fullName.WriteString(this.employeeName)

	return fullName.String()
}

func (this *foundEmployee) makeManagerName(maxIndex int) string {
	var fullName strings.Builder

	for i := 0; i <= maxIndex; i++ {
		orgUnit := this.managersOrgUnits[i]

		fullName.WriteString("/")
		fullName.WriteString(orgUnit.Name)
	}

	fullName.WriteString("/")
	fullName.WriteString(this.managersOrgUnits[maxIndex].Manager.Name)

	return fullName.String()
}
