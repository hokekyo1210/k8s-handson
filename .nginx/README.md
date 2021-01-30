# podのIPとnodeのIPを返すnginx

```
.
├── Dockerfile
├── README.md
├── docker-compose.yaml
├── entrypoint.sh
├── nginx.conf
└── server.conf
```

## 動作確認

```
docker-compose up -d
curl localhost:7777/podinfo
{"pod_ip":"127.0.0.2", "node_ip": "127.0.0.1"}% 
```