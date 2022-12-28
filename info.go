package main

import (
	"bytes"
	"crypto/sha1"
	"crypto/tls"
	"errors"
	"fmt"
	"math/big"
	"regexp"
	"strings"
	"time"

	"github.com/urfave/cli/v2"
)

func getCertInfo(c *cli.Context) error {
	// exit if no host argument
	if (c.String("host") == "") && (c.Args().First() == "") {
		return errors.New("missing argument host")
	}
	if (c.String("host") == "") && (c.Args().First() != "") {
		_ = c.Set("host", c.Args().First())
	}
	var address string
	var cfg tls.Config
	var versionLookup = map[uint16]string{
		tls.VersionTLS10: `VersionTLS1.0`,
		tls.VersionTLS11: `VersionTLS1.1`,
		tls.VersionTLS12: `VersionTLS1.2`,
		tls.VersionTLS13: `VersionTLS.1.3`,
	}

	// force client version, if specified
	switch c.String("tls") {
	case "1.0":
		cfg.MinVersion = tls.VersionTLS10
		cfg.MaxVersion = tls.VersionTLS10
	case "1.1":
		cfg.MinVersion = tls.VersionTLS11
		cfg.MaxVersion = tls.VersionTLS11
	case "1.2":
		cfg.MinVersion = tls.VersionTLS12
		cfg.MaxVersion = tls.VersionTLS12
	case "1.3":
		cfg.MinVersion = tls.VersionTLS13
		cfg.MaxVersion = tls.VersionTLS13
	}

	if c.String("host") != "" {
		address = fmt.Sprintf("%s:%d", c.String("host"), c.Int("port"))
	} else {
		address = fmt.Sprintf("%s:%d", c.Args().First(), c.Int("port"))
	}

	// ignore certificate errors if insecure flag is set
	if c.Bool("insecure") {
		cfg.InsecureSkipVerify = true
	}

	conn, err := tls.Dial("tcp", address, &cfg)
	if conn == nil && err != nil {
		return err
	}
	defer conn.Close()

	err = conn.VerifyHostname(c.String("host"))
	fmt.Printf("Certificate Verification\n------------------------\n")
	if err != nil {
		fmt.Printf("%s\n\n", err.Error())
	} else {
		fmt.Printf("Verified OK\n\n")
	}

	fmt.Printf("Connection Info\n---------------\n")
	fmt.Printf("Connection Protocol: %s \n", versionLookup[conn.ConnectionState().Version])
	fmt.Printf("Connection Ciphersuit: %s\n\n", tls.CipherSuiteName(conn.ConnectionState().CipherSuite))
	conn.OCSPResponse()

	fmt.Printf("Certificate Info\n----------------\n")
	certs := conn.ConnectionState().PeerCertificates
	for i, cert := range certs {

		if i == 0 {
			fmt.Printf("Server Name: %s\n", conn.ConnectionState().ServerName)

			// cert.SerialNumber.Text(16) would be sufficient, however for colon
			// seperated serial number using custom function
			//fmt.Printf("Serial Number: %s\n", cert.SerialNumber.Text(16))
			fmt.Printf("Serial Number: %s\n", colonedSerial(cert.SerialNumber))
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
		if i == 1 {
			fmt.Printf("\nCertificate Chain\n-----------------\n")
		}

		fmt.Printf("Common Name: %s \n", cert.Issuer.CommonName)
		fmt.Printf("Issuer Name: %s\n", cert.Issuer)
		fmt.Printf("Expiry: %s \n\n", cert.NotAfter.Format(time.RFC850))
	}
	return nil
}

func colonedSerial(i *big.Int) string {
	re := regexp.MustCompile("..")

	hex := fmt.Sprintf("%x", i)
	if len(hex)%2 == 1 {
		hex = "0" + hex
	}
	return strings.TrimRight(re.ReplaceAllString(hex, "$0:"), ":")
}
