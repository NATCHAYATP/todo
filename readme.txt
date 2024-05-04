How to use gin

step 1
create file 
-> mkdir todo && cd todo

setup project
-> git init
-> go mod init main.go

import gin
-> import "github.com/gin-gonic/gin"
-> go get -u github.com/gin-gonic/gin

create func main 

use gin follow pattern
-> gin.Default() 
-> r.GET("/", func(c *gin.Context){
    c.JSON(200, gin.H{
        "message":"test",
    })
})
-> r.Run()

--------------------------------
note 
- I delete "" in local.env for make 21, old code ->
DB_CONN="root:my-secret-pw@tcp(127.0.0.1:3306)/myapp?charset-utf8mb4&parseTime=True&loc=Local"

How to do TodoList 

1. create struct in todo.go
-> gorm.model

2. create gorm type for use db
3. create func handler for frontend send
4. create func newTask for get data form path
-> implement handler
-> use gin.Context 
-> var for get data
-> blindJson
-> if err, use c.JSON for set error
    -> should use http and err.Error()

5. database
in todo.go
-> use db.create

in main.go
-> open
-> if err use panic
-> migrate
-> path and method
-> handler
-> use func newTask

6. return in todo.go
-> Create
-> error json
-> statusOK, id

---finish---

7. jwt
-> create directory auth and file auth.go
-> create func AccessToken(put context form gin)
-> jwt.NewWithClaims like this for set ExpiresAt (token time)
token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(5 * time.Minute).Unix(),
})
-> put license use SignedString(if use string must have []byte(string))
token.SignedString()
-> if have error send statusInternalServerError with JSON message
-> send token to frontend with JSON message

8. call jwt in main.go for give token when call path
-> you can add under r.GET("/)
-> add path for call AccessToken

9. todo.go add check(time and same license) token 
-> create file auth/protect.go for check token
-> create func Protect() paramiter is tokenString return error
-> use jwt.Parse paramiter is tokenString and func (Parse help to validate token)
-> func for check it same jwt.SigningMethodHMSC if not same return nil, error but same return license
func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte("==signature=="), nil
})

10. add check token in todo.go func newTask
-> call func Protect but Protect need have token 
-> token can call form context header 
s := c.Request.Header.Get("Authorization")
-> if user don't know token or Bearer token you can use trimPrefix for cut Bearer 
-> fix func Protect when have error must call c.AbortWithStatus(http.StatusUnauthorized) use Abort help to stop, not call any middleware

---- protect before have middleware
package auth

import (
	"fmt"

	jwt "github.com/dgrijalva/jwt-go"
)

func Protect(tokenString string) error {
	// use jwt.Parse
	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// check it same jwt.SigningMethodHMSC 
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte("==signature=="), nil
	})

	return err
}
------

11. use middleware
-> auth/protect adjust code, func Protect paramiter is signature (signature sented from main) 
-> func Protect return gin.HandlerFunc
-> move code in {}
-> edit return err to c.Next() for call other middleware or handler
-> move Request.Header and TrimPrefix in todo.go to auth/protect func Protect
s := c.Request.Header.Get("Authorization")
tokenString := strings.TrimPrefix(s, "Bearer ")
-> edit return of jwt.Parse to return signature, nil
-> check err of jwt.Parse if err, c.AbortWithStatus

12. fix auth.go func AccessToken use higher-order function(have func is argument or return func)
-> take signature is paramiter and return gin.HandlerFunc
-> use other code in {}
-> add return func(c *gin.Context) in all code
-> change ==signature== to signature
-> delete Request.Header, TrimPrefix and check error in todo.go

13. fix main.go use middleware
-> call auth.AccessToken() 
r.GET("/tokenz", auth.AccessToken("==signature==")) 
-> add r.Group
protected := r.Group("", auth.Protect([]byte("==signature==")))
-> use protected insteal r.POST("/todos", handler.NewTask)
protected.POST("/todos", handler.NewTask)

---- AccessToken before have middleware
func AccessToken(c *gin.Context) {
	// jwt.NewWithClaims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(5 * time.Minute).Unix(),
	})
	// put license use SignedString
	ss, err := token.SignedString([]byte("==signature=="))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// send token to frontend
	c.JSON(http.StatusOK, gin.H{
		"token": ss,
	})
}
------

---- NewTask before have middleware
func (t *TodoHandler) NewTask(c *gin.Context){
	// token can call form context header
	s := c.Request.Header.Get("Authorization")
	// use trimPrefix for cut Bearer
	tokenString := strings.TrimPrefix(s, "Bearer ")

	// call func Protect
	if err := auth.Protect(tokenString); err!= nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	var todo Todo
	if err := c.ShouldBindJSON(&todo); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	// database
	r := t.db.Create(&todo)
	// return
	if err := r.Error; err != nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ID": todo.Model.ID,
	})
}
------

