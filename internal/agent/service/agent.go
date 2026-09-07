package service

import (
	"context"
	"e-commerce/internal/agent/dto"
	"e-commerce/internal/agent/repository"
	"encoding/json"
	"fmt"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/tmc/langchaingo/agents"
	"github.com/tmc/langchaingo/llms/openai"
	agentTools "github.com/tmc/langchaingo/tools"
)

type AgentService struct {
	executor *agents.Executor
}

func NewAgentService(productRepo repository.ProductRepository) *AgentService {
	llm, err := openai.New(
		openai.WithBaseURL("https://dashscope.aliyuncs.com/compatible-mode/v1"), // 你的Spark/OpenAI
		openai.WithToken("sk-6aed0708bae24083b313632de87796d9"),
		openai.WithModel("qwen-plus"))

	if err != nil {
		logger.Fatalf("Agent.New fail: %w", err)
		return nil
	}

	tools := []agentTools.Tool{
		NewProductListTool(productRepo),
	}

	agent := agents.NewOneShotAgent(llm, tools)
	executor := agents.NewExecutor(agent)
	return &AgentService{
		executor: executor,
	}
}

func (s *AgentService) Chat(ctx context.Context, userID string, msg string) (*dto.ChatRes, error) {
	// 1. 执行 Agent
	resultMap, err := s.executor.Call(ctx, map[string]any{
		"input": fmt.Sprintf("用户ID:%d，问题:%s", userID, msg),
	})
	if err != nil {
		return nil, err
	}

	// 2. 获取 AI 最终回答
	output := resultMap["output"].(string)

	// ===================================================================
	// 关键：
	// 让 Tool 直接返回结构化数据，我们在这里把 JSON → any
	// 这样 Data 就是对象/数组，不是字符串！
	// ===================================================================

	var data any
	err = json.Unmarshal([]byte(output), &data)

	var replyMsg string
	if err != nil {
		// 解析失败 → 说明是纯文本，没有数据
		replyMsg = output
		data = nil
	} else {
		// 解析成功 → 说明 output 是 Tool 返回的 JSON
		// 让 AI 生成一句自然话术（你也可以自定义）
		replyMsg = "已为你找到相关信息"
		
	}

	return &dto.ChatRes{
		Msg:  replyMsg, // AI 说话内容
		Data: data,     // 自动变成对象/数组
	}, nil
}
