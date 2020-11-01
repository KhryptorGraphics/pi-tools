package messagesservice

import (
	tripspb "github.com/mjm/pi-tools/detect-presence/proto/trips"
	"github.com/mjm/pi-tools/homebase/bot/database"
	"github.com/mjm/pi-tools/homebase/bot/telegram"
	"github.com/mjm/pi-tools/storage"
)

type Server struct {
	db     storage.DB
	q      *database.Queries
	t      *telegram.Client
	trips  tripspb.TripsServiceClient
	chatID int
}

func New(db storage.DB, t *telegram.Client, trips tripspb.TripsServiceClient, chatID int) *Server {
	return &Server{
		db:     db,
		q:      database.New(db),
		t:      t,
		trips:  trips,
		chatID: chatID,
	}
}
