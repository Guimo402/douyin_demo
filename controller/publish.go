package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
	"bytes"
	"github.com/disintegration/imaging"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"log"
	"os"
)

type VideoListResponse struct {
	Response
	VideoList []Video `json:"video_list"`
}

func GetSnapshot(videoPath, snapshotPath string, frameNum int) (string, error) {
	buf := bytes.NewBuffer(nil)
	err := ffmpeg.Input(videoPath).
		Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", frameNum)}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(buf, os.Stdout).
		Run()
	if err != nil {
		log.Fatal("generate failed:", err)
		return "", err
	}
	img, err := imaging.Decode(buf)
	if err != nil {
		log.Fatal("generate failed:", err)
		return "", err
	}
	
	coverName := filepath.Base(snapshotPath)
	coverName = coverName[:len(coverName)-len(filepath.Ext(coverName))]
	
	newCoverName := coverName + ".png"
	
	err = imaging.Save(img, filepath.Join("./public/videos", newCoverName))
	if err != nil {
		log.Fatal("generate failed:", err)
		return "", err
	}
	return filepath.Join("static/images", newCoverName), nil
}

// Publish check token then save upload file to public directory
func Publish(c *gin.Context) {
	token := c.PostForm("token")
	db, err := dbinit()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	db.AutoMigrate(&User{})
	db.AutoMigrate(&Video{})

	var user User
	if err := db.Find(&user, "token = ?", token).Error; err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}

	data, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	filename := filepath.Base(data.Filename)
	finalName := fmt.Sprintf("%d_%s", user.Id, filename)
	saveFile := filepath.Join("./public/", finalName)
	if err := c.SaveUploadedFile(data, saveFile); err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	coverPath := filepath.Join("./public/", finalName)

	var cover_path string

	cover_path, err = GetSnapshot(saveFile, coverPath, 1)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	new_video := Video{
		AuthorID: user.Id,
		PlayUrl:  "http://39.101.1.113:8080/static/" + finalName,
		CoverUrl: "http://39.101.1.113:8080/" + cover_path,
	}

	err = db.Create(&new_video).Error
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  finalName + " uploaded successfully",
	})
}

// PublishList all users have same publish video list
func PublishList(c *gin.Context) {
	user_id := c.Query("user_id")

	db, err := dbinit()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var videos []Video
	var user User
	db.Find(&user, "id = ?", user_id)
	err = db.Find(&videos, "author_id = ?", user_id).Error
	if err != nil {
		panic(err)
	}
	
	for i := range videos {
		videos[i].Author = user
	}

	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: videos,
	})
}