14. add config with local.env because i want to take signature with config
-> godotenv.Load("local.env")
-> create .gitignore and add local.env and db.test in file
touch .gitignore
-> change ==signature== to os.Getenv(key string)
r.GET("/tokenz", auth.AccessToken(os.Getenv("SIGN")))

15. create file local.env for set value when lazy alway input data ex port=8081 go run main.go
-> touch local.env
-> set value in local.env
PORT=8080
SIGN=mysignature

======= main.go before graceful shutdown ======
package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"main.go/auth"
	"main.go/todo"
)

func main() {
	err := godotenv.Load("local.env")
	if err != nil {
		log.Printf("please consider environment variables: %s", err)
	}

	// open
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// migrate
	db.AutoMigrate(&todo.Todo{})

	// path and method
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "test",
		})
	})

	r.GET("/tokenz", auth.AccessToken(os.Getenv("SIGN")))
	protected := r.Group("", auth.Protect([]byte(os.Getenv("SIGN"))))

	// handler
	handler := todo.NewHandler(db)
	// use func newTask
	protected.POST("/todos", handler.NewTask)
	r.Run()
}

============
16. make graceful shutdown
-> create instant server
s := &http.Server{
		Addr: ":" + os.Getenv("PORT"),
		Handler: r,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second,
		MaxHeaderBytes: 1 << 20,
}
-> use s.ListenAndServe for stop. call in go rutein because wait to handler signal, don't want stop here
go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
}()
-> handle signal
ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
defer stop()
-> ctx.Done when have signal
<-ctx.Done()
stop()
fmt.Println("shutting down gracefully, press Ctrl+C again to force")

timeoutCtx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
defer cancel()

if err := s.Shutdown(timeoutCtx); err != nil {
	fmt.Println(err)
}
-> delete r.Run
-> test it can work run and crtl + C you can look ^Cshutting down gracefully, press Ctrl+C again to force

17. don't do but note for use in future
idflag for create path output information
-> main.go upper func main create var
var {
    buildcommit = "dev"
    buildtime = time.Now().String()
}
-> create handler /x
r.GET("/x", func(c *gin.Context) {
    c.JSON(200, gin.H{
        "buildcommit": buildcommit,
        "buildtime": buildtime,
    })
})
-> when call path you will give dev and time but you run build buildcommit is commit you will have commit number. you can read more in goodnote

18. don't do but note for use in future
liveness Probe Readiness Probe
-> [liveness Probe] add in func main. create file and delete file
_, err := os.Create("/tmp/live")
if err != nil {
    log.Fatal(err)
}
defer os.Remove("/tmp/live")
-> [Readiness Probe] check set data
-> in code func main we have gorm.Open for set database (it is Readiness Probe)
-> you can create path for test success
r.GET("/healthz", func(c *gin.Context) {
    c.Status(200)
})
-> you can check on kub but you can easy check it on service
cat /tmp/live
echo $? (can see = 0, not see = 1)

19. don't do but note for use in future
rate limit
-> mani.go func main downest
var limiter = rate.NewLimiter(5, 5)

func limiteHanddler(c *gin.Context) {
    if !limiter.Allow() {
        c.AbortWithStatus(http.StatusTooManyRequests)
        return
    }
    c.JSON(200, gin.H{
        "message": "pong",
    })
}
-> create path
r.GET("/limitz", limitedHandler)
-> install vegeta
go install github.com/tsenart/vegeta@latest
-> use vegeta for call path
echo "GET http://:8081/limitz" | vegeta attack -rate=10/s -duration=1s | vegeta report

20. change DB and move to config file
-> open local.env and add DB_CONN="ex.db"
-> main.go change from sqlite.Open("test.db") to
sqlite.Open(os.Getenv("DB_CONN"))
-> I create makeFile you can run for build db maria in docker
make maria
-> fix local.env follow makefile na
DB_CONN="root:my-secret-pw@tcp(127.0.0.1:3306)/myapp?charset-utf8mb4&parseTime=True&loc=Local"
-> main.go gorm.open... change sqlite to mysql
-> go get gorm.io/driver/mysql
-> test by postman and open mysql

21. docker
-> add Dockerfile 
-> add script in Makefile 
image:
	docker build -t todo:test -f Dockerfile .

container:
	docker run  -p:8080:8080 --evn-file ./local.env --link some-mariadb:db \
	--name myapp todo:test
-> run image
make image
-> run container but ตอนแรกเราเรียก local.env โดยตรง แต่ตอนนี้เราเอาไปใช้เรียกใน Dockerfile 
โค้ดนี้อ่ะที่เราไม่ได้เรียก local.env โดยตรง docker run  -p:8080:8080 --evn-file ./local.env
มันจะกลายไปเป็นค่าใน container ทำให้มันเรียก environment ผ่าน os อีกทีแล้วมันไม่รู้จัก "" ต้องลบ "" ใน local.env ออก

