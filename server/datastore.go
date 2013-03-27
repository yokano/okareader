/**
 * データモデルの定義とデータストアへのアクセス
 */
package okareader
import (
	"appengine"
	"appengine/datastore"
	"appengine/user"
	"log"
)

type Folder struct {
	Type string  // root or other
	Title string
	Children []string  // encoded string array
	Owner string
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
 * 既に存在するものは上書きされる
 * @methodOf DAO
 * @param {appengine.Context} c コンテキスト
 * @param {*Atom} feed 登録するフィードオブジェクト
 * @param {string} to 追加先のフォルダのキー
 * @returns {*datastore.Key} 追加したフィードのキー
 */
func (this *DAO) RegisterFeed(c appengine.Context, atom *Atom, to string) *datastore.Key {
	var key *datastore.Key
	var err error
	var parentFolderKey *datastore.Key
	var parentFolder *Folder
	
	// フィード保存
	log.Printf("Id:%s\n", atom.Id)
	key = datastore.NewKey(c, "feed", atom.Id, 0, nil)
	key, err = datastore.Put(c, key, atom)
	Check(c, err)
	
	// フォルダに追加
	parentFolderKey, err = datastore.DecodeKey(to)
	Check(c, err)
	parentFolder = new(Folder)
	err = datastore.Get(c, parentFolderKey, parentFolder)
	Check(c, err)
	parentFolder.Children = append(parentFolder.Children, atom.Id)
	_, err = datastore.Put(c, parentFolderKey, parentFolder)
	Check(c, err)
	
	return key
}

/**
 * エントリをフィードに追加する
 * @methodOf DAO
 * @param {appengine.Context} c コンテキスト
 * @param {[]*Entry} entries 追加するエントリ配列
 * @param {string} to 追加先のフィードのキー
 */
func (this *DAO) RegisterEntries(c appengine.Context, entries []*Entry, to string) {
	// エントリ保存
//	for _, entry_db = range entries_db {
//		key = datastore.NewKey(c, "entry", entry_db.Id, 0, nil)
//		_, err = datastore.Put(c, key, entry_db)
//		Check(c, err)
//	}
}

/**
 * フィードをデータストアから読み出す
 * @param c {appengine.Context} コンテキスト
 * @param id {string} フィードのID
 * @retruns {*Atom} 読みだしたフィード
 */
func (this *DAO) GetFeed(c appengine.Context, id string) *Atom {
	var atom *Atom
	var err error
	var query *datastore.Query
	var iterator *datastore.Iterator

	atom = new(Atom)
	
	// フィードを読み込む
	query = datastore.NewQuery("feed").Filter("Id =", id)
	iterator = query.Run(c)
	_, err = iterator.Next(atom)
	Check(c, err)
		
	return atom
}
