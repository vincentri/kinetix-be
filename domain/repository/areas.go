package repository

import (
	"kinetix.com/db"
	"kinetix.com/lib"
)

type IAreasRepository interface {
	Db() *db.PrismaClient
}

type AreasRepository struct {
	db lib.Prisma
}

func NewAreasRepository(db lib.Prisma) IAreasRepository {
	return &AreasRepository{
		db: db,
	}
}

func (c AreasRepository) Db() *db.PrismaClient {
	return c.db.Client
}
