package services

import (
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"short-url/db"
	"short-url/models"
	"short-url/utils"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func ResolveUrl(ctx *gin.Context) {
	url := ctx.Params.ByName("url")

	redisClient := db.CreateClient(0)

	defer redisClient.Close()

	val, err := redisClient.Get(db.Context, url).Result()

	if err == redis.Nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "not found",
		})
		return
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "something went wrong",
		})
	}

	rateLimitDB := db.CreateClient(1)

	rateLimitDB.Incr(db.Context, "counter")

	defer rateLimitDB.Close()

	ctx.Redirect(301, val)

}

func AddUrl(ctx *gin.Context) {
	var url models.Url

	err := ctx.BindJSON(&url)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid data provided",
		})
		return
	}

	rateLimitDB := db.CreateClient(1)

	defer rateLimitDB.Close()

	val, err := rateLimitDB.Get(db.Context, ctx.ClientIP()).Result()
	limit, _ := rateLimitDB.TTL(db.Context, ctx.ClientIP()).Result()

	if err == redis.Nil {
		_ = rateLimitDB.Set(db.Context, ctx.ClientIP(), os.Getenv("API_QUOTA"), 30*60*time.Second).Err()
	} else if err == nil {
		valInt, _ := strconv.Atoi(val)
		if valInt <= 0 {
			ctx.JSON(http.StatusServiceUnavailable, gin.H{
				"error":            "Rate limit exceeded",
				"rate_limit_reset": limit / time.Nanosecond / time.Hour,
			})
			return
		}
	}

	if !utils.RemoveDomainError(url.Url) {
		ctx.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "Inappropriate url",
		})
		return
	}

	url.Url = utils.EnforceHTTP(url.Url)

	var id string

	if url.CustomShort != "" {
		id = url.CustomShort
	} else {
		id = utils.Base62Encode(rand.Uint64())
	}

	redisUrlDB := db.CreateClient(0)
	defer redisUrlDB.Close()

	val, _ = redisUrlDB.Get(db.Context, id).Result()

	if val != "" {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "Custom URL already in use",
		})
		return
	}

	if url.Expiry == 0 {
		url.Expiry = 24
	}

	err = redisUrlDB.Set(db.Context, id, url.Url, url.Expiry*3600*time.Second).Err()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}

	// apiQuota, _ := strconv.Atoi(os.Getenv("API_QUOTA"))

	remainingQuota, _ := redisUrlDB.Decr(db.Context, ctx.ClientIP()).Result()

	ctx.JSON(http.StatusOK, gin.H{
		"url": url.Url,
		"custom_short": os.Getenv("DOMAIN") + "/" + string(id),
		"expiry": url.Expiry,
		"x-rate-limiting": int(remainingQuota),
		"x-rate-limiting-reset": int(limit / time.Nanosecond / time.Minute),
	})
}
