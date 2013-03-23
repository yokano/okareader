/**
 * データモデルの定義とデータストアへのアクセス
 */
package okareader
import (
	"appengine"
	"appengine/datastore"
	"appengine/user"
)

// データの表示時に使用するデータモデル
type Entry struct {
	Id string `xml:"id"`
	Link struct {
		Href string `xml:"href,attr"`
	} `xml:"link"`
	Summary string `xml:"summary"`
	Title string `xml:"title"`
	Updated string `xml:"updated"`
	Owner string
}

type Atom struct {
	Id string `xml:"id"`
	Title string `xml:"title"`
	Entries []*Entry `xml:"entry"`
	Owner string
}

type Folder struct {
	Type string  // root or other
	Name string
	Children []*datastore.Key
	Owner string
	Key string  // encoded string key
}

// データストアに保存する時のデータモデル
// データストアはネストされた構造体が使用できないためアクセス前後に変換する
type Entry_DB struct {
	Id string
	Link string
	Summary string
	Title string
	Updated string
}

type Atom_DB struct {
	Id string
	Title string
	Entries []string
}

// DAO
type DAO struct {
}

/**
 * フォルダの登録
 * @param c {Context} コンテクスト
 * @param u {User} ユーザ
 * @param name {string} フォルダ名
 * @param root {bool} ルートフォルダならtrue
 * @returns {string} 追加したフォルダのキーをエンコードした文字列
 */
func (this *DAO) RegisterFolder(c appengine.Context, u *user.User, name string, root bool) string {
	var folder *Folder
	var key *datastore.Key
	var err error
	var encodedKey string
	
	folder = new(Folder)
	folder.Owner = u.ID
	folder.Children = make([]*datastore.Key, 0)
	if root {
		folder.Type = "root"
		folder.Name = "root"
	} else {
		folder.Type = "other"
		folder.Name = name
	}
	
	key = datastore.NewIncompleteKey(c, "folder", nil)
	folder.Key = key.Encode()
	key, err = datastore.Put(c, key, folder)
	Check(c, err)
	
	encodedKey = key.Encode()
	return encodedKey
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
 * ルートフォルダを取得
 */
func (this *DAO) GetRootFolder(c appengine.Context, u *user.User) *Folder {
	var root *Folder
	var query *datastore.Query
	var iterator *datastore.Iterator
	var err error
	
	root = new(Folder)
	query = datastore.NewQuery("folder").Filter("Type =", "root").Filter("Owner =", u.ID)
	iterator = query.Run(c)
	_, err = iterator.Next(root)
	Check(c, err)
	
	return root
}

/**
 * フォルダの中身を取得する
 * @param {*Folder} folder 親フォルダ
 * @returns {[]interface{}} フォルダの中身を配列化したもの
 */
func (this *DAO) GetChildren(c appengine.Context, folder *Folder) []interface{} {
	var err error
	var children []interface{}
	
	err = datastore.GetMulti(c, folder.Children, children)
	Check(c, err)
	return children
}

/**
 * フィードとそれに含まれるエントリをデータストアに追加
 * 既に存在するものは上書きされる
 * @param feed {Atom} 登録するフィードオブジェクト
 * @param {*Folder} to 追加するフォルダ
 */
func (this *DAO) RegisterFeed(c appengine.Context, feed *Atom, to *Folder) {
	var key *datastore.Key
	var err error
	var atom_db *Atom_DB
	var entry_db *Entry_DB
	var entries_db []*Entry_DB
	
	// データストア用に変換
	atom_db, entries_db = this.Encode(feed)
	
	// フィード保存
	key = datastore.NewKey(c, "feed", atom_db.Id, 0, nil)
	_, err = datastore.Put(c, key, atom_db)
	Check(c, err)
	
	// エントリ保存
	for _, entry_db = range entries_db {
		key = datastore.NewKey(c, "entry", entry_db.Id, 0, nil)
		_, err = datastore.Put(c, key, entry_db)
		Check(c, err)
	}
}

/**
 * フィードをデータストアから読み出す
 * @param c {appengine.Context} コンテキスト
 * @param id {string} フィードのID
 * @retruns {*Atom} 読みだしたフィード
 */
func (this *DAO) GetFeed(c appengine.Context, id string) *Atom {
	var atom_db *Atom_DB
	var entry_db *Entry_DB
	var err error
	var query *datastore.Query
	var iterator *datastore.Iterator
	var atom *Atom
	var entryId string
	var entries_db []*Entry_DB
	
	atom = new(Atom)
	atom_db = new(Atom_DB)
	entries_db = make([]*Entry_DB, 0)
	
	// フィードを読み込む
	query = datastore.NewQuery("feed").Filter("Id =", id)
	iterator = query.Run(c)
	_, err = iterator.Next(atom_db)
	Check(c, err)
	
	// エントリの読み出し
	for _, entryId = range atom_db.Entries {
		query = datastore.NewQuery("entry").Filter("Id =", entryId)
		iterator = query.Run(c)
		entry_db = new(Entry_DB)
		_, err = iterator.Next(entry_db)
		Check(c, err)
		entries_db = append(entries_db, entry_db)
	}

	// 変換
	atom = this.Decode(atom_db, entries_db)
	
	return atom
}

/**
 * データストアに保存できる形式に変換する
 * @param atom {Atom} 変換するフィード
 * @returns entries {[]Entry_DB} 変換後のエントリ
 * @returns atom {Atom_DB} 変換後のAtom
 */
func (this *DAO) Encode(atom *Atom) (*Atom_DB, []*Entry_DB){
	var atom_db *Atom_DB
	var entry_db *Entry_DB
	var entries_db []*Entry_DB
	var entry *Entry
	atom_db = new(Atom_DB)
	atom_db.Entries = make([]string, 0)
	entries_db = make([]*Entry_DB, 0)
	
	// エントリの変換
	for _, entry = range atom.Entries {
		entry_db = new(Entry_DB)
		entry_db.Id = entry.Id
		entry_db.Link = entry.Link.Href
		entry_db.Summary = entry.Summary
		entry_db.Title = entry.Title
		entry_db.Updated = entry.Updated
		
		entries_db = append(entries_db, entry_db)
		atom_db.Entries = append(atom_db.Entries, entry_db.Id)
	}
	
	// Atomの変換
	atom_db.Id = atom.Id
	atom_db.Title = atom.Title
	
	return atom_db, entries_db
}

/**
 * 変換したデータをもとに戻す
 */
func (this *DAO) Decode(atom_db *Atom_DB, entries_db []*Entry_DB) *Atom {
	var atom *Atom
	var entry *Entry
	var entry_db *Entry_DB
	var entries []*Entry
	
	// エントリの変換
	for _, entry_db = range entries_db {
		entry = new(Entry)
		entry.Id = entry_db.Id
		entry.Link.Href = entry_db.Link
		entry.Summary = entry_db.Summary
		entry.Title = entry_db.Title
		entry.Updated = entry_db.Updated
		entries = append(entries, entry)
	}

	// フィードの変換
	atom = new(Atom)
	atom.Id = atom_db.Id
	atom.Title = atom_db.Title
	atom.Entries = entries

	return atom
}