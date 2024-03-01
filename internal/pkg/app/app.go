package app

import (
	"context"
	"errors"
	"fmt"
	pluralize "github.com/gertd/go-pluralize"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"time"
)

type NotifyContext struct {
	gorm.Model
	ChatId                 int64
	OriginalMessageId      int
	RescheduleCallbackData string
	TimeableName           string
	TriggerAt              int64
}

type UserOffset struct {
	UserId int64
	Offset int
}

type App struct {
	bot         *bot.Bot
	db          *gorm.DB
	rescheduled chan struct{}
	ctx         context.Context
	pluralize   *pluralize.Client
}

func New(ctx context.Context, b *bot.Bot, dsn string) App {
	return App{
		ctx:         ctx,
		bot:         b,
		db:          GetDatabase(dsn),
		rescheduled: make(chan struct{}),
		pluralize:   pluralize.NewClient(),
	}
}

func (a *App) Start() {
	sleep := a.getSleepChannel()
	go func() {
		for {
			select {
			case _ = <-a.rescheduled:
				sleep = a.getSleepChannel()
				continue
			case _ = <-sleep:
				a.sendNotification()
				sleep = a.getSleepChannel()
				continue
			case _ = <-a.ctx.Done():
				return
			}
		}
	}()
}

func (a *App) sendNotification() {
	var ctx NotifyContext

	err := a.db.Order("trigger_at").First(&ctx).Error
	if err != nil {
		log.Println("could not get top notification context:", err)
		return
	}

	log.Println(fmt.Sprintf("sending notification for: %+v", ctx))

	defer a.db.Delete(&ctx)

	// Delete the original timer message
	_, err = a.bot.DeleteMessage(a.ctx, &bot.DeleteMessageParams{
		ChatID:    ctx.ChatId,
		MessageID: ctx.OriginalMessageId,
	})
	if err != nil {
		log.Println("unable to delete message:", err)
	}

	// Send telegram notification
	ikb := [][]models.InlineKeyboardButton{
		{{Text: "Restart", CallbackData: ctx.RescheduleCallbackData}},
	}

	_, err = a.bot.SendMessage(a.ctx, &bot.SendMessageParams{
		ChatID:      ctx.ChatId,
		Text:        fmt.Sprintf(`Your *%s* are ready\!`, a.pluralize.Plural(ctx.TimeableName)),
		ParseMode:   models.ParseModeMarkdown,
		ReplyMarkup: models.InlineKeyboardMarkup{InlineKeyboard: ikb},
	})

	if err != nil {
		log.Println("could not send notification:", err)
	}
}

func (a *App) cancelNotification(id uint) {
	err := a.db.Delete(&NotifyContext{}, id).Error
	if err != nil {
		log.Println("could delete notification from database:", err)
	} else {
		log.Println(fmt.Sprintf("deleted notification (id: %d) from database", id))
	}

	a.rescheduled <- struct{}{}
}

func (a *App) getSleepChannel() <-chan time.Time {
	var ctx NotifyContext

	err := a.db.Order("trigger_at").First(&ctx).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		log.Panic("unable to determine next sleep time:", err)
	}

	log.Println("next waking at:", time.Unix(ctx.TriggerAt, 0).Format("15:04"))

	return time.After(time.Unix(ctx.TriggerAt, 0).Sub(time.Now().UTC()))
}

func (a *App) ScheduleNotification(ctx NotifyContext) uint {
	a.db.Create(&ctx)
	a.rescheduled <- struct{}{}
	log.Println("scheduled notification for:", ctx)
	return ctx.ID
}

func (a *App) setUserOffset(userId int64, value int) {
	uo := UserOffset{
		UserId: userId,
		Offset: value,
	}

	err := a.db.Where(UserOffset{UserId: userId}).Assign(UserOffset{Offset: value}).FirstOrCreate(&uo).Error
	if err != nil {
		log.Println("unable to set user offset:", err)
	}
}

func (a *App) getUserOffset(userId int64) int {
	var uo UserOffset
	err := a.db.Where(UserOffset{UserId: userId}).First(&uo).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println(fmt.Sprintf("user with id: '%d' has no offset configured", userId))
			return 0
		}
	}
	return uo.Offset
}

func GetDatabase(dsn string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Panic("failed to connect database", err)
	}

	err = db.AutoMigrate(&NotifyContext{}, &UserOffset{})
	if err != nil {
		log.Panic("failed to auto migrate database", err)
	}

	return db
}
