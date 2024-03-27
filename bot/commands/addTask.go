package commands

import (
	"log"
	"root/database"
	"root/database/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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
	methodDB := database.MethodDB()
	db, err := methodDB.Connect()
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
		// команды бота как ооп метод чтобы было прощи
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
	subjectName, task = "", ""
	step = 0
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Пожалуйста, введите предмет:")
	_, err := hb.bot.Send(msg)
	if err == nil {
		step = 1
	}

}

func (hb *HomeworkBot) processSubject(update tgbotapi.Update) {
	subjectName = update.Message.Text
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Пожалуйста, введите задание:")
	_, err := hb.bot.Send(msg)
	if err == nil {
		step = 2
	}
}

func (hb *HomeworkBot) processTask(update tgbotapi.Update) {
	task = update.Message.Text
	homework := &models.Homework{
		SubjectName: subjectName,
		Task:        task,
	}
	err := hb.db.Create(&homework).Error
	if err != nil {
		log.Fatal("ошибка: не удалось создать задание. \n", err)
	}
	confirmMsg := tgbotapi.NewMessage(update.Message.Chat.ID, "Задание добавлено!")
	if _, err := hb.bot.Send(confirmMsg); err != nil {
		log.Fatal("ошибка: не удалось отправить сообщение. \n", err)
	}

	// сброс значений
	subjectName = ""
	task = ""
	step = 0

}
