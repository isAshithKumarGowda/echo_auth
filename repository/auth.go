package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/isAshithKumarGowda/Echo_Auth/internals/models"
	"github.com/isAshithKumarGowda/Echo_Auth/pkg/database"
	"github.com/isAshithKumarGowda/Echo_Auth/pkg/utils"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"log"
	"strconv"
	"time"
)

const OTP_EXPIRE_TIME = 10

var ctx = context.Background()

var rdb = redis.NewClient(&redis.Options{
	Addr: "localhost:6379",
})

type AuthRepo struct {
	db *sql.DB
}

func NewAuthRepo(db *sql.DB) *AuthRepo {
	return &AuthRepo{
		db: db,
	}
}

func (ar *AuthRepo) Register(e echo.Context) (string, string, error) {
	userType := e.Param("type")

	var newUser models.AuthModel

	query := database.NewDBinstance(ar.db)

	if err := e.Bind(&newUser); err != nil {
		return "", "", fmt.Errorf("failed to decode request %v", err)
	}

	if err := query.IfUserExists(newUser.Email); err != nil {
		return "", "", fmt.Errorf("email already exists")
	}

	newUser.ID = uuid.NewString()
	if b := utils.StrongPasswordValidator(newUser.Password); !b {
		return "", "", fmt.Errorf("use a strong password")
	}

	limiter := utils.GetRateLimiter(newUser.Email)

	if !limiter.Allow() {
		return "", "", fmt.Errorf("Too many OTP requests. Please try again after 24hrs.")
	}

	otp, err := utils.GenerateOtp()

	if err != nil {
		return "", "", fmt.Errorf("error while generating the otp %v", err)
	}

	key := fmt.Sprintf("otp:%s", newUser.Email)

	err = rdb.Set(ctx, key, otp, time.Minute*OTP_EXPIRE_TIME).Err()
	if err != nil {
		return "", "", fmt.Errorf("error while seting otp in redis %v", err)
	}

	err = rdb.MSet(ctx, "user_name", newUser.Name, "user_id", newUser.ID).Err()
	if err != nil {
		return "", "", fmt.Errorf("error while seting user details in redis %v", err)
	}

	// if err = query.StoreOtp(newUser.Email, otp); err != nil {
	// 	return "", "", fmt.Errorf("error while storing the otp")
	// }

	hash, err := utils.HashPassword(newUser.Password)
	if err != nil {
		return "", "", fmt.Errorf("error while hashing password")
	}

	newUser.Password = hash

	token, err := utils.GenerateToken(newUser.ID, newUser.Name, userType, time.Now().Local().Add(365*24*time.Hour).Unix())
	if err != nil {
		return "", "", fmt.Errorf("error while generating token")
	}

	err = query.Register(newUser, userType)
	if err != nil {
		log.Println(err)
		return "", "", fmt.Errorf("error while registering the user in the database")
	}

	// go func() {
	// 	time.Sleep(time.Minute * OTP_EXPIRE_TIME)
	// 	if err = query.ClearOtp(newUser.Email, otp); err != nil {
	// 		log.Fatal(err)
	// 	}
	// }()

	if err = utils.SendOTP(newUser.Email, otp, strconv.Itoa(OTP_EXPIRE_TIME)); err != nil {
		return "", "", fmt.Errorf("error while sending otp")
	}
	return token, "Registration successful", nil
}

func (ar *AuthRepo) VeirfyEmail(e echo.Context) (string, string, error) {
	var user models.AuthVerifyModel

	param := e.Param("type")

	var dbUser models.AuthModel

	if err := e.Bind(&user); err != nil {
		return "", "", fmt.Errorf("failed to decode request: %v", err)
	}

	key := fmt.Sprintf("otp:%s", user.Email)

	val, err := rdb.Get(ctx, key).Result()

	if err == redis.Nil {
		return "", "", fmt.Errorf("otp not found")
	} else if err != nil {
		return "", "", err
	}

	if user.Otp != val {
		return "", "", fmt.Errorf("wrong otp")
	}

	dbu, err := rdb.MGet(ctx, "user_name", "user_id").Result()
	if err != nil {
		return "", "", err
	}

	if dbu[0] != nil {
		dbUser.Name = dbu[0].(string)
	}
	if dbu[1] != nil {
		dbUser.ID = dbu[1].(string)
	}

	go func() {
		query := database.NewDBinstance(ar.db)
		query.SetVerifiedEmail(user.Email)
	}()
	// b, err := query.IsOtpValid(user.Email, user.Otp)
	// if err != nil || !b {
	// 	return "", "", err
	// }

	// dbUser, err := query.GetUserDetails(user.Email, param)
	// if err != nil {
	// 	return "", "", err
	// }

	token, err := utils.GenerateToken(dbUser.ID, dbUser.Name, param, time.Now().Local().Add(10*time.Minute).Unix())
	if err != nil {
		return "", "", err
	}

	return token, "Email Verified", nil
}

func (ar *AuthRepo) Login(e echo.Context) (string, string, time.Time, error) {
	var user models.AuthLoginModel

	userType := e.Param("type")

	if err := e.Bind(&user); err != nil {
		return "", "", time.Now(), fmt.Errorf("failed to decode request %v", err)
	}

	db := database.NewDBinstance(ar.db)
	dbUser, login_time, err := db.Login(user, userType)
	if err != nil {
		return "", "", login_time, err
	}

	if err = utils.CheckPassword(user.Password, dbUser.Password); err != nil {
		return "", "", login_time, fmt.Errorf("wrong password %v", err)
	}

	access_token, err := utils.GenerateToken(dbUser.ID, dbUser.Name, userType, time.Now().Local().Add(time.Hour).Unix())
	if err != nil {
		return "", "", login_time, fmt.Errorf("error while generating token %v", err)
	}

	refresh_token, err := utils.GenerateToken(dbUser.ID, dbUser.Name, userType, time.Now().Local().Add(5*24*time.Hour).Unix())
	if err != nil {
		return "", "", login_time, fmt.Errorf("error while genereting token %v", err)
	}

	return access_token, refresh_token, login_time, nil
}

func (ar *AuthRepo) Logout(e echo.Context) error {
	var user models.AuthLogoutModel

	query := database.NewDBinstance(ar.db)

	if err := e.Bind(&user); err != nil {
		return fmt.Errorf("failed to decode request %v", err)
	}

	userType := e.Param("type")

	if err := query.Logout(userType, user.Email, user.Login); err != nil {
		return fmt.Errorf("error while Loging out")
	}

	return nil
}

func (ar *AuthRepo) GetLoginHistory(e echo.Context) ([]models.AuthLoginHistoryModel, error) {
	query := database.NewDBinstance(ar.db)
	userType := e.Param("type")

	total, err := query.GetTotalCount(userType)
	if err != nil {
		return nil, fmt.Errorf("error while counting total no. of logs")
	}

	page, _ := strconv.Atoi(e.QueryParam("page"))
	limit, _ := strconv.Atoi(e.QueryParam("limit"))

	if page <= 0 {
		page = 1
	}

	if limit <= 0 {
		limit = 10
	}

	offset := (page - 1) * limit

	if offset >= total {
		return nil, fmt.Errorf("error bad request")
	}

	users, err := query.GetLoginHistory(limit, offset, userType)
	if err != nil {
		return nil, fmt.Errorf("error while retreving data")
	}

	return users, err

}
