package dao

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/jmoiron/sqlx"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

const DRIVER = "mysql"

var SqlSession *gorm.DB

var DB *sqlx.DB

type Config struct {
	Url      string `yaml:"url"`
	UserName string `yaml:"userName"`
	PassWord string `yaml:"passWord"`
	DbName   string `yaml:"dbname"`
	Port     string `yaml:"post"`
}

func InitDB() {
	database, err := sqlx.Open("mysql", "root:1040422188@tcp(127.0.0.1:3306)/douyin")
	if err != nil {
		fmt.Println("数据库初始化失败:" + err.Error())
		return
	}
	DB = database
}

func (c *Config) GetConfig() *Config {
	//读取yaml配置文件
	YamlFile, err := ioutil.ReadFile("resource/appliaction.yaml")
	if err != nil {
		fmt.Println(err.Error())
	}
	err = yaml.Unmarshal(YamlFile, c)
	if err != nil {
		fmt.Println(err.Error())
	}

	return c
}

// InitMysql 初始化mysql连接
func InitMysql() (err error) {
	var c Config
	conf := c.GetConfig()

	DBUrl := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		conf.UserName,
		conf.PassWord,
		conf.Url,
		conf.Port,
		conf.DbName,
	)

	sqlSession, err := gorm.Open(DRIVER, DBUrl)
	if err != nil {
		fmt.Println(err.Error())
	}
	//正常链接则返回
	return sqlSession.DB().Ping()

}

func DBClose() {
	err := SqlSession.Close()
	if err != nil {
		return
	}

}
func CloseDB() {
	DB.Close()
}
