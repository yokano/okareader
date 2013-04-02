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
});
