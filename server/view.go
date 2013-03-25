/**
 *  @file データを受け取ってページを表示する処理
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
 */
func (this *View) ShowFolder(c appengine.Context, key string, folder *Folder, w http.ResponseWriter) {
	var content map[string]interface{}
	var err error
	var t *template.Template
	var children []interface{}
	var child interface{}
	var encodedKey string
	var dao *DAO
	
	dao = new(DAO)
	
	content = make(map[string]interface{}, 0)
	content["LogoutURL"], err = user.LogoutURL(c, "/")
	content["FolderKey"] = key
	Check(c, err)
	
	children = make([]interface{}, 0)
	for _, encodedKey = range folder.Children {
		child = dao.GetFolder(c, encodedKey)
		children = append(children, child)
	}
	content["Children"] = children
	
	t, err = template.ParseFiles("server/home.html")
	Check(c, err)
	
	t.Execute(w, content)
}

/**
 * ログインを促す画面を表示
 * @methodOf View
 */
func (this *View) ShowLogin(c appengine.Context, w http.ResponseWriter) {
	var content map[string]interface{}
	var err error
	var t *template.Template
	
	t, err = template.ParseFiles("server/login.html")
	Check(c, err)
	
	content = make(map[string]interface{}, 0)
	content["LoginURL"], err = user.LoginURL(c, "/")
	Check(c, err)
	
	t.Execute(w, content)
}