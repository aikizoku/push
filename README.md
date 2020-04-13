# これは何？
iOS, Android, Webにプッシュ通知を送信するサーバーです。

各サービスのユーザーID毎に各プラットフォームに良い感じに負荷分散しながら送信してくれます。

プッシュ通知送信機能を導入したいプロジェクトに個別でデプロイして使用してください。

# 対応状況

## 機能
- ユーザー＆Token登録
- 即時送信
- 予約送信
- 定期送信 ToBe...

# セットアップ

## 準備
```bash
cd appengine/push
go get -u ...
go mod tidy
```

## 実行
```bash
cd appengine/push
make run
```

## デプロイ
```bash
cd appengine/push
make deploy
make deploy-prod
```

# API
JSONRPC2.0を使用しています。
詳細はrundoc/docs内を参照してください。
