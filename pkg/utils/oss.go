package utils

import (
	"fmt"
	URL "net/url"
	"strings"
	"twitta/global"
)

func FulfillImageOSSPrefix(relativePath string) string {
	return getOSSPath(global.ServerConfig.StaticOSS.Domain, relativePath)
}

func TrimDomainPrefix(url string) string {
	//re, _ := regexp.Compile(`http(s)?.+\.com\/`)
	//return re.ReplaceAllString(url, "")

	// 去除域名前缀和?后缀
	urlObj, _ := URL.Parse(url)
	return urlObj.Path
}

func getOSSPath(domain, relativePath string) string {
	if relativePath == "" {
		return relativePath
	}
	dm := strings.TrimSuffix(domain, "/")
	rPath := strings.TrimPrefix(relativePath, "/")
	return fmt.Sprintf("%s/%s", dm, rPath)
}
