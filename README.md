certinfo
A tool to display server certificate details.  With an option to choose Client TLS version. 

# Overview
Openssl can be used to obtain the certificate info, however it needs one to remember openssl command options and quriks.
The tool makes it simpler.  

`certinfo` is useful to check TLS Version supported by remote servers, you can also specify client TLS version to use.


# Usage
```
Usage of certinfo:
  -host string
        FQDN of the server
  -port int
        Port number (default 443)
```
# Alternative

Same info obtained from openssl commands (e.g)
```
echo | openssl s_client -connect example.com:443 2>&1 
```
e.g. Extracting SAN from certificate
```
echo | openssl s_client -connect example.com:443 2>&1 | openssl x509 -noout -text |  awk -F, -v OFS="\n" '/DNS:/{gsub(/ *DNS:/, ""); $1=$1; print $0}'
```
