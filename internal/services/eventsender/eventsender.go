package eventsender

import (
	"context"
	"log/slog"
	"time"
	"url-shortener/internal/domain"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/storage/sqlite"
)

type Sender struct {
	storage *sqlite.Storage
	log     *slog.Logger
}

func New(storage *sqlite.Storage, log *slog.Logger) *Sender {
	return &Sender{
		storage: storage,
		log:     log,
	}
}

func (s *Sender) StartProcessingEvents(ctx context.Context, handlePeriod time.Duration) {
	const op = "services.eventsender.StartProcessingEvents"

	log := s.log.With(slog.String("op", op))

	ticker := time.NewTicker(handlePeriod)

	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Info("stopping event processing")
				return
			case <-ticker.C:
				// noop тикер блокирует горутину пока не будет получен тик
			}

			event, err := s.storage.GetNewEvent()
			if err != nil {
				log.Error("failed to get new event", sl.Err(err))
				continue
			}
			if event.ID == 0 {
				log.Debug("no new events")
				continue
			}

			s.SendMessage(event)

			if err := s.storage.SetDone(event.ID); err != nil {
				log.Error("failed to set event done", sl.Err(err))
				continue
			}

		}
	}()
}

func (s *Sender) SendMessage(event domain.Event) {
	const op = "services.eventsender.SendMessage"

	log := s.log.With(slog.String("op", op))
	log.Info("sending message", slog.Any("event", event))

	// TODO: implement sending message to the Kafka
}
