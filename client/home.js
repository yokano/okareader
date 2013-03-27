/**
 * @file フォルダ画面のJavaScript
 */

$(function() {
	var contents = $('#contents');
	
	// フォルダを追加するボタン
	$('#add_folder_button').tap(function() {
		var name = $('#folder_name').val();
		
		$.ajax('/api/addfolder', {
			data: {
				folder_name: name,
				folder_key: $('#folder_key').val()
			},
			success: function() {
				contents.append($('<li><a href="#">' + name + '</a></li>')).listview('refresh');
				$('#add_folder').popup('close');
				$('#folder_name').val('');
			},
			error: function() {
				console.log('error');
			}
		});
	});
	
	// フィードを追加するボタン
	$('#add_feed_button').tap(function() {
		console.log('add feed');
		var url = $('#feed_url').val();
		
		$.ajax('/api/addfeed', {
			data: {
				url: url,
				folder_key: $('#folder_key').val()
			},
			success: function() {
				$('#add_feed').popup('close');
				$('#feed_url').val('');
			},
			error: function() {
				console.log('error');
			}
		});
	});
});

