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

type EntryTemplate struct {
	Id string `xml:"id"`
	Link struct {
		Href string `xml:"href,attr"`
	} `xml:"link"`
	Title string `xml:"title"`
	Updated string `xml:"updated"`
	Owner string
}

type AtomTemplate struct {
	Id string `xml:"id"`
	Title string `xml:"title"`
	Entries []*EntryTemplate `xml:"entry"`
	Owner string
}

/**
 * urlからatomファイルを受信して解析結果を返す
 * @function
 * @param c  コンテキスト
 * @param url atomファイルの場所
 */
func (this *AtomTemplate) get(c appengine.Context, url string) (*Feed, []*Entry) {
	var client *http.Client
	var response *http.Response
	var err error
	var encoded []byte
	var feed *Feed
	var entries []*Entry
	
	// URLからatomを取得
	client = urlfetch.Client(c)
	response, err = client.Get(url)
	check(c, err)
	
	// atomを受信
	encoded = make([]byte, response.ContentLength)
	_, err = response.Body.Read(encoded)
	check(c, err)
	
	// atomを解析
	err = xml.Unmarshal(encoded, this)
	check(c, err)
	
	// atomを変換
	feed, entries = this.encode()
	
	return feed, entries
}

/**
 * データストアに保存できる形式に変換する
 * @methodOf AtomTemplate
 * @returns entries {[]Entry} 変換後のエントリ
 * @returns feed {Feed} 変換後のFeed
 */
func (this *AtomTemplate) encode() (*Feed, []*Entry){
	var feed *Feed
	var entryTemplate *EntryTemplate
	var entries []*Entry
	var entry *Entry
	feed = new(Feed)
	feed.Entries = make([]string, 0)
	entries = make([]*Entry, 0)
	
	// エントリの変換
	for _, entryTemplate = range this.Entries {
		entry = new(Entry)
		entry.Id = entryTemplate.Id
		entry.Link = entryTemplate.Link.Href
		entry.Title = entryTemplate.Title
		entry.Updated = entryTemplate.Updated
		
		entries = append(entries, entry)
	}
	
	// Atomの変換
	feed.Id = this.Id
	feed.Title = this.Title
	
	return feed, entries
}
