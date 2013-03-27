okareader
=========

PC,スマホ対応RSSリーダーのWebアプリ  
MitLicense  

現在開発中です  

okareader/  
├── README.md  このファイル  
├── app.yaml   アプリの設定ファイル  
├── client     クライアント(htmlファイル)が使うファイル  
│   └── home.js
│  
└── server  サーバが使用するファイル　クライアントから見えない
    ├── controller.go  リクエストの振り分け
    ├── model.go       データモデル
    ├── view.go        画面表示関係
    ├── test.go        テスト  
    ├── html           HTMLテンプレート Viewが使う
    │   ├── feed.html
    │   ├── folder.html
    │   ├── home.html
    │   └── login.html
    └─── library
	    ├── atom.go
        └── lib.go
  
  
client/ は静的ディレクトリとして設定しているため、  
goからアクセスすることができないので注意。  
基本的にMVCパターンになっています。  