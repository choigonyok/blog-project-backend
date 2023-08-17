package model

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/google/uuid"
)

type Comment struct {
	ID       int    `json:"id"`
	PostID   int    `json:"postid"`
	Admin    int    `json:"admin"`
	Text     string `json:"text"`
	WriterID string `json:"writerid"`
	WriterPW string `json:"writerpw"`
}

type Post struct {
	ID        int       `json:"id"`
	Tag       string    `json:"tag"`
	Title     string    `json:"title"`
	Text      string    `json:"text"`
	WriteTime string 	`json:"writetime"`
	ImagePath string    `json:"imagepath"`
	ImageNum  int       `json:"imagenum"`
}

type Reply struct {
	ID        int    `json:"id"`
	Admin     int    `json:"admin"`
	WriterID  string `json:"writerid"`
	WriterPW  string `json:"writerpw"`
	Text      string `json:"text"`
	CommentID int    `json:"commentid"`
}

type Cookie struct {
	Value string `json:"value"`
}

type Visitor struct {
	Today int `json:"today"`
	Total int `json:"total"`
}

var db *sql.DB

func GetCookieValue(inputValue string) (string, error) {
	r, err := db.Query("SELECT value FROM cookie")
	if err != nil {
		return "", err
	}
	defer r.Close()
	var cookieValue string
	r.Next()
	r.Scan(&cookieValue)
	return cookieValue, nil
}

func UpdateCookieRecord() (uuid.UUID, error) {
	tx, err := db.Begin()
	if err != nil {
		return uuid.Nil, err
	}
	_, err =  db.Exec("DELETE FROM cookie")
	if err != nil {
		tx.Rollback()
		return uuid.Nil, err
	}
	cookieValue := uuid.New()
	_, err = db.Exec(`INSERT INTO cookie (value) VALUES ("`+cookieValue.String()+`")`)
	if err != nil {
		tx.Rollback()
		return uuid.Nil, err
	}
	err = tx.Commit()
	if err != nil {
		return uuid.Nil, err
	}
	return cookieValue, nil
}

func OpenDB(driverName, dataSourceName string) error {
	fmt.Println(driverName)
	fmt.Println(dataSourceName)
	database, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return err
	}

	db = database

	// DB와 서버가 연결 되었는지 확인
	err = db.Ping()
	if err != nil {
		return err
	}

	return nil
}

func CloseDB() error {
	err := db.Close()
	return err
}

func GetRecentPostID() (int, error) {
	var idnum int
	r, err := db.Query("SELECT id FROM post order by id desc limit 1")
	if err != nil {
		return 0, err
	}
	defer r.Close()
	for r.Next() {
		r.Scan(
			&idnum)
	}
	return idnum, nil
}
func AddPost(tag, title, text, writetime string) error {
	_, err := db.Exec(`INSERT INTO post (tag, writetime, title, text) values ('` + tag + `', '`+writetime+`' ,'` + title + `','` + text + `')`)
	return err
}

func UpdatePostImagePath(recentID int, imagename string) error {
	_, err := db.Query(`UPDATE post SET imgpath = '` + imagename + `' where id = ` + strconv.Itoa(recentID))
	return err
}

func UpdatePost(title, text, tag, postID string) error {
	_, err := db.Query(`UPDATE post SET title = '` + title + `', text = '` + text + `', tag ='` + tag + `'  where id = ` + postID)
	return err
}

func SelectPostByTag(tagSlice []string) ([]Post, error) {
	var data Post
	var datas []Post
	for _, v := range tagSlice {
		r, err := db.Query("SELECT id,tag,title,text,writetime,imgpath FROM post where tag LIKE '%" + v + "%' order by id desc")
		if err != nil {
			return nil, err
		}
		defer r.Close()
		for r.Next() {
			r.Scan(&data.ID, &data.Tag, &data.Title, &data.Text, &data.WriteTime, &data.ImagePath)
			datas = append(datas, data)
		}
	}

	return datas, nil
}

func GetEveryTagAsString() (string, error) {
	r, err := db.Query("SELECT tag FROM post group by tag")
	if err != nil {
		return "", err
	}
	defer r.Close()
	tagdata := Post{}
	sum := ""
	for r.Next() {
		r.Scan(&tagdata.Tag)
		sum += " " + tagdata.Tag
	}
	return sum, nil
}

func DeleteRecentPost() error {
	_, err := db.Query("DELETE FROM post ORDER BY id DESC LIMIT 1")
	return err
}

func DeletePostByPostID(postID string) error {
	_, err := db.Query("DELETE FROM post WHERE id = " + postID)
	return err
}

func SelectEveryCommentIDByPostID(postID string) ([]string, error) {
	r, err := db.Query("SELECT id FROM comment WHERE postid = " + postID)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	var commentsSlice []string
	var tempString string
	for r.Next() {
		r.Scan(&tempString)
		commentsSlice = append(commentsSlice, tempString)
	}
	return commentsSlice, nil
}

