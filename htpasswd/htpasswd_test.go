package htpasswd

import (
	"io/ioutil"
	"testing"
)

func poe(err error) {
	if err != nil {
		panic(err)
	}
}

func tFile(name string) string {
	f, err := ioutil.TempFile("./", "test."+name+"")
	poe(err)
	return f.Name()
}

func fileContentsAre(t *testing.T, file string, contents string) {
	fileBytes, err := ioutil.ReadFile(file)
	poe(err)
	if contents != string(fileBytes) {
		t.Fatal("unexpected file contents", "should have been", contents, "was \""+string(fileBytes)+"\"")
	}
}
func TestEmptyHtpasswdFile(t *testing.T) {
	f := tFile("empty")
	SetPassword(f, "sha", "sha", HashSHA)
	SetPassword(f, "a", "a", HashBCrypt)
	SetPassword(f, "b", "b", HashMD5)
	//fileContentsAre(t, f, "sha:{SHA}2PRZAyDhNDqRW2OUFwZQqPNdaSY=\n")
}
func TestSetPasswordHash(t *testing.T) {
	f := tFile("set-hashes")

	poe(SetPasswordHash(f, "a", "a"))
	poe(SetPasswordHash(f, "b", "b"))
	poe(SetPasswordHash(f, "c", "c"))
	poe(RemoveUser(f, "b"))
	passwords, err := ParseHtpasswdFile(f)
	poe(err)
	if passwords["a"] != "a" {
		t.Fatal("a failed")
	}
	if passwords["b"] != "" {
		t.Fatal("b failed")
	}
	if passwords["c"] != "c" {
		t.Fatal("c failed")
	}
}
