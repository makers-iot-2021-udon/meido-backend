package main

import (
	"log"
)

//現在のステータスを作るやつ
type CurrentStatusMessage struct {
	Action          string `json:"action"`
	ConnectingCount int64  `json:"connect_count"`
	AcceptUserCount int64  `json:"accept_count"`
	DeniedUserCount int64  `json:"denied_count"`
	ErrorLogCount   int64  `json:"error_count"`
	SystemStatus    string `json:"system_status"`
	AuthStatus      string `json:"auth_status`
}

var errorStatusMessageResponse = CurrentStatusMessage{
	Action:          "NOTIFY_CURRENT_STATUS",
	ConnectingCount: 0,
	AcceptUserCount: 0,
	DeniedUserCount: 0,
	ErrorLogCount:   0,
	SystemStatus:    "Available",
	AuthStatus:      "not working",
}

func currentStatus() CurrentStatusMessage {

	//現在の接続ユーザーのカウント
	err := addValue(CLIENT_NUM)
	connectingCount, err := declValue(CLIENT_NUM)
	if err != nil {
		log.Println(err)
		return errorStatusMessageResponse
	}

	acceptCount, err := countUser(acceptTarget)

	if err != nil {
		log.Println(err)
		return errorStatusMessageResponse

	}

	deniedCount, err := countUser(deniedTarget)
	if err != nil {
		log.Println(err)
		return errorStatusMessageResponse
	}

	const errorLogCount int64 = 0

	const systemStatus = "FINE"
	const authStatus = "Not working"

	//メッセージを作成
	r := CurrentStatusMessage{
		Action:          "NOTIFY_CURRENT_STATUS",
		ConnectingCount: connectingCount,
		AcceptUserCount: acceptCount,
		DeniedUserCount: deniedCount,
		ErrorLogCount:   errorLogCount,
		SystemStatus:    systemStatus,
		AuthStatus:      authStatus,
	}

	// b, err := json.Marshal(r)
	// if err != nil {
	// 	log.Println("cannot marshal struct: %v", err)
	// 	return errorResponse
	// }
	return r
}
