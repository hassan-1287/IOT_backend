package handlers

import (
	"context"
	"encoding/json"
	"ginframework/intarnel/auth"
	gconfig "ginframework/intarnel/auth/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
*  داله تسجيل الدخول ثم رابط التوجيه الى كوكل
 */
func Login(ctx *gin.Context) {
	state, err := auth.GenerateState(64)
	if err != nil {
		return
	}
	url := gconfig.GoogleConfig.AuthCodeURL(state)
	ctx.Redirect(http.StatusTemporaryRedirect, url)

}

// *الداله تعمل بعد العوده من كوكل وتمرير البيانات الى السيرفر الخاص بي
func CallBackFromGoogle(c *gin.Context) {
	var ctx context.Context = context.Background()
	stateFromURL := c.Query("state")
	codeFromURL := c.Query("code")

	if stateFromURL == "" || codeFromURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "meassing value from url"})
	}

	code := c.Query("code")
	token, err := gconfig.GoogleConfig.Exchange(ctx, code)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	client := gconfig.GoogleConfig.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return

	}
	defer resp.Body.Close()

	var userInfo auth.GoogleUserInfo
	
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error",
			"msg":   "no data err!!!!",
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "تم تسجيل الدخول بنجاح!",
		"user":    userInfo,
	})
}
