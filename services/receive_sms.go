package services

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/endpoints"
	"ec/english-corner/common"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dybaseapi"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dybaseapi/mns"
	"fmt"
	"encoding/base64"
	"encoding/json"
)

type QueueResp struct {
	PhoneNumber string `json:"phone_number"`
	Content     string `json:"content"`
	DestCode    string `json:"dest_code"`
	SendTime    string `json:"send_time"`
	SignName    string `json:"sign_name"`
	SequenceId  int64  `json:"sequence_id"`
}

func QueueWatcher() (msgArr []string) {
	endpoints.AddEndpointMapping(common.Credentials.Region, common.PRODUCT_ID, common.SMS_BASE_DOMAIN)

	client, err := dybaseapi.NewClientWithAccessKey(
		common.Credentials.Region,
		common.Credentials.AccessKey,
		common.Credentials.AccessKeySecret)
	if err != nil {
		panic(err)
	}

	var token *dybaseapi.MessageTokenDTO

	for {
		if token == nil {
			request := dybaseapi.CreateQueryTokenForMnsQueueRequest()
			request.MessageType = common.SMS_UP_MESSAGE_TYPE
			request.QueueName = common.Credentials.SmsQueueName
			response, err := client.QueryTokenForMnsQueue(request)
			if err != nil {
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

		receiptHandles := make([]string, len(mnsResponse.Message))
		for i, message := range mnsResponse.Message {
			messageBody, decodeErr := base64.StdEncoding.DecodeString(message.MessageBody)
			if decodeErr != nil {
				panic(decodeErr)
			}
			messageBodyStr := string(messageBody)
			fmt.Println(messageBodyStr)
			receiptHandles[i] = message.ReceiptHandle
			var queueRsp QueueResp
			err := json.Unmarshal(messageBody, &queueRsp)
			if err == nil && queueRsp.Content == "1" {
				msgArr = append(msgArr, queueRsp.PhoneNumber)
			}
		}

		if len(receiptHandles) > 0 {
			mnsDeleteRequest := mns.CreateBatchDeleteMessageRequest()
			mnsDeleteRequest.Domain = common.MNS_DOMAIN
			mnsDeleteRequest.QueueName = common.Credentials.SmsQueueName
			mnsDeleteRequest.SetReceiptHandles(receiptHandles)
			_, err = mnsClient.BatchDeleteMessage(mnsDeleteRequest) // 取消注释将删除队列中的消息
			if err != nil {
				panic(err)
			}
		}
		return msgArr
	}
}
