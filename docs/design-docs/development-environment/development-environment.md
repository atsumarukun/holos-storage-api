# 概要

開発環境の構築を行う.

# 対象範囲

## 達成基準

- 開発用コンテナが立ち上がっている状態
- ローカルでLintやTestの実行ができる状態

## 除外項目

- 本番デプロイ用の環境構築は行わない
- CIによるLintやTestの実行は行わない

# 利用方法

開発環境には[DevContainer](https://code.visualstudio.com/docs/devcontainers/containers)を利用する.

Visual Studio Code拡張機能:<br />
https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers

## コンテナ一覧

| サービス | イメージ | ポート |
| --- | --- | --- |
| account-api | golang:1.24 | 8001:8000 |
| account-db | mysql:9.2 | |

## ネットワーク

`nw-holos`を利用する.

ネットワークが作成されていない場合は以下コマンドで作成を行う.
```bash
docker network create nw-holos
```

# 詳細設計

## VSCode拡張機能

### Golang

Golang:<br />
https://marketplace.visualstudio.com/items?itemName=golang.Go

### Swagger

Swagger Editor:<br />
https://marketplace.visualstudio.com/items?itemName=42Crunch.vscode-openapi

### Markdown

Mermaid:<br />
https://marketplace.visualstudio.com/items?itemName=bierner.markdown-mermaid

### Code Formatter

SQL Formatter:<br />
https://marketplace.visualstudio.com/items?itemName=ReneSaarsoo.sql-formatter-vsc

### Code Helper

Tailing Spaces:<br />
https://marketplace.visualstudio.com/items?itemName=shardulm94.trailing-spaces

## パッケージ

### ホットリロード

Air:<br />
https://github.com/air-verse/air

### フレームワーク

Gin:<br />
https://gin-gonic.com/ja/

### ORM

sqlx:<br />
https://github.com/jmoiron/sqlx

### マイグレーション

golang-migrate:<br />
https://github.com/golang-migrate/migrate

### モック生成

uber-go/mock:<br />
https://github.com/uber-go/mock

### データベースモック

go-sqlmock:<br />
https://github.com/DATA-DOG/go-sqlmock

### Linter

golangci-lint:<br />
https://github.com/golangci/golangci-lint

# その他の手法

# 参考文献

# 変更履歴

| 変更日 | 変更者 | 変更内容 |
| --- | --- | --- |
| 2025/04/08 | @atsumarukun | 初版 |
