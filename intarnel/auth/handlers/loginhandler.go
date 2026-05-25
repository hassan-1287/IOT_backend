package handlers

import (
	"context"
	"encoding/json"
	"ginframework/intarnel/auth"
	gconfig "ginframework/intarnel/auth/config"
	"ginframework/intarnel/auth/middleware"
	"ginframework/intarnel/auth/repository"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
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
func CallBackFromGoogle(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		//? يعتبر هذا التعريف افضل لانه مرتبط بالطلب بشكل مباشر
		ctx, cancel := context.WithTimeout(
			c.Request.Context(),
			5*time.Second,
		)
		defer cancel()

		stateFromURL := c.Query("state")
		codeFromURL := c.Query("code")

		if stateFromURL == "" || codeFromURL == "" {
			c.JSON(400, gin.H{"error": "missing values"})
			return
		}

		if stateFromURL == "" || codeFromURL == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "meassing value from url"})
			return
		}

		code := c.Query("code")
		token, err := gconfig.GoogleConfig.Exchange(ctx, code)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid authorization code"})
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

		// 1. محاولة إنشاء المستخدم أو جلب بياناته
		user, err := repository.CreateNewUser(ctx, pool, userInfo)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "error",
				"msg":    "فشل حفظ المستخدم في قاعدة البيانات",
				"error":  err.Error(),
			})
			return // إيقاف التنفيذ فوراً في حال حدوث خطأ
		}

		// 2. فحص إضافي للأمان للتأكد أن المؤشر ليس nil
		if user == nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "error",
				"msg":    "لم يتم العثور على بيانات المستخدم",
			})
			return
		}

		// 3. الآن يمكنك توليد التوكن بأمان تالم
		jwtToken, err := middleware.GenerateToken(user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "error",
				"msg":    "failed to generate token",
				"error":  err.Error(),
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
