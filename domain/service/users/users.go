package usersService

import (
	"context"
	"net/http"

	"kinetix.com/db"
	"kinetix.com/domain/repository"
	"kinetix.com/domain/service"
)

type IUsersService interface {
	FindAll(w http.ResponseWriter, r *http.Request)
	service.CoreInterface
}

type UsersService struct {
	service.CoreService
	UserRepository repository.IUsersRepository
}

func NewUsersService(UsersRepository repository.IUsersRepository) IUsersService {
	return &UsersService{
		UserRepository: UsersRepository,
		CoreService:    service.CoreService{},
	}
}

func (c UsersService) FindAll(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	data, errData := c.UserRepository.Db().Users.FindMany(
		db.Users.Status.Equals(db.ActiveStatusActive),
	).OrderBy(
		db.Users.Name.Order(db.SortOrderDesc),
	).Exec(ctx)
	if errData != nil {
		c.CoreService.ErrorResponse(w, errData, http.StatusInternalServerError, "")
		return
	}
	c.CoreService.SuccessResponse(w, map[string]interface{}{
		"data": data,
	})
}
