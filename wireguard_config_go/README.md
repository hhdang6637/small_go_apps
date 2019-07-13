# wireguard_config_go

**Build:**

> go build -o wg_key_mgt

This application help to generate wg config for all peer in your wireguard subnet:

**Create empty profile:**

> echo "[]" > wg.json

**Add master node:**

> ./wg_key_mgt add master
```
{
    "Name": "master",
    "PrivateKey": "wDw4otn7b2TUDRTls4QTTeWSJQHvMA8cFuceIY4/9Gg=",
    "PublicKey": "W24QTkA6UeBK5MMNdAnGFGbNdwqm8m4BQfm0yOIMQdQ=",
    "IPAddr": "10.0.0.1"
}
```

**Add peer node:**

> ./wg_key_mgt add peer_1

```
{
    "Name": "peer_1",
    "PrivateKey": "uHvB96Vd92wsMgnM8q9k56sNCGMcMrzDgS2qgScsY30=",
    "PublicKey": "z63tM5qo7gudjFZ4G9km3l8900JveT0BpsimZ9JRQaw=",
    "IPAddr": "10.0.0.2"
}

```


**Show all config data:**

> ./wg_key_mgt showconf
```
Global config:
{
    "Subnet": "10.0.0.0",
    "SubnetMask": 24,
    "ServerDomainName": "example.wireguard.com",
    "WgPort": 51820
}

Peer config:
[
    {
        "Name": "master",
        "PrivateKey": "wDw4otn7b2TUDRTls4QTTeWSJQHvMA8cFuceIY4/9Gg=",
        "PublicKey": "W24QTkA6UeBK5MMNdAnGFGbNdwqm8m4BQfm0yOIMQdQ=",
        "IPAddr": "10.0.0.1"
    },
    {
        "Name": "peer_1",
        "PrivateKey": "uHvB96Vd92wsMgnM8q9k56sNCGMcMrzDgS2qgScsY30=",
        "PublicKey": "z63tM5qo7gudjFZ4G9km3l8900JveT0BpsimZ9JRQaw=",
        "IPAddr": "10.0.0.2"
    }
]
```

**get master wg0.conf:**

> ./wg_key_mgt genconf master
```
# master
[Interface]
ListenPort = 51820
PrivateKey = wDw4otn7b2TUDRTls4QTTeWSJQHvMA8cFuceIY4/9Gg=
#PublicKey = W24QTkA6UeBK5MMNdAnGFGbNdwqm8m4BQfm0yOIMQdQ=
Address = 10.0.0.1/24
PostUp = iptables -I FORWARD 2 -i wg0 -j ACCEPT;
PostDown = iptables -D FORWARD 2 -i wg0 -j ACCEPT;


# peer_1
[Peer]
#PrivateKey = uHvB96Vd92wsMgnM8q9k56sNCGMcMrzDgS2qgScsY30=
PublicKey = z63tM5qo7gudjFZ4G9km3l8900JveT0BpsimZ9JRQaw=
AllowedIPs = 10.0.0.2/32
PersistentKeepalive = 30
```
**get peer_1 wg0.conf:**

> ./wg_key_mgt genconf  peer_1
```
[Interface]
PrivateKey = uHvB96Vd92wsMgnM8q9k56sNCGMcMrzDgS2qgScsY30=
#PublicKey = z63tM5qo7gudjFZ4G9km3l8900JveT0BpsimZ9JRQaw=
Address = 10.0.0.2/24

[Peer]
PublicKey = W24QTkA6UeBK5MMNdAnGFGbNdwqm8m4BQfm0yOIMQdQ=
AllowedIPs = 10.0.0.0/24
Endpoint = example.wireguard.com:51820
PersistentKeepalive = 30
```
