# k8s-sync

将api server的pod 数据同步到 mysql

原因：用户发布项目实例后，想查询一下实例的信息（比如ip等），通常需要与 静态的项目信息 进行关联查询、模糊查询等。这时，实例信息在apiserver 上， 项目信息在mysql 上，处理起来就不太方便。

## mysql 数据库表结构

```
k8s.tb_pod

1. id
2. env
3. app_name         // 建有索引
4. isolation
5. name             // 建有索引
6. container_id
7. node
8. ip
9. status
10. create_at
11. update_at
```

## 功能

1. 从配置文件中读取 apiserver token
2. 从命令行读取日志级别，日志打到当前文件夹下
3. 监听apiserver 保证数据同步的实时性
4. 定时根据mysql 中的数据反查 apiserver，清理掉漏网之鱼
