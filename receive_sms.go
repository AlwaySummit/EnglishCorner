package main

import (
	"encoding/base64"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/endpoints"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dybaseapi"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dybaseapi/mns"
)

const (
	mnsDomain = "1943695596114318.mns.cn-hangzhou.aliyuncs.com"
)

func main() {
	endpoints.AddEndpointMapping("cn-hangzhou", "Dybaseapi", "dybaseapi.aliyuncs.com")

	// 创建client实例
	client, err := dybaseapi.NewClientWithAccessKey(
		"cn-hangzhou",           // 您的可用区ID
		"<AccessKeyId>",         // 您的Access Key ID
		"<AccessKeySecret>")     // 您的Access Key Secret
	if err != nil {
		// 异常处理
		panic(err)
	}

	queueName := "<QueueName>"
	messageType := "<MessageType>"

	var token *dybaseapi.MessageTokenDTO

	for {
		if token == nil  {
			// 创建API请求并设置参数
			request := dybaseapi.CreateQueryTokenForMnsQueueRequest()
			request.MessageType = messageType
			request.QueueName = queueName
			// 发起请求并处理异常
			response, err := client.QueryTokenForMnsQueue(request)
			if err != nil {
				// 异常处理
				panic(err)
			}

			token = &response.MessageTokenDTO
		}

		mnsClient, err := mns.NewClientWithStsToken(
			"cn-hangzhou",
			token.AccessKeyId,
			token.AccessKeySecret,
			token.SecurityToken,
		)

		if err != nil {
			panic(err)
		}

		mnsRequest := mns.CreateBatchReceiveMessageRequest()
		mnsRequest.Domain = mnsDomain
		mnsRequest.QueueName = queueName
		mnsRequest.NumOfMessages = "10"
		mnsRequest.WaitSeconds = "5"

		mnsResponse, err := mnsClient.BatchReceiveMessage(mnsRequest)
		if err != nil {
			panic(err)
		}
		// fmt.Println(mnsResponse)

		receiptHandles := make([]string, len(mnsResponse.Message))
		for i, message := range mnsResponse.Message {
			messageBody, decodeErr := base64.StdEncoding.DecodeString(message.MessageBody)
			if decodeErr != nil {
				panic(decodeErr)
			}
			fmt.Println(string(messageBody))
			receiptHandles[i] = message.ReceiptHandle
		}
		if len(receiptHandles) > 0 {
			mnsDeleteRequest := mns.CreateBatchDeleteMessageRequest()
			mnsDeleteRequest.Domain = mnsDomain
			mnsDeleteRequest.QueueName = queueName
			mnsDeleteRequest.SetReceiptHandles(receiptHandles)
			//_, err = mnsClient.BatchDeleteMessage(mnsDeleteRequest) // 取消注释将删除队列中的消息
			if err != nil {
				panic(err)
			}
		}
	}
}