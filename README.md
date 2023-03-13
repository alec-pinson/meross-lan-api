# meross-lan-api
 Allows control of meross devices over LAN via Rest API


## Configuration
Setup devices within `configuration.yaml`, see example.  
Also an environment variable is required for `KEY` which is your meross device key.  

## API Endpoints

`/deviceList`  
`/status/<device-name>`  
`/turnOn/<device-name>`  
`/turnOff/<device-name>`  

## Meross commands via Curl

`MessageId` - Random hex string characters in length.  
`Sign` - Sign is an MD5 hash of (MessageId + Your Key + Timestamp).  

**GET**  
```
curl -X POST -H "Content-Type: application/json; charset=UTF-8" http://192.168.1.175/config \
  -d '{
        "header": {
          "messageId": "af3fc9534c0e4727b5816bd5470fbb51",
          "namespace": "Appliance.System.All",
          "method": "GET",
          "payloadVersion": 1,
          "from": "Meross",
          "timestamp": 1678292549,
          "timestampMs": 0,
          "sign": "lll51743a3c75528k097ebfa921e0790"
        },
        "payload": {
          "all": {}
        }
      }'
```

**SET**  
```
curl -X POST -H "Content-Type: application/json; charset=UTF-8" http://192.168.1.45/config \
  -d '{
        "header": {
          "messageId": "a3c648b943cfbbd4c00b4800ce367c91",
          "method": "SET",
          "namespace": "Appliance.Control.ToggleX",
          "payloadVersion": 1,
          "sign": "h0436f1fe4523a5609623e3b9f881e61",
          "timestamp": 1678279320
        },
        "payload": {
          "togglex": {
            "channel": 0,
            "onoff": 0
          }
        }
      }'
```

## References
https://github.com/arandall/meross/blob/main/doc/protocol.md#headers
