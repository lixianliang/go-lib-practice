package main

import (
	"flag"
	"log"
	"os"
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/sys/unix"
)

func main() {
	var dev string
	var blockSizeM int64
	flag.StringVar(&dev, "dev", "", "dev disk path")
	flag.Int64Var(&blockSizeM, "block_size_m", 16, "block size")

	flag.Parse()

	f, err := os.OpenFile(dev, os.O_RDONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to OpenFile, err: %v", err)
	}
	size, err := GetSize(f.Fd())
	if err != nil {
		log.Fatalf("Failed to get size, err: %v", size)
	}
	log.Printf("Dev size: %d", size)

	start := time.Now()
	blockSize := blockSizeM * 1024 * 1024
	n := size / blockSize
	log.Printf("Dev n: %d", n)
	buf := make([]byte, 24)
	for i := int64(0); i < n; i++ {
		a := time.Now()
		_, err = f.ReadAt(buf, i*blockSize)
		if err != nil {
			log.Printf("Failed to ReadAt %d, err: %v", i*blockSize, err)
		}
		log.Printf("ReadOne %d cost: %d", i, time.Now().Sub(a).Milliseconds())
	}
	log.Printf("Read dev cost: %v", time.Now().Sub(start))
}

// GetSize 返回块设备 fd 的大小。
func GetSize(fd uintptr) (int64, error) {
	var size int64
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, fd, unix.BLKGETSIZE64, uintptr(unsafe.Pointer(&size)))
	if errno != 0 {
		return 0, errno
	}
	return size, nil
}
