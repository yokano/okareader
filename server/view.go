/**
 * Controllerの命令に従いページを描画する
 * ここからデータストアへ直接アクセスしてはいけない
 * データが必要な場合は必ず DAO(model.go) に頼むこと
 * htmlファイルへのアクセスはここからだけ行なってよい
 */
package okareader
import (
	"html/template"
	"appengine"
	"appengine/user"
	"net/http"
)

/**
 * ページの表示関係を行うオブジェクト
 */
type View struct {

}

/**
 * フォルダの中身を一覧表示
 * @methodOf View
 * @param {appengine.Context} c コンテキスト
 * @param {string} key エンコード済みのフォルダのキー
 * @param {string} from 遷移前のフォルダのキー
 * @param {http.ResponseWriter} w HTMLの出力先
 */
func (this *View) showFolder(c appengine.Context, key string, from string, w http.ResponseWriter) {
	type ListItem struct {
		Key string
		Item interface{}
		ItemType string
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
	contents["From"] = from
	
	folder = new(Folder)
	folder = dao.getFolder(c, key)
	contents["Title"] = folder.Title
	
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
 * @param {string} from 遷移前のフォルダのキー
 * @param {http.ResponseWriter} w HTMLの出力先
 */
func (this *View) showFeed(c appengine.Context, feedKey string, from string, w http.ResponseWriter) {
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
	contents["From"] = from
	contents["LogoutURL"], err = user.LogoutURL(c, "/")
	check(c, err)
	
	t.Execute(w, contents)
}

/**
 * ログインを促す画面を表示
 * @methodOf View
 */
func (this *View) showLogin(c appengine.Context, w http.ResponseWriter) {
	var content map[string]interface{}
	var err error
	var t *template.Template
	
	t, err = template.ParseFiles("server/html/login.html")
	check(c, err)
	
	content = make(map[string]interface{}, 0)
	content["LoginURL"], err = user.LoginURL(c, "/")
	check(c, err)
	
	t.Execute(w, content)
}