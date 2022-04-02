# FOR PIPI

## 初始化

### 1. 安装`Go`环境

### 2. 下载代码`git clone`
```
git clone https://github.com/waiter/pipi-gf-demo
```

### 3. 将`manifest/sql/create.sql`导入数据库，创建user表

### 4. 修改`manifest/config/config.yaml`中数据库配置

### 5. 安装以来`go mod tidy`

### 6. 使用`go run main.go`或者`make dev`来启动，正常情况可以看到以下日志：
```
  ADDRESS | METHOD |        ROUTE         |                                  HANDLER                                   |         MIDDLEWARE           
----------|--------|----------------------|----------------------------------------------------------------------------|------------------------------
  :8199   | ALL    | /*                   | github.com/gogf/gf/v2/net/ghttp.internalMiddlewareServerTracing            | GLOBAL MIDDLEWARE            
----------|--------|----------------------|----------------------------------------------------------------------------|------------------------------
  :8199   | ALL    | /*                   | github.com/gogf/gf/v2/net/ghttp.MiddlewareHandlerResponse                  | GLOBAL MIDDLEWARE            
----------|--------|----------------------|----------------------------------------------------------------------------|------------------------------
  :8199   | ALL    | /api.json            | github.com/gogf/gf/v2/net/ghttp.(*Server).openapiSpec                      |                              
----------|--------|----------------------|----------------------------------------------------------------------------|------------------------------
  :8199   | ALL    | /swagger/*           | github.com/gogf/gf/v2/net/ghttp.(*Server).swaggerUI                        | HOOK_BEFORE_SERVE            
----------|--------|----------------------|----------------------------------------------------------------------------|------------------------------
  :8199   | POST   | /user/check-passport | pipi.com/gogf/pipi-gf-demo/internal/controller.(*cUser).CheckPassport | service.(*sMiddleware).Ctx   
          |        |                      |                                                                            | service.(*sMiddleware).CORS  
----------|--------|----------------------|----------------------------------------------------------------------------|------------------------------
  :8199   | POST   | /user/is-signed-in   | pipi.com/gogf/pipi-gf-demo/internal/controller.(*cUser).IsSignedIn    | service.(*sMiddleware).Ctx   
          |        |                      |                                                                            | service.(*sMiddleware).CORS  
----------|--------|----------------------|----------------------------------------------------------------------------|------------------------------
  :8199   | GET    | /user/profile        | pipi.com/gogf/pipi-gf-demo/internal/controller.(*cUser).Profile       | service.(*sMiddleware).Ctx   
          |        |                      |                                                                            | service.(*sMiddleware).CORS  
          |        |                      |                                                                            | service.(*sMiddleware).Auth  
----------|--------|----------------------|----------------------------------------------------------------------------|------------------------------
  :8199   | POST   | /user/sign-in        | pipi.com/gogf/pipi-gf-demo/internal/controller.(*cUser).SignIn        | service.(*sMiddleware).Ctx   
          |        |                      |                                                                            | service.(*sMiddleware).CORS  
----------|--------|----------------------|----------------------------------------------------------------------------|------------------------------
  :8199   | POST   | /user/sign-out       | pipi.com/gogf/pipi-gf-demo/internal/controller.(*cUser).SignOut       | service.(*sMiddleware).Ctx   
          |        |                      |                                                                            | service.(*sMiddleware).CORS  
----------|--------|----------------------|----------------------------------------------------------------------------|------------------------------
  :8199   | POST   | /user/sign-up        | pipi.com/gogf/pipi-gf-demo/internal/controller.(*cUser).SignUp        | service.(*sMiddleware).Ctx   
          |        |                      |                                                                            | service.(*sMiddleware).CORS  
----------|--------|----------------------|----------------------------------------------------------------------------|------------------------------
```

## 多环境配置

当前支持两种环境：

- 开发环境：
  - 配置文件：`manifest/config/config.yaml`
  - 启动命令：`go run main.go`或者`make dev`
- 正式环境：
  - 配置文件：`manifest/config/config.prod.yaml`
  - 启动命令：`go run main.go --gf.gcfg.file=config.prod.yaml`或者`make prod`

## `WebSocket`相关说明

### 如何调试

启动服务后，访问`http://localhost:8199`即可

### 关于心跳

- 当前需要前端每隔一定时间（比`internal/consts/consts.go`中定义的`WebSocketHeartBeatTime`短即可）给服务端发送心跳包，直接发送内容为`ping`的包即可
- 另外，所有前端发送的包（不只是心跳包）都会刷新后端记录本`WebSocket`的超时时间

### 消息格式

#### 前端
前端给后端主要以**JSON形式传递**：
- `cmd`：本次命令类型
- `data`: 本次携带数据

例如：
```JSON
{"cmd":"testAdd","data":[1,78,3]}
```
#### 后端
后端回包格式与前端一致，会额外增加：
- `unique`：本WebSocket唯一ID
- `pack`：服务端发出时间

例如
```JSON
{"cmd":"add","data":82,"pack":1648906839,"unique":"ws-RFbD56T"}
```
另外，后端只需要在`internal/service/ws/logic.go`中添加相关逻辑即可