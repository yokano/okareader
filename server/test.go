package okareader
import (
	"appengine"
	"net/http"
	"text/template"
)

func atom_test(w http.ResponseWriter, r *http.Request) {
	var c appengine.Context
	var t *template.Template
	var err error
	var atom *Atom
	
	c = appengine.NewContext(r)
	atom = get(c, "http://feed.rssad.jp/rss/gigazine/rss_atom")
	t, err = template.ParseFiles("server/feed.html")
	Check(c, err)
	
	err = t.Execute(w, atom)
	Check(c, err)
}

func dao_test(w http.ResponseWriter, r *http.Request) {
	var dao *DAO
	var c appengine.Context
	var atom *Atom

	c = appengine.NewContext(r)
	atom = get(c, "http://feed.rssad.jp/rss/gigazine/rss_atom")
	dao = new(DAO)
	dao.RegisterFeed(c, atom)
}