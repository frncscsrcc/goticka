package services

import (
	"goticka/pkg/domain/article"
	"goticka/pkg/domain/queue"
	"goticka/pkg/domain/ticket"
	"goticka/pkg/domain/user"
	"goticka/testUtils"
	"testing"
)

func TestCreateTicket(t *testing.T) {
	testUtils.ResetTestDependencies()

	q1, _ := NewQueueService().Create(
		queue.Queue{
			Name: "Queue1",
		},
	)
	us := NewUserService()
	u1, _ := us.Create(
		user.User{
			UserName: "u1",
			Password: "p1",
		},
	)
	u2, _ := us.Create(
		user.User{
			UserName: "u2",
			Password: "p2",
		},
	)

	createdTicket, err := NewTicketService().Create(
		ticket.Ticket{
			Queue:   q1,
			Subject: "subject",
			Articles: []article.Article{
				{
					From: u1,
					To:   u2,
					Body: "BODY",
				},
			},
		},
	)

	if err != nil {
		t.Errorf("unexpected error %s", err)
	}
	if createdTicket.ID == 0 {
		t.Error("ticket ID should be initialized")
	}

}
