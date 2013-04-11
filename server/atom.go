/**
 * Atomを読み込んでデータに変換する
 */
package okareader
import (
	"appengine"
	"encoding/xml"
)

/**
 * Atomファイルの操作
 * @class
 */
type Atom struct {
}

/**
 * データストアに保存できる形式に変換する
 * @methodOf Atom
 * @param {appengine.Context} c コンテキスト
 * @param {[]byte} xmldata
 * @returns {Feed} feed フィードリスト
 * @returns {[]Entry} entries エントリリスト
 */
func (this *Atom) encode(c appengine.Context, xmldata []byte) (*Feed, []*Entry) {
	type EntryTemplate struct {
		Link struct {
			Href string `xml:"href,attr"`
		} `xml:"link"`
		Title string `xml:"title"`
		Owner string
	}
	type FeedLink struct {
		Rel string `xml:"rel,attr"`
		Href string `xml:"href,attr"`
	}
	type AtomTemplate struct {
		Id string `xml:"id"`
		Title string `xml:"title"`
		Link []FeedLink `xml:"link"`
		Entries []*EntryTemplate `xml:"entry"`
		Owner string
	}
	var feed *Feed
	var entryTemplate *EntryTemplate
	var entries []*Entry
	var entry *Entry
	var err error
	var link FeedLink
	
	feed = new(Feed)
	feed.Entries = make([]string, 0)
	entries = make([]*Entry, 0)
	
	// atomを解析
	var atomTemplate = new(AtomTemplate)
	err = xml.Unmarshal(xmldata, atomTemplate)
	check(c, err)
	
	// エントリの変換
	for _, entryTemplate = range atomTemplate.Entries {
		entry = new(Entry)
		entry.Link = entryTemplate.Link.Href
		entry.Title = entryTemplate.Title
		
		entries = append(entries, entry)
	}
	
	// Atomの変換
	for _, link = range atomTemplate.Link {
		if link.Rel == "alternate" {
			feed.SiteURL = link.Href
		} else if link.Rel == "self" {
			feed.URL = link.Href
		}
	}
	feed.Title = atomTemplate.Title
	feed.Standard = "Atom"
	
	return feed, entries
}
