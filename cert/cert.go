package cert

import (
	"crypto/dsa"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"math/big"
	rd "math/rand"
	"os"
	"strings"
	"time"
)

// A StructuralError suggests that the ASN.1 data is valid, but the Go type
// which is receiving it doesn't match.
type StructuralError struct {
	Msg string
}

func (e StructuralError) Error() string { return "asn1: structure error: " + e.Msg }

// A SyntaxError suggests that the ASN.1 data is invalid.
type SyntaxError struct {
	Msg string
}

func (e SyntaxError) Error() string { return "asn1: syntax error: " + e.Msg }

// parseBase128Int parses a base-128 encoded int from the given offset in the
// given byte slice. It returns the value and the new offset.
func parseBase128Int(bytes []byte, initOffset int) (ret, offset int, err error) {
	offset = initOffset
	for shifted := 0; offset < len(bytes); shifted++ {
		if shifted == 4 {
			err = StructuralError{"base 128 integer too large"}
			return
		}
		ret <<= 7
		b := bytes[offset]
		ret |= int(b & 0x7f)
		offset++
		if b&0x80 == 0 {
			return
		}
	}
	err = SyntaxError{"truncated base 128 integer"}
	return
}

// parseObjectIdentifier parses an OBJECT IDENTIFIER from the given bytes and
// returns it. An object identifier is a sequence of variable length integers
// that are assigned in a hierarchy.
func parseObjectIdentifier(bytes []byte) (s []int, err error) {
	if len(bytes) == 0 {
		err = SyntaxError{"zero length OBJECT IDENTIFIER"}
		return
	}

	// In the worst case, we get two elements from the first byte (which is
	// encoded differently) and then every varint is a single byte long.
	s = make([]int, len(bytes)+1)

	// The first varint is 40*value1 + value2:
	// According to this packing, value1 can take the values 0, 1 and 2 only.
	// When value1 = 0 or value1 = 1, then value2 is <= 39. When value1 = 2,
	// then there are no restrictions on value2.
	v, offset, err := parseBase128Int(bytes, 0)
	if err != nil {
		return
	}
	if v < 80 {
		s[0] = v / 40
		s[1] = v % 40
	} else {
		s[0] = 2
		s[1] = v - 80
	}

	i := 2
	for ; offset < len(bytes); i++ {
		v, offset, err = parseBase128Int(bytes, offset)
		if err != nil {
			return
		}
		s[i] = v
	}
	s = s[0:i]
	return
}
func init() {
	rd.Seed(time.Now().UnixNano())
}

type CertInformation struct {
	Country            []string
	Organization       []string
	OrganizationalUnit []string
	EmailAddress       []string
	Province           []string
	Locality           []string
	CommonName         string
	CrtName, KeyName   string
	IsCA               bool
	Names              map[string]string
}

func GetCertExtProperty(cert *x509.Certificate, key string) []byte {
	ext := cert.Extensions
	if len(ext) > 0 {
		for i := 0; i < len(ext); i++ {
			key, err := parseObjectIdentifier([]byte(key))
			if err != nil {
				return nil
			}
			if ext[i].Id.Equal(key) {
				return ext[i].Value
			}
		}
	}
	return nil
}

func CreateCRT(RootCa *x509.Certificate, RootKey *rsa.PrivateKey, info CertInformation) error {
	Crt := newCertificate(info)
	Key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	var buf []byte
	if RootCa == nil || RootKey == nil {
		//创建自签名证书
		buf, err = x509.CreateCertificate(rand.Reader, Crt, Crt, &Key.PublicKey, Key)
	} else {
		//使用根证书签名
		buf, err = x509.CreateCertificate(rand.Reader, Crt, RootCa, &Key.PublicKey, RootKey)
	}
	if err != nil {
		return err
	}

	err = write(info.CrtName, "CERTIFICATE", buf)
	if err != nil {
		return err
	}

	buf = x509.MarshalPKCS1PrivateKey(Key)
	return write(info.KeyName, "PRIVATE KEY", buf)
}

//编码写入文件
func write(filename, Type string, p []byte) error {
	File, err := os.Create(filename)
	defer File.Close()
	if err != nil {
		return err
	}
	var b *pem.Block = &pem.Block{Bytes: p, Type: Type}
	return pem.Encode(File, b)
}

