package telego

import (
	"fmt"
	"testing"

	"github.com/SakoDroid/telego/v2/configs"
	objs "github.com/SakoDroid/telego/v2/objects"
)

func TestBot(t *testing.T) {
	bot, err := NewBot(configs.Default("5078468473:AAHLCQfMnJTIFM25rFlbU2k422EZKmARK0s"))

	if err != nil {
		panic(err)
	}

	// The general update channel.
	updateChannel := *(bot.GetUpdateChannel())

	// Adding a handler. Everytime the bot receives message "hi" in a private chat, it will respond "hi to you too".
	bot.AddHandler("hi", func(u *objs.Update) {
		_, err := bot.SendMessage(u.Message.Chat.Id, "hi to you too", "", u.Message.MessageId, false, false)
		if err != nil {
			fmt.Println(err)
		}
	}, "private")

	// Adding a handler. Everytime the bot receives bloater message in a private chat, it will respond "default".
	bot.AddHandler("key", func(u *objs.Update) {

		//Create the custom keyboard.
		kb := bot.CreateKeyboard(false, false, false, false, "menu")

		//Add buttons to it. First argument is the button's text and the second one is the row number that the button will be added to it.
		kb.AddButton("button1", 1)
		kb.AddButton("button2", 1)
		kb.AddButton("button3", 2)

		//Sends the message along with the keyboard.
		_, err := bot.AdvancedMode().ASendMessage(u.Message.Chat.Id, "invalid keywords, please select to continue", "", u.Message.MessageId, 0, false, false, nil, false, false, kb)
		if err != nil {
			fmt.Println(err)
		}

	}, "private")

	// Monitores any other update. (Updates that don't contain text message "hi" in a private chat)
	go func() {
		for {
			update := <-updateChannel
			fmt.Println(update.Update_id)

			//Some processing on the update
		}
	}()

	bot.Run(true)
}
