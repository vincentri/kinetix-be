package CaseCommentsService

type CreateRequest struct {
	CaseId  int    `validate:"required"`
	UserId  int    `validate:"required"`
	Message string `validate:"required"`
}
