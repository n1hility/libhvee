package ignition

import (
	"bytes"
)

const (
	kvpValueMaxLen = int(990)
)

type segment []byte
type Segments []segment

func Dice(k *bytes.Reader) (Segments, error) {
	var (
		done  bool
		parts Segments
	)
	for {
		sl := make([]byte, kvpValueMaxLen)
		n, err := k.Read(sl)
		if err != nil {
			return nil, err
		}
		if n < kvpValueMaxLen {
			sl = sl[0:n]
			done = true
		}
		parts = append(parts, sl)
		if done {
			break
		}
	}
	return parts, nil
}

func Glue(parts Segments) (b []byte) {
	for _, p := range parts {
		b = append(b, p...)
	}
	return
}
