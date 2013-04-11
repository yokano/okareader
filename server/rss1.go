/**
 * RSS1.0の読み込み
 */
package okareader
import (
	"appengine"
	"encoding/xml"
)

/**
 * RSS1.0
 * @class
 */
type RSS1 struct {

}

/**
 * RSS1.0のXMLをFeedに変換する
 * @methodOf RSS1
 * @param {appengine.Context} c コンテキスト
 * @param {[]byte} xmldata XMLのバイト配列
 * @returns {*Feed} 変換したフェード
 * @returns {[]*Entries} 変換したエントリ
 */
func (this *RSS1) encode(c appengine.Context, xmldata []byte) (*Feed, []*Entry) {
	type Item struct {
		Title string `xml:"title"`
		Link string `xml:"link"`
		Date string `xml:"date"`
	}
	type Channel struct {
		Title string `xml:"title"`
		Link string `xml:"link"`
		Date string `xml:"date"`
		About string `xml:"about,attr"`
	}
	type RDF struct {
		Channel *Channel `xml:"channel"`
		Item []*Item `xml:"item"`
	}
	var feed *Feed
	var entries []*Entry
	var rdf *RDF
	var err error
	var i int
	var item *Item
	
	rdf = new(RDF)
	err = xml.Unmarshal(xmldata, rdf)
	check(c, err)
	
	feed = new(Feed)
	feed.URL = rdf.Channel.About
	feed.Title = rdf.Channel.Title
	feed.Standard = "RSS1.0"
	
	entries = make([]*Entry, len(rdf.Item))
	for i, item = range rdf.Item {
		entries[i] = new(Entry)
		entries[i].Title = item.Title
		entries[i].Link = item.Link
	}
	
	return feed, entries
}