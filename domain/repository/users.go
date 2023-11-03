package repository

import (
	"kinetix.com/db"
	"kinetix.com/lib"
)

type IUsersRepository interface {
	Db() *db.PrismaClient
}

type UsersRepository struct {
	db lib.Prisma
}

func NewUsersRepository(db lib.Prisma) IUsersRepository {
	return &UsersRepository{
		db: db,
	}
}

func (c UsersRepository) Db() *db.PrismaClient {
	return c.db.Client
}
