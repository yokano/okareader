/**
 * フォルダ画面のJavaScript
 * ボタンを押した時のポップアップ操作や
 * サーバのAPIを呼び出しを行う
 */
$(document).on('pageinit', '.folder_page', function() {
	var contents = $(this).find('#contents');
	var addFolderButton = $(this).find('#add_folder_button');
	var folderName = $(this).find('#folder_name');
	var folderKey = $(this).attr('folder_key');
	var addFolder = $(this).find('#add_folder');
	var addFeedButton = $(this).find('#add_feed_button');
	var feedURL = $(this).find('#feed_url');
	var addFeed = $(this).find('#add_feed');
	var editButton = $(this).find('#edit');
	var feedName = $(this).find('#feed_name');
	var feedMenu = $(this).find('#feed_menu');
	var folderNewName = $(this).find('#folder_new_name');
	var folderMenu = $(this).find('#folder_menu');
	var editFeed = $(this).find('#edit_feed');
	var editMode = false;
	var editTarget = null;
	var busy = false;
	
	// フォルダを追加するボタン
	addFolderButton.on('tap', function() {
		var name = folderName.val();
		
		$.ajax('/api/addfolder', {
			data: {
				folder_name: name,
				folder_key: folderKey
			},
			dataType: 'json',
			success: function(data) {
				contents.append($('<li><div class="folder_icon"></div><a class="item" href="/folder?key=' + data.key + '" key="' + data.key + '" type="folder"><span class="title">' + name + '</span></a></li>')).listview('refresh');
				addFolder.popup('close');
				folderName.val('');
			},
			error: function() {
				console.log('error');
			}
		});
	});
	
	// フィードを追加するボタン
	addFeedButton.on('tap', function() {
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
					if(data.result == 'nothing_file') {
						alert('指定されたURLに配信用のファイルが見つかりませんでした。Atom, RSS2.0, RSS1.0 に対応したファイルの場所を指定してください。');
					} else if(data.result == 'duplicated') {
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
	editButton.on('tap', function() {
		if(editMode) {

			// 編集モード終了時の処理 //

			editMode = false;
			editTarget = null;
			$(this).find('.ui-btn-text').html('編集');
			$(this).find('.ui-icon').removeClass('ui-icon-check').addClass('ui-icon-edit');
			
			// リンクを有効化
			contents.find('a').off('tap');
			
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
				$(data).find('a').on('tap', function() {
					editTarget = $(this);

					if(editTarget.attr('type') == 'feed') {
						feedName.val($(this).find('.title').html());
						feedMenu.popup('open', {
							transition: 'pop',
							positionTo: 'window'
						});
					} else if(editTarget.attr('type') == 'folder') {
						folderNewName.val($(this).find('.title').html());
						folderMenu.popup('open', {
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
	$(this).find('#feed_name_button').on('tap', function() {
		var name = feedName.val();
		var key = editTarget.attr('key');
		
		if(name == '') {
			alert('名前を入力してください');
			return;
		}
		$.ajax('/api/renamefeed', {
			data: {
				name: name,
				key: key
			},
			success: function() {
				editTarget.find('.title').html(name);
				feedName.val('');
				editFeed.popup('close');
			},
			error: function() {
				console.log('error');
			}
		});
	});
	
	// フィード削除ボタン
	$(this).find('#remove_feed').on('tap', function() {
		var key = editTarget.attr('key');
		$.ajax('/api/removefeed', {
			data: {
				key: key
			},
			success: function() {
				editTarget.parent().parent().parent().remove();
				feedMenu.popup('close');
				contents.listview('refresh');
			},
			error: function() {
				console.log('error');
			}
		});
	});
	
	// フォルダ名変更ボタン
	$(this).find('#folder_name_button').on('tap', function() {
		var name = folderNewName.val();
		var key = editTarget.attr('key');
		
		if(name == '') {
			alert('名前を入力してください');
			return;
		}
		
		$.ajax('/api/renamefolder', {
			data: {
				key: key,
				name: name
			},
			success: function() {
				editTarget.find('.title').html(name);
				$('#edit_folder').popup('close');
				folderNewName.val('');
			},
			error: function() {
				console.log('error');
			}
		});
	});
	
	// フォルダ削除ボタン
	$(this).find('#remove_folder').on('tap', function() {
		var key = editTarget.attr('key');
		$.ajax('/api/removefolder', {
			data: {
				key: key
			},
			success: function() {
				editTarget.parent().parent().parent().remove();
				folderMenu.popup('close');
				contents.listview('refresh');
			},
			error: function() {
				console.log('error');
			}
		});
	});
	
	// フォルダの既読化ボタン
	$(this).find('#read').on('tap', function() {
		if(busy) {
			return;
		}
		busy = true;
		if(confirm('フォルダの中身をすべて既読化しますか？')) {
			$.ajax('/api/readfolder', {
				data: {
					key: folderKey
				},
				success: function() {
					$('.ui-li-count').remove();
					alert('フォルダを既読化しました');
				},
				error: function() {
					console.log('error');
				},
				complete: function() {
					busy = false;
				}
			});
		}
	});
	
	// フォルダの更新ボタン
	$(this).find('#reload').on('tap', function() {
		if(busy) {
			return;
		}
		busy = true;
		
		if(!confirm('フォルダ内のフィードを更新しますか？')) {
			return;
		}
		
		$.ajax('/api/updatefolder', {
			data: {
				key: folderKey
			},
			dataType: 'json',
			success: function(data) {
				var updated = false;
				for(var key in data) {
					var count = $('[key=' + key + ']').find('.ui-li-count');
					if(count.html() < data[key]) {
						updated = true;
					}
					count.html(data[key]);
				}
				if(updated) {
					alert('新着エントリを追加しました');
				} else {
					alert('新着はありませんでした');
				}
			},
			error: function() {
				console.log('error');
			},
			complete: function() {
				busy = false;
			}
		});
	});
	
	// XMLアップロード時のチェック
	$(this).find('#uploadxml').on('tap', function() {
		var filename = $('#xmlfile').val();
		if(!filename.match(/\.xml$/)) {
			alert('xmlファイルが選択されていません');
			return false
		}
	});
});

