package net

import (
	"log"
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInterface(t *testing.T) {
	// 可以获取网卡信息，根据网卡再获取 ip 地址。
	interfaces, err := net.Interfaces()
	require.NoError(t, err)
	for _, item := range interfaces {
		log.Println(item.Name, item)
		addrs, err := item.Addrs()
		require.NoError(t, err)
		log.Printf("addrs: %v", addrs)
		multiAddrs, err := item.MulticastAddrs()
		require.NoError(t, err)
		log.Printf("multiAddrs: %v", multiAddrs)
	}

	// 获取 ip 方式。
	addrs, err := net.InterfaceAddrs()
	require.NoError(t, err)
	for _, item := range addrs {
		log.Println(item.String(), item.Network(), item)
	}
}
