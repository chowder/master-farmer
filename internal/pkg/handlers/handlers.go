package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/chowder/master-farmer/internal/pkg/farming"
	"github.com/chowder/master-farmer/internal/pkg/utils"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"log"
	"strings"
)

const (
	ExpandCategoryPrefix  = "expand:"
	CreateTimerPrefix     = "create:"
	CancelTimerPrefix     = "cancel:"
	ReturnPrefix          = "return:"
	EndConversationPrefix = "end:"
)

type TimerContext struct {
	Category string
	Name     string
}

func (t TimerContext) GetTimeable() (farming.Timeable, error) {
	timeables, ok := farming.TimeablesByCategory[t.Category]
	if !ok {
		return nil, errors.New(fmt.Sprintf("could not find timables for category: %s", t.Category))
	}

	for _, timeable := range timeables {
		if timeable.GetName() == t.Name {
			return timeable, nil
		}
	}

	return nil, errors.New(fmt.Sprintf("could not find timables for category: %s with name: %s", t.Category, t.Name))
}

func ExpandCategoryHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	utils.AnswerCallbackQuery(ctx, b, update)

	callback := update.CallbackQuery
	category, ok := strings.CutPrefix(callback.Data, ExpandCategoryPrefix)
	if !ok {
		log.Println(fmt.Sprintf("callback callback data: '%s' did not contain '%s'", callback.Data, ExpandCategoryPrefix))
		return
	}

	timeables, ok := farming.TimeablesByCategory[category]
	if !ok {
		log.Println(fmt.Sprintf("unknown category: '%s'", category))
		return
	}

	var ikb [][]models.InlineKeyboardButton
	for _, t := range timeables {
		timerCtx := TimerContext{
			Category: category,
			Name:     t.GetName(),
		}

		timerCtxJson, _ := json.Marshal(timerCtx)

		ikb = append(ikb, []models.InlineKeyboardButton{
			{Text: t.GetName(), CallbackData: CreateTimerPrefix + string(timerCtxJson)},
		})
	}

	ikb = append(ikb, []models.InlineKeyboardButton{
		{Text: "« Back", CallbackData: ReturnPrefix},
	})

	_, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
		Text:        "Choose the crop to time!",
		ChatID:      callback.Message.Chat.ID,
		MessageID:   callback.Message.MessageID,
		ReplyMarkup: models.InlineKeyboardMarkup{InlineKeyboard: ikb},
	})

	if err != nil {
		log.Println(err)
	}
}

func ShowCategory(ctx context.Context, b *bot.Bot, update *models.Update) {
	var ikb [][]models.InlineKeyboardButton
	for _, k := range farming.Categories {
		ikb = append(ikb, []models.InlineKeyboardButton{
			{Text: k, CallbackData: ExpandCategoryPrefix + k},
		})
	}

	ikb = append(ikb, []models.InlineKeyboardButton{
		{Text: "« Cancel", CallbackData: EndConversationPrefix},
	})

	kb := &models.InlineKeyboardMarkup{
		InlineKeyboard: ikb,
	}

	callback := update.CallbackQuery

	var err error
	if callback != nil {
		_, err = b.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:      callback.Message.Chat.ID,
			MessageID:   callback.Message.MessageID,
			Text:        "Choose a type of crop to time!",
			ReplyMarkup: kb,
		})
	} else {
		_, err = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      update.Message.Chat.ID,
			Text:        "Choose a type of crop to time!",
			ReplyMarkup: kb,
		})
	}

	if err != nil {
		log.Println(err)
	}
}

func EndConversationHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	callback := update.CallbackQuery
	_, err := b.DeleteMessage(ctx, &bot.DeleteMessageParams{
		ChatID:    callback.Message.Chat.ID,
		MessageID: callback.Message.MessageID,
	})

	if err != nil {
		log.Println(err)
	}
}
