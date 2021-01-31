# ExternalName Service

## 概要

ExternalName Serviceを利用することで、外部のドメイン宛のCNAMEを返すクラスタ内DNSを作成出来ます。
これは、クラスタ内アプリケーションと外部サービスを疎結合にしたい場合などに利用します。
- 例えば、クラスタ内アプリケーションから直接外部サービスのアドレスを叩いてしまうと、外部サービスのアドレスが変更された際にアプリケーションを全て書き換えなければいけなくなってしまいますが、ExternalName Serviceを一段噛ませることで、Serviceの設定を書き換えるだけで済むようになります。

> ![ExternalName Serviceを利用した疎結合の確保](https://thinkit.co.jp/sites/default/files/article_node/1373901.jpg)
> 引用元: [KubernetesのDiscovery＆LBリソース（その2）](https://thinkit.co.jp/article/13739)

### マニフェスト例

```
apiVersion: v1
kind: Service
metadata:
  name: sample-externalname
spec:
  type: ExternalName
  externalName: external.example.com #外部ドメインのCNAME
```

## 動作の確認

ExternalNameServiceとdig用のPodを作成します。
```
> kubectl apply -f external-name-service.yaml 
namespace/handson created
service/sample-externalname created
deployment.apps/bastion created

> kubectl get pods --namespace=handson
NAME                       READY   STATUS    RESTARTS   AGE
bastion-79d4945957-2vv4f   1/1     Running   0          10s
```

ExternalNameServiceに対してdigを実行し、CNAMEを確認します。
```
> kubectl exec -it deployment/bastion --namespace=handson -- dig sample-externalname.handson.svc.cluster.local CNAME | grep external.example.com
sample-externalname.handson.svc.cluster.local. 30 IN CNAME external.example.com.
```

CNAMEレコードが内部DNSに登録されていることが確認出来ます。

(実際に外部サービスを用意してCNAMEを書き換えるところまでやりたかったけどHost Headerの問題を突破出来なかった)
> HTTPやHTTPSなどの一般的なプロトコルでExternalNameを使用する際に問題が発生する場合があります。ExternalNameを使用する場合、クラスター内のクライアントが使用するホスト名は、ExternalNameが参照する名前とは異なります。
> ホスト名を使用するプロトコルの場合、この違いによりエラーまたは予期しない応答が発生する場合があります。HTTPリクエストがオリジンサーバーが認識しないHost:ヘッダーを持っていたなら、TLSサーバーはクライアントが接続したホスト名に一致する証明書を提供できません。
> 引用元: [kubernetes公式ドキュメント Service](https://kubernetes.io/ja/docs/concepts/services-networking/service/#externalname)

## 最後に

リソースを削除しておきましょう。
```
kubectl delete namespace handson
```