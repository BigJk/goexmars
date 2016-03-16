package goexmars

/*
#include <stdio.h>
#include <stdlib.h>
#include "pmars.h"
*/
import "C"
import "unsafe"

// Fight1Warrior lets one warrior fight against himself
func Fight1Warrior(w1 string, coresize int, cycles int, maxprocess int, rounds int, maxwarriorlen int, fixpos int) (resultWin1 int, resultWin2 int, resultEqual int) {
	var win1 int
	var win2 int
	var equal int
	cstrW1 := C.CString(w1)
	pWin1 := (*C.int)(unsafe.Pointer(&win1))
	pWin2 := (*C.int)(unsafe.Pointer(&win2))
	pEqual := (*C.int)(unsafe.Pointer(&equal))
	C.Fight1Warrior(cstrW1, C.int(coresize), C.int(cycles), C.int(maxprocess), C.int(rounds), C.int(maxwarriorlen), C.int(fixpos), pWin1, pWin2, pEqual)
	C.free(unsafe.Pointer(cstrW1))
	return win1, win2, equal
}

// Fight2Warriors lets two warrior fight
func Fight2Warriors(w1 string, w2 string, coresize int, cycles int, maxprocess int, rounds int, maxwarriorlen int, fixpos int) (resultWin1 int, resultWin2 int, resultEqual int) {
	var win1 int
	var win2 int
	var equal int
	cstrW1 := C.CString(w1)
	cstrW2 := C.CString(w2)
	pWin1 := (*C.int)(unsafe.Pointer(&win1))
	pWin2 := (*C.int)(unsafe.Pointer(&win2))
	pEqual := (*C.int)(unsafe.Pointer(&equal))
	C.Fight2Warriors(cstrW1, cstrW2, C.int(coresize), C.int(cycles), C.int(maxprocess), C.int(rounds), C.int(maxwarriorlen), C.int(fixpos), pWin1, pWin2, pEqual)
	C.free(unsafe.Pointer(cstrW1))
	C.free(unsafe.Pointer(cstrW2))
	return win1, win2, equal
}
