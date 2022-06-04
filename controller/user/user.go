package user

import (
	
	"melivecode/go-jwt-api/orm"
	"net/http"
	

	"github.com/gin-gonic/gin"
)

func ReadAll(c *gin.Context){
	/*
	header:=c.Request.Header.Get("Authorization")
	tokenString:=strings.Replace(header, "Bearer ","",1)
	hmacSampleSecret :=[]byte(os.Getenv("JWT_SECRET_KEY"))
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return hmacSampleSecret, nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims["userId"])
		var users []orm.User
		orm.Db.Find(&users)
		c.JSON(http.StatusOK,gin.H{"status":"ok","message":"user read success","users":users})
	} else {
		// fmt.Println(err)
		c.JSON(http.StatusUnauthorized,gin.H{"status":"error","message":err.Error()})
		return 
	}
	*/
	var users []orm.User
	orm.Db.Find(&users)
	c.JSON(http.StatusOK,gin.H{"status":"ok","message":"user read success","users":users})

}