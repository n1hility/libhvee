package main

import (
	"fmt"
	"os"

	"github.com/baude/hyperv_kvp/hyperv"
)

func main() {
	sample := "/home/baude/kvp_sample"
	f, err := os.Open(sample)
	if err != nil {
		fmt.Println(err)
		return
	}
	kvpf, err := hyperv.NewKVPFile(f)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(kvpf.Keys())
	val, exists := kvpf.Get("HostinSystemOsMinor")
	fmt.Println(exists)
	fmt.Println(val)
}
