package commands

import (
	"log"
	"root/database/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (hb *HomeworkBot) getAllTask(update tgbotapi.Update) {
	if update.Message.Text == "/getAllTask" || update.Message.Text == "Получить все ДЗ" {
		var tasks []models.Homework
		err := hb.db.Find(&tasks).Error
		if err != nil {
			log.Fatal("ошибка: не удалось получить задания. \n", err)
		}
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ваши задания")
		for _, task := range tasks {
			//добавить форматированный импорт если будет не похуй
			msg.Text += task.SubjectName + " - " + task.Task + "\n"

		}
		if _, err := hb.bot.Send(msg); err != nil {
			log.Fatal("ошибка: не удалось отправить сообщение. \n", err)
		}
	}
}
