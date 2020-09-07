package directory_test

import (
	. "github.com/directory/internal/directory"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

)
import "github.com/stretchr/testify/assert"

func TestOnlyCEO(t *testing.T) {
	dir := NewDirectory(
		&OrgUnit{
			Name:    "Bureaucr.at",
			Manager: Manager{Employee{Name: "Claire"}},
		})
	assert.NotNil(t, dir)

	cm, err := dir.FindClosestCommonManager("Claire", "Claire")
	require.NoError(t, err)
	assert.Equal(t, 0, len(cm))

	cm, err = dir.FindClosestCommonManager("Jonh", "Claire")
	require.Error(t, err)
	assert.Equal(t, 0, len(cm))
}

func checkTwoEmployes(t *testing.T, employee1, employee2 string, cm CommonManager) {
	expected := map[string]int{}
	expected[employee1]++
	expected[employee2]++

	received := map[string]int{}
	received[cm.Employee1]++
	received[cm.Employee2]++

	assert.True(t, reflect.DeepEqual(expected, received), expected, received)
}

func makeCommonManagerMap(cm []CommonManager) map[CommonManager]int {
	result := map[CommonManager]int{}
	for _, m := range cm {
		result[m]++
	}

	return result
}

func TestTwoOrgUnits(t *testing.T) {
	dir := NewDirectory(
		&OrgUnit{
			Name:    "Bureaucr.at",
			Manager: Manager{Employee{Name: "Claire"}},
			OrgUnits: []*OrgUnit{
				{
					Name:    "department 1",
					Manager: Manager{Employee{Name: "Bob"}},
					Reports: []*Employee{
						{Name: "Bill"},
						{Name: "John"},
						{Name: "Joseph"},
						{Name: "Kate"},
						{Name: "Monica"},
						{Name: "Jane"},
						{Name: "Jane"},
					},
				},
				{
					Name:    "department 2",
					Manager: Manager{Employee{Name: "Alice"}},
					Reports: []*Employee{
						{Name: "Fred"},
						{Name: "Donald"},
						{Name: "Bill"},
						{Name: "Monica"},
					},
				},
			},
			Reports: []*Employee{
				{Name: "Ann"},
				{Name: "Julia"},
				{Name: "John"},
				{Name: "John"},
			},
		},
	)
	assert.NotNil(t, dir)

	var cm []CommonManager
	var err error

	cm, err = dir.FindClosestCommonManager("Ann", "Julia")
	require.NoError(t, err)
	require.Equal(t, 1, len(cm))
	checkTwoEmployes(t, "/Bureaucr.at/Ann", "/Bureaucr.at/Julia", cm[0])
	assert.Equal(t, "/Bureaucr.at/Claire", cm[0].Manager)

	cm, err = dir.FindClosestCommonManager("Bob", "Ann")
	require.NoError(t, err)
	require.Equal(t, 1, len(cm))
	checkTwoEmployes(t, "/Bureaucr.at/department 1/Bob", "/Bureaucr.at/Ann", cm[0])
	assert.Equal(t, "/Bureaucr.at/Claire", cm[0].Manager)

	cm, err = dir.FindClosestCommonManager("Joseph", "Kate")
	require.NoError(t, err)
	require.Equal(t, 1, len(cm))
	checkTwoEmployes(t, "/Bureaucr.at/department 1/Joseph", "/Bureaucr.at/department 1/Kate", cm[0])
	assert.Equal(t, "/Bureaucr.at/department 1/Bob", cm[0].Manager)

	cm, err = dir.FindClosestCommonManager("Joseph", "Donald")
	require.NoError(t, err)
	require.Equal(t, 1, len(cm))
	checkTwoEmployes(t, "/Bureaucr.at/department 1/Joseph", "/Bureaucr.at/department 2/Donald", cm[0])
	assert.Equal(t, "/Bureaucr.at/Claire", cm[0].Manager)

	cm, err = dir.FindClosestCommonManager("Joseph", "Ann")
	require.NoError(t, err)
	require.Equal(t, 1, len(cm))
	checkTwoEmployes(t, "/Bureaucr.at/department 1/Joseph", "/Bureaucr.at/Ann", cm[0])
	assert.Equal(t, "/Bureaucr.at/Claire", cm[0].Manager)

	cm, err = dir.FindClosestCommonManager("Monica", "Monica")
	require.NoError(t, err)
	require.Equal(t, 1, len(cm))
	checkTwoEmployes(t, "/Bureaucr.at/department 2/Monica", "/Bureaucr.at/department 1/Monica", cm[0])
	assert.Equal(t, "/Bureaucr.at/Claire", cm[0].Manager)

	cm, err = dir.FindClosestCommonManager("Jane", "Jane")
	require.NoError(t, err)
	require.Equal(t, 1, len(cm))
	checkTwoEmployes(t, "/Bureaucr.at/department 1/Jane", "/Bureaucr.at/department 1/Jane", cm[0])
	assert.Equal(t, "/Bureaucr.at/department 1/Bob", cm[0].Manager)

	cm, err = dir.FindClosestCommonManager("Bill", "Bill")
	require.NoError(t, err)
	require.Equal(t, 1, len(cm))
	checkTwoEmployes(t, "/Bureaucr.at/department 2/Bill", "/Bureaucr.at/department 1/Bill", cm[0])
	assert.Equal(t, "/Bureaucr.at/Claire", cm[0].Manager)

	cm, err = dir.FindClosestCommonManager("John", "John")
	require.NoError(t, err)
	require.Equal(t, 3, len(cm))

	expected := makeCommonManagerMap([]CommonManager{
		{
			Employee1: "/Bureaucr.at/John",
			Employee2: "/Bureaucr.at/John",
			Manager:   "/Bureaucr.at/Claire",
		},
		{
			Employee1: "/Bureaucr.at/John",
			Employee2: "/Bureaucr.at/department 1/John",
			Manager:   "/Bureaucr.at/Claire",
		},
		{
			Employee1: "/Bureaucr.at/John",
			Employee2: "/Bureaucr.at/department 1/John",
			Manager:   "/Bureaucr.at/Claire",
		},
	})
	received := makeCommonManagerMap(cm)
	assert.True(t, reflect.DeepEqual(expected, received), expected, received)
}

