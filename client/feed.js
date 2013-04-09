/**
 * フィード画面のJavaScript
 */
$('.feed_page').live('pageinit', function() {
	var feedKey = $(this).attr('key');

	// エントリをタップしたら既読化
	$('.entry').bind('tap', function() {
		var self = $(this);
		$.ajax('/api/read', {
			data: {
				id: self.attr('id'),
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
	$('#read_all').bind('tap', function() {
		if(window.confirm('すべてのエントリを既読化しますか？')) {
			$.ajax('/api/readall', {
				data: {
					key: feedKey
				},
				error: function() {
					console.log('network error');
				},
				success: function() {
					$('#entries').empty();
				}
			});
		}
	});
	
	// 更新ボタンをタップしたらフィードを更新
	$('#reload').bind('tap', function() {
		var loading_div = $('<div class="loading"></div>').appendTo(contents);
		$.ajax('/api/updatefeed', {
			data: {
				key: feedKey
			},
			dataType: 'json',
			success: function(data) {
				loading_div.remove();
				var entries = $('#entries');
				for(var i = data.length - 1; i >= 0; i++) {
					var li = $('<li><a href="' + data[i].Link + '" id=' + data[i].Link + ' class="entry" target="_blank">' + data[i].Title + '</a></li>');
					li.prependTo(entries)
				}
				entries.listview('refresh');
			},
			error: function() {
				console.log('error');
			}
		});
	});
});
