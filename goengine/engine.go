package main

/*
#include <stdlib.h>
#include <string.h>
#include <stdint.h>
*/
import "C"
import (
	"log"
	"unsafe"
)

//export CallEngine
func CallEngine(inPtr unsafe.Pointer, inLen C.uint64_t, outPtr **C.uchar, outLen *C.uint64_t) C.int {

	if inPtr == nil || inLen == 0 {
		return 1
	}

	buf := C.GoBytes(inPtr, C.int(inLen))
	resp, err := dispatchFunction(buf)
	if err != nil || len(resp) == 0 {
		log.Printf("CallEngine: dispatchFunction error: %v", err)
		return 1
	}

	cbuf := C.malloc(C.size_t(len(resp)))
	if cbuf == nil {
		log.Printf("CallEngine: C.malloc failed")
		return 1
	}
	C.memcpy(cbuf, unsafe.Pointer(&resp[0]), C.size_t(len(resp)))

	*outPtr = (*C.uchar)(cbuf)
	*outLen = C.uint64_t(len(resp))
	return 0
}

//export FreeBuffer
func FreeBuffer(ptr unsafe.Pointer) {
	C.free(ptr)
}

func main() {}
