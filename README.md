# k8s-handson

このハンズオンの大部分は「[青山真也著: Kubernetes完全ガイド 第2版](https://github.com/MasayaAoyama/kubernetes-perfect-guide)」を参考に作成しています。

## 事前準備

ハンズオンを始める前に、ツールのインストールとk8sクラスタの作成を行います。
([こちら](https://github.com/OriishiTakahiro/k8s_handson)を参考にさせていただきました)

### ツール

ローカルでk8s環境を作成するために以下をインストールします。

- [Docker](https://docs.docker.com/get-docker/)
- [kubectl](https://kubernetes.io/ja/docs/tasks/tools/install-kubectl/)
- [kind](https://kind.sigs.k8s.io/docs/user/quick-start/#installation)

### クラスタとkubectl

`kind create cluster --config cluster.yaml`コマンドを実行してk8sクラスタを作成します。
- ワーカーノード3台のクラスターを作成します

```
kind create cluster --config cluster.yaml
Creating cluster "kind" ...
 ✓ Ensuring node image (kindest/node:v1.19.1) 🖼 
 ✓ Preparing nodes 📦 📦 📦 📦  
 ✓ Writing configuration 📜 
 ✓ Starting control-plane 🕹️ 
 ✓ Installing CNI 🔌 
 ✓ Installing StorageClass 💾 
 ✓ Joining worker nodes 🚜 
Set kubectl context to "kind-kind"
You can now use your cluster with:

kubectl cluster-info --context kind-kind

Thanks for using kind! 😊
```

また、このコマンドによりkubectlのコンテキストも自動的に設定されます。

## ハンズオンが終わったら

ハンズオンが終わったらクラスターを削除します。

```
kind delete cluster
```

## Tips

### k8sリソースの削除

ハンズオンでは全てのk8sリソースに`namespace: handson`を設定しています。
そのため、`kubectl delete namespace handon`でk8sリソースを削除出来ます。