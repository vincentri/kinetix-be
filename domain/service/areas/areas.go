package areasService

import (
	"context"
	"net/http"

	"kinetix.com/db"
	"kinetix.com/domain/repository"
	"kinetix.com/domain/service"
)

type IAreasService interface {
	FindAll(w http.ResponseWriter, r *http.Request)
	service.CoreInterface
}

type AreasService struct {
	service.CoreService
	AreaRepository repository.IAreasRepository
}

func NewAreasService(AreasRepository repository.IAreasRepository) IAreasService {
	return &AreasService{
		AreaRepository: AreasRepository,
		CoreService:    service.CoreService{},
	}
}

func (c AreasService) FindAll(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	data, errData := c.AreaRepository.Db().Areas.FindMany(
		db.Areas.Status.Equals(db.ActiveStatusActive),
	).OrderBy(
		db.Areas.Name.Order(db.SortOrderDesc),
	).Exec(ctx)
	if errData != nil {
		c.CoreService.ErrorResponse(w, errData, http.StatusInternalServerError, "")
		return
	}
	c.CoreService.SuccessResponse(w, map[string]interface{}{
		"data": data,
	})
}
