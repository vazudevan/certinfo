package main

import (
	"bytes"
	"crypto/sha1"
	"crypto/tls"
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {

	//cfg := tls.Config{MinVersion: tls.VersionTLS10, MaxVersion: tls.VersionTLS10}
  // for converting TLS version int constants to string
	var versionLookup = map[uint16]string{
		tls.VersionTLS10: `VersionTLS1.0`,
		tls.VersionTLS11: `VersionTLS1.1`,
		tls.VersionTLS12: `VersionTLS1.2`,
		tls.VersionTLS13: `VersionTLS.1.3`,
	}
  // declare flags and parse flag
	host := flag.String("host", "", "FQDN of the server")
	port := flag.Int("port", 443, "Port number")

	flag.Parse()
	if *host == "" {
		fmt.Printf("Usage: %s\n", os.Args[0])
		flag.PrintDefaults()
		return
	}

	address := fmt.Sprintf("%s:%d", *host, *port)
	conn, err := tls.Dial("tcp", address, nil)
	if err != nil {
		panic("Server doesn't support SSL certificate err: " + err.Error())
	}
	defer conn.Close()
  //verify certificate
	err = conn.VerifyHostname(*host)
	if err != nil {
		panic("Hostname doesn't match with certificate: " + err.Error())
	}
	fmt.Printf("Connection Info\n---------------\n")
	fmt.Printf("Connection Protocol: %s \n", versionLookup[conn.ConnectionState().Version])
	fmt.Printf("Connection Ciphersuit: %s\n\n", tls.CipherSuiteName(conn.ConnectionState().CipherSuite))

	fmt.Printf("Certificate Info\n----------------\n")
	certs := conn.ConnectionState().PeerCertificates
	for i, cert := range certs {
    // Leaf certificate
		if i == 0 {
			fmt.Printf("Serial Number: %s\n", cert.SerialNumber.Text(16))
			fingerprint := sha1.Sum(cert.Raw)

			var buf bytes.Buffer
			for i, f := range fingerprint {
				if i > 0 {
					fmt.Fprintf(&buf, ":")
				}
				fmt.Fprintf(&buf, "%02X", f)
			}
			fmt.Printf("Fingerprint SHA1: %s\n", buf.String())

			if len(cert.DNSNames) > 0 {
				fmt.Printf("SubjectAlternativeNames:\n")
				for _, n := range cert.DNSNames {
					fmt.Printf("  %s\n", n)
				}
			}
		}
    // Chain
		if i == 1 {
			fmt.Printf("\nCertificate Chain\n-----------------\n")
		}

		fmt.Printf("Common Name: %s \n", cert.Issuer.CommonName)
		fmt.Printf("Issuer Name: %s\n", cert.Issuer)
		fmt.Printf("Expiry: %s \n\n", cert.NotAfter.Format(time.RFC850))
	}
}
