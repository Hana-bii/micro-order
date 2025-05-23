# micro-order
一个**微服务订单支付系统**

## 描述
项目分4个service、order、stock、payment微服务，代码解耦，独立容器运行，本地方便开发所以集合在一起。
- 基于**DDD领域驱动设计**、CQRS模式
- 微服务间使用**gRPC**通信，order与payment之间采用**rabbitMQ**异步通信。
- 使用**protobuf**, **protoc**, **openAPI**, **oapi-codegen** 生成grpc服务、http服务api框架代码
- 使用**consul**和**viper**实现服务发现和注册
- 使用成**OpenTelemetry**和**Jaeger**实现全链路追踪，解决了异步追踪断裂问题
- 集成**Prometheus**和**Grafana**实现关键指标监控，使用**logrus**实现结构化日志
- 实现基于**redis**的分布式锁
- 使用**stripe**完成支付

## request流转如下
1. 用户下单发起http请求，order服务向调用stock的grpc服务校验库存并扣减，返回订单状态。
2. order服务创建并存储订单信息，向消息队列中发送订单。
3. payment从消息队列中接收消息，调用stripe api生成支付链接，并调用order grpc服务更新订单状态信息为等待支付和返回支付链接。
4. 前端轮询order服务api，Get订单信息，获取支付url并跳转。用户点击链接支付后，payment使用webhook从stripe服务器获取支付状态，写入消息队列。
5. order获取消息并更新订单状态，根据状态信息跳转支付成功。
