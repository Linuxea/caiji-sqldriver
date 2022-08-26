## sql driver proxy

### 背景

本地实现读写分离, 将部分数据库流量按照配置打到指定连接, 对某连接的稳定性进行测试

### 方案
- 一 gorm db 上

对原来的 gorm.DB 进行封装, 对原来代码有侵入性, 且gorm.DB 的设计不是接口, 需要覆盖的方法众多

- 二 gorm db 中

通过 gorm 支持的插件形式嵌入,目前使用 gorm 版本过低, 不能使用作者开发的 <code>dbresolver</code> 直接使用, 需要对 <code>gorm</code> 代码进行修改

- 三 gorm db 下

通过驱动层进行接口实现重写, 优点:实现接口少,标准统一. 缺点:需要同时维护多个数据库连接

### 实现

采用方案三

通过代理的设计模式,实现数据库相关 <code>driver</code>, <code>connection</code>, <code>tx</code> 基本 <code>interface</code>

通过一对多关系维护真实的数据库连接

根据 sql 的读写类型, 连接对象的权重配置来选择唯一一条连接,将调用委托给真实的数据库连接对象


![连接示意图](https://res.caijiyouxi.com/static/activity/pre/2022-06-10-17-22-02.6694.png)



### 示例

#### 配置多个 dsn
```golang
    // 代理
	proxy := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Cfg.DbRW.User,
		config.Cfg.DbRW.Password,
		config.Cfg.DbRW.Host,
		config.Cfg.DbRW.Port,
		config.Cfg.DbRW.Database,
	)

	// 直连读
	directWrite := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Cfg.DbW.User,
		config.Cfg.DbW.Password,
		config.Cfg.DbW.Host,
		config.Cfg.DbW.Port,
		config.Cfg.DbW.Database,
	)

    // 只读库 r1
	r1 := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Cfg.DbR2.User,
		config.Cfg.DbR2.Password,
		"rr-wz9h6ed1065k1l4se.mysql.rds.aliyuncs.com",
		config.Cfg.DbR2.Port,
		config.Cfg.DbR2.Database,
	)

    // 只读库 r2
	r2 := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Cfg.DbR2.User,
		config.Cfg.DbR2.Password,
		config.Cfg.DbR2.Host,
		config.Cfg.DbR2.Port,
		config.Cfg.DbR2.Database,
	)
```


#### 整合成一个 dsn
```golang
    // 根据需求进行读,写权重的配置,以及 `Flag` 配置(flag 用于日志输出时观测对应的执行数据库, 可自定义)
	dsns := []*cjsqldriver.Dsn{
		{
			ReadWeight:  1,
			WriteWeight: 1,
			Dsn:         proxy,
			Flag:        "proxy",
		},
		{
			ReadWeight:  0,
			WriteWeight: 1,
			Dsn:         directWrite,
			Flag:        "directWrite",
		},
		{
			ReadWeight:  1,
			WriteWeight: 0,
			Dsn:         r1,
			Flag:        "r1",
		},
		{
			ReadWeight:  2,
			WriteWeight: 0,
			Dsn:         r2,
			Flag:        "r2",
		},
	}

    //进行 json 序列化
	dsnsJson, _ := json.Marshal(dsns)
```


####
```golang
    // 使用 `cjmysql` 驱动创建连接
    con, err = gorm.Open("cjmysql", dsnsJson)
```