# grpc-gateway-example

grpc-gateway的使用示例，其中集成了各类库。

# 集成

* grpc-gateway
* viper
* zap
* opentracing
* gorm
* redigo

# 目录结构
```
.
├── cmd                                 // 命令
├── comsumer                            // 消费者
├── cronjob                             // 定时任务
├── definition                          // 全局变量定义
├── hack                                // 自动化相关
├── handler                             // 处理器
├── logic                               // 业务逻辑层
├── model                               // 数据模型、DAO层
├── pkg
│     ├── bootstrap                     // 引导
│     ├── lib                           // 核心库
│     └── util                          // 公共工具库
└── proto
    ├── gen
    └── src                             // proto源码
```