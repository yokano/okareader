/**
 * 汎用関数
 */
package okareader
import (
	"appengine"
)

/**
 * エラーチェック
 * エラーがあればコンソールに出力する
 * @func
 * @param {appengine.Context} c コンテキスト
 * @param {error} err チェックするエラーオブジェクト
 */
func Check(c appengine.Context, err error) {
	if err != nil {
		c.Errorf(err.Error())
	}
}
