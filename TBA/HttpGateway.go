package TBA

import (
	"bytes"
	"io"
	"net/http"
	"net/textproto"
	"os"
	"strconv"

	mp "mime/multipart"

	errs "github.com/SakoDroid/telebot/Errors"
	objs "github.com/SakoDroid/telebot/objects"
)

type HttpSenderClient struct {
	BotApi, APIKey string
}

type HttpRecevierClient struct {
	BotApi, APIKey string
}

/*This method sends an http request (without processing the response) as application/json. Returns the body of the response.*/
func (hsc *HttpSenderClient) SendHttpReqJson(method string, args objs.MethodArguments) (io.Reader, error) {
	bd := args.ToJson()
	res, err2 := hsc.sendHttpReq(method, "application/json", bd)
	if err2 != nil {
		return nil, &errs.MethodNotSentError{Method: method, Reason: err2.Error()}
	}
	if res.StatusCode < 300 {
		return res.Body, nil
	} else {
		return nil, &errs.MethodNotSentError{Method: method, Reason: "received status code " + strconv.Itoa(res.StatusCode)}
	}
}

func (hsc *HttpSenderClient) SendHttpReqMultiPart(method string, file *os.File, args objs.MethodArguments) (io.Reader, error) {
	body := &bytes.Buffer{}
	writer := mp.NewWriter(body)
	args.ToMultiPart(writer)
	err := hsc.addFileToMultiPartForm(file, writer, args.GetMediaType())
	if err == nil {
		_ = writer.Close()
		bts := body.Bytes()
		res, err2 := hsc.sendHttpReq(method, writer.FormDataContentType(), bts)
		if err2 != nil {
			return nil, &errs.MethodNotSentError{Method: method, Reason: err2.Error()}
		}
		if res.StatusCode < 300 {
			return res.Body, nil
		} else {
			return nil, &errs.MethodNotSentError{Method: method, Reason: "received status code " + strconv.Itoa(res.StatusCode)}
		}
	} else {
		return nil, &errs.MethodNotSentError{Method: method, Reason: "unable to add file to the multipart form. " + err.Error()}
	}
}

func (hsc *HttpSenderClient) addFileToMultiPartForm(file *os.File, wr *mp.Writer, fieldName string) error {
	fileStat, err := file.Stat()
	if err != nil {
		return err
	}
	fw, err2 := wr.CreateFormFile(fieldName, fileStat.Name())
	if err2 != nil {
		return err2
	}
	_, err3 := io.Copy(fw, file)
	if err3 != nil {
		return err3
	}
	return nil
}

func (hsc *HttpSenderClient) sendHttpReq(method, contetType string, body []byte) (*http.Response, error) {
	cl := http.Client{}
	req, err := http.NewRequest("POST", hsc.BotApi+hsc.APIKey+"/"+method, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Add(textproto.CanonicalMIMEHeaderKey("content-type"), contetType)
	req.Header.Add(textproto.CanonicalMIMEHeaderKey("content-length"), strconv.Itoa(len(body)))
	return cl.Do(req)
}

func (hrc *HttpRecevierClient) ReceiveUpdates() (io.Reader, error) {
	res, err := http.Get(hrc.BotApi + hrc.APIKey + "/getUpdates")
	if err != nil {
		return nil, err
	}
	return res.Body, nil
}
