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
