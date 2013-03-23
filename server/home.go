/**
 * ホーム画面
 */

package okareader
import (
	"appengine"
	"appengine/user"
	"net/http"
	"html/template"
)

func home(w http.ResponseWriter, r *http.Request) {
	var t *template.Template
	var err error
	var c appengine.Context
	var u *user.User
	var html string
	var content map[string]interface{}
	var root *Folder
	var dao *DAO

	c = appengine.NewContext(r)
	u = user.Current(c)
	dao = new(DAO)
	content = make(map[string]interface{})
	
	if u == nil {
		html = "server/login.html"
		content["LoginURL"], err = user.LoginURL(c, "/")
		Check(c, err)
	} else {
		html = "server/home.html"
		content["LogoutURL"], err = user.LogoutURL(c, "/")
		Check(c, err)
		
		// ルートフォルダが存在しなければ作成
		root = dao.GetRootFolder(c, u)
		if root.Type == "" {
			dao.RegisterFolder(c, u, "root", true)
		}
		
		// ルートフォルダを表示
		content["Children"] = dao.GetChildren(c, root)
	}
	t, err = template.ParseFiles(html)
	Check(c, err)
	
	t.Execute(w, content)
}