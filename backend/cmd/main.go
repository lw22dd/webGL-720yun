package main

func main() {
	// // 1. 初始化配置
	// if err := config.Init(); err != nil {
	// 	log.Fatalf("配置加载失败: %v\n", err)
	// }

	// // 2. 初始化日志系统和Redis
	// logger.Setup(&config.Conf.Logger)
	// redis.Init()

	// // 3. 创建Gin实例（根据配置设置模式）
	// if config.Conf.App.Debug {
	// 	gin.SetMode(gin.DebugMode)
	// } else {
	// 	gin.SetMode(gin.ReleaseMode)
	// }
	// router := gin.New()

	// // 4. 注册全局中间件
	// registerCoreMiddleware(router)

	// // 5. 注册路由
	// routes.RegisterAPIRoutes(router)

	// // 6. 配置HTTP服务器
	// srv := &http.Server{
	// 	Addr:           ":" + strconv.Itoa(config.Conf.Server.Port),
	// 	Handler:        router,
	// 	ReadTimeout:    30 * time.Second,
	// 	WriteTimeout:   30 * time.Second,
	// 	MaxHeaderBytes: 1 << 20,
	// }

	// // 7. 优雅关机处理
	// go func() {
	// 	logger.Info("服务器启动", zap.String("address", srv.Addr))
	// 	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
	// 		logger.Fatal("服务器异常退出", zap.Error(err))
	// 	}
	// }()

	// quit := make(chan os.Signal, 1)
	// signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	// sig := <-quit
	// logger.Info("正在关闭服务...", zap.String("signal", sig.String()))

	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()

	// if err := srv.Shutdown(ctx); err != nil {
	// 	logger.Error("强制关闭服务", zap.Error(err))
	// }
	// logger.Info("服务已安全停止")
}

// // 核心中间件注册（与网关功能强相关）
// func registerCoreMiddleware(r *gin.Engine) {
// 	r.Use(
// 		gin.Recovery(),        // 官方恢复中间件
// 		logger.GinZapLogger(), // 自定义日志中间件
// 	)
// }
