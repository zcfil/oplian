## oplian项目结构

| 文件夹       | 说明                    | 描述                        |
| ------------ | ----------------------- | --------------------------- |
| `api`        | api层                   | api层 |
| `--v1`       | v1版本接口              | v1版本接口                  |
| `config`     | 配置包                  | config.yaml对应的配置结构体 |
| `core`       | 核心文件                | 核心组件(zap，viper，server)的初始化 |
| `global`     | 全局对象                | 全局对象 |
| `initialize` | 初始化 | router，gorm，validator，timer的初始化 |
| `--internal` | 初始化内部函数 | gorm 的 longger 自定义，在此文件夹的函数只能由 `initialize` 层进行调用 |
| `middleware` | 中间件层 | 用于存放 `gin` 中间件代码 |
| `model`      | 模型层                  | 模型对应数据表              |
| `--request`  | 入参结构体              | 接收前端发送到后端的数据。  |
| `--response` | 出参结构体              | 返回给前端的数据结构体      |
| `router`     | 路由层                  | 路由层 |
| `service`    | service层               | 存放业务逻辑问题 |
| `source` | source层 | 存放初始化数据的函数 |
| `utils`      | 工具包                  | 工具函数封装            |
| `--timer` | timer | 定时器接口封装 |
| `--upload`      | oss                  | oss接口封装        |

