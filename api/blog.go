package api

import (
	"bloggist/static"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/boltdb/bolt"
	"github.com/gin-gonic/gin"
)

func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

func GetBlog(c *gin.Context) {
	tokenString := c.Query("token")
	claim, tokenErr := ParseToken(tokenString)
	if tokenErr != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  "fail",
			"err_msg": tokenErr,
		})
	} else if !(claim.Username == c.Param("name")) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "fail",
			"err_msg": "No match username",
		})
	} else {
		name := c.Param("name")
		blogid := c.Param("blogid")

		//get dataBase
		db, err := bolt.Open("server.db", 0600, nil)
		if err != nil {
			fmt.Println("Open DB error!")
		}
		defer db.Close()
		//get table blog
		err = db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("BlogInfo"))
			if b == nil {
				_, tableErr := tx.CreateBucket([]byte("BlogInfo"))
				if tableErr != nil {
					fmt.Println("Open table error!")
				}
			}
			//operation on the table
			b = tx.Bucket([]byte("BlogInfo"))
			blogIDInt, _ := strconv.Atoi(blogid)
			blogIDByte := itob(blogIDInt)
			data := b.Get(blogIDByte)
			if data == nil {
				c.JSON(http.StatusOK, gin.H{
					"status":  "fail",
					"err_msg": "No such blog",
				})
			} else {
				var tempBlog static.BlogInfo
				if castErr := json.Unmarshal(data, &tempBlog); castErr == nil {
					if tempBlog.Author == name {
						v := b.Get([]byte(itob(tempBlog.BlogID)))
						var tempBlog static.BlogInfo
						if castErr := json.Unmarshal(v, &tempBlog); castErr == nil {
							blogTitle := tempBlog.BlogTitle
							blogContent := tempBlog.BlogContent
							Liked := tempBlog.LikedNum
							c.JSON(http.StatusOK, gin.H{
								"title":   blogTitle,
								"content": blogContent,
								"liked":   Liked,
							})
						} else {
							fmt.Println("cast Error!")
						}
					} else {
						c.JSON(http.StatusOK, gin.H{
							"status":  "fail",
							"err_msg": "Unmatch author and blog",
						})
					}
				}
			}
			return nil
		})
	}
}

func GetBlogs(c *gin.Context) {
	tokenString := c.Query("token")
	claim, tokenErr := ParseToken(tokenString)
	if tokenErr != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  "fail",
			"err_msg": "token wrong",
		})
	} else if !(claim.Username == c.Param("name")) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "fail",
			"err_msg": "No match username",
		})
	} else {
		name := c.Param("name")
		//get dataBase
		db, err := bolt.Open("server.db", 0600, nil)
		if err != nil {
			fmt.Println("Open DB error!")
		}
		defer db.Close()
		err = db.Update(func(tx *bolt.Tx) error {
			//get table blog
			fmt.Println(tx)
			b := tx.Bucket([]byte("BlogInfo"))
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
				var blogsTitle []string
				var blogsID []int
				var Liked []int
				for k, v := temp.First(); k != nil; k, v = temp.Next() {
					fmt.Printf("key=%s, value=%s\n", k, v)
					var tempBlog static.BlogInfo
					if castErr := json.Unmarshal(v, &tempBlog); castErr == nil {
						if tempBlog.Author == name {
							blogsTitle = append(blogsTitle, tempBlog.BlogTitle)
							blogsID = append(blogsID, tempBlog.BlogID)
							Liked = append(Liked, tempBlog.LikedNum)
						}
					}
				}
				c.JSON(http.StatusOK, gin.H{
					"status":   "success",
					"blogs":    blogsTitle,
					"blog_ids": blogsID,
					"liked":    Liked,
				})
				return nil
			})

			return nil
		})
	}
}

func DeleteBlog(c *gin.Context) {
	tokenString := c.Query("token")
	claim, tokenErr := ParseToken(tokenString)
	if tokenErr != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  "fail",
			"err_msg": tokenErr,
		})
	} else if !(claim.Username == c.Param("name")) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "fail",
			"err_msg": "No match username",
		})
	} else {
		name := c.Param("name")
		blogid := c.Param("blogid")

		//get dataBase
		db, err := bolt.Open("server.db", 0600, nil)
		if err != nil {
			fmt.Println("Open DB error!")
		}
		defer db.Close()
		//get table blog
		err = db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("BlogInfo"))
			if b == nil {
				_, tableErr := tx.CreateBucket([]byte("BlogInfo"))
				if tableErr != nil {
					fmt.Println("Open table error!")
				}
			}
			//operation on the table
			b = tx.Bucket([]byte("BlogInfo"))
			data := b.Get([]byte(blogid))
			if data == nil {
				c.JSON(http.StatusOK, gin.H{
					"status":  "fail",
					"err_msg": "No such a blog",
				})
			} else {
				var tempBlog static.BlogInfo
				if castErr := json.Unmarshal(data, &tempBlog); castErr == nil {
					if tempBlog.Author == name {
						b.Delete(itob(tempBlog.BlogID))
						currentNum := 0
						db.View(func(tx *bolt.Tx) error {
							temp := b.Cursor()
							for k, _ := temp.First(); k != nil; k, _ = temp.Next() {
								currentNum = currentNum + 1
							}
							return nil
						})
						c.JSON(http.StatusOK, gin.H{
							"status":   "ok",
							"blog_num": currentNum,
						})
					} else {
						c.JSON(http.StatusOK, gin.H{
							"status":  "fail",
							"err_msg": "Unmatch author and blog",
						})
					}
				}
			}
			return nil
		})
	}
}

