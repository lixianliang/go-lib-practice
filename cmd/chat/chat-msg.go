package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/sashabaranov/go-openai"
)

// 阿里云配置
const (
	baseURL   = "https://dashscope.aliyuncs.com/compatible-mode/v1"
	apiKey    = "YOUR_ALIYUN_API_KEY" // 替换为阿里云API Key
	modelName = "qwen-turbo"
)

// 定义可调用工具
var tools = []openai.Tool{
	{
		Type: openai.ToolTypeFunction,
		Function: &openai.FunctionDefinition{
			Name:        "get_current_weather",
			Description: "获取指定地点的当前天气",
			Parameters: json.RawMessage(`{
				"type": "object",
				"properties": {
					"location": {"type": "string", "description": "城市名称，如：北京"},
					"unit": {"type": "string", "enum": ["celsius", "fahrenheit"]}
				},
				"required": ["location"]
			}`),
		},
	},
	{
		Type: openai.ToolTypeFunction,
		Function: &openai.FunctionDefinition{
			Name:        "calculate_math",
			Description: "执行数学计算",
			Parameters: json.RawMessage(`{
				"type": "object",
				"properties": {
					"expression": {"type": "string", "description": "数学表达式，如：(12+3)*4/5"}
				},
				"required": ["expression"]
			}`),
		},
	},
}

func main() {
	// 1. 创建阿里云专用客户端
	config := openai.DefaultConfig(apiKey)
	config.BaseURL = baseURL
	client := openai.NewClientWithConfig(config)

	// 2. 用户消息
	messages := []openai.ChatCompletionMessage{
		{Role: openai.ChatMessageRoleUser, Content: "北京现在的天气怎么样？顺便计算(25-3)*4的结果"},
	}

	// 3. 发送请求（启用函数调用）
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    modelName,
			Messages: messages,
			Tools:    tools,
		},
	)
	if err != nil {
		log.Fatalf("请求失败: %v", err)
	}

	// 4. 检查响应
	if len(resp.Choices) == 0 {
		log.Fatalf("未收到响应")
	}

	assistantMsg := resp.Choices[0].Message
	toolCalls := assistantMsg.ToolCalls

	// 5. 处理函数调用
	if len(toolCalls) == 0 {
		fmt.Println("模型直接回复:", assistantMsg.Content)
		return
	}

	fmt.Printf("检测到 %d 个函数调用:\n", len(toolCalls))

	// 6. 执行函数调用
	messages = append(messages, assistantMsg) // 添加助手消息

	for _, call := range toolCalls {
		fmt.Printf("调用函数: %s (ID: %s)\n", call.Function.Name, call.ID)
		fmt.Printf("参数: %s\n", call.Function.Arguments)

		var result string
		switch call.Function.Name {
		case "get_current_weather":
			result = handleWeatherCall(call.Function.Arguments)
		case "calculate_math":
			result = handleMathCall(call.Function.Arguments)
		default:
			result = `{"error": "未知函数"}`
		}

		// 7. 添加工具响应消息
		messages = append(messages, openai.ChatCompletionMessage{
			Role:       openai.ChatMessageRoleTool,
			Content:    result,
			ToolCallID: call.ID,
			Name:       call.Function.Name,
		})
	}

	// 8. 获取最终回复
	finalResp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    modelName,
			Messages: messages,
		},
	)
	if err != nil {
		log.Fatalf("最终请求失败: %v", err)
	}

	if len(finalResp.Choices) > 0 {
		fmt.Println("\n最终回复:", finalResp.Choices[0].Message.Content)
	} else {
		fmt.Println("未收到最终回复")
	}
}

// 处理天气函数调用
func handleWeatherCall(args string) string {
	// 解析参数
	var params struct {
		Location string `json:"location"`
		Unit     string `json:"unit,omitempty"`
	}
	if err := json.Unmarshal([]byte(args), &params); err != nil {
		return `{"error": "参数解析失败"}`
	}

	// 模拟天气数据
	weatherData := map[string]string{
		"beijing":   "晴天, 25°C",
		"上海":        "多云, 28°C",
		"guangzhou": "阵雨, 30°C",
		"hangzhou":  "晴转多云, 26°C",
	}

	key := strings.ToLower(params.Location)
	if temp, ok := weatherData[key]; ok {
		return fmt.Sprintf(`{"location": "%s", "temperature": "%s", "unit": "celsius"}`, params.Location, temp)
	}
	return fmt.Sprintf(`{"error": "未知地点: %s"}`, params.Location)
}

