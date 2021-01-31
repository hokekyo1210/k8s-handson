# topology-aware-service-routing

## 概要

マルチリージョンやマルチゾーンクラスタにおいてClusterIPを利用すると、物理的に遠いPodにトラフィックが流れてしまいパフォーマンスが低下する場合があります。
このような場合`spec.TopologyKeys`オプションを利用することで、同一ノードのPodを優先、同一ゾーンを優先といった優先度付きのルーティングが可能になります。

実際には以下のような設定を行います。

```
spec:
  type: ClusterIP
  selector:
    app: sample-app
  ports:
  - name: http
    protocol: TCP
    port: 8080
    targetPort: 80
  topologyKeys:
  - kubernetes.io/hostname #同一ノード
  - kubernetes.io/zone #同一ゾーン
  - "*"
```

topologyKeysの中で上にかかれているものほど優先度が高くなります。
1. 同一ノード内で有効なPodを探す
2. なければ、同一ゾーン内で有効なPodを探す
3. なければ、クラスタ全体から有効なPodを探す
といったルーティングを行います。

なお、`kubernetes.io/zone`, `kubernetes.io/region`はクラウドプロバイダによって自動的に付与されるラベルなので、kindでは利用出来ません。

## topologyKeysを実際に試してみる

ClusterIPに利用する設定は以下です。
```
  topologyKeys:
    - kubernetes.io/hostname #同一ノード優先
    - "*"
```
それでは、同一ノード内にPodがある場合と無い場合でClusterIPの実際の動きを確認してみます。

### 同一ノードにPodがある場合

ClusterIP, podのIPを返すnginx(6台), bastion(1台)を追加します。

```
> kubectl apply -f topology-aware-service-routing-same-node.yaml
namespace/handson created
service/sample-clusterip created
deployment.apps/sample-deployment created
deployment.apps/bastion created

> kubectl get po -o wide --namespace=handson
NAME                       READY   STATUS    RESTARTS   AGE   IP            NODE           NOMINATED NODE   READINESS GATES
bastion-588db45795-r2r99   1/1     Running   0          15s   10.244.1.9    kind-worker    <none>           <none>
nginx-856d8456bc-58k4b     1/1     Running   0          15s   10.244.2.9    kind-worker2   <none>           <none>
nginx-856d8456bc-jgj65     1/1     Running   0          15s   10.244.3.7    kind-worker3   <none>           <none>
nginx-856d8456bc-lhbjw     1/1     Running   0          15s   10.244.1.10   kind-worker    <none>           <none>
nginx-856d8456bc-x4qjz     1/1     Running   0          15s   10.244.2.10   kind-worker2   <none>           <none>
nginx-856d8456bc-x62gf     1/1     Running   0          15s   10.244.1.8    kind-worker    <none>           <none>
nginx-856d8456bc-x7b5s     1/1     Running   0          15s   10.244.3.8    kind-worker3   <none>           <none>
```

bastionがkind-workerノードに配置され、nginxが全てのノードにまんべんなく配置されています。

bastionからClusterIPに何回かcurlを飛ばしてみましょう。

```
> kubectl exec -it deployment/bastion --namespace=handson -- curl sample-clusterip/podinfo
{"pod_ip":"10.244.1.10", "node_ip": "192.168.160.2", "client_ip": "10.244.1.9"}

> kubectl exec -it deployment/bastion --namespace=handson -- curl sample-clusterip/podinfo
{"pod_ip":"10.244.1.8", "node_ip": "192.168.160.2", "client_ip": "10.244.1.9"} 
```

bastionと同一ノード(kind-worker)のPodだけにリクエストが送信されていることが分かります。

### 同一ノードにPodが無い場合

Podの配置を変更するために、さきほど作成したリソースを削除します。
```
> kubectl delete namespace handson
```

再度リソースを追加します。(先ほどとはマニフェストが若干違います)
```
> kubectl apply -f topology-aware-service-routing-other-node.yaml 
namespace/handson created
service/sample-clusterip created
deployment.apps/nginx created
deployment.apps/bastion created

> kubectl get po -o wide --namespace=handson                                              
NAME                       READY   STATUS    RESTARTS   AGE   IP            NODE           NOMINATED NODE   READINESS GATES
bastion-588db45795-vn9sh   1/1     Running   0          8s    10.244.1.11   kind-worker    <none>           <none>
nginx-7559fbc8d-7lk2g      1/1     Running   0          8s    10.244.2.12   kind-worker2   <none>           <none>
nginx-7559fbc8d-8wf7v      1/1     Running   0          8s    10.244.3.9    kind-worker3   <none>           <none>
nginx-7559fbc8d-9h2kq      1/1     Running   0          8s    10.244.3.11   kind-worker3   <none>           <none>
nginx-7559fbc8d-ksbs8      1/1     Running   0          8s    10.244.2.13   kind-worker2   <none>           <none>
nginx-7559fbc8d-pnlzx      1/1     Running   0          8s    10.244.2.11   kind-worker2   <none>           <none>
nginx-7559fbc8d-sgrtf      1/1     Running   0          8s    10.244.3.10   kind-worker3   <none>           <none>
```

bastionがkind-workerノードに配置され、nginxがkind-worker以外のノードに配置されています。

bastionからClusterIPに何回かcurlを飛ばしてみましょう。

```
> kubectl exec -it deployment/bastion --namespace=handson -- curl sample-clusterip/podinfo
{"pod_ip":"10.244.2.12", "node_ip": "192.168.160.4", "client_ip": "10.244.1.11"}

> kubectl exec -it deployment/bastion --namespace=handson -- curl sample-clusterip/podinfo
{"pod_ip":"10.244.3.9", "node_ip": "192.168.160.5", "client_ip": "10.244.1.11"}
```

今回はbastionと同一ノード(kind-worker)に宛先Podが存在しないため、topologyKeysの2番目の設定が反映された結果、全てのPodにリクエストが送られるようになりました。

## 最後に

リソースを削除しておきましょう。
```
kubectl delete namespace handson
```