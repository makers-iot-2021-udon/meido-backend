package main

import (
	"encoding/json"
	"log"
	"main/persistence/redis"
	"os"
	"time"

	"github.com/pkg/errors"
)

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

// Count Accept user and denied.
func countUpUserHandler(target string, actionType string) ([]byte, error) {
	count, err := countUser(target)
	if err != nil {
		log.Println(err)
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
