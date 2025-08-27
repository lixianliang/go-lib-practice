package main

import (
	"context"
	"log"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
)

func main() {
	// 正确初始化客户端（带认证）

	// 创建SSE传输层
	// 通过 http 协议抓包  sudo tcpdump -i en0   'tcp port 80 and host mcp.amap.com' -A -s 0
	//c, err := client.NewSSEMCPClient("https://mcp.amap.com/sse?key=f9614b9763dcebe8a0717add00a0f2be")
	c, err := client.NewSSEMCPClient("http://mcp.amap.com/sse")
	if err != nil {
		log.Fatalf("Failed to new sse mcp client, err: %v", err)
	}

	// 必须先 start
	// 然后 Initialize => tool list
	ctx := context.Background()
	err = c.Start(ctx)
	if err != nil {
		log.Fatalf("Failed to start transport: %v", err)
	}

	initResult, err := c.Initialize(ctx, mcp.InitializeRequest{})
	if err != nil {
		log.Fatalf("Failed to init, err: %v", err)
	}
	log.Printf("init result: %v, result: %v", initResult, initResult.Result)

	/*toolResult, err := c.ListTools(context.Background(), mcp.ListToolsRequest{})
	if err != nil {
		log.Fatalf("Failed to list tool, err: %v", err)
	}

	log.Println("tool list:")
	for _, tool := range toolResult.Tools {
		log.Printf("name: %s, 描述：%s", tool.Name, tool.Description)
		log.Printf("input schema: %v, meta: %v", tool.InputSchema, tool.Meta)
	}*/

	log.Println("resouce list:")
	resResult, err := c.ListResources(ctx, mcp.ListResourcesRequest{})
	if err != nil {
		log.Fatalf("Failed to list resource, err: %v", err)
	}
	for _, res := range resResult.Resources {
		log.Printf("res: %v", res)
	}
}
