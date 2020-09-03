package internal

import (
	"fmt"
	"reflect"
	"strings"
)

type directoryImpl struct {
	top *OrgUnit
}

type duplicateKeys = map[string]struct{}
func (this *directoryImpl) FindClosestCommonManager(employeeName1, employeeName2 string) []CommonManager {
	allEmployees1 := this.findEmployee(employeeName1)
	allEmployees2 := this.findEmployee(employeeName2)

	commonManagers := []CommonManager{}

	duplicateKeys := duplicateKeys{}
	for _, e1 := range allEmployees1 {
		for _, e2 := range allEmployees2 {
			cm := this.findCommonManager(&e1, &e2, duplicateKeys)
			if cm != nil {
				commonManagers = append(commonManagers, *cm)
			}
		}
	}

	return commonManagers
}

type orgUnits []*OrgUnit

type orgUnit2Traverse struct {
	currentOrgUnit *OrgUnit
	parentOrgUnits orgUnits
}

func (this *directoryImpl) findEmployee(employeeName string) []foundEmployee {
	result := []foundEmployee{}
	items := []*orgUnit2Traverse{{this.top, orgUnits{}}}

	for len(items) > 0 {
		currentItem := items[len(items)-1]
		items = items[:len(items)-1]

		if currentItem.currentOrgUnit.Manager.Name == employeeName {
			e := foundEmployee{
				managersOrgUnits: orgUnits{},
				ownUnit:          orgUnits{},
				employeeName: employeeName,
			}
			e.managersOrgUnits = append(e.managersOrgUnits, currentItem.parentOrgUnits...)
			e.ownUnit = append(e.ownUnit, currentItem.parentOrgUnits...)
			e.ownUnit = append(e.ownUnit, currentItem.currentOrgUnit)
			result = append(result, e)
		}

		currentItem.parentOrgUnits = append(currentItem.parentOrgUnits, currentItem.currentOrgUnit)
		for _, employee := range currentItem.currentOrgUnit.Reports {
			if employee.Name == employeeName {
				e := foundEmployee{
					managersOrgUnits: orgUnits{},
					ownUnit:          orgUnits{},
					employeeName:     employeeName,
				}
				e.managersOrgUnits = append(e.managersOrgUnits, currentItem.parentOrgUnits...)
				e.ownUnit = append(e.ownUnit, currentItem.parentOrgUnits...)
				result = append(result, e)
			}
		}

		for _, orgUnit := range currentItem.currentOrgUnit.OrgUnits {
			items = append(items, &orgUnit2Traverse{orgUnit, currentItem.parentOrgUnits})
		}
	}

	return result
}

func (this *directoryImpl) findCommonManager(e1 *foundEmployee, e2 *foundEmployee,
	duplicateKeys duplicateKeys) *CommonManager {
	maxLen := len(e1.managersOrgUnits)

	if len(e2.managersOrgUnits) < maxLen {
		maxLen = len(e2.managersOrgUnits)
	}

	if maxLen == 0 {
		return nil
	}

	if e1.employeeName == e2.employeeName && reflect.DeepEqual(e1.ownUnit, e2.ownUnit) {
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
		Manager:   e1.makeManagerName(i-1),
	}

	const pattern4Key = "%s:%s"

	if e1.employeeName == e2.employeeName {
		if _, ok := duplicateKeys[fmt.Sprintf(pattern4Key, commonManager.Employee1, commonManager.Employee2)]; ok {
			return nil
		}

		duplicateKeys[fmt.Sprintf(pattern4Key, commonManager.Employee1, commonManager.Employee2)] = struct{}{}
		duplicateKeys[fmt.Sprintf(pattern4Key, commonManager.Employee2, commonManager.Employee1)] = struct{}{}
	}

	return commonManager
}

type foundEmployee struct {
	managersOrgUnits orgUnits
	ownUnit          orgUnits
	employeeName     string
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