[//]: # (protoc --go_out=./service/lotus/pb/ ./service/lotus/proto/op.proto)

[//]: # (protoc --go-grpc_out=require_unimplemented_servers=false:./service/lotus/pb/ ./service/lotus/proto/op.proto)

[//]: # ()
[//]: # (protoc --go_out=./service/sysrpc/ ./service/sysrpc/proto/gateway.proto)

[//]: # (protoc --go-grpc_out=require_unimplemented_servers=false:./service/sysrpc/ ./service/sysrpc/proto/gateway.proto)


protoc --go_out=./service/pb/ ./service/proto/op.proto
protoc --go-grpc_out=require_unimplemented_servers=false:./service/pb/ ./service/proto/op.proto

protoc --go_out=./service/pb/ ./service/proto/gateway.proto
protoc --go-grpc_out=require_unimplemented_servers=false:./service/pb/ ./service/proto/gateway.proto

protoc --go_out=./service/pb/ ./service/proto/header.proto
protoc --go_out=./service/pb/ ./service/proto/slot/slot_header.proto

protoc --go_out=./service/pb/ ./service/proto/slot/slot_op.proto
protoc --go-grpc_out=require_unimplemented_servers=false:./service/pb/ ./service/proto/slot/slot_op.proto

protoc --go_out=./service/pb/ ./service/proto/slot/slot_gateway.proto
protoc --go-grpc_out=require_unimplemented_servers=false:./service/pb/ ./service/proto/slot/slot_gateway.proto

## 安装部署

<font color = sandybrown >注: 以运行环境为linux下Ubuntu系统，使用`root`用户进行操作</font>


首先，从 `GitHub` 克隆仓库: <br />
`git clone https://github.com/zcfil/oplian.git` <br/>

生成对应的安装包 <br/>
`make clean all`<br/>

生成的安装包为:<br/>
`oplian oplian-gateway oplian-op oplian-op-c2`


## oplian部署

<font color = sandybrown >注: 需要保证部署`oplian`的机器有公网`IP`</font>

创建`/root/oplian`目录 <br />
`mkdir -p /root/oplian/`<br />

将生成的对应`oplian`包放到该目录下<br />

执行授权 `chmod 777 oplian`

运行`oplian` <br /> 
`./oplian`<br />

系统会自动创建 `/root/oplian/config`目录，并创建`config.yaml，config_room.yaml`文件<br />

#### config.yaml文件
```
mysql:
  path: 127.0.0.1         // 数据库IP
  port: "3306"
  config: charset=utf8mb4&parseTime=True&loc=Local
  db-name: oplian_test    // 数据库名
  username: root          // 用户名
  password: 123456        // 密码
  max-idle-conns: 10
  max-open-conns: 100
  log-mode: error
  log-zap: false
system:
  env: public
  addr: 50005               // 对应的oplian端口
  db-type: mysql
  oss-type: local
  use-multipoint: false
  iplimit-count: 15000
  iplimit-time: 3600

```
在数据修改无误后再次运行`oplian`即可<br />
`./oplian` or `nohup ./oplian > log/oplian.log &`<br/>

## oplian-gateway部署

<font color = sandybrown >注: 需要保证部署`oplian-gateway`的机器有公网`IP`</font>


创建`/root/oplian`目录 <br/>
`mkdir -p /root/oplian/`<br/>

将生成的对应`oplian-gateway`包放到该目录下<br/>

执行授权 `chmod 777 oplian-gateway`<br/>

运行`oplian-gateway`<br/>
`./oplian-gateway run`<br/>

系统会自动创建 `/root/oplian/config`目录，并创建`config.yaml，config_room.yaml`文件<br />
#### config.yaml文件
```
mysql:
  path: 127.0.0.1         // 数据库IP
  port: "3306"
  config: charset=utf8mb4&parseTime=True&loc=Local
  db-name: oplian_test    // 数据库名
  username: root          // 用户名
  password: 123456        // 密码
  max-idle-conns: 10
  max-open-conns: 100
  log-mode: error
  log-zap: false
```
#### config_room.yaml文件
```
web:
  addr: 127.0.0.1:50005         // 填写对应前面的oplian地址和端口
  token: slfsdaklfhasldfjda
gateway:
  gateWayId: 8b356e3d-5e35-47a6-a933-07baa8cc0fcb
  port: 50006       
  ip: 127.0.0.1     // 填写对应机器内网IP
  token: slfsdaklfhasldfjda
```

在数据无误后再次运行`oplian-gateway`<br/>

其中对应命令`./oplian-gateway run --help`<br/>

```
NAME:
   oplian-gateway run - 运行

USAGE:
   oplian-gateway run [command options] [arguments...]

OPTIONS:
   --listen-ip value  内网IP  
   --help,-h         show help
```

系统会自动检测内网`IP`，当有多个或者检测不到时，需要手动录入对应内网`IP`<br/>
启动命令变为<br/>
`./oplian-gateway run --listen-ip 对应内网IP`

## oplian-op部署

<font color = sandybrown >注: 需要保证部署的`op`机器在内网上和对应的`gateway`机器连通</font><br/>

在基于`oplian-gateway`已创建为前提下，借助指令即可获取对应的`oplian`文件和创建对应目录<br/>
`curl -fsSL http://10.0.1.77:50009/download/sh | bash && source ~/.bashrc`<br/>

其中`10.0.1.77`为对应`gateway`的内网地址，根据不同`gateway`进行替换即可<br/>

如果系统未安装对应信息，可以借助命令进行系统初始化，安装必备软件<br/>
其中对应命令`./oplian-op init`<br/>

确保系统无误后
执行对应命令`./oplian-op run --help`<br/>

```
NAME:
   oplian-op run - 运行

USAGE:
   oplian-op run [command options] [arguments...]

OPTIONS:
   --dc-type               是否是DC原值主机，true为设置该机器为原值主机类型，false为不设置，默认为非原值主机 (default: false)
   --storage               是否是存储机 (default: false)
   --listen-ip value       内网IP 
   --paramters-path value  指定证明参数路径
   --not-proof-parameters  如果不存在证明参数，程序会到文件管理平台上下载。但由于批量下载会出现问题，建议提前在/mnt/md0中准备filecoin-proof-parameters (default: false)
   --worker                是否是算力机 (default: false)
   --miner                 是否是miner机 (default: false)
   --help,-h              show help
```

运行命令<br/>
`./oplian-op run` 根据需要拼接对应的命令<br/>

即可正常启动`oplian-op`<br/>
