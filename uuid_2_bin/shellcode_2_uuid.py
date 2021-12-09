#coding=utf-8
import uuid

#Input your shellcode like:\xfc\x48\x83\xe4\xf0\xe8\xxx
buf = b"""Your shellcode"""
import uuid

def convertToUUID(shellcode):
    # If shellcode is not in multiples of 16, then add some nullbytes at the end
    if len(shellcode) % 16 != 0:
        print("[-] Shellcode's length not multiplies of 16 bytes")
        print("[-] Adding nullbytes at the end of shellcode, this might break your shellcode.")
        print("\n[*] Modified shellcode length: ", len(shellcode) + (16 - (len(shellcode) % 16)))

        addNullbyte = b"\x00" * (16 - (len(shellcode) % 16))
        shellcode += addNullbyte

    uuids = []
    for i in range(0, len(shellcode), 16):
        uuidString = str(uuid.UUID(bytes_le=shellcode[i:i + 16]))
        uuids.append(uuidString.replace("'", "\""))
    return uuids

u = convertToUUID(buf)
print(str(u).replace("'", "\""))
