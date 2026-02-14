<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**目录**

- [Go LDAP Admin (Dashug Fork)](#go-ldap-admin-dashug-fork)
  - [致谢](#%E8%87%B4%E8%B0%A2)
  - [核心功能](#%E6%A0%B8%E5%BF%83%E5%8A%9F%E8%83%BD)
  - [接口能力（配置相关）](#%E6%8E%A5%E5%8F%A3%E8%83%BD%E5%8A%9B%E9%85%8D%E7%BD%AE%E7%9B%B8%E5%85%B3)
  - [目录服务配置示例](#%E7%9B%AE%E5%BD%95%E6%9C%8D%E5%8A%A1%E9%85%8D%E7%BD%AE%E7%A4%BA%E4%BE%8B)
  - [配置流程说明（小白版）](#%E9%85%8D%E7%BD%AE%E6%B5%81%E7%A8%8B%E8%AF%B4%E6%98%8E%E5%B0%8F%E7%99%BD%E7%89%88)
  - [图文功能说明](#%E5%9B%BE%E6%96%87%E5%8A%9F%E8%83%BD%E8%AF%B4%E6%98%8E)
  - [维护仓库](#%E7%BB%B4%E6%8A%A4%E4%BB%93%E5%BA%93)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

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

## 配置流程说明（小白版）

### 1. 配置目录服务（OpenLDAP/AD）

1. 登录系统后进入 `人员管理 -> 用户管理`。
2. 点击 `目录快速配置`。
3. 选择目录类型：`OpenLDAP` 或 `Windows AD`。
4. 填写 `LDAP地址`、`Base DN`、`管理员 DN`、`用户 OU DN` 等字段。
5. 点击保存。

建议：

- AD 场景优先使用 `ldaps://域名:636`。
- `管理员密码` 可留空，表示不覆盖当前密码。

### 2. 配置第三方平台（钉钉/企微/飞书）

1. 进入 `人员管理 -> 用户管理`。
2. 点击 `平台对接向导`。
3. 选择平台标签页，填写平台凭证（如 `AppKey/AppSecret` 或 `CorpId/CorpSecret`）。
4. 先点击 `测试连接`。
5. 连接成功后点击 `保存`。

### 3. 首次同步建议顺序

1. 先同步部门（避免用户找不到部门）。
2. 再同步用户。
3. 最后检查用户和分组列表中的同步状态。

## 图文功能说明

### 1. 功能总览图（首页 / 用户管理 / 分组管理）

![首页](https://cdn.jsdelivr.net/gh/eryajf/tu/img/image_20220724_165545.png)
![用户管理](https://cdn.jsdelivr.net/gh/eryajf/tu/img/image_20220724_165623.png)
![分组管理](https://cdn.jsdelivr.net/gh/eryajf/tu/img/image_20220724_165701.png)

### 2. 目录快速配置（步骤图 1-2-3）

```mermaid
flowchart LR
  A["进入 用户管理"] --> B["点击 目录快速配置"]
  B --> C["选择目录类型(OpenLDAP/AD)"]
  C --> D["填写地址与DN参数"]
  D --> E["点击保存"]
```

### 3. 平台对接向导（钉钉 / 企微 / 飞书）

```mermaid
flowchart LR
  A["点击 平台对接向导"] --> B["选择平台标签页"]
  B --> C["填写平台凭证"]
  C --> D["测试连接"]
  D --> E["保存配置"]
  E --> F["执行同步任务"]
```

### 4. 测试连接与保存结果示例

测试连接成功返回示例：

```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "platform": "feishu",
    "ok": true
  }
}
```

保存配置成功返回示例：

```json
{
  "code": 0,
  "msg": "success",
  "data": null
}
```

### 5. 常见错误提示（可选）

- `钉钉连接测试失败`：检查 `AppKey/AppSecret`、应用权限。
- `企微连接测试失败`：检查 `CorpId/CorpSecret/AgentId`。
- `飞书连接测试失败`：检查 `AppId/AppSecret`、通讯录权限。
- `LDAP连接异常`：检查 `url/base-dn/admin-dn`，AD 场景优先 `ldaps://`。

## 维护仓库

- 后端：<https://github.com/dashug/go-ldap-admin>
- 前端：<https://github.com/dashug/go-ldap-admin-ui>

## 部署说明

### 快速部署（推荐，Docker Compose）

后端仓库已提供可直接使用的编排文件：

- `/docs/docker-compose/docker-compose.yaml`

执行步骤：

```bash
cd docs/docker-compose
docker compose up -d
```

默认端口：

- 后端 API：`http://<你的服务器IP>:8888`

说明：

- 编排中默认带 OpenLDAP 容器，适合快速试用。
- 如你使用外部 LDAP/AD，请修改编排里的 `config` 内容或挂载自定义 `config.yml`。

### 本地开发运行（源码）

```bash
go mod download
cp config.yml config.local.yml  # 可选，建议保留一份本地配置
go run main.go
```

或使用 Makefile：

```bash
make run
```
