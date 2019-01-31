# これは何？
iOS, Android, Webにプッシュ通知を送信するサーバーです。

各サービスのユーザーID毎に各プラットフォームに良い感じに負荷分散しながら送信してくれます。

プッシュ通知送信機能を導入したいプロジェクトに個別でデプロイして使用してください。

# 対応状況

## データベース
- Cloud Datastore
- Cloud Firestore
- Cloud SQL(MySQL) ToBe...

## 機能
- ユーザー＆Token登録
- 即時送信
- 予約送信 ToBe...
- 定期送信 ToBe...

# セットアップ

## 準備
```bash
cp env.example.mk env.mk
cp appengine/env/credentials_local.example.json appengine/env/credentials_local.json
cp appengine/env/credentials_staging.example.json appengine/env/credentials_staging.json
cp appengine/env/credentials_production.example.json appengine/env/credentials_production.json
cp appengine/env/values_local.example.yaml appengine/env/values_local.yaml
cp appengine/env/values_staging.example.yaml appengine/env/values_staging.yaml
cp appengine/env/values_production.example.yaml appengine/env/values_production.yaml

dep ensure
```

## 実行
```bash
make run app=push
```

## デプロイ
```bash
make deploy app=push
make deploy-production app=push
```

# API
JSONRPC2.0を使用しています。
詳細はdoc内を参照してください。
