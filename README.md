# smtptest
Simple tool to test SMTP mail send with various settings including TLS1.1 downgrade
- All settings are configurable in the config.yaml file

```yaml
server: mymailserver.com
port: 25
user:
password:
TLS: # "", "StartTLS", "TLS"
TLSvalid: false # Validate TLS or allow any certificate
TLSmin: "SSL" # "SSL", "1.0", "1.1", "1.2", "1.3"
TLSmax: "1.3" # "SSL", "1.0", "1.1", "1.2", "1.3"
auth: # "", "PLAIN", "LOGIN"
from: myaddress@gmail.com
to: someaddress@gmail.com
subject: My test Email
body: Test Email notifications
```

