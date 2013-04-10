/**
 * フォルダ画面のJavaScript
 * ボタンを押した時のポップアップ操作や
 * サーバのAPIを呼び出しを行う
 */
$('.folder_page').live('pageinit', function() {
	var contents = $(this).find('#contents');
	var addFolderButton = $(this).find('#add_folder_button');
	var folderName = $(this).find('#folder_name');
	var folderKey = $(this).attr('folder_key');
	var addFolder = $(this).find('#add_folder');
	var addFeedButton = $(this).find('#add_feed_button');
	var feedURL = $(this).find('#feed_url');
	var addFeed = $(this).find('#add_feed');
	var editButton = $(this).find('#edit');
	var editMode = false;
	var editTarget = null;
	
	// フォルダを追加するボタン
	addFolderButton.bind('tap', function() {
		var name = folderName.val();
		
		$.ajax('/api/addfolder', {
			data: {
				folder_name: name,
				folder_key: folderKey
			},
			dataType: 'json',
			success: function(data) {
				contents.append($('<li><div class="folder_icon"></div><a class="item" href="/folder?key=' + data.key + '" key="' + data.key + '" type="folder">' + name + '</a></li>')).listview('refresh');
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
		var url = feedURL.val();
		
		if(!url.match(/http:\/\/.+/)) {
			alert('HTTPのURLを入力してください');
		} else {
			$.ajax('/api/addfeed', {
				data: {
					url: url,
					folder_key: folderKey
				},
				dataType: 'json',
				success: function(data) {
					if(data.duplicated) {
						alert('既に登録済みのフィードです')
					} else {
						contents.append($('<li><div class="feed_icon"></div><a class="item" href="/feed?key=' + data.key + '"  key="' + data.key + '" type="feed"><span class="title">' + data.name + '</span><span class="ui-li-count">' + data.count + '</span></a></li>')).listview('refresh');
					}
					addFeed.popup('close');
					feedURL.val('');
				},
				error: function() {
					console.log('error');
				}
			});
		}
		
		feedURL.val('');
	});
	
	// 編集ボタン
	editButton.bind('tap', function() {
		if(editMode) {

			// 編集モード終了時の処理 //

			editMode = false;
			editTarget = null;
			$(this).find('.ui-btn-text').html('編集');
			$(this).find('.ui-icon').removeClass('ui-icon-check').addClass('ui-icon-edit');
			
			// リンクを有効化
			contents.find('a').unbind('tap');
			
			// アイコンを戻す
			contents.find('.ui-icon').addClass('ui-icon-arrow-r').removeClass('ui-icon-gear');
			
			// メッセージ非表示
			$('#edit_message').remove();
			
		} else {
		
			// 編集モード開始時の処理 //
		
			editMode = true;
			$(this).find('.ui-btn-text').html('完了');
			$(this).find('.ui-icon').removeClass('ui-icon-edit').addClass('ui-icon-check');
			
			// リンクを無効化と編集ポップアップの表示
			$.each(contents.children(), function(i, data) {
				$(data).find('a').bind('tap', function() {
					editTarget = $(this);

					if(editTarget.attr('type') == 'feed') {
						$('#feed_name').val($(this).find('.title').html());
						$('#feed_menu').popup('open', {
							transition: 'pop',
							positionTo: 'window'
						});
					} else if(editTarget.attr('type') == 'folder') {
						$('#folder_new_name').val($(this).find('.title').html());
						$('#folder_menu').popup('open', {
							transition: 'pop',
							positionTo: 'window'
						});
					}
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
	
	// フィード名変更ボタン
	$('#feed_name_button').bind('tap', function() {
		var name = $('#feed_name').val();
		var key = editTarget.attr('key');
		$.ajax('/api/renamefeed', {
			data: {
				name: name,
				key: key
			},
			success: function() {
				editTarget.find('.title').html(name);
				$('#feed_name').val('');
				$('#edit_feed').popup('close');
			},
			error: function() {
				console.log('error');
			}
		});
	});
	
	// フィード削除ボタン
	$('#remove_feed').bind('tap', function() {
		var key = editTarget.attr('key');
		$.ajax('/api/removefeed', {
			data: {
				key: key
			},
			success: function() {
				editTarget.parent().parent().parent().remove();
				$('#feed_menu').popup('close');
			},
			error: function() {
				console.log('error');
			}
		});
	});
	
	// フォルダ名変更ボタン
	$('#folder_name_button').bind('tap', function() {
		var name = $('#folder_new_name').val();
		var key = editTarget.attr('key');
		
		$.ajax('/api/renamefolder', {
			data: {
				key: key,
				name: name
			},
			success: function() {
				editTarget.find('.title').html(name);
				$('#edit_folder').popup('close');
				$('#folder_new_name').val('');
			},
			error: function() {
				console.log('error');
			}
		});
	});
	
	// フォルダ削除ボタン
	$('#remove_folder').bind('tap', function() {
		var key = editTarget.attr('key');
		$.ajax('/api/removefolder', {
			data: {
				key: key
			},
			success: function() {
				editTarget.parent().parent().parent().remove();
				$('#folder_menu').popup('close');
			},
			error: function() {
				console.log('error');
			}
		});
	});
	
	// フォルダの既読化ボタン
	$('#read').bind('tap', function() {
		var loading_div = $('<div class="loading"></div>').appendTo(contents);
		if(confirm('フォルダの中身をすべて既読化しますか？')) {
			$.ajax('/api/readfolder', {
				data: {
					key: folderKey
				},
				success: function() {
					$('.ui-li-count').remove();
					loading_div.remove();
				},
				error: function() {
					console.log('error');
				}
			});
		}
	});
	
	// フォルダの更新ボタン
	$('#reload').bind('tap', function() {
		var loading_div = $('<div class="loading"></div>').appendTo(contents);
		$.ajax('/api/updatefolder', {
			data: {
				key: folderKey
			},
			dataType: 'json',
			success: function(data) {
				for(var key in data) {
					$('[key=' + key + ']').find('.ui-li-count').html(data[key]);
				}
				loading_div.remove();
			},
			error: function() {
				console.log('error');
			}
		});
	});
	
	// XMLアップロード時のチェック
	$('#uploadxml').bind('tap', function() {
		var filename = $('#xmlfile').val();
		if(!filename.match(/\.xml$/)) {
			alert('xmlファイルが選択されていません');
			return false
		}
	});
});
