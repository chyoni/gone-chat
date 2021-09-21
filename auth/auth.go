package auth

import (
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v7"
	"github.com/twinj/uuid"
)

var Client *redis.Client

func Start() {
	dsn := os.Getenv("REDIS_DSN")
	Client = redis.NewClient(&redis.Options{
		Addr: dsn,
	})
	_, err := Client.Ping().Result()
	if err != nil {
		panic(err)
	}
}

func CreateToken(userID uint) (*TokenDetails, error) {
	td := &TokenDetails{}
	td.AtExpires = time.Now().Add(time.Minute * 15).Unix()
	td.AccessUUID = uuid.NewV4().String()

	td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	td.RefreshUUID = uuid.NewV4().String()

	var err error
	// Create Access Token
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUUID
	atClaims["user_id"] = userID
	atClaims["exp"] = td.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(os.Getenv("TOKEN_SECRET")))
	if err != nil {
		return nil, err
	}
	// Create Refresh Token
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUUID
	rtClaims["user_id"] = userID
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(os.Getenv("TOKEN_SECRET")))
	if err != nil {
		return nil, err
	}
	return td, nil
}

func CreateAuth(userID uint, td *TokenDetails) error {
	at := time.Unix(td.AtExpires, 0)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	err := Client.Set(td.AccessUUID, strconv.Itoa(int(userID)), at.Sub(now)).Err()
	if err != nil {
		return err
	}
	err = Client.Set(td.RefreshUUID, strconv.Itoa(int(userID)), rt.Sub(now)).Err()
	if err != nil {
		return err
	}
	return nil
}
