package main

import (
	"fmt"
	"github.com/baude/hyperv_kvp/hyperv/ignition"
	"os"
	"unsafe"
)

//#include <stdlib.h>
//#include <string.h>
//extern char *read_kvp_data(void);
import "C"

func readString(ptr unsafe.Pointer, len uint32) string {
	return C.GoStringN((*C.char)(ptr), C.int(len))
}

func readData() map[string]string {
	cdata := C.read_kvp_data()
	if cdata == nil {
		return nil
	}
	defer C.free(unsafe.Pointer(cdata))

	ptr := unsafe.Pointer(cdata)

	values := make(map[string]string)
	for {
		keyLen := *((*uint32)(ptr))
		if keyLen == 0 {
			break
		}
		ptr = unsafe.Add(ptr, 4)

		key := readString(ptr, keyLen)

		ptr = unsafe.Add(ptr, keyLen)

		dataLen := *((*uint32)(ptr))
		ptr = unsafe.Add(ptr, 4)

		data := readString(ptr, dataLen)
		values[key] = data

		ptr = unsafe.Add(ptr, dataLen)
	}
	return values
}

func main() {
	m := readData()
	var (
		parts   ignition.Segments
		counter int
	)
	//for k, v := range readData() {
	//	fmt.Printf("%s: %s\n", k, v)
	//}

	for {
		key := fmt.Sprintf("com_coreos_ignition_kvp_%d", counter)
		p, exists := m[key]
		if !exists {
			break
		}
		parts = append(parts, []byte(p))
		counter++
	}
	if len(parts) < 1 {
		fmt.Println("no keys for ignition were found")
		os.Exit(1)
	}
	b := ignition.Glue(parts)
	fmt.Println(string(b))
}
