package TBA

import (
	"encoding/json"
	"errors"
	"time"

	bot "github.com/SakoDroid/telebot"
	errs "github.com/SakoDroid/telebot/Errors"
	cfgs "github.com/SakoDroid/telebot/configs"
	objs "github.com/SakoDroid/telebot/objects"
)

var interfaceCreated = false

type botAPIInterface struct {
	botConfigs           *cfgs.BotConfigs
	updateRoutineRunning bool
	updateChannel        *chan *objs.Update
	updateRoutineChannel chan bool
	lastOffset           int
}

/*Starts the update routine to receive updates from api sever*/
func (bai *botAPIInterface) StartUpdateRoutine() error {
	if !bai.botConfigs.Webhook {
		if bai.updateRoutineRunning {
			return &errs.UpdateRoutineAlreadyStarted{}
		}
		bai.updateRoutineRunning = true
		ch := make(chan *objs.Update)
		bai.updateChannel = &ch
		bai.updateRoutineChannel = make(chan bool)
		go bai.startReceiving()
		return nil
	} else {
		return errors.New("Webhook option is true.")
	}
}

/*Stops the update routine*/
func (bai *botAPIInterface) StopUpdateRoutine() {
	if bai.updateRoutineRunning {
		bai.updateRoutineRunning = false
		bai.updateRoutineChannel <- true
	}
}

func (bai *botAPIInterface) startReceiving() {
	cl := httpSenderClient{botApi: bai.botConfigs.BotAPI, apiKey: bai.botConfigs.APIKey}
loop:
	for {
		select {
		case <-bai.updateRoutineChannel:
			break loop
		default:
			args := objs.GetUpdatesArgs{Offset: bai.lastOffset, Limit: bai.botConfigs.UpdateConfigs.Limit, Timeout: bai.botConfigs.UpdateConfigs.Timeout}
			if bai.botConfigs.UpdateConfigs.AllowedUpdates != nil {
				args.AllowedUpdates = bai.botConfigs.UpdateConfigs.AllowedUpdates
			}
			res, err := cl.sendHttpReqJson("getUpdates", &args)
			if err != nil {
				bot.Logger.Println("Error receiving updates.", err)
				continue
			}
			err = bai.parseUpdateresults(res)
			if err != nil {
				bot.Logger.Println("Error parsing the result of the update. " + err.Error())
			}
		}
		time.Sleep(bai.botConfigs.UpdateConfigs.UpdateFrequency)
	}
}

func (bai *botAPIInterface) parseUpdateresults(body []byte) error {
	ur := &objs.UpdateResult{}
	err := json.Unmarshal(body, ur)
	if err != nil {
		return err
	}
	if !ur.Ok {
		return &errs.UpdateNotOk{Offset: bai.lastOffset}
	}
	for _, val := range ur.Result {
		if val.Update_id > bai.lastOffset {
			bai.lastOffset = val.Update_id
		}
		(*bai.updateChannel) <- &val
	}
	return nil
}

/*This method returns an iterface to communicate with the bot api.
If the updateFrequency argument is not nil, the update routine begins automtically*/
func CreateInterface(botCfg *cfgs.BotConfigs) (*botAPIInterface, error) {
	if interfaceCreated {
		return nil, &errs.BotInterfaceAlreadyCreated{}
	}
	interfaceCreated = true
	temp := &botAPIInterface{botConfigs: botCfg}
	return temp, nil
}
