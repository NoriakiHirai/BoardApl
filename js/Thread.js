var url_host="http://localhost:5000/";

//初期処理
function init() {
  var threadName = getThreadName();
  $('.ThreadName').text(threadName);

  //掲示板更新
  RefreshBoardLatest();
}

//掲示板更新(最新)
function RefreshBoardLatest() {
  var funcname = "RefreshBoardLatest";
  //投稿情報を読み込む開始位置と終了位置を設定
  //最新の投稿を取得するので、今回は開始位置と終了位置は使用しないため、0を設定する
  var startmsgnum = "0";
  var endmsgnum = "0";

  RefreshBoard(funcname, startmsgnum, endmsgnum);
}

//前の30件を表示させる
function Pre30Contribution() {
  //投稿情報を読み込む開始位置と終了位置を設定
  var tmpstartmsgnum = $(".StartMsgNum").text() - 31;
  if (tmpstartmsgnum < 1) {
    tmpstartmsgnum = 1;
    tmpendmsgnum = 30;
  } else {
    tmpendmsgnum = tmpstartmsgnum + 29;
  }
  //型変換
  var startmsgnum = tmpstartmsgnum + '';
  var endmsgnum = tmpendmsgnum + '';

  var funcname = "RefreshBordSelect";

  RefreshBoard(funcname, startmsgnum, endmsgnum);
}

//次の30件を表示させる
function Rear30Contribution() {
  //投稿情報を読み込む開始位置と終了位置を設定
  var tmpstartmsgnum = $(".StartMsgNum").text() + 30;
  var tmpendmsgnum = tmpstartmsgnum + 30;

  //型変換
  var startmsgnum = tmpstartmsgnum + '';
  var endmsgnum = tmpendmsgnum + '';

  var funcname = "RefreshBordSelect";

  RefreshBoard(funcname, startmsgnum, endmsgnum);
}

//掲示板更新処理実行
function RefreshBoard(funcname, startmsgnum, endmsgnum) {
  //既存の投稿を削除
  $('#BoardTable').empty();

  var threadName = getThreadName();
  getContribution(funcname, threadName,  startmsgnum, endmsgnum);
}

//スレッド名取得
function getThreadName() {
  //スレッド一覧ページから渡されたスレッド名をタイトルに設定する
  var threadName = window.sessionStorage.getItem(['ThreadName']);
  return threadName;
}

// 投稿実行
function Contribution() {
  var url = url_host + "chaincode";
  //スレッド名を取得
  var threadName = $(".ThreadName").text();
  //ユーザー名を設定
  var user_name = window.sessionStorage.getItem(['USER_NAME']);
  //メッセージを設定
  var message = $('#InputArea').val();
  //改行コードを置換
  message = message.replace("\n", "<br>")
  var JSONdata = createJSONdataForBoardApp("invoke", "contribution", threadName, "0", "0", user_name, message, 3);
  executeJsonRpc(url, JSONdata,
    function success(data) {
        console.log("Contribution Success");
    },
    function error(data) {
      console.log("Contribution Error");
    }
  );
  //投稿フォームを初期化
  $('#InputArea').val('');
}
//投稿情報取得
function getContribution(funcname, threadName, startmsgnum, endmsgnum) {
  var url = url_host + "chaincode";
  //ユーザー名を設定
  var user_name = window.sessionStorage.getItem(['USER_NAME']);
  var JSONdata = createJSONdataForBoardApp("query", funcname, threadName, startmsgnum, endmsgnum, user_name, "", 5);
  executeJsonRpc(url, JSONdata,
    function success(data) {
      contributionList = JSON.parse(data.result.message);
      var getstartmsgnum = 0;
      var getendmsgnum = 0;
      if (contributionList.toString() != "[{No Contribution.}]") {
        //投稿情報の開始番号を設定
        getstartmsgnum = contributionList[0].msgnumber;

        for (var i = 0; i < contributionList.length; i++) {
          // スレッド情報の取得
          var contribution = contributionList[i];
          if (contribution.threadName == "") {
            break;
          }
          //### HTML編集 table行の追加、編集 ここから ###
          var temp = "";
          temp += "<table>";
          temp += "<thead>";
          temp += "<tr>";
          temp += "<th align=\"left\" width=\"30\">" + contribution.msgnumber + "</th>";
          temp += "<th align=\"left\">" + contribution.userID + "</th>";
          temp += "</tr>";
          temp += "</thead>";
          temp += "</table>";
          temp += "<table width=\"750\" style=\"table-layout:fixed;margin-left:10px;\">";
          temp += "<tr>";
          temp += "<td style=\"word-wrap:break-word;\" align=\"left\">" + contribution.message + "</td>";
          temp += "</tr><tr></tr>";
          temp += "</table>";

          $("#BoardTable").append(temp);
          //### HTML編集 table行の追加、編集 ここまで ###
          getendmsgnum = contribution.msgnumber;
        }
      } else if (contributionList.toString() == "[{No Contribution.}]") {
        window.alert("まだ1件も投稿されていません。\n投稿をお願いします。");
      }
      $(".StartMsgNum").text(getstartmsgnum);
      $(".EndMsgNum").text(getendmsgnum);
      console.log("contributionList Refresh Success");
    },
    function error(data) {
      console.log("contributionList Refresh Error");
    }
  );
}

//JSONメッセージ生成
function createJSONdataForBoardApp(method, functionName, threadName, startmsgnum, endmsgnum, user_name, message, id) {
  var ccId = window.sessionStorage.getItem(['CCID']);
  var user_name = window.sessionStorage.getItem(['USER_NAME']);
  var JSONdata = {
    jsonrpc: "2.0",
    method: method,
    params: {
      type: 1,
      ctorMsg: {
        function: functionName,
        args: [
          threadName,
          startmsgnum,
          endmsgnum,
          user_name,
          message
        ]
      },
      secureContext: user_name,
    },
    id: id
    };
    //チェーンコードIDを設定
    if (functionName == "init") {
      JSONdata.params["chaincodeID"] = {
        path: "github.com/hyperledger/fabric/examples/chaincode/go/chaincode_board"
      };
    } else {
      JSONdata.params["chaincodeID"] = {
        name: ccId
      };
    }
  return JSONdata;
}
