/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"fmt"
	"goticka/pkg/db"
	"goticka/pkg/db/migrations"
	"goticka/pkg/domain/article"
	"goticka/pkg/domain/attachment"
	"goticka/pkg/domain/queue"
	"goticka/pkg/domain/ticket"
	"goticka/pkg/domain/user"
	"goticka/pkg/services"
	"goticka/pkg/version"
	"log"
)

func main() {
	//cmd.Execute()

	fmt.Printf("Running version %s\n", version.GetVersion())

	db := db.GetDB()

	migrations.Migrate(db)

	us := services.NewUserService()
	u1, u1err := us.Create(user.User{
		UserName: "USERNAME1",
		Password: "PASSWORD1",
	})
	if u1err != nil {
		panic(u1err)
	}

	u2, u2err := us.Create(user.User{
		UserName: "USERNAME2",
		Password: "PASSWORD2",
	})
	if u2err != nil {
		panic(u2err)
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
				Body: "BODY",
				From: u1,
				To:   u2,
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
	_, err1 := ticket_service.Create(t)
	log.Print(err1)

	retrivedUser, retrivedUserError := us.GetByID(1)
	log.Printf("%+v, %s", retrivedUser, retrivedUserError)
}
