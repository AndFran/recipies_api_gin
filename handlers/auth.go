package handlers

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"os"
	"recipies_api_gin/models"
	"time"
)

type AuthHandler struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewAuthHandler(collection *mongo.Collection, ctx context.Context) *AuthHandler {
	return &AuthHandler{
		collection: collection,
		ctx:        ctx,
	}
}

func (h *AuthHandler) SignUpHandler(c *gin.Context) {
	hash := sha256.New()
	var user models.User
	if c.ShouldBind(&user) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid signup credentials"})
		return
	}

	if len(user.Password) < 5 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password is too short"})
		return
	}

	hash.Write([]byte(user.Password))
	pwd := hex.EncodeToString(hash.Sum(nil))
	cur := h.collection.FindOne(h.ctx, bson.M{"username": user.Username, "password": pwd})
	if cur.Err() == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user already exists"})
		return
	} else if !errors.Is(cur.Err(), mongo.ErrNoDocuments) {
		log.Println(cur.Err())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	id := primitive.NewObjectID()
	res, err := h.collection.InsertOne(h.ctx, bson.M{"username": user.Username, "password": pwd, "_id": id})
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Println("new user signup with id", res.InsertedID)

	expirationTime := time.Now().Add(10 * time.Minute)
	claims := &jwt.MapClaims{
		"sub": user.Username,
		"iss": "recipes_api_gin",
		"exp": expirationTime.Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"token": tokenString})
}

func (h *AuthHandler) SignInHandler(c *gin.Context) {
	// receives post and returns the token as json
	hash := sha256.New()
	var user models.User
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	/*
		if user.Username != "and" || user.Password != "fran" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}*/

	hash.Write([]byte(user.Password))
	pwd := hex.EncodeToString(hash.Sum(nil))

	cur := h.collection.FindOne(h.ctx, bson.M{"username": user.Username,
		"password": pwd})
	if cur.Err() != nil {
		log.Println(cur.Err())
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid username or password"})
		return
	}

	expirationTime := time.Now().Add(10 * time.Minute)
	/*claims := &Claims{
		Username: user.Username,
		Claims: map[string]string{
			"exp": strconv.FormatInt(expirationTime.Unix(), 10),
		},
	}*/

	claims := &jwt.MapClaims{
		"sub": user.Username,
		"iss": "recipes_api_gin",
		"exp": expirationTime.Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func (h *AuthHandler) RefreshJWTHandler(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil || token == nil || !token.Valid {
		if err != nil {
			log.Println("error refreshing token", err)
		}
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	claims := token.Claims.(jwt.MapClaims) // might need to map sub-claims
	if time.Unix(claims["exp"].(int64), 0).Sub(time.Now()) > 30*time.Second {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token not yet expired"})
		return
	}
	newExp := time.Now().Add(5 * time.Minute)
	claims["exp"] = newExp.Unix()
	refreshedToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err = refreshedToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
