package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/pkg/errors"
)

//リクエストデータタイプ

type RequestBody struct {
	Message string `json: "message"`
}
type ResponseBody struct {
	Messages []string `json:"messages"`
}

func flaskHandler(message string) ([]byte, error) {
	flaskPath := os.Getenv("FLASK_PATH")

	body := new(RequestBody)
	body.Message = message

	body_json, err := json.Marshal(body)
	if err != nil {
		log.Println(err)
		return nil, errors.Wrap(err, "failed to parse json")
	}

	res, err := http.Post(flaskPath, "application/json", bytes.NewBuffer(body_json))

	defer res.Body.Close()

	if err != nil {
		fmt.Println("[!]" + err.Error())
		return nil, errors.Wrap(err, "failed to request API")
	}
	res_body, err := ioutil.ReadAll(res.Body)

	str_json := string(res_body)
	fmt.Println(str_json)

	messages := new(ResponseBody)
	err = json.Unmarshal([]byte(str_json), &messages)

	if err != nil {
		fmt.Println(err)
		return nil, errors.Wrapf(err, "failed to convert string to json")
	}

	//レスポンスを作成

	r := Messages{
		Messages: messages.Messages,
		Action:   "LOVE_MESSAGE",
	}

	b, err := json.Marshal(r)
	if err != nil {
		log.Println("cannot marshal struct: %v", err)
		return nil, err
	}
	return b, nil

}
