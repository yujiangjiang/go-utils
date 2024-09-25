package jwt

import (
	"fmt"
	"testing"
)

func Test_authenticateNode(t *testing.T) {
	type args struct {
		accessKey string
		secretKey string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{name: "sss", args: args{accessKey: "minio", secretKey: "minio123"}, want: "test", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := authenticateNode(tt.args.accessKey, tt.args.secretKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("authenticateNode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			fmt.Print(got)
			if got != tt.want {
				t.Errorf("authenticateNode() got = %v, want %v", got, tt.want)
			}
		})
	}
}
