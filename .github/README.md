# Dawn
`Dawn` is an opinionated `web` framework that provides rapid development capabilities. It provides basic services such as logging, configuration, `fiber` extension, `gorm` extension, and event system. 

The core idea of ​​Dawn is modularity. High-level business modules can call low-level basic modules, such as databases and so on. 

Each module needs to implement its own two core methods of `Init` and `Boot`, and then register it in `Server`. General business modules need to implement its `Register Routes` method to register routes and provide `http` services.

The modules should be based on the principle of not recreating the wheel, and directly provides the original structure and method of the dependent library.

The libraries currently used are
- [klog](https://github.com/kubernetes/klog)
- [viper](https://github.com/spf13/viper)
- [godotenv](https://github.com/joho/godotenv)
- [fiber](https://github.com/gofiber/fiber)
- [gorm](https://github.com/go-gorm/gorm)
- [go-redis](https://github.com/go-redis/redis)

# Notice
**This project is still under development, please do not use it in a production environment.**

# Why dawn?
Tribute to the first episode of one piece romance dawn. Let us set sail towards romance with the sloop.