func LikeBlog(c *gin.Context) {
	tokenString := c.Query("token")
	claim, tokenErr := ParseToken(tokenString)
	if tokenErr != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  "fail",
			"err_msg": tokenErr,
		})
	} else if !(claim.Username == c.Param("name")) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "fail",
			"err_msg": "No match username",
		})
	} else {
		name := c.Param("name")
		blogid := c.Param("blogid")

		//get dataBase
		db, err := bolt.Open("server.db", 0600, nil)
		if err != nil {
			fmt.Println("Open DB error!")
		}
		defer db.Close()
		//get table blog
		err = db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("BlogInfo"))
			if b == nil {
				_, tableErr := tx.CreateBucket([]byte("BlogInfo"))
				if tableErr != nil {
					fmt.Println("Open table error!")
				}
			}
			//operation on the table
			b = tx.Bucket([]byte("BlogInfo"))
			blogIDInt, _ := strconv.Atoi(blogid)
			blogIDByte := itob(blogIDInt)
			data := b.Get(blogIDByte)

			var tempBlog static.BlogInfo
			if data == nil {
				c.JSON(http.StatusOK, gin.H{
					"status":  "fail",
					"err_msg": "No such blog",
				})
			} else {
				if castErr := json.Unmarshal(data, &tempBlog); castErr == nil {
					if tempBlog.Author == name {
						b.Delete([]byte(itob(tempBlog.BlogID)))
						tempBlog.LikedNum = tempBlog.LikedNum + 1
						v, castErr := json.Marshal(tempBlog)
						if castErr != nil {
							fmt.Println("Cast Error!")
						}
						putErr := b.Put(itob(tempBlog.BlogID), v)
						if putErr != nil {
							fmt.Println("Put Error!")
						}
						c.JSON(http.StatusOK, gin.H{
							"status": "ok",
							"liked":  tempBlog.LikedNum,
						})
					} else {
						c.JSON(http.StatusOK, gin.H{
							"status":  "fail",
							"err_msg": "Unmatch author and blog",
						})
					}
				} else {
					fmt.Println(castErr)
				}
			}
			return nil
		})
		if err != nil {
			fmt.Println(err)
		}
	}
}

func PublishBlog(c *gin.Context) {
	tokenString := c.Query("token")
	claim, tokenErr := ParseToken(tokenString)
	if tokenErr != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  "fail",
			"err_msg": tokenErr,
		})
	} else if !(claim.Username == c.Param("name")) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "fail",
			"err_msg": "No match username",
		})
	} else {
		name := c.Param("name")
		var receive static.PublishInfo
		c.BindJSON(&receive)
		//get dataBase
		db, err := bolt.Open("server.db", 0600, nil)
		if err != nil {
			fmt.Println("Open DB error!")
		}
		defer db.Close()
		//get table blog
		err = db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("BlogInfo"))
			if b == nil {
				_, tableErr := tx.CreateBucket([]byte("UserInfo"))
				if tableErr != nil {
					fmt.Println("Open table error!")
				}
			}

			//operation on the table
			b = tx.Bucket([]byte("BlogInfo"))
			var tempInput static.BlogInfo
			id, _ := b.NextSequence()
			tempInput.BlogID = int(id)
			tempInput.Author = name
			tempInput.LikedNum = 0
			tempInput.BlogTitle = receive.Title
			tempInput.BlogContent = receive.Content
			v, castErr := json.Marshal(tempInput)
			if castErr != nil {
				fmt.Println("cast Error!")
			} else {
				putErr := b.Put(itob(int(id)), v)
				if putErr != nil {
					fmt.Println("put Error!")
				} else {
					c.JSON(http.StatusOK, gin.H{
						"status":  "ok",
						"blog_id": int(id),
					})
				}
			}
			return nil
		})
		if err != nil {
			fmt.Println("DB update error!")
		}
	}
}
