package src

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"github.com/pquerna/ffjson/ffjson"
	"log"
)

//缓存装饰器
func CacheDecorator(h gin.HandlerFunc, param string, redKeyPattern string, empty interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		//redis判断
		getID := c.Param(param)
		redisKey := fmt.Sprint(redKeyPattern, getID)
		conn := RedisDefaultPool.Get()
		defer conn.Close()
		ret, err := redis.Bytes(conn.Do("get", redisKey))
		if err != nil {
			h(c) //执行目标方法
			dbResult, exists := c.Get("dbResult")
			if !exists {
				dbResult = empty
			}
			retData, _ := ffjson.Marshal(dbResult)
			conn.Do("setex", redisKey, 20, retData)
			c.JSON(200, dbResult)
			log.Println("从数据库读取")
		} else { //缓存有 直接抛出
			log.Println("从Redis读取")
			ffjson.Unmarshal(ret, &empty)
			c.JSON(200, empty)
		}
	}
}
