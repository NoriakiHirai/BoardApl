package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	_"time"
)

type BoardChaincode struct {
}

// スレッド数(キー："TotalThread"(固定))
type TotalThread struct {
	Counts	int `json:"counts"`
}

// スレッド情報(キー：連番(0,1,2,・・・))
type Thread struct {
	ThreadID	string `json:"threadId"`
	ThreadName	string `json:"threadName"`
	MsgNumber	string `json:"msgnumber"`	//メッセージの総数
}

// 投稿情報(キー：スレッド名⁺メッセージ)
type ContributionInfo struct {
	MsgNumber	string `json:"msgnumber"`
	UserID		string `json:"userID"`
	Message		string `json:"message"`
}

// スレッド情報の初期値を設定
func (cc *BoardChaincode) Init (stub *shim.ChaincodeStub, function string, args []string)([]byte, error) {
	var totalThread TotalThread
	var threadName string

	totalThread.Counts = 0
	defaultThreadName := [4]string{"News", "Economics", "Sports", "Culture"}

	// すでにスレッドが作成されていないか、確認する
	valAsbytes, _ := stub.GetState("0")
	if valAsbytes == nil {
		for i :=0; i < len(defaultThreadName); i++ {
			threadName = defaultThreadName[i]
			cc.addThread(stub, threadName)
		}
	} else if valAsbytes != nil {
		//初期スレッドが登録されていれば、何もしない
		return nil, nil
	}

	return nil, nil
}

func (cc *BoardChaincode) Invoke(stub *shim.ChaincodeStub, function string, args [] string) ([]byte, error) {

	//投稿処理実行
	if function == "contribution" {
		threadName := args[0]
		threadID := args[1]
		userID := args[4]
		msg := args[5]

		cc.contribution(stub, threadID, threadName, userID, msg)
		return nil, nil

	//スレッド追加
	} else if function == "AddThread" {
		threadName := args[0]
		cc.addThread(stub, threadName)
		return nil, nil
	}

	return nil, errors.New("Received unknown function")
}

// スレッド情報および投稿情報を参照
func (cc *BoardChaincode) Query (stub *shim.ChaincodeStub, function string, args []string)([]byte, error) {
	// スレッド一覧画面の更新
	if function == "GetThread" {
		return cc.getThread(stub)

	// 個別スレッド画面の更新
	} else if function == "GetContribution" {
		threadName := args[0]
		threadID := args[1]
		firstMsgNum := args[2]
		endMsgNum := args[3]

		// 投稿情報を取得
		return cc.getContribution(stub, threadID, threadName, firstMsgNum, endMsgNum)
	}

	return nil, errors.New("Received unknown function")
}

// スレッド追加
func (cc *BoardChaincode) addThread(stub *shim.ChaincodeStub, threadName string) ([]byte, error) {
	var threads Thread
	var tmpBytesSet [2][]byte

	// 登録されているスレッド総数を取得
	tmptotalThreadBytes, _ := stub.GetState("TotalThread")

	// 取得したバイト形式の情報をTotalThread型に変換
	totalThread := TotalThread{}
	json.Unmarshal(tmptotalThreadBytes, &totalThread)

	// スレッド情報を作成
	totalThread.Counts++
	threadID := strconv.Itoa(totalThread.Counts)
	threads = Thread{ThreadID: threadID, ThreadName: threadName, MsgNumber: "0"}

	// スレッド情報をワールドステートに追加
	// バイト形式に変換
	tmpBytesSet[0], _ = json.Marshal(totalThread)
	tmpBytesSet[1], _ = json.Marshal(threads)

	//strTotalcnt := strconv.Itoa(totalThread.Counts -1 )

	// ワールドステートに追加
	stub.PutState("TotalThread", tmpBytesSet[0])
	stub.PutState(threadID, tmpBytesSet[1])

	return nil, nil
}

// 投稿処理実行
func (cc *BoardChaincode) contribution(stub *shim.ChaincodeStub, threadID string, threadName string,
userID string, msg string) ([]byte, error) {
	var cntinfo ContributionInfo
	var threads Thread
	var tmpBytesSet [2][]byte

	// ワールドステートから最新のメッセージNoを取得する
	threadJson, _ := stub.GetState(threadID)

	//取得したJSON形式バイト型の情報をThread形式に変換
	threads = Thread{}
	json.Unmarshal(threadJson, &threads)

	msgNumber := threads.MsgNumber
	intmsgNum, _ := strconv.Atoi(msgNumber)

	// 1件も投稿がない場合
	if msgNumber == "" {
		intmsgNum = 1

	// 1件以上投稿がある場合
	} else if msgNumber != "" {
		// 最新の投稿Noに⁺1
		intmsgNum++
	}

	strmsgNum := strconv.Itoa(intmsgNum)
	cntIndex := threadName + strmsgNum

	// ワールドステートの更新
	// スレッド情報を作成
	threads = Thread{ThreadID: strmsgNum, ThreadName: threadName, MsgNumber: strmsgNum}

	// 投稿情報を作成
	cntinfo = ContributionInfo{MsgNumber: strmsgNum,UserID: userID, Message: msg}

	// バイト形式に変換
	tmpBytesSet[0], _ = json.Marshal(threads)
	tmpBytesSet[1], _ = json.Marshal(cntinfo)

	// ワールドステートを更新
	stub.PutState(threadID, tmpBytesSet[0])
	stub.PutState(cntIndex, tmpBytesSet[1])

	return nil, nil
}

