/**
 * Atomを読み込んでデータに変換する
 */
package okareader
import (
	"appengine"
	"appengine/urlfetch"
	"net/http"
	"encoding/xml"
)

/**
 * urlからatomファイルを受信して解析結果を返す
 * @function
 * @param c  コンテキスト
 * @param url atomファイルの場所
 */
func get(c appengine.Context, url string) *Atom {
	var client *http.Client
	var response *http.Response
	var err error
	var atom *Atom
	var encoded []byte
	
	// URLからatomを取得
	client = urlfetch.Client(c)
	response, err = client.Get(url)
	Check(c, err)
	
	// atomを受信
	atom = new(Atom)
	encoded = make([]byte, response.ContentLength)
	_, err = response.Body.Read(encoded)
	Check(c, err)

	// atomを解析
	err = xml.Unmarshal(encoded, atom)
	Check(c, err)
	
	return atom
}