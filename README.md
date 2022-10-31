# apfw

Activitypub Forwarder

一个简单的 Activitypub 转发器，用于将 Activitypub 消息转发到其他 Activitypub 服务器。

## 配置文件

配置文件为 `config.json`

```json
{
    "token": "your token",
    "address": "127.0.0.1",
    "port": 8080
}
```

与二进制文件放在同一目录下即可。

`apfw.js`为cloudflare workers版本，需要在cloudflare worker中使用，但无法绕过1020。其中`TOKEN`需要在环境变量里设置。

## 基本原理

由于各种原因，不同的 Activitypub 服务器之间可能无法直接通信，因此需要一个中间服务器来转发消息。由于 Activitypub 请求会携带签名，需要把签名相关的`header`全部转发。在转发时，需要将`host`改为转发器的hostname以解决https证书匹配问题，同时将用于身份验证的token和最终接收消息的服务器地址放在`header`中。
