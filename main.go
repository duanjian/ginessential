package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type User struct {
	gorm.Model
	Name string `gorm:"type:varchar(20);not null"`
	Telephone string `gorm:"type:varchar(11); not null"`
	Password string `gorm:"type:varchar(32);not null"`
}

func main() {

	db := InitDB()
	defer  db.Close()

	r := gin.Default()
	r.POST("/api/auth/register", func(ctx *gin.Context) {
		//获取参数
		name := ctx.PostForm("name")
		telephone := ctx.PostForm("telephone")
		password := ctx.PostForm("password")

		//数据验证
		if len(telephone) != 11 {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{ "code": 200, "msg": "手机号必须为11位"})
			return
		}

		if len(password) < 6 {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 200, "msg": "密码不能少于6位"})
			return
		}

		if len(name) == 0 {
			name = RandomString(10)
		}

		log.Println(name, telephone, password)
		
		//判断手机号
		if isTelephoneExist(db, telephone) {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 200, "msg": "手机号已存在"})
			return
		}
		
		//创建用户
		newUser := User{
			Name:      name,
			Telephone: telephone,
			Password:  password,
		}

		db.Create(&newUser)
		

		//返回结果
		ctx.JSON(200, gin.H{
			"message": "注册成功",
		})
	})

	r.Run()
	fmt.Println("hello world")

}

func isTelephoneExist(db *gorm.DB, telephone string) bool {
	var user User
	db.Where("telephone = ?", telephone).First(&user)
	if user.ID != 0 {
		return true
	}

	return false
}

func RandomString(n int)  string {
	var letters = []byte("asdfghjklqwertyuiopASDFGHJKLQWERTYUIOP")
	result := make([]byte, n)

	rand.Seed(time.Now().Unix())

	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}

	return string(result)
}

func InitDB() *gorm.DB {
	driverName := "mysql"
	host := "localhost"
	port := "3306"
	database := "ginessential"
	username := "root"
	password := "sasasa"
	charset := "utf8"
	args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true",
		username,
		password,
		host,
		port,
		database,
		charset)

	db, err := gorm.Open(driverName, args)
	if err != nil {
		panic("failed to connect databse, err: "+ err.Error())
	}

	db.AutoMigrate(&User{})

	return db
}