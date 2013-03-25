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
	var parentFolder *Folder
	var encodedParentKey string
	var childKey string
	
	c = appengine.NewContext(r)
	u = user.Current(c)
	dao = new(DAO)
	
	// フォームデータ取得
	name = r.FormValue("folder_name")
	encodedParentKey = r.FormValue("folder_key")
	
	// 親フォルダを取得
	parentFolder = dao.GetFolder(c, encodedParentKey)
	
	// フォルダ新規作成
	childKey = dao.RegisterFolder(c, u, name, false, encodedParentKey)
	
	// 親フォルダに関連付ける
	parentFolder.Children = append(parentFolder.Children, childKey)
	dao.UpdateFolder(c, encodedParentKey, parentFolder)
}