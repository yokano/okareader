/**
 * フォルダ画面のJavaScript
 * ボタンを押した時のポップアップ操作や
 * サーバのAPIを呼び出しを行う
 */
$('.folder_page').live('pageinit', function() {
	var contents = $(this).find('#contents');
	var addFolderButton = $(this).find('#add_folder_button');
	var folderName = $(this).find('#folder_name');
	var folderKey = $(this).find('#folder_key');
	var addFolder = $(this).find('#add_folder');
	var addFeedButton = $(this).find('#add_feed_button');
	var feedURL = $(this).find('#feed_url');
	var addFeed = $(this).find('#add_feed');
	
	// フォルダを追加するボタン
	addFolderButton.bind('tap', function() {
		var name = folderName.val();
		
		$.ajax('/api/addfolder', {
			data: {
				folder_name: name,
				folder_key: folderKey.val()
			},
			dataType: 'json',
			success: function(data) {
				contents.append($('<li><a href="/folder?key=' + data.key + '">' + name + '</a></li>')).listview('refresh');
				addFolder.popup('close');
				folderName.val('');
			},
			error: function() {
				console.log('error');
			}
		});
	});
	
	// フィードを追加するボタン
	addFeedButton.bind('tap', function() {
		console.log('add feed');
		var url = feedURL.val();
		
		$.ajax('/api/addfeed', {
			data: {
				url: url,
				folder_key: folderKey.val()
			},
			dataType: 'json',
			success: function(data) {
				if(data.duplicated) {
					alert('既に登録済みのフィードです')
				} else {
					contents.append($('<li><a href="/feed?key=' + data.key + '">' + data.name + '</a></li>')).listview('refresh');
				}
				addFeed.popup('close');
				feedURL.val('');
			},
			error: function() {
				console.log('error');
			}
		});
	});
});
