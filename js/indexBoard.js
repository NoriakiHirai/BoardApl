var url_host="http://localhost:5000/";

//スレッド一覧更新
function RefreshThreadList() {
  //既存のスレッド一覧を削除
  $('table#ThreadTBL tbody *').remove();

  //スレッド一覧情報を取得
  getThreadList();
}

//スレッド追加
function AddThread() {
  var url = url_host + "chaincode";
  var NewThreadName = $('#InputThreadName').val();
  var JSONdata = createJSONdataForThread("invoke", "AddThread", [NewThreadName], 3);
  executeJsonRpc(url, JSONdata,
    function success(data) {
        console.log("AddThread Success");
    },
    function error(data) {
      console.log("AddThread Error");
    }
  );
  //入力フィールドを初期化
  $('#InputThreadName').val('');
}

//スレッド一覧作成
function getThreadList() {
  var url = url_host + "chaincode";
  var JSONdata = createJSONdataForThread("query", "GetThread", [], 5);
  executeJsonRpc(url, JSONdata,
    function success(data) {
      var DataList = JSON.parse(data.result.message);
      make1layerTable(DataList);
      console.log("ThreadList Refresh Success");
    },
    function error(data) {
      console.log("ThreadList Refresh Error");
    }
  );
}

//JSONメッセージ生成(スレッド関連)
function createJSONdataForThread(method, functionName, args, id) {
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
  name: ccId
  }
  // console.log(JSONdata);
  return JSONdata;
}

$(function(){
  $("#ThreadTBL tbody").click(function(e){
      var threadid = e.target.id;
      var threadname = e.target.innerText;
      window.sessionStorage.setItem(['ThreadId'], [threadid]);
      window.sessionStorage.setItem(['ThreadName'], [threadname]);
      window.location.href = "Thread.html";
  });
});
