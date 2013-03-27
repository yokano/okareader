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
	Summary string `xml:"summary"`
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

type Entry struct {
	Id string
	Link string
	Summary string
	Title string
	Updated string
}

type Atom struct {
	Id string
	Title string
	Entries []string
}

/**
 * urlからatomファイルを受信して解析結果を返す
 * @function
 * @param c  コンテキスト
 * @param url atomファイルの場所
 */
func (this *AtomTemplate) Get(c appengine.Context, url string) (*Atom, []*Entry) {
	var client *http.Client
	var response *http.Response
	var err error
	var encoded []byte
	var atom *Atom
	var entries []*Entry
	
	// URLからatomを取得
	client = urlfetch.Client(c)
	response, err = client.Get(url)
	Check(c, err)
	
	// atomを受信
	encoded = make([]byte, response.ContentLength)
	_, err = response.Body.Read(encoded)
	Check(c, err)
	
	// atomを解析
	err = xml.Unmarshal(encoded, this)
	Check(c, err)
	
	// atomを変換
	atom, entries = this.encode()
	
	return atom, entries
}

/**
 * データストアに保存できる形式に変換する
 * @methodOf AtomTemplate
 * @returns entries {[]Entry} 変換後のエントリ
 * @returns atom {Atom} 変換後のAtom
 */
func (this *AtomTemplate) encode() (*Atom, []*Entry){
	var atom *Atom
	var entryTemplate *EntryTemplate
	var entries []*Entry
	var entry *Entry
	atom = new(Atom)
	atom.Entries = make([]string, 0)
	entries = make([]*Entry, 0)
	
	// エントリの変換
	for _, entryTemplate = range this.Entries {
		entry = new(Entry)
		entry.Id = entryTemplate.Id
		entry.Link = entryTemplate.Link.Href
		entry.Summary = entryTemplate.Summary
		entry.Title = entryTemplate.Title
		entry.Updated = entryTemplate.Updated
		
		entries = append(entries, entry)
		atom.Entries = append(atom.Entries, entry.Id)
	}
	
	// Atomの変換
	atom.Id = this.Id
	atom.Title = this.Title
	
	return atom, entries
}
