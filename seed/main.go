package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/joho/godotenv"
	"kinetix.com/db"
	"kinetix.com/lib"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbClient := lib.NewPrisma()

	ctx := context.Background()

	var wg sync.WaitGroup
	wg.Add(40)

	go func() {
		for i := 0; i < 5; i++ {
			defer wg.Done()
			dbClient.Client.Teams.CreateOne(
				db.Teams.Name.Set(faker.DomainName()),
				db.Teams.Status.Set(db.ActiveStatusActive),
				db.Teams.UpdatedAt.Set(time.Now()),
			).Exec(ctx)
		}
	}()

	go func() {
		for i := 0; i < 5; i++ {
			defer wg.Done()
			dbClient.Client.Areas.CreateOne(
				db.Areas.Name.Set(faker.GetRealAddress().Address),
				db.Areas.Status.Set(db.ActiveStatusActive),
				db.Areas.UpdatedAt.Set(time.Now()),
			).Exec(ctx)
		}
	}()

	var authorities []db.UsersModel
	var staffs []db.UsersModel
	go func() {
		roles := []db.UserRole{db.UserRoleAuthority, db.UserRoleStaff}
		for i := 0; i < 30; i++ {
			defer wg.Done()
			role := roles[rand.Intn(len(roles))]
			result, _ := dbClient.Client.Users.CreateOne(
				db.Users.Name.Set(faker.FirstName()+" "+faker.LastName()),
				db.Users.Role.Set(role),
				db.Users.Status.Set(db.ActiveStatusActive),
				db.Users.UpdatedAt.Set(time.Now()),
			).Exec(ctx)
			if role == db.UserRoleAuthority {
				authorities = append(authorities, *result)
				continue
			}
			staffs = append(staffs, *result)
		}
	}()

	wg.Wait()
	time.Sleep(5)

	var wg1 sync.WaitGroup
	totalCases := 200
	wg1.Add(totalCases)

	go func() {
		risk := []db.Risk{db.RiskLow, db.RiskMedium, db.RiskHigh}
		ctx := context.Background()
		for i := 0; i < totalCases; i++ {
			defer wg1.Done()
			caseAt, _ := time.Parse("2006-01-02 15:04:05", faker.Timestamp())
			statuses := []db.Status{db.StatusReopen, db.StatusInProgress, db.StatusReviewSubmitted, db.StatusComplete}
			data, _ := dbClient.Client.Cases.CreateOne(
				db.Cases.Alert.Set(faker.Word()),
				db.Cases.CaseAt.Set(caseAt),
				db.Cases.Zone.Set(faker.GetRealAddress().City),
				db.Cases.Camera.Set(faker.GetRealAddress().State),
				db.Cases.Authority.Link(
					db.Users.ID.Equals(authorities[rand.Intn(len(authorities)-1)].ID),
				),
				db.Cases.Area.Link(
					db.Areas.ID.Equals(rand.Intn(4-1)+1),
				),
				db.Cases.Team.Link(
					db.Teams.ID.Equals(rand.Intn(4-1)+1),
				),
				db.Cases.Assign.Link(
					db.Users.ID.Equals(staffs[rand.Intn(len(staffs)-1)].ID),
				),
				db.Cases.Status.Set(statuses[rand.Intn(len(statuses)-1)]),
				db.Cases.Risk.Set(risk[rand.Intn(len(risk)-1)]),
				db.Cases.Status.Set(statuses[rand.Intn(len(statuses)-1)]),
				db.Cases.ReviewDescription.Set(faker.Sentence()),
				db.Cases.UpdatedAt.Set(time.Now()),
			).Exec(ctx)

			for j := 0; j < 10; j++ {
				dbClient.Client.CaseComments.CreateOne(
					db.CaseComments.Message.Set(faker.Sentence()),
					db.CaseComments.User.Link(
						db.Users.ID.Equals(data.AssignID),
					),
					db.CaseComments.Case.Link(
						db.Cases.ID.Equals(data.ID),
					),
					db.CaseComments.UpdatedAt.Set(time.Now()),
				).Exec(ctx)
			}

		}
	}()

	wg1.Wait()

	defer func() {
		if err := dbClient.Client.Prisma.Disconnect(); err != nil {
			panic(err)
		}
	}()
	fmt.Println("Seed finish")
}
