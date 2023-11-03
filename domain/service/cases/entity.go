package casesService

type FindAllRequest struct {
	Page  int
	Limit int
}

type FindAllResponse struct {
	ID int
}

type UpdateRequest struct {
	ID          int    `validate:"required"`
	AuthorityId int    `validate:"required"`
	AssignId    int    `validate:"required"`
	TeamId      int    `validate:"required"`
	Status      string `validate:"required,oneof='Reopen' 'ReviewSubmitted' 'InProgress' 'Completed'"`
	Description string
}