func TestThreeOrgUnits(t *testing.T) {
	dir := NewDirectory(
		&OrgUnit{
			Name:    "Bureaucr.at",
			Manager: Manager{Employee{Name: "Claire"}},
			OrgUnits: []*OrgUnit{
				{
					Name:    "department 1",
					Manager: Manager{Employee{Name: "Bob"}},
					OrgUnits: []*OrgUnit{
						{
							Name:    "group 1.1",
							Manager: Manager{Employee{Name: "Mark"}},
							Reports: []*Employee{
								{Name: "John"},
								{Name: "Elsa"},
							},
						},
						{
							Name:    "group 1.2",
							Manager: Manager{Employee{Name: "Paul"}},
							Reports: []*Employee{
								{Name: "Mike"},
								{Name: "Bill"},
							},
						},
					},
				},
				{
					Name:    "department 2",
					Manager: Manager{Employee{Name: "Bob"}},
					OrgUnits: []*OrgUnit{
						{
							Name:    "group 2.1",
							Manager: Manager{Employee{Name: "Henry"}},
							Reports: []*Employee{
								{Name: "Andrew"},
								{Name: "Fil"},
							},
						},
						{
							Name:    "group 2.2",
							Manager: Manager{Employee{Name: "Joan"}},
							Reports: []*Employee{
								{Name: "Helmut"},
								{Name: "Ann"},
							},
						},
					},
				},
			},
		},
	)
	assert.NotNil(t, dir)

	cm, err := dir.FindClosestCommonManager("John", "Elsa")
	require.NoError(t, err)
	require.Equal(t, 1, len(cm))

	expected := makeCommonManagerMap([]CommonManager{
		{
			Employee1: "/Bureaucr.at/department 1/group 1.1/John",
			Employee2: "/Bureaucr.at/department 1/group 1.1/Elsa",
			Manager:   "/Bureaucr.at/department 1/group 1.1/Mark",
		},
	})
	received := makeCommonManagerMap(cm)
	assert.True(t, reflect.DeepEqual(expected, received), expected, received)

	cm , err= dir.FindClosestCommonManager("John", "Ann")
	require.Equal(t, 1, len(cm))

	expected = makeCommonManagerMap([]CommonManager{
		{
			Employee1: "/Bureaucr.at/department 1/group 1.1/John",
			Employee2: "/Bureaucr.at/department 2/group 2.2/Ann",
			Manager:   "/Bureaucr.at/Claire",
		},
	})
	received = makeCommonManagerMap(cm)
	assert.True(t, reflect.DeepEqual(expected, received), expected, received)

	cm, err = dir.FindClosestCommonManager("John", "Paul")
	require.NoError(t, err)
	require.Equal(t, 1, len(cm))

	expected = makeCommonManagerMap([]CommonManager{
		{
			Employee1: "/Bureaucr.at/department 1/group 1.1/John",
			Employee2: "/Bureaucr.at/department 1/group 1.2/Paul",
			Manager:   "/Bureaucr.at/department 1/Bob",
		},
	})
	received = makeCommonManagerMap(cm)
	assert.True(t, reflect.DeepEqual(expected, received), expected, received)
}
