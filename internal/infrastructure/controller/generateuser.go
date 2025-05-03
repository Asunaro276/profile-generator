package controller

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ryuhei/randomuser-go/internal/config"
	"github.com/ryuhei/randomuser-go/internal/model"
)

// UserResult はランダムユーザーの生成結果
type userResponse struct {
	Results []model.User `json:"results"`
	Info    info         `json:"info"`
}

type info struct {
	Seed    string `json:"seed"`
	Results int    `json:"results"`
	Page    int    `json:"page"`
}

// UserGenerator はユーザー生成インターフェース
type UserGenerator interface {
	Generate(results int, seed int64, page int, gender string) ([]model.User, error)
}

func GenerateUser(c *gin.Context, gen UserGenerator, cfg *config.Config) {
	seed := time.Now().UnixNano()
	if seedParam := c.DefaultQuery("seed", ""); seedParam != "" {
		seedInt, err := strconv.ParseInt(seedParam, 10, 64)
		if err == nil {
			seed = seedInt
		}
	}

	page := 1
	if pageparam := c.DefaultQuery("page", "1"); pageparam != "" {
		pagenum, err := strconv.Atoi(pageparam)
		if err == nil && pagenum > 0 {
			page = pagenum
		}
	}
	seed += int64(page)

	resultsStr := c.DefaultQuery("results", "1")
	results, err := strconv.Atoi(resultsStr)
	if err != nil || results < 1 || results > cfg.MaxResults {
		results = 1
	}

	gender := c.DefaultQuery("gender", "")

	output, err := gen.Generate(results, seed, page, gender)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	res := userResponse{
		Results: output,
		Info: info{
			Seed:    strconv.FormatInt(seed, 10),
			Results: results,
			Page:    page,
		},
	}

	c.Header("Content-Type", "application/json; charset=utf-8")
	c.JSON(http.StatusOK, res)
}
