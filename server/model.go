/**
 * データモデルの定義とデータストアへのアクセス
 * データストアへのアクセスはここからだけ行う
 */
package okareader
import (
	"appengine"
	"appengine/datastore"
	"appengine/user"
)

type Folder struct {
	Type string  // root or other
	Title string
	Children []string  // encoded string array
	Owner string
}

type Entry struct {
	Id string
	Link string
	Title string
	Updated string
	Owner string
}

type Feed struct {
	Id string
	Title string
	Entries []string
	Owner string
}

type DAO struct {
}

/**
 * フォルダの新規登録
 * @param c {Context} コンテクスト
 * @param u {User} ユーザ
 * @param title {string} フォルダ名
 * @param root {bool} ルートフォルダならtrue
 * @param encodedParentKey {string} 追加先の親フォルダのキー
 * @returns {string} 追加したフォルダのキーをエンコードした文字列
 */
func (this *DAO) registerFolder(c appengine.Context, u *user.User, title string, root bool, encodedParentKey string) string {
	var folder *Folder
	var key *datastore.Key
	var err error
	var encodedKey string
	var parentFolder *Folder
	var parentKey *datastore.Key
	
	// 追加するフォルダの作成
	folder = new(Folder)
	folder.Owner = u.ID
	folder.Children = make([]string, 0)
	if root {
		folder.Type = "root"
		folder.Title = "root"
	} else {
		folder.Type = "other"
		folder.Title = title
	}
	
	// 追加するフォルダをデータストアに保存
	key = datastore.NewIncompleteKey(c, "folder", nil)
	key, err = datastore.Put(c, key, folder)
	check(c, err)
	
	encodedKey = key.Encode()
	
	// 親フォルダの子に登録
	if !root {
		
		// 親のChildrenに子のキーを追加して上書きする
		parentKey, err = datastore.DecodeKey(encodedParentKey)
		check(c, err)
		
		parentFolder = new(Folder)
		err = datastore.Get(c, parentKey, parentFolder)
		check(c, err)
		
		parentFolder.Children = append(parentFolder.Children, encodedKey)
		
		_, err = datastore.Put(c, parentKey, parentFolder)
		check(c, err)
	}
	
	return encodedKey
}

/**
 * フォルダの更新
 * @param {appengine.Context} c コンテキスト
 * @param {string} encodedKey 文字列
 */
func (this *DAO) updateFolder(c appengine.Context, encodedKey string, folder *Folder) {
	var key *datastore.Key
	var err error
	
	key, err = datastore.DecodeKey(encodedKey)
	check(c, err)
	
	_, err = datastore.Put(c, key, folder)
	check(c, err)
}

/**
 * フォルダの削除
 * 中身も全て削除する
 * rootフォルダは削除不可
 * @param encodedKey {string} 削除するフォルダのキーをエンコードした文字列
 */
func (this *DAO) removeFolder(c appengine.Context, encodedKey string) {
	var err error
	var key *datastore.Key
	
	key, err = datastore.DecodeKey(encodedKey)
	check(c, err)
	
	err = datastore.Delete(c, key)
	check(c, err)
}

/**
 * フォルダの取得
 * @param {appengine.Context} c コンテキスト
 * @param {string} encodedKey 取得したいフォルダのキーをエンコードした文字列
 */
func (this *DAO) getFolder(c appengine.Context, encodedKey string) *Folder {
	var key *datastore.Key
	var err error
	var folder *Folder
	
	key, err = datastore.DecodeKey(encodedKey)
	check(c, err)
	
	folder = new(Folder)
	err = datastore.Get(c, key, folder)
	check(c, err)
	
	return folder
}

/**
 * フォルダ・フィードの取得
 * フォルダの中身を表示するときなど取り出す対象がどちらかわからないときに使用する
 * @methodOf DAO
 * @param {appengine.Context} c コンテキスト
 * @param {string} encodedKey アイテムのエンコード済みのキー
 * @returns {string} 取得したアイテムがフォルダなら"folder",フィードなら"feed"
 * @returns {*Folder or *Feed} 取得したフォルダまたはフィードオブジェクト
 */
