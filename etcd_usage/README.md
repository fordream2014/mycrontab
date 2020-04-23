* 启动etcd单机版
```cassandraql
nohup ./etcd --listen-client-urls http://0.0.0.0:2379 --advertise-client-urls http://0.0.0.0:2379 &
```

* 解决clientv3报错
```
# github.com/coreos/etcd/clientv3/balancer/picker
../../../pkg/mod/github.com/coreos/etcd@v3.3.20+incompatible/clientv3/balancer/picker/err.go:37:44: undefined: balancer.PickOptions
../../../pkg/mod/github.com/coreos/etcd@v3.3.20+incompatible/clientv3/balancer/picker/roundrobin_balanced.go:55:54: undefined: balancer.PickOptions
# github.com/coreos/etcd/clientv3/balancer/resolver/endpoint
../../../pkg/mod/github.com/coreos/etcd@v3.3.20+incompatible/clientv3/balancer/resolver/endpoint/endpoint.go:114:78: undefined: resolver.BuildOption
../../../pkg/mod/github.com/coreos/etcd@v3.3.20+incompatible/clientv3/balancer/resolver/endpoint/endpoint.go:182:31: undefined: resolver.ResolveNowOption

```
golang版本： go version go1.11.13 darwin/amd64

解决方法：
```cassandraql
1、修改依赖为v1.26.0
go mod edit -require=google.golang.org/grpc@v1.26.0

2、下载grpc
go get -u -x google.golang.org/grpc@v1.26.0
```