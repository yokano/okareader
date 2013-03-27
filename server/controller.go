/**
 * クライアントからのリクエストを振り分ける
 * ブラウザからのアクセスやAjaxによるアクセスに対応する
 */

package okareader

import(
	"appengine"
	"appengine/user"
	"net/http"
)

/**
 * リクエストURLによる処理の振り分け
 * /api/*** はAjaxによるアクセスに対応するもので api.go 内に処理を記述
 * それ以外のURLはこのファイルで処理する
 * @func
 */
func init() {
	http.HandleFunc("/", home)
	http.HandleFunc("/atom", atom_test)
	http.HandleFunc("/api/addfolder", addFolder)
	http.HandleFunc("/api/addfeed", addFeed)
	http.HandleFunc("/folder", folder)
}

/**
 * http://okareader.appspot.com/ へアクセスした時の処理
 * ログインしていなければログインさせる
 * ログインしていればルートフォルダを表示
 * @func
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
		view.ShowFolder(c, key, root, w)
	}
}

/**
 * フォルダの追加
 * @param {http.ResponseWriter} w 応答先
 * @param {*http.Request} r HTTPリクエスト
 */
func addFolder(w http.ResponseWriter, r *http.Request) {
	var c appengine.Context
	var u *user.User
	var dao *DAO
	var title string
	var encodedParentKey string
	
	c = appengine.NewContext(r)
	u = user.Current(c)
	dao = new(DAO)
	
	// フォームデータ取得
	title = r.FormValue("folder_name")
	encodedParentKey = r.FormValue("folder_key")
	
	// フォルダ新規作成
	dao.RegisterFolder(c, u, title, false, encodedParentKey)
}

/**
 * フィードの追加
 * @param {http.ResponseWriter} w 応答先
 * @param {*http.Request}
 */
func addFeed(w http.ResponseWriter, r *http.Request) {
	 var url string
	 var encodedParentKey string
	 var dao *DAO
	 var c appengine.Context
	 var atom *Atom
	 var atomTemplate *AtomTemplate
	 
	 c = appengine.NewContext(r)
	 dao = new(DAO)
	 atomTemplate = new(AtomTemplate)
	 
	 // フォームデータ取得
	 url = r.FormValue("url")
	 encodedParentKey = r.FormValue("folder_key")
	 
	 // フィード取得
	 atom, _ = atomTemplate.Get(c, url)
	 
	 // フィード追加
	 dao.RegisterFeed(c, atom, encodedParentKey)
}

/**
 * http://okareader.appspot.com/folder へアクセスしたらフォルダを表示
 * フォルダをデータストアから取得するためのキーをGETで受け取る
 * @func
 * @param {http.ResponseWriter} 応答先
 * @param {*http.Request} リクエスト
 * 
 */
func folder(w http.ResponseWriter, r *http.Request) {

}
