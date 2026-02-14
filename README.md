# Go LDAP Admin (Dashug Fork)

基于原项目二次开发的企业目录管理后台，当前聚焦可落地能力：

- 支持 `OpenLDAP` / `Windows AD`
- 支持钉钉 / 企业微信 / 飞书组织与用户同步
- 支持可视化配置向导（目录配置、平台对接、测试连接）

## 致谢

本项目基于 [eryajf/go-ldap-admin](https://github.com/eryajf/go-ldap-admin) 深度二次开发，感谢原作者与全部贡献者。

## 核心功能

- 目录服务双模式：`openldap` / `ad`
- AD 兼容：用户/组模型、DN 规则、成员属性、改密逻辑
- 目录快速配置：可视化配置 LDAP 地址、DN、类型、同步开关
- 第三方平台向导：钉钉/企微/飞书参数配置、测试连接、保存
- 组织与用户同步：平台到 LDAP、SQL 到 LDAP
- 权限与审计：角色权限、接口权限、操作日志

## 接口能力（配置相关）

- `GET /api/base/config`：读取当前系统配置
- `POST /api/base/directoryConfig`：保存目录服务配置
- `POST /api/base/thirdPartyConfig`：保存平台配置
- `POST /api/base/thirdPartyConfig/test`：测试平台连接

## 目录服务配置示例

```yaml
ldap:
  directory-type: "openldap" # openldap / ad
  url: ldap://localhost:389
  base-dn: "dc=example,dc=com"
  admin-dn: "cn=admin,dc=example,dc=com"
  admin-pass: "your-password"
  user-dn: "ou=people,dc=example,dc=com"
  user-init-password: "123456"
  default-email-suffix: "example.com"
  enable-sync: false
```

AD 场景建议使用 `ldaps://...:636`。

## 维护仓库

- 后端：<https://github.com/dashug/go-ldap-admin>
- 前端：<https://github.com/dashug/go-ldap-admin-ui>
