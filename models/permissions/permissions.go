package models

import "strings"

const (
	PermissionRoot   = "root"
	PermissionSystem = "system"
	PermissionAdmin  = "admin"
	PermissionUser   = "user"

	permissionSeparator = ";"
)

func GeneratePermissions(permissions ...string) string {
	return strings.Join(permissions, permissionSeparator)
}

func ParsePermissions(permissions string) []string {
	return strings.Split(permissions, permissionSeparator)
}

func PermissionContain(permissions, requiredPermission string) bool {
	for _, permission := range strings.Split(permissions, permissionSeparator) {
		if permission == requiredPermission {
			return true
		}
	}
	return false
}

func UpsertPermissions(p1 []string, p2 []string) []string {
	resultMap := map[string]bool{}
	for _, v := range p1 {
		resultMap[v] = true
	}

	for _, v := range p2 {
		resultMap[v] = true
	}

	keys := make([]string, 0, len(resultMap))
	for k := range resultMap {
		keys = append(keys, k)
	}
	return keys
}
