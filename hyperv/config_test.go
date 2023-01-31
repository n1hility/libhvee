package hyperv

import (
	"bytes"
	"embed"
	"fmt"
	"reflect"
	"testing"
)

//go:embed test/default.data
var data embed.FS

var default_data_keys = []string{"HostName", "HostingSystemEditionId", "HostingSystemNestedLevel", "HostingSystemOsMajor", "HostingSystemOsMinor", "HostingSystemProcessorArchitecture", "HostingSystemProcessorIdleStateMax", "HostingSystemProcessorThrottleMax", "HostingSystemProcessorThrottleMin", "HostingSystemSpMajor", "HostingSystemSpMinor", "PhysicalHostName", "PhysicalHostNameFullyQualified", "VirtualMachineDynamicMemoryBalancingEnabled", "VirtualMachineId", "VirtualMachineName"}

func TestKvpFile_Get(t *testing.T) {
	f, _ := data.ReadFile("test/default.data")
	rdr := bytes.NewReader(f)
	o, _ := NewKVPFile(rdr)

	type fields struct {
		Reader   *bytes.Reader
		kvps     []kvp
		position int64
		records  int64
	}
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
		want1  bool
	}{
		{
			name: "valid key",
			fields: fields{
				Reader:   rdr,
				kvps:     o.kvps,
				position: 0,
				records:  int64(len(o.kvps)),
			},
			args: args{
				key: "VirtualMachineId",
			},
			want:  "8E3ECE81-7EF3-4581-A6B4-92F152ABE267",
			want1: true,
		},
		{
			name: "invalid key",
			fields: fields{
				Reader:   rdr,
				kvps:     o.kvps,
				position: 0,
				records:  int64(len(o.kvps)),
			},
			args: args{
				key: "VirtualMachineIdZZZZZZZZZZZ",
			},
			want:  "",
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &KvpFile{
				Reader:   tt.fields.Reader,
				kvps:     tt.fields.kvps,
				position: tt.fields.position,
				records:  tt.fields.records,
			}
			got, got1 := f.Get(tt.args.key)
			if got != tt.want {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Get() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestKvpFile_Keys(t *testing.T) {
	f, _ := data.ReadFile("test/default.data")
	rdr := bytes.NewReader(f)
	o, _ := NewKVPFile(rdr)
	type fields struct {
		Reader   *bytes.Reader
		kvps     []kvp
		position int64
		records  int64
	}
	var (
		tests = []struct {
			name   string
			fields fields
			want   []string
		}{
			{
				name: "default key equality",
				fields: fields{
					Reader:   rdr,
					kvps:     o.kvps,
					position: 0,
					records:  int64(len(o.kvps)),
				},
				want: default_data_keys,
			},
		}
	)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &KvpFile{
				Reader:   tt.fields.Reader,
				kvps:     tt.fields.kvps,
				position: tt.fields.position,
				records:  tt.fields.records,
			}
			if got := f.Keys(); !reflect.DeepEqual(got, tt.want) {
				fmt.Println(len(f.Keys()), len(tt.want))
				t.Errorf("Keys() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKvpFile_decode(t *testing.T) {
	type fields struct {
		Reader   *bytes.Reader
		kvps     []kvp
		position int64
		records  int64
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &KvpFile{
				Reader:   tt.fields.Reader,
				kvps:     tt.fields.kvps,
				position: tt.fields.position,
				records:  tt.fields.records,
			}
			if err := f.decode(); (err != nil) != tt.wantErr {
				t.Errorf("decode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_newKvp(t *testing.T) {
	type args struct {
		key []byte
		val []byte
	}
	tests := []struct {
		name string
		args args
		want kvp
	}{
		{
			name: "valid key and val",
			args: args{
				key: makeBytesWithPadding("foo", 10),
				val: makeBytesWithPadding("bar", 999),
			},
			want: kvp{
				key: []byte("foo"),
				val: []byte("bar"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newKvp(tt.args.key, tt.args.val); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newKvp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func makeBytesWithPadding(s string, l int) []byte {
	b := make([]byte, l)
	for i := 0; i < len(s); i++ {
		b[i] = s[i]
	}

	for i := len(s); i < l-len(s); i++ {
		b[i] = 0
	}

	return b
}