func (this *DAO) getItem(c appengine.Context, encodedKey string) (string, interface{}) {
	var key *datastore.Key
	var err error
	type Item struct {
		Title string
		Owner string
		Entries []string
		Type string
		Id string 
	}
	var item *Item
	var result interface{}
	var itemType string
	
	key, err = datastore.DecodeKey(encodedKey)
	check(c, err)
	
	item = new(Item)
	err = datastore.Get(c, key, item)
	check(c, err)
	
	if item.Type == "" {
		// 要素はFeed
		result = new(Feed)
		result = item
		itemType = "feed"
	} else {
		// 要素はフォルダ
		result = new(Folder)
		result = item
		itemType = "folder"
	}
	
	return itemType, result
}

/**
 * ルートフォルダを取得
 */
func (this *DAO) getRootFolder(c appengine.Context, u *user.User) (string, *Folder) {
	var root *Folder
	var query *datastore.Query
	var iterator *datastore.Iterator
	var err error
	var key *datastore.Key
	var encodedKey string
	
	root = new(Folder)
	query = datastore.NewQuery("folder").Filter("Type =", "root").Filter("Owner =", u.ID)
	iterator = query.Run(c)
	key, err = iterator.Next(root)
	check(c, err)
	
	if key != nil {
		encodedKey = key.Encode()
	}
	
	return encodedKey, root
}

/**
 * フォルダの中身を取得する
 * @param {*Folder} folder 親フォルダ
 * @returns {[]interface{}} フォルダの中身を配列化したもの
 */
func (this *DAO) getChildren(c appengine.Context, folder *Folder) []interface{} {
	var err error
	var children []interface{}
	var keys []*datastore.Key
	var encodedKey string
	var key *datastore.Key
	
	keys = make([]*datastore.Key, 0)
	for _, encodedKey = range folder.Children {
		key, err = datastore.DecodeKey(encodedKey)
		check(c, err)
		keys = append(keys, key)
	}
	
	err = datastore.GetMulti(c, keys, children)
	check(c, err)
	return children
}

/**
 * フィードをデータストアに追加
 * 既に存在するフィードは無視する
 * @methodOf DAO
 * @param {appengine.Context} c コンテキスト
 * @param {*Feed} feed 登録するフィードオブジェクト
 * @param {string} to 追加先のフォルダのキー
 * @returns {string} 追加したフィードのキーをエンコードしたもの　重複していたら空文字列
 * @returnss {bool} 重複していた場合はtrue
 */
func (this *DAO) registerFeed(c appengine.Context, feed *Feed, to string) (string, bool) {
	var key *datastore.Key
	var encodedKey string
	var err error
	var parentFolderKey *datastore.Key
	var parentFolder *Folder
	var duplicated bool
	var u *user.User
	
	key = datastore.NewIncompleteKey(c, "feed", nil)
	
	// 重複していたら登録しない
	duplicated = this.exist(c, feed)
	if duplicated {
		encodedKey = ""
	} else {
		// ユーザID追加
		u = user.Current(c)
		feed.Owner = u.ID
		
		// フィード保存
		key, err = datastore.Put(c, key, feed)
		check(c, err)
		encodedKey = key.Encode()
		
		// 親フォルダ取得
		parentFolderKey, err = datastore.DecodeKey(to)
		check(c, err)
		parentFolder = new(Folder)
		err = datastore.Get(c, parentFolderKey, parentFolder)
		check(c, err)
		
		// 親フォルダの子に追加
		parentFolder.Children = append(parentFolder.Children, encodedKey)
		_, err = datastore.Put(c, parentFolderKey, parentFolder)
		check(c, err)
	}
	
	return encodedKey, duplicated
}

/**
 * エントリをフィードに追加する
 * @methodOf DAO
 * @param {appengine.Context} c コンテキスト
 * @param {[]*Entry} entries 追加するエントリ配列
 * @param {string} to 追加先のフィードのキー
 * @returns {[]string} 追加したエントリのキー配列
 */
func (this *DAO) registerEntries(c appengine.Context, entries []*Entry, to string) []string {
	var entry *Entry
	var key *datastore.Key
	var result []string
	var err error
	var i int
	var feed *Feed
	var feedKey *datastore.Key
	var u *user.User
	
	feedKey, err = datastore.DecodeKey(to)
	feed = this.getFeed(c, to)
	
	u = user.Current(c)
	
	result = make([]string, len(entries))
	for i, entry = range entries {
		entry.Owner = u.ID
		key = datastore.NewIncompleteKey(c, "entry", nil)
		key, err = datastore.Put(c, key, entry)
		check(c, err)
		result[i] = key.Encode()
		feed.Entries = append(feed.Entries, result[i])
	}
	
	_, err = datastore.Put(c, feedKey, feed)
	
	return result
}

