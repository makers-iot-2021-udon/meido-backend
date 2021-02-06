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
	Action  string `json:"action"`
	Message string `json:"message"`
	Name    string `json:"name"`
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

// type Message struct {
// 	Action   string   `json:"action"`
// 	Messages []string `json:"messages"`
// }

type Message struct {
	Action  string `json:"action"`
	Message string `json:"message"`
}

type CountMessage struct {
	Action string `json:"action"`
	Count  int64  `json:"count"`
}

type CertStatusMessage struct {
	Action string `json:"action"`
	Error  bool   `json:"error"`
	Status string `json:"status"`
	Count  int64  `json:"count"`
	Name   string `json:"name"`
}

const connectionTarget = "connections"
const messageTarget = "messages"
const doorTarget = "doorTarget"
const acceptTarget = "acceptTarget"
const deniedTarget = "deniedTarget"

var errorResponse = []byte(`{"action":"ERROR_MESSAGE","status":"NG","error": true}`)
var defaultMeidoStatus = []byte(`{"action":"MEIDO_STATUS","status":"FINE","error":false}`)

func handler(s []byte) ([]byte, bool) {
	var r Request
	if err := json.Unmarshal(s, &r); err != nil {
		return errorResponse, false
	}

	switch {
	case r.Action == "POST_DOOR":
		r, err := doorHandler(r.Message)
		if err != nil {
			return errorResponse, false
		}
		return r, true

	// case r.Action == "GET_DOOR":
	// 	r,err:=getDoorHandler()
	// 	if err != nil{
	// 		return errorResponse,
	// 	}

	// 	return r,true
	case r.Action == "POST_ACCEPT_USER":
		r, err := certUserHandler(acceptTarget, r.Name, "POST_ACCEPT_USER")
		if err != nil {
			return errorResponse, false
		}
		return r, true
	case r.Action == "POST_DENIED_USER":
		r, err := certUserHandler(deniedTarget, r.Name, "POST_DENIED_USER")
		if err != nil {
			return errorResponse, false
		}
		return r, true

	case r.Action == "ACCEPT_USER":
		r, err := countUpUserHandler(acceptTarget, "ACCEPT_USER")
		if err != nil {
			return errorResponse, false
		}
		return r, false

	case r.Action == "DENIED_USER":
		r, err := countUpUserHandler(deniedTarget, "DENIED_USER")
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
		r, err := countPeopleHandler(connectionTarget)
		if err != nil {
			return errorResponse, false
		}
		return r, true

	case r.Action == "POST_MESSAGE":
		r, err := messageHandler(r.Message)
		if err != nil {
			return errorResponse, false
		}
		return r, false

	// こいつ使わんでもよさげ
	// case r.Action == "MEIDO_FUN":
	// 	r, err := doorHandler()
	// 	if err != nil {
	// 		return errorResponse, false
	// 	}
	// 	return r, false
	case r.Action == "MEIDO_MESSAGE":
		r, err := connectHandler()
		if err != nil {
			return errorResponse, false
		}
		return r, false
	}
	return errorResponse, false
}

//Todo これはオウム返しを全体配信するだけ
func messageHandler(message string) ([]byte, error) {
	//DBに記録する
	err := saveMessage(message)
	if err != nil {
		return nil, err
	}
	r := Message{
		Message: message,
		Action:  "POST_MESSAGE",
	}
	b, err := json.Marshal(r)
	if err != nil {
		log.Println("cannot marshal struct: %v", err)
		return nil, err
	}
	return b, nil

}

func countUpUserHandler(target string, actionType string) ([]byte, error) {
	count, err := countUser(target)
	if err != nil {
		return nil, err
	}

	//メッセージを作成
	r := CountMessage{
		Action: actionType,
		Count:  count,
	}

	b, err := json.Marshal(r)
	if err != nil {
		log.Println("cannot marshal struct: %v", err)
		return nil, err
	}
	return b, nil
}

func certUserHandler(target string, name string, actionType string) ([]byte, error) {
	count, err := addCertUser(target, name)
	if err != nil {
		return nil, err
	}
	r := CertStatusMessage{
		Action: actionType,
		Status: "SUCCESS",
		Error:  false,
		Count:  count,
		Name:   name,
	}
	b, err := json.Marshal(r)
	if err != nil {
		log.Println("cannot marshal struct: %v", err)
		return nil, err
	}
	return b, nil
}

