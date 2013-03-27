package okareader
import (
	"appengine"
	"net/http"
	"log"
)

func atom_test(w http.ResponseWriter, r *http.Request) {
	var entry *Entry
	var i int
	var c appengine.Context
	var atom *Atom
	var entries []*Entry
	var atomTemplate *AtomTemplate
	
	c = appengine.NewContext(r)
	atom = new(Atom)
	atomTemplate = new(AtomTemplate)
	atom, entries = atomTemplate.Get(c, "http://feed.rssad.jp/rss/gigazine/rss_atom")
	
	log.Printf("Title:%s", atom.Title)
	log.Printf("EntriesNum:%d", len(entries))
	for i, entry = range entries {
		log.Printf("[%d]:%s\n", i, entry.Title)
	}
}
