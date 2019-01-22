package main

import (
	"encoding/base64"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/endpoints"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dybaseapi"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dybaseapi/mns"
	"ec/english-corner/common"
)

func main() {
	endpoints.AddEndpointMapping(common.REGION_ID, common.PRODUCT_ID, common.SMS_DOMAIN)

	// 创建client实例
	client, err := dybaseapi.NewClientWithAccessKey(
		common.REGION_ID,           // 您的可用区ID
		common.ACCESS_KEY,         // 您的Access Key ID
		common.ACCESS_SECRET)     // 您的Access Key Secret
	if err != nil {
		// 异常处理
		panic(err)
	}

	var token *dybaseapi.MessageTokenDTO

	for {
		if token == nil  {
			// 创建API请求并设置参数
			request := dybaseapi.CreateQueryTokenForMnsQueueRequest()
			request.MessageType = common.SMS_UP_MESSAGE_TYPE
			request.QueueName = common.SMS_QUEUE_NAME
			// 发起请求并处理异常
			response, err := client.QueryTokenForMnsQueue(request)
			if err != nil {
				// 异常处理
				panic(err)
			}

			token = &response.MessageTokenDTO
		}

		mnsClient, err := mns.NewClientWithStsToken(
			common.REGION_ID,
			token.AccessKeyId,
			token.AccessKeySecret,
			token.SecurityToken,
		)

		if err != nil {
			panic(err)
		}

		mnsRequest := mns.CreateBatchReceiveMessageRequest()
		mnsRequest.Domain = common.MNS_DOMAIN
		mnsRequest.QueueName = common.SMS_QUEUE_NAME
		mnsRequest.NumOfMessages = "10"
		mnsRequest.WaitSeconds = "5"

		mnsResponse, err := mnsClient.BatchReceiveMessage(mnsRequest)
		if err != nil {
			fmt.Printf("err: %+v", err)
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
			mnsDeleteRequest.Domain = common.MNS_DOMAIN
			mnsDeleteRequest.QueueName = common.SMS_QUEUE_NAME
			mnsDeleteRequest.SetReceiptHandles(receiptHandles)
			//_, err = mnsClient.BatchDeleteMessage(mnsDeleteRequest) // 取消注释将删除队列中的消息
			if err != nil {
				panic(err)
			}
		}
	}
}