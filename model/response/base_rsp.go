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
}
