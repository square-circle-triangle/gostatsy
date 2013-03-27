package statsy

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type Statsy struct {
	ApiKey    string
	SecretKey string
	BaseUrl   string
	expires   int64
}

type EventSubmit struct {
	ApiKey    string       `json:"api_key"`
	Expires   int64        `json:"expires"`
	Signature string       `json:"signature"`
	Events    []*JsonEvent `json:"events"`
}

type StatsyResponse struct {
	Error string `json:"error"`
}

func New(ApiKey string, SecretKey string) *Statsy {
	client := Statsy{ApiKey: ApiKey, SecretKey: SecretKey, BaseUrl: "http://statsyapp.com/"}

	return &client
}

func (s *Statsy) Sign(expires int64, stream string) string {
	h := sha1.New()
	signString := s.SecretKey + "POST" + stream + strconv.Itoa(int(expires))
	io.WriteString(h, signString)

	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func (s *Statsy) Send(events []Event) error {
	// convert events to json events, so we can re-arrange the struct in to something that marshalls correctly
	expires := time.Now().Add(5 * time.Minute).Unix()
	if s.expires > 0 {
		expires = s.expires
	}

	eventSubmit := EventSubmit{ApiKey: s.ApiKey, Expires: expires, Signature: s.Sign(expires, "")}
	for _, event := range events {
		eventSubmit.Events = append(eventSubmit.Events, event.JsonEvent())
	}
	b, err := json.Marshal(eventSubmit)
	if err != nil {
		return err
	}

	res, err := http.Post(s.BaseUrl+"api/v1/multiple_event", "application/json", bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		statsyResponse := StatsyResponse{}
		responseBytes, err := ioutil.ReadAll(res.Body)

		if err != nil {
			return err
		}
		err = json.Unmarshal(responseBytes, &statsyResponse)

		if err != nil {
			return err
		}

		err = errors.New(statsyResponse.Error)
		return err
	}

	return nil
}

func (s *Statsy) Increment(stream string, weight float32) error {
	events := make([]Event, 0)
	events = append(events, Event{Stream: stream, Weight: weight})
	return s.Send(events)
}
