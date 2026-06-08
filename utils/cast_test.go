package utils

import (
	"testing"
)

// ===================== int =====================

func TestConvert_ToInt(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		want    int
		wantErr bool
	}{
		{"int -> int", 10, 10, false},
		{"int64 -> int", int64(20), 20, false},
		{"float64 -> int", 12.8, 12, false},
		{"string -> int", "123", 123, false},
		{"string invalid", "abc", 0, true},
		{"bool invalid", true, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Convert[int](tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("got = %v, want %v", got, tt.want)
			}
		})
	}
}

// ===================== float64 =====================

func TestConvert_ToFloat64(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		want    float64
		wantErr bool
	}{
		{"float64 -> float64", 1.23, 1.23, false},
		{"int -> float64", 10, 10.0, false},
		{"string -> float64", "3.14", 3.14, false},
		{"string invalid", "abc", 0, true},
		{"bool invalid", true, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Convert[float64](tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("got = %v, want %v", got, tt.want)
			}
		})
	}
}

// ===================== string =====================

func TestConvert_ToString(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		want    string
		wantErr bool
	}{
		{"string -> string", "abc", "abc", false},
		{"int -> string", 123, "123", false},
		{"float -> string", 1.5, "1.500000", false},
		{"bool -> string", true, "true", false},
		{"unsupported", struct{}{}, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Convert[string](tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("got = %v, want %v", got, tt.want)
			}
		})
	}
}

// ===================== bool =====================

func TestConvert_ToBool(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		want    bool
		wantErr bool
	}{
		{"bool -> bool", true, true, false},
		{"string true", "true", true, false},
		{"string false", "false", false, false},
		{"string invalid", "abc", false, true},
		{"int invalid", 1, false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Convert[bool](tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("got = %v, want %v", got, tt.want)
			}
		})
	}
}

// ===================== 不支持类型 =====================

func TestConvert_UnsupportedType(t *testing.T) {
	type MyStruct struct{}
	//_, err := Convert
	//if err == nil {
	//	t.Errorf("expected error for unsupported type")
	//}
}
