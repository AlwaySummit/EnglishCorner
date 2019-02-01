package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"time"
	"github.com/gorilla/mux"
	"io/ioutil"
	"encoding/json"
	"net/url"
	"ec/english-corner/common"
)

const (
	INVALID_BODY = 1
	INVALID_PARA = 2
	INNER_ERROR  = 3
)

type Req struct {
	Msg string `json:"msg"`
}

type Rsp struct {
	Status int64       `json:"status"`
	Msg    string      `json:"msg"`
	Data   interface{} `json:"data,omitempty"`
}

type Process func(query url.Values, body []byte, rsp *Rsp)

func main() {
	defer func() {
		if err := recover(); err != nil {
			debug.PrintStack()
		}
	}()

	common.Retrieve()

	Router := mux.NewRouter()
	//cgi router
	Router.HandleFunc("/", GenHandler(CheckoutToken))
	Router.HandleFunc("/hello", GenHandler(HelloServer))
	Router.HandleFunc("/attenders", GenHandler(AttendersHandler))
	Router.HandleFunc("/send_sms", GenHandler(SmsHandler))
	Router.HandleFunc("/send_receive_sms", GenHandler(ReceiveSmsHandler))

		Router.PathPrefix("/englishcorner/").Handler(http.StripPrefix("/englishcorner/",http.FileServer(http.Dir("englishcorner"))))
	svr := http.Server{
		Addr:         ":80",
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		Handler:      Router,
	}
	svr.ListenAndServe()
}

func GenHandler(pro Process) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, No-Cache, X-Requested-With, If-Modified-Since, Pragma, Last-Modified, Cache-Control, Expires, Content-Type, X-E4M-With")
		//url param
		query := r.URL.Query()
		//post data
		sbody, e := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if e != nil {
			w.WriteHeader(500)
			return
		}
		//process
		var rsp Rsp
		fmt.Printf("Req : %v\n", string(sbody))
		pro(query, sbody, &rsp)
		fmt.Printf("Rsp : %d, %v\n", rsp.Status, rsp.Data)

		//rsp
		buf, e := json.Marshal(&rsp)
		if e != nil {
			w.WriteHeader(500)
		}
		w.Write([]byte(buf))
	}
}

func ProcessEcho(query url.Values, body []byte, rsp *Rsp) {
	//解析url参数
	var isEcho bool
	echo := query["echo"]
	if echo != nil {
		//整型值使用strconv转换
		//state, e := strconv.ParseInt(echo[0], 10, 32)
		if echo[0] == "true" {
			isEcho = true
		} else {
			isEcho = false
		}
	}

	//解析post数据
	var req Req
	if e := json.Unmarshal(body, &req); e != nil {
		fmt.Printf("unmarshal fail %v [%v]\n", e, string(body))
		rsp.Status = INVALID_BODY
		rsp.Msg = "Invalid Req Body"
		return
	}

	//其他服务操作

	//回复数据
	if isEcho {
		rsp.Data = req.Msg
	}

	return
}
