package types

type Message struct {
	ID      string                 `json:"id" bson:"_id"`
	AppID   string                 `json:"app_id" bson:"app_id"`
	RegID   string                 `json:"reg_id" bson:"reg_id"`
	Content string                 `json:"content"`
	TTL     uint64                 `json:"ttl" bson:"ttl"`
	Extras  map[string]interface{} `json:"extras"`
}
