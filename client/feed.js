/**
 * フィード画面のJavaScript
 */
$('.feed_page').live('pageinit', function() {
	var feedKey = $(this).attr('key');

	// エントリをタップしたら既読化
	$(this).find('.entry').bind('tap', function() {
		var self = $(this);
		$.ajax('/api/read', {
			data: {
				link: self.attr('href'),
				feed_key: feedKey
			},
			error: function() {
				console.log('network error');
			},
			success: function() {
				self.parent().parent().parent().remove();
			}
		});
	});
	
	// 既読化ボタンをタップしたらすべて既読化
	$(this).find('#read_all').bind('tap', function() {
		if(window.confirm('すべてのエントリを既読化しますか？')) {
			var loading_div = $('<div class="loading"></div>').appendTo($('#contents'));
			$.ajax('/api/readall', {
				data: {
					key: feedKey
				},
				async: false,
				error: function() {
					console.log('network error');
				},
				success: function() {
					$('#entries').empty();
				},
				complete: function() {
					loading_div.remove();
				}
			});
		}
	});
	
	// 更新ボタンをタップしたらフィードを更新
	$(this).find('#reload').bind('tap', function() {
		var loading_div = $('<div class="loading"></div>').appendTo($('#contents'));
		$.ajax('/api/updatefeed', {
			data: {
				key: feedKey
			},
			dataType: 'json',
			async: false,
			success: function(data) {
				if(data.length == 0) {
					alert('新着はありませんでした');
					return;
				}
				var entries = $('#entries');
				for(var i = data.length - 1; i >= 0; i--) {
					var li = $('<li><a href="' + data[i].Link + '" class="entry" target="_blank">' + data[i].Title + '</a></li>');
					li.prependTo(entries)
				}
				entries.listview('refresh');
				alert(data.length + '件の新着を追加しました');
			},
			error: function() {
				console.log('error');
			},
			complete: function() {
				loading_div.remove();
			}
		});
	});
});
