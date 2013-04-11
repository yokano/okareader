/**
 * データモデルの定義とデータストアへのアクセス
 * データストアへのアクセスはここからだけ行う
 */
package okareader
import (
	"appengine"
	"appengine/datastore"
	"appengine/user"
	"encoding/xml"
	"log"
)

/**
 * MVCでいうModelの役割
 * データ操作全般
 * @class
 */
type DAO struct {
}

/**
 * フォルダ
 * @class
 * @member {string} Type ルートフォルダなら"root" それ以外は"other"
 * @member {string} Title フォルダのタイトル
 * @member {[]string} Children 子への参照キーリスト
 * @member {string} Owner フォルダ作成者のユーザID
 * @member {string} Parent 親フォルダへの参照キー
 */
type Folder struct {
	Type string
	Title string
	Children []string
	Owner string
	Parent string
}

/**
 * フィード
 * @class
 * @member {string} Title フィードのタイトル
 * @member {[]string} Entries エントリのキーリスト
 * @member {string} Owner 所有者のユーザID
 * @member {string} Parent 親フォルダへの参照キー
 * @member {string} Standard フィードの規格("Atom"/"RSS1.0"/"RSS2.0"のいずれか)
 * @member {string} FinalEntry 最後に取得したエントリのキー
 * @member {string} URL フィードファイルの場所
 * @member {string} SiteURL ウェブページの場所
 */
type Feed struct {
	Title string
	Entries []string
	Owner string
	Parent string
	Standard string
	FinalEntry string
	URL string
	SiteURL string
}

/**
 * エントリ
 * @class
 * @member {string} Link エントリのURL
 * @member {string} Title エントリのタイトル
 * @member {string} Owner 所有者のユーザID
 */
type Entry struct {
	Link string
	Title string
	Owner string
}

/**
 * XMLインポート用
 * フォルダまたはフィードを表す
 * @class
 * @member {string} kind "folder" または "feed"
 * @member {string} title フォルダ、フィードのタイトル
 * @member {string} children 包含するノードリスト
 * @member {string} xmlURL 配信URL(フィードのみ)
 * @member {string} htmlURL ウェブサイトのURL(フィードのみ)
 */
type Node struct {
	kind string
	title string
	children []*Node
	xmlURL string
	htmlURL string
}

/**
 * フォルダの新規登録
 * @methofOf DAO
 * @param c {appengine.Context} コンテクスト
 * @param u {*user.User} ユーザ
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
	folder.Parent = encodedParentKey
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
 * フォルダの削除
 * 中身も全て削除する
 * rootフォルダは削除不可
 * @methodOf DAO
 * @param {appengine.Context} c コンテキスト
 * @param {string} encodedKey 削除するフォルダのキーをエンコードした文字列
 */
func (this *DAO) removeFolder(c appengine.Context, encodedKey string) {
	var err error
	var key *datastore.Key
	var folder *Folder
	var childKey string
	var childType string
	var encodedParentKey string
	var parentKey *datastore.Key
	var parentFolder *Folder
	
	key, err = datastore.DecodeKey(encodedKey)
	check(c, err)
	
	folder = new(Folder)
	err = datastore.Get(c, key, folder)
	check(c, err)
	
	// 親からの参照を削除
	encodedParentKey = folder.Parent
	parentKey, err = datastore.DecodeKey(encodedParentKey)
	check(c, err)
	parentFolder = new(Folder)
	err = datastore.Get(c, parentKey, parentFolder)
	check(c, err)
	parentFolder.Children = removeItem(parentFolder.Children, encodedKey)
	_, err = datastore.Put(c, parentKey, parentFolder)
	
	// 子を削除
	for _, childKey = range folder.Children {
		childType, _ = this.getItem(c, childKey)
		if childType == "folder" {
			this.removeFolder(c, childKey)
		} else if childType == "feed" {
			this.removeFeed(c, childKey)
		}
	}
	
	err = datastore.Delete(c, key)
	check(c, err)
}

/**
 * フォルダの取得
 * @methodOf DAO
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
 * フォルダ名の変更
 * @methodOf DAO
 * @param {appengine.Context} c コンテキスト
 * @param {string} encodedKey フォルダのキー
 * @param {string} name 新しいフォルダ名
 */