// スレッド情報の取得
func (cc *BoardChaincode) getThread (stub *shim.ChaincodeStub)([]byte, error) {

	//スレッドの総数を取得
	tmptotalThreadBytes, _ := stub.GetState("TotalThread")

	// 取得したバイト形式の情報をTotalThread型に変換
	totalThread := TotalThread{}
	json.Unmarshal(tmptotalThreadBytes, &totalThread)

	//スレッド情報を格納するスライスを定義
	threads := make([]Thread, totalThread.Counts)

	for i :=1; i <= totalThread.Counts; i++ {
		//スレッド情報を取得する
		tmpthreadBytes, _ := stub.GetState(strconv.Itoa(i))
		threads[i-1] = Thread{}
		json.Unmarshal(tmpthreadBytes, &threads[i-1])
	}

	// json形式に変換
	return json.Marshal(threads)
}

// 投稿情報を取得する
func (cc *BoardChaincode) getContribution (stub *shim.ChaincodeStub, threadID string, threadName string,
firstMsgNum string, endMsgNum string)([]byte, error) {

	//投稿情報を格納するスライスを定義
	cntinfo := make([]ContributionInfo, 30)

	//スレッド情報を格納する構造体を定義
	var threads Thread

	msgFrom, _ := strconv.Atoi(firstMsgNum)
	msgTo, _ := strconv.Atoi(endMsgNum)

	// ワールドステートから最新のメッセージNoを取得する
	threadJson, _ := stub.GetState(threadID)
	//取得したJSON形式バイト型の情報をThread形式に変換
	threads = Thread{}
	json.Unmarshal(threadJson, &threads)

	msgNumber := threads.MsgNumber
	intmsgNumber, _ := strconv.Atoi(msgNumber)
	//fmt.Println(msgNumber)
	//fmt.Println(msgFrom)
	//fmt.Println(msgTo)
	//fmt.Println(intmsgNumber)

	// 1件も投稿がない場合、処理終了
	if msgNumber == "0" {
		//投稿がなければ、メッセージを表示する
		jsonResp := "[{No Contribution.}]"
		return json.Marshal(jsonResp)

	// 1件以上投稿がある場合
	} else if msgNumber != "0" {
		// if  取得開始位置 > 最新の投稿No then 取得終了位置 = 最新の投稿No
		if msgFrom > intmsgNumber {
			//取得終了位置に最新の投稿Noを設定
			msgTo = intmsgNumber
			// if 最新投稿No > 30 then 取得開始位置 = msgTo - 29
			// if 最新投稿No > 30 else 取得開始位置 = 0
			if intmsgNumber > 30 {
				msgFrom = msgTo - 30
			} else {
				msgFrom = 0
			}
		// if  取得開始位置 <= 最新の投稿No
		} else {
			// if 0 < 取得終了位置 <= 最新投稿No
			if msgTo < intmsgNumber {
				// 取得開始位置 = 取得開始位置 - 1
				msgFrom = msgFrom - 1

				// if 取得終了位置 = 0 (最新情報取得) then 取得終了位置 = 最新投稿No
				// else 取得終了位置 ⇒ そのまま
				if msgTo == 0 {
					msgTo = intmsgNumber

					// 取得開始位置 = 最新投稿No - 30
					msgFrom = msgTo - 30

					// 取得開始位置 < 0 then 取得開始位置 = 0
					if msgFrom < 0 {
						msgFrom = 0
					}
				}

			// if 取得終了位置 >= 最新投稿No
			} else {
				// 取得終了位置 = 最新投稿No
				msgTo = intmsgNumber

				// if 0 < 最新投稿No - 30 then 取得開始位置 = 最新投稿No - 29
				// if 最新投稿No - 30 <= 0 then 取得開始位置 = 0
				if intmsgNumber > 30 {
					msgFrom = intmsgNumber - 30
				} else {
					msgFrom = 0
				}
			}
		}

		// スライスのインデックスを定義
		index := 0

		// 任意の位置の投稿情報を取得する
		for i := msgFrom; i < msgTo; i++ {
			// ワールドステートのインデックスを生成
			cntIndex := threadName + strconv.Itoa(i+1)

			// ワールドステートから投稿情報を取得
			cntJson, _ := stub.GetState(cntIndex)

			//取得したJSON形式バイト型の情報をContributionInfo形式に変換
			cntinfo[index] = ContributionInfo{}
			json.Unmarshal(cntJson, &cntinfo[index])
			index = index + 1
		}

		// json形式に変換
		return json.Marshal(cntinfo)
	}

	return nil, errors.New("Unknown Error Occurs")
}

//Validating Peerに接続し、チェーンコードを実行
func main() {
	err := shim.Start(new(BoardChaincode))
	if err != nil {
		//err内容をそのまま出力
		fmt.Printf("Error starting chaincode: %s", err)
	}
}
