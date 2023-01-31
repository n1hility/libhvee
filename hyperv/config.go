package hyperv

import (
	"bytes"
	"errors"
	"fmt"
)

const (
	// Microsoft defines the key in the binary file to be up to 512 bytes.  If the
	// key is less than 512 bytes, it is padded with nulls to retain the 512 byte size.
	keyLen = 512
	// Microsoft defines the value in the binary file to be 2048 bytes.  Unused bytes will
	// be padded as described in the key description.
	valLen = 2048

	// KvpPool3 is where the kvp daemon writes the data it is
	// provided by hyperv for describing the guest vm.
	KvpPool3 = "/var/lib/hyperv/.kvp_pool_3"
)

// kvp represents a single key/value pairs
type kvp struct {
	// key is always 512 bytes
	key []byte
	// val is always 2048 bytes
	val []byte
}

// print is a debug only function rn
func (k kvp) print() {
	fmt.Printf("key: %s value: %s\n", k.key, k.val)
}

// newkvp creates a new instance of the key value pair in the
// predecided lengths
func newKvp(key []byte, val []byte) kvp {
	return kvp{
		key: bytes.Trim(key, "\x00"),
		val: bytes.Trim(val, "\x00"),
	}
}

// KvpFile represents a key value pair file and keeps information
// about the file and its records.
type KvpFile struct {
	// File being read
	*bytes.Reader
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
	for i, s := range f.kvps {
		keys[i] = fmt.Sprintf("%s", s.key)
	}
	return keys
}

// Get returns the value of a key.  The bool is returned
// to signify if the key was found.
func (f *KvpFile) Get(key string) (string, bool) {
	for _, s := range f.kvps {
		iterKey := fmt.Sprintf("%s", s.key)
		if key == iterKey {
			return fmt.Sprintf("%s", s.val), true
		}
	}
	return "", false
}

// parseNext is an internal function used by decode to read and record
// the next record in the file
func (f *KvpFile) parseNext() error {
	tmpKey := make([]byte, keyLen)
	n1, err := f.ReadAt(tmpKey, f.position)
	if err != nil {
		return err
	}
	if n1 != keyLen {
		return errors.New("unable to read full key length")
	}

	// remove cruft/padding
	f.position += keyLen

	tmpVal := make([]byte, valLen)

	n2, err := f.ReadAt(tmpVal, f.position)
	if err != nil {
		return err
	}
	if n2 != valLen {
		return errors.New("unable to read full value")
	}

	newkvp := newKvp(tmpKey, tmpVal)
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
func NewKVPFile(f *bytes.Reader) (KvpFile, error) {
	k := KvpFile{
		f,
		nil,
		0,
		// chunk out the file for easy way to iterate
		f.Size() / (keyLen + valLen),
	}
	err := k.decode()
	return k, err
}
