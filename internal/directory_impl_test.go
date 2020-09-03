package internal

import (
	"errors"
	"testing"
)
import "github.com/stretchr/testify/assert"

func TestOnlyCEO(t *testing.T) {
	dir := NewDirectory(&OrgUnit{Manager: Manager{Employee{Name: "Claire"}}})
	assert.NotNil(t, dir)
	_, err := dir.FindClosestCommonManager("Claire", "Claire")
	assert.Error(t, err)
	assert.True(t, errors.Is(err, errNoCommonManager))
}

func TestTwoOrgUnits(t *testing.T) {
	dir := NewDirectory(
		&OrgUnit{
			Manager: Manager{Employee{Name: "Claire"}},
			OrgUnits: []*OrgUnit{
				&OrgUnit{
					Manager: Manager{Employee{Name: "Bob"}},
					Reports: []*Employee{
						&Employee{Name: "Bill"},
						&Employee{Name: "John"},
					},
				},
				&OrgUnit{
					Manager: Manager{Employee{Name: "Alice"}},
					Reports: []*Employee{
						&Employee{Name: "Fred"},
						&Employee{Name: "Donald"},
					},
				},
			},
			Reports: []*Employee{
				&Employee{Name: "Ann"},
				&Employee{Name: "Julia"},
			},
		},
	)
	assert.NotNil(t, dir)

	var err error
	var m *Manager

	m, err = dir.FindClosestCommonManager("Ann", "Julia")
	assert.NoError(t, err)
	assert.Equal(t, "Claire", m.Name)

	m, err = dir.FindClosestCommonManager("Ann", "Claire")
	assert.Error(t, err)
	assert.True(t, errors.Is(err, errNoCommonManager))

	m, err = dir.FindClosestCommonManager("Fred", "Donald")
	assert.NoError(t, err)
	assert.Equal(t, "Alice", m.Name)

	m, err = dir.FindClosestCommonManager("Fred", "Bill")
	assert.NoError(t, err)
	assert.Equal(t, "Claire", m.Name)
}

func TestThreeOrgUnits(t *testing.T) {
	dir := NewDirectory(
		&OrgUnit{
			Manager: Manager{Employee{Name: "Claire"}},
			OrgUnits: []*OrgUnit{
				&OrgUnit{
					Manager: Manager{Employee{Name: "Bob"}},
					OrgUnits: []*OrgUnit{
						&OrgUnit{
							Manager: Manager{Employee{Name: "Mark"}},
							Reports: []*Employee{
								&Employee{Name: "A1"},
								&Employee{Name: "A2"},
							},
						},
						&OrgUnit{
							Manager: Manager{Employee{Name: "Paul"}},
							Reports: []*Employee{
								&Employee{Name: "B1"},
								&Employee{Name: "B2"},
							},
						},
					},
				},
				&OrgUnit{
					Manager: Manager{Employee{Name: "Alice"}},
					Reports: []*Employee{
						&Employee{Name: "Fred"},
						&Employee{Name: "Donald"},
					},
					OrgUnits: []*OrgUnit{
						&OrgUnit{
							Manager: Manager{Employee{Name: "Boris"}},
							Reports: []*Employee{
								&Employee{Name: "C1"},
								&Employee{Name: "C2"},
							},
						},
						&OrgUnit{
							Manager: Manager{Employee{Name: "Pablo"}},
							Reports: []*Employee{
								&Employee{Name: "D1"},
								&Employee{Name: "D2"},
							},
						},
					},
				},
			},
			Reports: []*Employee{
				&Employee{Name: "Ann"},
				&Employee{Name: "Julia"},
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
