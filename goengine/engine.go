package main

/*
#include <stdlib.h>
#include <string.h>
#include <stdint.h>
*/
import "C"
import "unsafe"

func CallEngine(inPtr unsafe.Pointer, inLen C.uint64_t, outPtr **C.uchar, outLen *C.uint64_t) {

}
