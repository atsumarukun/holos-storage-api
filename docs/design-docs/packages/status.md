# 概要

テータスパッケージを実装する.

# 対象範囲

## 達成基準

- テータスを呼び出せる状態にする

## 除外項目

- Interface層でのレスポンスマッピングは考慮しない

# 利用方法

ステータスの初期化は以下を利用する.

```golang
st := status.New(code, message)
```

エラーステータスの初期化は以下を利用する.

```golang
st := status.Error(code, message)
```

エラーからのステータス初期化は以下を利用する.

```golang
st := status.FromError(err)
```

# 詳細設計

## ステータスコード

| コード | 説明 |
| --- | --- |
| BadRequest | 不正なリクエスト |
| Unauthorized | 認証失敗 |
| Forbidden | 認可失敗 |
| Conflict | リソースの重複 |
| Internal | サーバーの内部エラー |

エラーからステータスを初期化する際のコードは`Internal`にする.

# その他の手法

# 参考文献

# 変更履歴

| 変更日 | 変更者 | 変更内容 |
| --- | --- | --- |
| 2025/04/13 | @atsumarukun | 初版 |
