package utils

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"github.com/objcoding/wxpay"
	"sort"
	"strings"
)

func Sign(params map[string]string) string {
	// 创建切片
	var keys = make([]string, 0, len(params))
	// 遍历签名参数
	for k := range params {
		if k != "sign" { // 排除sign字段
			keys = append(keys, k)
		}
	}
	// 由于切片的元素顺序是不固定，所以这里强制给切片元素加个顺序
	sort.Strings(keys)
	//创建字符缓冲
	var buf bytes.Buffer
	for _, k := range keys {
		if len(params[k]) > 0 {
			buf.WriteString(k)
			buf.WriteString(`=`)
			buf.WriteString(params[k])
			buf.WriteString(`&`)
		}
	}
	// 加入apiKey作加密密钥
	buf.WriteString(`key=D3A010313F9E4DD3`)
	sha1 := sha1.New()
	sha1.Write([]byte(buf.String()))
	return strings.ToUpper(hex.EncodeToString(sha1.Sum([]byte(""))))
}

//签名函数 (包内的签名函数的改写)wx
func WxSign(params wxpay.Params) string {
	//params := make(wxpay.Params)
	/*params.SetString("package", "prepay_id="+prepayId).
	SetString("nonceStr", noncestr).
	SetString("timeStamp", timestamp).
	SetString("appId", viper.GetString("mini_program.app_id")).
	SetString("signType", "MD5")*/
	var keys = make([]string, 0, len(params))
	for k := range params {
		if k != "sign" { // 排除sign字段
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	//创建字符缓冲
	var buf bytes.Buffer
	for _, k := range keys {
		if len(params.GetString(k)) > 0 {
			buf.WriteString(k)
			buf.WriteString(`=`)
			buf.WriteString(params.GetString(k))
			buf.WriteString(`&`)
		}
	}
	// 加入apiKey作加密密钥
	buf.WriteString(`app_key=`)
	//
	//a95eceb1ac8c24ee28b70f7dbba912bf
	buf.WriteString("tfY3Ht2OrvkZ1zyt")
	var (
		dataMd5 [16]byte
		str     string
	)
	dataMd5 = md5.Sum(buf.Bytes())
	str = hex.EncodeToString(dataMd5[:])
	return strings.ToUpper(str)
}
