package ildap

import (
	"fmt"
	"strings"

	"github.com/eryajf/go-ldap-admin/config"
)

func IsActiveDirectory() bool {
	return strings.EqualFold(strings.TrimSpace(config.Conf.Ldap.DirectoryType), "ad")
}

func BuildUserDN(username string) string {
	return fmt.Sprintf("%s=%s,%s", userRDNAttr(), username, config.Conf.Ldap.UserDN)
}

func IsOUGroupDN(dn string) bool {
	return strings.HasPrefix(strings.ToLower(strings.TrimSpace(dn)), "ou=")
}

func userRDNAttr() string {
	if IsActiveDirectory() {
		return "cn"
	}
	return "uid"
}

func groupMemberAttr() string {
	if IsActiveDirectory() {
		return "member"
	}
	return "uniqueMember"
}

func userSearchFilter(filter string) string {
	if IsActiveDirectory() {
		return fmt.Sprintf("(&(objectClass=user)%s)", filter)
	}
	return fmt.Sprintf("(&(|(objectClass=inetOrgPerson)(objectClass=simpleSecurityObject))%s)", filter)
}

func userListFilter() string {
	if IsActiveDirectory() {
		return "(objectClass=user)"
	}
	return "(|(objectClass=inetOrgPerson)(objectClass=simpleSecurityObject))"
}

func groupListFilter() string {
	if IsActiveDirectory() {
		return "(|(objectClass=organizationalUnit)(objectClass=group))"
	}
	return "(|(objectClass=organizationalUnit)(objectClass=groupOfUniqueNames))"
}

func buildADUPN(username string) string {
	domainParts := make([]string, 0)
	for _, p := range strings.Split(config.Conf.Ldap.BaseDN, ",") {
		kv := strings.SplitN(strings.TrimSpace(p), "=", 2)
		if len(kv) == 2 && strings.EqualFold(kv[0], "dc") && kv[1] != "" {
			domainParts = append(domainParts, kv[1])
		}
	}
	if len(domainParts) == 0 {
		return username
	}
	return fmt.Sprintf("%s@%s", username, strings.Join(domainParts, "."))
}
