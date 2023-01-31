package main

import (
	"net/http"
	"os"

	"github.com/Edgar-Avila/go-sqlmql/parser"
	"github.com/Edgar-Avila/go-sqlmql/translate"
	"github.com/gin-gonic/gin"
)

type TranslateBody struct {
	Text string `binding:"required" json:"text"`
}

func main() {
	router := gin.Default()
	router.SetTrustedProxies(nil)
	router.Use(func(ctx *gin.Context) {
		ctx.Writer.Header().Set("Access-Control-Allow-Origin", os.Getenv("FRONT_URL"))
		ctx.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST")
		ctx.Next()
	})
	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Api is running",
		})
	})
	router.POST("/translate", func(ctx *gin.Context) {
		var body TranslateBody
		if err := ctx.ShouldBindJSON(&body); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"ok":      false,
				"message": "Property text required",
			})
			return
		}
		p, err := parser.NewParser()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"ok":      false,
				"message": "Internal Server Error",
			})
			return
		}
		sqlFile, err := p.ParseString("", body.Text)
		if err != nil {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{
				"ok":      false,
				"message": err.Error(),
			})
			return
		}
		translated, err := translate.TranslateSqlFile(sqlFile)
		if err != nil {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{
				"ok":      false,
				"message": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"ok":         true,
			"message":    "Successfully translated",
			"translated": translated,
		})
	})
	router.Run()
}