func (this *DAO) renameFolder(c appengine.Context, encodedKey string, name string) {
	var key *datastore.Key
	var folder *Folder
	var err error
	
	key, err = datastore.DecodeKey(encodedKey)
	check(c, err)
	
	folder = new(Folder)
	err = datastore.Get(c, key, folder)
	check(c, err)
	
	folder.Title = name
	_, err = datastore.Put(c, key, folder)
	check(c, err)
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
		Children []string
		Type string
		Id string
		Count int
		URL string
		Standard string
		Parent string
		SiteURL string
		FinalEntry string
	}
	var item *Item
	var itemType string
	
	key, err = datastore.DecodeKey(encodedKey)
	check(c, err)
	
	item = new(Item)
	err = datastore.Get(c, key, item)
	check(c, err)
	
	if item.Type == "" {
		// 要素はFeed
		item.Count = len(item.Entries)
		itemType = "feed"
	} else {
		// 要素はフォルダ
		itemType = "folder"
		item.Count = this.getEntriesCount(c, encodedKey)
	}
	
	return itemType, item
}

/**
 * 指定されたフォルダ以下にあるエントリの総数を返す
 * @methodOf DAO
 * @param {appengine.Context} c コンテキスト
 * @param {string} folderKey エンコード済みのフォルダキー
 * @returns {int} エントリの総数
 */
func (this *DAO) getEntriesCount(c appengine.Context, folderKey string) int {
	var key *datastore.Key
	var err error
	var folder *Folder
	var sum int
	var itemType string
	var childKey string

	key, err = datastore.DecodeKey(folderKey)
	check(c, err)
	
	folder = new(Folder)
	err = datastore.Get(c, key, folder)
	check(c, err)

	sum = 0
	for _, childKey = range folder.Children {
		itemType, _ = this.getItem(c, childKey)
		
		if itemType == "folder" {
			sum = sum + this.getEntriesCount(c, childKey)
		} else if itemType == "feed" {
			sum = sum + len(this.getEntries(c, childKey))
		}
	}
	
	return sum
}

/**
 * ルートフォルダを取得
 * @methodOf DAO
 * @param {appengine.Context} c コンテキスト
 * @param {*user.User} u ユーザオブジェクト
 * @returns {string} ルートフォルダのキー
 * @returns {*Folder} ルートフォルダ
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
 * @methodOf DAO
 * @param {appengine.Context} c コンテキスト
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
 * @param {[]*Entry} entries フィードのエントリリスト
 * @param {string} to 追加先のフォルダのキー
 * @returns {string} 追加したフィードのキーをエンコードしたもの　重複していたら空文字列
 * @returnss {bool} 重複していた場合はtrue
 */
func (this *DAO) registerFeed(c appengine.Context, feed *Feed, entries []*Entry, to string) (string, bool) {
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
		feed.Parent = to
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
		
		// エントリを追加
		this.registerEntries(c, entries, encodedKey)
	}
	
	return encodedKey, duplicated
}

/**
 * フィード名を変更する
 * @methodOf DAO
 * @param {appengine.Context} c コンテキスト
 * @param {string} encodedKey エンコード済みのフィードキー
 * @param {string} name 新しい名前
 */
func (this *DAO) renameFeed(c appengine.Context, encodedKey string, name string) {
	var key *datastore.Key
	var err error
	var feed *Feed
	
	key, err = datastore.DecodeKey(encodedKey)
	check(c, err)
	
	feed = new(Feed)
	err = datastore.Get(c, key, feed)
	check(c, err)
	
	feed.Title = name
	_, err = datastore.Put(c, key, feed)
	check(c, err)
}

/**
 * フィードの削除
 * @methodOf DAO
 * @param {appengine.Context} c コンテキスト
 * @param {string} encodedKey エンコード済みのフィードキー
 */
func (this *DAO) removeFeed(c appengine.Context, encodedKey string) {
	var key *datastore.Key
	var err error
	var feed *Feed
	var parent *Folder
	var parentKey *datastore.Key
	var encodedParentKey string
	var encodedEntryKey string
	var entryKey *datastore.Key
	
	// フィードを取得
	key, err = datastore.DecodeKey(encodedKey)
	check(c, err)
	
	feed = new(Feed)
	err = datastore.Get(c, key, feed)
	
	// 親フォルダからの参照を削除
	encodedParentKey = feed.Parent
	parentKey, err = datastore.DecodeKey(encodedParentKey)
	check(c, err)
	
	parent = new(Folder)
	err = datastore.Get(c, parentKey, parent)
	check(c, err)
	
	parent.Children = removeItem(parent.Children, encodedKey)
	_, err = datastore.Put(c, parentKey, parent)
	check(c, err)
	
	// フィードに含まれるエントリを削除
	for _, encodedEntryKey = range feed.Entries {
		entryKey, err = datastore.DecodeKey(encodedEntryKey)
		check(c, err)
		err = datastore.Delete(c, entryKey)
		check(c, err)
	}
	
	// フィードを削除
	err = datastore.Delete(c, key)
	check(c, err)
}

