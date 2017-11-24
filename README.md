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

### 架构

1. session信息按照Hash表结构存放在Redis里面，具体形式为{"app_id":"","reg_id":"","gate_server_ip":"","gate_server_port":""},哈希KEY是app_id:reg_id。

2. 服务端对每个链接都生成一个对应的session对象进行保存相关信息，调用session对象的start方法会开启三个goroutine,一个goroutine负责读数据，一个goroutine负责写数据,另外一个goroutine专门用来检查客户端链接情况，也就是说检查session情况，如果session失效(找不到session)，自动关闭本session，关闭session的顺序是：首先把此session设置成关闭状态，然后直接关闭socket，等待读写Gorotine都退出后，调用closecallback方法来清除在此session上的所有用户在线缓存信息。session自身健康检查的频率是每120秒一次。因为每次客户端心跳来临时会同时更新redis中的session信息和gate服务内存中的最新心跳时间戳，所以自检的方法很简单，速度也很快，直接检查session对象的touchtime就可以了。

3. Redis中的session信息会自动失效的如果没有心跳来更新过期时间，过期时间默认为是300秒。

4. rest_api service只对推送请求进行安全性检查。通过请求合法，参数正确验证后，为了以后消息的追踪，会把消息首先存入数据库，然后直接丢入NSQ队列，让Notifer服务进行处理。

5. notifer service作为NSQ队列的消费者，处理来自rest_api的推送请求。消费处理函数失败返回error时nsq会自动重发。服务本身很简单，具体流程：获取异步消息后根据app_id:reg_id为key查询Redis来获取session信息，如果session有效（存在），根据得到的session获取gate_server_ip和gate_server_port，并通过grpc调用gate　service的推送服务。为了提高性能会利用到grpc链接池。如果session失效（不存在）的话，也意味着用户离线，所以只是打印出日志，完成消费。

6. gate　servic负责和客户端的通信，接收到notifer service的推送请求后会根据得到的app_id:reg_id查询用户在线状态，获得相应的session，把消息直接放到outChannel中等待WriteGoroutine进行读取发送给客户端。gate服务会缓存连接到本服务器用户session信息，查询用户状态时不用调用redis-session的服务。用户断开连接时，也只是清除本地session缓存，并不会主动通知redis-session进行更新，redis-session的每个session都有一个过期失效时间设置。收到客户端发送的message ack消息后会对消息做一次数据库操作来标记消息已经送到。如果没有收到ack消息，不会进行重发，而是在用户登录时作为离线消息发送给用户。消息都有一个失效时间，过了失效时间，将不再发送。

7. 整个服务的设计是作为一个平台来设计的，所以每一个用户都需要通过app_id:reg_id连个参数进行确认。为了减轻服务器压力并减少客户端的好点和流量等，一台手机上多个app只开启一个tcp链接即可，然后多个app的多个用户同用此链接。
