package main

import (
	"log"
	"os"
	"root/database"
	"root/database/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type HomeworkBot struct {
	bot *tgbotapi.BotAPI
	db  *gorm.DB
}

var (
	subjectName string
	task        string
	step        int
)

func NewHomeworkBot(token string) (*HomeworkBot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	bot.Debug = true
	method := database.MethodDB()
	db, err := method.Connect()
	if err != nil {
		return nil, err
	}

	return &HomeworkBot{
		bot: bot,
		db:  db,
	}, nil
}

func (hb *HomeworkBot) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := hb.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil || update.Message.From.ID == hb.bot.Self.ID {
			continue
		}
		hb.processTaskUpdate(update)
	}
}

func (hb *HomeworkBot) processTaskUpdate(update tgbotapi.Update) {
	if update.Message.Text == "/task" {
		hb.resetVariablesAndRequestSubject(update)
	} else if step == 1 {
		hb.processSubject(update)
	} else if step == 2 {
		hb.processTask(update) // Pass the update parameter here
	}
}

func (hb *HomeworkBot) resetVariablesAndRequestSubject(update tgbotapi.Update) {
	// Reset variables for the new homework entry
	subjectName = ""
	task = ""
	step = 0
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Please enter the subject name:")
	_, err := hb.bot.Send(msg)
	if err == nil {
		step = 1
	}
}

func (hb *HomeworkBot) processSubject(update tgbotapi.Update) {
	subjectName = update.Message.Text
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Please enter the task:")
	_, err := hb.bot.Send(msg)
	if err == nil {
		step = 2
	}
}

func (hb *HomeworkBot) processTask(update tgbotapi.Update) {
	task = update.Message.Text
	homework := models.Homework{
		SubjectName: subjectName,
		Task:        task,
	}
	err := hb.db.Create(&homework).Error
	if err == nil {
		confirmMsg := tgbotapi.NewMessage(update.Message.Chat.ID, "Homework saved successfully!")
		hb.bot.Send(confirmMsg)
	}
	// Reset variables for the next homework entry
	subjectName = ""
	task = ""
	step = 0
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	token := os.Getenv("TOKEN")

	bot, err := NewHomeworkBot(token)
	if err != nil {
		log.Fatal("епта не работает")
	}

	bot.Start()
}
