# Ingress

## 概要

k8sにおいて、ServiceがL4のロードバランシングを提供するリソースであるのに対し、Ingressは**L7ロードバランシング**を提供するリソースです。
アプリケーションレベルのロードバランシングを提供することが出来るため、例えば
- httpリクエストのリクエストパスやhostヘッダに応じてルーティングするサービスを切り替える
- SSLの終端処理を行う
などの機能が利用出来ます。

また、Ingressの種類にはクラウドベンダが提供するL7LBを利用するものと、クラスタ内にIngress用のPodをデプロイするものの2種類があります。
ここでは、後者にあたるローカル環境のkindでも利用可能なNginxIngressを紹介します。

## ハンズオン

まずはじめに、NginxIngressリソースを利用するために必要な「NginxIngressController」をkindにapplyしましょう。
```
> kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/master/deploy/static/provider/kind/deploy.yaml
```
また、ここで作業の必要はありませんが、NginxIngressをkindで利用する際はクラスターのコントロールプレーンに以下の設定を入れる必要があります。
```
nodes:
  - role: control-plane
    kubeadmConfigPatches:
      - |
        kind: InitConfiguration
        nodeRegistration:
          kubeletExtraArgs:
            node-labels: "ingress-ready=true"
```

### Hostヘッダによるルーティング

NginxIngressを作成し、HostヘッダによってルーティングされるPodが変化することを確認してみましょう。

Deployment, NginxIngress, ClusterIPリソースをapplyします。
```
> kubectl apply -f ingress-host-base-routing.yaml 
namespace/handson created
ingress.networking.k8s.io/sample-ingress-by-nginx created
service/sample-clusterip-a created
service/sample-clusterip-b created
deployment.apps/sample-deployment-a created
deployment.apps/sample-deployment-b created

> kubectl get pods --namespace=handson
NAME                                   READY   STATUS    RESTARTS   AGE
sample-deployment-a-56ccb6fb74-5rr4q   1/1     Running   0          11s
sample-deployment-a-56ccb6fb74-b2xzn   1/1     Running   0          11s
sample-deployment-a-56ccb6fb74-fzlfm   1/1     Running   0          11s
sample-deployment-b-5b8f4497c8-6crwk   1/1     Running   0          11s
sample-deployment-b-5b8f4497c8-96hcr   1/1     Running   0          11s
sample-deployment-b-5b8f4497c8-lgwkf   1/1     Running   0          11s

> kubectl get ingress --namespace=handson
Warning: extensions/v1beta1 Ingress is deprecated in v1.14+, unavailable in v1.22+; use networking.k8s.io/v1 Ingress
NAME                      CLASS    HOSTS                         ADDRESS   PORTS   AGE
sample-ingress-by-nginx   <none>   a.foo.bar.com,b.foo.bar.com             80      19s
```

`kubectl get ingress --namespace=handson`のHOSTSを見ると、
- a.foo.bar.com
- b.foo.bar.com

の2種類のHostが登録されていることが確認出来ます。

それではHostヘッダを載せてcurlを投げてみましょう。

```
> curl localhost/podinfo -H 'Host:a.foo.bar.com'
{"pod_ip":"sample-app-a", "node_ip": "192.168.160.3", "client_ip": "10.244.0.5"}

> curl localhost/podinfo -H 'Host:b.foo.bar.com'
{"pod_ip":"sample-app-b", "node_ip": "192.168.160.3", "client_ip": "10.244.0.5"}
```

Hostヘッダベースでサービスを切り変えられていることが確認出来ました。
(ここでは環境変数を差し込んでpod_ipにアプリケーション名を識別出来るような値を入れています)

### リクエストパスによるルーティング

次にリクエストパスによるルーティングを確認してみます。

Deployment, NginxIngress, ClusterIPリソースをapplyします。
```
> kubectl apply -f ingress-path-base-routing.yaml
namespace/handson unchanged
ingress.networking.k8s.io/sample-ingress-by-nginx configured
service/sample-clusterip-a unchanged
service/sample-clusterip-b unchanged
deployment.apps/sample-deployment-a unchanged
deployment.apps/sample-deployment-b unchanged
```

以下のように従来のパスの前に「/a」「/b」を付けてcurlを投げることでサービスを切り替えることが出来ます。

```
> curl localhost/a/podinfo                       
{"pod_ip":"sample-app-a", "node_ip": "192.168.160.5", "client_ip": "10.244.0.5"}

> curl localhost/b/podinfo
{"pod_ip":"sample-app-b", "node_ip": "192.168.160.5", "client_ip": "10.244.0.5"}
```

また、リクエストの「/a/podinfo」を「/podinfo」に書き換えてプロキシする動作を実現するために、[NginxIngressのRewrite機能](https://kubernetes.github.io/ingress-nginx/examples/rewrite/)を活用しています。
- この機能が例えばGKEIngressなどで使えるかどうかは確認していませんが...

## 最後に

リソースを削除しておきましょう。
```
kubectl delete namespace handson
```