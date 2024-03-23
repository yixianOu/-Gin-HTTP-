# 使用Gin框架的HTTP博客

是学习Gin的练手项目，会在源码中作大量解释代码逻辑的注释，并且在README文件中分析项目的各个模块及其流程

主要是要实现两大块的基础业务功能，功能点分别如下：
 * 标签管理：文章所归属的分类，也就是标签。我们平时都会针对文章的内容打上好几个标签，用于标识文章内容的要点要素，这样子便于读者的识别。
 * 文章管理：整个文章内容的管理，并且需要将文章和标签进行关联。

  目录结构：
* configs：配置文件。
   * config.yaml:对Server，App，Database，Email,JWT的默认配置

* docs：文档集合。
   * docs.go:将项目信息和接口路由信息按规范生成到包全局变量 doc 中
   * swagger.json:默认指向当前应用所启动的域名下的 swagger/doc.json 路径
   * swagger.yaml:swagger默认配置
   * sql包：内含blog.sql保存创建数据库表的MySQL语句，blog_tag、blog_auth、blog_article_tag、blog_article

* global：全局变量。
   * setting.go:声明全局对象，将配置信息和应用程序关联起来


* internal：内部模块。
   * dao：数据访问层（Database Access Object），所有与数据相关的操作都会在 dao 层进行，如 MySQL
      * dao.go:新建dao方法。调用了model包各模块的方法，是对数据库操作的第二层封装。目的是对数据访问对象的封装，上层访问数据的对象是Dao结构体而非model包的各结构体
      * tag.go：处理标签模块的 dao 方法，其根据参数生成tag并调用对应方法提供Dao属性作为入参
      * article.go:处理文章模块的 dao 方法，其根据参数生成article并调用对应方法提供Dao属性作为入参
      * auth.go:包含处理Auth模块的 dao 方法，其根据参数生成auth并调用Get方法提供Dao属性作为入参。

   * middleware：HTTP 中间件。
      * translations.go：用于编写针对 validator 的语言包翻译的相关功能。借助了多语言包locales、通用翻译器universal—translator和validator自带翻译器,通过GetHeader获取语言类型。
      * jwt.go:JWT 通过 GetHeader 方法从 gin.Context 中获取 token 参数，然后调用 app.ParseToken 对其进行解析，成功则执行c.Next，失败则根据返回的错误类型进行断言判定，然后响应并执行c.Abort回退。
      * access_log.go：访问日志AccessLogWriter结构体，其方法Write实现了对访问日志和响应体的双写。其方法AccessLog用AccessLogWriter代替gin.Context的Writer，此后对回复体的写入会被记录。还会自动将请求类型，响应状态，处理的开始结束时间记入日志

   * model：模型层，用于存放 model 对象。为上层提供直接操作数据库的方法
      * model.go:公共字段结构体; 借助GORM实现NewDBEngine方法；注册回调函数实现公共字段的处理，如新增行为，更新行为，删除行为，都会触发对应的回调函数
      * tag.go:标签结构体；操作标签模块，是对数据库的增删改查函数的第一层封装，并且只与实体产生关系
      * article.go:文章结构体；操作文章模块，是对数据库的增删改查函数的第一层封装，并且只与实体产生关系
      * auth.go:Auth结构体，其Get方法用于判断能否根据客户端传来的app_key和app_secret，在数据库blog_auth表中查到记录，查到则返回该记录。

   * routers：路由相关逻辑处理。
      * router.go:注册路由，apiv1路由组,Swagger，upload/file,static，auth；使用中间件,Logger,Recovery，Translations，JWT
      * api:解析唯一入参gin.Contex的各字段（如request）、完成入参绑定和判断、根据request的Context字段创建service并调用其方法、序列化结果响应到gin.Contex中，集四大功能板块的逻辑串联；日志、错误处理。
         * v1：直接调用service中封装好的操作数据库的函数和app中的参数校验函数和响应体函数
            * tag.go:标签模块的接口，包含相关路由的handler。借助service包实现
            * artical.go:文章模块的接口，包含相关路由的handler。借助service包实现
         * upload.go：声明Upload结构体，其方法UploadFile读取Context内的file后，调用service的UploadFile方法上传文件，未出错则返回文件访问地址
         * auth.go：auth模块的接口，包含auth路由的handle。借助service包的CheckAuth方法和app包的GenerateToken函数实现

   * service：项目核心业务逻辑，为api层提供Service结构体及其已封装入参校验功能的各个方法。
      * tag.go:针对业务接口中定义的的增删改查统行为进行了 Request 结构体编写,利用标签实现参数绑定和参数校验。对blog_tag的增删改查操作做第三层封装。
      * article.go:针对业务接口中定义的的增删改查统行为进行了 Request 结构体编写,利用标签实现参数绑定和参数校验。对blog_article的增删改查操作做第三层封装。
      * service.go:定义服务结构体，并用Context和数据库Engine实例化一个服务
      * upload.go：将上传文件工具库与具体的业务接口结合起来，作为Service的UploadFile方法提供给api
      * auth.go:声明AuthRequest结构体用于校验接口的入参，利用标签实现参数绑定和参数校验。其方法CheckAuth调用了dao.FetAuth，是对blog_auth查询的第三层封装。


 * pkg：项目相关的模块包。
    * errcode:错误码标准化
       * common_code.go:预定义项目中的一些公共错误码，便于引导和规范大家的使用
       * errcode.go:Error结构体；全局错误码的存储载体codes；错误处理公共方法，标准化错误输
       * module_code.go：针对标签,文章和上传模块，用业务错误码区分不同的失败行为
出；将错误码转换为http状态码

    * setting:借助viper处理配置的读取
       * setting.go:针对读取配置的行为进行封装，便于应用程序的使用
       * section.go:用于声明配置属性的结构体,编写读取区段配置的配置方法

    * logger:借助lumberjack进行日志写入
       * logger.go:日志分级;日志的实例初始化和标准化参数绑定;日志格式化和输出的方法；日志分级输出

    * convert：类型转换
       * convet.go:为StrTo结构体提供类型转换的方法

    * app:应用模块。组合第三方库提供的 API，向MiddleWare或api层提供根据需求进行再封装的函数。
       * pagination.go:分页处理
       * app.go:响应处理。实现响应结构体（has a gin.Context）及其方法
       * form.go:引入validator库。实现绑定和判断的函数BindAndValid是对shouldBind方法进行的二次封装，发生错误则使用Translator翻译错误消息体。声明了 ValidError 相关的结构体和方法
       * jwt.go：GenerateToken 根据appKey和appSecret生成 JWT Token；ParseToken 解析和校验 Token返回JWT的属性

    * util:一个上传文件的工具库，功能是针对上传文件时的一些相关处理。
       * md5.go:实现函数，对上传的文件名进行MD5处理后再返回，防止暴露原始名称

    * upload:处理上传操作
       * file.go:实现文件相关参数获取的函数（文件名，文件后缀，保存地址）、检查文件的函数（目标目录是否存在，文件后缀是否匹配，文件大小是否超出，是否允许写入目录）

 * storage：项目生成的临时文件。
    * logs:内含app.log，记录项目的日志信息
    * uploads:存储前端上传的图片

 * main.go:启动文件。init调用初始化方法，配置公共组件；setupFlag设置编译信息，setupSetting调用setting包的ReadSection方法读取yaml中的各部分数据并给对应的全局变量赋值，setupLogger设置日志，setupDBEngine连接数据库，setupValidator设置校验器，setupTracer设置追踪器




本项目是对《Go 语言编程之旅》的第二章教程的复现。
