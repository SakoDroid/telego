package tba

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	cfg "github.com/SakoDroid/telego/v2/configs"
	log "github.com/SakoDroid/telego/v2/logger"
	objs "github.com/SakoDroid/telego/v2/objects"
	up "github.com/SakoDroid/telego/v2/parser"
)

type Webhook struct {
	configs          *cfg.BotConfigs
	isSecretTokenSet bool
	parser           *up.UpdateParser
}

// StartWebHook starts the webhook.
func (w *Webhook) StartWebHook(cfg *cfg.BotConfigs, parser *up.UpdateParser) error {
	w.configs = cfg
	// interfaceUpdateChannel = iuc
	// chatUpdateChannel = cuc
	w.isSecretTokenSet = cfg.WebHookConfigs.SecretToken != ""
	w.parser = parser
	w.startTheServer()
	return nil
}

func (w *Webhook) startTheServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", w.mainHandler)
	mux.HandleFunc("/"+w.configs.APIKey, w.handleReq)
	go func() {
		err := http.ListenAndServeTLS(":"+strconv.Itoa(w.configs.WebHookConfigs.Port), w.configs.WebHookConfigs.CertFile, w.configs.WebHookConfigs.KeyFile, mux)
		if err != nil {
			log.Logger.Fatalln("Webhook : Failed to start the HTTPS server.", err)
		}
	}()
}

func (w *Webhook) mainHandler(wr http.ResponseWriter, req *http.Request) {
	wr.WriteHeader(404)
	wr.Write([]byte{})
}

func (w *Webhook) handleReq(wr http.ResponseWriter, req *http.Request) {
	if w.isSecretTokenSet {
		token := req.Header.Get("X-Telegram-Bot-Api-Secret-Token")
		if token != w.configs.WebHookConfigs.SecretToken {
			wr.WriteHeader(403)
			wr.Write([]byte{})
			return
		}
	}
	contentType := req.Header.Get("Content-Type")
	if contentType != "" && strings.HasSuffix(contentType, "json") {
		if req.ContentLength > 0 {
			body := make([]byte, req.ContentLength)
			_, err := req.Body.Read(body)
			if err == nil {
				update := &objs.Update{}
				jsonErr := json.Unmarshal(body, update)
				if jsonErr == nil {
					// up.ParseSingleUpdate(update, interfaceUpdateChannel, chatUpdateChannel, configs)
					w.parser.ExecuteChain(update)
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
			w.send400(&wr, "Request has no body")
		}
	} else {
		log.Logger.Println("Webhook : \"Content-Type\" header is not json or it's missing. Address :", req.RemoteAddr)
		w.send400(&wr, " \"Content-Type\" header is not json or it's missing")
	}
}

func (w *Webhook) send400(wr *http.ResponseWriter, reason string) {
	(*wr).Header().Add("Content-Type", "text/plain")
	(*wr).Header().Add("Content-Length", strconv.Itoa(len(reason)))
	(*wr).WriteHeader(400)
	(*wr).Write([]byte(reason))
}
