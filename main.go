/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"fmt"
	"goticka/pkg/db"
	"goticka/pkg/db/migrations"
	"goticka/pkg/dependencies"
	"goticka/pkg/domain/article"
	"goticka/pkg/domain/attachment"
	"goticka/pkg/domain/audit"
	"goticka/pkg/domain/queue"
	"goticka/pkg/domain/ticket"
	"goticka/pkg/domain/user"
	"goticka/pkg/events"
	"goticka/pkg/services"
	"goticka/pkg/version"
	"log"
)

func main() {
	//cmd.Execute()

	fmt.Printf("Running version %s\n", version.GetVersion())

	db := db.GetDB()

	migrations.Migrate(db)

	handler := events.Handler()
	handler.RegisterSyncCallBack(events.TICKET_CREATED, func(event events.LocalEvent) error {
		fmt.Println("NEW_TICKET1", event)
		return nil
	})
	handler.RegisterSyncCallBack(events.TICKET_CREATED, func(event events.LocalEvent) error {
		fmt.Println("NEW_TICKET2", event)
		return nil
	})

	us := services.NewUserService()
	agent, agentErr := us.CreateAgent(user.User{
		UserName: "USERNAME1",
		Password: "PASSWORD1",
		Email:    "Email1",
	})
	if agentErr != nil {
		panic(agentErr)
	}

	customer, customerErr := us.CreateCustomer(user.User{
		UserName: "USERNAME2",
		Password: "PASSWORD2",
		Email:    "Email2",
	})
	if customerErr != nil {
		panic(customerErr)
	}

	q1, newQueueError := services.NewQueueService().Create(queue.Queue{Name: "Queue1"})
	if newQueueError != nil {
		panic(newQueueError)
	}

	t := ticket.Ticket{
		Subject: "AAA",
		Queue:   q1,
		Articles: []article.Article{
			{
				Body:     "BODY",
				External: true,
				From:     customer,
				To:       agent,
				Attachments: []attachment.Attachment{
					{
						FileName:    "file1.txt",
						ContentType: "text/txt",
						URI:         "...",
						Size:        10,
						Raw:         []byte("content1"),
					},
					{
						FileName:    "file2.txt",
						ContentType: "text/txt",
						URI:         "...",
						Size:        20,
						Raw:         []byte("content2"),
					},
				},
			},
		},
	}

	ticket_service := services.NewTicketService()
	createdTicket, err1 := ticket_service.Create(t)
	log.Print(err1)

	retrivedUser, retrivedUserError := us.GetByID(1)
	log.Printf("%+v, %s", retrivedUser, retrivedUserError)

	ticket_service.GetByID(createdTicket.ID)
	ticket_service.GetByID(createdTicket.ID)
	ticket_service.GetByID(createdTicket.ID)

	dependencies.DI().AuditRepository.Save(audit.Audit{
		Message:  "Message",
		TicketID: 1,
	})

	JWT, authError := services.NewAuthService().PasswordAuthentication(
		agent.UserName,
		agent.Password,
	)
	if authError != nil {
		panic(authError)
	}
	fmt.Println(JWT)
	fmt.Println(services.NewAuthService().VerifyJWT(JWT))
}
