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

// スレッド数
type TotalThread struct {
	Counts	int `json:"counts"`
}

// スレッド名
type ThreadName struct {
	ThreadName	string `json:"threadName"`
}

// スレッド情報
type Thread struct {
	ThreadName	string `json:"threadName"`
	MsgNumber	string `json:"msgnumber"`	//メッセージの総数
}

// 投稿情報
type ContributionInfo struct {
	ThreadName	string `json:"threadName"`
	MsgNumber	string `json:"msgnumber"`
	UserID		string `json:"userID"`
	Message		string `json:"message"`
}

// スレッド情報の初期値を設定
func (cc *BoardChaincode) Init (stub *shim.ChaincodeStub, function string, args []string)([]byte, error) {
	var totalThread TotalThread
	var threadNames ThreadName
	var threads Thread

	totalThread.Counts = 0
	defaultThreadName := [4]string{"News", "Economics", "Sports", "Culture"}

	// すでにスレッドが作成されていないか、確認する
	valAsbytes, _ := stub.GetState(defaultThreadName[0])
	if valAsbytes == nil {
		for i :=0; i < len(defaultThreadName); i++ {
			// スレッド情報を作成
			totalThread.Counts++
			threadNames = ThreadName{ThreadName: defaultThreadName[i]}
			threads = Thread{ThreadName: defaultThreadName[i], MsgNumber: "0"}

			// スレッド情報をワールドステートに追加
			// バイト形式に変換
			tmptotalThreadBytes, _ := json.Marshal(totalThread)
			tmpthreadNamesBytes, _ := json.Marshal(threadNames)
			tmpthreadsBytes, _ := json.Marshal(threads)

			// ワールドステートに追加
			stub.PutState("TotalThread", tmptotalThreadBytes)
			stub.PutState(strconv.Itoa(i), tmpthreadNamesBytes)
			stub.PutState(defaultThreadName[i], tmpthreadsBytes)
		}
	} else if valAsbytes != nil {
		//初期スレッドが登録されていれば、何もしない
		return nil, nil
	}

	return nil, nil
}

func (cc *BoardChaincode) Invoke(stub *shim.ChaincodeStub, function string, args [] string) ([]byte, error) {
	//function名でハンドリング
	//投稿処理実行
	if function == "contribution" {
		return cc.contribution(stub, args)
	}
	//スレッド追加
	if function == "AddThread" {
		return cc.AddThread(stub, args)
	}

	return nil, errors.New("Received unknown function")
}

// スレッド情報および投稿情報を参照
func (cc *BoardChaincode) Query (stub *shim.ChaincodeStub, function string, args []string)([]byte, error) {

	// function名でハンドリング
	// スレッド一覧画面の更新
	if function == "RefreshThreadList" {
		return cc.getThread(stub, args)

	// 投稿情報表示画面の更新(最新情報)
	} else if function == "RefreshBoardLatest" {
		// 投稿情報を取得
		return cc.getLatestContribution(stub, args)

	// 投稿情報表示画面の更新(表示メッセージNo指定)
	} else if function == "RefreshBordSelect" {
		// 投稿情報を取得
		return cc.getContribution(stub, args)
	}

	return nil, errors.New("Received unknown function")
}

