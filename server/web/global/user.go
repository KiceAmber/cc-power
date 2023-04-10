package global

import (
	"fmt"

	"powerTrading/web/model"
)

var UserList = []*model.User{}

var UserMap = map[string]*model.User{}

// 查询所有用户
func QueryAllUser() []*model.User {
	return UserList
}

// userId 查询用户
func QueryUser(userId string) *model.User {
	if _, exits := UserMap[userId]; !exits {
		fmt.Println("用户不存在.")
		return nil
	}
	return UserMap[userId]
}

// 插入用户
func InsertUser(user *model.User) int {
	tmp := len(UserList)
	UserList = append(UserList, user)
	UserMap[user.UserId] = user
	return len(UserList) - tmp
}

// username 查询用户
func QueryUserByName(username string, password string) *model.User {
	for _, user := range UserList {
		if user.Username == username && user.Password == password {
			return user
		}
	}
	return nil
}
