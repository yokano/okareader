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
func (this *View) ShowFolder(c appengine.Context, folder *Folder, w http.ResponseWriter) {
	var content map[string]interface{}
	var err error
	var t *template.Template
	
	content = make(map[string]interface{}, 0)
	content["LogoutURL"], err = user.LogoutURL(c, "/")
	Check(c, err)
	
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