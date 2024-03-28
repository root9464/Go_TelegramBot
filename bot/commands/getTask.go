package commands

import (
	"log"
	"root/database/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (hb *HomeworkBot) getTask(update tgbotapi.Update) {
	state := &hb.state

	if update.Message.Text == "/getTask" || update.Message.Text == "Получить ДЗ" {
		state.IsSubjectInput = true
		state.IsTaskInput = false
		state.SubjectName = ""
		state.Task = ""
		state.IsAddTask = false

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите предмет:")
		if _, err := hb.bot.Send(msg); err != nil {
			log.Println("ошибка: не удалось отправить сообщение. \n", err)
		}
	} else if state.IsSubjectInput && !state.IsTaskInput && !state.IsAddTask {
		state.SubjectName = update.Message.Text

		var homework models.Homework
		err := hb.db.Where("subject_name = ?", state.SubjectName).First(&homework).Error
		if err != nil {
			log.Println("ошибка: не удалось получить задание. \n", err)
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Задание: "+homework.Task)
		if _, err := hb.bot.Send(msg); err != nil {
			log.Println("ошибка: не удалось отправить сообщение. \n", err)
		}

		state.IsSubjectInput = true
		state.IsTaskInput = false
	}
}
