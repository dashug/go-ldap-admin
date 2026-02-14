package request

// BaseSendCodeReq 发送验证码
type BaseSendCodeReq struct {
	Mail string `json:"mail" validate:"required,min=0,max=100"`
}

// BaseChangePwdReq 修改密码结构体
type BaseChangePwdReq struct {
	Mail string `json:"mail" validate:"required,min=0,max=100"`
	Code string `json:"code" validate:"required,len=6"`
}

// BaseDashboardReq  系统首页展示数据结构体
type BaseDashboardReq struct {
}

// EncryptPasswdReq
type EncryptPasswdReq struct {
	Passwd string `json:"passwd" form:"passwd" validate:"required"`
}

// DecryptPasswdReq
type DecryptPasswdReq struct {
	Passwd string `json:"passwd" form:"passwd" validate:"required"`
}

// BaseConfigReq 获取系统配置结构体
type BaseConfigReq struct {
}

// BaseVersionReq 获取版本信息结构体
type BaseVersionReq struct {
}

// BaseUpdateDirectoryConfigReq 更新目录服务配置
type BaseUpdateDirectoryConfigReq struct {
	DirectoryType      string `json:"directoryType" validate:"required,oneof=openldap ad"`
	Url                string `json:"url" validate:"required,min=1,max=255"`
	BaseDN             string `json:"baseDN" validate:"required,min=1,max=255"`
	AdminDN            string `json:"adminDN" validate:"required,min=1,max=255"`
	AdminPass          string `json:"adminPass" validate:"omitempty,min=0,max=255"`
	UserDN             string `json:"userDN" validate:"required,min=1,max=255"`
	UserInitPassword   string `json:"userInitPassword" validate:"required,min=1,max=255"`
	DefaultEmailSuffix string `json:"defaultEmailSuffix" validate:"required,min=1,max=100"`
	LdapEnableSync     bool   `json:"ldapEnableSync"`
}
