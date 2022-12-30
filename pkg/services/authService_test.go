package services

import (
	"goticka/pkg/config"
	"goticka/pkg/domain/role"
	"goticka/pkg/domain/user"
	"goticka/testUtils"
	"testing"
	"time"
)

func TestJWTAuth(t *testing.T) {
	// Setting JWT expiration date to 1 sec
	originalConfig := config.GetConfig()
	shortTTLConfig := originalConfig
	shortTTLConfig.Secrets.JWTTTL = 1 * time.Second
	config.OverwriteConfig(shortTTLConfig)
	defer config.OverwriteConfig(originalConfig)

	testUtils.ResetTestDependencies()

	rs := NewRoleService()
	role1, _ := rs.Create(role.Role{Name: "ROLE1"})
	role2, _ := rs.Create(role.Role{Name: "ROLE2"})
	role3, _ := rs.Create(role.Role{Name: "ROLE3"})

	username := "username"
	password := "password"
	email := "email"

	us := NewUserService()
	createdUser, createdUserError := us.Create(user.User{
		UserName: username,
		Email:    email,
		Password: password,
	})
	if createdUserError != nil {
		t.Errorf("unexpected error %s", createdUserError)
	}
	us.AddRole(createdUser, role1)
	us.AddRole(createdUser, role2)
	us.AddRole(createdUser, role3)

	auth := NewAuthService()

	// CORRECT USERNAME AND PASSWORD
	validJwt, errAuth := auth.PasswordAuthentication(username, password)
	if errAuth != nil {
		t.Errorf("unexpected error %s", errAuth)
	}

	authData, errVerification := auth.VerifyJWT(validJwt)
	if errVerification != nil {
		t.Errorf("unexpected error %s", errVerification)
	}
	if authData.Username != username {
		t.Errorf("wrong username, expected %s got %s", username, authData.Username)
	}
	if authData.ID != createdUser.ID {
		t.Errorf("wrong id, expected %d got %d", createdUser.ID, authData.ID)
	}
	if len(authData.Roles) != 3 {
		t.Errorf("wrong role number, expected %d got %d", 3, len(authData.Roles))
	}

	// WRONG USERNAME / PASSWORD
	jwtWrong, wrongAuth := auth.PasswordAuthentication("XXX", "YYY")
	if wrongAuth == nil {
		t.Error("expected error but got nothing")
	}
	if jwtWrong != "" {
		t.Error("jwt should be an empty string!")
	}

	// INVALID JWT
	inValidJwt := validJwt + "INVALID!!"
	_, inValidJwtError := auth.VerifyJWT(inValidJwt)
	if inValidJwtError == nil {
		t.Errorf("expected error but got nothing")
	}

	// Expired JWT (TTW was set to 1 sec)
	time.Sleep(time.Second)
	authData2, errVerification2 := auth.VerifyJWT(validJwt)
	if errVerification2 == nil {
		t.Errorf("expected error (jwt expired) but got nothing")
	}
	if authData2.ID != 0 {
		t.Errorf("authData should be empty")
	}
}
