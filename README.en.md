[Russian version](README.md)

# fast_mtproxy
Fast deployment of the official MTProxy server with support for Fake TLS.

## Working with a script
### Download executable file
```bash
curl -L -o install https://git.io/JeOSr && chmod +x install
```

### Quick installation
```bash
./install
```

Default: port 443, random secret key generated, domain www.google.com

### Example
```bash
./install -p=8888 -H=443 -S=25c8dfee81acdadaff3a338a10db8497 -D=www.google.com
```
`-p` - local port for statistics (optional);  
`-H` - port to connect;  
`-S` - secret key;  
`-P` - ad tag (get here @MTProxybot);  
`-D` - TLS 1.3 domain (you can check <a href="https://www.cdn77.com/tls-test">here</a>).

For CentOS/RHEL holders you need to add a flag `-centos`

### Server removal
```bash
./install -uninstall=443
```

Flag `-uninstall` indicating the server port.

## Recommendations
The protocol, due to its features, mimics under TLS 1.3, therefore, it is necessary to use the appropriate ports for greater credibility (domain examples are shown in brackets):
* 261
* 271
* 324
* 443 (www.google.com, www.youtube.com)
* 465
* 563
* 636
* 853 (dns.google, cloudflare-dns.com)
* 989
* 990
* 992
* 993
* 994
* 995
* 4843
* 5061
* 5085
* 5349
* 5671
* 6513
* 6514
* 6619
* 8883

The best solution in terms of logic would be to use the hostname of your hosting provider, for example:
```
>nslookup 116.203.235.0

Name:    static.0.235.203.116.clients.your-server.de
Address:  116.203.235.0
```

`static.0.235.203.116.clients.your-server.de` - indicated in domain.