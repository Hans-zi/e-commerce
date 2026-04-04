package mq

import (
	"e-commerce/internal/order/dto"
	"encoding/json"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestKafka(t *testing.T) {
	topic := "test"
	x := 0

	producer := NewKafka()

	// 业务函数：前两次失败，第三次成功
	f := func(data []byte) error {
		fmt.Println("======================================")
		fmt.Println("当前 x =", x)

		if x != 2 {
			err := errors.New(fmt.Sprintf("失败！x = %d", x))
			x++
			return err
		}

		// 成功
		var req dto.PlaceOrderReq
		_ = json.Unmarshal(data, &req)
		fmt.Println("✅ 成功！消息内容：", req)
		return nil
	}
	groupID := "test-group" + time.Now().String()
	// 启动消费者
	go producer.StartConsumeMessages(topic, groupID, f)

	time.Sleep(1 * time.Second)

	// 发送消息
	var req dto.PlaceOrderReq
	for i := 0; i < 5; i++ {
		req.Lines = append(req.Lines, dto.PlaceOrderLineReq{
			ProductID: fmt.Sprintf("%d", i),
			Quantity:  i,
		})
	}

	err := producer.ProduceMessage(topic, req)
	require.NoError(t, err)
	fmt.Printf("消息发送成功: %v\n", req)

	// 等待足够久，让重试发生
	time.Sleep(10 * time.Second)
}
