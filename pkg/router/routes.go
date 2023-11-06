package router

import (
	"github.com/choigonyok/techlog/pkg/handler"
	"github.com/gin-gonic/gin"
)

type Route struct {
	Path    string
	Method  string
	Handler gin.HandlerFunc
}

const (
	POST   = "post"
	GET    = "get"
	PUT    = "put"
	DELETE = "delete"
)

func (r *Router) NewRoutes(prefix string) []Route {
	h := []Route{
		// Post
		{
			Path:    prefix + "post",
			Method:  POST,
			Handler: handler.CreatePost,
		},
		{
			Path:    prefix + "posts",
			Method:  GET,
			Handler: handler.GetPosts,
		},
		{
			Path:    prefix + "/posts/:postid/thumbnail",
			Method:  GET,
			Handler: handler.GetThumbnailByPostID,
		},
		{
			Path:    prefix + "posts/:postid",
			Method:  GET,
			Handler: handler.GetPost,
		},
		{
			Path:    prefix + "posts/:postid",
			Method:  PUT,
			Handler: handler.UpdatePostByPostID,
		},
		{
			Path:    prefix + "posts/:postid",
			Method:  DELETE,
			Handler: handler.DeletePostByPostID,
		},

		// Visitor
		{
			Path:    prefix + "visitor",
			Method:  GET,
			Handler: handler.GetVisitorCounts,
		},

		// Login
		{
			Path:    prefix + "login",
			Method:  POST,
			Handler: handler.VerifyAdminIDAndPW,
		},
		{
			Path:    prefix + "login",
			Method:  GET,
			Handler: handler.VerifyAdminUser,
		},

		// Tagee
		{
			Path:    prefix + "tags",
			Method:  GET,
			Handler: handler.GetTags,
		},

		// Comment
		{
			Path:    prefix + "posts/:postid/comment",
			Method:  POST,
			Handler: handler.CreateComment,
		},
		{
			Path:    prefix + "comments",
			Method:  GET,
			Handler: handler.GetComments,
		},
		{
			Path:    prefix + "posts/:postid/comments",
			Method:  GET,
			Handler: handler.GetCommentsByPostID,
		},
		{
			Path:    prefix + "posts/:postid/comments/:commentid",
			Method:  DELETE,
			Handler: handler.DeleteCommentByCommentID,
		},

		// Reply
		{
			Path:    prefix + "posts/:postid/replies",
			Method:  GET,
			Handler: handler.GetRepliesByPostID,
		},
		{
			Path:    prefix + "posts/:postid/comments/:commentid/reply",
			Method:  POST,
			Handler: handler.CreateReply,
		},
		{
			Path:    prefix + "posts/:postid/comments/:commentid/replies/:replyid",
			Method:  DELETE,
			Handler: handler.DeleteReplyByReplyID,
		},
	}
	return h
}
