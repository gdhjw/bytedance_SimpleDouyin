package service

import (
	"bytedance_SimpleDouyin/common"
	"bytedance_SimpleDouyin/dao"
	"bytedance_SimpleDouyin/model"
	"fmt"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

const (
	MaxUsernameLength = 32 //用户名最大长度
	MaxPasswordLength = 32 //密码最大长度
	MinPasswordLength = 6  //密码最小长度
)

func CreateRegisterUser(userName string, passWord string) (model.User, error) {
	//1.Following数据模型准备
	newPassword, _ := HashAndSalt(passWord)
	newUser := model.User{
		Name:     userName,
		Password: newPassword,
	}
	//2.模型关联到数据库表users //可注释
	dao.SqlSession.AutoMigrate(&model.User{})
	//3.新建user
	if IsUserExistByName(userName) {
		//用户已存在
		return newUser, common.ErrorUserExit
	} else {
		//用户不存在，新建用户
		if err := dao.SqlSession.Model(&model.User{}).Create(&newUser).Error; err != nil {
			//错误处理
			//fmt.Println(err)
			panic(err)
			return newUser, err
		}
	}
	return newUser, nil
}

// CreateRegisterUser2 注册用户
func CreateRegisterUser2(username string, password string) bool {
	newPassword, err := HashAndSalt(password)
	registerSQL, err := dao.DB.Exec("insert into user(username,password)values(?,?)", username, newPassword)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	id, err := registerSQL.LastInsertId()
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	fmt.Println("注册成功:", id)

	return true

}

func FindPasswordByName(username string) string {
	var user model.Users

	//查询操作
	rows := dao.DB.QueryRowx("select * from user where username=? ", strings.TrimSpace(username))

	rows.Scan(&user)
	fmt.Println(username)
	fmt.Println(user)

	return user.Password
}

// FindIdByName 通过username查找id
func FindIdByName(username string) uint {
	//构建数据存储
	var user model.Users

	//查询操作
	err := dao.DB.Select(&user, "select * from user where username=? ", username)
	if err != nil {
		fmt.Println("查询失败2")
		return 0
	}
	return user.Id

}

func IsUserExistByName(name string) bool {
	var userExist = &model.User{}
	if err := dao.SqlSession.Model(&model.User{}).Where("name=?", name).First(&userExist).Error; gorm.IsRecordNotFoundError(err) {
		//关注不存在
		return false
	}
	//关注存在
	return true
}

// IsUserExist 检查登录用户是否存在
func IsUserExist(username string, password string, login *model.User) error {
	if login == nil {
		return common.ErrorNullPointer
	}
	dao.SqlSession.Where("name=?", username).First(login)
	if !ComparePasswords(login.Password, password) {
		return common.ErrorPasswordFalse
	}
	if login.Model.ID == 0 {
		return common.ErrorPasswordNull
	}
	return nil

}

// GetUserById 通过id查询用户信息
func GetUserById(UserId uint) (model.User, error) {
	var user model.User
	if err := dao.SqlSession.Model(&model.User{}).Where("id=?", UserId).Find(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}

func GetUserById2(UserId uint) (model.Users, error) {
	var user model.Users

	err := dao.DB.Select(&user, "select * from user where id=? ", UserId)
	if err != nil {
		fmt.Println("查询失败3")
		return user, err
	}
	return user, nil

}

// GetUserByIdInfo 通过id检索用户信息,userinfo
func GetUserByIdInfo(UserId uint, user *model.User) error {
	if user == nil {
		return common.ErrorNullPointer
	}
	dao.SqlSession.Where("id=?", UserId).First(user)
	return nil
}

// ISUserLegal 用户输入合法性检验
func ISUserLegal(username string, password string) error {
	//用户名输入检查
	//用户名为空
	if username == "" {
		return common.ErrorUserNameNull
	}
	//用户名过长
	if len(username) > MaxUsernameLength {
		return common.ErrorPasswordLength
	}
	//密码为空
	if password == "" {
		return common.ErrorPasswordNull
	}
	//密码长度不合法
	if len(password) > MaxPasswordLength || len(password) < MinPasswordLength {
		return common.ErrorPasswordLength
	}
	return nil

}

// HashAndSalt MD5加盐哈希
/*func HashAndSalt(pwdStr string, saltStr string) (pwdHash string) {

	//密码
	pwd := []byte(pwdStr)
	//盐值
	salt := []byte(saltStr)
	md5pwd := md5.New()
	//先写盐值
	md5pwd.Write(salt)
	md5pwd.Write(pwd)
	pwdHash = hex.EncodeToString(md5pwd.Sum(nil))
	return pwdHash
}*/
func HashAndSalt(pwdStr string) (pwdHash string, err error) {
	pwd := []byte(pwdStr)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		return
	}
	pwdHash = string(hash)
	return
}

// ComparePasswords 验证密码
func ComparePasswords(hashedPwd string, plainPwd string) bool {
	byteHash := []byte(hashedPwd)
	bytePwd := []byte(plainPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, bytePwd)
	if err != nil {
		return false
	}
	return true
}