/**
 * フォルダの既読化
 * @methodOf DAO
 * @param {appengine.Context} c コンテキスト
 * @param {string} encodedKey
 */
func (this *DAO) readFolder(c appengine.Context, encodedKey string) {
	var key *datastore.Key
	var err error
	var folder *Folder
	var childKey string
	var childType string
	
	// フォルダを取得する
	key, err = datastore.DecodeKey(encodedKey)
	check(c, err)
	
	folder = new(Folder)
	err = datastore.Get(c, key, folder)
	check(c, err)
	
	// フォルダ以下にあるすべてのフィードを既読化
	for _, childKey = range folder.Children {
		childType, _ = this.getItem(c, childKey)
		if childType == "feed" {
			this.readFeed(c, childKey)
		} else if childType == "folder" {
			this.readFolder(c, childKey)
		}
	}
}

/**
 * フィードの既読化
 * @methodOf DAO
 * @param {appengine.Context} c コンテキスト
 * @param {string} encodedKey フィードのキー
 */
func (this *DAO) readFeed(c appengine.Context, encodedKey string) {
	var entries []*Entry
	var entry *Entry
	
	entries = this.getEntries(c, encodedKey)
	for _, entry = range entries {
		this.removeEntry(c, entry.Link, encodedKey)
	}
}

/**
 * 複数のエントリをフィードに一括で新規追加する
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
	
	if len(entries) == 0 {
		return nil
	}
	
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
	}
	feed.Entries = prepend(feed.Entries, result)
	
	// 最新のエントリを保存
	feed.FinalEntry = entries[0].Link
		
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
 * @param {string} link 削除するエントリのURL
 * @param {string} feedKey エントリが登録されているフィードのキー
 */
