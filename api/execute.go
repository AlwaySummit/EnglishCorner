package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"sort"
	"net/url"
)

func HelloServer(query url.Values, body []byte, rsp *Rsp) {
	rsp.Data = "hello, world!\n"
}

func SmsHandler(query url.Values, body []byte, rsp *Rsp) {
	//SendSms()
	rsp.Data = "hello, world!\n"
}

func AttendersHandler(query url.Values, body []byte, rsp *Rsp) {
	//retrieve data from kv
	//TODO

	//return date
	rsp.Status = 200
	rsp.Msg = "Success"
	rsp.Data = GetAttenderResults()
	return
}

func CheckoutToken(query url.Values, body []byte, rsp *Rsp) {
	// token
	var token string = "iwuqing"
	// 获取参数
	signature := query.Get("signature")
	timestamp := query.Get("timestamp")
	nonce := query.Get("nonce")
	//将token、timestamp、nonce三个参数进行字典序排序
	var tempArray = []string{token, timestamp, nonce}
	sort.Strings(tempArray)
	//将三个参数字符串拼接成一个字符串进行sha1加密
	var sha1String string = ""
	for _, v := range tempArray {
		sha1String += v
	}
	h := sha1.New()
	h.Write([]byte(sha1String))
	sha1String = hex.EncodeToString(h.Sum([]byte("")))
	//获得加密后的字符串可与signature对比
	if sha1String != signature {
		fmt.Println("Authenticate failed")
	}
}
