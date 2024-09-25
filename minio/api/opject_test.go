package api

import "testing"

func TestPutObject(t *testing.T) {
	type args struct {
		bucketName string
		objectName string
		filePath   string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "ssss", args: args{bucketName: "cube-01", objectName: "test-png", filePath: "D:\\downloads\\图片\\3bc064a3-9cc2-4e4c-be23-d9c61feea139@1x.png"}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := PutObject(tt.args.bucketName, tt.args.objectName, tt.args.filePath); (err != nil) != tt.wantErr {
				t.Errorf("PutObject() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
