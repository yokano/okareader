/**
 * 汎用関数
 */
package okareader
import (
	"appengine"
	"appengine/urlfetch"
	"net/http"
)

/**
 * エラーチェック
 * エラーがあればコンソールに出力する
 * @function
 * @param {appengine.Context} c コンテキスト
 * @param {error} err チェックするエラーオブジェクト
 */
func check(c appengine.Context, err error) {
	if err != nil {
		c.Errorf(err.Error())
	}
}

/**
 * スライスから指定された要素を削除して返す
 * 存在しなければ何もしない
 * 削除するのは最初に出現した１つのみ
 * @function
 * @param {[]string} s 対象のスライス
 * @param {string} target 削除する文字列
 * @returns {[]string} 削除済みのスライス
 */
func removeItem(s []string, target string) []string {
	var i int
	var str string
	var result []string
	
	result = make([]string, len(s))
	copy(result, s)
	for i, str = range s {
		if str == target {
			result = append(s[:i], s[i+1:]...)
			break
		}
	}
	
	return result
}

/**
 * 指定されたURLからXMLファイルを受信して返す
 * @param {appengine.Context} c コンテキスト
 * @param {string} url URL
 * @returns {[]byte} 受信したXMLデータ
 */
func getXML(c appengine.Context, url string) []byte {
	var client *http.Client
	var response *http.Response
	var err error
	var result []byte
	
	client = urlfetch.Client(c)
	response, err = client.Get(url)
	check(c, err)
	
	result = make([]byte, response.ContentLength)
	_, err = response.Body.Read(result)
	check(c, err)
	
	return result
}