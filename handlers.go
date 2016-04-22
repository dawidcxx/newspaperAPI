package main

import "net/http"
import "time"
import "github.com/gin-gonic/gin"
import jwt "github.com/dgrijalva/jwt-go"
import "strconv"
import "golang.org/x/crypto/bcrypt"

type standardCreatedResponse struct {
	ID int `json:"id"`
}

//<AUTH>
type authRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

//PostAuthRegister POST /auth/register
func PostAuthRegister(c *gin.Context) {
	var input authRequest

	if err := c.BindJSON(&input); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	newUser := User{Name: input.Username, Hash: string(hash)}

	res, err := DB.NamedQuery("INSERT INTO users (name, hash) VALUES (:name, :hash) RETURNING users.id", &newUser)

	if err != nil {
		c.Status(http.StatusConflict)
		return
	}

	// we expect only 1 row containting the insterd ID
	res.Next()

	var insertedID int

	if err := res.Scan(&insertedID); err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, standardCreatedResponse{ID: insertedID})

}

//PostAuthLogin POST /auth/login
func PostAuthLogin(c *gin.Context) {
	var input authRequest

	if err := c.BindJSON(&input); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	var user User

	err := DB.QueryRowx("SELECT * FROM users WHERE name=$1", input.Username).StructScan(&user)

	if err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Hash), []byte(input.Password)) != nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims["UserID"] = user.ID
	tokenStr, err := token.SignedString([]byte(SECRET))

	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenStr})

}

//</AUTH>

//<ARTICLE>
type articleRequest struct {
	Title     string    `json:"title" binding:"required"`
	Body      string    `json:"body" binding:"required"`
	PublishAt time.Time `json:"publishAt" binding:"required"`
}

//PostAPIArticle POST /api/article
func PostAPIArticle(c *gin.Context) {
	var input articleRequest

	if err := c.BindJSON(&input); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	var newArticleID int

	currentUserID, _ := c.Get("UserID")

	err := DB.QueryRow(`
    INSERT INTO articles (title, body, published_at, user_id)
    VALUES ($1, $2, $3, $4)
    RETURNING articles.id;
  `, input.Title, input.Body, input.PublishAt, currentUserID).Scan(&newArticleID)

	if err != nil {
		c.Status(http.StatusConflict)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": newArticleID})

}

//GetAPIArticle GET /api/article/:id<int>
func GetAPIArticle(c *gin.Context) {
	input, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	var out Article

	if err := DB.QueryRowx("SELECT * FROM articles WHERE id=$1", input).StructScan(&out); err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, out)

}

//PutAPIArticle PUT /api/article/:id<int>
func PutAPIArticle(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	var updatedArticle articleRequest

	if err := c.BindJSON(&updatedArticle); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	DB.Exec(`
    UPDATE articles 
    SET title=$1, body=$2, published_at=$3
    WHERE id=$4
  `, updatedArticle.Title, updatedArticle.Body, updatedArticle.PublishAt, id)

	c.Status(http.StatusNoContent)

}

//DeleteAPIArticle DELETE /api/article/:id<int>
func DeleteAPIArticle(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	var articleOwnerID int

	if err := DB.QueryRowx("SELECT user_id FROM articles WHERE id=$1", id).Scan(&articleOwnerID); err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	currentUserID, _ := c.Get("UserID")

	if articleOwnerID != currentUserID {
		c.Status(http.StatusUnauthorized)
		return
	}

	DB.Exec("DELETE FROM articles where id=$1", id)

	c.Status(http.StatusNoContent)

}

//</ARTICLE>
