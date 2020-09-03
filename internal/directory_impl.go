package internal

import "fmt"

type directoryImpl struct {
	top *OrgUnit
}

func (this *directoryImpl) FindClosestCommonManager(employee1, employee2 string) (*Manager, error) {
	managersOfEmpl1, err := this.findManagers(employee1)
	if err != nil {
		return nil, fmt.Errorf("FindClosestCommonManager: %w", err)
	}

	managersOfEmpl2, err := this.findManagers(employee2)
	if err != nil {
		return nil, fmt.Errorf("FindClosestCommonManager: %w", err)
	}

	maxLen := len(managersOfEmpl1)
	if len(managersOfEmpl2) < maxLen {
		maxLen = len(managersOfEmpl2)
	}

	if maxLen == 0 {
		return nil, errNoCommonManager
	}

	i := 0
	for ; i < maxLen; i++ {
		if managersOfEmpl1[i] != managersOfEmpl2[i] {
			if i == 0 {
				return nil, errNoCommonManager
			}

			break
		}
	}

	commonManager := *managersOfEmpl1[i-1]

	return &commonManager, nil
}

type orgUnit2Traverse struct {
	orgUnit  *OrgUnit
	managers []*Manager
}

func (this *directoryImpl) findManagers(employeeName string) ([]*Manager, error) {
	items := []*orgUnit2Traverse{{this.top, []*Manager{}}}
	for len(items) > 0 {
		currentItem := items[len(items)-1]
		items = items[:len(items)-1]

		if currentItem.orgUnit.Manager.Name == employeeName {
			return currentItem.managers, nil
		}

		for _, employee := range currentItem.orgUnit.Reports {
			if employee.Name == employeeName {
				currentItem.managers = append(currentItem.managers, &currentItem.orgUnit.Manager)
				return currentItem.managers, nil
			}
		}

		currentItem.managers = append(currentItem.managers, &currentItem.orgUnit.Manager)
		for _, orgUnit := range currentItem.orgUnit.OrgUnits {
			items = append(items, &orgUnit2Traverse{orgUnit, currentItem.managers})
		}
	}

	return nil, fmt.Errorf("%w: %v", errEmployeeNotFound, employeeName)
}
