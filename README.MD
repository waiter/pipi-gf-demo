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
