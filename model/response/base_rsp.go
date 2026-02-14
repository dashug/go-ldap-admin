package response

type DashboardList struct {
	DataType  string `json:"dataType"`
	DataName  string `json:"dataName"`
	DataCount int64  `json:"dataCount"`
	Icon      string `json:"icon"`
	Path      string `json:"path"`
}

type BaseConfigRsp struct {
	LdapEnableSync     bool   `json:"ldapEnableSync"`
	DingTalkEnableSync bool   `json:"dingTalkEnableSync"`
	FeiShuEnableSync   bool   `json:"feiShuEnableSync"`
	WeComEnableSync    bool   `json:"weComEnableSync"`
	DirectoryType      string `json:"directoryType"`
	Url                string `json:"url"`
	BaseDN             string `json:"baseDN"`
	AdminDN            string `json:"adminDN"`
	AdminPass          string `json:"adminPass"`
	UserDN             string `json:"userDN"`
	UserInitPassword   string `json:"userInitPassword"`
	DefaultEmailSuffix string `json:"defaultEmailSuffix"`

	DingTalkFlag      string `json:"dingTalkFlag"`
	DingTalkAppKey    string `json:"dingTalkAppKey"`
	DingTalkAppSecret string `json:"dingTalkAppSecret"`
	DingTalkAgentID   string `json:"dingTalkAgentId"`

	WeComFlag       string `json:"weComFlag"`
	WeComCorpID     string `json:"weComCorpId"`
	WeComCorpSecret string `json:"weComCorpSecret"`
	WeComAgentID    int    `json:"weComAgentId"`

	FeiShuFlag      string `json:"feiShuFlag"`
	FeiShuAppID     string `json:"feiShuAppId"`
	FeiShuAppSecret string `json:"feiShuAppSecret"`
}
