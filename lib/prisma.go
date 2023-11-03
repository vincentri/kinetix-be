package lib

import (
	"log"

	"kinetix.com/db"
)

type Prisma struct {
	Client *db.PrismaClient
}

func NewPrisma() Prisma {
	client := db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		log.Fatal(err)
	}

	return Prisma{
		Client: client,
	}
}

func (m Prisma) client() *db.PrismaClient {
	return m.Client
}
