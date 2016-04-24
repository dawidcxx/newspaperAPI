package main

import "time"

//User struct representing a row in the users table
type User struct {
	ID   int    `db:"id" json:"id"`
	Name string `db:"name" json:"username"`
	Hash string `db:"hash" json:"-"`
}

//Article struct representing a row in the articles table
type Article struct {
	ID     int    `db:"id" json:"id"`
	Title  string `db:"title" json:"title"`
	Body   string `db:"body" json:"body"`
	UserID int    `db:"user_id" json:"-"`
	//the consumer should treat the publishing date as the creation date
	PublishAt time.Time `db:"published_at" json:"createdAt"`
	// don't leak the real creation date
	CreatedAt time.Time `db:"created_at" json:"-"`
}

//Comment struct representing a row in the comments table
type Comment struct {
	ID        int       `db:"id" json:"id"`
	Body      string    `db:"body" json:"body"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UserID    int       `db:"user_id" json:"-"`
}
