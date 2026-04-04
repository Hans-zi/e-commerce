package utils

import "fmt"

func GenerateSeckillProductCode(productID string) string {
	return fmt.Sprintf("seckill:%s", productID)
}
func GenerateSeckillOrderCode(userID, productID string) string {
	return fmt.Sprintf("%s:%s", userID, productID)
}
