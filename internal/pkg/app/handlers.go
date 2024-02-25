package app

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/chowder/master-farmer/internal/pkg/handlers"
	"github.com/chowder/master-farmer/internal/pkg/utils"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"log"
	"strconv"
	"strings"
	"time"
)

func (a *App) CreateTimerHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	utils.AnswerCallbackQuery(ctx, b, update)

	callback := update.CallbackQuery
	timerCtxJson, ok := strings.CutPrefix(callback.Data, handlers.CreateTimerPrefix)
	if !ok {
		log.Println(fmt.Sprintf("callback callback data: '%s' did not contain '%s'", callback.Data, handlers.ExpandCategoryPrefix))
		return
	}

	var timerCtx handlers.TimerContext
	err := json.Unmarshal([]byte(timerCtxJson), &timerCtx)
	if err != nil {
		log.Println("could not unmarshall timer context:", err)
		return
	}

	timeable, err := timerCtx.GetTimeable()
	if err != nil {
		log.Println("could not determine timeable:", err)
		return
	}

	offset := a.getUserOffset(callback.From.ID)

	triggerTime, err := timeable.GetTriggerTime(time.Now().UTC(), time.Duration(offset)*time.Minute)
	if err != nil {
		log.Println("could not determine next trigger time:", err)
	}

	notifyCtx := NotifyContext{
		ChatId:                 callback.Message.Chat.ID,
		OriginalMessageId:      callback.Message.MessageID,
		RescheduleCallbackData: callback.Data,
		TimeableName:           timeable.GetName(),
		TriggerAt:              triggerTime,
	}

	notificationId := a.ScheduleNotification(notifyCtx)

	// Feed back to user
	ikb := [][]models.InlineKeyboardButton{
		{
			{Text: "Cancel", CallbackData: fmt.Sprintf("%s%d", handlers.CancelTimerPrefix, notificationId)},
		},
	}
	_, err = b.EditMessageText(ctx, &bot.EditMessageTextParams{
		Text:        fmt.Sprintf("Timer: *%s*\n\nCompletes at: `%s UTC`", timeable.GetName(), time.Unix(triggerTime, 0).Format("Mon 15:04")),
		ParseMode:   models.ParseModeMarkdown,
		ChatID:      callback.Message.Chat.ID,
		MessageID:   callback.Message.MessageID,
		ReplyMarkup: models.InlineKeyboardMarkup{InlineKeyboard: ikb},
	})

	if err != nil {
		log.Println("could not update message for user:", err)
	}
}

func (a *App) CancelTimerHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	utils.AnswerCallbackQuery(ctx, b, update)

	callback := update.CallbackQuery
	data, ok := strings.CutPrefix(callback.Data, handlers.CancelTimerPrefix)
	if !ok {
		log.Println(fmt.Sprintf("callback callback data: '%s' did not contain '%s'", callback.Data, handlers.CancelTimerPrefix))
		return
	}

	nid, err := utils.ToUint(data)
	if err != nil {
		log.Println(fmt.Sprintf("could not parse notification id: '%s' as uint: %s", data, err))
	}

	a.cancelNotification(nid)

	_, err = b.DeleteMessage(ctx, &bot.DeleteMessageParams{
		ChatID:    callback.Message.Chat.ID,
		MessageID: callback.Message.MessageID,
	})

	if err != nil {
		log.Println("could not delete message:", err)
	}
}

func (a *App) SetOffsetHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	words := strings.Fields(update.Message.Text)
	if len(words) < 2 {
		sendOffsetHelp(ctx, b, update, "No offset value provided")
		return
	}

	data := words[1]
	value, err := strconv.Atoi(data)
	if err != nil {
		sendOffsetHelp(ctx, b, update, "Offset value was not an integer")
		return
	}

	if value <= 0 || value >= 30 {
		sendOffsetHelp(ctx, b, update, "Offset value must be between 1 and 29")
		return
	}

	a.setUserOffset(update.Message.From.ID, value)

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Successfully set your farming cycle offset",
	})

	if err != nil {
		log.Println("could not set message", err)
	}
}

func sendOffsetHelp(ctx context.Context, b *bot.Bot, update *models.Update, reason string) {
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   reason + "\n\nSend '/offset <1-29>' to set your farming cycle offset",
	})

	if err != nil {
		log.Println("could not send offset command help:", err)
	}
}
