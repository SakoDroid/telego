package tba

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	cfg "github.com/SakoDroid/telego/configs"
	log "github.com/SakoDroid/telego/logger"
	objs "github.com/SakoDroid/telego/objects"
	up "github.com/SakoDroid/telego/parser"
)

var configs *cfg.BotConfigs
var interfaceUpdateChannel *chan *objs.Update
var chatUpdateChannel *chan *objs.ChatUpdate

//StartWebHook starts the webhook.
func StartWebHook(cfg *cfg.BotConfigs, iuc *chan *objs.Update, cuc *chan *objs.ChatUpdate) error {
	configs = cfg
	interfaceUpdateChannel = iuc
	chatUpdateChannel = cuc
	startTheServer()
	return nil
}

func startTheServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", mainHandler)
	mux.HandleFunc("/"+configs.APIKey, handleReq)
	go func() {
		err := http.ListenAndServeTLS(":"+strconv.Itoa(configs.WebHookConfigs.Port), configs.WebHookConfigs.CertFile, configs.WebHookConfigs.KeyFile, mux)
		if err != nil {
			log.Logger.Fatalln("Webhook : Failed to start the HTTPS server.", err)
		}
	}()
}

func mainHandler(wr http.ResponseWriter, req *http.Request) {
	wr.WriteHeader(404)
	wr.Write([]byte{})
}

func handleReq(wr http.ResponseWriter, req *http.Request) {
	contentType := req.Header.Get("Content-Type")
	if contentType != "" && strings.HasSuffix(contentType, "json") {
		if req.ContentLength > 0 {
			body := make([]byte, req.ContentLength)
			_, err := req.Body.Read(body)
			if err == nil {
				update := &objs.Update{}
				jsonErr := json.Unmarshal(body, update)
				if jsonErr == nil {
					up.ParseSingleUpdate(update, interfaceUpdateChannel, chatUpdateChannel, configs)
				} else {
					log.Logger.Println("Webhook : Error parsing the update. Address :", req.RemoteAddr, ". Error :", jsonErr)
				}
			} else {
				log.Logger.Println("Webhook : Error reading the body. Address :", req.RemoteAddr, ". Error :", err)
			}
			wr.WriteHeader(200)
			wr.Write([]byte{})
		} else {
			log.Logger.Println("Webhook : Request has no body. Address :", req.RemoteAddr)
			send400(&wr, "Request has no body")
		}
	} else {
		log.Logger.Println("Webhook : \"Content-Type\" header is not json or it's missing. Address :", req.RemoteAddr)
		send400(&wr, " \"Content-Type\" header is not json or it's missing")
	}
}

func send400(wr *http.ResponseWriter, reason string) {
	(*wr).Header().Add("Content-Type", "text/plain")
	(*wr).Header().Add("Content-Length", strconv.Itoa(len(reason)))
	(*wr).WriteHeader(400)
	(*wr).Write([]byte(reason))
}
