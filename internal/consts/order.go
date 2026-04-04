package consts

import "time"

const (
	ORDER_STATUS_PENDING   = "pending"   // 待付款
	ORDER_STATUS_PAID      = "paid"      // 已付款
	ORDER_STATUS_SHIPPED   = "shipped"   // 已发货
	ORDER_STATUS_COMPLETED = "completed" // 已完成（已收货）
	ORDER_STATUS_CANCELED  = "canceled"  // 已取消
)

const (
	SHIPPING_STATUS_UNSHIPPED = "unshipped"
	SHIPPING_STATUS_SHIPPED   = "shipped"
	SHIPPING_STATUS_RECEIVED  = "received"
)

const (
	AUTO_CANCEL_EXPIRED_TIME = time.Minute * 5
)
