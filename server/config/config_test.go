package config

import (
	"encoding/json"
	"reflect"
	"testing"
)

type exampleConfig struct {
	IntValue    int    `mapstructure:"int_value"`
	StringValue string `mapstructure:"string_value"`
}

func TestMustJsonRawMessage(t *testing.T) {
	type args struct {
		v any
	}
	tests := []struct {
		name string
		args args
		want json.RawMessage
	}{
		{
			name: "simple struct",
			args: args{
				v: exampleConfig{
					IntValue:    1,
					StringValue: "test",
				},
			},
			want: json.RawMessage(`{"int_value":1,"string_value":"test"}`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MustJsonRawMessage(tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MustJsonRawMessage() = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}