// 处理数学计算
func handleMathCall(args string) string {
	// 解析参数
	var params struct {
		Expression string `json:"expression"`
	}
	if err := json.Unmarshal([]byte(args), &params); err != nil {
		return `{"error": "参数解析失败"}`
	}

	// 简单计算逻辑
	result := evaluateExpression(params.Expression)
	if result != "" {
		return fmt.Sprintf(`{"expression": "%s", "result": %s}`, params.Expression, result)
	}
	return fmt.Sprintf(`{"error": "无法计算表达式: %s"}`, params.Expression)
}

// 表达式计算器
func evaluateExpression(expr string) string {
	// 移除空格
	expr = strings.ReplaceAll(expr, " ", "")

	// 处理简单运算
	switch {
	case strings.Contains(expr, "+"):
		parts := strings.Split(expr, "+")
		if len(parts) != 2 {
			return ""
		}
		return fmt.Sprintf("%.2f", parseFloat(parts[0])+parseFloat(parts[1]))

	case strings.Contains(expr, "-"):
		parts := strings.Split(expr, "-")
		if len(parts) != 2 {
			return ""
		}
		return fmt.Sprintf("%.2f", parseFloat(parts[0])-parseFloat(parts[1]))

	case strings.Contains(expr, "*"):
		parts := strings.Split(expr, "*")
		if len(parts) != 2 {
			return ""
		}
		return fmt.Sprintf("%.2f", parseFloat(parts[0])*parseFloat(parts[1]))

	case strings.Contains(expr, "/"):
		parts := strings.Split(expr, "/")
		if len(parts) != 2 {
			return ""
		}
		divisor := parseFloat(parts[1])
		if divisor == 0 {
			return "0" // 避免除零错误
		}
		return fmt.Sprintf("%.2f", parseFloat(parts[0])/divisor)

	case strings.Contains(expr, "(") && strings.Contains(expr, ")"):
		// 简单处理括号表达式
		inner := expr[strings.Index(expr, "(")+1 : strings.Index(expr, ")")]
		innerResult := evaluateExpression(inner)
		if innerResult == "" {
			return ""
		}
		outer := strings.Replace(expr, "("+inner+")", innerResult, 1)
		return evaluateExpression(outer)

	default:
		// 尝试解析简单数字
		if val := parseFloat(expr); val != 0 {
			return fmt.Sprintf("%.2f", val)
		}
		return ""
	}
}

// 安全解析浮点数
func parseFloat(s string) float64 {
	var result float64
	_, err := fmt.Sscanf(s, "%f", &result)
	if err != nil {
		return 0
	}
	return result
}

/*
// chat message 结果
type ChatCompletionMessage struct {
	Role         string `json:"role"`
	Content      string `json:"content,omitempty"`
	Refusal      string `json:"refusal,omitempty"`
	MultiContent []ChatMessagePart

	// This property isn't in the official documentation, but it's in
	// the documentation for the official library for python:
	// - https://github.com/openai/openai-python/blob/main/chatml.md
	// - https://github.com/openai/openai-cookbook/blob/main/examples/How_to_count_tokens_with_tiktoken.ipynb
	Name string `json:"name,omitempty"`

	// This property is used for the "reasoning" feature supported by deepseek-reasoner
	// which is not in the official documentation.
	// the doc from deepseek:
	// - https://api-docs.deepseek.com/api/create-chat-completion#responses
	ReasoningContent string `json:"reasoning_content,omitempty"`

	FunctionCall *FunctionCall `json:"function_call,omitempty"`

	// For Role=assistant prompts this may be set to the tool calls generated by the model, such as function calls.
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`

	// For Role=tool prompts this should be set to the ID given in the assistant's prior request to call a tool.
	ToolCallID string `json:"tool_call_id,omitempty"`
}
*/

/*
简化为
Role
Content
ToolCalls
ToolCallID == role=tool tool 结果

ToolTypeFunction ToolType = "function"

FinishReasonStop          FinishReason = "stop"
FinishReasonLength        FinishReason = "length"
FinishReasonFunctionCall  FinishReason = "function_call"
FinishReasonToolCalls     FinishReason = "tool_calls"
FinishReasonContentFilter FinishReason = "content_filter"
FinishReasonNull          FinishReason = "null"

function_call 为老的版本使用方式，已经废弃，tools calls 支持多个工具的使用
*/
