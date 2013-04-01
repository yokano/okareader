/**
 * エントリポイント
 * クライアントからのリクエストをコントローラへ投げる
 * @function
 */
package okareader

func init() {
	var controller *Controller
	controller = new(Controller)
	controller.handle()	
}