//単純な数え上げを許容
func countPeopleHandler(target string) ([]byte, error) {
	err := addValue(target)
	if err != nil {
		return nil, err
	}
	var _count int64 = 0
	_count, err = declValue(target)
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
func doorHandler(message string) ([]byte, error) {
	message, err := getDoorState(message)

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
	err := saveMessage(message)
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
		err = saveMessage("HELLO_WORLD")
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

func saveMessage(message string) error {
	redisPath := os.Getenv("REDIS_PATH")
	client, err := redis.New(redisPath)

	if err != nil {
		return errors.Wrap(err, "failed to get redis client")
	}

	defer client.Close()

	err = client.RPush(messageTarget, message).Err()
	if err != nil {
		return errors.Wrapf(err, "failed to get redis client")
	} else {
		err = client.Do("DEL", messageTarget).Err()
		if err != nil {
			log.Println(err, "failed to delete message from redis")
		}
		return nil
	}
	return nil
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

//ユーザーの追加と値返し

//絶対に増えない人
func countUser(target string) (int64, error) {
	redisPath := os.Getenv("REDIS_PATH")
	log.Println(redisPath)
	client, err := redis.New(redisPath)

	if err != nil {
		return -1, errors.Wrap(err, "failed to get redis client")
	}

	defer client.Close()
	err = client.Get(target).Err()

	if err != nil {
		return 0, errors.Wrapf(err, "failed to get %s", target)
	} else {
		err := client.SCard(target).Err()
		currentNum := client.SCard(target).Val()
		if err != nil {
			return -1, errors.Wrapf(err, "failed to count %s", target)
		} else {
			log.Printf("currentNum is %d\n", currentNum)
			return currentNum, nil
		}
	}
}

func addCertUser(target string, name string) (int64, error) {
	redisPath := os.Getenv("REDIS_PATH")
	log.Println(redisPath)
	client, err := redis.New(redisPath)

	if err != nil {
		return -1, errors.Wrap(err, "failed to get redis client")
	}

	defer client.Close()

	err = client.SRandMember(target).Err()

	if err == redis.Nil {
		// err = client.Set(target, name, time.Hour*24).Err()
		err = client.SAdd(target, name).Err()
		if err != nil {
			log.Println(err)
			return -1, errors.Wrap(err, "failed to get redis client")
		}
		err = client.Expire(target, 24*time.Hour).Err()
		if err != nil {
			log.Println(err)
			log.Println("Set Expired")
			return -1, errors.Wrap(err, "failed to get redis client")
		}
	} else if err != nil {
		log.Println(err)
		return -1, errors.Wrapf(err, "failed to get %s", target)
	} else {
		err := client.SAdd(target, name).Err()
		currentNum := client.SCard(target).Val()
		if err != nil {
			log.Println(err)
			return -1, errors.Wrapf(err, "failed to incr %s", target)
		}
		log.Printf("currentNum is %d\n", currentNum)
		return currentNum, nil
	}
	return 1, nil
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

// func postMessage(message string) error {

// 	redisPath := os.Getenv("REDIS_PATH")
// 	client, err := redis.New(redisPath)
// 	if err != nil {
// 		return errors.Wrap(err, "failed to get redis client")
// 	}

// 	defer client.Close()

// 	err = client.Get(messageTarget).Err()

// 	if err == redis.Nil {
// 		err = client.Set(messageTarget, message, time.Hour*24).Err()
// 		if err != nil {
// 			return errors.Wrap(err, "failed to get redis client")
// 		}
// 	} else if err != nil {
// 		return errors.Wrapf(err, "failed to get %s", messageTarget)
// 	} else {
// 		err = client.Append(messageTarget, message).Err()
// 		if err != nil {
// 			return errors.Wrapf(err, "failed to incr %s", messageTarget)
// 		}
// 	}
// 	return nil
// }

func getDoorState(message string) (string, error) {

	redisPath := os.Getenv("REDIS_PATH")
	client, err := redis.New(redisPath)

	if err != nil {
		return "", errors.Wrap(err, "failed to get redis client")
	}

	defer client.Close()

	err = client.Get(doorTarget).Err()
	if err == redis.Nil {
		err = client.Set(doorTarget, message, time.Hour*24).Err()
		if err != nil {
			return "ERROR!", errors.Wrap(err, "failed to set client")
		}
		return "CLOSED", nil
	} else if err != nil {
		return "ERROR!", errors.Wrap(err, "failed to connect")
	} else {
		err = client.Set(doorTarget, message, time.Hour*1).Err()
		//Todo オウム返ししているだけやん！
		return message, nil
	}

}
