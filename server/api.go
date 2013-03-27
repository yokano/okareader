/**
 * JavaScriptから呼び出すことができるAPI
 */
package okareader
import (
	"net/http"
	"appengine"
	"appengine/user"
)

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