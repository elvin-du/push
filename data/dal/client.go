package dal

import (
	"push/meta"
)

var (
	DefaultMysqlClient = &Client{dal: NewMysql(), provider: "mysql"}
	DefaultRedisClient = &Client{dal: &Redis{}, provider: "redis"}
)

//只是一个包装器，可以选择不同的持久化工具，例如：mysql,redis
type Client struct {
	dal      DAL
	provider string //mysql,redis
}

//上线
func (c *Client) Online(req *meta.DataOnlineRequest) (*meta.DataOnlineResponse, error) {
	return c.dal.Online(req)
}

//下线
func (c *Client) Offline(req *meta.DataOfflineRequest) (*meta.DataOfflineResponse, error) {
	return c.dal.Offline(req)
}

func (c *Client) SaveOfflineMsg(req *meta.SaveOfflineMsgRequest) (*meta.SaveOfflineMsgResponse, error) {
	return c.dal.SaveOfflineMsg(req)
}

func (c *Client) GetOfflineMsgs(req *meta.GetOfflineMsgsRequest) (*meta.GetOfflineMsgsResponse, error) {
	return c.dal.GetOfflineMsgs(req)
}

func (c *Client) DelOfflineMsgs(req *meta.DelOfflineMsgsRequest) (*meta.DelOfflineMsgsResponse, error) {
	return c.dal.DelOfflineMsgs(req)
}

func (c *Client) GetClientInfo(req *meta.GetClientInfoRequest) (*meta.GetClientInfoResponse, error) {
	return c.dal.GetClientInfo(req)
}

func (c *Client) UpdateClientInfo(req *meta.UpdateClientInfoRequest) (*meta.UpdateClientInfoResponse, error) {
	return c.dal.UpdateClientInfo(req)
}
