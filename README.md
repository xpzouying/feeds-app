# feeds-app

feeds app demo.

## DDD模式


## 架构设计

包含如下，

1. `server` - Go实现的服务端。

2. 基础设施

    - `prometheus` - metrics统计。

    - `OpenTracing`
    
        安装
        
        ```bash
        docker run -d -p 5775:5775/udp -p 16686:16686 jaegertracing/all-in-one:latest

        # open URL of jaeger UI
        # http://docker.zy.local:16686/
        ```

    - `zipkin` - https://github.com/openzipkin/zipkin

        安装

        ```bash
        docker run -d -p 9411:9411 openzipkin/zipkin

        # http://docker.zy.local:9411/zipkin/
        ```

## 运行

### 服务端

```bash
# 运行服务

# Use zipkin for Tracing
go run server/cmd/main.go -zipkin.addr="http://localhost:9411/api/v2/spans"

# 请求
curl -XPOST -d '{"page": 0, "count": 10}' http://localhost:8080/feeding/feeds

# 浏览器查看metrics信息
http://localhost:8080/metrics
```

## 参考资料

### 关于基础设施

- [Prometheus Guide - Go Application](https://prometheus.io/docs/guides/go-application/)

- [OpenTracing常用Tags/Logs字段](https://github.com/opentracing/specification/blob/master/semantic_conventions.md)

