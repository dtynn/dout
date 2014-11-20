### dout
=====

##### 说明
尝试用一个比较简单的方式来完成邮件(激活邮件, 重置密码邮件等)发送.  
使用这个服务会有一定的门槛, 比如使用你的开发语言来发起一个 http 请求, 以及添加 dns 记录等.  

##### 使用方法:
1.  下载安装:  
    ```
    go get -u github.com/dtynn/dout
    ```
    如果 `$GOPATH/bin` 还不在你的 `$PATH` 中, 那么你可能需要执行
    ```
    export PATH=$PATH:$GOPATH/bin
    ```
    或者将之加入你的终端配置文件中
    

1.  启动:
    ```
    >>> dout
    2014/11/20 16:59:38 [INFO][github.com/dtynn/dout] main.go:11: Dout Listen on :10025
    ```
    
    
1.  接口:  
    请求规格:
    ```
    Endpoint: http://127.0.0.1:10025/send
    Method: POST
    Content-Type: application/x-www-form-urlencoded;charset=utf-8
    Request Body:
        from=noreply@mydomain.com&to=a@example1.com,b@example2.com&body=mail_content
    ```
    请求参数:  
    
    参数名 | 类型 | 必选 | 说明 
    :--- | :------ | :--- | :---
    from | string  | 是   | 发送人地址
    to   | string  | 是   | 收件人地址, 多个收件地址以`,`分隔; 虽然接口允许多个收件人, 但由于会同步发送, 因此不建议这么做. 
    body | string  | 否   | 邮件正文, 不提供则会发送一封空白的邮件...
    
    使用 curl 的范例:
    ```
    curl -X POST -d "from=noreply@mydomain.com&to=a@example1.com,b@example2.com&body=test" http://127.0.0.1:10025/send
    ```
    
    响应格式:
    ```
    HTTP/1.1 298 status code 298
    Content-Length: 103
    Content-Type: application/json
    Date: Thu, 20 Nov 2014 08:22:34 GMT

    {
        "failed": [
            {
                "detail": "451 System Error.", 
                "email": "b@example2.com", 
                "error": 3
            }
        ], 
        "msg": "some failed"
    }
    ```
    
    可能出现的状态码含义:   
    
    状态码 | 含义  
    :--- | :-----
    200 | 请求成功, 邮件发送成功
    298 | 请求成功, 部分邮件发送失败, 发送失败的地址及原因会在`failed`字段中提供
    400 | 请求失败, 参数格式不正确, 如发送地址不合规格, 没有收件地址等
    
    `failed` 中出现的错误码及含义
    
    错误码 | 含义  
    :--- | :-----
    1    | 不合规格的收件人地址
    2    |  dns 查询失败
    3    |  smtp 发送失败
    
    
##### 附
尝试给 gmail 等较为严谨的邮件服务商发送邮件时, 需要至少对你的发送域名进行 spf 的设置, 参考:  
[https://support.google.com/a/answer/33786](https://support.google.com/a/answer/33786)  
及  
[如何正确发送Email—SPF,DKIM介绍与配置](http://mengzhuo.org/blog/%E5%A6%82%E4%BD%95%E6%AD%A3%E7%A1%AE%E5%8F%91%E9%80%81email-spfdkim%E4%BB%8B%E7%BB%8D%E4%B8%8E%E9%85%8D%E7%BD%AE.html)  

更为详尽的 gmail 设置指南参考:  
[群发邮件发件人指南](https://support.google.com/mail/answer/81126)
    