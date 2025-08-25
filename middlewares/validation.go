package middlewares

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ValidationHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		fmt.Println("Validation handler")
		// get option param from requests and check if it is a valid year & after 2018
		year := ctx.Param("year")
		fmt.Println(year)
		if year != "" {
			yearInt, err := strconv.Atoi(year)
			fmt.Println(year)
			if err != nil || yearInt < 2018 {
				fmt.Println("Year must be a valid integer after 2018")
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "Year must be a valid integer after 2018"})
				ctx.Abort()
				return
			}
			//set the year param to the context
			ctx.Set("year", yearInt)
		}

		ctx.Next()
	}
}
