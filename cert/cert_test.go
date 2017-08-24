package cert

import (
	"encoding/asn1"
	"fmt"
	"testing"
)

func TestCert05(t *testing.T) {
	a, x := asn1.Marshal("ssss")
	if x != nil {
		t.Log("Parse crt error,Error info:", x)
		return
	}
	fmt.Printf("%s", a)
	cert, err := ParseCrt("test_server.crt")
	if err != nil {
		t.Log("Parse crt error,Error info:", err)
		return
	}
	fmt.Printf("====%s====", GetCertExtProperty(cert, "p"))
	//PrintCertInfo(cert)
}
func TestCert04(t *testing.T) {
	fmt.Printf("==1==")
	cert, err := ParseCrt("test_server.crt")
	if err != nil {
		t.Log("Parse crt error,Error info:", err)
		return
	}
	fmt.Printf("==2==")
	err = cert.CheckSignature(cert.SignatureAlgorithm, cert.RawTBSCertificate, cert.Signature)
	if err != nil {
		fmt.Printf("Parse crt error,Error info:", err)
		return
	}
	fmt.Printf("==4==")
}
func TestCert02(t *testing.T) {
	fmt.Printf("==1==")
	crt, err := ParseCrt("test_server.crt")
	if err != nil {
		t.Log("Parse crt error,Error info:", err)
		return
	}
	fmt.Printf("==2==")
	parent, err := ParseCrt("test_root.crt")
	if err != nil {
		t.Log("Parse crt error,Error info:", err)
		return
	}
	fmt.Printf("==3==")
	err = crt.CheckSignatureFrom(parent)
	if err != nil {
		fmt.Printf("Parse crt error,Error info:", err)
		return
	}
	fmt.Printf("==4==")
}
func TestCert03(t *testing.T) {
	cert, err := ParseCrt("test_root.crt")
	//parent, err := ParseCrt("root_server.crt")
	if err != nil {
		t.Log("Parse crt error,Error info:", err)
		return
	}
	PrintCertInfo(cert)
	err = cert.CheckSignature(cert.SignatureAlgorithm, cert.RawTBSCertificate, cert.Signature)
	if err != nil {
		fmt.Printf("Parse crt error,Error info:", err)
		return
	}
	fmt.Printf("==")
	fmt.Printf(string(cert.RawTBSCertificate))
	fmt.Printf("==")
}
func TestCert01(t *testing.T) {
	baseinfo := CertInformation{Country: []string{"CN"}, Organization: []string{"WS"}, IsCA: true,
		OrganizationalUnit: []string{"work-stacks"}, EmailAddress: []string{"czxichen@163.com"},
		Locality: []string{"SuZhou"}, Province: []string{"JiangSu"}, CommonName: "Work-Stacks",
		CrtName: "test_root.crt", KeyName: "test_root.key"}

	err := CreateCRT(nil, nil, baseinfo)
	if err != nil {
		t.Log("Create crt error,Error info:", err)
		return
	}
	crtinfo := baseinfo
	crtinfo.IsCA = false
	crtinfo.CrtName = "test_server.crt"
	crtinfo.KeyName = "test_server.key"
	crtinfo.Names = map[string]string{"abcdefghijklmnopqrstuvwxyz": "p", "p": "bbbbbb"} //添加扩展字段用来做自定义使用

	crt, pri, err := Parse(baseinfo.CrtName, baseinfo.KeyName)
	if err != nil {
		t.Log("Parse crt error,Error info:", err)
		return
	}
	err = CreateCRT(crt, pri, crtinfo)
	if err != nil {
		fmt.Printf("Create crt error,Error info:", err)
	}
	//os.Remove(baseinfo.CrtName)
	//os.Remove(baseinfo.KeyName)
	//os.Remove(crtinfo.CrtName)
	//os.Remove(crtinfo.KeyName)
}
