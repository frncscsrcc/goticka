package services

import (
	"goticka/pkg/domain/role"
	"goticka/testUtils"
	"testing"
)

func TestCreateAndRetriveRole(t *testing.T) {
	testUtils.ResetTestDependencies()

	rs := NewRoleService()

	roleName := "ROLE1"
	roleDescription := "DESCRIPTION"

	createdRole, err := rs.Create(
		role.Role{
			Name:        roleName,
			Description: roleDescription,
		},
	)
	if err != nil {
		t.Errorf("unexpected error %s", err)
	}
	if createdRole.ID == 0 {
		t.Error("queue ID should be initialized")
	}

	retrivedRole, retrivedRoleError := rs.GetByID(createdRole.ID)
	if retrivedRoleError != nil {
		t.Errorf("unexpected error %s", retrivedRoleError)
	}
	if retrivedRole.ID != createdRole.ID {
		t.Errorf("wrong id, expected %d but got %d", createdRole.ID, retrivedRole.ID)
	}

	if retrivedRole.Name != roleName {
		t.Errorf("wrong role name, expected %s but got %s", roleName, retrivedRole.Name)
	}
	if retrivedRole.Description != roleDescription {
		t.Errorf("wrong role description, expected %s but got %s", roleDescription, retrivedRole.Description)
	}
	if retrivedRole.Created.IsZero() {
		t.Error("created should not be empty")
	}
	if retrivedRole.Changed.IsZero() {
		t.Error("changed should not be empty")
	}
	if !retrivedRole.Deleted.IsZero() {
		t.Error("deleted should be empty")
	}

	retrivedByName, retrivedByNameError := rs.GetByName(roleName)
	if retrivedByNameError != nil {
		t.Errorf("unexpected error %s", retrivedByNameError)
	}
	if retrivedByName.ID != createdRole.ID {
		t.Errorf("wrong id, expected %d but got %d", createdRole.ID, retrivedByName.ID)
	}
}
