# WemonFit

Weight monitoring for Fitbit, using InfluxDB and Grafana


## Generating a new certificate

```
openssl req  -new  -newkey rsa:2048  -nodes  -keyout localhost.key  -out localhost.csr
openssl  x509  -req  -days 365  -in localhost.csr  -signkey localhost.key  -out localhost.crt
```