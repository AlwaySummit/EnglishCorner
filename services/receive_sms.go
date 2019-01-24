package services

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/endpoints"
	"ec/english-corner/common"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dybaseapi"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dybaseapi/mns"
	"fmt"
	"encoding/base64"
)

func main() {
	endpoints.AddEndpointMapping(common.Credentials.Region, common.PRODUCT_ID, common.SMS_DOMAIN)

	// 创建client实例
	client, err := dybaseapi.NewClientWithAccessKey(
		common.Credentials.Region,           // 您的可用区ID
		common.Credentials.AccessKey,         // 您的Access Key ID
		common.Credentials.AccessKeySecret)     // 您的Access Key Secret
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
			request.QueueName = common.Credentials.SmsQueueName
			// 发起请求并处理异常
			response, err := client.QueryTokenForMnsQueue(request)
			if err != nil {
				// 异常处理
				panic(err)
			}

			token = &response.MessageTokenDTO
		}

		mnsClient, err := mns.NewClientWithStsToken(
			common.Credentials.Region,
			token.AccessKeyId,
			token.AccessKeySecret,
			token.SecurityToken,
		)

		if err != nil {
			panic(err)
		}

		mnsRequest := mns.CreateBatchReceiveMessageRequest()
		mnsRequest.Domain = common.MNS_DOMAIN
		mnsRequest.QueueName = common.Credentials.SmsQueueName
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
			mnsDeleteRequest.QueueName = common.Credentials.SmsQueueName
			mnsDeleteRequest.SetReceiptHandles(receiptHandles)
			//_, err = mnsClient.BatchDeleteMessage(mnsDeleteRequest) // 取消注释将删除队列中的消息
			if err != nil {
				panic(err)
			}
		}
	}
}