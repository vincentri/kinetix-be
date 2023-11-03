package casesService

import (
	"context"
	"math"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"kinetix.com/db"
	"kinetix.com/domain/repository"
	"kinetix.com/domain/service"
)

type ICasesService interface {
	FindAll(w http.ResponseWriter, r *http.Request)
	Find(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	service.CoreInterface
}

type CasesService struct {
	service.CoreService
	CaseRepository repository.ICasesRepository
}

func NewCasesService(casesRepository repository.ICasesRepository) ICasesService {
	return &CasesService{
		CaseRepository: casesRepository,
		CoreService:    service.CoreService{},
	}
}

func (c CasesService) FindAll(w http.ResponseWriter, r *http.Request) {
	var req FindAllRequest

	errDecode := c.CoreService.DecodeQueryString(&req, r.URL)
	if errDecode != nil {
		c.CoreService.ErrorResponse(w, errDecode, http.StatusInternalServerError, "")
		return
	}

	ctx := context.Background()
	skip := 0
	limit := 30
	if req.Limit > 0 && req.Limit <= 50 {
		limit = req.Limit
	}
	if req.Page > 1 {
		skip = (req.Page - 1) * req.Limit
	}

	data, errData := c.CaseRepository.Db().Cases.FindMany().OrderBy(
		db.Cases.ID.Order(db.SortOrderDesc),
	).Take(limit).Skip(skip).Exec(ctx)

	var countData []struct {
		Total db.RawString `json:"total"`
	}
	c.CaseRepository.Db().Prisma.QueryRaw("select count(*) as total from cases").Exec(ctx, &countData)

	totalToInt, _ := strconv.Atoi(string(countData[0].Total))
	if errData != nil {
		c.CoreService.ErrorResponse(w, errData, http.StatusNotFound, "Case not found")
		return
	}

	totalPage := math.Floor(float64(totalToInt / limit))
	hasNextPage := true
	if float64(req.Page) >= totalPage {
		hasNextPage = false
	}

	hasPrevPage := true
	if req.Page == 1 {
		hasPrevPage = false
	}

	c.CoreService.SuccessResponse(w, map[string]interface{}{
		"data": data,
		"pagination": map[string]interface{}{
			"totalData":        totalToInt,
			"totalPage":        totalPage,
			"limit":            limit,
			"page":             req.Page,
			"currentStartPage": req.Page * limit,
			"hasNextPage":      hasNextPage,
			"hasPrevPage":      hasPrevPage,
		},
	})
}

func (c CasesService) Find(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	idInt, _ := strconv.Atoi(id)

	ctx := context.Background()
	data, errData := c.CaseRepository.Db().Cases.FindFirst(
		db.Cases.ID.Equals(idInt),
	).With(
		db.Cases.Area.Fetch(),
		db.Cases.Authority.Fetch(),
		db.Cases.Assign.Fetch(),
		db.Cases.Team.Fetch(),
	).Exec(ctx)

	if errData != nil {
		c.CoreService.ErrorResponse(w, errData, http.StatusNotFound, "Case not found")
		return
	}

	var nextId int
	nextData, errNextData := c.CaseRepository.Db().Cases.FindFirst(
		db.Cases.ID.GT(idInt),
	).Exec(ctx)
	if errNextData == nil {
		nextId = nextData.ID
	}

	var prevId int
	prevData, errPrevData := c.CaseRepository.Db().Cases.FindFirst(
		db.Cases.ID.LT(idInt),
	).Exec(ctx)
	if errPrevData == nil {
		prevId = prevData.ID
	}

	c.CoreService.SuccessResponse(w, map[string]interface{}{
		"data":       data,
		"nextCursor": nextId,
		"prevCursor": prevId,
	})
}

func (c CasesService) Update(w http.ResponseWriter, r *http.Request) {
	req := UpdateRequest{}
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
	data, errData := c.CaseRepository.Db().Cases.FindUnique(
		db.Cases.ID.Equals(req.ID),
	).Update(
		db.Cases.Authority.Link(
			db.Users.ID.Equals(req.AuthorityId),
		),
		db.Cases.Assign.Link(
			db.Users.ID.Equals(req.AssignId),
		),
		db.Cases.Team.Link(
			db.Teams.ID.Equals(req.TeamId),
		),
		db.Cases.Status.Set("ReviewSubmitted"),
	).Exec(ctx)
	if errData != nil {
		c.CoreService.ErrorResponse(w, errData, http.StatusNotFound, "Update fail. Please check again your data")
		return
	}

	c.CoreService.SuccessResponse(w, data)
}
