/**
 * ブラウザやAjaxのリクエストを適切な処理へ振り分ける
 * 直接データストアへアクセスしたり画面を描画したりしてはいけない
 * データストアへのアクセスはModelに,画面の描画はViewに頼むこと.
 */

package okareader

import(
	"appengine"
	"appengine/user"
	"net/http"
	"fmt"
)

/**
 * リクエストURLによる処理の振り分け
 * /api/*** はAjaxによるAPIへのアクセスであり画面の描画は不用
 * それ以外はページ遷移を表し画面を描画する
 * @function
 */
func init() {
	http.HandleFunc("/", home)
	http.HandleFunc("/folder", folder)
	http.HandleFunc("/feed", feed)
	http.HandleFunc("/test", atom_test)
	http.HandleFunc("/api/addfolder", addFolder)
	http.HandleFunc("/api/addfeed", addFeed)
}

/**
 * http://okareader.appspot.com/ へアクセスした時の処理
 * ログインしていなければログインさせる
 * ログインしていればルートフォルダを表示
 * @function
 * @param {http.ResponseWriter} 応答先
 * @param {*http.Request} リクエスト
 */
func home(w http.ResponseWriter, r *http.Request) {
	var c appengine.Context
	var u *user.User
	var root *Folder
	var dao *DAO
	var view *View
	var key string

	c = appengine.NewContext(r)
	u = user.Current(c)
	dao = new(DAO)
	view = new(View)
	
	if u == nil {
		view.ShowLogin(c, w)
	} else {
		key, root = dao.GetRootFolder(c, u)
		if root.Type == "" {
			key = dao.RegisterFolder(c, u, "root", true, "")
		}
		view.ShowFolder(c, key, w)
	}
}

/**
 * http://okareader.appspot.com/folder へアクセスしたらフォルダを表示
 * フォルダをデータストアから取得するためのキーをGETで受け取る
 * @function
 * @param {http.ResponseWriter} w 応答先
 * @param {*http.Request} r リクエスト
 * @param {HTTP GET} key エンコード済みのフォルダキー
 */
func folder(w http.ResponseWriter, r *http.Request) {
	var c appengine.Context
	var view *View
	var encodedKey string
	
	c = appengine.NewContext(r)
	encodedKey = r.FormValue("key")

	view = new(View)
	view.ShowFolder(c, encodedKey, w)
}

/**
 * http://okareader.appspot.com/feed へアクセスしたらフィードを表示
 * フィードのキーはGETで渡される
 * @function
 * @param {http.ResponseWriter} w 応答先
 * @param {*http.Request} r リクエスト
 * @param {HTTP GET} key エンコード済みのフィードキー
 */
func feed(w http.ResponseWriter, r *http.Request) {
	var c appengine.Context
	var view *View
	var feedKey string
	
	c = appengine.NewContext(r)
	feedKey = r.FormValue("key")
	
	view = new(View)
	view.ShowFeed(c, feedKey, w)
}

/**
 * API:フォルダの新規追加
 * @function
 * @param {http.ResponseWriter} w 応答先
 * @param {*http.Request} r HTTPリクエスト
 * @returns {AJAX:JSON} 追加したフォルダのキーを含むJSON
 */
func addFolder(w http.ResponseWriter, r *http.Request) {
	var c appengine.Context
	var u *user.User
	var dao *DAO
	var title string
	var encodedParentKey string
	var resultKey string
	
	c = appengine.NewContext(r)
	u = user.Current(c)
	dao = new(DAO)
	
	title = r.FormValue("folder_name")
	encodedParentKey = r.FormValue("folder_key")
	
	resultKey = dao.RegisterFolder(c, u, title, false, encodedParentKey)
	fmt.Fprintf(w, `{"key":"%s"}`, resultKey)
}

/**
 * API:フィードの登録
 * @function
 * @param {http.ResponseWriter} w 応答先
 * @param {*http.Request} r リクエスト
 */
func addFeed(w http.ResponseWriter, r *http.Request) {
	var url string
	var folderKey string
	var dao *DAO
	var c appengine.Context
	var atom *Atom
	var atomTemplate *AtomTemplate
	var entries []*Entry
	var feedKey string
	var duplicated bool
	
	c = appengine.NewContext(r)
	dao = new(DAO)
	atomTemplate = new(AtomTemplate)

	// フォームデータ取得
	url = r.FormValue("url")
	folderKey = r.FormValue("folder_key")

	// フィード取得
	atom, entries = atomTemplate.Get(c, url)
	
	// フィード追加を試みる
	feedKey, duplicated = dao.RegisterFeed(c, atom, folderKey)
	if duplicated {
		fmt.Fprintf(w, `{"duplicated":true}`)
	} else {
		dao.RegisterEntries(c, entries, feedKey)
		fmt.Fprintf(w, `{"duplicated":false, "key":"%s", "name":"%s"}`, feedKey, atom.Title)
	}

}
