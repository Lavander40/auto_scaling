package telegram

import (
	"auto_scaling/clients/telegram"
	"auto_scaling/events"
	"auto_scaling/lib/e"
	"auto_scaling/scaler"
	"auto_scaling/storage"
	"errors"
)

type Processor struct {
	tg      *telegram.Client
	offset  int
	storage storage.Storage
	scaler scaler.Scaler
}

type Meta struct {
	ChatId   int
	Username string
}

func New(client *telegram.Client, storage storage.Storage, scaler scaler.Scaler) *Processor {
	return &Processor{
		tg:      client,
		offset:  0,
		storage: storage,
		scaler: scaler,
	}
}

func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	updates, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, e.WrapErr("can't fetch", err)
	}

	if len(updates) == 0 {
		return nil, nil
	}

	res := make([]events.Event, 0, len(updates))

	for _, u := range updates {
		res = append(res, event(u))
	}

	p.offset = updates[len(updates)-1].Id + 1

	return res, nil
}

func (p *Processor) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.processMessage(event)
	default:
		return errors.New("unknown event type")
	}
}

func (p *Processor) processMessage(event events.Event) error {
	meta, err := fetchMeta(event)
	if err != nil {
		return e.WrapErr("can't process message", err)
	}
	if err := p.doCmd(event.Text, meta.ChatId, meta.Username); err != nil {
		return e.WrapErr("can't process message", err)
	}

	return nil
}

func fetchMeta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, errors.New("can't fetch meta")
	}

	return res, nil
}

func event(update telegram.Update) events.Event {
	uType := fetchType(update)

	res := events.Event{
		Type: uType,
		Text: fetchText(update),
	}

	if uType == events.Message {
		res.Meta = Meta{
			ChatId:   update.Message.Chat.Id,
			Username: update.Message.From.Username,
		}
	}

	return res
}

func fetchType(update telegram.Update) events.Type {
	if update.Message == nil {
		return events.Unknown
	}
	return events.Message
}

func fetchText(update telegram.Update) string {
	if update.Message == nil {
		return ""
	}
	return update.Message.Text
}
