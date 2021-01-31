# headless-service

## 概要

クラスター内部からDeploymentのPodに疎通したい場合、ClusterIPを利用してプロキシする方法が一般的です。
しかしながら、StatefulSetを利用する場合など、Podを特定して疎通することは出来ません。

そのような場合、Headless Serviceを利用することでPodを特定して疎通することが出来ます。
実際には
```
curl sample-statefulset-headless-0
curl sample-statefulset-headless-1
curl sample-statefulset-headless-2
```
のように、Pod名で直接名前解決が出来るようになります。

### マニフェストの書き方

```
---
apiVersion: v1
kind: Service
metadata:
  name: sample-headless
spec:
  type: ClusterIP
  clusterIP: None
  ports:
  - name: "http-port"
    protocol: "TCP"
    port: 80
    targetPort: 80
  selector:
    app: sample-app
---
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: sample-statefulset-headless
spec:
  serviceName: sample-headless #HeadlessServiceの名前を指定
  replicas: 3
  selector:
    matchLabels:
      app: sample-app
  template:
    metadata:
      labels:
        app: sample-app
    spec:
      containers: #省略
```
HeadlessServiceをStatefulSetで利用する場合、
- `spec.Type: ClusterIP`にする
- `spec.clusterIP: None`
- StatefulSetの`spec.serviceName`はHeadlessServiceの`metadata.name`と一致させる
必要があります。

## 動作を確認

HeadlessServiceとStatefulSetを作成します。

```
> kubectl apply -f headless-service.yaml
namespace/handson created
service/sample-headless created
statefulset.apps/sample-statefulset-headless created

> kubectl get po -o wide --namespace=handson
NAME                            READY   STATUS    RESTARTS   AGE   IP            NODE           NOMINATED NODE   READINESS GATES
sample-statefulset-headless-0   1/1     Running   0          5s    10.244.2.15   kind-worker2   <none>           <none>
sample-statefulset-headless-1   1/1     Running   0          3s    10.244.3.13   kind-worker3   <none>           <none>
sample-statefulset-headless-2   1/1     Running   0          2s    10.244.1.13   kind-worker    <none>           <none>
```

Podが3台それぞれ「sample-statefulset-headless-*」という名前で作成されていることが確認出来ます。
早速クラスタ内からcurlを投げてみましょう。

```
> kubectl exec -it statefulset/sample-statefulset-headless --namespace=handson -- curl sample-statefulset-headless-0.sample-headless/podinfo
{"pod_ip":"10.244.2.15", "node_ip": "192.168.160.4", "client_ip": "10.244.2.15"}

> kubectl exec -it statefulset/sample-statefulset-headless --namespace=handson -- curl sample-statefulset-headless-1.sample-headless/podinfo
{"pod_ip":"10.244.3.13", "node_ip": "192.168.160.5", "client_ip": "10.244.2.15"}

> kubectl exec -it statefulset/sample-statefulset-headless --namespace=handson -- curl sample-statefulset-headless-2.sample-headless/podinfo
{"pod_ip":"10.244.1.13", "node_ip": "192.168.160.2", "client_ip": "10.244.2.15"}
```

HeadlessServiceを作成することで、`sample-statefulset-headless-0.sample-headless`を利用して特定のPodに疎通出来ていることが分かります。

## 最後に

リソースを削除しておきましょう。
```
kubectl delete namespace handson
```