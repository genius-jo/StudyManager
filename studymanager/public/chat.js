$(function () {

  //이벤트 소스 지원 하는지 확인
  if (!window.EventSource) {
      alert("not support EventSource")
      return
  }

  var $chatlog = $('#chat-log')
  var $chatmsg = $('#chat-msg')

  //sumbit이 발생했을때 해당 function 실행
  $('#input-form').on('submit', function (e) {
      
    //post메세지를 보낸다
    $.post('/messages', {
        msg: $chatmsg.val(),
    });
    
    //메세지 보내고, 비우고 포커스 맞추기
    $chatmsg.val("");
    $chatmsg.focus();
    
    return false; //다른 페이지로 넘어가지 않기 위해
  });
  
  var addMessage = function (data) {
    var text = "";
    
    text = '<strong>' + data.name + ' : </strong>'
    text += data.msg;
    
    $chatlog.prepend('<div><span>' + text + '</span></div>');
  };

  //이벤트 소스를 요청할 경로를 지정해서 이벤트 소스를 연다
  var es = new EventSource('/stream')

  es.onmessage = function (e) {
    var msg = JSON.parse(e.data);
    addMessage(msg);
  }

})