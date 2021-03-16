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
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

//リクエストデータタイプ

type RequestBody struct {
	MeidoMessage string `json:"message"`
}
type ResponseBody struct {
	Messages []string `json:"messages"`
	Score    int      `json:"score"`
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
	uuidV4 := uuid.New()
	if 0 <= randNum && randNum < PARAM {
		//好意的な文章

		messages, err := readMessage(LIKE_FILE_NAME)
		if err != nil {
			return "だいすきだよ", LIKE
		}

		//DB書き込み
		addCertUser(acceptTarget, uuidV4.String())

		return messages[rand.Intn(len(messages))], LIKE

	} else {
		//否定的な文章を読み取る
		messages, err := readMessage(DISLIKE_FILE_NAME)
		addCertUser(deniedTarget, uuidV4.String())

		//サイコロゲーム
		if err != nil {
			return "だいきらいだよ", DISLIKE
		}

		return messages[rand.Intn(len(messages))], DISLIKE
	}
	//return "これつくったひとむのう", DISLIKE
}
func getRuneAt(s string, i int) rune {
	rs := []rune(s)
	return rs[i]
}
func flaskHandler(message string) ([]byte, error) {

	flaskPath := os.Getenv("FLASK_URL")
	body := new(RequestBody)
	generateMessage, likeType := selectMessage()

	fmt.Println(generateMessage)
	fmt.Println(likeType)

	body.MeidoMessage = generateMessage

	body_json, err := json.Marshal(body)
	fmt.Println("test")
	fmt.Println(body_json)
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
	var certMessage string
	if likeType == LIKE {
		certMessage = acceptMessage
	} else {
		certMessage = deniedMessage
	}

	for i := 0; i < utf8.RuneCountInString(generateMessage); i++ {
		fmt.Println(string(getRuneAt(generateMessage, i)) + "：　" + messages.Messages[i])
		messages.Messages[i] = string(getRuneAt(generateMessage, i)) + "：　" + messages.Messages[i]
	}

	r := FlaskMessages{
		Messages:        messages.Messages,
		OriginalMessage: generateMessage,
		SendMessage:     message,
		Score:           messages.Score,
		CertMessage:     certMessage,
		Action:          "LOVE_MESSAGE",
	}

	b, err := json.Marshal(r)
	if err != nil {
		log.Println("cannot marshal struct: %v", err)
		return nil, err
	}
	return b, nil

}
