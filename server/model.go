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
}

type Feed struct {
	Id string
	Title string
	Entries []string
}

// DAO
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
func (this *DAO) RegisterFolder(c appengine.Context, u *user.User, title string, root bool, encodedParentKey string) string {
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
	Check(c, err)
	
	encodedKey = key.Encode()
	
	// 親フォルダの子に登録
	if !root {
		
		// 親のChildrenに子のキーを追加して上書きする
		parentKey, err = datastore.DecodeKey(encodedParentKey)
		Check(c, err)
		
		parentFolder = new(Folder)
		err = datastore.Get(c, parentKey, parentFolder)
		Check(c, err)
		
		parentFolder.Children = append(parentFolder.Children, encodedKey)
		
		_, err = datastore.Put(c, parentKey, parentFolder)
		Check(c, err)
	}
	
	return encodedKey
}

/**
 * フォルダの更新
 * @param {appengine.Context} c コンテキスト
 * @param {string} encodedKey 文字列
 */
func (this *DAO) UpdateFolder(c appengine.Context, encodedKey string, folder *Folder) {
	var key *datastore.Key
	var err error
	
	key, err = datastore.DecodeKey(encodedKey)
	Check(c, err)
	
	_, err = datastore.Put(c, key, folder)
	Check(c, err)
}

/**
 * フォルダの削除
 * 中身も全て削除する
 * rootフォルダは削除不可
 * @param encodedKey {string} 削除するフォルダのキーをエンコードした文字列
 */
func (this *DAO) RemoveFolder(c appengine.Context, encodedKey string) {
	var err error
	var key *datastore.Key
	
	key, err = datastore.DecodeKey(encodedKey)
	Check(c, err)
	
	err = datastore.Delete(c, key)
	Check(c, err)
}

/**
 * フォルダの取得
 * @param {appengine.Context} c コンテキスト
 * @param {string} encodedKey 取得したいフォルダのキーをエンコードした文字列
 */
func (this *DAO) GetFolder(c appengine.Context, encodedKey string) *Folder {
	var key *datastore.Key
	var err error
	var folder *Folder
	
	key, err = datastore.DecodeKey(encodedKey)
	Check(c, err)
	
	folder = new(Folder)
	err = datastore.Get(c, key, folder)
	Check(c, err)
	
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
func (this *DAO) GetItem(c appengine.Context, encodedKey string) (string, interface{}) {
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
	Check(c, err)
	
	item = new(Item)
	err = datastore.Get(c, key, item)
	Check(c, err)
	
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
func (this *DAO) GetRootFolder(c appengine.Context, u *user.User) (string, *Folder) {
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
	Check(c, err)
	
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
func (this *DAO) GetChildren(c appengine.Context, folder *Folder) []interface{} {
	var err error
	var children []interface{}
	var keys []*datastore.Key
	var encodedKey string
	var key *datastore.Key
	
	keys = make([]*datastore.Key, 0)
	for _, encodedKey = range folder.Children {
		key, err = datastore.DecodeKey(encodedKey)
		Check(c, err)
		keys = append(keys, key)
	}
	
	err = datastore.GetMulti(c, keys, children)
	Check(c, err)
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
func (this *DAO) RegisterFeed(c appengine.Context, feed *Feed, to string) (string, bool) {
	var key *datastore.Key
	var encodedKey string
	var err error
	var parentFolderKey *datastore.Key
	var parentFolder *Folder
	var duplicated bool
	
	key = datastore.NewKey(c, "feed", feed.Id, 0, nil)
	encodedKey = key.Encode()
	
	// 重複していたら登録しない
	duplicated = this.Exist(c, encodedKey)
	if duplicated {
		encodedKey = ""
	} else {
		// フィード保存
		_, err = datastore.Put(c, key, feed)
		Check(c, err)
		
		// 親フォルダ取得
		parentFolderKey, err = datastore.DecodeKey(to)
		Check(c, err)
		parentFolder = new(Folder)
		err = datastore.Get(c, parentFolderKey, parentFolder)
		Check(c, err)
		
		// 親フォルダの子に追加
		parentFolder.Children = append(parentFolder.Children, encodedKey)
		_, err = datastore.Put(c, parentFolderKey, parentFolder)
		Check(c, err)
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
func (this *DAO) RegisterEntries(c appengine.Context, entries []*Entry, to string) []string {
	var entry *Entry
	var key *datastore.Key
	var result []string
	var err error
	var i int
	var feed *Feed
	var feedKey *datastore.Key
	
	feedKey, err = datastore.DecodeKey(to)
	feed = this.GetFeed(c, to)
	
	result = make([]string, len(entries))
	for i, entry = range entries {
		key = datastore.NewKey(c, "entry", entry.Id, 0, nil)
		_, err = datastore.Put(c, key, entry)
		Check(c, err)
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
func (this *DAO) GetEntries(c appengine.Context, feedKey string) []*Entry {
	var feed *Feed
	var entryKey string
	var entry *Entry
	var key *datastore.Key
	var err error
	var entries []*Entry
	
	entries = make([]*Entry, 0)
	
	feed = this.GetFeed(c, feedKey)
	for _, entryKey = range feed.Entries {
		key, err = datastore.DecodeKey(entryKey)
		Check(c, err)
		
		entry = new(Entry)
		err = datastore.Get(c, key, entry)
		Check(c, err)
		
		entries = append(entries, entry)
	}
	
	return entries
}

/**
 * フィードをデータストアから読み出す
 * @param {appengine.Context} c コンテキスト
 * @param {string} feedKey エンコード済みのフィードキー
 * @retruns {*Feed} フィード
 */
func (this *DAO) GetFeed(c appengine.Context, feedKey string) *Feed {
	var feed *Feed
	var err error
	var key *datastore.Key

	feed = new(Feed)
	key, err = datastore.DecodeKey(feedKey)
	Check(c, err)
	
	err = datastore.Get(c, key, feed)
	Check(c, err)
		
	return feed
}

/**
 * 指定されたキーのデータが既に存在するか調べる
 * フィードやエントリなど重複させたくないデータはこの関数を使ってチェックする
 * @methodOf DAO
 * @param {appengine.Context} c コンテキスト
 * @param {string} encodedKey エンコード済みのキー
 * @returns {bool} 重複していたらtrue
 */
func (this *DAO) Exist(c appengine.Context, encodedKey string) bool {
	var result bool
	var key *datastore.Key
	var item interface{}
	var err error
	
	key, err = datastore.DecodeKey(encodedKey)
	Check(c, err)
	
	item = new(interface{})
	err = datastore.Get(c, key, item)
	if err == datastore.ErrNoSuchEntity {
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
		Check(c, err)
		datastore.DeleteMulti(c, keys)
	}
}
