package commands

import (
	"root/database"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
)

type HomeworkBot struct {
	bot *tgbotapi.BotAPI
	db  *gorm.DB
}

var (
	isSubjectInput bool
	isTaskInput    bool
	subjectName    string
	task           string
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
		hb.processesBotUpdate(update)
	}
}

func (hb *HomeworkBot) processesBotUpdate(update tgbotapi.Update) {
	hb.addTask(update)
	hb.getAllTask(update)
}
