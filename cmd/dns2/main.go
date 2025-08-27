package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"
)

func main() {
	var localIP, dnsServer, domain string
	flag.StringVar(&domain, "domain", "", "resolver domain")
	flag.StringVar(&localIP, "local_ip", "", "local ip")
	flag.StringVar(&dnsServer, "dns_server", "", "dns server, 8.8.8.8:53")

	flag.Parse()
	if domain == "" {
		log.Fatalf("domain invalid")
	}
	log.Printf("args %s %s %s", domain, localIP, dnsServer)

	// localIP 本地绑定的IP地址
	// dnsServer  指定DNS服务器地址
	// 自定义Resolver
	resolver := &net.Resolver{
		PreferGo: true, // 使用Go实现的DNS解析器
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			log.Printf("use remote addr: %s:%s", network, address)
			// 创建自定义Dialer，绑定本地IP
			dialer := &net.Dialer{
				Timeout: 5 * time.Second,
			}
			if localIP != "" {
				// 系统自动分配 port。
				dialer.LocalAddr = &net.UDPAddr{IP: net.ParseIP(localIP)}
			}
			// 明确连接到指定的DNS服务器，ddress 带 port
			if dnsServer != "" {
				return dialer.DialContext(ctx, network, dnsServer)
			}
			return dialer.DialContext(ctx, network, address)
		},
	}

	// 使用自定义Resolver解析域名
	start := time.Now()
	ips, err := resolver.LookupIPAddr(context.Background(), domain)
	if err != nil {
		fmt.Printf("Failed lookup, cost: %v, err: %v", time.Since(start), err)
		return
	}

	// 输出结果
	fmt.Printf("Success lookup, cost: %v", time.Since(start))
	for _, ip := range ips {
		fmt.Println(ip.IP)
	}
}
