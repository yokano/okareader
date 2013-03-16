/**
 * Atomを読み込んでデータに変換する
 */
package okareader
import (
	"appengine"
	"appengine/urlfetch"
	"net/http"
	"log"
)

func get(c appengine.Context, url string) {
	var client *http.Client
	var response *http.Response
	var err error
	
	client = urlfetch.Client(c)
	response, err = client.Get(url)
	Check(c, err)
	log.Printf("%s", response.Body)
}