```
白山redis集群go语言版本使用规范库  
使用方法：
    //获取redis连接  
		rClient, err := bsredis.MustNewClient(redisOption)  
    //进行命令操作
    if res := rClient.Set("test_1", "hello redis test1", 0); res.Err() != nil {
        return err
    }
    rClient.Get("test_1")
    
```

