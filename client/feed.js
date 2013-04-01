/**
 * フィード画面のJavaScript
 */
$('.feed_page').live('pageinit', function() {
	var feedKey = $(this).attr('key');
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
});
