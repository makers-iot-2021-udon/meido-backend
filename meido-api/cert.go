package main

import (
	"encoding/json"
	"log"
	"main/persistence/redis"
	"os"
	"time"

	"github.com/pkg/errors"
)

func certUserHandler(target string, name string, uid string, certMessage string, actionType string) ([]byte, error) {
	count, err := addCertUser(target, uid)
	if err != nil {
		return nil, err
	}
	r := CertStatusMessage{
		Action: actionType,
		Status: certMessage,
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
		//err := client.SAdd(target, name).Err()
		currentNum, err := client.SCard(target).Result()
		if err != nil {
			log.Println(err)
			return -1, errors.Wrapf(err, "failed to incr %s", target)
		}
		log.Printf("currentNum is %d\n", currentNum)
		return currentNum, nil
	}
	return 1, nil
}

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

	if err == redis.Nil {
		//テーブルが作成されていない間はユーザーが存在しないので超法規的に0を返す
		return 0, nil
	} else if err != nil {
		return 0, errors.Wrapf(err, "failed to get %s", target)
	} else {
		// err := client.SCard(target).Err()
		currentNum, err := client.SCard(target).Result()
		if err != nil {
			return -1, errors.Wrapf(err, "failed to count %s", target)
		} else {
			log.Printf("currentNum is %d\n", currentNum)
			return currentNum, nil
		}
	}
}
