/**
 * Atomを読み込んでデータに変換する
 */
package okareader
import (
	"appengine"
	"encoding/xml"
)

type Atom struct {

}

/**
 * データストアに保存できる形式に変換する
 * @methodOf AtomTemplate
 * @returns entries {[]Entry} 変換後のエントリ
 * @returns feed {Feed} 変換後のFeed
 */
func (this *Atom) encode(c appengine.Context, xmldata []byte) (*Feed, []*Entry) {
	type EntryTemplate struct {
		Id string `xml:"id"`
		Link struct {
			Href string `xml:"href,attr"`
		} `xml:"link"`
		Title string `xml:"title"`
		Owner string
	}
	type AtomTemplate struct {
		Id string `xml:"id"`
		Title string `xml:"title"`
		Link struct {
			Href string `xml:"href,attr"`
		} `xml:"link"`
		Entries []*EntryTemplate `xml:"entry"`
		Owner string
	}
	var feed *Feed
	var entryTemplate *EntryTemplate
	var entries []*Entry
	var entry *Entry
	var err error
	
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
		entry.Id = entryTemplate.Id
		entry.Link = entryTemplate.Link.Href
		entry.Title = entryTemplate.Title
		
		entries = append(entries, entry)
	}
	
	// Atomの変換
	feed.Id = atomTemplate.Link.Href
	feed.Title = atomTemplate.Title
	feed.Standard = "Atom"
	
	return feed, entries
}
