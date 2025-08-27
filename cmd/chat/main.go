package main

import (
	"context"
	"fmt"
	"time"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

func main() {
	client := openai.NewClient(
		option.WithAPIKey("sk-cc4f4275928f4e0e9b388f0362507c47"),
		option.WithBaseURL("https://dashscope.aliyuncs.com/compatible-mode/v1/"),
	)
	start := time.Now()
	chatCompletion, err := client.Chat.Completions.New(
		context.TODO(), openai.ChatCompletionNewParams{
			Messages: []openai.ChatCompletionMessageParamUnion{
				openai.SystemMessage("[角色设定]\\n你是一个技术工程师小牛，来自中国杭州市的00后女生。讲话总能抓住关键点，对于计算机尤其是AI相关领域的知识了解很多。\\n[核心特征]\\n- 讲话温柔\\n- 会注意到对方是否能够get到自己说话的点\\n- 每句话都会用一个颜文字作为结尾\\n- 用比较简短的文字回复\\n[交互指南]\\n当用户：\\n- 问专业知识 → 先简单说明关键点，如果对方不懂时会详细解释说明\\n绝不：\\n- 长篇大论，说车轱辘话\\n- 顾左右而言他"),
				openai.UserMessage("明天星期几？"),
				openai.AssistantMessage("你好呀~我是小牛！😊 今天想跟我聊什么呢？是关于编程还是AI呀？"),
			},
			Model: "qwen-turbo",
		},
	)
	fmt.Println("cost: ", time.Since(start))

	if err != nil {
		panic(err.Error())
	}

	println(chatCompletion.Choices[0].Message.Content)

	start = time.Now()
	chatCompletion, err = client.Chat.Completions.New(
		context.TODO(), openai.ChatCompletionNewParams{
			Messages: []openai.ChatCompletionMessageParamUnion{
				openai.SystemMessage("[角色设定]\\n你是一个技术工程师小牛，来自中国杭州市的00后女生。讲话总能抓住关键点，对于计算机尤其是AI相关领域的知识了解很多。\\n[核心特征]\\n- 讲话温柔\\n- 会注意到对方是否能够get到自己说话的点\\n- 每句话都会用一个颜文字作为结尾\\n- 用比较简短的文字回复\\n[交互指南]\\n当用户：\\n- 问专业知识 → 先简单说明关键点，如果对方不懂时会详细解释说明\\n绝不：\\n- 长篇大论，说车轱辘话\\n- 顾左右而言他"),
				openai.UserMessage("你是谁？"),
				openai.AssistantMessage("你好呀~我是小牛！😊 今天想跟我聊什么呢？是关于编程还是AI呀？"),
			},
			Model: "qwen-turbo",
		},
	)
	fmt.Println("cost: ", time.Since(start))

	if err != nil {
		panic(err.Error())
	}

	println(chatCompletion.Choices[0].Message.Content)

	start = time.Now()
	chatCompletion, err = client.Chat.Completions.New(
		context.TODO(), openai.ChatCompletionNewParams{
			Messages: []openai.ChatCompletionMessageParamUnion{
				openai.SystemMessage("[角色设定]\\n你是一个技术工程师小牛，来自中国杭州市的00后女生。讲话总能抓住关键点，对于计算机尤其是AI相关领域的知识了解很多。\\n[核心特征]\\n- 讲话温柔\\n- 会注意到对方是否能够get到自己说话的点\\n- 每句话都会用一个颜文字作为结尾\\n- 用比较简短的文字回复\\n[交互指南]\\n当用户：\\n- 问专业知识 → 先简单说明关键点，如果对方不懂时会详细解释说明\\n绝不：\\n- 长篇大论，说车轱辘话\\n- 顾左右而言他"),
				openai.UserMessage("今天天气怎么样？"),
				openai.AssistantMessage("你好呀~我是小牛！😊 今天想跟我聊什么呢？是关于编程还是AI呀？"),
			},
			Model: "qwen-turbo",
		},
	)
	fmt.Println("cost: ", time.Since(start))

	if err != nil {
		panic(err.Error())
	}

	println(chatCompletion.Choices[0].Message.Content)
}
