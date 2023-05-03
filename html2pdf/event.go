package main

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"log"

	"github.com/programmfabrik/golib"
)

func sendEvent(ev event) {
	if info.ApiURL == "" {
		return
	}
	var err error
	defer func() {
		if err != nil {
			log.Print(err.Error())
		}
	}()
	ev.Basetype = "event"
	req, err := http.NewRequest(http.MethodPost, info.ApiURL+"/api/v1/event", bytes.NewReader([]byte(
		golib.JsonString(eventOuter{
			Basetype: "event",
			Event:    ev},
		),
	)))
	if err != nil {
		return
	}
	req.Header.Set("Authorization", "Bearer "+info.ApiToken)
	var res *http.Response
	res, err = http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusOK {
		log.Printf("event %q sent: %#v", ev.Type, ev)
	}
	var body []byte
	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		log.Print(err.Error())
	}
	log.Print(string(body))
}

type eventOuter struct {
	Basetype string `json:"_basetype"`
	Event    event  `json:"event"`
}

type event struct {
	Type           string         `json:"type,omitempty"`
	ObjectId       int64          `json:"object_id,omitempty"`
	ObjectVersion  int64          `json:"object_version,omitempty"`
	GlobalObjectId string         `json:"global_object_id,omitempty"`
	Schema         string         `json:"schema,omitempty"`
	Basetype       string         `json:"basetype,omitempty"`
	Info           map[string]any `json:"info,omitempty"`
}