func Parse(crtPath, keyPath string) (rootcertificate *x509.Certificate, rootPrivateKey *rsa.PrivateKey, err error) {
	rootcertificate, err = ParseCrt(crtPath)
	if err != nil {
		return
	}
	rootPrivateKey, err = ParseKey(keyPath)
	return
}

func ParseCrt(path string) (*x509.Certificate, error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	p := &pem.Block{}
	p, buf = pem.Decode(buf)
	return x509.ParseCertificate(p.Bytes)
}

func ParseKey(path string) (*rsa.PrivateKey, error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	p, buf := pem.Decode(buf)
	return x509.ParsePKCS1PrivateKey(p.Bytes)
}

func newCertificate(info CertInformation) *x509.Certificate {
	cert := &x509.Certificate{
		SerialNumber: big.NewInt(rd.Int63()),
		Subject: pkix.Name{
			Country:            info.Country,
			Organization:       info.Organization,
			OrganizationalUnit: info.OrganizationalUnit,
			Province:           info.Province,
			CommonName:         info.CommonName,
			Locality:           info.Locality,
		},
		NotBefore:             time.Now(),                   //证书的开始时间
		NotAfter:              time.Now().AddDate(20, 0, 0), //证书的结束时间
		BasicConstraintsValid: true,                         //基本的有效性约束
		IsCA:           info.IsCA,                                                                  //是否是根证书
		ExtKeyUsage:    []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth}, //证书用途
		KeyUsage:       x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		EmailAddresses: info.EmailAddress,
	}
	for key, value := range info.Names {
		xi, err := parseObjectIdentifier([]byte(key))
		if err != nil {
			fmt.Printf("error=:", err.Error())
		}
		cert.ExtraExtensions = append(cert.ExtraExtensions, pkix.Extension{
			Id:       xi,
			Critical: true,
			Value:    []byte(value),
		})
	}
	return cert
}

// PrintCertInfo does a poor imitation of OpenSSL's `-text` output for certificates.
func PrintCertInfo(cert *x509.Certificate) {
	fmt.Printf("Certificate:\n")
	fmt.Printf("    Data:\n")
	fmt.Printf("        Version: %d\n", cert.Version)
	fmt.Printf("        Serial Number:\n")
	fmt.Printf("                % s\n", PrettyPrintBytes(cert.SerialNumber.Bytes(), 256)[0])
	fmt.Printf("    Signature Algorithm: %s\n", cert.SignatureAlgorithm)
	fmt.Printf("        Issuer: C=%s, O=%s, CN=%s, OU=%s\n", cert.Issuer.Country[0],
		cert.Issuer.Organization[0], cert.Issuer.CommonName, cert.Issuer.OrganizationalUnit[0])
	fmt.Printf("        Validity\n")
	fmt.Printf("            Not Before: %s\n", cert.NotBefore.UTC().Format(time.UnixDate))
	fmt.Printf("            Not After : %s\n", cert.NotAfter.UTC().Format(time.UnixDate))
	fmt.Printf("        Subject: %s\n", cert.Subject.CommonName)
	fmt.Printf("        Subject Public Key Info:\n")
	fmt.Printf("            Public Key Algorithm: %s\n", PublicKeyAlgorithm(int(cert.PublicKeyAlgorithm)))
	PublicKeyPrint(cert.PublicKey)
	fmt.Printf("        X509v3 extensions:\n")
	fmt.Printf("            X509v3 Key Usage: %v\n", cert.KeyUsage)
	fmt.Printf("            X509v3 Extended Key Usage: %v\n", cert.ExtKeyUsage)
	fmt.Printf("            X509v3 Basic Constraints: %v\n", cert.BasicConstraintsValid)
	if len(cert.SubjectKeyId) > 0 {
		fmt.Printf("            X509v3 Subject Key Identifier:\n")
		fmt.Printf("                % s\n", PrettyPrintBytes(cert.SubjectKeyId, 256)[0])
	}
	if len(cert.AuthorityKeyId) > 0 {
		fmt.Printf("            X509v3 Authority Key Identifier:\n")
		fmt.Printf("                % s\n", PrettyPrintBytes(cert.AuthorityKeyId, 256)[0])
	}
	fmt.Printf("        X509v3 Subject Alternative Name:\n")
	ext := cert.ExtraExtensions
	if len(ext) > 0 {
		for i := 0; i < len(ext); i++ {
			fmt.Printf("           ------- ExtraExtensions:%s=%s\n", ext[i].Id.String(), ext[i].Value)
		}
	}
	ids := cert.UnhandledCriticalExtensions
	if len(ids) > 0 {
		for i := 0; i < len(ids); i++ {
			fmt.Printf("           ------- ids:%s=%s\n", ids[i].String())
		}
	}
	ext = cert.Extensions
	if len(ext) > 0 {
		for i := 0; i < len(ext); i++ {
			fmt.Printf("           ------- Extensions:%s=%s\n", ext[i].Id.String(), string(ext[i].Value))
		}
	}
	if len(cert.DNSNames) > 0 {
		for i := 0; i < len(cert.DNSNames); i++ {
			fmt.Printf("            DNS:%s\n", cert.DNSNames[i])
		}
	}
	if len(cert.EmailAddresses) > 0 {
		for i := 0; i < len(cert.EmailAddresses); i++ {
			fmt.Printf("            Email:%s\n", cert.EmailAddresses[i])
		}
	}
	if len(cert.IPAddresses) > 0 {
		for i := 0; i < len(cert.IPAddresses); i++ {
			fmt.Printf("            DNS:%s\n", cert.IPAddresses[i])
		}
	}
	fmt.Printf("        X509v3 Certificate Policies:\n")
	if len(cert.PolicyIdentifiers) > 0 {
		for i := 0; i < len(cert.PolicyIdentifiers); i++ {
			fmt.Printf("            Policy: %v\n", cert.PolicyIdentifiers[i])
		}
	}
	fmt.Printf("        Authority Information Access:\n")
	if len(cert.OCSPServer) > 0 {
		for i := 0; i < len(cert.OCSPServer); i++ {
			fmt.Printf("            OCSP - %s\n", cert.OCSPServer[0])
		}
	}
	if len(cert.IssuingCertificateURL) > 0 {
		for i := 0; i < len(cert.IssuingCertificateURL); i++ {
			fmt.Printf("            CA Issuers - %s\n", cert.IssuingCertificateURL[0])
		}
	}
	fmt.Printf("     Signature Algorithm: %s\n", cert.SignatureAlgorithm)
	sig := PrettyPrintBytes(cert.Signature, 18)
	for i := 0; i < len(sig); i++ {
		fmt.Printf("         %s\n", sig[i])
	}

}

