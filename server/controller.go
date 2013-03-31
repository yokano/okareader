/**
 * ブラウザやAjaxのリクエストを適切な処理へ振り分ける
 * 直接データストアへアクセスしたり画面を描画したりしてはいけない
 * データストアへのアクセスはModelに,画面の描画はViewに頼むこと.
 */

package okareader

import(
	"appengine"
	"appengine/user"
	"net/http"
	"fmt"
	"appengine/urlfetch"
	"encoding/xml"
)

type Controller struct {

}

/**
 * リクエストURLによる処理の振り分け
 * /api/*** はAjaxによるAPIへのアクセスであり画面の描画は不用
 * それ以外はページ遷移を表し画面を描画する
 * @function
 */
func init() {
	var controller *Controller
	controller = new(Controller)
	
	// ルートフォルダの表示
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		controller.home(w, r)
	})
	
	// フォルダ画面
	http.HandleFunc("/folder", func(w http.ResponseWriter, r *http.Request) {
		controller.folder(w, r)
	})
	
	// フィード画面
	http.HandleFunc("/feed", func(w http.ResponseWriter, r *http.Request) {
		controller.feed(w, r)
	})
	
	// フォルダの追加API
	http.HandleFunc("/api/addfolder", func(w http.ResponseWriter, r *http.Request) {
		controller.addFolder(w, r)
	})
	
	// フィードの追加API
	http.HandleFunc("/api/addfeed", func(w http.ResponseWriter, r *http.Request) {
		controller.addFeed(w, r)
	})
	
	// データ削除
	http.HandleFunc("/clear", func(w http.ResponseWriter, r *http.Request) {
		controller.clear(w, r)
	})
}

/**
 * http://okareader.appspot.com/ へアクセスした時の処理
 * ログインしていなければログインさせる
 * ログインしていればルートフォルダを表示
 * @methodOf Controller
 * @param {http.ResponseWriter} 応答先
 * @param {*http.Request} リクエスト
 */
func (this *Controller)home(w http.ResponseWriter, r *http.Request) {
	var c appengine.Context
	var u *user.User
	var root *Folder
	var dao *DAO
	var view *View
	var key string

	c = appengine.NewContext(r)
	u = user.Current(c)
	dao = new(DAO)
	view = new(View)
	
	if u == nil {
		view.ShowLogin(c, w)
	} else {
		key, root = dao.GetRootFolder(c, u)
		if root.Type == "" {
			key = dao.RegisterFolder(c, u, "root", true, "")
		}
		view.ShowFolder(c, key, "", w)
	}
}

/**
 * http://okareader.appspot.com/folder へアクセスしたらフォルダを表示
 * フォルダをデータストアから取得するためのキーをGETで受け取る
 * @methodOf Controller
 * @param {http.ResponseWriter} w 応答先
 * @param {*http.Request} r リクエスト
 * @param {HTTP GET} key エンコード済みのフォルダキー
 * @param {HTTP GET} from 遷移前のフォルダのキー
 */
func (this *Controller) folder(w http.ResponseWriter, r *http.Request) {
	var c appengine.Context
	var view *View
	var encodedKey string
	var from string
	
	c = appengine.NewContext(r)
	encodedKey = r.FormValue("key")
	from = r.FormValue("from")

	view = new(View)
	view.ShowFolder(c, encodedKey, from, w)
}

/**
 * http://okareader.appspot.com/feed へアクセスしたらフィードを表示
 * フィードのキーはGETで渡される
 * @methodOf Controller
 * @param {http.ResponseWriter} w 応答先
 * @param {*http.Request} r リクエスト
 * @param {HTTP GET} key エンコード済みのフィードキー
 * @param {HTTP GET} from 遷移前のフォルダのキー
 */
func (this *Controller) feed(w http.ResponseWriter, r *http.Request) {
	var c appengine.Context
	var view *View
	var feedKey string
	var from string
	
	c = appengine.NewContext(r)
	feedKey = r.FormValue("key")
	from = r.FormValue("from")
	
	view = new(View)
	view.ShowFeed(c, feedKey, from, w)
}

