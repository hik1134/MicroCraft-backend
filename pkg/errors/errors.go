package errors

import (
	"errors"
	"net/http"
)

type Code string

const (
	OK Code = "OK"

	//通用
	INVALID_PARAM  Code = "INVALID_PARAM"
	INTERNAL_ERROR Code = "INTERNAL_ERROR"

	//Auth-EmailCode
	EMAIL_INVALID           Code = "EMAIL_INVALID"
	EMAIL_CODE_TOO_FREQUENT Code = "EMAIL_CODE_TOO_FREQUENT"
	EMAIL_CODE_SEND_FAIL    Code = "EMAIL_CODE_SEND_FAIL"
	EMAIL_CODE_GEN_FAIL     Code = "EMAIL_CODE_GEN_FAIL"
	CODE_EXPIRED            Code = "CODE_EXPIRED"
	CODE_INVALID            Code = "CODE_INVALID"

	//Auth-Register
	PASSWORD_LENGTH_INVALID Code = "PASSWORD_LENGTH_INVALID"
	USER_EXISTS             Code = "USER_EXISTS"
	PASSWORD_HASH_FAIL      Code = "PASSWORD_HASH_FAIL"
	DB_QUERY_FAIL           Code = "DB_QUERY_FAIL"
	DB_CREATE_FAIL          Code = "DB_CREATE_FAIL"

	//Auth-Login
	USER_NOT_FOUND     Code = "USER_NOT_FOUND"
	PASSWORD_INCORRECT Code = "PASSWORD_INCORRECT"
	JWT_SECRET_EMPTY   Code = "JWT_SECRET_EMPTY"
	JWT_GEN_FAIL       Code = "JWT_GEN_FAIL"

	//Auth-JWT
	AUTH_TOKEN_MISSING Code = "AUTH_TOKEN_MISSING"
	AUTH_TOKEN_INVALID Code = "AUTH_TOKEN_INVALID"
	AUTH_TOKEN_EXPIRED Code = "AUTH_TOKEN_EXPIRED"
	AUTH_FORBIDDEN     Code = "AUTH_FORBIDDEN"

	//基础设施
	CONFIG_NOT_INIT Code = "CONFIG_NOT_INIT"
	REDIS_NOT_INIT  Code = "REDIS_NOT_INIT"
	REDIS_OP_FAIL   Code = "REDIS_OP_FAIL"
	MYSQL_NOT_INIT  Code = "MYSQL_NOT_INIT"
	MYSQL_OPEN_FAIL Code = "MYSQL_OPEN_FAIL"
	DB_CONNECT_FAIL Code = "DB_CONNECT_FAIL"
	REDIS_INIT_FAIL Code = "REDIS_INIT_FAIL"
	REDIS_CONNECT_FAIL Code = "REDIS_CONNECT_FAIL"

	//Config
	CONFIG_LOAD_FAIL  Code = "CONFIG_LOAD_FAIL"
	CONFIG_PARSE_FAIL Code = "CONFIG_PARSE_FAIL"

	//Works-Upload
	WORK_TYPE_INVALID   Code = "WORK_TYPE_INVALID"
	WORK_FILE_REQUIRED  Code = "WORK_FILE_REQUIRED"
	WORK_SAVE_FILE_FAIL Code = "WORK_SAVE_FILE_FAIL"
	WORK_CREATE_FAIL    Code = "WORK_CREATE_FAIL"
	FILE_MISSING        Code = "FILE_MISSING"
	FILE_TOO_LARGE      Code = "FILE_TOO_LARGE"
	FILE_TYPE_NOT_ALLOW Code = "FILE_TYPE_NOT_ALLOW"
	SAVE_FILE_FAIL      Code = "SAVE_FILE_FAIL"
	DB_UPDATE_FAIL      Code = "DB_UPDATE_FAIL"

	//Works
	WORK_NOT_FOUND Code = "WORK_NOT_FOUND"
	WORK_FORBIDDEN Code = "WORK_FORBIDDEN"

	WORK_ALREADY_DELETED Code = "WORK_ALREADY_DELETED"

	//Posts
	POST_NOT_FOUND  Code = "POST_NOT_FOUND"
	POST_FORBIDDEN  Code = "POST_FORBIDDEN"

)

type AppError struct {
	Code Code
	Err  error 
}

