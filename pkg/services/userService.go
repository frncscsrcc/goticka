package services

import (
	"goticka/pkg/adapters/repositories"
	"goticka/pkg/dependencies"
	"goticka/pkg/domain/user"
	"log"
)

type UserService struct {
	userRepository repositories.UserRepositoryInterface
}

func NewUserService() UserService {
	return UserService{}
}

func (us UserService) Create(u user.User) (user.User, error) {
	// if validationError := u.Validate(); validationError != nil {
	// 	return ticket.Ticket{}, validationError
	// }
	createdUser, err := dependencies.DI().UserRepository.CreateUser(u)
	if err != nil {
		return user.User{}, err
	}
	log.Printf("created user %d\n", createdUser.ID)
	return createdUser, err
}

func (us UserService) GetByID(ID int64) (user.User, error) {
	user, err := dependencies.DI().UserRepository.GetByID(ID)
	if err != nil {
		log.Printf("[ERROR] User ID=%d not found!\n", ID)
	}
	return user, err
}

func (us UserService) GetByUserName(userName string) (user.User, error) {
	user, err := dependencies.DI().UserRepository.GetByUserName(userName)
	if err != nil {
		log.Printf("[ERROR] User UserName=%s not found!\n", userName)
	}
	return user, err
}

func (us UserService) GetByUserNameAndPassword(userName string, password string) (user.User, error) {
	user, err := dependencies.DI().UserRepository.GetByUserNameAndPassword(userName, password)
	if err != nil {
		log.Printf("[ERROR] User UserName=%s Password=.... not found!\n", userName)
	}
	return user, err
}
