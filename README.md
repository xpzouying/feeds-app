# feeds-app

feeds app demo.

## DDD模式


## 架构设计

包含如下，

1. `server` - Go实现的服务端。

2. 基础设施

    - `prometheus` - metrics统计。


## 运行

### 服务端

```bash
# 运行服务
go run server/cmd/main.go

# 请求
curl -XPOST -d '{"page": 0, "count": 10}' http://localhost:8080/feeding/feeds

# 浏览器查看metrics信息
http://localhost:8080/metrics
```

## 参考资料

### 关于基础设施

- [Prometheus Guide - Go Application](https://prometheus.io/docs/guides/go-application/)
