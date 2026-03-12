package service

import (
	"context"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"
	"MicroCraft/internal/config"
	"MicroCraft/internal/dao/mysql"
	red "MicroCraft/internal/dao/redis"
	"MicroCraft/internal/model"
	perr "MicroCraft/pkg/errors"
	"MicroCraft/pkg/utils"

	goredis "github.com/redis/go-redis/v9"
)

const (
	codeTTL        = 5 * time.Minute
	rateLimitTTL   = 60 * time.Second
	codeExpireSecs = 300
)

type RegisterReq struct {
	Email    string `json:"email"`
	Code     string `json:"code"`
	Password string `json:"password"`
}

type RegisterResp struct {
	UserID   uint   `json:"user_id"`
	Email    string `json:"email"`
	Nickname string `json:"nickname"`
}

type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResp struct {
	Token         string `json:"token"`
	ExpireSeconds int64  `json:"expire_seconds"`
	User          struct {
		UserID   uint   `json:"user_id"`
		Email    string `json:"email"`
		Nickname string `json:"nickname"`
	} `json:"user"`
}

func redisKeyCode(email string) string {
	return fmt.Sprintf("microcraft:email_code:%s", email)
}
func redisKeyRate(email string) string {
	return fmt.Sprintf("microcraft:email_rate:%s", email)
}

func SendEmailCode(ctx context.Context, email string) (expireSeconds int, err error) {
	if red.RDB == nil {
		return 0, perr.New(perr.REDIS_NOT_INIT)
	}
	if config.Conf == nil {
		return 0, perr.New(perr.CONFIG_NOT_INIT)
	}
	rk := redisKeyRate(email)
	exists, err := red.RDB.Exists(ctx, rk).Result()
	if err != nil {
		return 0, perr.Wrap(perr.REDIS_OP_FAIL, err)
	}
	if exists == 1 {
		return 0, perr.New(perr.EMAIL_CODE_TOO_FREQUENT)
	}
	code, err := utils.Gen6DigitCode()
	if err != nil {
		return 0, perr.Wrap(perr.EMAIL_CODE_GEN_FAIL, err)
	}
	ck := redisKeyCode(email)
	if err := red.RDB.Set(ctx, ck, code, codeTTL).Err(); err != nil {
		return 0, perr.Wrap(perr.REDIS_OP_FAIL, err)
	}
	if err := red.RDB.Set(ctx, rk, "1", rateLimitTTL).Err(); err != nil {
		return 0, perr.Wrap(perr.REDIS_OP_FAIL, err)
	}
	ec := config.Conf.Email
	subject := "MicroCraft 邮箱验证码"
	body := fmt.Sprintf("你的验证码是：%s（5分钟内有效）", code)

	if err := utils.SendEmail(
		ec.Host, ec.Port,
		ec.Username, ec.Password,
		ec.From, email,
		subject, body,
	); err != nil {
		return 0, perr.Wrap(perr.EMAIL_CODE_SEND_FAIL, err)
	}

	return codeExpireSecs, nil
}

func Register(ctx context.Context, req RegisterReq) (*RegisterResp, error) {
	if config.Conf == nil {
		return nil, perr.New(perr.CONFIG_NOT_INIT)
	}
	if red.RDB == nil {
		return nil, perr.New(perr.REDIS_NOT_INIT)
	}
	email := strings.TrimSpace(req.Email)
	code := strings.TrimSpace(req.Code)
	pwd := strings.TrimSpace(req.Password)
	if email == "" || code == "" || pwd == "" {
		return nil, perr.New(perr.INVALID_PARAM)
	}
	if utf8.RuneCountInString(pwd) < 8 || utf8.RuneCountInString(pwd) > 20 {
		return nil, perr.New(perr.PASSWORD_LENGTH_INVALID)
	}
	ck := redisKeyCode(email)
	val, err := red.RDB.Get(ctx, ck).Result()
	if err != nil {
		if err == goredis.Nil {
			return nil, perr.New(perr.CODE_EXPIRED)
		}
		return nil, perr.Wrap(perr.REDIS_OP_FAIL, err)
	}
	if val != code {
		return nil, perr.New(perr.CODE_INVALID)
	}
	_, err = mysql.GetUserByEmail(email)
	if err == nil {
		return nil, perr.New(perr.USER_EXISTS)
	}
	if err != nil && !mysql.IsNotFound(err) {
		return nil, perr.Wrap(perr.DB_QUERY_FAIL, err)
	}
	hash, err := utils.HashPassword(pwd)
	if err != nil {
		return nil, perr.Wrap(perr.PASSWORD_HASH_FAIL, err)
	}
	u := &model.User{
		Email:        email,
		PasswordHash: hash,
		Nickname:     utils.NicknameFromEmail(email),
	}
	if err := mysql.CreateUser(u); err != nil {
		return nil, perr.Wrap(perr.DB_CREATE_FAIL, err)
	}
	_ = red.RDB.Del(ctx, ck).Err()

	return &RegisterResp{
		UserID:   u.ID,
		Email:    u.Email,
		Nickname: u.Nickname,
	}, nil
}

func Login(ctx context.Context, req LoginReq) (*LoginResp, error) {
	if config.Conf == nil {
		return nil, perr.New(perr.CONFIG_NOT_INIT)
	}

	email := strings.TrimSpace(req.Email)
	pwd := strings.TrimSpace(req.Password)
	if email == "" || pwd == "" {
		return nil, perr.New(perr.INVALID_PARAM)
	}

	u, err := mysql.GetUserByEmail(email)
	if err != nil {
		if mysql.IsNotFound(err) {
			return nil, perr.New(perr.USER_NOT_FOUND)
		}
		return nil, perr.Wrap(perr.DB_QUERY_FAIL, err)
	}

	if !utils.CheckPassword(pwd, u.PasswordHash) {
		return nil, perr.New(perr.PASSWORD_INCORRECT)
	}

	jc := config.Conf.Jwt
	if jc.Secret == "" {
		return nil, perr.New(perr.JWT_SECRET_EMPTY)
	}
	exp := jc.ExpireSeconds
	if exp <= 0 {
		exp = 7200
	}

	token, err := utils.GenToken(jc.Secret, u.ID, exp)
	if err != nil {
		return nil, perr.Wrap(perr.JWT_GEN_FAIL, err)
	}

	var resp LoginResp
	resp.Token = token
	resp.ExpireSeconds = exp
	resp.User.UserID = u.ID
	resp.User.Email = u.Email
	resp.User.Nickname = u.Nickname
	return &resp, nil
}