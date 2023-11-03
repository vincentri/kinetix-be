package repository

import (
	"kinetix.com/db"
	"kinetix.com/lib"
)

type ICaseCommentsRepository interface {
	Db() *db.PrismaClient
}

type CaseCommentsRepository struct {
	db lib.Prisma
}

func NewCaseCommentsRepository(db lib.Prisma) ICaseCommentsRepository {
	return &CaseCommentsRepository{
		db: db,
	}
}

func (c CaseCommentsRepository) Db() *db.PrismaClient {
	return c.db.Client
}