func DeleteCommentByCommentID(commentID string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec("DELETE FROM comment WHERE id = " + commentID)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = db.Exec("DELETE FROM reply WHERE commentid = " + commentID)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func GetRecentCommentID() (int, error) {
	r, err := db.Query("SELECT id FROM comment order by id desc limit 1")
	if err != nil {
		return 0, err
	}
	defer r.Close()
	var recentCommentID int
	for r.Next() {
		r.Scan(&recentCommentID)
	}
	return recentCommentID, nil
}

func InsertComment(postID, admin int, text, writerID, writerPW string) error {
	_, err := db.Query(`INSERT INTO comment(postid, text, writerid, writerpw, admin) values (` + strconv.Itoa(postID) + ",'" + text + "','" + writerID + "','" + writerPW + "'," + strconv.Itoa(admin) + ")")
	return err
}

func GetCommentWriterPWByCommentID(commentID string) (string, error) {
	r, err := db.Query("SELECT writerpw FROM comment WHERE id =" + commentID)
	if err != nil {
		return "", err
	}
	defer r.Close()
	var writerPW string
	r.Next()
	err = r.Scan(&writerPW)
	return writerPW, err
}

func SelectNotAdminWriterComment(postID int) ([]Comment, error) {
	r, err := db.Query(`SELECT writerid, writerpw, text, admin, id FROM comment WHERE admin != 1`)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	datas := []Comment{}
	data := Comment{}
	for r.Next() {
		r.Scan(&data.WriterID, &data.WriterPW, &data.Text, &data.Admin, &data.ID)
		data.PostID = postID
		datas = append(datas, data)
	}
	return datas, nil
}

func SelectCommentByPostID(postID int) ([]Comment, error) {
	r, err := db.Query(`SELECT writerid, writerpw, text, admin, id FROM comment WHERE postid = ` + strconv.Itoa(postID))
	if err != nil {
		return nil, err
	}
	defer r.Close()
	datas := []Comment{}
	data := Comment{}
	for r.Next() {
		r.Scan(&data.WriterID, &data.WriterPW, &data.Text, &data.Admin, &data.ID)
		data.PostID = postID
		datas = append(datas, data)
	}
	return datas, nil
}

func SelectReplyByCommentID(commentID string) ([]Reply, error) {
	r, err := db.Query("SELECT id, admin, writerid, writerpw, text FROM reply WHERE commentid = " + commentID + " order by id asc")
	if err != nil {
		return nil, err
	}
	defer r.Close()
	var datas []Reply
	var data Reply
	for r.Next() {
		r.Scan(&data.ID, &data.Admin ,&data.WriterID, &data.WriterPW, &data.Text)
		datas = append(datas, data)
	}
	return datas, nil
}

func GetRecentReplyID() (int, error) {
	r, err := db.Query("SELECT id FROM reply order by id desc limit 1")
	if err != nil {
		return 0, err
	}
	defer r.Close()
	var recentReplyID int
	for r.Next() {
		r.Scan(&recentReplyID)
	}
	return recentReplyID, nil
}

func InsertReply(admin int, commentID, text, writerID, writerPW string) error {
	_, err := db.Query("INSERT INTO reply (commentid, text, writerid, writerpw, admin) values (" + commentID + ",'" + text + "','" + writerID + "','" + writerPW + `',` + strconv.Itoa(admin) + `)`)
	return err
}

func GetReplyPWByReplyID(replyID string) (string, error) {
	r, err := db.Query("SELECT writerpw FROM reply WHERE id = " + replyID)
	if err != nil {
		return "", err
	}
	defer r.Close()
	var replyPW string
	r.Next()
	r.Scan(&replyPW)
	return replyPW, nil
}

func DeleteReplyByReplyID(replyID string) error {
	_, err := db.Query("DELETE FROM reply WHERE id = " + replyID)
	return err
}

func GetEveryPost() ([]Post, error) {
	r, err := db.Query("SELECT id, tag,title,text,writetime,imgpath FROM post ORDER BY id desc")
	if err != nil {
		return nil, err
	}
	defer r.Close()
	var datas []Post
	var data Post
	for r.Next() {
		r.Scan(&data.ID, &data.Tag, &data.Title, &data.Text, &data.WriteTime, &data.ImagePath)
		datas = append(datas, data)
	}
	return datas, nil
}
func GetPostByPostID(postID string) ([]Post, error) {
	r, err := db.Query("SELECT id, tag,title,text,writetime,imgpath FROM post where id = " + postID)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	var datas []Post
	var data Post
	for r.Next() {
		r.Scan(&data.ID, &data.Tag, &data.Title, &data.Text, &data.WriteTime, &data.ImagePath)
		datas = append(datas, data)
	}
	return datas, nil
}

func AddVisitorCount(visitor Visitor) error {
	_, err := db.Exec(`UPDATE visitor SET today = `+strconv.Itoa(visitor.Today+1)+`, total = `+strconv.Itoa(visitor.Total+1))
	if err != nil {
		return err
	}	
	return nil
}

func GetVisitorCount() (Visitor, error) {
	var visitor Visitor
	r, err := db.Query("SELECT today, total FROM visitor")
	if err != nil {
		return visitor, err
	}
	defer r.Close()
	if !r.Next() {
		_, err := db.Exec(`INSERT INTO visitor (today, total) VALUES (0, 0)`)
		if err != nil {
			return visitor, err
		}	
	}
	r.Scan(&visitor.Today, &visitor.Total)
	return visitor, nil
}

func GetTodayRecord() (string, error) {
	r, err := db.Query("SELECT date FROM visitor")
	if err != nil {
		return "", err
	}
	var today string
	r.Next()
	r.Scan(&today)
	return today, nil
}

func ResetTodayVisitorNum(date string) error {
	_, err := db.Exec(`UPDATE visitor SET today = 1, date = "`+date+`"`)
	return err
}