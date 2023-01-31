package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/baude/hyperv_kvp/hyperv"
)

func main() {

	file_content, err := os.ReadFile(hyperv.KvpPool3)
	rdr := bytes.NewReader(file_content)

	kvpf, err := hyperv.NewKVPFile(rdr)
	if err != nil {
		fmt.Println(err)
		return
	}

	keys := kvpf.Keys()

	for _, k := range keys {
		val, _ := kvpf.Get(k)
		fmt.Println(k, val)
	}
}
