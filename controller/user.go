package controller

import (
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"net/http"
	"sync/atomic"
	"strings"
	"fmt"
)

// usersLoginInfo use map to store user info, and key is username+password for demo
// user data will be cleared every time the server starts
// test data: username=zhanglei, password=douyin
var usersLoginInfo = map[string]User{
	"zhangleidouyin": {
		Id:            1,
		Name:          "zhanglei",
		FollowCount:   10,
		FollowerCount: 5,
		IsFollow:      true,
	},
}

func getUsersSum() int {
	db, err := dbinit()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	
	db.AutoMigrate(&User{})
	var count int
	err = db.Model(&User{}).Count(&count).Error
	if err != nil {
		fmt.Println("exec failed", err)
	}
	return count
}

var userIdSequence = int64(getUsersSum())


type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	Response
	User User `json:"user"`
}

func Register(c *gin.Context) {
	db, err := dbinit()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	username := c.Query("username")
	password := c.Query("password")

	if strings.Contains(username,"#") {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "Sorry, '#' could not be used in username"},
		})
		return
	}

	token := username + "#" + password

	db.AutoMigrate(&User{})
	
	var f User
	err = db.Find(&f, "token=?", token).Error
	if err == nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User already exist"},
		})
	} else {
		atomic.AddInt64(&userIdSequence, 1)
		newUser := User{
			Id:    userIdSequence,
			Name:  username,
			Token: token,
		}
		db.Create(&newUser)
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   userIdSequence,
			Token:    token,
		})
	}
}

func Login(c *gin.Context) {
	db, err := dbinit()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	username := c.Query("username")
	password := c.Query("password")

	if strings.Contains(username,"#") {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "Sorry, '#' could not be used in username"},
		})
		return
	}

	token := username + "#" + password

	db.AutoMigrate(&User{})

	var u User
	err = db.Find(&u, "token=?", token).Error

	if err == nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   u.Id,
			Token:    token,
		})
	} else {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
	}
}

func UserInfo(c *gin.Context) {
	db, err := dbinit()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	token := c.Query("token")

	db.AutoMigrate(&User{})

	var u User
	err = db.Find(&u, "token=?", token).Error
	if err == nil {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 0},
			User:     u,
		})
	} else {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
	}
}