// スレッド追加
func (cc *BoardChaincode) AddThread(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	var threadNames ThreadName
	var threads Thread

	// ワールドステートのインデックスを作成
	threadName := args[0]

	// すでにスレッドが作成されていないか、確認する
	valAsbytes, _ := stub.GetState(threadName)
	if valAsbytes == nil {
		// 登録されているスレッド総数を取得
		tmptotalThreadBytes, _ := stub.GetState("TotalThread")

		// 取得したバイト形式の情報をTotalThread型に変換
		totalThread := TotalThread{}
		json.Unmarshal(tmptotalThreadBytes, &totalThread)

		// スレッド情報を作成
		totalThread.Counts++
		threadNames = ThreadName{ThreadName: threadName}
		threads = Thread{ThreadName: threadName, MsgNumber: "0"}

		// スレッド情報をワールドステートに追加
		// バイト形式に変換
		tmptotalThreadBytes, _ = json.Marshal(totalThread)
		tmpthreadNamesBytes, _ := json.Marshal(threadNames)
		tmpthreadsBytes, _ := json.Marshal(threads)

		strTotalcnt := strconv.Itoa(totalThread.Counts -1 )

		// ワールドステートに追加
		stub.PutState("TotalThread", tmptotalThreadBytes)
		stub.PutState(strTotalcnt, tmpthreadNamesBytes)
		stub.PutState(threadName, tmpthreadsBytes)

	} else if valAsbytes != nil {
		//スレッドが登録されていれば、メッセージを表示する
		jsonResp := "[{Thread Exists.}]"
		return json.Marshal(jsonResp)
	}

	return nil, nil
}

// 投稿処理実行
func (cc *BoardChaincode) contribution(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	//投稿情報を格納する構造体を定義
	var cntinfo ContributionInfo

	// スレッド情報
	var threads Thread

	var cntIndex string

	// ワールドステートのインデックスを作成
	threadName := args[0]

	// ワールドステートから最新のメッセージNoを取得する
	threadJson, _ := stub.GetState(threadName)

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
	cntIndex = threadName + strmsgNum

	// ワールドステートへの登録情報を設定する
	userId := args[3]
	msg := args[4]

	// １．投稿情報の更新
	// 投稿情報を作成
	cntinfo = ContributionInfo{ThreadName: threadName, MsgNumber: strmsgNum, UserID: userId, Message: msg}

	// バイト形式に変換
	tmpcntinfoBytes, _ := json.Marshal(cntinfo)

	// 投稿情報をワールドステートに追加
	stub.PutState(cntIndex, tmpcntinfoBytes)

	// ２．スレッド一覧の更新
	// スレッド情報を作成
	threads = Thread{ThreadName: threadName, MsgNumber: strmsgNum}

	// ワールドステートのスレッド情報を更新
	// バイト形式に変換
	tmpthreadsBytes, _ := json.Marshal(threads)

	// ワールドステートを更新
	stub.PutState(threadName, tmpthreadsBytes)

	return nil, nil
}

// スレッド情報の取得
func (cc *BoardChaincode) getThread (stub *shim.ChaincodeStub, args []string)([]byte, error) {

	//スレッドの総数を取得
	tmptotalThreadBytes, _ := stub.GetState("TotalThread")

	// 取得したバイト形式の情報をTotalThread型に変換
	totalThread := TotalThread{}
	json.Unmarshal(tmptotalThreadBytes, &totalThread)

	//スレッド情報を格納するスライスを定義
	threads := make([]Thread, totalThread.Counts)
	threadNames := make([]ThreadName, totalThread.Counts)

	for i :=0; i < totalThread.Counts; i++ {
		//スレッド名を取得
		threadNameBytes, _ := stub.GetState(strconv.Itoa(i))

		//取得したJSON形式バイト型の情報をThreadName形式に変換
		threadNames[i] = ThreadName{}
		json.Unmarshal(threadNameBytes, &threadNames[i])

		//スレッド情報を取得する
		tmpthreadBytes, _ := stub.GetState(threadNames[i].ThreadName)
		threads[i] = Thread{}
		json.Unmarshal(tmpthreadBytes, &threads[i])
	}

	// json形式に変換
	return json.Marshal(threads)
}

