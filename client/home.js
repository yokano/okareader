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
	var editButton = $(this).find('#edit');
	var editMode = false;
	
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
	
	// 編集ボタン
	editButton.bind('tap', function() {
		if(editMode) {
			editMode = false;
			$(this).find('.ui-btn-text').html('編集');
			$(this).find('.ui-icon').removeClass('ui-icon-check').addClass('ui-icon-edit');
			
			// リンクを有効化
			contents.find('a').unbind('tap');
			
			// アイコンを戻す
			contents.find('.ui-icon').addClass('ui-icon-arrow-r').removeClass('ui-icon-gear');
			
			// メッセージ非表示
			$('#edit_message').remove();
			
		} else {
			editMode = true;
			$(this).find('.ui-btn-text').html('完了');
			$(this).find('.ui-icon').removeClass('ui-icon-edit').addClass('ui-icon-check');
	
			// リンクを無効化
			$.each(contents.children(), function(i, data) {
				$(data).find('a').bind('tap', function() {
					return false;
				});
			});
			
			// アイコンを変更
			$.each(contents.children(), function(i, data) {
				$(data).find('.ui-icon').removeClass('ui-icon-arrow-r').addClass('ui-icon-gear');
			});
			
			// メッセージ表示
			contents.prepend($('<li id="edit_message" data-role="list-divider">編集したいタイトルをタップ</li>'));
			
			contents.listview('refresh');
		}
	});
});
