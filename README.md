# 使用Gin框架的HTTP博客

是学习Gin的练手项目，会在源码中作大量解释代码逻辑的注释，并且在README文件中分析项目的各个模块及其流程

主要是要实现两大块的基础业务功能，功能点分别如下：
 * 标签管理：文章所归属的分类，也就是标签。我们平时都会针对文章的内容打上好几个标签，用于标识文章内容的要点要素，这样子便于读者的识别。
 * 文章管理：整个文章内容的管理，并且需要将文章和标签进行关联。

  目录结构：
* configs：配置文件。
   * config.yaml:对Server，App，Database的默认配置
* docs：文档集合。
   * docs.go:将项目信息和接口路由信息按规范生成到包全局变量 doc 中
   * swagger.json:默认指向当前应用所启动的域名下的 swagger/doc.json 路径
   * swagger.yaml:swagger默认配置
* global：全局变量。
   * setting.go:将配置信息和应用程序关联起来
* internal：内部模块。
   * dao：数据访问层（Database Access Object），所有与数据相关的操作都会在 dao 层进行，例如 MySQL、ElasticSearch 等。
   * middleware：HTTP 中间件。
      * translations.go：用于编写针对 validator 的语言包翻译的相关功能。引入了多语言包locales、通用翻译器universal—translator和validator自带翻译器,通过GetHeader获取语言类型
   * model：模型层，用于存放 model 对象。
      * model.go:公共字段结构体; 借助GORM实现NewDBEngine方法
      * tag.go:标签结构体；操作标签模块，是对数据库的增删改查函数的第一层封装，并且只与实体产生关系
      * article.go:文章结构体；操作文章模块，是对数据库的增删改查函数的第一层封装，并且只与实体产生关系
   * routers：路由相关逻辑处理。
      * router.go:注册路由，apiv1路由组,Swagger；中间件,Logger,Recovery，Translations
      * api:路由handler
         * v1：直接调用service中封装好的操作数据库的函数
            * tag.go:标签模块的接口编写
            * artical.go:文章模块的接口编写
   * service：项目核心业务逻辑。
      * tag.go:针对业务接口中定义的的增删改查和统计行为进行了 Request 结构体编写,应用标签实现参数绑定和参数校验。对数据库的增删改查操作做第三层封装，函数会对入参进行参数校验
      * article.go:针对业务接口中定义的的增删改查和统计行为进行了 Request 结构体编写,应用标签实现参数绑定和参数校验，对数据库的增删改查操作做第三层封装，函数会对入参进行参数校验
 * pkg：项目相关的模块包。
    * errcode:错误码标准化
       * common_code.go:预定义项目中的一些公共错误码，便于引导和规范大家的使用
       * errcode.go:Error结构体；全局错误码的存储载体codes；错误处理公共方法，标准化错误输
出；将错误码转换为http状态码
    * setting:借助viper处理配置的读取
       * setting.go:针对读取配置的行为进行封装，便于应用程序的使用
       * section.go:用于声明配置属性的结构体并编写读取区段配置的配置方法
    * logger:借助lumberjack进行日志写入
       * logger.go:日志分级;日志的实例初始化和标准化参数绑定;日志格式化和输出的方法；日志分级输出
    * convert：类型转换
       * convet.go:为StrTo结构体提供类型转换的方法
    * app:应用模块
       * pagination.go:分页处理
       * app.go:响应处理
       * form.go:针对shouldBind方法进行了二次封装，发生错误则使用Translator翻译错误消息体
 * storage：项目生成的临时文件。
 * scripts：各类构建，安装，分析等操作的脚本。

 * main.go:启动文件。init调用初始化方法，配置公共组件；
