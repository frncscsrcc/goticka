package services

import (
	"goticka/pkg/domain/role"
	"goticka/pkg/domain/user"
	"goticka/testUtils"
	"testing"
)

func TestCreateUser(t *testing.T) {
	testUtils.ResetTestDependencies()
	us := NewUserService()

	// Create an user
	createdUser, err := us.Create(
		user.User{
			UserName: "TesUserCreate",
			Password: "Password",
		},
	)
	if err != nil {
		t.Errorf("unexpected error %s", err)
	}
	if createdUser.ID == 0 {
		t.Error("user ID should be initialized")
	}
}

func TestGetUserByID(t *testing.T) {
	testUtils.ResetTestDependencies()
	us := NewUserService()

	createdUser, _ := us.Create(
		user.User{
			UserName: "TestUserGetByID",
			Password: "Password",
		},
	)

	retrivedUserByID, retrivedUserByIDError := us.GetByID(createdUser.ID)
	if retrivedUserByIDError != nil {
		t.Errorf("unexpected error %s", retrivedUserByIDError)
	}
	if retrivedUserByID.ID != 1 {
		t.Errorf("wrong user id, expected %d, got %d", createdUser.ID, retrivedUserByID.ID)
	}
	if retrivedUserByID.UserName != "TestUserGetByID" {
		t.Errorf("wrong username, expected %s, got %s", createdUser.UserName, retrivedUserByID.UserName)
	}
	if retrivedUserByID.Password != "" {
		t.Errorf("the service should never return the user password")
	}

	notExistUser, notExistUserError := us.GetByID(999)
	if notExistUserError == nil {
		t.Errorf("expected error, but not fot any error")
	}
	if notExistUser.ID != 0 {
		t.Errorf("wrong user id, expected %d, got %d", 0, notExistUser.ID)
	}
}

func TestGetUserByUserName(t *testing.T) {
	testUtils.ResetTestDependencies()
	us := NewUserService()

	createdUser, _ := us.Create(
		user.User{
			UserName: "TestUserGetByUserName",
			Password: "Password",
		},
	)

	retrivedUser, retrivedUserError := us.GetByUserName(createdUser.UserName)
	if retrivedUserError != nil {
		t.Errorf("unexpected error %s", retrivedUserError)
	}
	if retrivedUser.UserName != createdUser.UserName {
		t.Errorf("wrong username , expected %s, got %s", createdUser.UserName, retrivedUser.UserName)
	}

	notExistUser, notExistUserError := us.GetByUserName("NotExistingUsername")
	if notExistUserError == nil {
		t.Errorf("expected error, but not fot any error")
	}
	if notExistUser.ID != 0 {
		t.Errorf("wrong user id, expected %d, got %d", 0, notExistUser.ID)
	}
}

func TestGetUserByUserNameAndPassword(t *testing.T) {
	testUtils.ResetTestDependencies()
	us := NewUserService()

	createdUser, _ := us.Create(
		user.User{
			UserName: "TestUserGetByUserNameAndPassword",
			Password: "TestPassword",
		},
	)

	retrivedUser, retrivedUserError := us.GetByUserNameAndPassword(
		createdUser.UserName,
		createdUser.Password,
	)

	if retrivedUserError != nil {
		t.Errorf("unexpected error %s", retrivedUserError)
	}
	if retrivedUser.ID != createdUser.ID {
		t.Errorf("wrong user id, expected %d, got %d", createdUser.ID, retrivedUser.ID)
	}

	wrongPswUser, wrongPswUserError := us.GetByUserNameAndPassword(
		createdUser.UserName,
		"WRONG_PASSWORD",
	)
	if wrongPswUserError == nil {
		t.Errorf("expected error, but dod not get any error")
	}
	if wrongPswUser.ID != 0 {
		t.Errorf("wrong user id, expected %d, got %d", 0, wrongPswUser.ID)
	}

	// Marking the user as deleted
	deleteError := us.Delete(createdUser)
	if deleteError != nil {
		t.Errorf("unexpected error deleting an user, got %s", deleteError)
	}

	// Check deleted users are not found
	deletedUser, deletedUserError := us.GetByUserNameAndPassword(
		createdUser.UserName,
		createdUser.Password,
	)
	if deletedUserError == nil {
		t.Errorf("expected error, but dod not get any error")
	}
	if deletedUser.ID != 0 {
		t.Errorf("wrong user id, expected %d, got %d", 0, deletedUser.ID)
	}
}

