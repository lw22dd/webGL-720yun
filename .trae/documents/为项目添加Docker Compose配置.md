1. 创建docker-compose.yml文件，包含MySQL和Redis服务
2. 配置MySQL服务：

   * 版本：8.0

   * 端口映射：3306:3306

   * 环境变量：设置root密码123456、创建720db数据库

   * 数据持久化：挂载volumes保存数据
3. 配置Redis服务：

   * 版本：7.0

   * 端口映射：6379:6379

   * 环境变量：设置空密码

   * 数据持久化：挂载volumes保存数据
4. 配置网络：使用默认bridge网络，确保容器间通信
5. 确保配置与项目的config.dev.yaml文件中的数据库连接参数匹配

这个配置将允许项目通过localhost访问Docker中的MySQL和Redis服务，与当前配置无缝兼容。
