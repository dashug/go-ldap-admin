package ildap

import (
	"fmt"
	"unicode/utf16"

	"github.com/eryajf/go-ldap-admin/config"
	"github.com/eryajf/go-ldap-admin/model"
	"github.com/eryajf/go-ldap-admin/public/common"
	"github.com/eryajf/go-ldap-admin/public/tools"

	ldap "github.com/go-ldap/ldap/v3"
)

type UserService struct{}

// 创建资源
func (x UserService) Add(user *model.User) error {
	// 获取 LDAP 连接
	conn, err := common.GetLDAPConn()
	defer common.PutLADPConn(conn)
	if err != nil {
		return err
	}

	if IsActiveDirectory() {
		return x.addADUser(conn, user)
	}

	add := ldap.NewAddRequest(user.UserDN, nil)
	add.Attribute("objectClass", []string{"inetOrgPerson"})
	add.Attribute("cn", []string{user.Username})
	add.Attribute("sn", []string{user.Nickname})
	add.Attribute("businessCategory", []string{user.Departments})
	add.Attribute("departmentNumber", []string{user.Position})
	add.Attribute("description", []string{user.Introduction})
	add.Attribute("displayName", []string{user.Nickname})
	add.Attribute("mail", []string{user.Mail})
	add.Attribute("employeeNumber", []string{user.JobNumber})
	add.Attribute("givenName", []string{user.GivenName})
	add.Attribute("postalAddress", []string{user.PostalAddress})
	add.Attribute("mobile", []string{user.Mobile})
	add.Attribute("uid", []string{user.Username})
	var pass string
	if config.Conf.Ldap.UserPasswordEncryptionType == "clear" {
		pass = tools.NewParPasswd(user.Password)
	} else {
		pass = tools.EncodePass([]byte(tools.NewParPasswd(user.Password)))
	}
	add.Attribute("userPassword", []string{pass})

	return conn.Add(add)
}

func (x UserService) addADUser(conn *ldap.Conn, user *model.User) error {
	if user.UserDN == "" {
		user.UserDN = BuildUserDN(user.Username)
	}
	sn := user.Nickname
	if sn == "" {
		sn = user.Username
	}
	add := ldap.NewAddRequest(user.UserDN, nil)
	add.Attribute("objectClass", []string{"top", "person", "organizationalPerson", "user"})
	add.Attribute("cn", []string{user.Username})
	add.Attribute("sn", []string{sn})
	add.Attribute("sAMAccountName", []string{user.Username})
	if user.Mail != "" {
		add.Attribute("userPrincipalName", []string{user.Mail})
	} else {
		add.Attribute("userPrincipalName", []string{buildADUPN(user.Username)})
	}
	if user.Nickname != "" {
		add.Attribute("displayName", []string{user.Nickname})
	}
	if user.GivenName != "" {
		add.Attribute("givenName", []string{user.GivenName})
	}
	if user.Mail != "" {
		add.Attribute("mail", []string{user.Mail})
	}
	if user.Mobile != "" {
		add.Attribute("mobile", []string{user.Mobile})
	}
	if user.Introduction != "" {
		add.Attribute("description", []string{user.Introduction})
	}
	if user.PostalAddress != "" {
		add.Attribute("postalAddress", []string{user.PostalAddress})
	}
	if user.Departments != "" {
		add.Attribute("department", []string{user.Departments})
	}
	if user.Position != "" {
		add.Attribute("title", []string{user.Position})
	}
	if user.JobNumber != "" {
		add.Attribute("employeeID", []string{user.JobNumber})
	}

	if err := conn.Add(add); err != nil {
		return err
	}

	if err := updatePasswordAD(conn, user.UserDN, tools.NewParPasswd(user.Password)); err != nil {
		return err
	}

	enableReq := ldap.NewModifyRequest(user.UserDN, nil)
	enableReq.Replace("userAccountControl", []string{"512"})
	if err := conn.Modify(enableReq); err != nil {
		return err
	}

	return nil
}

