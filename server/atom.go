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

type Entry struct {
	Id string `xml:"id"`
	Content string `xml:"content"`
	Link struct {
		Href string `xml:"href,attr"`
	} `xml:"link"`
	Summary string `xml:"summary"`
	Title string `xml:"title"`
	Updated string `xml:"updated"`
}
type Atom struct {
	ID string `xml:"id"`
	Title string `xml:"title"`
	Entries []Entry `xml:"entry"`
}


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
	atom.ID = "none"
	encoded = make([]byte, response.ContentLength)
	_, err = response.Body.Read(encoded)
	Check(c, err)

	// atomを解析
	err = xml.Unmarshal(encoded, atom)
	Check(c, err)
	
	return atom
}