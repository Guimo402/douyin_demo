package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"net/http"
	"time"
)

type FeedResponse struct {
	Response
	VideoList []Video `json:"video_list,omitempty"`
	NextTime  int64   `json:"next_time,omitempty"`
}

// Feed same demo video list for every request
func Feed(c *gin.Context) {
	db, err := gorm.Open("mysql", "guest:guest123@/new?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}

	defer db.Close()
	var videos []Video
	err = db.Order("id desc").Limit(3).Find(&videos).Error
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, FeedResponse{
		Response:  Response{StatusCode: 0},
		VideoList: videos,
		NextTime:  time.Now().Unix(),
	})
}
