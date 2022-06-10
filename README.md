# middleware-controller-manager



// 应用运行平台
// 基于oam实现。demo.caf
// 使用configmap存储demo.caf，
// 安装应用即将demo.caf转换成Application资源。
// 启用应用，vela up

// 应用定义：
// name，
// image，
// port，
// config，每个应用提供一个，可配置路径
// env，注入用户在系统内定义的所有环境变量，运行跨应用访问
// secret，每个应用一个，不允许跨应用访问
// volume，挂载一个标准地址，每个应用一个

// 系统设置
// ingress，配置一个通配符域名，需要外部访问的应用可以创建一个ingress路由到目标应用
// volume
// SSL证书管理
// 账户管理

// 应用设置
// config
// env
// secret
// 权限控制

// 应用开发平台
// todo

## middleware

type: Selfhosted/CloudProduct
默认只提供Selfhosted的处理逻辑，CloudProduct类型的资源由由云平台自行实现controller
