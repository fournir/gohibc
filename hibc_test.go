package gohibc

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    *HIBC
		wantErr bool
	}{
		{
			name: "no_code",
			args: args{
				s: "+",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "starts_w_secondary",
			args: args{
				s: "+$",
			},
			want: &HIBC{
				supplierFlag: '+',
			},
			wantErr: false,
		},
		{
			name: "invalid_supplier_code",
			args: args{
				s: "*A99912345/$$52001510X33",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "LIC_wrong",
			args: args{
				s: "+099912345/$$52001510X33",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "LIC_wrong",
			args: args{
				s: "+A-9912345/$$52001510X33",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "LIC_short",
			args: args{
				s: "+A99",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "LIC_ok",
			args: args{
				s: "+A99912345/$$52001510X33",
			},
			want: &HIBC{
				supplierFlag: '+',
				di: deviceIdentifier{
					lic: []rune("A999"),
					pcn: []rune("1234"),
					um:  '5',
				},
			},
			wantErr: false,
		},
		{
			name: "PCN_short",
			args: args{
				s: "+A99912",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "PCN_non_alphanumeric",
			args: args{
				s: "+A9991-*45/$$52001510X33",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "PCN_ok",
			args: args{
				s: "+A99912345C",
			},
			want: &HIBC{
				supplierFlag: '+',
				di: deviceIdentifier{
					lic: []rune("A999"),
					pcn: []rune("1234"),
					um:  '5',
				},
			},
			wantErr: false,
		},
		{
			name: "UM_non_numeric",
			args: args{
				s: "+A9991234NC",
			},
			want:    nil,
			wantErr: true,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}
