package auth

import (
	"net/http"
	"os"
	"time"

	"melivecode/go-jwt-api/orm"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

var hmacSampleSecret []byte
//Binding from JSON
type RegisterBody struct {
	User     string `form:"user"  binding:"required"`
	Password string `form:"password"  blinding:"required"`
	Fullname string `form:"fullname"  blinding:"required"`
	Avatar   string `form:"avatar"  blinding:"required"`
}
type LoginBody struct {
	User     string `form:"user"  binding:"required"`
	Password string `form:"password"  blinding:"required"`
}

func Register(c *gin.Context) {
	json := RegisterBody{}
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//check user
	userExist := orm.User{}
	orm.Db.Where("user=?", json.User).First(&userExist)
	if userExist.ID > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "user exists",
		})
		return
	}
	encryptPassword, _ := bcrypt.GenerateFromPassword([]byte(json.Password), 10)
	user := orm.User{
		User:     json.User,
		Password: string(encryptPassword),
		Fullname: json.Fullname,
		Avatar:   json.Avatar,
	}
	orm.Db.Create(&user)
	if user.ID > 0 {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "user create success",
			"userId":  user.ID,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":  "error",
			"message": "user create failed",
		})
	}
}
func Login(c *gin.Context) {
	json := LoginBody{}
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//Check User
	userExist :=orm.User{}
	orm.Db.Where("user = ?",json.User).First(&userExist)
	if userExist.ID == 0{
		c.JSON(http.StatusOK,gin.H{"status":"error","message":"User Does Not Exists"})
		return
	}
	e:=bcrypt.CompareHashAndPassword([]byte(userExist.Password),[]byte(json.Password))
	if e == nil {
		hmacSampleSecret=[]byte(os.Getenv("JWT_SECRET_KEY"))

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"userId": userExist.ID,
			"ExpiresAt":time.Now().Add(time.Minute*1).Unix(),
		})
		// Sign and get the complete encoded token as a string using the secret
		tokenString, _ := token.SignedString(hmacSampleSecret)
		c.JSON(http.StatusOK,gin.H{"status":"ok","message":"Login successed","token": tokenString})
	} else {
		c.JSON(http.StatusBadRequest,gin.H{"status":"err","message":"Login Failed"})
	}	
}
