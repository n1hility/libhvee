package hyperv

import (
	"errors"
	"fmt"
	"os"
)

const (
	// Microsoft defines the key in the binary file to be up to 512 bytes.  If the
	// key is less than 512 bytes, it is padded with nulls to retain the 512 byte size.
	keyLen = 512
	// Microsoft defines the value in the binary file to be 2048 bytes.  Unused bytes will
	// be padded as described in the key description.
	valLen = 2048
)

// kvp represents a single key/value pairs
type kvp struct {
	// key is always 512 bytes
	key []byte
	// val is always 2048 bytes
	val []byte
}

func (k kvp) print() {
	fmt.Printf("key: %s value: %s\n", k.key, k.val)
}

// newkvp creates a new instance of the key value pair in the
// predecided lengths
func newKvp() kvp {
	return kvp{
		key: make([]byte, 512),
		val: make([]byte, 2048),
	}
}

// KvpFile represents a key value pair file and keeps information
// about the file and its records.
type KvpFile struct {
	// File being read
	*os.File
	// once decoded, kvps are an array of key/value pairs
	kvps []kvp
	// read position, byte number
	position int64
	// number of key/value pairs found
	records int64
}

// Keys returns all the keys in the key-value file
func (f *KvpFile) Keys() []string {
	keys := make([]string, len(f.kvps))
	for _, s := range f.kvps {
		keys = append(keys, fmt.Sprintf("%s", s.key))
	}
	return keys
}

// Get returns the value of a key.  The bool is returned
// to signify if the key was found.
func (f *KvpFile) Get(key string) (string, bool) {
	for _, s := range f.kvps {
		iterKey := fmt.Sprintf("%s", s.key)
		// we need to trim the kvp key string because it is always
		// 512 bytes long; we trim to len of the key being searched for and
		// then comparison is made
		if key == iterKey[0:len(key)] {
			return iterKey, true
		}
	}
	return "", false
}

// parseNext is an internal function used by decode to read and record
// the next record in the file
func (f *KvpFile) parseNext() error {
	newkvp := newKvp()
	n1, err := f.ReadAt(newkvp.key, f.position)
	if err != nil {
		return err
	}
	if n1 != keyLen {
		return errors.New("unable to read full key length")
	}
	f.position += keyLen
	n2, err := f.ReadAt(newkvp.val, f.position)
	if err != nil {
		return err
	}
	if n2 != valLen {
		return errors.New("unable to read full value")
	}
	f.kvps = append(f.kvps, newkvp)
	f.position += valLen
	return nil
}

// decode is an internal function that iterates the binary file
// and creates records from the key/value pairs
func (f *KvpFile) decode() error {
	if f.kvps == nil {
		f.kvps = make([]kvp, 0)
	}
	for i := int64(0); i < f.records; i++ {
		if err := f.parseNext(); err != nil {
			return err
		}
	}
	return nil
}

// NewKVPFile reads and processes a key-value-pair file exposed by hyperv to Linux
// in a "binary" form.
func NewKVPFile(f *os.File) (KvpFile, error) {
	fileInfo, err := f.Stat()
	if err != nil {
		return KvpFile{}, err
	}
	k := KvpFile{
		File:     f,
		kvps:     nil,
		position: 0,
		// chunk out the file for easy way to iterate
		records: fileInfo.Size() / (keyLen + valLen),
	}
	err = k.decode()
	return k, err
}
