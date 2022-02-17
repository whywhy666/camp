package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	. "topic/src"
)

func main() {

	router := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("topicurl", TopicUrl)
	} //自定义判定规则

	// 成员管理
	v1 := router.Group("7/member")
	{
		v1.Use(MustLogin())
		{
			v1.POST("/create", MemberCreate)
			v1.GET("", Member)
			v1.GET("/list", MemberList)
			v1.POST("/update", MemberUpdate)
			v1.POST("/delete", MemberDelete)
		}
	}

	// 登录
	v2 := router.Group("7/auth")
	{
		v2.POST("login", AuthLogin)
		v2.POST("logout", AuthLogout)
		v2.GET("whoami", Authwhoami)
	}

	// 排课
	v3 := router.Group("7/course")
	{
		v3.Use(MustLogin())
		{
			v3.POST("/create", CreateCourse)
			v3.GET("/get", GetCourse)
			v3.POST("/schedule")
		}
	}

	v4 := router.Group("7/teacher")
	{
		v4.POST("/bind_course", TeacherBindCourse)
		v4.POST("/unbind_course", TeacherUnbindCourse)
		v4.GET("/get_course")
	}

	// 抢课
	v5 := router.Group("7/student")
	{
		v5.POST("/book_course")
		v5.GET("/course")
	}

	v6 := router.Group("v1/topics")
	{
		v6.GET("", GetTopicList)

		v6.GET("/:topic_id", CacheDecorator(GetTopicDetail, "topic_id", "topic_%s", Topic{}))

		v6.Use(MustLogin())
		{
			v6.POST("", NewTopic)
			v6.DELETE("/:topic_id", DeleteTopic)
		}
	}

	router.Run()

}
