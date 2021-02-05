package main

import (
	"encoding/json"
	"log"
	"main/persistence/redis"
	"os"
	"time"

	"github.com/pkg/errors"
)

type Request struct {
	Action  string `json:"action`
	Message string `json:"message`
	Uid     string `json:"uid"`
}

// APIやサーバーの状態を扱う(ドア・メイド・認証で用いる）

type StatusMessage struct {
	Action string `json:"action"`
	Error  bool   `json:"error"`
	Status string `json:"status"`
}

type MeidoMessage struct {
	Action  string `json:"action"`
	Message string `json:"message"`
	Status  string `json:"status"`
	Error   bool   `json:"error"`
}

type Message struct {
	Action   string   `json:"action"`
	Messages []string `json:"messages"`
}

type CountMessage struct {
	Action string `json:"action"`
	Count  int64  `json:"count"`
}

const connectionTarget = "connections"
const messageTarget = "messages"
const doorTarget = "doorTarget"

var errorResponse = []byte(`{"action":"ERROR_MESSAGE","status":"NG","error": true}`)
var defaultMeidoStatus = []byte(`{"action":"MEIDO_STATUS","status":"FINE","error":false}`)

func handler(s []byte) ([]byte, bool) {
	var r Request
	if err := json.Unmarshal(s, &r); err != nil {
		return errorResponse, false
	}

	switch {
	case r.Action == "POST_DOOR":
		r, err := doorHandler()
		if err != nil {
			return errorResponse, false
		}

		return r, false
	case r.Action == "MEIDO_VOTE":
		r, err := voteHandler(r.Message)
		if err != nil {
			return errorResponse, false
		}
		return r, false

	//Todo 何らかの形で実装したい
	case r.Action == "MEIDO_STAUTS":
		return defaultMeidoStatus, false

	case r.Action == "SYSTEM_STATUS":
		return []byte(`{"action":"SYSTEM_STATUS","status":"FINE","error":false}`), false

	case r.Action == "MEIDO_COUNT":
		r, err := connectionCountHandler()
		if err != nil {
			return errorResponse, false
		}
		return r, true

	case r.Action == "GET_MESSAGE":
		r, err := messageHandler()
		if err != nil {
			return errorResponse, false
		}
		return r, false

	// こいつ使わんでもよさげ
	case r.Action == "MEIDO_FUN":
		r, err := doorHandler()
		if err != nil {
			return errorResponse, false
		}
		return r, false
	case r.Action == "MEIDO_MESSAGE":
		r, err := connectHandler()
		if err != nil {
			return errorResponse, false
		}
		return r, false
	}
	return errorResponse, false
}

func messageHandler() ([]byte, error) {
	messages, err := getMessages()
	if err != nil {
		return nil, err
	}

	r := Message{
		Messages: messages,
		Action:   "MEIDO_MESSAGE",
	}
	b, err := json.Marshal(r)
	if err != nil {
		log.Println("cannot marshal struct: %v", err)
		return nil, err
	}
	return b, nil

}

func connectionCountHandler() ([]byte, error) {
	err := addValue(connectionTarget)
	if err != nil {
		return nil, err
	}
	var _count int64 = 0
	_count, err = declValue(connectionTarget)
	//fmt.Println(_count)
	r := CountMessage{
		Action: "MEIDO_COUNT",
		Count:  _count,
	}
	b, err := json.Marshal(r)
	if err != nil {
		log.Println("cannot marshal struct: %v", err)
		return nil, err
	}
	return b, nil
}

//ドアのステータスを取得する
func doorHandler() ([]byte, error) {
	message, err := getDoorState()

	if err != nil {
		log.Println(err)
		return nil, err
	}

	r := StatusMessage{
		Action: "POST_DOOR",
		Status: message,
		Error:  false,
	}
	b, err := json.Marshal(r)
	if err != nil {

		log.Println("cannot marshal struct: %v", err)
		return nil, err
	}
	return b, nil
}

func connectHandler() ([]byte, error) {
	err := addValue(connectionTarget)
	if err != nil {
		return nil, errors.Wrap(err, "failed")
	}
	r := MeidoMessage{Action: "MEIDO_MESSAGE", Error: false, Status: "OK", Message: "ごしゅじんさま～、よおこそ"}
	b, err := json.Marshal(r)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return b, nil
}

//推しのメイドに愛のメッセージを投下する
func voteHandler(message string) ([]byte, error) {
	log.Println(message)

	//ここでポスト
	err := postMessage(message)
	if err != nil {
		return nil, errors.Wrap(err, "failed")
	}
	r := MeidoMessage{Action: "MEIDO_VOTE", Error: false, Status: "OK", Message: "ごしゅじんさま～、ありがとなす！"}
	b, err := json.Marshal(r)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return b, nil
}

