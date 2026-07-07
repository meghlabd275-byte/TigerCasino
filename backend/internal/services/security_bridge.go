package services

/*
#cgo LDFLAGS: -L${SRCDIR}/../../../security/target/release -ltigercasino_security
#include <stdlib.h>
#include <stdbool.h>

char* rust_hash_password(const char* password);
bool rust_verify_password(const char* password, const char* hash);
unsigned long long rust_generate_random(unsigned long long min, unsigned long long max);
void rust_free_string(char* s);
double rust_generate_outcome(const char* server_seed, const char* client_seed, int nonce);
*/
import "C"
import (
	"unsafe"
)

type SecurityBridge struct{}

func NewSecurityBridge() *SecurityBridge {
	return &SecurityBridge{}
}

func (b *SecurityBridge) HashPassword(password string) string {
	cPassword := C.CString(password)
	defer C.free(unsafe.Pointer(cPassword))

	cHash := C.rust_hash_password(cPassword)
	if cHash == nil {
		return ""
	}
	defer C.rust_free_string(cHash)

	return C.GoString(cHash)
}

func (b *SecurityBridge) VerifyPassword(password, hash string) bool {
	cPassword := C.CString(password)
	cHash := C.CString(hash)
	defer C.free(unsafe.Pointer(cPassword))
	defer C.free(unsafe.Pointer(cHash))

	return bool(C.rust_verify_password(cPassword, cHash))
}

func (b *SecurityBridge) GenerateRandom(min, max uint64) uint64 {
	return uint64(C.rust_generate_random(C.ulonglong(min), C.ulonglong(max)))
}

func (b *SecurityBridge) GenerateOutcome(serverSeed, clientSeed string, nonce int) float64 {
	cServer := C.CString(serverSeed)
	cClient := C.CString(clientSeed)
	defer C.free(unsafe.Pointer(cServer))
	defer C.free(unsafe.Pointer(cClient))

	return float64(C.rust_generate_outcome(cServer, cClient, C.int(nonce)))
}
