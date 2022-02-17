package src

import (
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
)

func MustLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		if val, err := c.Cookie("userID"); err != nil {
			c.String(http.StatusUnauthorized, val)
			c.Abort()
		} else {
			fmt.Println(val)
			c.Next()
		}
	}
}

func MemberCreate(c *gin.Context) {
	creatememberrequest := CreateMemberRequest{}
	if err := c.BindJSON(&creatememberrequest); err != nil {
		c.JSON(200, CreateMemberResponse{
			Code: 1,
			Data: struct{ UserID string }{UserID: ""},
		})
	} else {
		tmember := TMember{}
		DBHelper.Where("username=?", creatememberrequest.Username).Find(&tmember)
		if tmember.UserID != "" {
			c.JSON(200, CreateMemberResponse{
				Code: UserHasExisted,
				Data: struct{ UserID string }{UserID: ""},
			})
		} else {
			tmember.Username = creatememberrequest.Username
			tmember.Nickname = creatememberrequest.Nickname
			tmember.Password = creatememberrequest.Password
			tmember.UserType = creatememberrequest.UserType
			u := make([]byte, 16)
			_, err := rand.Read(u)
			if err != nil {
				return
			}
			u[8] = (u[8] | 0x80) & 0xBF // what does this do?
			u[6] = (u[6] | 0x40) & 0x4F // what does this do?
			tmember.UserID = hex.EncodeToString(u)
			DBHelper.Create(&tmember)
			c.JSON(200, CreateMemberResponse{
				Code: OK,
				Data: struct{ UserID string }{UserID: tmember.UserID},
			})
		}
	}
}

func Member(c *gin.Context) {
	getmemberrequest := GetMemberRequest{}
	tmember := TMember{}
	if err := c.BindQuery(&getmemberrequest); err != nil {
		c.JSON(200, GetMemberResponse{
			Code: 1,
			Data: tmember,
		})
	} else {
		DBHelper.Where("userID=?", getmemberrequest.UserID).Find(&tmember)
		if tmember.UserID == "" {
			c.JSON(200, GetMemberResponse{
				Code: UserNotExisted,
				Data: tmember,
			})
		} else {
			c.JSON(200, GetMemberResponse{
				Code: OK,
				Data: tmember,
			})
		} //此处少一个用户已被删除的处理
	}
}

func MemberList(c *gin.Context) {
	getmemberlistrequest := GetMemberListRequest{}
	var tmemberlist []TMember
	if err := c.BindQuery(getmemberlistrequest); err != nil {
		c.JSON(200, GetMemberListResponse{
			Code: 1,
			Data: struct{ MemberList []TMember }{MemberList: tmemberlist},
		})
	} else {
		DBHelper.Limit(getmemberlistrequest.Limit).Offset(getmemberlistrequest.Offset).Find(&tmemberlist)
		c.JSON(200, GetMemberListResponse{
			Code: OK,
			Data: struct{ MemberList []TMember }{MemberList: tmemberlist},
		})
	}
}

func MemberUpdate(c *gin.Context) {
	updatememberrequest := UpdateMemberRequest{}
	tmember := TMember{}
	if err := c.BindJSON(updatememberrequest); err != nil {
		c.JSON(200, UpdateMemberResponse{
			Code: 1,
		})
	} else {
		DBHelper.Where("userID=?", updatememberrequest.UserID).Find(&tmember)
		if tmember.UserID == "" {
			c.JSON(200, UpdateMemberResponse{
				Code: UserNotExisted,
			})
		} else {
			DBHelper.Where("userID=?", updatememberrequest.UserID).Update("nickname", updatememberrequest.Nickname)
			c.JSON(200, UpdateMemberResponse{
				Code: 1,
			})
		}
	}
}

func MemberDelete(c *gin.Context) {
	deletememberrequest := DeleteMemberRequest{}
	tmember := TMember{}
	if err := c.BindJSON(&deletememberrequest); err != nil {
		c.JSON(200, DeleteMemberResponse{
			Code: 1,
		})
	} else {
		DBHelper.Where("userID=?", deletememberrequest.UserID).Update("IsExisted", 1).Find(&tmember)
		c.JSON(200, DeleteMemberResponse{
			Code: OK,
		})
	}
}