//メッセージ取得
func getMessages() ([]string, error) {
	redisPath := os.Getenv("REDIS_PATH")
	client, err := redis.New(redisPath)

	if err != nil {
		return nil, errors.Wrap(err, "failed to get redis client")
	}

	defer client.Close()

	lrangeVal, err := client.LRange(messageTarget, 0, -1).Result()
	log.Println(lrangeVal)

	if err == redis.Nil {
		err = postMessage("Hello")
		if err != nil {
			return nil, err
		} else {
			getMessages()
		}
	} else if err != nil {
		return nil, errors.Wrapf(err, "failed to get redis client")
	} else {
		err = client.Do("DEL", messageTarget).Err()
		if err != nil {
			log.Println(err, "failed to delete message from redis")
		}
		return lrangeVal, nil

	}
	return nil, nil
}

//ユーザーのカウント
func addValue(target string) error {
	redisPath := os.Getenv("REDIS_PATH")
	log.Println(redisPath)
	client, err := redis.New(redisPath)

	if err != nil {
		return errors.Wrap(err, "failed to get redis client")
	}

	defer client.Close()

	err = client.Get(target).Err()

	if err == redis.Nil {

		err = client.Set(target, 1, time.Hour*24).Err()
		if err != nil {
			return errors.Wrap(err, "failed to get redis client")
		}
	} else if err != nil {
		return errors.Wrapf(err, "failed to get %s", target)
	} else {
		currentNum, err := client.Incr(target).Result()
		if err != nil {
			return errors.Wrapf(err, "failed to incr %s", target)
		}
		log.Printf("currentNum is %d\n", currentNum)
	}
	return nil
}

// //ユーザーの削除
// func declValue(target string) (int64, error) {
// 	redisPath := os.Getenv("REDIS_PATH")
// 	client, err := redis.New(redisPath)

// 	if err != nil {
// 		return 0, errors.Wrap(err, "failed to get redis client")
// 	}

// 	defer client.Close()

// 	err = client.Get(target).Err()

// 	if err == redis.Nil {

// 		err = client.Set(target, 0, time.Hour*24).Err()
// 		if err != nil {
// 			return -1, errors.Wrap(err, "failed to get redis client")
// 		}
// 	} else if err != nil {
// 		return -1, errors.Wrapf(err, "failed to get %s", target)
// 	} else {
// 		currentNum, err := client.Decr(target).Result()
// 		if err != nil {
// 			return currentNum, errors.Wrapf(err, "failed to incr %s", target)
// 		}
// 		log.Printf("currentNum is %d\n", currentNum)
// 	}
// 	return 0, nil
// }
func declValue(target string) (int64, error) {
	redisPath := os.Getenv("REDIS_PATH")
	client, err := redis.New(redisPath)
	if err != nil {
		return 0, errors.Wrap(err, "failed to get redis client")
	}
	defer client.Close()
	currentNum, err := client.Decr(target).Result()
	if err != nil {
		return 0, errors.Wrap(err, "failed to decr CLIENT_NUM")
	}
	return currentNum, nil
}

func postMessage(message string) error {

	redisPath := os.Getenv("REDIS_PATH")
	client, err := redis.New(redisPath)
	if err != nil {
		return errors.Wrap(err, "failed to get redis client")
	}

	defer client.Close()

	err = client.Get(messageTarget).Err()

	if err == redis.Nil {
		err = client.Set(messageTarget, message, time.Hour*24).Err()
		if err != nil {
			return errors.Wrap(err, "failed to get redis client")
		}
	} else if err != nil {
		return errors.Wrapf(err, "failed to get %s", messageTarget)
	} else {
		err = client.Append(messageTarget, message).Err()
		if err != nil {
			return errors.Wrapf(err, "failed to incr %s", messageTarget)
		}
	}
	return nil
}

func getDoorState() (string, error) {

	redisPath := os.Getenv("REDIS_PATH")
	client, err := redis.New(redisPath)

	if err != nil {
		return "", errors.Wrap(err, "failed to get redis client")
	}

	defer client.Close()

	message, err := client.Get(doorTarget).Result()
	if err == redis.Nil {
		err = client.Set(doorTarget, "CLOSED", time.Hour*24).Err()
		if err != nil {
			return "ERROR!", errors.Wrap(err, "failed to set client")
		}
		return "CLOSED", nil
	} else if err != nil {
		return "ERROR!", errors.Wrap(err, "failed to connect")
	} else {
		return message, nil
	}

}
