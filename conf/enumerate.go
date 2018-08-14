// package conf 为配置相关.
package conf

import (
	"strings"

	"path/filepath"
	"strconv"

	"github.com/astaxie/beego"
)

// 登录用户的Session名
const LoginSessionName = "LoginSessionName"

const CaptchaSessionName = "__captcha__"

const RegexpEmail = "^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"

//允许用户名中出现点号

const RegexpAccount = `^[a-zA-Z][a-zA-z0-9\.]{2,50}$`

// PageSize 默认分页条数.
const PageSize = 10

// 用户权限
const (
	// 超级管理员.
	MemberSuperRole = 0
	//普通管理员.
	MemberAdminRole = 1
	//普通用户.
	MemberGeneralRole = 2
)

const (
	// 创始人.
	BookFounder = 0
	//管理者
	BookAdmin = 1
	//编辑者.
	BookEditor = 2
	//观察者
	BookObserver = 3
)

const (
	LoggerOperate   = "operate"
	LoggerSystem    = "system"
	LoggerException = "exception"
	LoggerDocument  = "document"
)
const (
	//本地账户校验
	AuthMethodLocal = "local"
	//LDAP用户校验
	AuthMethodLDAP = "ldap"
)

var (
	VERSION    string
	BUILD_TIME string
	GO_VERSION string
)

var (
	ConfigurationFile = "./conf/app.conf"
	WorkingDirectory  = "./"
	LogFile           = "./runtime/logs"
	BaseUrl           = ""
	AutoLoadDelay     = 0
	UploadDirectory   = "/tmp/mindoc"
)

// app_key
func GetAppKey() string {
	return beego.AppConfig.DefaultString("app_key", "godoc")
}

func GetDatabasePrefix() string {
	return beego.AppConfig.DefaultString("db_prefix", "md_")
}

//获取默认头像
func GetDefaultAvatar() string {
	return URLForWithCdnImage(beego.AppConfig.DefaultString("avatar", "/static/images/headimgurl.jpg"))
}

//获取阅读令牌长度.
func GetTokenSize() int {
	return beego.AppConfig.DefaultInt("token_size", 12)
}

//获取默认文档封面.
func GetDefaultCover() string {

	return URLForWithCdnImage(beego.AppConfig.DefaultString("cover", "/static/images/book.jpg"))
}

//获取允许的商城文件的类型.
func GetUploadFileExt() []string {
	ext := beego.AppConfig.DefaultString("upload_file_ext", "png|jpg|jpeg|gif|txt|doc|docx|pdf")

	temp := strings.Split(ext, "|")

	exts := make([]string, len(temp))

	i := 0
	for _, item := range temp {
		if item != "" {
			exts[i] = item
			i++
		}
	}
	return exts
}

// 获取上传文件允许的最大值
func GetUploadFileSize() int64 {
	size := beego.AppConfig.DefaultString("upload_file_size", "0")

	if strings.HasSuffix(size, "MB") {
		if s, e := strconv.ParseInt(size[0:len(size)-2], 10, 64); e == nil {
			return s * 1024 * 1024
		}
	}
	if strings.HasSuffix(size, "GB") {
		if s, e := strconv.ParseInt(size[0:len(size)-2], 10, 64); e == nil {
			return s * 1024 * 1024 * 1024
		}
	}
	if strings.HasSuffix(size, "KB") {
		if s, e := strconv.ParseInt(size[0:len(size)-2], 10, 64); e == nil {
			return s * 1024
		}
	}
	if s, e := strconv.ParseInt(size, 10, 64); e == nil {
		return s * 1024
	}
	return 0
}

//是否启用导出
func GetEnableExport() bool {
	return beego.AppConfig.DefaultBool("enable_export", true)
}

//同一项目导出线程的并发数
func GetExportProcessNum() int {
	exportProcessNum := beego.AppConfig.DefaultInt("export_process_num", 1)

	if exportProcessNum <= 0 || exportProcessNum > 4 {
		exportProcessNum = 1
	}
	return exportProcessNum
}

//导出项目队列的并发数量
func GetExportLimitNum() int {
	exportLimitNum := beego.AppConfig.DefaultInt("export_limit_num", 1)

	if exportLimitNum < 0 {
		exportLimitNum = 1
	}
	return exportLimitNum
}

//等待导出队列的长度
func GetExportQueueLimitNum() int {
	exportQueueLimitNum := beego.AppConfig.DefaultInt("export_queue_limit_num", 10)

	if exportQueueLimitNum <= 0 {
		exportQueueLimitNum = 100
	}
	return exportQueueLimitNum
}

