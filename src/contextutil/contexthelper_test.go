package contextutil

import (
	"context"
	"testing"
)

func TestCheckCtx(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"not_done", args{context.Background()}, false},
		{"done", args{ctx}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CheckCtx(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("CheckCtx() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
