# Super Simple CORS Proxy Server


### Requirement
go version 1.15.5
### Quick Start

#### 1. Configuration .env file
```bash
# Reverse Proxy Target Host Information
TARGET_HOST=
# This Proxy Server Port
LOCAL_PORT=
# CORS Allow Origin Information. It is Client host information (i.e https://192.168.93.1:4000)
ACCESS_CONTROL_ALLOWS_ORIGIN=
# Credentials Flag for Sharing Cookie between origin resource server, Must Set as true
WITH_CREDENTIALS={The flag of share or not cookie value}
# Optional, Session Cookie Name In Resource Server (If, not using cookie, let that empty)
SESSION_COOKIE_NAME={Optional, Session Cookie name}
# SSL Key File Path For Running Proxy Server as Https
SSL_KEY_PATH=
# SSL Cert File Path For Running Proxy Server as Https
SSL_CERT_PATH=
```


#### 2. Run Server
```bash
go run main.go
```