func TestDeleteUser(t *testing.T) {
	testUtils.ResetTestDependencies()
	us := NewUserService()

	createdUser, _ := us.Create(
		user.User{
			UserName: "TestUserGetByUserNameAndPassword",
			Password: "TestPassword",
		},
	)

	retrivedUser, _ := us.GetByUserName(createdUser.UserName)

	if !retrivedUser.Deleted.IsZero() {
		t.Errorf("deleted field should be empty, but got %s", retrivedUser.Deleted)
	}

	// Marking the user as deleted
	deleteError := us.Delete(createdUser)
	if deleteError != nil {
		t.Errorf("unexpected error deleting an user, got %s", deleteError)
	}

	retrivedUser, _ = us.GetByUserName(createdUser.UserName)

	if retrivedUser.Deleted.IsZero() {
		t.Error("deleted field should now not be empty")
	}

}

func TestUserRoles(t *testing.T) {
	testUtils.ResetTestDependencies()
	us := NewUserService()
	rs := NewRoleService()

	createdUser, _ := us.Create(
		user.User{
			UserName: "TestUserGetByUserNameAndPassword",
			Password: "TestPassword",
		},
	)

	role1 := role.Role{
		Name: "ROLE1",
	}
	role2 := role.Role{
		Name: "ROLE2",
	}
	role1, _ = rs.Create(role1)
	role2, _ = rs.Create(role2)

	var retrivedUser user.User
	var err error

	retrivedUser, err = us.GetByID(createdUser.ID)
	if err != nil {
		t.Errorf("unexpected error %s", err)
	}
	if len(retrivedUser.Roles) != 0 {
		t.Errorf("wrong number of roles, expected %d but got %d", 0, len(retrivedUser.Roles))
	}

	err = us.AddRole(createdUser, role1)
	if err != nil {
		t.Errorf("unexpected error %s", err)
	}
	retrivedUser, err = us.GetByID(createdUser.ID)
	if err != nil {
		t.Errorf("unexpected error %s", err)
	}
	if len(retrivedUser.Roles) != 1 {
		t.Errorf("wrong number of roles, expected %d but got %d", 1, len(retrivedUser.Roles))
	} else if retrivedUser.Roles[0].Name != role1.Name {
		t.Errorf("wrong role, expected %s but got %s", role1.Name, retrivedUser.Roles[0].Name)
	}

	err = us.AddRole(createdUser, role2)
	if err != nil {
		t.Errorf("unexpected error %s", err)
	}
	retrivedUser, err = us.GetByID(createdUser.ID)
	if err != nil {
		t.Errorf("unexpected error %s", err)
	}
	if len(retrivedUser.Roles) != 2 {
		t.Errorf("wrong number of roles, expected %d but got %d", 2, len(retrivedUser.Roles))
	}

	err = us.RemoveRole(createdUser, role1)
	if err != nil {
		t.Errorf("unexpected error %s", err)
	}
	retrivedUser, err = us.GetByID(createdUser.ID)
	if err != nil {
		t.Errorf("unexpected error %s", err)
	}
	if len(retrivedUser.Roles) != 1 {
		t.Errorf("wrong number of roles, expected %d but got %d", 1, len(retrivedUser.Roles))
	} else if retrivedUser.Roles[0].Name != role2.Name {
		t.Errorf("wrong role, expected %s but got %s", role2.Name, retrivedUser.Roles[0].Name)
	}

}