/**
 * 指定されたフィードのエントリをすべて返す
 * @methodOf DAO
 * @param {appengine.Context} c コンテキスト
 * @param {string} feedKey エンコード済みのフィードキー
 * @returns {[]*Entry} エントリ配列
 */
func (this *DAO) getEntries(c appengine.Context, feedKey string) []*Entry {
	var feed *Feed
	var entryKey string
	var entry *Entry
	var key *datastore.Key
	var err error
	var entries []*Entry
	
	entries = make([]*Entry, 0)
	
	feed = this.getFeed(c, feedKey)
	for _, entryKey = range feed.Entries {
		key, err = datastore.DecodeKey(entryKey)
		check(c, err)
		
		entry = new(Entry)
		err = datastore.Get(c, key, entry)
		check(c, err)
		
		entries = append(entries, entry)
	}
	
	return entries
}

/**
 * 指定されたエントリを削除する
 * @methodOf DAO
 * @param {appengine.Context} c コンテキスト
 * @param {string} id 削除するエントリのID
 * @param {string} feedKey エントリが登録されているフィードのキー
 */
func (this *DAO) removeEntry(c appengine.Context, id string, feedKey string) {
	var query *datastore.Query
	var iterator *datastore.Iterator
	var key *datastore.Key
	var err error
	var feed *Feed
	var encodedEntryKey string
	
	query = datastore.NewQuery("entry").Filter("Id =", id)
	iterator = query.Run(c)
	key, err = iterator.Next(nil)
	check(c, err)
	
	err = datastore.Delete(c, key)
	check(c, err)
	
	encodedEntryKey = key.Encode()
	
	key, err = datastore.DecodeKey(feedKey)
	check(c, err)
	
	feed = new(Feed)
	err = datastore.Get(c, key, feed)
	check(c, err)
	
	feed.Entries = removeItem(feed.Entries, encodedEntryKey)
	datastore.Put(c, key, feed)
}

/**
 * フィードをデータストアから読み出す
 * @param {appengine.Context} c コンテキスト
 * @param {string} feedKey エンコード済みのフィードキー
 * @retruns {*Feed} フィード
 */
func (this *DAO) getFeed(c appengine.Context, feedKey string) *Feed {
	var feed *Feed
	var err error
	var key *datastore.Key

	feed = new(Feed)
	key, err = datastore.DecodeKey(feedKey)
	check(c, err)
	
	err = datastore.Get(c, key, feed)
	check(c, err)
	
	return feed
}

/**
 * 指定されたキーのデータが既に存在するか調べる
 * フィードやエントリなど重複させたくないデータはこの関数を使ってチェックする
 * @methodOf DAO
 * @param {appengine.Context} c コンテキスト
 * @param {*Feed} feed フィード
 * @returns {bool} 重複していたらtrue
 */
func (this *DAO) exist(c appengine.Context, feed *Feed) bool {
	var result bool
	var err error
	var query *datastore.Query
	var u *user.User
	var count int
	
	u = user.Current(c)
	query = datastore.NewQuery("feed").Filter("Id =", feed.Id).Filter("Owner =", u.ID)
	count, err = query.Count(c)
	check(c, err)
	
	if count == 0 {
		result = false
	} else {
		result = true
	}
	
	return result
}

/**
 * すべてのデータを削除する
 * デバッグ用
 * @methodOf DAO
 * @param {appengine.Context} c
 */
func (this *DAO) clear(c appengine.Context) {
	var keys []*datastore.Key
	var query *datastore.Query
	var err error
	var kinds [3]string
	var kind string
	
	keys = make([]*datastore.Key, 0)
	kinds = [3]string{"folder", "feed", "entry"}
	
	for _, kind = range kinds {
		query = datastore.NewQuery(kind).KeysOnly()
		keys, err = query.GetAll(c, nil)
		check(c, err)
		datastore.DeleteMulti(c, keys)
	}
}
