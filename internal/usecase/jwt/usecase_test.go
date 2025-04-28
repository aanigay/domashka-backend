package jwt

import (
	"domashka-backend/config"
	"reflect"
	"testing"
)

func TestUseCase_ValidateJWT(t *testing.T) {
	type fields struct {
		cfg *config.JWTConfig
	}
	type args struct {
		tokenString string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[string]interface{}
		wantErr bool
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UseCase{
				cfg: tt.fields.cfg,
			}
			got, err := u.ValidateJWT(tt.args.tokenString)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ValidateJWT() got = %v, want %v", got, tt.want)
			}
		})
	}
}
