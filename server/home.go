/**
 * ホーム画面
 */

package okareader
import (
	"appengine"
	"net/http"
	"html/template"
)

func home(w http.ResponseWriter, r *http.Request) {
	var t *template.Template
	var err error
	var c appengine.Context

	c = appengine.NewContext(r)
	t, err = template.ParseFiles("server/home.html")
	Check(c, err)
	
	t.Execute(w, []string{})
}