//默认导出项目的缓存目录
func GetExportOutputPath() string {
	exportOutputPath := filepath.Join(beego.AppConfig.DefaultString("export_output_path", filepath.Join(WorkingDirectory, "cache")), "books")

	return exportOutputPath
}

//判断是否是允许商城的文件类型.
func IsAllowUploadFileExt(ext string) bool {

	if strings.HasPrefix(ext, ".") {
		ext = string(ext[1:])
	}
	exts := GetUploadFileExt()

	for _, item := range exts {
		if strings.EqualFold(item, ext) {
			return true
		}
	}
	return false
}

//重写生成URL的方法，加上完整的域名
func URLFor(endpoint string, values ...interface{}) string {
	baseUrl := beego.AppConfig.DefaultString("baseurl", "")
	pathUrl := beego.URLFor(endpoint, values...)

	if baseUrl == "" {
		baseUrl = BaseUrl
	}
	if strings.HasPrefix(pathUrl, "http://") {
		return pathUrl
	}
	if strings.HasPrefix(pathUrl, "/") && strings.HasSuffix(baseUrl, "/") {
		return baseUrl + pathUrl[1:]
	}
	if !strings.HasPrefix(pathUrl, "/") && !strings.HasSuffix(baseUrl, "/") {
		return baseUrl + "/" + pathUrl
	}
	return baseUrl + beego.URLFor(endpoint, values...)
}

func URLForWithCdnImage(p string) string {
	if strings.HasPrefix(p, "http://") || strings.HasPrefix(p, "https://") {
		return p
	}
	cdn := beego.AppConfig.DefaultString("cdnimg", "")
	//如果没有设置cdn，则使用baseURL拼接
	if cdn == "" {
		baseUrl := beego.AppConfig.DefaultString("baseurl", "")
		if baseUrl == "" {
			baseUrl = BaseUrl
		}
		if strings.HasPrefix(p, "/") && strings.HasSuffix(baseUrl, "/") {
			return baseUrl + p[1:]
		}
		if !strings.HasPrefix(p, "/") && !strings.HasSuffix(baseUrl, "/") {
			return baseUrl + "/" + p
		}
		return baseUrl + p
	}
	if strings.HasPrefix(p, "/") && strings.HasSuffix(cdn, "/") {
		return cdn + string(p[1:])
	}
	if !strings.HasPrefix(p, "/") && !strings.HasSuffix(cdn, "/") {
		return cdn + "/" + p
	}
	return cdn + p
}

func URLForWithCdnCss(p string) string {
	cdn := beego.AppConfig.DefaultString("cdncss", "")
	if strings.HasPrefix(p, "http://") || strings.HasPrefix(p, "https://") {
		return p
	}
	//如果没有设置cdn，则使用baseURL拼接
	if cdn == "" {
		baseUrl := beego.AppConfig.DefaultString("baseurl", "")
		if baseUrl == "" {
			baseUrl = BaseUrl
		}
		if strings.HasPrefix(p, "/") && strings.HasSuffix(baseUrl, "/") {
			return baseUrl + p[1:]
		}
		if !strings.HasPrefix(p, "/") && !strings.HasSuffix(baseUrl, "/") {
			return baseUrl + "/" + p
		}
		return baseUrl + p
	}
	if strings.HasPrefix(p, "/") && strings.HasSuffix(cdn, "/") {
		return cdn + string(p[1:])
	}
	if !strings.HasPrefix(p, "/") && !strings.HasSuffix(cdn, "/") {
		return cdn + "/" + p
	}
	return cdn + p
}

func URLForWithCdnJs(p string) string {
	cdn := beego.AppConfig.DefaultString("cdnjs", "")
	if strings.HasPrefix(p, "http://") || strings.HasPrefix(p, "https://") {
		return p
	}
	//如果没有设置cdn，则使用baseURL拼接
	if cdn == "" {
		baseUrl := beego.AppConfig.DefaultString("baseurl", "")
		if baseUrl == "" {
			baseUrl = BaseUrl
		}
		if strings.HasPrefix(p, "/") && strings.HasSuffix(baseUrl, "/") {
			return baseUrl + p[1:]
		}
		if !strings.HasPrefix(p, "/") && !strings.HasSuffix(baseUrl, "/") {
			return baseUrl + "/" + p
		}
		return baseUrl + p
	}
	if strings.HasPrefix(p, "/") && strings.HasSuffix(cdn, "/") {
		return cdn + string(p[1:])
	}
	if !strings.HasPrefix(p, "/") && !strings.HasSuffix(cdn, "/") {
		return cdn + "/" + p
	}
	return cdn + p
}

func WorkingDir(elem ...string) string {

	elems := append([]string{WorkingDirectory}, elem...)

	return filepath.Join(elems...)
}
