# CMD
* Cobra：按照Cobra标准要创建main.go和一个rootCmd文件，另外我们有子命令server
  1. rootCmd： rootCmd表示在没有任何子命令的情况下的基本命令
  2. &cobra.Command：
     - Use：Command的用法，Use是一个行用法消息
     - Short：Short是help命令输出中显示的简短描述
     - Run：运行:典型的实际工作功能。大多数命令只会实现这一点；另外还有PreRun、PreRunE、PostRun、PostRunE等等不同时期的运行命令，但比较少用，具体使用时再查看亦可
  3. rootCmd.AddCommand：AddCommand向这父命令（rootCmd）添加一个或多个命令
  4. serverCmd.Flags().StringVarP()：

# 服务端模块server
* server流程
  1. 启动监听
  2. 获取TLS
  3. 创建内部服务
  4. 创建tls.NewListener
  5. 服务开始接受请求