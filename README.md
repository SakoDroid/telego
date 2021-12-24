# Telebot

[![Go Reference](https://pkg.go.dev/badge/github.com/SakoDroid/telebot.svg)](https://pkg.go.dev/github.com/SakoDroid/telebot)
[![telegram bot api](https://img.shields.io/badge/telegram-telegram%20bot%20api-blue)](https://core.telegram.org/bots/api)

A Go library for creating telegram bots.

![telebot logo inspired by Golang logo](https://github.com/SakoDroid/telebot/blob/master/telebot-logo.jpg?raw=true)

* [Features](#features)
* [Requirements](#requirements)
* [Installation](#installation)
* [Usage](#usage)
    * [Quick start](#quick-start)
    * [Step by step](#step-by-step)
        * [Creating the bot](#creating-the-bot)
        * [Methods](#methods)
            * [Text messages](#text-messages)
            * [Media messages](#media-messages)
            * [Media group messages](#media-group-messages)
            * [Polls](#polls)
            * [Files](#files)
        * [Special channels](#special-channels)
            * [Update type channels](#update-type-channels)
            * [Chat channels](#chat-channels)
            * [Channels priority](#channels-priority)
* [License](#license)

---------------------------------

## Features
* Fast and reliable
* Full support for [telegram bot api](https://core.telegram.org/bots/api)
* Highly customizable.
* Automatic poll management
* [Special channels](#special-channels) : [Update type channels](#update-type-channels) for each type of update (polls,inline queries,messages, ...) and [Chat channels](#chat-channels) to manage each chat seperately and track the progress in each chat.
* Webhook support. (in development, not released yet) 

---------------------------------

## Requirements
  * Go 1.17 or higher.

---------------------------------

## Installation
 Install the package into your [$GOPATH](https://github.com/golang/go/wiki/GOPATH "GOPATH") with the [go command](https://golang.org/cmd/go/ "go command") from terminal :
 ```
 $ go get -u github.com/SakoDroid/telebot
 ```
 Git needs to be installed on your computer.

 --------------------------------

 ## Usage

 ### Quick start

 The following code creates a bot and starts receving updates. If the update is a text message that contains "hi" the bot will respond "hi to you too!".

 ```
 import (
    "fmt"
    
	bt "github.com/SakoDroid/telebot"
	cfg "github.com/SakoDroid/telebot/configs"
	objs "github.com/SakoDroid/telebot/objects"
 )

 func main(){
    up := cfg.DefaultUpdateConfigs()
    
    cf := cfg.BotConfigs{BotAPI: cfg.DefaultBotAPI, APIKey: "your api key", UpdateConfigs: up, Webhook: false, LogFileAddress: cfg.DefaultLogFile}

    bot, err := bt.NewBot(&cf)

    if err == nil{

        err == bot.Run()

        if err == nil{
            go start(bot)
        }
    }
 }

 func start(bot *bt.Bot){

     updateChannel := bot.GetUpdateChannel()

     for {
         update := <- updateChannel

         if update.Message.Text == "hi" {
            chatId := update.Message.Chat.Id
	        messageId := update.Message.MessageId
            _,err := bot.SendMessage(chatId,"hi to you too!","",messageId,false)
            if err != nil{
                fmt.println(err)
            }
         }
     }
 }
 ```
 ### Step by step

#### **Creating the bot**
 First you need to import required libraries :

 ```
 import (
    bt "github.com/SakoDroid/telebot"
    cfg "github.com/SakoDroid/telebot/configs"
    objs "github.com/SakoDroid/telebot/objects"
 )
 ```

 Then you need to create bot configs. **BotConfigs** struct is located in configs package and contains these fields :

 ```
 /*This is the bot api server. If you dont have a local bot api server, use "configs.DefaultBotAPI" for this field.*/

 BotAPI string

 /*The API key for your bot. You can get the api key (token) from botfather*/

 APIKey string

 /*The settings related to getting updates from the api server. This field shoud only be populated when Webhook field is false, otherwise it is ignored.*/

 UpdateConfigs *UpdateConfigs

 /*This field idicates if webhook should be used for receiving updates or not.
 Recommend : false*/

 Webhook bool

 /*All the logs related to bot will be written in this file. You can use configs.DefaultLogFile for default value*/

 LogFileAddress string
```

 * **Note** : telebot library currently does not support webhooks so Webhook field should always be *false*.

To create bot configs you need an UpdateConfigs to populate related field in BotConfigs. **UpdateConfigs** struct contains following fields :

```
/*Limits the number of updates to be retrieved. Values between 1-100 are accepted. Defaults to 100.*/

 Limit int

 /*Timeout in seconds for long polling. Defaults to 0, i.e. usual short polling. Should be positive, short polling should be used for testing purposes only.*/

 Timeout int

 /*List of the update types you want your bot to receive. For example, specify [“message”, “edited_channel_post”, “callback_query”] to only receive updates of these types. See Update for a complete list of available update types. Specify an empty list to receive all update types except chat_member (default). If not specified, the previous setting will be used.
 Please note that this parameter doesnt affect updates created before the call to the getUpdates, so unwanted updates may be received for a short period of time.*/

 AllowedUpdates []string

 /*This field indicates the frequency to call getUpdates method. Default is one second*/

 UpdateFrequency time.Duration
 ```
 You can use **`configs.DefaultUpdateConfigs()`** to create default update configs. Otherwise you can create your own custom update configs.

 After you have created BotConfigs you can create the bot by passing the `BotConfigs` struct you've created to **NewBot** method located in **telebot** package. After bot is created call **Run()** method and your bot will start working and will receive updates from the api server: 
 ```
 import (
	bt "github.com/SakoDroid/telebot"
	cfg "github.com/SakoDroid/telebot/configs"
	objs "github.com/SakoDroid/telebot/objects"
 )

 func main(){
    up := cfg.DefaultUpdateConfigs()
    
    cf := cfg.BotConfigs{BotAPI: cfg.DefaultBotAPI, APIKey: "your api key", UpdateConfigs: up, Webhook: false, LogFileAddress: cfg.DefaultLogFile}

    bot, err := bt.NewBot(&cf)
    if err == nil{
        err == bot.Run()
        if err == nil{
            //Do anything you want with the bot.
        }
    }
 }
```

Now that the bot is running it will receive updates from api server and passes them into UpdateeChannel. So you can use this channel to know if an update is received from api server. You can get the channel via **GetUpdateChannel()** method of the bot :

 ```
 import (
	bt "github.com/SakoDroid/telebot"
	cfg "github.com/SakoDroid/telebot/configs"
	objs "github.com/SakoDroid/telebot/objects"
 )

 func main(){
    up := cfg.DefaultUpdateConfigs()
    
    cf := cfg.BotConfigs{BotAPI: cfg.DefaultBotAPI, APIKey: "your api key", UpdateConfigs: up, Webhook: false, LogFileAddress: cfg.DefaultLogFile}

    bot, err := bt.NewBot(&cf)
    if err == nil{
        err == bot.Run()
        if err == nil{
            go start(bot)
        }
    }
 }

 func start(bot *bt.Bot){
     updateChannel := bot.GetUpdateChannel()
     for {
         update := <- updateChannel
         //Do your own processing.
     }
 }
```
#### **Methods**

 To send back text or media (such as photo, video, gif, ...) you can use Send methods. There are several send methods such as **SendMessage** and **SendPhoto**. There is two ways to send back data to the client. First way is using unique chat ids (which are integers that are unique for each chat) to send data to private chats, groups and supergroups. Second way is using chat username which can be used to send back data to supergroups (with username) and channels. Methods that use username as chat identificator end with `UN`.
 
 We will cover some of the methods below. All these methods are fully documented in the source code and will be described here briefly. In all methods you can ignore `number` arguments (int or float) by passing 0 and ignore `string` arguments by passing empty string ("").
  * **Note** : All bot methods are simplified to avoid unnecessary arguments. To access more options for each method you can call `AdvancedMode()` method of the bot that will return an advanced version of bot which will give you full access.

 #### **Text messages**

 To send back text you can use **SendMessage** (chat id) or **SendMessageUN** (username). 

 #### **Media messages**

 To send media types such as photo,video,gif,audio,voice,videonote,mpeg4 gif,sticker and document you can use their specified method. In general there are three ways to send media :
 
 1. **By file id** : File id is a unique id for a file that already exists in telegram servers. [Telegram bot api documentation](https://core.telegram.org/bots/api) recommends using file id.
 2. **By URL** : You can pass an HTTP url to send. The file will be downloaded in telegram servers and then it will be sent to the specified chat.
 3. **By file** : You can send a file on your computer. The file will be uploaded to telegram servers and then it will be sent to the specified chat.

 Calling each media sending related method returnes a MediaSender. MediaSender has all methods that are needed to send a media. For example lets send photo in our computer :

 ```
 photoFile,err := os.Open("photo.jpg")

 if err == nil{

    ms := bot.SendPhoto(chatId, messageId, "custom caption", "")

    _,err = ms.SendByFile(photoFile,false)

    if err != nil{
        fmt.Println(err)
    }

 }
 ```
 
 #### **Media group messages**

 To send a group of medias (aka albums) first you need to create a *`MediaGroup`* by calling `CreateAlbum(replyto int)` method of the bot. MediaGroup has several methods for adding photo,video,audio and other media types to the album. Keep in mind that according to [Telegram bot api documentation about media groups](https://core.telegram.org/bots/api#sendmediagroup), documents and audio files can be only grouped in an album with messages of the same type. Also the media group must include 2-10 items. The code below shows how to create a media group, add some photo to it and send it :

 ```
 mg := bot.CreateAlbum(messageId)

//Add a file on the computer.
fl,_ := os.Open("file.jpg")
 pa1,_ := mg.AddPhoto("", "", nil)
 err := pa1.AddByFile(fl)
 if err != nil{
     fmt.Println(err)
 }

//Add a photo by file id or url.
 pa2,_ ;= mg.AddPhoto("","",nil)
 err = pa2.AddByFileIdOrURL("fileId or HTTP url")
 if err != nil{
     fmt.Println(err)
 }

//Send the media group
_, err = mg.Send(chatId, false)
if err != nil {
    fmt.Println(err)
}
```

#### **Polls**

telebot library offers automatic poll management. When you create a poll and send the poll bot will receive updates about the poll. Whene you create a poll by **`CreatePoll`** method, it will return a Poll which has methods for managing the poll. You should keep the returned pointer (to Poll) somewhere because everytime an update about a poll is received the bot will process the update and update the related poll and notifies user through a [bool]channel (which you can get by calling `GetUpdateChannel` method of the poll). 

* **Note** : If an update is received that contains update about a poll and the poll is not registered with the Polls map, the given update is passed into *UpdateChannel* of the bot. Otherwise as described above, the related poll will be updated.

Let's see an example :

```

//A custom function that creates and sends a poll and listens to its updates.
func pollTest(chatId int) {

    //Creates the poll
	poll, _ := bot.CreatePoll(chatId, "How are you?", "regular")

    //Adds some options
	poll.AddOption("good")
	poll.AddOption("not bad")
	poll.AddOption("alright")
	poll.AddOption("bad")

    //Adds an explanation for the poll.
	poll.SetExplanation("This is just a test for telebot framework", "", nil)

    //Sends the poll
	err := poll.Send(false, 0)

	if err != nil {
		fmt.Println(err)
	} else {

        //Starts waiting for updates and when poll is updated, the updated result of the bot is printed.
		ch := poll.GetUpdateChannel()
		for {
			<-*ch
			fmt.Println("poll updated.")
			for _, val := range poll.GetResult() {
				fmt.Println(val.Text, ":", val.VoterCount)
			}
		}
	}
}
```

#### **Files**

You can get informations of a file that is stored in telegram servers and download it into your computer by calling **`GetFile`** method. If you want to download the file, pass true for *download* argument of the method. The below example downloads a received sticker from the user and saves it into the given file (read full documentation of the method for more information) :

```
//Receives upadate
update := <- updateChannel

//Get sticker file id
fi := update.Message.Sticker.FileId

//Open a file in the computer.
fl, _ := os.OpenFile("sticker.webp", os.O_CREATE|os.O_WRONLY, 0666)

//Gets the file info and downloads it.
_, err := bot.GetFile(fi, true, fl)
if err != nil {
    fmt.Println(err)
}
fl.Close()

```
#### **Special channels**

In telebot you can register special channels. Special channels are channels for a specified update. Currently there are two types of special channels :
1. Update type channels : These are channels for a specified update type received from api server. ( read fully in [update type channels](#update-type-channels) )
2. Chat channels : These channels can be used to get only updates for a specified chat. ( read fully in [chat channels](#chat-channels))

#### **Update type channels**
In telebot you can register special channels for a specified type of update.Updates received from api server can have `message` field, `edited_message` field, `inline_query` field and some other fields ( you can see them [here](https://core.telegram.org/bots/api#update) ). As described in [telegram bot api](https://core.telegram.org/bots/api) updates received from api server will have only one of these fields. To have easier processing and erase the part where you have to check all of the field to see what kind of update is received, we have created special channels. Special channels can be used to get notified whenever a specified kind of update is received.  This feature let's you have easier processing for each type of update. This feature is included in the advanced mode so for activating it follow these steps : 
1. Call the `AdvancedMode()` to have AdvancedBot. 
2. Then call the `Register[field name]Channel()` ( like `RegisterMessageChannel()` or `RegisterInlineQueryChannel()` ) method of the AdvancedBot. This methods will register a channel for that update type and return the channel.

**Notes :**
1. When you register a channel for a specified update type, all the received updates that contain that field will be passed into this channel and none of them will be passed into general update channel anymore.
2. When a register method is called, returned channel is permenant. Meaning that further calls of the same method will return the *same channel* not a new one.
3. Read about [channels priority](#channels-priority).

**An example :**

For example you want to get notified whenever an update is received that contains message field. To do this we follow the above steps and write the below code :

```
//General update channel.
ch := *bot.GetUpdateChannel()

//Register message channel.
mch := *bot.AdvancedMode().RegisterMessageChannel()

for {
    select {

    //This will be triggered when an update is passed into general update channel.
    case update := <-ch:
        fmt.Println("this is update channel")
        //Processing the update ...
    
    //This will be triggered when an update containing message field is received and the message field is passed into the channel.
    case message := <-mch:
        fmt.Println("this is message channel")
        fmt.Println(message.Text)
    }
}
```

#### **Chat channels**

Chat channels can be used to get the updates for a specified chat. When an update is received from api server, if it belongs to a chat (some updates don't contain chat info) and that chat has a registered channel the update will be passed into the chat channel. This feature belongs to advanced mode so for using it first you need to activate the advanced mode. To use this feature follow the following steps :
1. Call the `AdvancedMode()` to have AdvancedBot. 
2. Then call the `RegisterChatChannel(chatId string)` method of the AdvancedBot. This methods will register a channel for the specified chat and return the channel.
3. Use the channel to receive updates for the chat.

**Notes :**
1. When you register a channel for a specified chat, all the received updates that belong to that chat be passed into this channel and none of them will be passed into general update channel or update type channels anymore.
2. When a register method is called, returned channel is permenant. Meaning that further calls of the same method will return the *same channel* not a new one.
3. Read about [channels priority](#channels-priority).

**An example :**

In the below code first an update is received from the general update channel. If the update is for a private chat then a method (startPrivateChat) is called to do further processing. The *startPrivateChat* method will ask for a location from the user and prints the received location and if location is not received warns the user to send a location :

```
import (
	bt "github.com/SakoDroid/telebot"
	cfg "github.com/SakoDroid/telebot/configs"
    objs "github.com/SakoDroid/telebot/objects"
)

func main(){
    //Start the bot, described above.
    up := cfg.DefaultUpdateConfigs()

    cf := cfg.BotConfigs{BotAPI: cfg.DefaultBotAPI, APIKey: "your api key", UpdateConfigs: up, Webhook: false, LogFileAddress: cfg.DefaultLogFile}

    bot, err := bt.NewBot(&cf)
    if err == nil{
        err == bot.Run()
        if err == nil{
            go start(bot)
        }
    }
}

func start(bot *bt.Bot){
    //Get the update channel
    updateChannel := bot.GetUpdateChannel()

    for {
        //Receives an update from general update channel.
        update := <- updateChannel
        
        //Checks if the update is a message
        if update.Message != nil {

            //Checks if the message is a text message.
            if update.Message.Text == "" {
                continue
            }

            //Checks if the chat is a private chat
            if update.Message.Chat.Type == "private" {
                go startPrivateChat(update)
            }
        } else {
            fmt.Println("not message received")
        }

    }
}

func specialUser(update *objs.Update) {
    //Informations of the message and chat
	chatId := update.Message.Chat.Id
	messageId := update.Message.MessageId

    //Registers a channel for the chat.
	ch := bot.AdvancedMode().RegisterChatChannel(strconv.Itoa(chatId))

    //Sends a message to the user and asks for a location
	_, err := bot.SendMessage(chatId, "send a location", "", messageId, false)
	if err != nil {
		fmt.Println(err)
	}

	for up := <-*ch;; {

        //Gets the location
		loc := up.Message.Location

        //Checks if the message is a location.
		if loc != nil{

            //Prints the latitude and longitude
            fmt.Println("latitude :",loc.Latitude,", longitude :",loc.Longitude)
			break

		}else{
			_, err := bot.SendMessage(chatId, "Please send a location", "", up.Message.MessageId, false)
			if err != nil {
				fmt.Println(err)
			}
		}

	}
}

```

#### **Channels priority :**

Since different types of channels may get involved it's important to know the priority of channels.Meaning when an update is received which channels have higher priority to have the update passed into them. Basically this is how channels are prioritized :

1. Chat channels
2. Update type channels
3. General channels

When an update is received, first it is checked to see if it has chat information and if it does it will be passed into the related channel if it is registered. If this step fails (does not have chat information or no channel is registered for the chat) then the *update type channels* are checked and if the update contains a field that does have a channel registered for it the related field will be passed into the channel.(For example if the update contains message field and you have called `RegisterMessageChannel()` method, the message field will be passed into the channel). If this step fails too then the update will be passed into general update channel. 

To summarize :

```
Update is received -> Chat channels
                             |
                             |
if chat channel check fails  |
                             |
                             |----------> Update type channels
                                                   |
                                                   |
                    if update channel check fails  |
                                                   |
                                                   |----------> General update channel
                              
```
---------------------------

## License

telebot is licensed under [MIT lisence](https://en.wikipedia.org/wiki/MIT_License). Which means it can be used for commerical and private apps and can be modified.

---------------------------

![telebot logo inspired by Golang logo](https://github.com/SakoDroid/telebot/blob/master/telebot-logo.jpg?raw=true)
