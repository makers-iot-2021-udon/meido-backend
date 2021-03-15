package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/pkg/errors"
)

//リクエストデータタイプ

type RequestBody struct {
	Message string `json: "message"`
}
type ResponseBody struct {
	Messages []string `json:"messages"`
}

const LIKE = 0
const DISLIKE = 1
const PARAM = 7
const LIKE_FILE_NAME = "./like.txt"
const DISLIKE_FILE_NAME = "./dislike.txt"

func readMessage(filename string) ([]string, error) {
	f, err := os.Open(filename)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	defer f.Close()
	var strSlice []string
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := scanner.Text()
		//fmt.Println(line)
		strSlice = append(strSlice, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return strSlice, nil
}

func selectMessage() (string, int) {

	rand.Seed(time.Now().UnixNano())

	var randNum = rand.Intn(10)

	if 0 <= randNum && randNum < PARAM {
		//好意的な文章

		messages, err := readMessage(LIKE_FILE_NAME)
		if err != nil {
			return "だいすきだよ", LIKE
		}
		return messages[rand.Intn(len(messages))], LIKE

	} else {
		//否定的な文章を読み取る
		messages, err := readMessage(DISLIKE_FILE_NAME)
		//サイコロゲーム
		if err != nil {
			return "だいきらいだよ", DISLIKE
		}

		return messages[rand.Intn(len(messages))], DISLIKE
	}
	//return "これつくったひとむのう", DISLIKE
}

func flaskHandler(message string) ([]byte, error) {
	flaskPath := os.Getenv("FLASK_URL")

	body := new(RequestBody)
	generateMessage, likeType := selectMessage()

	fmt.Println(generateMessage)
	fmt.Println(likeType)

	body.Message = generateMessage

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
