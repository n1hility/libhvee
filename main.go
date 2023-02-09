package main

import (
	"bytes"
	"fmt"
	"github.com/baude/hyperv_kvp/ignition"
	"os"
)

func main() {
	//dice_name := "com.coreos.ignition.kvp."
	fn := "hyperv/ignition/test_configs/ignition1.json"
	file_content, err := os.ReadFile(fn)
	if err != nil {
		panic(err)
	}
	rdr := bytes.NewReader(file_content)
	//rdrLen := rdr.Len()

	//kvpf, err := hyperv.NewKVPFile(rdr)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}

	bs, err := ignition.Dice(rdr)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(len(bs), len(bs[0]))
	//fmt.Println(bs)
	//for i, s := range bs {
	//	outfile := fmt.Sprintf("%s%d", dice_name, i)
	//	if err := os.WriteFile(outfile, s, 0777); err != nil {
	//		fmt.Println(err)
	//		return
	//	}
	//	fmt.Printf("Wrote: %s of len %d\n", outfile, len(s))
	//}
	//
	//redone := ignition.Glue(bs)
	//fmt.Printf("input: %d, output: %d\n", rdrLen, len(redone))
	//fmt.Println(string(redone) == string(file_content))

}
