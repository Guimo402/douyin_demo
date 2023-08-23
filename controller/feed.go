package controller

import (
	"github.com/gin-gonic/gin"
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
	db, err := dbinit()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	db.AutoMigrate(&Video{})

	var videos []Video
	err = db.Order("id desc").Limit(30).Find(&videos).Error
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, FeedResponse{
		Response:  Response{StatusCode: 0},
		VideoList: videos,
		NextTime:  time.Now().Unix(),
	})
}
