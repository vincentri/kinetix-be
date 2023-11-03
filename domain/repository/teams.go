package repository

import (
	"kinetix.com/db"
	"kinetix.com/lib"
)

type ITeamsRepository interface {
	Db() *db.PrismaClient
}

type TeamsRepository struct {
	db lib.Prisma
}

func NewTeamsRepository(db lib.Prisma) ITeamsRepository {
	return &TeamsRepository{
		db: db,
	}
}

func (c TeamsRepository) Db() *db.PrismaClient {
	return c.db.Client
}