func (this *DAO) removeEntry(c appengine.Context, link string, feedKey string) {
	var query *datastore.Query
	var iterator *datastore.Iterator
	var key *datastore.Key
	var err error
	var feed *Feed
	var encodedEntryKey string
	
	query = datastore.NewQuery("entry").Filter("Link =", link)
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
 * @methodOf DAO
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
	query = datastore.NewQuery("feed").Filter("URL =", feed.URL).Filter("Owner =", u.ID)
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

/**
 * フィードの更新
 * 指定されたフィードに新しく追加されたエントリをデータストアに追加する
 * @methodOf DAO
 * @param {appengine.Context} c コンテキスト
 * @param {string} encodedFeedKey フィードのキー
 * @param {chan bool} 処理が完了したことを報告するチャネル　フォルダ更新から呼び出された場合に使用する
 * @returns {[]*Entry} 追加したエントリ一覧
 */
func (this *DAO) updateFeed(c appengine.Context, encodedFeedKey string, parentChannel chan bool) []*Entry {
	var feed *Feed
	var err error
	var feedKey *datastore.Key
	var currentEntries []*Entry
	var newEntries []*Entry
	var xml []byte
	var i int
	
	// フィードの取得
	feedKey, err = datastore.DecodeKey(encodedFeedKey)
	check(c, err)
	feed = new(Feed)
	err = datastore.Get(c, feedKey, feed)
	check(c, err)
	
	// URLからエントリをフェッチする
	xml = getXML(c, feed.URL)
	currentEntries = make([]*Entry, 0)
	switch feed.Standard {
		case "Atom":
			var atom *Atom
			atom = new(Atom)
			_, currentEntries = atom.encode(c, xml)

		case "RSS2.0":
			var rss2 *RSS2
			rss2 = new(RSS2)
			_, currentEntries = rss2.encode(c, xml)
			
		case "RSS1.0":
			var rss1 *RSS1
			rss1 = new(RSS1)
			_, currentEntries = rss1.encode(c, xml)
	}
	
	// エントリ一覧から最新エントリと同じURLを探す
	newEntries = make([]*Entry, 0)
	for i = 0; i < len(currentEntries); i++ {
		if currentEntries[i].Link == feed.FinalEntry {
			break
		}
		newEntries = append(newEntries, currentEntries[i])
	}
	this.registerEntries(c, newEntries, encodedFeedKey)
	
	if parentChannel != nil {
		parentChannel <- true
	}
	
	return newEntries
}

/**
 * フォルダの更新
 * @methodOf DAO
 * @param {appengine.Context} c コンテキスト
 * @param {string} FolderKey フォルダのキー
 * @param {chan bool} parentChannel アップデートが完了したことを報告するチャネル　マルチスレッドで使う
 * @returns {map[string]int} 更新後の各フォルダ、フィードのエントリ件数
 */
func (this *DAO) updateFolder(c appengine.Context, folderKey string, parentChannel chan bool) map[string]int {
	var folder *Folder
	var childKey string
	var childType string
	var result map[string]int
	var feed *Feed
	var childrenChannel chan bool
	var i int
	
	folder = this.getFolder(c, folderKey)
		
	// 新規エントリをマルチスレッドで一斉に取得・追加する
	// 各URLフェッチに時間がかかるため
	childrenChannel = make(chan bool)
	for _, childKey = range folder.Children {
		childType, _ = this.getItem(c, childKey)
		if childType == "folder" {
			go this.updateFolder(c, childKey, childrenChannel)
		} else if childType == "feed" {
			go this.updateFeed(c, childKey, childrenChannel)
		}
	}
	
	// すべてのスレッドが完了するまで待機
	for i = 0; i < len(folder.Children); i++ {
		<- childrenChannel
	}
	
	// 呼び出し元直下のエントリの数をカウントして返す
	result = make(map[string]int)
	if parentChannel != nil {
		parentChannel <- true
	} else {
		for _, childKey = range folder.Children {
			childType, _ = this.getItem(c, childKey)
			if childType == "folder" {
				result[childKey] = this.getEntriesCount(c, childKey)
			} else if childType == "feed" {
				feed = this.getFeed(c, childKey)
				result[childKey] = len(feed.Entries)
			}
		}
	}
	
	return result
}

/**
 * XMLファイルを解析してフォルダ・フィードツリーを返す
 * @methodOf DAO
 * @param {appengine.Context} c コンテキスト
 * @param {[]byte} xmldata XMLデータ
 * @returns {[]interface{}} フォルダ・フィードツリー
 */
func (this *DAO) getTreeFromXML(c appengine.Context, xmldata []byte) []*Node {
	type OUTLINE struct {
		Outline []OUTLINE `xml:"outline"`
		Title string `xml:"title,attr"`
		XMLURL string `xml:"xmlUrl,attr"`
		HTMLURL string `xml:"htmlUrl,attr"`
	}
	type OPML struct {
		Outline []OUTLINE `xml:"body>outline"`
	}
	var opml *OPML
	var err error
	var depth1 OUTLINE
	var depth2 OUTLINE
	var i int
	var j int
	var tree []*Node
	var folder *Node
	var feed *Node
	
	opml = new(OPML)
	err = xml.Unmarshal(xmldata, opml)
	check(c, err)
	
	tree = make([]*Node, len(opml.Outline))
	for i, depth1 = range opml.Outline {
		if depth1.XMLURL == "" {
			folder = new(Node)
			folder.kind = "folder"
			folder.title = depth1.Title
			folder.children = make([]*Node, len(depth1.Outline))
			for j, depth2 = range depth1.Outline {
				feed = new(Node)
				feed.kind = "feed"
				feed.title = depth2.Title
				feed.xmlURL = depth2.XMLURL
				feed.htmlURL = depth2.HTMLURL
				folder.children[j] = feed
			}
			tree[i] = folder
		} else {
			feed = new(Node)
			feed.kind = "feed"
			feed.title = depth1.Title
			feed.xmlURL = depth1.XMLURL
			feed.htmlURL = depth1.HTMLURL
			tree[i] = feed
		}
	}

	return tree
}

/**
 * XMLファイルをデータストアにインポートする
 * @methodOf DAO
 * @param {appengine.Context} c コンテキスト
 * @param {[]byte} xmldate XMLファイル
 * @param {string} folderKey 追加先のフォルダのキー
 */
func (this *DAO) importXML(c appengine.Context, tree []*Node, folderKey string) {
	var feed *Feed
	var u *user.User
	var parentKey string
	var entries []*Entry
	var depth1 *Node
	var depth2 *Node
	
	u = user.Current(c)
	
	for _, depth1 = range tree {
		if depth1.kind == "folder" {
			parentKey = this.registerFolder(c, u, depth1.title, false, folderKey)
			for _, depth2 = range depth1.children {
				feed = new(Feed)
				entries = make([]*Entry, 0)
				feed, entries = this.getFeedFromXML(c, depth2.xmlURL)
				feed.URL = depth2.xmlURL
				feed.SiteURL = depth2.htmlURL
				this.registerFeed(c, feed, entries, parentKey)
			}
		} else {
			feed = new(Feed)
			entries = make([]*Entry, 0)
			feed, entries = this.getFeedFromXML(c, depth1.xmlURL)
			feed.URL = depth1.xmlURL
			feed.SiteURL = depth1.htmlURL
			this.registerFeed(c, feed, entries, parentKey)
		}
	}
}

/**
 * XMLのURLからフィードを取得する
 * @methodOf DAO
 * @param {appengine.Context} c コンテキスト
 * @param {string} url XMLファイルの場所
 * @returns {*Feed} フィード
 * @returns {[]*Entries} フィードのエントリリスト
 */
func (this *DAO) getFeedFromXML(c appengine.Context, url string) (*Feed, []*Entry) {
	var feedXML []byte
	var feedType string
	var feed *Feed
	var entries []*Entry
	
	feedXML = getXML(c, url)
	
	feed = new(Feed)
	entries = make([]*Entry, 0)
	feedType = this.getType(c, feedXML)
	switch feedType {
		case "Atom":
			var atom *Atom
			atom = new(Atom)
			feed, entries = atom.encode(c, feedXML)
		case "RSS2.0":
			var rss2 *RSS2
			rss2 = new(RSS2)
			feed, entries = rss2.encode(c, feedXML)
		case "RSS1.0":
			var rss1 *RSS1
			rss1 = new(RSS1)
			feed, entries = rss1.encode(c, feedXML)
		case "etc":
	}
	feed.URL = url
	
	return feed, entries
}

/**
 * XMLデータの規格を判断する
 * @methodOf DAO
 * @param {appengine.Context} c コンテキスト
 * @param {[]byte} bytes XMLデータ
 * @returns {string} フィードの規格(RSS1.0 / RSS2.0 / Atom / etc)
 */
func (this *DAO) getType(c appengine.Context, bytes []byte) string {
	type Checker struct {
		XMLName xml.Name
	}
	var checker *Checker
	var err error
	var result string
	
	checker = new(Checker)
	err = xml.Unmarshal(bytes, checker)
	check(c, err)
	
	switch checker.XMLName.Local {
		case "feed":
			result = "Atom"
		case "rss":
			result = "RSS2.0"
		case "RDF":
			result = "RSS1.0"
		default:
			result = "etc"
	}
	
	return result
}

/**
 * XMLデータをデータストアに保存する
 * @methodOf DAO
 * @param {appengine.Context} c コンテキスト
 * @param {[]byte} xml XMLデータ
 */
func (this *DAO) saveXML(c appengine.Context, xml []byte) {
	var u *user.User
	var key *datastore.Key
	var err error
	type Entity struct {
		XML []byte
	}
	var entity *Entity
	
	u = user.Current(c)
	key = datastore.NewKey(c, "xml", u.ID, 0, nil)
	
	entity = new(Entity)
	entity.XML = xml
	_, err = datastore.Put(c, key, entity)
	check(c, err)
}

/**
 * XMLデータをデータストアから取得する
 * @methodOf DAO
 * @param {appengine.Context} c コンテキスト
 * @returns {[]byte} ロードしたXMLデータ
 */
func (this *DAO) loadXML(c appengine.Context) []byte {
	var u *user.User
	var key *datastore.Key
	var err error
	type Entity struct {
		XML []byte
	}
	var entity *Entity
	
	u = user.Current(c)
	key = datastore.NewKey(c, "xml", u.ID, 0, nil)
	entity = new(Entity)
	err = datastore.Get(c, key, entity)
	check(c, err)
	
	return entity.XML
}

/**
 * すべてのフォルダをアップデートする
 * @methodOf DAO
 * @param {appengine.Context} c コンテキスト
 */
func (this *DAO) updateAll(c appengine.Context) {
	var query *datastore.Query
	var iterator *datastore.Iterator
	var key *datastore.Key
	var err error
	var count int
	var i int
	
	query = datastore.NewQuery("folder").Filter("Type =", "root").KeysOnly()
	count, err = query.Count(c)
	check(c, err)
	iterator = query.Run(c)
	for i = 0; i < count; i++ {
		key, err = iterator.Next(nil)
		check(c, err)
		this.updateFolder(c, key.Encode(), nil)
	}
	log.Printf("update all folder")
}