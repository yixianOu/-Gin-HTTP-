# 使用Gin框架的HTTP博客

是学习Gin的练手项目，会在源码中作大量解释代码逻辑的注释，并且在README文件中分析项目的各个模块及其流程

主要是要实现两大块的基础业务功能，功能点分别如下：
 * 标签管理：文章所归属的分类，也就是标签。我们平时都会针对文章的内容打上好几个标签，用于标识文章内容的要点要素，这样子便于读者的识别。
 * 文章管理：整个文章内容的管理，并且需要将文章和标签进行关联。

  目录结构：
* configs：配置文件。
   * config.yaml:对Server，App，Database，Email，JWT的默认配置
   * config.go:由validator生成，用于数据绑定

* docs：文档集合。
   * docs.go:自动将项目信息和接口路由信息按规范生成到包全局变量 doc 中
   * swagger.json:默认指向当前应用所启动的域名下的 swagger/doc.json 路径
   * swagger.yaml:swagger默认配置
   * sql包：内含blog.sql保存创建数据库表的MySQL语句，blog_tag、blog_auth、blog_article_tag、blog_article

* global：全局变量。
   * setting.go:声明全局对象，将配置信息和应用程序关联起来
   * validator.go:声明校验器全局变量
   * db.go：声明数据库引擎全局变量


