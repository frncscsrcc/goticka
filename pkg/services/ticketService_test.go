package services

import (
	"fmt"
	"goticka/pkg/domain/article"
	"goticka/pkg/domain/attachment"
	"goticka/pkg/domain/queue"
	"goticka/pkg/domain/ticket"
	"goticka/pkg/domain/user"
	"goticka/testUtils"
	"testing"
)

func TestCreateAndRetrieveTicket(t *testing.T) {
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

	subject := "SUBJECT"

	ts := NewTicketService()

	createdTicket, err := ts.Create(
		ticket.Ticket{
			Queue:   q1,
			Subject: subject,
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
	if createdTicket.Created.IsZero() {
		t.Error("creation date should not be empty")
	}
	if !createdTicket.Deleted.IsZero() {
		t.Error("deletion date should be empty")
	}
	if createdTicket.Subject != subject {
		t.Errorf("wrong subject, expected '%s' but got '%s'", subject, createdTicket.Subject)
	}

	retrivedTicket, retrivedTicketError := ts.GetByID(createdTicket.ID)
	if retrivedTicketError != nil {
		t.Errorf("unexpected error %v", retrivedTicketError)
	}
	if retrivedTicket.ID != createdTicket.ID {
		t.Errorf("wrong ticket ID, expected %d but got %d", createdTicket.ID, retrivedTicket.ID)
	}
	if retrivedTicket.Created.IsZero() {
		t.Error("creation date should not be empty")
	}
	if !retrivedTicket.Deleted.IsZero() {
		t.Error("deletion date should be empty")
	}
	if retrivedTicket.Subject != subject {
		t.Errorf("wrong subject, expected '%s' but got '%s'", subject, retrivedTicket.Subject)
	}

	// Check minimal queue attributes
	if retrivedTicket.Queue.ID != q1.ID {
		t.Errorf("wrong queue ID, expected '%d' but got '%d'", q1.ID, retrivedTicket.Queue.ID)
	}

	deleteError := ts.Delete(createdTicket)
	if deleteError != nil {
		t.Errorf("unexpected error %s", deleteError)
	}

	deletedTicket, deletedTicketError := ts.GetByID(createdTicket.ID)
	if deletedTicketError != nil {
		t.Errorf("unexpected error %s", deletedTicketError)
	}
	if deletedTicket.Deleted.IsZero() {
		t.Error("expected deletion date but the field is empty")
	}

}

func TestEnrichedTicket(t *testing.T) {
	testUtils.ResetTestDependencies()

	queue, _ := NewQueueService().Create(
		queue.Queue{
			Name: "Queue1",
		},
	)
	us := NewUserService()
	from, _ := us.Create(
		user.User{
			UserName: "FROM_USER",
			Password: "p1",
		},
	)
	to, _ := us.Create(
		user.User{
			UserName: "TO_USER",
			Password: "p2",
		},
	)

	attachments := make([]attachment.Attachment, 0)
	attachments = append(attachments, attachment.Attachment{
		URI:         "URI1",
		FileName:    "FILENAME1",
		ContentType: "CONTENTTYPE1",
		Size:        1,
		Raw:         []byte("RAW1"),
	})
	attachments = append(attachments, attachment.Attachment{
		URI:         "URI2",
		FileName:    "FILENAME2",
		ContentType: "CONTENTTYPE2",
		Size:        2,
		Raw:         []byte("RAW2"),
	})

	articles := make([]article.Article, 0)
	articles = append(articles, article.Article{
		From:        from,
		To:          to,
		Body:        "BODY1",
		Attachments: attachments,
	})
	articles = append(articles, article.Article{
		From: from,
		To:   to,
		Body: "BODY2",
	})

	subject := "SUBJECT"

	ts := NewTicketService()

	createdTicket, err := ts.Create(
		ticket.Ticket{
			Queue:    queue,
			Subject:  subject,
			Articles: articles,
		},
	)

	if err != nil {
		t.Errorf("unexpected error %s", err)
	}

	ticketID := createdTicket.ID
	if ticketID == 0 {
		t.Error("expected ticket ID bigger than zero")
	}

	enrichedTicket, enrichedTicketError := ts.EnrichTicketInfo(createdTicket)

	if enrichedTicketError != nil {
		t.Errorf("unexpected error %s", enrichedTicketError)
	}

	fmt.Printf("%+v\n", enrichedTicket.Articles)

	// Check the queue
	if enrichedTicket.Queue.Name != queue.Name {
		t.Errorf("wrong queue name, expected %s but got %s",
			queue.Name,
			enrichedTicket.Queue.Name,
		)
	}

	// Check the articles
	if len(enrichedTicket.Articles) != 2 {
		t.Errorf("wrong number of article returned, expected %d but got %d",
			len(articles),
			len(enrichedTicket.Articles))
	}
	for i, article := range articles {

		// Check the body
		if enrichedTicket.Articles[i].Body != article.Body {
			t.Errorf("wrong article body for article %d, expected %s but got %s",
				i,
				article.Body,
				enrichedTicket.Articles[i].Body,
			)
		}

		// Check the from and to fields
		if enrichedTicket.Articles[i].From.UserName != from.UserName {
			t.Errorf("wrong from username for article %d, expected %s but got %s",
				i,
				from.UserName,
				enrichedTicket.Articles[i].From.UserName,
			)
		}
		if enrichedTicket.Articles[i].To.UserName != to.UserName {
			t.Errorf("wrong from username for article %d, expected %s but got %s",
				i,
				to.UserName,
				enrichedTicket.Articles[i].To.UserName,
			)
		}
	}

	// Check the attachments (expected only in the first article)
	attachments1 := enrichedTicket.Articles[0].Attachments
	if len(attachments1) != len(attachments) {
		t.Errorf("wrong number of attachments for article %d, expected %d but got %d",
			enrichedTicket.Articles[0].ID,
			len(attachments),
			len(attachments1),
		)
	}
	for i, attachment := range attachments {
		if attachments1[i].FileName != attachment.FileName {
			t.Errorf("wrong filenmae for attachment %d (article %d), expected %s but got %s",
				i,
				enrichedTicket.Articles[0].ID,
				attachment.FileName,
				attachments1[i].FileName,
			)
		}
	}

	attachments2 := enrichedTicket.Articles[1].Attachments
	if len(attachments2) != 0 {
		t.Errorf("wrong number of attachments for article %d, expected %d but got %d",
			enrichedTicket.Articles[1].ID,
			0,
			len(attachments2),
		)
	}

}
