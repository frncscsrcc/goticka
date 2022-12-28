package services

import (
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
		t.Errorf("expected error, but not fot any error")
	}
	if wrongPswUser.ID != 0 {
		t.Errorf("wrong user id, expected %d, got %d", 0, wrongPswUser.ID)
	}
}