* internal：内部模块。
   * dao：数据访问层（Database Access Object），所有与数据相关的操作都会在 dao 层进行，如 MySQL
      * dao.go:新建dao方法。是对数据库操作的第二层封装，目的是对数据访问对象的封装，service层访问数据的对象是Dao结构体而非model包的各结构体
      * tag.go：处理标签模块的 dao 方法，其根据参数生成tag，并调用model包中对应方法，然后提供Dao属性作为入参
      * article.go:处理文章模块的 dao 方法，其根据参数生成article，并调用model包中对应方法，然后提供Dao属性作为入参
      * article_tag:处理文章_标签模块的 dao 方法，其根据参数生成articleTag，并调用model包中对应方法，然后提供Dao属性作为入参
      * auth.go:包含处理Auth模块的 dao 方法，其根据参数生成auth并调用model包中的Get方法，然后提供Dao属性作为入参。

   * middleware：HTTP 中间件。
      * translations.go：用于编写针对 validator 的语言包翻译的相关功能。借助了多语言包locales、通用翻译器universal—translator和validator自带翻译器,通过GetHeader从gin.Context中获取语言类型。
      * jwt.go:JWT 通过 GetHeader 方法从 gin.Context 中获取 token 参数，然后调用 app包的ParseToken 对其进行解析，成功则执行c.Next，失败则根据返回的错误类型进行断言判定，然后响应错误并执行c.Abort回退。
      * access_log.go：声明访问日志AccessLogWriter结构体，其方法Write实现了对访问日志和响应体的双写、方法AccessLog用AccessLogWriter代替gin.Context的Writer，此后对回复体的写入操作会被记录，同时会自动将请求类型，响应状态，处理的开始结束时间记入日志
      * recovery.go:创建Email的饿汉单例模式defaultMailer，在捕获到异常后调用 SendMail 方法进行预警邮件发送
      * app_info.go：在进程内的Context自动添加一些内部消息，如应用名称和应用版本号
      * limiter.go:将app封装的限流器方法与对应的中间件逻辑串联起来。函数可接受不同的限流器入参，若令牌桶中缺少可发出的令牌，则Abort回退
      * context_timeout.go：上下文超时时间控制，设置Context.Request的超时属性。处理超时的request会返回err

   * model：模型层，用于存放 model 对象。为dao层提供直接操作数据库的方法
      * model.go:公共字段结构体; 借助GORM实现NewDBEngine方法；注册回调函数实现公共字段的处理，如新增行为，更新行为，删除行为，都会触发对应的回调函数
      * tag.go:标签结构体；操作标签模块，借助GORM对数据库的增删改查函数进行的第一层封装，并且只与实体产生关系
      * article.go:文章结构体；操作文章模块，借助GORM对数据库的增删改查函数进行的第一层封装，并且只与实体产生关系
      * article_tag.go：文章_标签结构体；操作文章_标签模块,借助GORM对数据库的增删改查函数进行的第一层封装，并且只与实体产生关系
      * auth.go:Auth结构体，其Get方法用于判断能否根据客户端传来的app_key和app_secret，在数据库blog_auth表中查到记录，查到则返回该记录。

   * routers：路由相关逻辑处理。
      * router.go:注册路由，apiv1路由组，Swagger，upload/file,static，auth；使用中间件,Logger,Recovery，Translations，JWT，AccessLog,Recovery，AppInfo，RateLimiter，ContextTimeout
      * api:解析唯一入参gin.Contex的各字段（如request）、完成入参绑定和判断、根据request的Context字段创建service并调用其方法、序列化结果响应到gin.Contex中，集四大功能板块的逻辑串联；日志、错误处理。
         * v1：直接调用service中封装好的操作数据库的函数和app中的参数校验函数和响应体函数
            * tag.go:标签模块的接口，包含相关路由的handler。
            * artical.go:文章模块的接口，包含相关路由的handler。
         * upload.go：声明Upload结构体，其方法UploadFile读取Context内的file后，调用service的UploadFile方法上传文件，未出错则将文件访问地址存入响应体中。
         * auth.go：auth模块的接口，包含auth路由的handle。借助service包的CheckAuth方法和app包的GenerateToken函数实现。调用了app包中的参数校验函数和响应体函数。

   * service：项目核心业务逻辑，为api层提供Service结构体及其已封装入参校验功能的各个方法。
      * tag.go:针对业务接口中定义的的增删改查统行为进行了 Request 结构体编写,利用标签实现参数绑定和参数校验。对blog_tag的增删改查操作做第三层封装。
      * article.go:针对业务接口中定义的的增删改查统行为进行了 Request 结构体编写,利用标签实现参数绑定和参数校验。对blog_article的增删改查操作做第三层封装。
      * service.go:定义服务结构体（has a Context and Dao），提供New用Context和数据库Engine实例化一个服务
      * upload.go：将上传文件工具库与具体的业务接口结合起来，将UploadFile方法提供给api层
      * auth.go:声明AuthRequest结构体用于校验接口的入参，利用标签实现参数绑定和参数校验。其方法CheckAuth调用了dao.GetAuth，是对blog_auth查询的第三层封装。


 * pkg：项目相关的模块包。
    * errcode:错误码标准化
       * common_code.go:预定义项目中的一些公共错误码，便于引导和规范
       * errcode.go:声明Error结构体和全局错误码的存储载体codes；实现Error的公共处理方法，将错误码转换为http状态码，并标准化错误输出
       * module_code.go：针对标签,文章和上传模块，用业务错误码区分不同的失败行为


    * setting:借助viper处理配置的读取
       * setting.go:针对读取配置的行为进行封装，便于应用程序的使用
       * section.go:用于声明各个配置模块的结构体,编写读取区段配置属性的配置方法

    * logger:借助lumberjack进行日志写入
       * logger.go:日志分级;日志的实例初始化和标准化参数绑定的方法;日志内容的格式化和输出json的方法；日志的分级输出方法

    * convert：类型转换
       * convet.go:为StrTo结构体提供类型转换的方法

    * app:应用模块。组合第三方库提供的 API，向MiddleWare或api层提供根据需求进行再封装的函数。
       * pagination.go:分页处理，GetPage获取文章在第几页，GetPageSize获取文章在当前页的第几个。GetPageOffset获取文章的偏移量
       * app.go:响应处理。实现响应结构体（has a gin.Context）及其响应处理方法
       * form.go:引入validator库。声明了 ValidError 相关的结构体和方法。实现绑定和判断的函数BindAndValid是对shouldBind方法进行的二次封装，发生错误则使用Translator翻译错误响应体。
       * jwt.go：GenerateToken 根据appKey和appSecret生成 JWT Token；ParseToken 解析和校验 Token返回JWT的属性Claims

    * util:一个上传文件的工具库，功能是针对上传文件时的一些相关处理。
       * md5.go:实现函数，对上传的文件名进行MD5处理后再返回，防止暴露原始名称

    * upload:处理上传文件操作
       * file.go:实现了获取上传所需相关参数的各函数（文件名，文件后缀，保存地址）、检查可否上传的各函数（目标目录是否存在，文件后缀是否匹配，文件大小是否超出，是否允许写入目录）

    * email：import第三方库Gomail支持使用SMTP服务器发送电子邮件，对发送电子邮件的行为进行封装
       * email.go：定义了SMTPInfo结构体用于传递发送邮箱所必须的信息。Email是对SMTPInfo的封装，其方法SendMail调用NewMessage创建邮件实例并对其赋值，然后调用NewDialer和DialAndSend创建拨号实例将邮件寄出
    * limiter:调包Ratelimit库 提供了一个简单高效的令牌桶实现，并提供大量方法帮助实现限流器的逻辑。
       * limiter.go:声明LimiterIface接口定义方法；声明Limiter令牌桶存储令牌与名称的映射关系；声明LimiterBuketRule定义令牌桶的各规则属性
       * method_limiter.go:针对LimiterIface实现MethodLimiter限流器，其向MiddleWare层提供 Key 方法（根据 RequestURI 切割出核心路由作为键值对名称）和GetBucket 方法（根据路由名称获取和设置 Bucket 的对应属性）。为router.go提供了NewMethodLimiter方法和AddBuckets方法以初始化RateLimiter中间件的入参。

 * storage：项目生成的临时文件。
    * logs:内含app.log，存储项目的日志信息
    * uploads:存储前端上传的文件

 * main.go:启动文件。
    * init调用初始化方法，配置公共组件；
    * 设置运行模式，创建路由实例并初始化，开始监听并服务，实现优雅重启和停止.
    * setupFlag设置编译信息；setupSetting调用setting包的ReadSection方法，将配置文件内容映射到应用配置结构体；setupLogger设置日志；setupDBEngine连接数据库；setupValidator设置校验器；setupTracer设置追踪器。




本项目是对《Go 语言编程之旅》的第二章教程的复现。
reference：github.com/go-programming-tour-book/blog-service
