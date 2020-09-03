package internal

import (
	"github.com/stretchr/testify/require"
	"testing"
)
import "github.com/stretchr/testify/assert"

func TestOnlyCEO(t *testing.T) {
	dir := NewDirectory(
		&OrgUnit{
			Name:    "Bureaucr.at",
			Manager: Manager{Employee{Name: "Claire"}},
		})
	assert.NotNil(t, dir)

	cm := dir.FindClosestCommonManager("Claire", "Claire")
	assert.Equal(t, 0, len(cm))
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

	cm = dir.FindClosestCommonManager("Ann", "Julia")
	require.Equal(t, 1, len(cm))
	assert.Equal(t, "/Bureaucr.at/Ann", cm[0].Employee1)
	assert.Equal(t, "/Bureaucr.at/Julia", cm[0].Employee2)
	assert.Equal(t, "/Bureaucr.at/Claire", cm[0].Manager)

	cm = dir.FindClosestCommonManager("Bob", "Ann")
	require.Equal(t, 1, len(cm))
	assert.Equal(t, "/Bureaucr.at/department 1/Bob", cm[0].Employee1)
	assert.Equal(t, "/Bureaucr.at/Ann", cm[0].Employee2)
	assert.Equal(t, "/Bureaucr.at/Claire", cm[0].Manager)

	cm = dir.FindClosestCommonManager("Joseph", "Kate")
	require.Equal(t, 1, len(cm))
	assert.Equal(t, "/Bureaucr.at/department 1/Joseph", cm[0].Employee1)
	assert.Equal(t, "/Bureaucr.at/department 1/Kate", cm[0].Employee2)
	assert.Equal(t, "/Bureaucr.at/department 1/Bob", cm[0].Manager)

	cm = dir.FindClosestCommonManager("Joseph", "Donald")
	require.Equal(t, 1, len(cm))
	assert.Equal(t, "/Bureaucr.at/department 1/Joseph", cm[0].Employee1)
	assert.Equal(t, "/Bureaucr.at/department 2/Donald", cm[0].Employee2)
	assert.Equal(t, "/Bureaucr.at/Claire", cm[0].Manager)

	cm = dir.FindClosestCommonManager("Joseph", "Ann")
	require.Equal(t, 1, len(cm))
	assert.Equal(t, "/Bureaucr.at/department 1/Joseph", cm[0].Employee1)
	assert.Equal(t, "/Bureaucr.at/Ann", cm[0].Employee2)
	assert.Equal(t, "/Bureaucr.at/Claire", cm[0].Manager)

	cm = dir.FindClosestCommonManager("Joseph", "Ann")
	require.Equal(t, 1, len(cm))
	assert.Equal(t, "/Bureaucr.at/department 1/Joseph", cm[0].Employee1)
	assert.Equal(t, "/Bureaucr.at/Ann", cm[0].Employee2)
	assert.Equal(t, "/Bureaucr.at/Claire", cm[0].Manager)

	cm = dir.FindClosestCommonManager("Monica", "Monica")
	require.Equal(t, 1, len(cm))
	assert.Equal(t, "/Bureaucr.at/department 2/Monica", cm[0].Employee1)
	assert.Equal(t, "/Bureaucr.at/department 1/Monica", cm[0].Employee2)
	assert.Equal(t, "/Bureaucr.at/Claire", cm[0].Manager)

	cm = dir.FindClosestCommonManager("Jane", "Jane")
	require.Equal(t, 1, len(cm))
	assert.Equal(t, "/Bureaucr.at/department 1/Jane", cm[0].Employee1)
	assert.Equal(t, "/Bureaucr.at/department 1/Jane", cm[0].Employee2)
	assert.Equal(t, "/Bureaucr.at/department 1/Bob", cm[0].Manager)
}

/*

func TestThreeOrgUnits(t *testing.T) {
	dir := NewDirectory(
		&OrgUnit{
			Manager: Manager{Employee{Name: "Claire"}},
			OrgUnits: []*OrgUnit{
				{
					Manager: Manager{Employee{Name: "Bob"}},
					OrgUnits: []*OrgUnit{
						{
							Manager: Manager{Employee{Name: "Mark"}},
							Reports: []*Employee{
								{Name: "A1"},
								{Name: "A2"},
							},
						},
						{
							Manager: Manager{Employee{Name: "Paul"}},
							Reports: []*Employee{
								{Name: "B1"},
								{Name: "B2"},
							},
						},
					},
				},
				{
					Manager: Manager{Employee{Name: "Alice"}},
					Reports: []*Employee{
						{Name: "Fred"},
						{Name: "Donald"},
					},
					OrgUnits: []*OrgUnit{
						{
							Manager: Manager{Employee{Name: "Boris"}},
							Reports: []*Employee{
								{Name: "C1"},
								{Name: "C2"},
							},
						},
						{
							Manager: Manager{Employee{Name: "Pablo"}},
							Reports: []*Employee{
								{Name: "D1"},
								{Name: "D2"},
							},
						},
					},
				},
			},
			Reports: []*Employee{
				{Name: "Ann"},
				{Name: "Julia"},
			},
		},
	)
	assert.NotNil(t, dir)

	var err error

	var m *Manager

	m, err = dir.FindClosestCommonManager("D1", "D2")
	assert.NoError(t, err)
	assert.Equal(t, "Pablo", m.Name)

	m, err = dir.FindClosestCommonManager("D1", "C2")
	assert.NoError(t, err)
	assert.Equal(t, "Alice", m.Name)

	m, err = dir.FindClosestCommonManager("D1", "Fred")
	assert.NoError(t, err)
	assert.Equal(t, "Alice", m.Name)

	m, err = dir.FindClosestCommonManager("D1", "A1")
	assert.NoError(t, err)
	assert.Equal(t, "Claire", m.Name)
}
*/
