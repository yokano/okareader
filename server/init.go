/**
 * init.go
 * クライアントからリクエストが来たら処理を振り分ける
 */

package okareader

import(
	"net/http"
	"appengine"
)

func init() {
	http.HandleFunc("/", home)
	http.HandleFunc("/atom", atom_test)
}

func atom_test(w http.ResponseWriter, r *http.Request) {
	var c appengine.Context
	
	c = appengine.NewContext(r)
	get(c, "http://feed.rssad.jp/rss/gigazine/rss_atom")
}