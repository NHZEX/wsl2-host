package hostsapi

import (
	"fmt"
	"reflect"
	"testing"
)

func TestFun(t *testing.T) {
	h, err := CreateAPI()

	fmt.Println(err)
	fmt.Println(h)
}

func TestCreateAPI(t *testing.T) {
	tests := []struct {
		name    string
		want    *HostsAPI
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateAPI()
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateAPI() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateAPI() got = %v, want %v", got, tt.want)
			}
		})
	}
}
