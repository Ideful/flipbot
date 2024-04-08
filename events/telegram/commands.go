package telegram

import (
	"log"
	"math/rand"
	"strings"
)

const (
	FlipCmd  = "/flip"
	HelpCmd  = "/help"
	StartCmd = "/start"
)

func (p *Processor) doCmd(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("got command:%v by %v", text, username)

	switch text {
	case FlipCmd:
		p.tg.SendMessage(chatID, flip())
	case HelpCmd:
		p.tg.SendMessage(chatID, msgHelp)
	case StartCmd:
		p.tg.SendMessage(chatID, msgHello)
	}

	// add page

	// rnd page /rnd

	// help /help

	// start /start
	return nil
}

func flip() string {
	n := rand.Intn(2)
	if n == 1 {
		return "tails"
	}
	return "heads"
}
