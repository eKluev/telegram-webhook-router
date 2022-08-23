# Telegram Webhook Router 1.0
#### It routes telegram requests to required IP:PORT
#### One SSL route for a couple of telegram bots.
https://hub.docker.com/r/kluev/telegram_webhook_router

## endpoints 👇
### /bot [POST]
default route for telegram requests
####
    request must include GET params: "route_ip", "route_port"
####
    example: https://..../bot?route_ip=127.0.0.1&route_port=8001
###
    moreover if you need to add extra params in request for your bot, then you can add its after required params
    theese params routes to your bot
####
    example: https://..../bot?route_ip=127.0.0.1&route_port=8001&token=...&user=...
###
### /setWebhook [POST]
simple setWebhook with a couple of params, use it for easy build a link and set webhook\
setWebhook by default sets up with address of webhook router
####
    Request must include POST params: "telegram_token", "ip", "port", "max_connections", "drop_pending_updates"
####
    example: 
    {
        "telegram_token": "0000000001:AAAAAAAAAA--TQLJvTpFLS6gLEgnC19ZOyR",
        "route_ip": "127.0.0.1",
        "route_port": 8001,
        "max_connections": 100,
        "drop_pending_updates": false,
        "extra_params": {"token": "sgDFSvvwWDFWVwe_01"}
    }

### /deleteWebhook [POST]

####
    Request must include POST param: "telegram_token"
####
    example: 
    {
        "telegram_token": "0000000001:AAAAAAAAAA--TQLJvTpFLS6gLEgnC19ZOyR",
    }

###
### response codes
+ 200: everything is ok
+ 400: validation error
+ 404: endpoint not found
+ 500: something wrong on server
