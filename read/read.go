package main

import (
	"errors"
	"fmt"
	"github.com/baude/hyperv_kvp/hyperv/ignition"
	"unsafe"

	"golang.org/x/sys/unix"
)

const TIMEOUT = 5000

const KVP_OP_REGISTER1 = 100
const HV_S_OK = 0

const HV_KVP_EXCHANGE_MAX_VALUE_SIZE = 2048
const HV_KVP_EXCHANGE_MAX_KEY_SIZE = 512

const KVP_OP_SET = 1

type hv_kvp_exchg_msg_value struct {
	value_type uint32
	key_size   uint32
	value_size uint32
	key        [HV_KVP_EXCHANGE_MAX_KEY_SIZE]uint8
	value      [HV_KVP_EXCHANGE_MAX_VALUE_SIZE]uint8
}

type hv_kvp_msg_set struct {
	data hv_kvp_exchg_msg_value
}

type hv_kvp_hdr struct {
	operation uint8
	pool      uint8
	pad       uint16
}

type hv_kvp_msg struct {
	kvp_hdr hv_kvp_hdr
	kvp_set hv_kvp_msg_set
	// unused is needed to get to the same struct size as the C version.
	unused [4856]byte
}

type hv_kvp_msg_ret struct {
	error   int
	kvp_set hv_kvp_msg_set
	// unused is needed to get to the same struct size as the C version.
	unused [4856]byte
}

func readKvpData() (map[string]string, error) {
	ret := make(map[string]string)

	kvp, err := unix.Open("/dev/vmbus/hv_kvp", unix.O_RDWR|unix.O_CLOEXEC, 0)
	if err != nil {
		return nil, err
	}
	defer unix.Close(kvp)

	var hv_msg hv_kvp_msg
	var hv_msg_ret hv_kvp_msg_ret

	hv_msg.kvp_hdr.operation = KVP_OP_REGISTER1

	const sizeof = int(unsafe.Sizeof(hv_msg))
	var asByteSlice []byte = (*(*[sizeof]byte)(unsafe.Pointer(&hv_msg)))[:]
	var retAsByteSlice []byte = (*(*[sizeof]byte)(unsafe.Pointer(&hv_msg_ret)))[:]

	l, err := unix.Write(kvp, asByteSlice)
	if err != nil {
		return nil, err
	}
	if l != int(sizeof) {
		return nil, fmt.Errorf("Failed to write to hv_kvp")
	}

next:
	for {
		var pfd unix.PollFd
		pfd.Fd = int32(kvp)
		pfd.Events = unix.POLLIN
		pfd.Revents = 0

		howMany, err := unix.Poll([]unix.PollFd{pfd}, TIMEOUT)
		if err != nil {
			if err == unix.EINVAL {
				return nil, err
			} else {
				continue
			}
		}

		if howMany == 0 {
			return ret, nil
		}

		l, err := unix.Read(kvp, asByteSlice)
		if l != sizeof {
			return nil, fmt.Errorf("Failed to read from hv_kvp")
		}

		switch hv_msg.kvp_hdr.operation {
		case KVP_OP_REGISTER1:
			continue next
		case KVP_OP_SET:
			// on the next two variables, we are cutting the last byte because otherwise
			// it is padded and key lookups fail
			key := []byte(hv_msg.kvp_set.data.key[:hv_msg.kvp_set.data.key_size-1])
			value := []byte(hv_msg.kvp_set.data.value[:hv_msg.kvp_set.data.value_size-1])
			ret[string(key)] = string(value)
		}

		hv_msg_ret.error = HV_S_OK

		l, err = unix.Write(kvp, retAsByteSlice)
		if err != nil {
			return nil, err
		}
		if l != int(sizeof) {
			return nil, fmt.Errorf("Failed to write to hv_kvp")
		}
	}
}

func main() {
	ret, err := readKvpData()
	if err != nil {
		panic(err)
	}
	var (
		counter int
		parts   ignition.Segments
		ign_key = "com_coreos_ignition_kvp_"
	)
	for {
		//fmt.Printf("Read %s -> %s\n", k, v)
		lookForKey := fmt.Sprintf("%s%d", ign_key, counter)
		val, exists := ret[lookForKey]
		if !exists {
			break
		}
		parts = append(parts, []byte(val))
	}
	if len(parts) < 1 {
		panic(errors.New("unable to find ignition configs in kvp"))
	}
	cfg := ignition.Glue(parts)
	fmt.Println(cfg)
}
