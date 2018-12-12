package ipnet

import (
	"encoding/binary"
	"net"
	"testing"
)

func TestIPv4_Int(t *testing.T) {
	type fields struct {
		IP        net.IP
		byteorder binary.ByteOrder
	}
	tests := []struct {
		name   string
		fields fields
		want   uint32
	}{
		{name: "big", fields: fields{IP: net.IPv4(1, 0, 0, 0), byteorder: binary.BigEndian}, want: 1 << 24},
		{name: "default big", fields: fields{IP: net.IPv4(1, 0, 0, 0)}, want: 1 << 24},
		{name: "default little", fields: fields{IP: net.IPv4(1, 0, 0, 0), byteorder: binary.LittleEndian}, want: 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &IPv4{
				IP:        tt.fields.IP,
				byteorder: tt.fields.byteorder,
			}
			if got := i.Int(); got != tt.want {
				t.Errorf("IPv4.Int() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIPv4_Equal(t *testing.T) {
	type fields struct {
		IP        net.IP
		byteorder binary.ByteOrder
	}
	type args struct {
		x *IPv4
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{name: "", fields: fields{
			IP: net.IPv4(1, 2, 3, 4)},
			args: args{
				&IPv4{
					IP: net.IPv4(1, 2, 3, 4),
				},
			},
			want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &IPv4{
				IP:        tt.fields.IP,
				byteorder: tt.fields.byteorder,
			}
			if got := i.Equal(tt.args.x); got != tt.want {
				t.Errorf("IPv4.Equal() = %v, want %v", got, tt.want)
			}
		})
	}
}
