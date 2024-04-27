//远程过程调用RPC
package main

import (
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
	/*Uses the Gorilla Mux and Gorilla RPC libraries 
	to implement a JSON-RPC based Remote Procedure Call (RPC) service.
        */
)

type SmsArgs struct {
	Number, Content string
}//Parameters of SMS

type EmailArgs struct {
	To, Subject, Content string
}//Parameters of email

type Response struct {
	Result string
}//Parameters of response

//Define the service
type SmsService struct{}
type EmailService struct{}

func (t *SmsService) SendSMS(r *http.Request, args *SmsArgs, result *Response) error {
	*result = Response{Result: fmt.Sprintf("Sms sent to %s", args.Number)}
	return nil
}
/*
args *SmsArgs :Arg for sending SMS, 
contains the number of the SMS recipient and the content of the SMS.
*/
func (t *EmailService) SendEmail(r *http.Request, args *EmailArgs, result *Response) error {
	*result = Response{Result: fmt.Sprintf("Email sent to %s", args.To)}
	return nil
}

func main() {
	rpcServer := rpc.NewServer()
	
	//Register Codecs
	rpcServer.RegisterCodec(json.NewCodec(), "application/json")
	rpcServer.RegisterCodec(json.NewCodec(), "application/json;charset=UTF-8")
	
       	//new():Creates a zero value of this type
	sms := new(SmsService)
	email := new(EmailService)
	
	//用于将请求路由到相应的服务处理函数。
	rpcServer.RegisterService(sms, "sms")
	rpcServer.RegisterService(email, "email")

	router := mux.NewRouter()
	router.Handle("/delivery", rpcServer)
	http.ListenAndServe(":1337", router)
}
