package teamsService

import (
	"context"
	"net/http"

	"kinetix.com/db"
	"kinetix.com/domain/repository"
	"kinetix.com/domain/service"
)

type ITeamsService interface {
	FindAll(w http.ResponseWriter, r *http.Request)
	service.CoreInterface
}

type TeamsService struct {
	service.CoreService
	TeamRepository repository.ITeamsRepository
}

func NewTeamsService(TeamsRepository repository.ITeamsRepository) ITeamsService {
	return &TeamsService{
		TeamRepository: TeamsRepository,
		CoreService:    service.CoreService{},
	}
}

func (c TeamsService) FindAll(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	data, errData := c.TeamRepository.Db().Teams.FindMany(
		db.Teams.Status.Equals(db.ActiveStatusActive),
	).OrderBy(
		db.Teams.Name.Order(db.SortOrderDesc),
	).Exec(ctx)
	if errData != nil {
		c.CoreService.ErrorResponse(w, errData, http.StatusInternalServerError, "")
		return
	}
	c.CoreService.SuccessResponse(w, map[string]interface{}{
		"data": data,
	})
}
