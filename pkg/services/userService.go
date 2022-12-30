package services

import (
	"errors"
	"goticka/pkg/adapters/repositories"
	"goticka/pkg/dependencies"
	"goticka/pkg/domain/role"
	"goticka/pkg/domain/user"
	"goticka/pkg/events"
	"log"
)

type UserService struct {
	userRepository repositories.UserRepositoryInterface
}

func NewUserService() UserService {
	return UserService{
		userRepository: dependencies.DI().UserRepository,
	}
}

func (us UserService) Create(u user.User) (user.User, error) {
	// if validationError := u.Validate(); validationError != nil {
	// 	return ticket.Ticket{}, validationError
	// }
	createdUser, err := us.userRepository.CreateUser(u)
	if err != nil {
		return user.User{}, err
	}
	log.Printf("created user %d\n", createdUser.ID)

	events.Handler().SendLocalEvent(events.LocalEvent{
		EventType: events.USER_CREATED,
		UserID:    createdUser.ID,
	})

	return createdUser, err
}

func (us UserService) GetByID(ID int64) (user.User, error) {
	user, err := us.userRepository.GetByID(ID)
	if err != nil {
		log.Printf("[ERROR] User ID=%d not found!\n", ID)
	}
	return user, err
}

func (us UserService) GetByUserName(userName string) (user.User, error) {
	user, err := us.userRepository.GetByUserName(userName)
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

func (us UserService) Delete(u user.User) error {
	if u.ID == 0 {
		return errors.New("can not delete an invalid user")
	}

	err := dependencies.DI().UserRepository.Delete(u)
	if err != nil {
		log.Printf("[ERROR] Can not delete User UserID=%d !\n", u.ID)
	}

	events.Handler().SendLocalEvent(events.LocalEvent{
		EventType: events.USER_DELETED,
		UserID:    u.ID,
	})
	return nil
}

func (us UserService) AddRole(u user.User, r role.Role) error {
	if u.ID == 0 {
		return errors.New("invalid user")
	}
	if r.ID == 0 {
		return errors.New("invalid role")
	}

	return us.userRepository.AddRole(u, r)
}

func (us UserService) RemoveRole(u user.User, r role.Role) error {
	if u.ID == 0 {
		return errors.New("invalid user")
	}
	if r.ID == 0 {
		return errors.New("invalid role")
	}

	return us.userRepository.RemoveRole(u, r)
}
