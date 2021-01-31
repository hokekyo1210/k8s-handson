# cloud-functions-go

ローカル環境でgo言語のcloud functions開発を加速させるためのボイラープレートです.

\[Use this template\]ボタンから自分のワークスペースにコピーして使ってください.

## できること

- go言語cloudfunctionsのローカル開発
- 環境変数の設定(docker-composeによる)
- ホットリロード([air](https://github.com/cosmtrek/air))

## 使い方

### Hello World

ローカルで動作を確認する.

```
cloud-functions-go % docker-compose up -d
cloud-functions-go % curl localhost:1323/time
Hello World!
```

cloud functionsにデプロイする.

```
gcloud functions deploy timeutc \
    --entry-point TimeUTC \
    --runtime go113 \
    --max-instances=3 \
    --trigger-http --region=asia-northeast1
curl https://asia-northeast1-*******.cloudfunctions.net/timeutc
Hello World!
```

cloud runにデプロイする.

gcloud run deploy time-utc \
            --project be-tsuchida-yuki1 \
            --image asia.gcr.io/be-tsuchida-yuki1/time:utc \
            --platform managed \
            --region asia-northeast1 \
            --allow-unauthenticated

gcloud run deploy time-jst \
            --project be-tsuchida-yuki1 \
            --image asia.gcr.io/be-tsuchida-yuki1/time:jst \
            --platform managed \
            --region asia-northeast1 \
            --allow-unauthenticated