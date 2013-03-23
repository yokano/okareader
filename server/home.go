/**
 * ホーム画面
 */

package okareader
import (
	"appengine"
	"appengine/user"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	var c appengine.Context
	var u *user.User
	var root *Folder
	var dao *DAO
	var view *View

	c = appengine.NewContext(r)
	u = user.Current(c)
	dao = new(DAO)
	view = new(View)
	
	// ログインしていなければログインページを表示
	if u == nil {
		view.ShowLogin(c, w)
	} else {
		// ルートフォルダが存在しなければ作成
		root = dao.GetRootFolder(c, u)
		if root.Type == "" {
			dao.RegisterFolder(c, u, "root", true)
		}
		
		// ルートフォルダを表示
		view.ShowFolder(c, root, w)
	}
}