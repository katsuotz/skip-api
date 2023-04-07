package helper

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func GetUserLocation(ctx *gin.Context) string {
	ip := ctx.ClientIP()

	resp, err := http.Get(fmt.Sprintf("http://ip-api.com/json/%s", ip))
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	var data map[string]interface{}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return ""
	}

	defer func() {
		if err := recover(); err != nil {
			log.Println("Recovered from panic:", err)
		}
	}()

	country := data["country"].(string)
	//region := data["regionName"].(string)
	city := data["city"].(string)

	return city + ", " + country
}
