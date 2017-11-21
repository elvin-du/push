### RPC默认端口：
* GATE :50001

### TCP默认端口：
* GATE :51001

### HTTP默认端口：
#### REST_API :52001
* POST /push

```	
//认证方式
authorization: Bearer BASE64(RC4(app_id:app_secret)
RC4_KEY: 01e9175ca8805cc2137c44eb86184922
```

```
//参数,JSON结构
{
	"alert":"hi test",//必选
	"audience":"5a0ea86008b62f0928970a52", //reg_id,必选
	"ttl":86400, //有效期，单位秒,必选
	"ios":{
		"sound":"", //可选
		"badge":1, //可选
		"production":false, //可选
		"extras":{"kind":2}, //JSON结构,可选
	},
	"android":{
		"extras":{"kind":2}, //JSON结构,可选
	},
	"extras":{"kind":2}, //JSON结构,可选
}
```

#### INTERNAL REST_API :52010 //内部使用，不可以公开端口
* GET /apps 所有APP信息
* POST /app 创建APP 
* PUT /app 更新APP信息,必须包含所有内容
* DELETE /app/:id 删除APP
		
#### REGISTRY REST_API :52020 
* POST /register 用户注册
* DELETE /register/:id  删除注册
		
```	
//认证方式
authorization: Bearer BASE64(RC4(app_id:app_secret)
RC4_KEY: 01e9175ca8805cc2137c44eb86184922
```

```
//参数,JSON结构
{
	"app_id": "63163c7b40f2abee", //必选
	"dev_token": "", //ios必选
	"platform": "android", //必选
}
```

### MQTT协议

#### connect 登录

```	
//认证方式
username: app_id
password: BASE64(RC4(app_id:app_secret)
```

#### ping 心跳

#### publish 通知ack,或退出

* 通知ack
```
//参数,JSON结构
{
	"kind": 1, //必选
	"content": {
		"msg_id":"45465464", //必选
	},
}
```

* 退出
```
{
	"kind": 2,//必选
	"content": {
		"app_id":"63163c7b40f2abee",//必选
		"reg_id":"",//必选
	},
}
```

### 开源工具

	1. NSQ 异步中间件
	2. GIN WEB框架
	3. Redis 保存session
	4. MongoDB 数据库
	


