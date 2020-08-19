# Mybed图床 

> 这是一个基于多家对象存储源的Golang开源图片托管程序。
> 本项目使用 Gin框架 搭建, 支持对接`Minio`、`阿里Oss`，`腾讯Cos`，`又拍云`，`七牛云` 等多家对象存储服务.
> 后台对用户管理。
> 支持配置多家存储源。
> 图片鉴黄配置等操作。

## 主要功能支持：

- 支持 图片拖拽、截图软件直接(Ctrl+V)地址上传。
- 支持对接Minio、阿里、又拍、七牛、腾讯等各大对象存储平台。
- 支持上传者IP记录
- 支持链接生成二维码。
- URL列表、缩略图。查看原图等功能。
- 图片鉴黄配置
- 游客、用户的上传管理
- 站点设置和上传规则配置等。

站点Demo：[http://mybed.cilidm.com/](http://mybed.cilidm.com/)

## 运行环境

- Go 1.12+
- MySQL5.5+

### 下载项目

```git
git clone https://github.com/cilidm/mybed.git
```

### 配置数据库

创建任意名称的数据库, 字符集选择 `utf8`, 排序规则选择 `utf8_general_ci`. 
修改conf文件夹下`app.ini`文件，`[database]` 下的`DBUser DBPwd DBHost DBTableName` 四项 
如果想使用sqlite3进行数据保存，可以修改`DBType = sqlite` 此时 `DBPath` 就是数据文件保存地址
注：如果是windows环境，使用sqlite需要自行配置gcc

### 配置文件

```app.ini
[runMode]
RunMode = debug     # 运行模式，生产模式请配置为release

[server]
HTTPPort = 8000     # 运行端口
ReadTimeout = 60    # 读超时时间 
WriteTimeout = 60   # 写超时时间

[app]
PageSize = 9        # 分页默认每页数据条数
JwtSecret = 23347$040412 # session加密
UploadTmpDir = "static/upload"  # 上传临时文件夹地址，请确保可以写入

[database]
DBType = mysql      # 数据库类型
DBUser = root       # 数据库账号名
DBPwd = 123456      # 数据库密码
DBHost = 127.0.0.1:3306 # 数据库地址+端口号
DBTableName = mytbed    # 表名
DBPath = mybed.db   # sqlite保存地址

[redis]
RedisAddr = 127.0.0.1:6379  # Redis地址
RedisPWD =      # Redis密码
RedisDB =       # RedisDB

```

### 启动项目

在完成了上述步骤后，执行go mod tidy下载依赖，启动main.go即可.

初始用户名: admin
初始密码: admin

启动后访问地址为：http://localhost:8000 , `8000`就是配置项`HTTPPort`的端口.

### 声明

本程序很多功能尚未完善，用做对外开放的图床时请注意配置图像审核，默认百度图像审核，请在设置里填写key secret。

### 反馈交流

**如果你遇到BUG可以在github反馈**

### Store配置说明

> 腾讯云Cos : 可以随便上传一张图片，然后打开图片的详情页，查看他的网址，网址的构成为：
> `http://[存储桶名称]-[AppID].cos.[COS区域].myqcloud.com/图片名称`
> 请求域名请填写: 
> `http://[存储桶名称]-[AppID].cos.[COS区域].myqcloud.com/`

### 开发计划

- [ ] Store功能调整测试
    - [x] 阿里云
    - [x] 七牛云
    - [x] 又拍云
    - [x] Minio
    - [x] 腾讯云
    - [ ] 百度云
    - [ ] 华为云
- [ ] 日志功能
- [ ] 水印功能
- [ ] IP黑名单拦截
- [ ] 游客每日上传数量限制(ip统计)
- [ ] 部分栏目增加搜索功能
- [ ] store实时管理
- [x] ssh服务器文件迁移
- [ ] 同一个store更换密钥后无法删除(加入验证文件是否存在，成功删除后才删除数据库)
- [ ] 鉴黄修改