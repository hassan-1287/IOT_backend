package handlers

import (
	"context"
	"encoding/json"
	"ginframework/intarnel/auth"
	gconfig "ginframework/intarnel/auth/config"
	"ginframework/intarnel/auth/middleware"
	"ginframework/intarnel/auth/repository"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

/*
*  داله تسجيل الدخول ثم رابط التوجيه الى كوكل
 */
var stateStore = map[string]bool{}

func Login(ctx *gin.Context) {
	state, err := auth.GenerateState(64)
	if err != nil {
		return
	}
	stateStore[state] = true

	url := gconfig.GoogleConfig.AuthCodeURL(state)

	ctx.Redirect(http.StatusTemporaryRedirect, url)

}

// *الداله تعمل بعد العوده من كوكل وتمرير البيانات الى السيرفر الخاص بي
func CallBackFromGoogle(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx context.Context = context.Background()

		stateFromURL := c.Query("state")
		codeFromURL := c.Query("code")

		if stateFromURL == "" || codeFromURL == "" {
			c.JSON(400, gin.H{"error": "missing values"})
			return
		}

		if !stateStore[stateFromURL] {
			c.JSON(401, gin.H{"error": "invalid state"})
			return
		}
		
		delete(stateStore, stateFromURL)

		if stateFromURL == "" || codeFromURL == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "meassing value from url"})
			return
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

		user, err := repository.CreateNewUser(pool, userInfo)
		jwtToken, err := middleware.GenerateToken(user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to generate token",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "تم تسجيل الدخول بنجاح",
			"token":  jwtToken,
			"user":   user,
		})
	}
}