// Update 更新资源
func (x UserService) Update(oldusername string, user *model.User) error {
	modify := ldap.NewModifyRequest(user.UserDN, nil)
	modify.Replace("cn", []string{user.Username})
	if IsActiveDirectory() {
		sn := user.Nickname
		if sn == "" {
			sn = user.Username
		}
		modify.Replace("sn", []string{sn})
		modify.Replace("department", []string{user.Departments})
		modify.Replace("title", []string{user.Position})
		modify.Replace("description", []string{user.Introduction})
		modify.Replace("displayName", []string{user.Nickname})
		modify.Replace("mail", []string{user.Mail})
		modify.Replace("employeeID", []string{user.JobNumber})
		modify.Replace("givenName", []string{user.GivenName})
		modify.Replace("postalAddress", []string{user.PostalAddress})
		modify.Replace("mobile", []string{user.Mobile})
		modify.Replace("sAMAccountName", []string{user.Username})
		if user.Mail != "" {
			modify.Replace("userPrincipalName", []string{user.Mail})
		} else {
			modify.Replace("userPrincipalName", []string{buildADUPN(user.Username)})
		}
	} else {
		modify.Replace("sn", []string{oldusername})
		modify.Replace("businessCategory", []string{user.Departments})
		modify.Replace("departmentNumber", []string{user.Position})
		modify.Replace("description", []string{user.Introduction})
		modify.Replace("displayName", []string{user.Nickname})
		modify.Replace("mail", []string{user.Mail})
		modify.Replace("employeeNumber", []string{user.JobNumber})
		modify.Replace("givenName", []string{user.GivenName})
		modify.Replace("postalAddress", []string{user.PostalAddress})
		modify.Replace("mobile", []string{user.Mobile})
	}

	// 获取 LDAP 连接
	conn, err := common.GetLDAPConn()
	defer common.PutLADPConn(conn)
	if err != nil {
		return err
	}

	err = conn.Modify(modify)
	if err != nil {
		return err
	}
	if config.Conf.Ldap.UserNameModify && oldusername != user.Username {
		modifyDn := ldap.NewModifyDNRequest(BuildUserDN(oldusername), fmt.Sprintf("%s=%s", userRDNAttr(), user.Username), true, "")
		return conn.ModifyDN(modifyDn)
	}
	return nil
}

func (x UserService) Exist(filter map[string]any) (bool, error) {
	filter_str := ""
	for key, value := range filter {
		filter_str += fmt.Sprintf("(%s=%s)", key, value)
	}
	search_filter := userSearchFilter(filter_str)
	// Construct query request
	searchRequest := ldap.NewSearchRequest(
		config.Conf.Ldap.BaseDN,                                     // This is basedn, we will start searching from this node.
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false, // Here several parameters are respectively scope, derefAliases, sizeLimit, timeLimit,  typesOnly
		search_filter,  // This is Filter for LDAP query
		[]string{"DN"}, // Here are the attributes returned by the query, provided as an array. If empty, all attributes are returned
		nil,
	)

	// 获取 LDAP 连接
	conn, err := common.GetLDAPConn()
	defer common.PutLADPConn(conn)
	if err != nil {
		return false, err
	}
	var sr *ldap.SearchResult
	// Search through ldap built-in search
	sr, err = conn.Search(searchRequest)
	if err != nil {
		return false, err
	}
	if len(sr.Entries) > 0 {
		return true, nil
	}
	return false, nil
}

// Delete 删除资源
func (x UserService) Delete(udn string) error {
	del := ldap.NewDelRequest(udn, nil)
	// 获取 LDAP 连接
	conn, err := common.GetLDAPConn()
	defer common.PutLADPConn(conn)
	if err != nil {
		return err
	}
	return conn.Del(del)
}

// ChangePwd 修改用户密码，此处旧密码也可以为空，ldap可以直接通过用户DN加上新密码来进行修改
func (x UserService) ChangePwd(udn, oldpasswd, newpasswd string) error {
	if IsActiveDirectory() {
		// AD 通常不支持 RFC3062 Password Modify，这里改为写 unicodePwd。
		conn, err := common.GetLDAPConn()
		defer common.PutLADPConn(conn)
		if err != nil {
			return err
		}
		return updatePasswordAD(conn, udn, newpasswd)
	}

	if config.Conf.Ldap.UserPasswordEncryptionType == "clear" {
		return updatePasswordClear(udn, newpasswd)
	}
	modifyPass := ldap.NewPasswordModifyRequest(udn, oldpasswd, newpasswd)

	// 获取 LDAP 连接
	conn, err := common.GetLDAPConn()
	defer common.PutLADPConn(conn)
	if err != nil {
		return err
	}

	_, err = conn.PasswordModify(modifyPass)
	if err != nil {
		return fmt.Errorf("password modify failed for %s, err: %v", udn, err)
	}
	return nil
}

