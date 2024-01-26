package admin

import (
	"Tally/common"
	"Tally/global"
	"fmt"
	"github.com/labstack/echo/v4"
	"strings"
)

// AddResource 给所有的角色添加权限可以访问的资源
func AddResource(c echo.Context) error {
	user := c.FormValue("user")
	method := strings.ToUpper(c.FormValue("method"))
	path := c.FormValue("path")
	if method == "" || path == "" || user == "" {
		return common.Fail(c, global.RootCode, global.QueryErr)
	}
	val := global.Global.Redis.SMembers(global.Global.Ctx, global.Role+user).Val()
	if len(val) == 0 {
		return common.Fail(c, global.RootCode, global.QueryErr)
	}
	//分配资源
	for i := 0; i < len(val); i++ {
		_, err := global.Global.CasBin.AddPermissionForUser(val[i], path, method)
		if err != nil {
			return common.Fail(c, global.RootCode, global.AddPermissionFail)
		}
	}
	return common.Ok(c, nil)
}

// AddRolesForUser 为用户分配角色
func AddRolesForUser(c echo.Context) error {
	role := c.FormValue("role")
	user := c.FormValue("user")
	if role == "" || user == "" {
		return common.Fail(c, global.RootCode, global.QueryErr)
	}
	//限制
	//写入redis
	_, err := global.Global.Redis.SAdd(global.Global.Ctx, global.Role+user, role).Result()
	if err != nil {
		return err
	}

	_, err = global.Global.CasBin.AddRoleForUser(user, role)
	if err != nil {
		return common.Fail(c, global.RootCode, global.AddRoleFail)
	}
	return common.Ok(c, nil)
}

// DeletePermissionForUser 删除用户能访问的资源
func DeletePermissionForUser(c echo.Context) error {
	role := c.FormValue("role")
	method := strings.ToUpper(c.FormValue("method"))
	path := strings.ReplaceAll(c.FormValue("path"), " ", "")
	if method == "" || path == "" || role == "" {
		return common.Fail(c, global.RootCode, global.QueryErr)
	}

	//判断是否存在
	if !global.Global.CasBin.HasPolicy(role, path, method) {
		return common.Fail(c, global.RootCode, global.PermissionNotFound)
	}
	_, err := global.Global.CasBin.DeletePermissionForUser(role, path, method)
	if err != nil {
		return common.Fail(c, global.RootCode, global.DelPermissionFail)
	}
	return common.Ok(c, nil)
}

// DeleteRoleForUser 删除用户的角色
func DeleteRoleForUser(c echo.Context) error {
	role := c.FormValue("role")
	user := c.FormValue("user")
	if role == "" || user == "" {
		return common.Fail(c, global.RootCode, global.QueryErr)
	}
	_, err := global.Global.Redis.SRem(global.Global.Ctx, global.Role+user, role).Result()
	if err != nil {
		return err
	}
	_, err = global.Global.CasBin.DeleteRoleForUser(user, role)
	if err != nil {
		return common.Fail(c, global.RootCode, global.DelRoleFail)
	}
	return common.Ok(c, nil)
}

// GetPermissionsForUser 查看用户能访问的资源
func GetPermissionsForUser(c echo.Context) error {
	role := c.FormValue("role")
	if role == "" {
		return common.Fail(c, global.RootCode, global.QueryErr)
	}
	list := global.Global.CasBin.GetPermissionsForUser(role)
	return common.Ok(c, list)
}

// GetAllNamedSubjects 查看所有用户权限
func GetAllNamedSubjects(c echo.Context) error {
	list := global.Global.CasBin.GetNamedPolicy("p")
	return common.Ok(c, list)
}

// UpdatePolicy 修改用户权限
func UpdatePolicy(c echo.Context) error {
	//todo 修改用户权限
	//获取用户的权限
	role := c.FormValue("role")
	path := c.FormValue("path")
	newPath := strings.ReplaceAll(c.FormValue("newPath"), " ", "")
	if role == "" || path == "" || newPath == "" {
		return common.Fail(c, global.RootCode, global.QueryErr)
	}
	fmt.Println(role, path, newPath)
	list := global.Global.CasBin.GetPermissionsForUser(role)
	fmt.Println(list)
	for i := 0; i < len(list); i++ {
		old := list[i]
		if list[i][1] == path {
			list[i][1] = newPath
			_, err := global.Global.CasBin.UpdatePolicy(old, list[i])
			if err != nil {
				return common.Fail(c, global.RootCode, global.UpdatePermissionFail)
			}
		}
	}
	err := global.Global.CasBin.SavePolicy()
	if err != nil {
		return common.Fail(c, global.RootCode, global.UpdatePermissionFail)
	}
	return common.Ok(c, nil)
}
