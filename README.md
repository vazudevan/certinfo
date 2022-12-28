certinfo
A tool to display server certificate details.  With an option to choose Client TLS version. 

# Overview
Openssl can be used to obtain the certificate info, however it needs one to remember openssl command options and quriks.
The tool makes it simpler.  

`certinfo` is useful to check TLS Version supported by remote servers, you can also specify client TLS version to use.


# Usage
```
NAME:
   certinfo - TLS certificate info tool

USAGE:
   certinfo [global options] command [command options] [arguments...]

VERSION:
   1.1

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --host value   FQDN of server to get certificate information
   --port value   Port number (default: 443)
   --tls value    Force client TLS version. Valid values are 1.0 to 1.3
   --insecure     Ignore certificate errors (default: false)
   --help, -h     show help (default: false)
   --version, -v  print the version (default: false)
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
Limiting to specific client tls version
```
# tls 1.0
echo | openssl s_client -tls1 -connect example.com:443 2>&1

# tls 1.1
echo | openssl s_client -tls1_1 -connect example.com:443 2>&1

# tls 1.2
echo | openssl s_client -tls1_2 -connect example.com:443 2>&1

#tls 1.3
echo | openssl s_client -tls1_3 -connect example.com:443 2>&1
```