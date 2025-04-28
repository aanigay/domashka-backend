package geo

import (
	"context"
	"testing"
)

func TestUseCase_PushClientAddress(t *testing.T) {
	type fields struct {
		repo GeoRepository
	}
	type args struct {
		ctx       context.Context
		addressID int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UseCase{
				repo: tt.fields.repo,
			}
			if err := u.PushClientAddress(tt.args.ctx, tt.args.addressID); (err != nil) != tt.wantErr {
				t.Errorf("PushClientAddress() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
