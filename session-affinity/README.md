# session-affinity

## 概要

- ClusterIP Service
- LoadBalancer Service
- NodePort Service

では、`spec.sessionAffinity`オプションを設定することで、セッションアフィニティを有効化できます(デフォルトでは無効)。

セッションアフィニティとは送信元のIP(ClientIP)によって、送信先のPodを固定する機能です。

注意点として、NodePort Serviceにおいては、どのNodeに転送するかによって、同じクライアントIPアドレスでも同じPodに転送されるとは限らないことに気をつけましょう。


## ClusterIP ServiceにおけるSessionAffinity

### SessionAffinity: None

まずはじめに、SessionAffinityが有効化されていないClusterIPの動作確認をしてみます。

Namespace, Service, Deploymentリソースを作成します。

```
> kubectl apply -f no-session-affinity-cluster-ip.yaml                    
namespace/handson unchanged
service/sample-clusterip configured
deployment.apps/sample-deployment configured

> kubectl get pods --namespace=handson
```

Pod内に入り、ClusterIP Serviceに対してcurlを投げてみます。

```
> kubectl exec -it deployment/sample-deployment --namespace=handson -- bash

> curl sample-clusterip/podinfo
{"pod_ip":"10.244.3.5", "node_ip": "192.168.160.3"}root@sample-deployment-b68856595-hs2h4:/# curl sample-clusterip/podinfo
{"pod_ip":"10.244.3.5", "node_ip": "192.168.160.3"}root@sample-deployment-b68856595-hs2h4:/# curl sample-clusterip/podinfo
{"pod_ip":"10.244.2.7", "node_ip": "192.168.160.5"}root@sample-deployment-b68856595-hs2h4:/# curl sample-clusterip/podinfo
{"pod_ip":"10.244.2.7", "node_ip": "192.168.160.5"}root@sample-deployment-b68856595-hs2h4:/# curl sample-clusterip/podinfo
{"pod_ip":"10.244.3.5", "node_ip": "192.168.160.3"}root@sample-deployment-b68856595-hs2h4:/# curl sample-clusterip/podinfo
{"pod_ip":"10.244.3.5", "node_ip": "192.168.160.3"}root@sample-deployment-b68856595-hs2h4:/# curl sample-clusterip/podinfo
{"pod_ip":"10.244.1.6", "node_ip": "192.168.160.2"}root@sample-deployment-b68856595-hs2h4:/# curl sample-clusterip/podinfo
{"pod_ip":"10.244.2.7", "node_ip": "192.168.160.5"}root@sample-deployment-b68856595-hs2h4:/# 

> exit
```

何度かhttpリクエストを送ると, 宛先のpodが固定されていないことが分かります。

### SessionAffinity: ClientIP

つぎに、SessionAffinityが有効化されているClusterIPの動作を確認してみます。

ここで使用する設定は以下です。
```
spec:
  type: ClusterIP
  ports:
    - name: "http-port"
      protocol: "TCP"
      port: 80
      targetPort: 80
  sessionAffinity: ClientIP
  sessionAffinityConfig:
    clientIP:
      timeoutSeconds: 5
  selector:
    app: sample-app
```
- 送信元IP(ClientIP)ベースで宛先を固定
- タイムアウト5秒

リソースを更新し、再度curlを投げてみます。

```
> kubectl apply -f session-affinity-cluster-ip.yaml 
namespace/handson unchanged
service/sample-clusterip configured
deployment.apps/sample-deployment unchanged

> kubectl exec -it deployment/sample-deployment --namespace=handson -- bash

> curl sample-clusterip/podinfo
{"pod_ip":"10.244.1.6", "node_ip": "192.168.160.2"}root@sample-deployment-b68856595-hs2h4:/# curl sample-clusterip/podinfo
{"pod_ip":"10.244.1.6", "node_ip": "192.168.160.2"}root@sample-deployment-b68856595-hs2h4:/# curl sample-clusterip/podinfo
{"pod_ip":"10.244.1.6", "node_ip": "192.168.160.2"}root@sample-deployment-b68856595-hs2h4:/# curl sample-clusterip/podinfo
{"pod_ip":"10.244.1.6", "node_ip": "192.168.160.2"}root@sample-deployment-b68856595-hs2h4:/# curl sample-clusterip/podinfo
{"pod_ip":"10.244.1.6", "node_ip": "192.168.160.2"}root@sample-deployment-b68856595-hs2h4:/# curl sample-clusterip/podinfo
{"pod_ip":"10.244.2.7", "node_ip": "192.168.160.5"}root@sample-deployment-b68856595-hs2h4:/# curl sample-clusterip/podinfo
{"pod_ip":"10.244.2.7", "node_ip": "192.168.160.5"}root@sample-deployment-b68856595-hs2h4:/# curl sample-clusterip/podinfo
{"pod_ip":"10.244.2.7", "node_ip": "192.168.160.5"}root@sample-deployment-b68856595-hs2h4:/# curl sample-clusterip/podinfo
{"pod_ip":"10.244.2.7", "node_ip": "192.168.160.5"}root@sample-deployment-b68856595-hs2h4:/# curl sample-clusterip/podinfo
{"pod_ip":"10.244.2.7", "node_ip": "192.168.160.5"}root@sample-deployment-b68856595-hs2h4:/# curl sample-clusterip/podinfo
{"pod_ip":"10.244.3.5", "node_ip": "192.168.160.3"}root@sample-deployment-b68856595-hs2h4:/# curl sample-clusterip/podinfo
{"pod_ip":"10.244.3.5", "node_ip": "192.168.160.3"}root@sample-deployment-b68856595-hs2h4:/# 

> exit
```

clusteripに対して連続的にhttpリクエストを送ると、宛先のpodが固定されていることが分かります。
また、リクエストを送らずに5秒待つと、次からは宛先が別のpodになっていることも確認出来ます。

## 最後に

リソースを削除しておきましょう。
```
kubectl delete namespace handson
```