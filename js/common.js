//JSON-RPC実行
function executeJsonRpc(url_exec, JSONdata, success, error) {
  $.ajax({
    type: 'post',
    url: url_exec,
    data: JSON.stringify(JSONdata),
    contentType: 'application/JSON;',
    dataType: 'JSON',
    scriptCharset: 'utf-8',
    success: function(data) {
        success(data)
    },
    error: function(data) {
      error(data);
    }
  });
  // console.log(JSON.stringify(JSONdata));
}

// テーブル作成(1段構造)
function make1layerTable(DataList) {
  for (var i = 0; i < DataList.length; i++) {
    // スレッド情報の取得
    var thread = DataList[i];
    //### HTML編集 table行の追加、編集 ここから ###
    var temp = "";
    temp += "<tr>";
    temp += "<td width=\"50\">" + (i + 1) + "</td>";
    temp += "<td width=\"200\"><a href=\"Thread.html\" onClick=\"setId(" + thread.threadId + ")\">" + thread.threadName + "</a></td>";
    temp += "<td width=\"50\">" + thread.msgnumber + "</td>";
    temp += "</tr>";

    $("#ThreadTBL tbody").append(temp);
    //### HTML編集 table行の追加、編集 ここまで ###
  }
}

// テーブル作成(2段構造)
function make2layerTable(DataList) {
  for (var i = 0; i < DataList.length; i++) {
    // スレッド情報の取得
    var contribution = DataList[i];
    if (contribution.userID == "") {
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
  }
  return DataList[i - 1].msgnumber;
}
