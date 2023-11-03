package repository

import (
	"kinetix.com/db"
	"kinetix.com/lib"
)

type ICasesRepository interface {
	Db() *db.PrismaClient
}

type CasesRepository struct {
	db lib.Prisma
}

func NewCasesRepository(db lib.Prisma) ICasesRepository {
	return &CasesRepository{
		db: db,
	}
}

func (c CasesRepository) Db() *db.PrismaClient {
	return c.db.Client
}
