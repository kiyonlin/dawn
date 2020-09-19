# Dawn
`Dawn`是一个提供快速开发能力的`web`框架。它提供了日志、配置、`fiber`扩展、`gorm`扩展、事件系统等基础服务。

`Dawn`的核心理念是模块化。高层的业务模块可以调用低层的基础模块，例如数据库等等。

每个模块都需要实现自己的`Init`，`Boot`这两个核心方法，然后注册到`Server`中。一般业务模块需要实现其`RegisterRoutes`方法，用于注册路由，提供`http`服务。

模块的封装本着不重复造轮子的原则，直接提供依赖库其原本的结构和方法。

目前用到的库有
- [klog](https://github.com/kubernetes/klog)
- [viper](https://github.com/spf13/viper)
- [godotenv](https://github.com/joho/godotenv)
- [fiber](https://github.com/gofiber/fiber)
- [gorm](https://github.com/go-gorm/gorm)
- [go-redis](https://github.com/go-redis/redis)

# 注意
**本项目还在开发中，请勿在生产环境中使用。**
