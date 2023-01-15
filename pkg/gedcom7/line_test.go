package gedcom7

import (
	"reflect"
	"strings"
	"testing"
)

func TestLine_String(t *testing.T) {
	type fields struct {
		Level   int
		Xref    string
		Tag     string
		Payload string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "simple",
			fields: fields{
				Level:   0,
				Xref:    "@I1@",
				Tag:     "INDI",
				Payload: "John Doe",
			},
			want: "0 @I1@ INDI John Doe",
		},
		{
			name:   "empty",
			fields: fields{},
			want:   "",
		},
		{
			name: "no tag",
			fields: fields{
				Level:   0,
				Xref:    "@I1@",
				Payload: "John Doe",
			},
			want: "",
		},
		{
			name: "no xref",
			fields: fields{
				Level:   0,
				Tag:     "INDI",
				Payload: "John Doe",
			},
			want: "0 INDI John Doe",
		},
		{
			name: "no Payload",
			fields: fields{
				Level: 0,
				Xref:  "@I1@",
				Tag:   "INDI",
			},
			want: "0 @I1@ INDI",
		},
		{
			name: "tag only",
			fields: fields{
				Tag: "INDI",
			},
			want: "0 INDI",
		},
		{
			name: "neg level",
			fields: fields{
				Level:   -1,
				Xref:    "@I1@",
				Tag:     "INDI",
				Payload: "John Doe",
			},
			want: "",
		},
		{
			name: "extra spaces",
			fields: fields{
				Level:   0,
				Xref:    "@I1@",
				Tag:     "INDI",
				Payload: " John  Doe ",
			},
			want: "0 @I1@ INDI  John  Doe ",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := Line{
				Level:   tt.fields.Level,
				Xref:    tt.fields.Xref,
				Tag:     tt.fields.Tag,
				Payload: tt.fields.Payload,
			}
			if got := l.String(); got != tt.want {
				t.Errorf("Stringify() = '%v', want '%v'", got, tt.want)
			}
		})
	}
}

func TestToLine(t *testing.T) {
	tests := []struct {
		name    string
		s       string
		want    *Line
		wantErr bool
	}{
		{
			name: "simple",
			s:    "0 @I1@ INDI John Doe",
			want: &Line{
				Level:   0,
				Xref:    "@I1@",
				Tag:     "INDI",
				Payload: "John Doe",
			},
			wantErr: false,
		},
		{
			name: "underscore",
			s:    "0 @I1@ _CUSTOM John Doe",
			want: &Line{
				Level:   0,
				Xref:    "@I1@",
				Tag:     "_CUSTOM",
				Payload: "John Doe",
			},
			wantErr: false,
		},
		{
			name: "proper @@",
			s:    "0 NOTE @@me is John Doe",
			want: &Line{
				Level:   0,
				Tag:     "NOTE",
				Payload: "@@me is John Doe",
			},
			wantErr: false,
		},
		{
			name: "improper @@",
			s:    "0 NOTE @me is John Doe",
			want: &Line{
				Level:   0,
				Tag:     "NOTE",
				Payload: "@@me is John Doe",
			},
			wantErr: false,
		},
		{
			name: "no xref",
			s:    "0 INDI John Doe",
			want: &Line{
				Level:   0,
				Tag:     "INDI",
				Payload: "John Doe",
			},
			wantErr: false,
		},
		{
			name: "extra spaces",
			s:    "0 @I1@ INDI  John  Doe ",
			want: &Line{
				Level:   0,
				Xref:    "@I1@",
				Tag:     "INDI",
				Payload: " John  Doe ",
			},
			wantErr: false,
		},
		{
			name: "no Payload",
			s:    "0 @I1@ INDI",
			want: &Line{
				Level: 0,
				Xref:  "@I1@",
				Tag:   "INDI",
			},
			wantErr: false,
		},
		{
			name: "most basic",
			s:    "1 DATA",
			want: &Line{
				Level: 1,
				Tag:   "DATA",
			},
			wantErr: false,
		},
		{
			name:    "empty",
			s:       "",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "missing tag",
			s:       "0 @I1@ John Doe",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "dbl underscore",
			s:       "0 __CUSTOM John Doe",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "void",
			s:       "0 @VOID@ INDI John Doe",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "non-ascii tag",
			s:       "0 NÖBIT John Doe",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "banned char",
			s:       "1 TEXT Is \u007f VALID?",
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.want != nil {
				tt.want.Text = tt.s
			}
			got, err := ToLine(tt.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToLine() error = %v, wantErr %v\n%+v", err, tt.wantErr, got)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToLine() got = '%+v', want '%+v'", got, tt.want)
				t.Logf("%d\t%d\n'%s'\t'%s'\n'%s'\t'%s'\n'%s'\t'%s'\n'%s'\t'%s'\n%t\t%t\n\n",
					got.Level, tt.want.Level,
					got.Xref, tt.want.Xref,
					got.Tag, tt.want.Tag,
					got.Payload, tt.want.Payload,
					got.Text, tt.want.Text,
					got.Deleted, tt.want.Deleted,
				)
			}
		})
	}
}

func TestNewlineFromFile(t *testing.T) {
	tests := []struct {
		name      string
		file      string
		total     int
		rootCount int
	}{
		{
			name:      "tgc551lf",
			file:      DocPath + "torture test/TGC551LF.ged",
			total:     2161,
			rootCount: 65,
		},
		{
			name:      "tgc551",
			file:      DocPath + "torture test/TGC551.ged",
			rootCount: 65,
			total:     2161,
		},
		{
			name:      "tgc55c",
			file:      DocPath + "torture test/TGC55C.ged",
			total:     2197,
			rootCount: 67,
		},
		{
			name:      "tgc55clf",
			file:      DocPath + "torture test/TGC55CLF.ged",
			total:     2197,
			rootCount: 67,
		},
		{
			name:      "Family Historian",
			file:      DocPath + "test/valid.ged",
			total:     93,
			rootCount: 10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lines, err := NewLinesFromFile(tt.file, WithMaxDeprecatedTags("5.5.1"))
			if err != nil {
				t.Errorf("NewLines() error opening %s: %v", tt.file, err)
			}

			if len(*lines) != tt.total {
				t.Errorf("NewLines() = wrong total number of lines. Wanted %d; got %d", tt.total, len(*lines))
			}

			var count int
			for _, v := range *lines {

				if strings.HasPrefix(v.Text, "0 ") {
					count++
				}
			}
			if count != tt.rootCount {
				t.Errorf("NewLines() = wrong number of root nodes. Wanted %d; got %d", tt.rootCount, count)
			}
		})
	}
}