// PublicKeyAlgorithm is essentially a stringer for x509's PublicKeyAlgorith const.
func PublicKeyAlgorithm(n int) string {
	switch n {
	case 1:
		return "RSA"
	case 2:
		return "DSA"
	case 3:
		return "ECDSA"
	default:
		return "UnknownPublicKeyAlgorithm"
	}
}

// PublicKeyPrint prints info about the public key depending on its type.
// Currenty only RSA is supported.
func PublicKeyPrint(pub interface{}) {
	switch k := pub.(type) {
	case *rsa.PublicKey:
		fmt.Printf("                Public-Key: (%d bits)\n", k.N.BitLen())
		fmt.Printf("                Modulus:\n")
		key := PrettyPrintBytes(k.N.Bytes(), 16)
		for i := 0; i < len(key); i++ {
			fmt.Printf("                    %s\n", key[i])
		}
		fmt.Printf("                Exponent: %d\n", k.E)
	case *dsa.PublicKey:
		fmt.Printf("                DSA Key Information Not Implemented\n")
	case *ecdsa.PublicKey:
		fmt.Printf("                DSA Key Information Not Implemented\n")
	}
}

// PrettyPrintBytes breaks a sequence of bytes into lineLength pieces
// converts them to hex pairs, colon separated and returns them as []string.
func PrettyPrintBytes(b []byte, lineLength int) []string {
	var t [][]byte

	for i := 0; i < len(b)/lineLength; i++ {
		t = append(t, b[i*lineLength:(i+1)*lineLength-1])
	}
	if len(b)%lineLength > 0 {
		last := len(b) / lineLength
		t = append(t, b[last*lineLength:])
	}

	var s []string
	for i := 0; i < len(t); i++ {
		s = append(s, fmt.Sprintf("% X", t[i]))
		s[i] = strings.Replace(s[i], " ", ":", -1)
	}

	return s
}
