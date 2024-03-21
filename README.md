# 使用Gin框架的HTTP博客

是学习Gin的练手项目，会在源码中作大量解释代码逻辑的注释，并且在README文件中分析项目的各个模块及其流程

主要是要实现两大块的基础业务功能，功能点分别如下：
 * 标签管理：文章所归属的分类，也就是标签。我们平时都会针对文章的内容打上好几个标签，用于标识文章内容的要点要素，这样子便于读者的识别和 SEO 的收录等。
 * 文章管理：整个文章内容的管理，并且需要将文章和标签进行关联。

  目录结构：
* configs：配置文件。
   * config.yaml:对Server，App，Database的默认配置
* docs：文档集合。
* global：全局变量。
* internal：内部模块。
   * dao：数据访问层（Database Access Object），所有与数据相关的操作都会在 dao 层进行，例如 MySQL、ElasticSearch 等。
   * middleware：HTTP 中间件。
   * model：模型层，用于存放 model 对象。
      * model.go:公共字段结构体
      * tag.go:标签结构体；handle方法
      * article.go:文章结构体；handle方法
   * routers：路由相关逻辑处理。
      * router.go:注册路由
   * service：项目核心业务逻辑。
 * pkg：项目相关的模块包。
    * errcode:错误码标准化
       * common_code.go:预定义项目中的一些公共错误码，便于引导和规范大家的使用
       * errcode.go:Error结构体；全局错误码的存储载体codes；错误处理公共方法，标准化错误输
出；将错误码转换为http状态码
    * setting:编写组件
       * setting.go:针对读取配置的行为进行封装，便于应用程序的使用
       * section.go:用于声明配置属性的结构体并编写读取区段配置的配置方法
 * storage：项目生成的临时文件。
 * scripts：各类构建，安装，分析等操作的脚本。
 * third_party：第三方的资源工具，例如 Swagger UI。
