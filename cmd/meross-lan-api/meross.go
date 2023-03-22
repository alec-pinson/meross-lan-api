package main

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"strconv"
	"time"
)

type Request struct {
	Header  Header         `json:"header"`
	Payload RequestPayload `json:"payload"`
}

type Header struct {
	// MessageID is used to identify a request/response pair.
	From           string `json:"from,omnitempty"`
	MessageID      string `json:"messageId,omitempty"`
	Method         string `json:"method,omitempty"`
	Namespace      string `json:"namespace,omitempty"`
	PayloadVersion int    `json:"payloadVersion,omitempty"`
	Signature      string `json:"sign,omitempty"`
	Timestamp      int64  `json:"timestamp,omitempty"`
	TimestampMs    int64  `json:"timestampMs,omitempty"`
}

type RequestPayload struct {
	ToggleX struct {
		Channel int `json:"channel,omitempty"`
		OnOff   int `json:"onoff,omitempty"`
	}
	All struct{}
}

type Response struct {
	Header  Header          `json:"header"`
	Payload ResponsePayload `json:"payload"`
}

type ResponsePayload struct {
	All struct {
		System struct {
			Type string `json:"type,omitempty"`
		}
		Digest struct {
			Togglex []struct {
				Channel int `json:"channel,omitempty"`
				OnOff   int `json:"onoff,omitempty"`
			}
		}
	}
}

func getStatus(deviceIp string, channel int, key string) int {
	body := generateBody(key)
	body.Header.Method = "GET"
	body.Header.Namespace = "Appliance.System.All"
	body.Header.From = "Meross"
	response, err := sendRequest(deviceIp, body)
	if err != nil {
		log.Print(err)
	}

	for i := 0; i < len(response.Payload.All.Digest.Togglex); i++ {
		if response.Payload.All.Digest.Togglex[i].Channel == channel {
			return response.Payload.All.Digest.Togglex[i].OnOff
		}
	}

	return 2 // should always return 0 or 1 so we'll class this as 'unknown'
}

func turnOn(deviceIp string, channel int, key string) {
	body := generateBody(key)
	body.Header.Method = "SET"
	body.Header.Namespace = "Appliance.Control.ToggleX"
	body.Payload.ToggleX.Channel = channel
	body.Payload.ToggleX.OnOff = 1
	_, err := sendRequest(deviceIp, body)
	if err != nil {
		log.Print(err)
	}
}

func turnOff(deviceIp string, channel int, key string) {
	body := generateBody(key)
	body.Header.Method = "SET"
	body.Header.Namespace = "Appliance.Control.ToggleX"
	body.Payload.ToggleX.Channel = channel
	body.Payload.ToggleX.OnOff = 0
	_, err := sendRequest(deviceIp, body)
	if err != nil {
		log.Print(err)
	}
}

func sendRequest(deviceIp string, body Request) (Response, error) {
	var result Response
	// convert body to bytes
	bodyJSON, _ := json.Marshal(body)

	if config.Debug {
		return result, nil
	}

	// create and send request
	req, err := http.NewRequest("POST", "http://"+deviceIp+"/config", bytes.NewBuffer(bodyJSON))
	req.Header.Add("Content-Type", "application/json; charset=UTF-8")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	// read response body
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return result, err
	}

	// if non-200 status code or debug enabled then dump out request/response info
	if config.Debug || resp.StatusCode != 200 {
		requestDump, err := httputil.DumpRequest(req, true)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("Request\n: %v", string(requestDump))
		fmt.Printf("%+v\n\n", body)

		responseDump, err := httputil.DumpResponse(resp, true)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("Response\n: %v\n\n", string(responseDump))
	}

	// convert response to Response type
	if err := json.Unmarshal(respBody, &result); err != nil { // Parse []byte to go struct pointer
		fmt.Println(err)
	}
	return result, nil
}

// random hex code string, length 32
func generateMessageId() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes)[:32], nil
}

func generateBody(key string) Request {
	var body Request

	body.Header.PayloadVersion = 1
	body.Header.MessageID, _ = generateMessageId()
	body.Header.Timestamp = time.Now().Unix()
	body.Header.Signature = fmt.Sprintf("%x", md5.Sum([]byte(body.Header.MessageID+key+strconv.Itoa(int(body.Header.Timestamp)))))

	if config.Debug {
		log.Printf("Body: %+v\n", body)
	}
	return body
}
