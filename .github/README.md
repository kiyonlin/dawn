# Dawn
`Dawn` is a `web` framework that provides rapid development capabilities. It provides basic services such as logging, configuration, `fiber` extension, `gorm` extension, and event system. 

The core idea of ​​Dawn is modularity. High-level business modules can call low-level basic modules, such as databases and so on. 

Each module needs to implement its own two core methods of `Init` and `Boot`, and then register it in `Server`. General business modules need to implement its `Register Routes` method to register routes and provide `http` services.

The modules should be based on the principle of not recreating the wheel, and directly provides the original structure and method of the dependent library.

The libraries currently used are
- klog
- viper
- godotenv
- fiber
- gorm
- go-redis
