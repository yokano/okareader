/**
 * @file フォルダ画面のJavaScript
 */

$(function() {
	var contents = $('#contents');
	
	// フォルダを追加するボタン
	$('#add_folder_button').tap(function() {
		var name = $('#folder_name').val();
		
		$.ajax('/addfolder', {
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
});

