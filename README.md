# 项目目录

```
MicroCraft/
├─ go.mod
├─ go.sum
├─ README.md
├─ configs/
│  └─ config.yaml                 #配置文件
├─ cmd/
│  └─ server/
│     └─ main.go                  #程序入口
├─ internal/
│  ├─ config/					  #配置层	
│  │  └─ config.go                
│  ├─ router/					  #路由注册层
│  │  └─ router.go                
│  ├─ middleware/				  #中间件层
│  │  └─ middleware.go            
│  ├─ controller/				  #控制器层
│  │  └─ user.go                  
│  ├─ service/					  #业务逻辑层
│  │  └─ user.go                  
│  ├─ dao/						  #数据访问层
│  │  └─ mysql/
│  │     ├─ mysql.go              #DB初始化连接
│  │     └─ user.go               #user CRUD
│  └─ model/					  #数据模型层
│     └─ user.go                  
└─ pkg/
   ├─ response/
   │  └─ response.go              #统一返回
   ├─ errors/
   │  └─ errors.go                #错误码
   └─ utils/
      └─ utils.go                 #工具函数
```