// NewPwd 新旧密码都是空，通过管理员可以修改成功并返回新的密码
func (x UserService) NewPwd(username string) (string, error) {
	udn := BuildUserDN(username)
	if username == "admin" {
		udn = config.Conf.Ldap.AdminDN
	}
	if IsActiveDirectory() {
		newpass := tools.GenerateRandomPassword()
		conn, err := common.GetLDAPConn()
		defer common.PutLADPConn(conn)
		if err != nil {
			return "", err
		}
		if err := updatePasswordAD(conn, udn, newpass); err != nil {
			return "", fmt.Errorf("password modify failed for %s, err: %v", username, err)
		}
		return newpass, nil
	}
	if config.Conf.Ldap.UserPasswordEncryptionType == "clear" {
		newpass := tools.GenerateRandomPassword()
		if err := updatePasswordClear(udn, newpass); err != nil {
			return "", fmt.Errorf("password modify failed for %s, err: %v", username, err)
		}
		return newpass, nil
	}
	modifyPass := ldap.NewPasswordModifyRequest(udn, "", "")

	// 获取 LDAP 连接
	conn, err := common.GetLDAPConn()
	defer common.PutLADPConn(conn)
	if err != nil {
		return "", err
	}

	newpass, err := conn.PasswordModify(modifyPass)
	if err != nil {
		return "", fmt.Errorf("password modify failed for %s, err: %v", username, err)
	}
	return newpass.GeneratedPassword, nil
}

func updatePasswordClear(udn, newpasswd string) error {
	modify := ldap.NewModifyRequest(udn, nil)
	modify.Replace("userPassword", []string{newpasswd})

	// 获取 LDAP 连接
	conn, err := common.GetLDAPConn()
	defer common.PutLADPConn(conn)
	if err != nil {
		return err
	}

	if err := conn.Modify(modify); err != nil {
		return fmt.Errorf("password modify failed for %s, err: %v", udn, err)
	}
	return nil
}

func updatePasswordAD(conn *ldap.Conn, udn, newpasswd string) error {
	modify := ldap.NewModifyRequest(udn, nil)
	modify.Replace("unicodePwd", []string{encodeADPassword(newpasswd)})
	if err := conn.Modify(modify); err != nil {
		return fmt.Errorf("password modify failed for %s, err: %v", udn, err)
	}
	return nil
}

func encodeADPassword(password string) string {
	quoted := fmt.Sprintf("\"%s\"", password)
	encoded := utf16.Encode([]rune(quoted))
	buf := make([]byte, len(encoded)*2)
	for i, v := range encoded {
		buf[i*2] = byte(v)
		buf[i*2+1] = byte(v >> 8)
	}
	return string(buf)
}

func (x UserService) ListUserDN() (users []*model.User, err error) {
	// Construct query request
	searchRequest := ldap.NewSearchRequest(
		config.Conf.Ldap.BaseDN,                                     // This is basedn, we will start searching from this node.
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false, // Here several parameters are respectively scope, derefAliases, sizeLimit, timeLimit,  typesOnly
		userListFilter(), // This is Filter for LDAP query
		[]string{"DN"},   // Here are the attributes returned by the query, provided as an array. If empty, all attributes are returned
		nil,
	)

	// 获取 LDAP 连接
	conn, err := common.GetLDAPConn()
	defer common.PutLADPConn(conn)
	if err != nil {
		return users, err
	}
	var sr *ldap.SearchResult
	// Search through ldap built-in search
	sr, err = conn.Search(searchRequest)
	if err != nil {
		return nil, err
	}
	if len(sr.Entries) > 0 {
		for _, v := range sr.Entries {
			users = append(users, &model.User{
				UserDN: v.DN,
			})
		}
	}
	return
}
