/**
 * インポート画面
 */
$(document).on('pageinit', '.confirm_page', function() {
	var key = $(this).attr('folder_key');
	var busy = false;
	
	// ボタン連打の防止
	$(this).find('#import').on('tap', function() {
		if(busy) {
			return;
		}
		busy = true;
		$.ajax('/api/importxml', {
			data: {
				key: key
			},
			error: function() {
				alert('インポートに失敗しました');
			},
			success: function() {
				alert('インポートが完了しました');
				$.mobile.changePage('/folder?key=' + key);
			},
			complete: function() {
				busy = false;
			}
		});
	});
});