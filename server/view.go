/**
 * Controllerの命令に従いページを描画する
 * ここからデータストアへ直接アクセスしてはいけない
 * データが必要な場合は必ず DAO(model.go) に頼むこと
 * htmlファイルへのアクセスはここからだけ行うこと
 */
package okareader
import (
	"html/template"
	"appengine"
	"appengine/user"
	"net/http"
	text "text/template"
)

/**
 * ページの表示関係を行うオブジェクト
 * @class
 */
type View struct {
}

/**
 * フォルダの中身を一覧表示
 * @methodOf View
 * @param {appengine.Context} c コンテキスト
 * @param {string} key エンコード済みのフォルダのキー
 * @param {http.ResponseWriter} w HTMLの出力先
 */
func (this *View) showFolder(c appengine.Context, key string, w http.ResponseWriter) {
	type ListItem struct {
		Key string
		Item interface{}
		ItemType string
		Count int
	}
	var contents map[string]interface{}
	var err error
	var t *template.Template
	var children []*ListItem
	var dao *DAO
	var folder *Folder
	var i int
	
	dao = new(DAO)
	
	contents = make(map[string]interface{}, 0)
	contents["LogoutURL"], err = user.LogoutURL(c, "/")
	check(c, err)

	contents["FolderKey"] = key
	
	folder = new(Folder)
	folder = dao.getFolder(c, key)
	contents["Title"] = folder.Title
	contents["Parent"] = folder.Parent
	
	children = make([]*ListItem, len(folder.Children))
	for i, key = range folder.Children {
		children[i] = new(ListItem)
		children[i].Key = key
		children[i].ItemType, children[i].Item = dao.getItem(c, key)
	}
	contents["Children"] = children
	
	t, err = template.ParseFiles("server/html/folder.html")
	check(c, err)
	
	t.Execute(w, contents)
}

/**
 * フィードのエントリを一覧表示
 * @methodOf View
 * @param {appengine.Context} c コンテキスト
 * @param {string} feedKey 表示するフィードのキー
 * @param {http.ResponseWriter} w HTMLの出力先
 */
func (this *View) showFeed(c appengine.Context, feedKey string, w http.ResponseWriter) {
	var dao *DAO
	var entries []*Entry
	var t *template.Template
	var err error
	var contents map[string]interface{}
	var feed *Feed
	
	dao = new(DAO)
	feed = dao.getFeed(c, feedKey)
	entries = dao.getEntries(c, feedKey)
	
	t, err = template.ParseFiles("server/html/feed.html")
	check(c, err)
	
	contents = make(map[string]interface{})
	contents["Title"] = feed.Title
	contents["Entries"] = entries
	contents["Parent"] = feed.Parent
	contents["FeedKey"] = feedKey
	contents["SiteURL"] = feed.SiteURL
	contents["LogoutURL"], err = user.LogoutURL(c, "/")
	check(c, err)
	
	t.Execute(w, contents)
}

/**
 * ログインを促す画面を表示
 * @methodOf View
 * @param {appengine.Context} c コンテキスト
 * @param {http.RespnoseWriter} w 応答先
 */
func (this *View) showLogin(c appengine.Context, w http.ResponseWriter) {
	var contents map[string]interface{}
	var err error
	var t *template.Template
	
	t, err = template.ParseFiles("server/html/login.html")
	check(c, err)
	
	contents = make(map[string]interface{}, 0)
	contents["LoginURL"], err = user.LoginURL(c, "/")
	check(c, err)
	
	t.Execute(w, contents)
}

/**
 * XMLファイルインポート前の確認画面
 * @param {appengine.Context} c コンテキスト
 * @param {http.ResponseWriter} w 応答先
 * @param {[]*Node} tree 追加するフォルダ・フィードツリー
 * @param {string} folderKey 追加先のフォルダのキー
 * @param {string} token 進捗をクライアントへ送信するためのチャネルで使用するトークン
 */
func (this *View) confirmImporting(c appengine.Context, w http.ResponseWriter, tree []*Node, folderKey string) {
	var t *text.Template
	var err error
	var contents map[string]string
	var node *Node
	var child *Node
	var html string
	
	t, err = text.ParseFiles("server/html/import.html")
	check(c, err)
	
	html = ""
	for _, node = range tree {
		if node.kind == "folder" {
			html = join(html, `<li>`, node.title, `</li>`)
			html = join(html, `<ul>`)
			for _, child = range node.children {
				html = join(html, `<li>`, child.title, `</li>`)
			}
			html = join(html, `</ul>`)
		} else if node.kind == "feed" {
			html = join(html, `<li>`, node.title, `</li>`)
		}
	}
	
	contents = make(map[string]string, 1)
	contents["tree"] = html
	contents["folder_key"] = folderKey
	contents["LogoutURL"], err = user.LogoutURL(c, "/")
	check(c, err)
	t.Execute(w, contents)
}
