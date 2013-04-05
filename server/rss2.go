/**
 * RSS2.0の読み込み
 * xmlデータを受け取ってFeedとEntryリストを返す
 */
package okareader

import(
	"appengine"
	"encoding/xml"
)

type RSS2 struct {

}

/**
 * xmlをFeedオブジェクトに変換する
 * @methodOf RSS2
 * @param {appengine.Context} c コンテキスト
 * @param {[]byte} xmldata 変換するXMLデータ
 * @returns {*Feed} 変換結果のフィード
 * @returns {[]*Entry} 変換結果のエントリ
 */
func (this *RSS2) encode(c appengine.Context, xmldata []byte) (*Feed, []*Entry) {
	type Item struct {
		Title string `xml:"title"`
		Link string `xml:"link"`
		Description string `xml:"description"`
		Date string `xml:"date"`
	}
	type Channel struct {
		Title string `xml:"channel>title"`
		Link string `xml:"channel>link"`
		Date string `xml:"channel>date"`
		Item []*Item `xml:"channel>item"`
	}
	var feed *Feed
	var entries []*Entry
	var channel *Channel
	var err error
	var item *Item
	var i int
	
	channel = new(Channel)
	err = xml.Unmarshal(xmldata, channel)
	check(c, err)
	
	feed = new(Feed)
	feed.URL = channel.Link
	feed.Title = channel.Title
	feed.Standard = "RSS2.0"
	
	entries = make([]*Entry, len(channel.Item))
	for i, item = range channel.Item {
		entries[i] = new(Entry)
		entries[i].Id = item.Link
		entries[i].Title = item.Title
		entries[i].Link = item.Link
	}
	
	return feed, entries
}
