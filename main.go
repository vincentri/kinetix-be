package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"kinetix.com/domain/repository"
	areasService "kinetix.com/domain/service/areas"
	caseCommentsService "kinetix.com/domain/service/caseComments"
	casesService "kinetix.com/domain/service/cases"
	teamsService "kinetix.com/domain/service/teams"
	usersService "kinetix.com/domain/service/users"
	"kinetix.com/lib"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbClient := lib.NewPrisma()
	areaRepo := repository.NewAreasRepository(dbClient)
	caseRepo := repository.NewCasesRepository(dbClient)
	teamRepo := repository.NewTeamsRepository(dbClient)
	userRepo := repository.NewUsersRepository(dbClient)
	caseCommentsRepo := repository.NewCaseCommentsRepository(dbClient)

	areaServices := areasService.NewAreasService(areaRepo)
	caseServices := casesService.NewCasesService(caseRepo)
	teamServices := teamsService.NewTeamsService(teamRepo)
	userServices := usersService.NewUsersService(userRepo)
	caseCommentsServices := caseCommentsService.NewCaseCommentsService(caseCommentsRepo)

	r := chi.NewRouter()

	port := os.Getenv("APP_PORT")
	r.Route("/cases", func(r chi.Router) {
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", caseServices.Find)
			r.Route("/comments", func(r chi.Router) {
				r.Get("/", caseCommentsServices.FindAll)
				r.Post("/", caseCommentsServices.Create)
			})
		})
		r.Get("/", caseServices.FindAll)
		r.Put("/", caseServices.Update)
	})
	r.Get("/areas", areaServices.FindAll)
	r.Get("/teams", teamServices.FindAll)
	r.Get("/users", userServices.FindAll)
	fmt.Println("Running in " + port)

	defer func() {
		if err := dbClient.Client.Prisma.Disconnect(); err != nil {
			panic(err)
		}
	}()
	http.ListenAndServe(":"+port, r)
}
