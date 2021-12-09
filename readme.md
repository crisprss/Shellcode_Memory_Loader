# Shellcode_Memory_Loader

## About 
相关资料和原理可以参考:
xxx

## Note
**注意**
*免杀效果不错,希望师傅们不跑沙箱,让这种方式存活的时间久一点!!!!*

## Description
基于Golang实现的Shellcode内存加载器,共实现3中内存加载shellcode方式,**UUID加载,MAC加载和IPv4加载**

结合`binject/universal`实现Golang的内存加载DLL方式,使用`AllocADsMem`实现内存申请,以加强免杀效果

简单的反沙箱机制,这里只是一个简单的Demo思路,后续再研究相关反沙箱的思路技术

## Usage
### UUID
在CS生成C版本的shellcode后填充到`shellcode_2_uuid.py`中:

![](https://md.byr.moe/uploads/upload_2b7c111c97ba77d8a854fd9e93c9b49f.png)

运行后得到转化后的UUID,全部填充到对应的go文件中:

![](https://md.byr.moe/uploads/upload_c548c00ef9e7f27e8139f82adc7306ab.png)

编译得到对应的可执行文件即可:
```golo
go build uuid_2_bin.go
```

**免杀效果**
![](https://md.byr.moe/uploads/upload_d575aa4b6bbad069384c4697aaef418a.png)

![](https://md.byr.moe/uploads/upload_f45720766c7c479e69be6b5f27c86367.png)

### MAC
在CS生成C版本的shellcode后填充到`shellcode_2_mac.py`中运行后会生成`mac_shell.txt`

![](https://md.byr.moe/uploads/upload_88e45ae35f6712a706f9e9ddeb3bfaba.png)

将其中的MAC地址填充到对应的go文件中:

![](https://md.byr.moe/uploads/upload_8c1af479b4685f5187cc13dc049db758.png)


编译得到对应的可执行文件即可:

```golo
go build mac_2_bin.go
```

**免杀效果**
![](https://md.byr.moe/uploads/upload_dfeb1e2b0716acb0cecbe57a134e2ea7.png)



![](https://md.byr.moe/uploads/upload_8a7f3abaaf6763ced0d9456bae79a3b4.png)

### IPv4
使用和MAC内存加载器一致,参考MAC加载器使用方式