/**
 * API:フォルダの新規追加
 * @methodOf Controller
 * @param {http.ResponseWriter} w 応答先
 * @param {*http.Request} r HTTPリクエスト
 * @returns {AJAX:JSON} 追加したフォルダのキーを含むJSON
 */
func (this *Controller) addFolder(w http.ResponseWriter, r *http.Request) {
	var c appengine.Context
	var u *user.User
	var dao *DAO
	var title string
	var encodedParentKey string
	var resultKey string
	
	c = appengine.NewContext(r)
	u = user.Current(c)
	dao = new(DAO)
	
	title = r.FormValue("folder_name")
	encodedParentKey = r.FormValue("folder_key")
	
	resultKey = dao.RegisterFolder(c, u, title, false, encodedParentKey)
	fmt.Fprintf(w, `{"key":"%s"}`, resultKey)
}

/**
 * API:フィードの登録
 * @methodOf Controller
 * @param {http.ResponseWriter} w 応答先
 * @param {*http.Request} r リクエスト
 */
func (this *Controller) addFeed(w http.ResponseWriter, r *http.Request) {
	var url string
	var folderKey string
	var dao *DAO
	var c appengine.Context
	var feed *Feed
	var atomTemplate *AtomTemplate
	var entries []*Entry
	var feedKey string
	var duplicated bool
	var xml []byte
	var feedType string
	
	c = appengine.NewContext(r)
	dao = new(DAO)
	atomTemplate = new(AtomTemplate)
	
	// フォームデータ取得
	url = r.FormValue("url")
	folderKey = r.FormValue("folder_key")
	
	// フィード取得
	feed, entries = atomTemplate.Get(c, url)
	
	// XML取得
	xml = this.getXML(c, url)
	feedType = this.getType(c, xml)
	
	// フィード追加を試みる
	feedKey, duplicated = dao.RegisterFeed(c, feed, folderKey)
	if duplicated {
		fmt.Fprintf(w, `{"duplicated":true}`)
	} else {
		dao.RegisterEntries(c, entries, feedKey)
		fmt.Fprintf(w, `{"duplicated":false, "key":"%s", "name":"%s"}`, feedKey, feed.Title)
	}
}

/**
 * 指定されたURLからXMLファイルを受信して返す
 * @methodOf Controller
 * @param {appengine.Context} c コンテキスト
 * @param {string} url URL
 * @returns {[]byte} 受信したXMLデータ
 */
func (this *Controller) getXML(c appengine.Context, url string) []byte {
	var client *http.Client
	var response *http.Response
	var err error
	var result []byte
	
	client = urlfetch.Client(c)
	response, err = client.Get(url)
	Check(c, err)
	
	result = make([]byte, response.ContentLength)
	_, err = response.Body.Read(result)
	Check(c, err)
	
	return result
}

/**
 * XMLデータの規格を判断する
 * @methodOf Controller
 * @param {[]byte} bytes XMLデータ
 * @returns フィードの規格(RSS1.0 / RSS2.0 / Atom / etc)
 */
func (this *Controller) getType(c appengine.Context, bytes []byte) string {
	type Checker struct {
		XMLName xml.Name
	}
	var checker *Checker
	var err error
	var result string
	
	checker = new(Checker)
	err = xml.Unmarshal(bytes, checker)
	Check(c, err)
	
	switch checker.XMLName.Local {
		case "feed":
			result = "Atom"
		case "rss":
			result = "RSS2.0"
		case "rdf":
			result = "RSS1.0"
		default:
			result = "etc"
	}
	
	return result
}

/**
 * データを削除する
 * @methodOf Controller
 * @param {http.ResponseWriter} w 応答先
 * @param {http.Request} r リクエスト
 */
func (this *Controller) clear(w http.ResponseWriter, r *http.Request) {
	var c appengine.Context
	var dao *DAO
	c = appengine.NewContext(r)
	dao = new(DAO)
	dao.clear(c)
	this.home(w, r)
}
