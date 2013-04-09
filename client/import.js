/**
 * XMLファイルのインポート画面のスクリプト
 * サーバと通信しながら進捗表示を行う
 */
$('.confirm_page').live('pageinit', function() {
	var folderKey = $(this).attr('key');
	var i = 0
	
	var token = $(this).attr("token");
	var channel = new goog.appengine.Channel(token)
	var socket = channel.open();
	socket.onopen = function() {
		console.log('open');
	};
	socket.onmessage = function(data) {
		console.log(data);
	};
	socket.onerror = function() {
		console.log('socket error');
	};
	socket.onclose = function() {
		console.log('socket close');
	}
	
	$('#import_xml').tap(function() {
		$.ajax('/api/importxml', {
			data: {
				key: folderKey
			},
			error: function() {
				console.log('error');
			},
			success: function() {
				console.log('success');
			}
		});
	});
});