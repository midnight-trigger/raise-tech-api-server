# 使用技術
* バックエンド：Golang1.14（Echo）
* ORM：GORM
* DB：MySQL8.0.22

# ディレクトリ構成
```
├── api
│   ├── controller
│   ├── definition
│   ├── domain
│   └── error_handing
├── configs
├── dev_tools
├── infra
├── logger
└── third_party
```
* `api` : API関係
   * controller： API Controller
   * definition： API リクエスト/レスポンス構造体定義
   * domain： API ビジネスロジック
   * error_handing： API エラー処理
   * route.go, route_gen/go: ルーティング関連
* `configs` : 環境毎（local, dev, qa, stg prod etc...）の定義ファイルなど
* `dev_tools` : 開発効率化ためのツールなど
* `infra` : DB、Cognito関連の定義と処理
* `logger` : log処理
* `third_party` : 使用するライブラリ関係の処理など

# マイグレーション方法
 https://github.com/golang-migrate/migrate

 golang-migrate/migrateを利用する

```
# Macでbrewをお使いの場合
brew install golang-migrate
```

マイグレーションコマンド
```
migrate -database 'mysql://root:password@tcp(0.0.0.0:3306)/<DB名>' -path migrations up
```
※ golang-migrate/migrateでまだ管理されていないテーブルに流し込む場合は、いったんDBを空にしてから上記コマンドを実行する

 schema_migrationsテーブルが作成されてそこで実行管理される。

 マイグレーションファイルの作成
 ```
 migrate create -ext sql -dir migrations/ -seq <ファイル名>
 ```
