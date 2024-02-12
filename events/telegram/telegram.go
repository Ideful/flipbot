package telegram2

import (
	"fmt"

	"github.com/Ideful/flipbot/clients/telegram"
	"github.com/Ideful/flipbot/events"
	"github.com/Ideful/flipbot/storage"
)

type Processor struct {
	tg      *telegram.Client
	offset  int
	storage storage.Storage
}

type Meta struct {
	ChatID   int
	UserName string
}

func New(client *telegram.Client, storage storage.Storage) *Processor {
	return &Processor{
		tg:      client,
		storage: storage,
	}
}

func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	updates, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, fmt.Errorf("error whiel gettin events:%v", err)
	}
	if len(updates) == 0 {
		return nil, fmt.Errorf("no updatess found")
	}

	res := make([]events.Event, 0, len(updates))

	for _, u := range updates {
		res = append(res, event(u))
	}

	p.offset = updates[len(updates)-1].ID + 1

	return res, nil
}

func (p *Processor) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		p.ProcessMessage(event)
	default:
		return fmt.Errorf("can't process message")
	}
	return nil
}

func (p *Processor) ProcessMessage(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return fmt.Errorf("can't process message")
	}

	if err := p.doCmd(event.Text, meta.ChatID, meta.UserName); err != nil {
		return fmt.Errorf("can't process message, %v", err)
	}
	return nil
}

func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, fmt.Errorf("uknown meta type")
	}
	return res, nil
}

func event(u telegram.Update) events.Event {
	updType := fetchType(u)
	res := events.Event{
		Type: int(updType),
		Text: fetchText(u),
	}
	if updType == events.Message {
		res.Meta = Meta{
			ChatID:   u.Message.Chat.ID,
			UserName: u.Message.From.Username,
		}
	}
	// chatID username
	return res
}

func fetchText(upd telegram.Update) string {
	if upd.Message == nil {
		return ""
	}
	return upd.Message.Text
}

func fetchType(upd telegram.Update) events.Type {
	if upd.Message == nil {
		return events.Unknown
	}
	return events.Message
}
