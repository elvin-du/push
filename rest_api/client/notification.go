package client

type Notification struct {
	Alert               string                 `json:"alert"`
	Audience            string                 `json:"audience,omitempty"` //all,reg_id
	TTL                 uint64                 `json:"ttl,omitempty"`      //sec
	Extras              map[string]interface{} `json:"extras,omitempty"`
	IosNotification     *IosNotification       `json:"ios,omitempty"`
	AndroidNotification *AndroidNotification   `json:"android,omitempty"`
}

func (n *Notification) AddExtra(key string, value interface{}) {
	if n.Extras == nil {
		n.Extras = make(map[string]interface{})
	}
	n.Extras[key] = value
}

type IosNotification struct {
	Sound      string                 `json:"sound,omitempty"`
	Badge      uint64                 `json:"badge,omitempty"`
	Production bool                   `json:"production"`
	Extras     map[string]interface{} `json:"extras,omitempty"`
}

func (in *IosNotification) AddExtra(key string, value interface{}) {
	if in.Extras == nil {
		in.Extras = make(map[string]interface{})
	}
	in.Extras[key] = value
}

type AndroidNotification struct {
	Extras map[string]interface{} `json:"extras,omitempty"`
}

func (an *AndroidNotification) AddExtra(key string, value interface{}) {
	if an.Extras == nil {
		an.Extras = make(map[string]interface{})
	}
	an.Extras[key] = value
}
