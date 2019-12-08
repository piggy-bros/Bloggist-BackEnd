package api

import (
	"bloggist/static"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/boltdb/bolt"
	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	var receive static.UserInfo
	c.BindJSON(&receive)
	//get dataBase
	db, err := bolt.Open("server.db", 0600, nil)
	if err != nil {
		fmt.Println("Open DB error!")
	}
	defer db.Close()
	//get table user
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("UserInfo"))
		if b == nil {
			_, tableErr := tx.CreateBucket([]byte("UserInfo"))
			if tableErr != nil {
				fmt.Println("Open table error!")
			}
		}

		//operation on the table
		b = tx.Bucket([]byte("UserInfo"))
		data := b.Get([]byte(receive.Username))
		if data == nil { //No such User
			c.JSON(http.StatusOK, gin.H{
				"status":  "fail",
				"err_msg": "no such user",
			})
		} else { //User exists
			if string(data) == receive.Password { //login succeed
				//get table blog
				b = tx.Bucket([]byte("BlogInfo"))
				if b == nil {
					_, tableErr := tx.CreateBucket([]byte("BlogInfo"))
					if tableErr != nil {
						fmt.Println("Open table error!")
					}
				}
				b = tx.Bucket([]byte("BlogInfo"))

				//scan every blog
				db.View(func(tx *bolt.Tx) error {
					temp := b.Cursor()
					blogNum := 0
					likedNum := 0
					for k, v := temp.First(); k != nil; k, v = temp.Next() {
						fmt.Printf("key=%s, value=%s\n", k, v)
						var tempBlog static.BlogInfo
						if castErr := json.Unmarshal(v, &tempBlog); castErr == nil {
							if tempBlog.Author == receive.Username {
								blogNum = blogNum + 1
								likedNum = likedNum + tempBlog.LikedNum
							}
						}
					}
					// add token
					tokenString, tokenError := GenerateToken(receive.Username, receive.Password)
					if tokenError != nil {
						fmt.Println(tokenError)
					} else {
						c.JSON(http.StatusOK, gin.H{
							"status":       "success",
							"blog_num":     blogNum,
							"liked_num":    likedNum,
							"token_string": tokenString,
						})
					}
					return nil
				})
			} else { //login fail
				c.JSON(http.StatusOK, gin.H{
					"status":  "fail",
					"err_msg": "wrong username or password",
				})
			}
		}
		return nil
	})
	if err != nil {
		fmt.Println("DB update error!")
	}
}

func Register(c *gin.Context) {
	var receive static.UserInfo
	c.BindJSON(&receive)
	//get dataBase
	db, err := bolt.Open("server.db", 0600, nil)
	if err != nil {
		fmt.Println("Open DB error!")
	}
	defer db.Close()
	//get table user
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("UserInfo"))
		if b == nil {
			_, tableErr := tx.CreateBucket([]byte("UserInfo"))
			if tableErr != nil {
				fmt.Println("Open table error!")
			}
		}

		//operation on the table
		b = tx.Bucket([]byte("UserInfo"))
		data := b.Get([]byte(receive.Username))
		if data == nil { //Succeed
			err := b.Put([]byte(receive.Username), []byte(receive.Password))
			if err != nil {
				fmt.Println("insert data error!")
			} else {
				c.JSON(http.StatusOK, gin.H{
					"status": "success",
				})
			}
		} else { //Duplicate
			c.JSON(http.StatusOK, gin.H{
				"status":  "fail",
				"err_msg": "duplicate",
			})
		}
		return nil
	})
	if err != nil {
		fmt.Println("DB update error!")
	}
}