func AuthLogin(c *gin.Context) {
	loginrequest := LoginRequest{}
	if err := c.BindJSON(&loginrequest); err != nil {
		c.JSON(200, LoginResponse{1, struct{ UserID string }{UserID: ""}})
	} else {
		tmember := TMember{}
		DBHelper.Where("username=?", loginrequest.Username).Where("password=?", loginrequest.Password).Find(&tmember)
		if tmember.UserID == "" {
			c.JSON(200, LoginResponse{5, struct{ UserID string }{UserID: ""}})
		} else {
			c.SetCookie("userID", tmember.UserID, 3600, "/", "localhost", false, true)
			c.JSON(200, LoginResponse{0, struct{ UserID string }{UserID: tmember.UserID}})
		}
	}
}

func AuthLogout(c *gin.Context) {
	c.SetCookie("userID", "", -1, "/", "localhost", false, true)
	c.JSON(200, LogoutResponse{Code: 6})
}

func Authwhoami(c *gin.Context) {
	if val, err := c.Cookie("userID"); err != nil {
		c.JSON(200, WhoAmIResponse{
			Code: 6,
			Data: TMember{},
		})
	} else {
		tmember := TMember{}
		DBHelper.Where("userID=?", val).Find(&tmember)
		c.JSON(200, WhoAmIResponse{
			Code: 0,
			Data: tmember,
		})
	}
}

func CreateCourse(c *gin.Context) {
	createcourserequest := CreateCourseRequest{}
	tcourse := TCourse{}
	if err := c.BindJSON(&createcourserequest); err != nil {
		c.JSON(200, CreateMemberResponse{
			Code: 1,
			Data: struct{ UserID string }{UserID: tcourse.CourseID},
		})
	} else {
		DBHelper.Where("name=?", createcourserequest.Name).Find(&tcourse)
		tcourse.Name = createcourserequest.Name
		tcourse.Capacity = createcourserequest.Cap
		DBHelper.Create(&tcourse)
	}
}

func GetCourse(c *gin.Context) {
	getcourserequest := GetCourseRequest{}
	tcourse := TCourse{}
	if err := c.BindQuery(&getcourserequest); err != nil {
		c.JSON(200, GetCourseResponse{
			Code: 1,
			Data: tcourse,
		})
	} else {
		DBHelper.Where("course_id=?", getcourserequest.CourseID).Find(&tcourse)
		c.JSON(200, GetCourseResponse{
			Code: OK,
			Data: tcourse,
		})
	}
}

func TeacherBindCourse(c *gin.Context) {
	bindcourserequest := BindCourseRequest{}
	tcourse := TCourse{}
	if err := c.BindJSON(&bindcourserequest); err != nil {
		c.JSON(200, BindCourseResponse{
			Code: 1,
		})
	} else {
		DBHelper.Where("course_id=?", bindcourserequest.CourseID).Find(&tcourse)
		tcourse.TeacherID = bindcourserequest.TeacherID
		c.JSON(200, BindCourseResponse{
			Code: OK,
		})
	}
}

func TeacherUnbindCourse(c *gin.Context) {
	unbindcourserequest := UnbindCourseRequest{}
	tcourse := TCourse{}
	if err := c.BindJSON(&unbindcourserequest); err != nil {
		c.JSON(200, UnbindCourseResponse{
			Code: 1,
		})
	} else {
		DBHelper.Where("course_id=?", unbindcourserequest.CourseID).Find(&tcourse)
		tcourse.TeacherID = ""
		c.JSON(200, BindCourseResponse{
			Code: OK,
		})
	}
}

func GetTopicDetail(c *gin.Context) {
	tid := c.Param("topic_id")
	topics := Topic{}
	DBHelper.Find(&topics, tid) //从数据库取
	c.Set("dbResult", topics)
}

func NewTopic(c *gin.Context) {
	//判断登录
	topic := Topic{}
	if err := c.BindJSON(&topic); err != nil {
		c.String(400, "参数错误:%s", err.Error())
	} else {
		c.JSON(200, topic)
	}
}

func DeleteTopic(c *gin.Context) {
	//判断登录
	c.String(200, "删除帖子")
}

func GetTopicList(c *gin.Context) {
	query := TopicQuery{}
	if err := c.BindQuery(&query); err != nil {
		c.String(400, "参数错误:%s", err.Error())
	} else {
		c.JSON(200, query)
	}
}
