okareader
=========

![スクリーンショット](https://raw.github.com/yokano/okareader/master/okareader_ss.png)

PC,スマホ対応RSSリーダーのWebアプリ  
Google App Engine + Go で開発しました
MitLicense です
Google Reader の xml ファイルをインポートすることも出来ます
ユーザ認証に Google アカウントを使用しています

	okareader/
	├── LICENSE.txt
	├── README.md
	├── app.yaml
	├── client
	│   ├── feed.js
	│   ├── feed.png
	│   ├── folder.js
	│   ├── folder.png
	│   ├── import.js
	│   ├── okareader.css
	│   └── okareader.png
	├── cron.yaml
	└── server
		├── atom.go
		├── controller.go
		├── html
		│   ├── feed.html
		│   ├── folder.html
		│   ├── import.html
		│   └── login.html
		├── lib.go
		├── main.go
		├── model.go
		├── rss1.go
		├── rss2.go
		└── view.go

基本的にMVCパターンになっています。  
外部ライブラリとして jQueryMobile を使用しています。

## 設定ファイル
* app.yaml　　アプリの設定
* cron.yaml　　フィードの定期的な自動更新の設定

## client/
このディレクトリはstatic_dirとして設定されています  
htmlファイルから参照するスクリプトや画像を保存するディレクトリです  
goから直接このファイルを参照することはできません

## server/
* main.go　　controllerの呼び出し
* controller.go　　クライアントからのリクエストをViewやModelに振り分けながら処理する
* model.go　　データ操作全般を行う
* view.go　　画面表示全般を行う
* rss1.go　　RSS1.0を読み込むための処理
* rss2.go　　RSS2.0を読み込むための処理
* atom.go　　Atomを読み込むための処理
* lib.go　　その他の汎用的な関数

## 連絡先
yuta.okano@gmail.com