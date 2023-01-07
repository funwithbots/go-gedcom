package uuid

import (
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		wantErr bool
	}{
		{
			"basic",
			[]Option{WithVersion(1), WithVariant("dce")},
			false,
		},
		{
			"missing name error",
			[]Option{WithVersion(2), WithVariant("dce")},
			true,
		},
		{
			"has namespace",
			[]Option{WithVersion(4), WithVariant("dce"), WithNamespace("2163d569-2c70-43d4-bb87-ff9c58814ad")},
			false,
		},
		{
			"invalid namespace",
			[]Option{WithVersion(2), WithVariant("dce"), WithName("test"), WithNamespace("invalid")},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got != nil) == tt.wantErr {
				if tt.wantErr {
					t.Errorf("New() got uuid, expected error")
				} else {
					t.Errorf("New() got nil, expected error. Error is %v", err)
				}
			}
		})
	}
}

func TestUUID_Generate(t *testing.T) {
	type fields struct {
		Version   int
		Variant   string
		Namespace string
		Name      string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "simple",
			fields: fields{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, _ := New(WithVersion(tt.fields.Version), WithVariant(tt.fields.Variant))
			if got := u.Generate(); got == "" {
				t.Errorf("Generate() = %v, want %v", got, tt.want)
			}
		})
	}
}
