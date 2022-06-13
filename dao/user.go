package dao

import "bytedance_SimpleDouyin/model"

// GetUser 根据userId获取用户信息
func GetUser(userId uint) (model.User, error) {
	//1.数据模型准备
	var user model.User
	//2.在users表中查对应user_id的user
	if err := db.Model(&model.User{}).Where("id = ?", userId).Find(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}
