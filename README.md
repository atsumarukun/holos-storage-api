# ストレージAPI

# アカウントAPI

## 開発環境

下記コマンドを用いて.envの作成とDockerNetworkのを行う.<br />
既にDockerNetworkが作成されている場合DockerNetWorkの作成は不要.
```bash
cp .env.example .env
docker network create nw-holos
```

その後[DevContainerを用いたコンテナの立ち上げ](https://code.visualstudio.com/docs/devcontainers/create-dev-container#:~:text=With%20the%20above%20devcontainer.json,Reopen%20in%20Container%20command)を行う.<br />
DevContainerを利用しない場合は下記コマンドでコンテナを立ち上げる.
```bash
docker compose up -d
docker compose exec storage-api air -c .air.toml
```

開発環境の詳細は以下のDesignDocに記載.<br />
[development-environment.md](https://github.com/atsumarukun/holos-storage-api/blob/main/docs/design-docs/development-environment/development-environment.md)

## ルール

### ブランチ名

Issue番号を用いたブランチ名にする.

```
issue-${ISSUE_NUMBER}   // issue-1
```

### コミットメッセージ

以下の表に記載するタグを利用する.

| タグ | 説明 |
| --- | --- |
| create | 機能作成 |
| update | 機能更新 |
| remove | 機能削除 |
| refactor | リファクタリング |
| fix | 不具合修正 |

```
${TAG}: ${MESSAGE}   // create: ファイル削除機能を作成.
```
