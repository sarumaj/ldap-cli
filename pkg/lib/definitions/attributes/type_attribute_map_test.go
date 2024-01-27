package attributes

import (
	"net"
	"reflect"
	"testing"
	"time"
)

func TestKeys(t *testing.T) {
	for _, tt := range []struct {
		name string
		args Map
		want Attributes
	}{
		{"test#2", Map{name: "", accountExpires: 0}, Attributes{accountExpires, name}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.args.Keys()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf(`(Map).Keys() failed: got: %v, want: %v`, got, tt.want)
			}
		})
	}
}

func TestParseBool(t *testing.T) {
	for _, tt := range []struct {
		name string
		args []string
		want any
	}{
		{"test#1", []string{"true"}, true},
		{"test#2", []string{"invalid"}, []string{"invalid"}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			v, a := make(Map), Attribute{LDAPDisplayName: "test"}
			v.ParseBool(a, tt.args)
			if got := v[a]; !reflect.DeepEqual(got, tt.want) {
				t.Errorf(`(*Map).ParseBool("test", %[1]v) failed: got: %[2]v (%[2]T), want: %[3]v (%[3]T)`, tt.args, got, tt.want)
			}
		})
	}
}

func TestParseDecimal(t *testing.T) {
	for _, tt := range []struct {
		name string
		args []string
		want any
	}{
		{"test#1", []string{"12345678.9"}, 12345678.9},
		{"test#2", []string{"invalid"}, []string{"invalid"}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			v, a := make(Map), Attribute{LDAPDisplayName: "test"}
			v.ParseDecimal(a, tt.args)
			if got := v[a]; !reflect.DeepEqual(got, tt.want) {
				t.Errorf(`(*Map).ParseDecimal("test", %[1]v) failed: got: %[2]v (%[2]T), want: %[3]v (%[3]T)`, tt.args, got, tt.want)
			}
		})
	}
}

func TestParseGroupType(t *testing.T) {
	for _, tt := range []struct {
		name string
		args []string
		want any
	}{
		{"test#1", []string{"10"}, []string{"DISTRIBUTION", "GLOBAL", "UNIVERSAL"}},
		{"test#2", []string{"invalid"}, []string{"invalid"}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			v, a := make(Map), Attribute{LDAPDisplayName: "test"}
			v.ParseGroupType(a, tt.args)
			if got := v[a]; !reflect.DeepEqual(got, tt.want) {
				t.Errorf(`(*Map).ParseGroupType("test", %[1]v) failed: got: %[2]v (%[2]T), want: %[3]v (%[3]T)`, tt.args, got, tt.want)
			}
		})
	}
}

func TestParseInt(t *testing.T) {
	for _, tt := range []struct {
		name string
		args []string
		want any
	}{
		{"test#1", []string{"10"}, int64(10)},
		{"test#2", []string{"invalid"}, []string{"invalid"}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			v, a := make(Map), Attribute{LDAPDisplayName: "test"}
			v.ParseInt(a, tt.args)
			if got := v[a]; !reflect.DeepEqual(got, tt.want) {
				t.Errorf(`(*Map).ParseInt("test", %[1]v) failed: got: %[2]v (%[2]T), want: %[3]v (%[3]T)`, tt.args, got, tt.want)
			}
		})
	}
}

func TestParseIPv4Address(t *testing.T) {
	for _, tt := range []struct {
		name string
		args []string
		want any
	}{
		{"test#1", []string{"2130706433"}, net.IP{127, 0, 0, 1}},
		{"test#2", []string{"invalid"}, []string{"invalid"}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			v, a := make(Map), Attribute{LDAPDisplayName: "test"}
			v.ParseIPv4Address(a, tt.args)
			if got := v[a]; !reflect.DeepEqual(got, tt.want) {
				t.Errorf(`(*Map).ParseIPv4Address("test", %[1]v) failed: got: %[2]v (%[2]T), want: %[3]v (%[3]T)`, tt.args, got, tt.want)
			}
		})
	}
}

func TestParseTime(t *testing.T) {
	for _, tt := range []struct {
		name string
		args []string
		want any
	}{
		{"test#1", []string{"72057594037927936"}, time.Date(1829, 5, 5, 23, 50, 3, 792799000, time.UTC)},
		{"test#2", []string{"invalid"}, []string{"invalid"}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			v, a := make(Map), Attribute{LDAPDisplayName: "test"}
			v.ParseTime(a, tt.args)
			if got := v[a]; !reflect.DeepEqual(got, tt.want) {
				t.Errorf(`(*Map).ParseTime("test", %[1]v) failed: got: %[2]v (%[2]T), want: %[3]v (%[3]T)`, tt.args, got, tt.want)
			}
		})
	}
}

func TestParseSAMAccountType(t *testing.T) {
	for _, tt := range []struct {
		name string
		args []string
		want any
	}{
		{"test#1", []string{"805306368"}, []string{"NORMAL_USER_ACCOUNT", "USER_OBJECT"}},
		{"test#2", []string{"invalid"}, []string{"invalid"}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			v, a := make(Map), Attribute{LDAPDisplayName: "test"}
			v.ParseSAMAccountType(a, tt.args)
			if got := v[a]; !reflect.DeepEqual(got, tt.want) {
				t.Errorf(`(*Map).ParseSAMAccountType("test", %[1]v) failed: got: %[2]v (%[2]T), want: %[3]v (%[3]T)`, tt.args, got, tt.want)
			}
		})
	}
}

func TestParseUserAccountControl(t *testing.T) {
	for _, tt := range []struct {
		name string
		args []string
		want any
	}{
		{"test#1", []string{"514"}, []string{"ACCOUNT_DISABLE", "NORMAL_ACCOUNT"}},
		{"test#2", []string{"invalid"}, []string{"invalid"}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			v, a := make(Map), Attribute{LDAPDisplayName: "test"}
			v.ParseUserAccountControl(a, tt.args)
			if got := v[a]; !reflect.DeepEqual(got, tt.want) {
				t.Errorf(`(*Map).ParseUserAccountControl("test", %[1]v) failed: got: %[2]v (%[2]T), want: %[3]v (%[3]T)`, tt.args, got, tt.want)
			}
		})
	}
}
