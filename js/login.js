var url_host="http://localhost:5000/";
var user_name;

//ログイン
function login(){
  var url = url_host + "registrar";
  var user_name = $("#userName").val();
  window.sessionStorage.setItem(['USER_NAME'], [user_name]);
  var password = $("#password").val();
  var JSONdata = {
    "enrollId": user_name,
    "enrollSecret": password
  };
  executeJsonRpc(url, JSONdata,
    function success(data) {
      //ログイン成功時
      console.log("login success!");
    },
    function error(data) {
      //ログインエラー
      console.log("login error");
    }
  );
}
//デプロイ
function deploy() {
  var url = url_host + "chaincode";
  var JSONdata = createJSONdataForDeploy("deploy", "init", [], 1);
  executeJsonRpc(url, JSONdata,
    function success(data) {
      // デプロイ成功時
      var ccId = data.result.message;
      console.log("deploy success!");
      window.sessionStorage.setItem(['CCID'], [ccId]);
      window.alert("デプロイ中です。30秒ほど待ってから次画面ボタンを押下してください。")
    },
    function error(data) {
      //デプロイエラー
      console.log("deploy error");
    }
  );
}

//JSONメッセージ生成
function createJSONdataForDeploy(method, functionName, args, id) {
  var ccId = window.sessionStorage.getItem(['CCID']);
  var user_name = window.sessionStorage.getItem(['USER_NAME']);
  var JSONdata = {
    jsonrpc: "2.0",
    method: method,
    params: {
      type: 1,
      ctorMsg: {
        function: functionName,
        args: args
      },
      secureContext: user_name,
    },
    id: id
  };
  //チェーンコードIDを設定
  JSONdata.params["chaincodeID"] = {
    path: "github.com/hyperledger/fabric/examples/chaincode/go/chaincode_board"
  }
  // console.log(JSONdata);
  return JSONdata;
}
