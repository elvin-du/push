package dal

import (
	"push/meta"
)

var (
	DefaultMysqlClient = &Client{DAL: NewMysql(), provider: "mysql"}
	DefaultRedisClient = &Client{DAL: &Redis{}, provider: "redis"}
)

//只是一个包装器，可以选择不同的持久化工具，例如：mysql,redis
type Client struct {
	DAL
	provider string //mysql,redis
}

//上线
func (c *Client) Online(req *meta.DataOnlineRequest) (*meta.DataOnlineResponse, error) {
	return c.Online(req)
}

//下线
func (c *Client) Offline(req *meta.DataOfflineRequest) (*meta.DataOfflineResponse, error) {
	return c.Offline(req)
}

func (c *Client) SaveOfflineMsg(req *meta.SaveOfflineMsgRequest) (*meta.SaveOfflineMsgResponse, error) {
	return c.SaveOfflineMsg(req)
}

func (c *Client) GetOfflineMsgs(req *meta.GetOfflineMsgsRequest) (*meta.GetOfflineMsgsResponse, error) {
	return c.GetOfflineMsgs(req)
}

func (c *Client) DelOfflineMsgs(req *meta.DelOfflineMsgsRequest) (*meta.DelOfflineMsgsResponse, error) {
	return c.DelOfflineMsgs(req)
}

func (c *Client) GetClientInfo(req *meta.GetClientInfoRequest) (*meta.GetClientInfoResponse, error) {
	return c.GetClientInfo(req)
}

func (c *Client) UpdateClientInfo(req *meta.UpdateClientInfoRequest) (*meta.UpdateClientInfoResponse, error) {
	return c.UpdateClientInfo(req)
}
