package CaseCommentsService

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"kinetix.com/db"
	"kinetix.com/domain/repository"
	"kinetix.com/domain/service"
)

type ICaseCommentsService interface {
	FindAll(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	service.CoreInterface
}

type CaseCommentsService struct {
	service.CoreService
	CaseCommentsRepository repository.ICaseCommentsRepository
}

func NewCaseCommentsService(CaseCommentsRepository repository.ICaseCommentsRepository) ICaseCommentsService {
	return &CaseCommentsService{
		CaseCommentsRepository: CaseCommentsRepository,
		CoreService:            service.CoreService{},
	}
}

func (c CaseCommentsService) FindAll(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	idInt, _ := strconv.Atoi(id)

	ctx := context.Background()
	data, errData := c.CaseCommentsRepository.Db().CaseComments.FindMany(
		db.CaseComments.CaseID.Equals(idInt),
	).With(db.CaseComments.User.Fetch()).OrderBy(db.CaseComments.ID.Order(db.SortOrderDesc)).Exec(ctx)

	if errData != nil {
		c.CoreService.ErrorResponse(w, errData, http.StatusNotFound, "Case not found")
		return
	}

	c.CoreService.SuccessResponse(w, map[string]interface{}{
		"data": data,
	})
}

func (c CaseCommentsService) Create(w http.ResponseWriter, r *http.Request) {
	req := CreateRequest{}
	errDecode := c.CoreService.DecodeBody(&req, r.Body)
	if errDecode != nil {
		c.CoreService.ErrorResponse(w, errDecode, http.StatusInternalServerError, "")
		return
	}
	validate := validator.New()
	err := validate.Struct(req)
	if err != nil {
		c.CoreService.ErrorResponse(w, err, http.StatusUnprocessableEntity, "")
		return
	}

	ctx := context.Background()
	data, errData := c.CaseCommentsRepository.Db().CaseComments.CreateOne(
		db.CaseComments.Message.Set(req.Message),
		db.CaseComments.User.Link(
			db.Users.ID.Equals(req.UserId),
		),
		db.CaseComments.Case.Link(
			db.Cases.ID.Equals(req.CaseId),
		),
	).Exec(ctx)
	if errData != nil {
		c.CoreService.ErrorResponse(w, errData, http.StatusNotFound, "Update fail. Please check again your data. "+errData.Error())
		return
	}

	c.CoreService.SuccessResponse(w, data)
}
