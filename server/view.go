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
//	"log"
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
 * @param {http.ResponseWriter} w HTMLの出力先
 */
func (this *View) ShowFolder(c appengine.Context, key string, w http.ResponseWriter) {
	type ListItem struct {
		Key string
		Item interface{}
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
	contents["FolderKey"] = key
	Check(c, err)
	
	folder = new(Folder)
	folder = dao.GetFolder(c, key)
	contents["Title"] = folder.Title
	
	children = make([]*ListItem, len(folder.Children))
	for i, key = range folder.Children {
		children[i] = new(ListItem)
		children[i].Key = key
		children[i].Item = dao.GetItem(c, key)
	}
	contents["Children"] = children
	
	t, err = template.ParseFiles("server/html/folder.html")
	Check(c, err)
	
	t.Execute(w, contents)
}

/**
 * ログインを促す画面を表示
 * @methodOf View
 */
func (this *View) ShowLogin(c appengine.Context, w http.ResponseWriter) {
	var content map[string]interface{}
	var err error
	var t *template.Template
	
	t, err = template.ParseFiles("server/html/login.html")
	Check(c, err)
	
	content = make(map[string]interface{}, 0)
	content["LoginURL"], err = user.LoginURL(c, "/")
	Check(c, err)
	
	t.Execute(w, content)
}