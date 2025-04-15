package main

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"unsafe"

	"filippo.io/age"
	"filippo.io/age/armor"
)

/*
#include <stdlib.h>
typedef const char constchar;
*/
import "C"

func main() {}

var (
  version = "unknown" // set with -X
  xcryptoversion string
  ageversion string
)

//export libcryptage_age_version
func libcryptage_age_version() *C.char {
   return C.CString(ageversion)
}
//export libcryptage_xcrypto_version
func libcryptage_xcrypto_version() *C.char {
   return C.CString(xcryptoversion)
}
//export libcryptage_version
func libcryptage_version() *C.char {
   return C.CString(version)
}

var _cerr *C.char = nil
var _ageerr error

//export ageerr
func ageerr() *C.char {
	if _cerr != nil {
		C.free(unsafe.Pointer(_cerr))
	}
	if _ageerr == nil {
		_cerr = nil
		return nil
	}
	_cerr = C.CString(_ageerr.Error())
	return _cerr
}

//export age_decrypt_armor
func age_decrypt_armor(privkey *C.constchar, text *C.constchar, flags C.int) *C.char {

	id, err := age.ParseX25519Identity(C.GoString(privkey))
	if err != nil {
		_ageerr = fmt.Errorf("Failed to parse public key %q: %v", len(C.GoString(privkey)), err)
		return nil
	}
	txt := C.GoString(text)
	var in io.Reader = strings.NewReader(txt)
	in = armor.NewReader(in)
	r, err := age.Decrypt(in, id)
	b, err := io.ReadAll(r)
	if err != nil {
		_ageerr = fmt.Errorf("reading: %v", err)
		return nil
	}
	return C.CString(string(b))

}

//export age_encrypt_armor
func age_encrypt_armor(pubkeys *C.constchar, text *C.constchar, flags C.int) *C.char {

	publicKeys := C.GoString(pubkeys)
	var recipients []age.Recipient
	for _, pubkey := range strings.Split(publicKeys, ",") {

		recipient, err := age.ParseX25519Recipient(pubkey)
		if err != nil {
			_ageerr = fmt.Errorf("Failed to parse public key %q: %v", publicKeys, err)
			return nil
		}

		recipients = append(recipients, recipient)
	}
	if len(recipients) == 0 {
		_ageerr = fmt.Errorf("no recip")
		return nil

	}
	var buf = &bytes.Buffer{}
	var armored = armor.NewWriter(buf)

	w, err := age.Encrypt(armored, recipients...)
	if err != nil {
		_ageerr = fmt.Errorf("Failed to create encrypted file: %v", err)
		return nil
	}
	if _, err := io.WriteString(w, C.GoString(text)); err != nil {
		_ageerr = fmt.Errorf("Failed to write to encrypted file: %v", err)
		return nil
	}
	if err := w.Close(); err != nil {
		_ageerr = fmt.Errorf("Failed to close encrypted file: %v", err)
		return nil
	}
	if err := armored.Close(); err != nil {
		_ageerr = fmt.Errorf("Failed to close encrypted file: %v", err)
		return nil
	}

	return C.CString(buf.String())
}

//export age_decrypt
func age_decrypt(privkey *C.constchar, text *C.constchar, flags C.int) *C.char {

	id, err := age.ParseX25519Identity(C.GoString(privkey))
	if err != nil {
		_ageerr = fmt.Errorf("Failed to parse public key %q: %v", len(C.GoString(privkey)), err)
		return nil
	}
	txt := C.GoString(text)
	var in io.Reader = strings.NewReader(txt)
	if strings.HasPrefix(txt, "age-encryption.org/v1\n") {
		in = armor.NewReader(in)
	}
	r, err := age.Decrypt(in, id)
	return nil
	b, err := io.ReadAll(r)
	if err != nil {
		return nil
	}
	return C.CString(string(b))

}

//export age_encrypt
func age_encrypt(pubkeys *C.constchar, text *C.constchar, flags C.int) *C.char {
	publicKeys := C.GoString(pubkeys)
	var recipients []age.Recipient
	for _, pubkey := range strings.Split(publicKeys, ",") {

		recipient, err := age.ParseX25519Recipient(pubkey)
		if err != nil {
			_ageerr = fmt.Errorf("Failed to parse public key %q: %v", publicKeys, err)
			return nil
		}

		recipients = append(recipients, recipient)
	}
	if len(recipients) == 0 {
		_ageerr = fmt.Errorf("no recip")
		return nil

	}
	var buf = &bytes.Buffer{}
	var out io.Writer = buf
	w, err := age.Encrypt(out, recipients...)
	if err != nil {
		_ageerr = fmt.Errorf("Failed to create encrypted file: %v", err)
		return nil
	}
	if _, err := io.WriteString(w, C.GoString(text)); err != nil {
		_ageerr = fmt.Errorf("Failed to write to encrypted file: %v", err)
		return nil
	}
	if err := w.Close(); err != nil {
		_ageerr = fmt.Errorf("Failed to close encrypted file: %v", err)
		return nil
	}
	return C.CString(buf.String())
}
