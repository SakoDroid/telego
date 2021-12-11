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

/*Client used for sending http requests to bot api*/
type httpSenderClient struct {
	botApi, apiKey string
}

/*This method sends an http request (without processing the response) as application/json. Returns the body of the response.*/
func (hsc *httpSenderClient) sendHttpReqJson(method string, args objs.MethodArguments) ([]byte, error) {
	bd := args.ToJson()
	res, err2 := hsc.sendHttpReq(method, "application/json", bd)
	if err2 != nil {
		return nil, &errs.MethodNotSentError{Method: method, Reason: err2.Error()}
	}
	if res.StatusCode < 300 {
		out := make([]byte, res.ContentLength)
		_, err3 := res.Body.Read(out)
		if err3 != nil {
			return nil, &errs.MethodNotSentError{Method: method, Reason: "unable to parse body into byte slice. " + err3.Error()}
		}
		return out, nil
	} else {
		return nil, &errs.MethodNotSentError{Method: method, Reason: "received status code " + strconv.Itoa(res.StatusCode)}
	}
}

/*This method sends an http request (without processing the response) as multipart/formdata. Returns the body of the response.
This method is only used for uploading files to bot api server.*/
func (hsc *httpSenderClient) sendHttpReqMultiPart(method string, file *os.File, args objs.MethodArguments) ([]byte, error) {
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
			out := make([]byte, res.ContentLength)
			_, err3 := res.Body.Read(out)
			if err3 != nil {
				return nil, &errs.MethodNotSentError{Method: method, Reason: "unable to parse body into byte slice. " + err3.Error()}
			}
			return out, nil
		} else {
			return nil, &errs.MethodNotSentError{Method: method, Reason: "received status code " + strconv.Itoa(res.StatusCode)}
		}
	} else {
		return nil, &errs.MethodNotSentError{Method: method, Reason: "unable to add file to the multipart form. " + err.Error()}
	}
}

func (hsc *httpSenderClient) addFileToMultiPartForm(file *os.File, wr *mp.Writer, fieldName string) error {
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

func (hsc *httpSenderClient) sendHttpReq(method, contetType string, body []byte) (*http.Response, error) {
	cl := http.Client{}
	req, err := http.NewRequest("POST", hsc.botApi+hsc.apiKey+"/"+method, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Add(textproto.CanonicalMIMEHeaderKey("content-type"), contetType)
	req.Header.Add(textproto.CanonicalMIMEHeaderKey("content-length"), strconv.Itoa(len(body)))
	return cl.Do(req)
}