// getLatestContribution
func (cc *BoardChaincode) getLatestContribution (stub *shim.ChaincodeStub, args []string)([]byte, error) {
	//投稿情報を格納するスライスを定義
	cntinfo := make([]ContributionInfo, 30)

	//スレッド情報を格納する構造体を定義
	var threads Thread

	var firstmsgNum int

	// 配列に格納されたスレッド名を受け取る
	threadName := args[0]

	// ワールドステートから最新のメッセージNoを取得する
	threadJson, _ := stub.GetState(threadName)

	//取得したJSON形式バイト型の情報をThread形式に変換
	threads = Thread{}
	json.Unmarshal(threadJson, &threads)

	msgNumber := threads.MsgNumber
	intmsgNum, _ := strconv.Atoi(msgNumber)

	// 1件も投稿がない場合、処理終了
	if msgNumber == "0" {
		//投稿がなければ、メッセージを表示する
		jsonResp := "[{No Contribution.}]"
		return json.Marshal(jsonResp)

	// 1件以上投稿がある場合
	} else if msgNumber != "0" {
		// 投稿されている書き込みが30件以下の場合、1件目から最新の投稿までを読み込む
		if intmsgNum < 31 {
			firstmsgNum = 1
		} else {
			firstmsgNum = intmsgNum - 29
		}


		// スライスのインデックスを定義
		index := 0
		// 最新の投稿情報を取得する
		for i := firstmsgNum; i <= intmsgNum; i++ {
			// ワールドステートのインデックスを生成
			cntIndex := threadName + strconv.Itoa(i)

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

// 任意の行数を取得する
func (cc *BoardChaincode) getContribution (stub *shim.ChaincodeStub, args []string)([]byte, error) {
	//投稿情報を格納するスライスを定義
	cntinfo := make([]ContributionInfo, 30)

	//スレッド情報を格納する構造体を定義
	var threads Thread

	// 配列に格納されたスレッド名を受け取る
	threadName := args[0]

	// 取得する投稿情報の開始位置と終了位置を設定する
	tmpNum1 := args[1]
	tmpNum2 := args[2]
	firstMsgNum, _ := strconv.Atoi(tmpNum1)
	endMsgNum, _ := strconv.Atoi(tmpNum2)

	// ワールドステートから最新のメッセージNoを取得する
	threadJson, _ := stub.GetState(threadName)
	//取得したJSON形式バイト型の情報をThread形式に変換
	threads = Thread{}
	json.Unmarshal(threadJson, &threads)

	msgNumber := threads.MsgNumber
	intmsgNumber, _ := strconv.Atoi(msgNumber)

	// 1件も投稿がない場合、処理終了
	if msgNumber == "0" {
		//投稿がなければ、メッセージを表示する
		jsonResp := "[{No Contribution.}]"
		return json.Marshal(jsonResp)

	// 1件以上投稿がある場合
	} else if msgNumber != "0" {
		// if  取得開始位置 > 最新の投稿No then 取得終了位置 = 最新の投稿No
		if firstMsgNum > intmsgNumber {
			//取得終了位置に最新の投稿Noを設定
			endMsgNum = intmsgNumber
			// if 最新投稿No > 30 then 取得開始位置 = endMsgNum - 29
			// if 最新投稿No > 30 else 取得開始位置 = 1
			if intmsgNumber > 30 {
				firstMsgNum = endMsgNum - 29
			} else {
				firstMsgNum = 1
			}
		// if  取得開始位置 <= 最新の投稿No
		} else {
			// if 0 < 取得終了位置 <= 最新投稿No
			if endMsgNum < intmsgNumber {
				// 取得開始位置,取得終了位置は引数のまま

			// if 取得終了位置 >= 最新投稿No
			} else {
				// 取得終了位置 = 最新投稿No
				endMsgNum = intmsgNumber

				// if 0 < 最新投稿No - 30 then 取得開始位置 = 最新投稿No - 29
				// if 最新投稿No - 30 <= 0 then 取得開始位置 = 1
				if intmsgNumber > 30 {
					firstMsgNum = intmsgNumber - 29
				} else {
					firstMsgNum = 1
				}
			}
		}

		// スライスのインデックスを定義
		index := 0

		// 任意の位置の投稿情報を取得する
		for i := firstMsgNum; i <= endMsgNum; i++ {
			// ワールドステートのインデックスを生成
			cntIndex := threadName + strconv.Itoa(i)

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
