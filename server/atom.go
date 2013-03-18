/**
 * Atomを読み込んでデータに変換する
 */
package okareader
import (
	"appengine"
	"appengine/urlfetch"
	"net/http"
	"log"
	"encoding/xml"
)

func get(c appengine.Context, url string) {
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
		Entries []Entry `xml:"entry"`
	}

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
	
	// 結果表示
	for i := 0; i < len(atom.Entries); i++ {
		log.Printf("%s\n", atom.Entries[i].Link.Href)
	}
}