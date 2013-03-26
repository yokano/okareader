/**
 * JavaScriptから呼び出すことができるAPI
 */
package okareader
import (
	"net/http"
	"appengine"
	"appengine/user"
)

func addFolder(w http.ResponseWriter, r *http.Request) {
	var c appengine.Context
	var u *user.User
	var dao *DAO
	var name string
	var encodedParentKey string
	
	c = appengine.NewContext(r)
	u = user.Current(c)
	dao = new(DAO)
	
	// フォームデータ取得
	name = r.FormValue("folder_name")
	encodedParentKey = r.FormValue("folder_key")
	
	// フォルダ新規作成
	dao.RegisterFolder(c, u, name, false, encodedParentKey)
}