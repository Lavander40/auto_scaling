package telegram

import (
	"auto_scaling/lib/e"
	"auto_scaling/storage"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	StartCmd = "/start"
	HelpCmd  = "/help"
	AddCmd   = "/add"
	RmCmd    = "/rm"
	LimitCmd = "/limit"
	GetLast  = "/last"
)

func (p *Processor) doCmd(text string, chatId int, userName string) error {
	// start + info
	// app call /edit {count}
	// get calls /get_last
	text = strings.TrimSpace(text)

	log.Printf("run commant %s, by %s", text, userName)

	if isCall(text) {
		return p.addCall(1, text, chatId, userName)
	}

	if isLimit(text) {
		return p.addCall(2, text, chatId, userName)
	}

	switch text {
	case StartCmd:
		return p.sendWelcome(chatId)
	case HelpCmd:
		return p.sendHelp(chatId)
	case AddCmd, LimitCmd:
		return p.tg.SendMessage(chatId, no_amount)
	case GetLast:
		return p.getLast(chatId, userName)
	default:
		return p.tg.SendMessage(chatId, unknown_msg)
	}
}

func isCall(text string) bool {
	match, _ := regexp.MatchString("^/add -?[0-9]+$", text)
	return match
}

func isLimit(text string) bool {
	match, _ := regexp.MatchString("^/limit ?[0-9]+$", text)
	return match
}

func (p *Processor) addCall(t int, text string, chatId int, userName string) error {
	amount, _ := strconv.Atoi(text[strings.LastIndex(text, " ")+1:])
	if t == 1 && (amount < -20 || amount > 20) {
		return p.tg.SendMessage(chatId, oversize_msg)
	}
	if t == 2 && (amount < 1 || amount > 100) {
		return p.tg.SendMessage(chatId, impossible_percent_msg)
	}

	call := &storage.Call{
		Type:      t,
		Amount:    amount,
		UserName:  userName,
		CreatedAt: time.Now(),
	}

	if err := p.scaler.ApplyCall(call); err != nil {
		_ = p.tg.SendMessage(chatId, unable_to_scale)
		return e.WrapErr("can't apply call", err)
	}

	if err := p.storage.Save(call); err != nil {
		_ = p.tg.SendMessage(chatId, fail_msg)
		return e.WrapErr("can't run cmd /add", err)
	}

	return p.tg.SendMessage(chatId, sucess_msg)
}

func (p *Processor) getLast(chatId int, userName string) error {
	calls, err := p.storage.PickLastCalls(userName)
	if errors.Is(err, storage.ErrEmpty) || errors.Is(err, storage.ErrNoDir) {
		return p.tg.SendMessage(chatId, nofound_msg)
	}
	if err != nil {
		return e.WrapErr("can't do cmd /last", err)
	}

	var res string
	for _, call := range calls {
		switch call.Type {
		case 1:
			res += fmt.Sprintf("Число добавленных Гб: %d\nПользователь: %s\nВремя: %v\n\n", call.Amount, call.UserName, call.CreatedAt)
		case 2:
			res += fmt.Sprintf("Установлен лимит загрузки ОЗУ (%%): %d\nПользователь: %s\nВремя: %v\n\n", call.Amount, call.UserName, call.CreatedAt)
		default:
			return storage.ErrUnknownType
		}
	}

	return p.tg.SendMessage(chatId, res)
}

func (p *Processor) sendHelp(chatId int) error {
	return p.tg.SendMessage(chatId, help_msg)
}

func (p *Processor) sendWelcome(chatId int) error {
	return p.tg.SendMessage(chatId, welcome_msg)
}