func (e *AppError) Error() string {
	return string(e.Code)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func New(code Code) error {
	return &AppError{Code: code}
}

func Wrap(code Code, err error) error {
	if err == nil {
		return &AppError{Code: code}
	}
	return &AppError{Code: code, Err: err}
}

type ErrMeta struct {
	HTTPStatus int
	Message    string
}

var meta = map[Code]ErrMeta{
	INVALID_PARAM:  {http.StatusBadRequest, "参数错误"},
	INTERNAL_ERROR: {http.StatusInternalServerError, "服务器内部错误"},

	EMAIL_INVALID:           {http.StatusBadRequest, "邮箱格式不正确"},
	EMAIL_CODE_TOO_FREQUENT: {http.StatusTooManyRequests, "发送过于频繁，请稍后再试"},
	EMAIL_CODE_SEND_FAIL:    {http.StatusInternalServerError, "发送验证码失败"},
	EMAIL_CODE_GEN_FAIL:     {http.StatusInternalServerError, "生成验证码失败"},
	CODE_EXPIRED:            {http.StatusBadRequest, "验证码已过期或不存在"},
	CODE_INVALID:            {http.StatusBadRequest, "验证码错误"},

	PASSWORD_LENGTH_INVALID: {http.StatusBadRequest, "密码长度需为 8-20 位"},
	USER_EXISTS:             {http.StatusConflict, "邮箱已注册"},
	PASSWORD_HASH_FAIL:      {http.StatusInternalServerError, "密码处理失败"},
	DB_QUERY_FAIL:           {http.StatusInternalServerError, "数据库查询失败"},
	DB_CREATE_FAIL:          {http.StatusInternalServerError, "注册失败"},

	USER_NOT_FOUND:     {http.StatusBadRequest, "用户不存在"},
	PASSWORD_INCORRECT: {http.StatusBadRequest, "密码错误"},
	JWT_SECRET_EMPTY:   {http.StatusInternalServerError, "JWT配置错误"},
	JWT_GEN_FAIL:       {http.StatusInternalServerError, "生成token失败"},

	CONFIG_NOT_INIT: {http.StatusInternalServerError, "系统配置未初始化"},
	REDIS_NOT_INIT:  {http.StatusInternalServerError, "Redis未初始化"},
	REDIS_OP_FAIL:   {http.StatusInternalServerError, "Redis操作失败"},
	REDIS_INIT_FAIL: {http.StatusInternalServerError, "Redis初始化失败"},

	MYSQL_NOT_INIT:  {http.StatusInternalServerError, "MySQL未初始化"},
	MYSQL_OPEN_FAIL: {http.StatusInternalServerError, "MySQL连接失败"},
	DB_CONNECT_FAIL: {http.StatusInternalServerError, "连接数据库失败"},

	AUTH_TOKEN_MISSING: {http.StatusUnauthorized, "缺少登录凭证"},
	AUTH_TOKEN_INVALID: {http.StatusUnauthorized, "登录凭证无效"},
	AUTH_TOKEN_EXPIRED: {http.StatusUnauthorized, "登录已过期"},
	AUTH_FORBIDDEN:     {http.StatusForbidden, "无权限访问"},
	REDIS_CONNECT_FAIL: {http.StatusInternalServerError, "连接 Redis 失败"},

	CONFIG_LOAD_FAIL:  {http.StatusInternalServerError, "读取配置文件失败"},
	CONFIG_PARSE_FAIL: {http.StatusInternalServerError, "解析配置文件失败"},

	WORK_TYPE_INVALID:   {http.StatusBadRequest, "作品类型不合法"},
	WORK_FILE_REQUIRED:  {http.StatusBadRequest, "请上传图片文件"},
	WORK_SAVE_FILE_FAIL: {http.StatusInternalServerError, "保存图片失败"},
	WORK_CREATE_FAIL:    {http.StatusInternalServerError, "创建作品失败"},

	FILE_MISSING:        {http.StatusBadRequest, "请选择要上传的图片文件"},
	FILE_TOO_LARGE:      {http.StatusBadRequest, "图片过大（最大 8MB）"},
	FILE_TYPE_NOT_ALLOW: {http.StatusBadRequest, "图片类型不支持（仅 jpg/png/webp）"},
	SAVE_FILE_FAIL:      {http.StatusInternalServerError, "保存图片失败"},
	DB_UPDATE_FAIL:      {http.StatusInternalServerError, "更新作品信息失败"},

	WORK_NOT_FOUND: {http.StatusNotFound, "作品不存在"},
	WORK_FORBIDDEN: {http.StatusForbidden, "无权限访问该作品"},

	WORK_ALREADY_DELETED: {http.StatusBadRequest, "作品已删除"},

	POST_NOT_FOUND: {http.StatusNotFound, "帖子不存在"},
	POST_FORBIDDEN: {http.StatusForbidden, "无权限访问该帖子"},
}

func GetMeta(code Code) ErrMeta {
	if m, ok := meta[code]; ok {
		return m
	}
	return meta[INTERNAL_ERROR]
}

func GetCode(err error) Code {
	if err == nil {
		return OK
	}
	var ae *AppError
	if errors.As(err, &ae) {
		return ae.Code
	}
	return INTERNAL_ERROR
}