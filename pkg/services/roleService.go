package services

import (
	"fmt"
	"goticka/pkg/adapters/cache"
	"goticka/pkg/adapters/repositories"
	"goticka/pkg/config"
	"goticka/pkg/dependencies"
	"goticka/pkg/domain/role"
	"goticka/pkg/events"
	"log"
	"strconv"
)

type RoleService struct {
	roleRepository repositories.RoleRepositoryInterface
}

func NewRoleService() RoleService {
	return RoleService{
		roleRepository: dependencies.DI().RoleRepository,
	}
}

func (rs RoleService) Create(q role.Role) (role.Role, error) {
	if validationError := q.Validate(); validationError != nil {
		return role.Role{}, validationError
	}

	createdQueue, err := rs.roleRepository.Create(q)
	if err != nil {
		return role.Role{}, err
	}
	log.Printf("created queue %d\n", createdQueue.ID)

	events.Handler().SendLocalEvent(events.LocalEvent{
		EventType: events.QUEUE_CREATED,
		QueueID:   createdQueue.ID,
	})

	return createdQueue, err
}

func (rs RoleService) GetByID(id int64) (role.Role, error) {
	// Check in the cache
	cached := dependencies.DI().Cache.Get(cache.Item{
		Type: "role",
		Key:  strconv.FormatInt(id, 10),
	})
	if cached.IsValid() {
		if value, ok := cached.Value.(role.Role); ok {
			return value, nil
		}
	}

	r, err := rs.roleRepository.GetByID(id)
	if err != nil {
		return role.Role{}, err
	}

	// Save in cache
	dependencies.DI().Cache.Set(cache.Item{
		Type:  "role",
		Key:   strconv.FormatInt(id, 10),
		Value: r,
		TTL:   config.GetConfig().Cache.RoleTTL,
	})

	return r, nil
}

func (rs RoleService) GetByName(roleName string) (role.Role, error) {
	// Check in the cache
	cacheKey := "role_name_" + roleName
	cached := dependencies.DI().Cache.Get(cache.Item{
		Type: "role",
		Key:  cacheKey,
	})
	if cached.IsValid() {
		if value, ok := cached.Value.(role.Role); ok {
			return value, nil
		}
	}

	r, err := rs.roleRepository.GetByName(roleName)
	if err != nil {
		return role.Role{}, err
	}

	// Save in cache
	dependencies.DI().Cache.Set(cache.Item{
		Type:  "role",
		Key:   cacheKey,
		Value: r,
		TTL:   config.GetConfig().Cache.RoleTTL,
	})

	return r, nil
}

func (rs RoleService) GetByUserID(userID int64) ([]role.Role, error) {
	// Check in the cache
	cacheKey := fmt.Sprintf("role_for_user_%s", strconv.FormatInt(userID, 10))
	cached := dependencies.DI().Cache.Get(cache.Item{
		Type: "role",
		Key:  cacheKey,
	})
	if cached.IsValid() {
		if value, ok := cached.Value.([]role.Role); ok {
			return value, nil
		}
	}

	roles, err := rs.roleRepository.GetByUserID(userID)
	if err != nil {
		return []role.Role{}, err
	}

	// Save in cache
	dependencies.DI().Cache.Set(cache.Item{
		Type:  "role",
		Key:   cacheKey,
		Value: roles,
		TTL:   config.GetConfig().Cache.RoleTTL,
	})

	return roles, nil
}
