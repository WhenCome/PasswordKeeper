# PasswordKeeper

这是一个简单的密码保存工具，用于管理个人密码，用于拥有大量密码需要维护，且经常在公共场合需要使用密码的情形。此工具可以直接将密码复制到剪贴板，避免密码被偷窥。
    
##  使用步骤：

以下假设编译的程序是pk。

* 使用 “pk init” 执行初始化操作，包括设置安全码、备份目录等，程序会自动生成RSA证书用于加解密；
* 使用 “pk set [item]” 添加新项（如果item已经存在则会覆盖设置），这里的item是一个标识，比如 sina；
* 使用 “pk get [item]” 可以获取密码并将密码复制到剪贴板，此命令不会展示密码内容；
* 使用 “pk get [item]” 可以直接获取密码并将密码展示到命令窗口；
* 使用 “pk help” 查看更多命令以及用法；

**说明:**
1. 此程序不会自动备份相关数据，需要手动调用 "pk sync" 进行同步备份；
2. 为了使用方便，可以将编译的程序放到环境变量“PATH”中，这样就可以直接调用。

## Dependencies

* github.com/bgentry/speakeasy    用于支持密码读取
* github.com/atotto/clipboard        用于复制密码到剪贴板
* github.com/mattn/go-sqlite3       用于操作sqlite数据库，如果编译时sqlite3.go中的context提示不存在或者报错，则需要将import中的context改为golang.org/x/net/context。

