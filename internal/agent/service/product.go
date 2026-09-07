package service

import (
	"context"
	"e-commerce/internal/agent/dto"
	"e-commerce/internal/agent/repository"
	"encoding/json"

	"github.com/bytedance/gopkg/util/logger"
)

type ProductListTool struct {
	repo repository.ProductRepository // 你的仓库
}

func NewProductListTool(repo repository.ProductRepository) *ProductListTool {
	return &ProductListTool{repo: repo}
}

// 工具名字
func (t *ProductListTool) Name() string {
	return "list_products"
}

// 👇 最重要：告诉AI要输出什么格式！
func (t *ProductListTool) Description() string {
	return `
你是一个商品查询工具，只负责解析用户查询条件并生成标准JSON。

规则：
1. 用户的查询条件必须全部转换为JSON字段，绝不遗漏！
2. 价格 > xxx  → 请在name或备注中带上，或直接解析为条件
3. 必须严格输出JSON，**不添加任何文字、解释、回答**
4. 工具返回什么，最终结果就输出什么
5. 绝对不要自己处理数据，绝对不要自己总结答案

支持字段：
page: 页码（默认1）
page_size: 每页条数（默认10）
name: 商品名称关键词
category_id: 分类ID
order_by: price/sales/create_time
order_desc: true/false
minprice: 价格筛选（如 >=3000）
`
}

// 👇 真正执行：input = AI输出的JSON字符串
func (t *ProductListTool) Call(ctx context.Context, input string) (string, error) {
	// 1. 把AI输出的JSON → 你的结构体

	logger.Info("AI 传给 Tool 的内容：", input)
	var req dto.ListProductsReq
	if err := json.Unmarshal([]byte(input), &req); err != nil {
		return "参数解析失败，请重试", nil
	}

	logger.Infof("filter: %+v", req)
	// 2. 默认值处理
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	// 3. 直接调用你的仓库函数！！！
	products, paging, err := t.repo.List(&req)
	if err != nil {
		return "查询商品失败", err
	}

	// 4. 结果转成JSON返回给AI
	result := map[string]any{
		"products":   products,
		"pagination": paging,
	}
	data, _ := json.MarshalIndent(result, "", "  ")
	return string(data), nil
}
