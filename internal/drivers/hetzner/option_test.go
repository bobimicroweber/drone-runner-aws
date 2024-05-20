package hetzner

import (
	"reflect"
	"testing"

	"github.com/drone-runners/drone-runner-aws/internal/oshelp"
	"github.com/drone-runners/drone-runner-aws/types"
)

func TestSetPlatformDefaults(t *testing.T) {
	tests := []struct {
		name     string
		platform *types.Platform
		want     *types.Platform
		wantErr  bool
	}{
		{
			name:     "happy path no defaults",
			platform: &types.Platform{},
			want: &types.Platform{
				Arch:   oshelp.ArchAMD64,
				OS:     oshelp.OSLinux,
				OSName: oshelp.Ubuntu,
			},
			wantErr: false,
		},
		{
			name: "happy path no defaults",
			platform: &types.Platform{
				Arch:   oshelp.ArchAMD64,
				OS:     oshelp.OSLinux,
				OSName: oshelp.Ubuntu,
			},
			want: &types.Platform{
				Arch:   oshelp.ArchAMD64,
				OS:     oshelp.OSLinux,
				OSName: oshelp.Ubuntu,
			},
			wantErr: false,
		},
		{
			name: "err on bad OS",
			platform: &types.Platform{
				Arch:   oshelp.ArchAMD64,
				OS:     oshelp.OSWindows,
				OSName: oshelp.Ubuntu,
			},
			want: &types.Platform{
				Arch:   oshelp.ArchAMD64,
				OS:     oshelp.OSWindows,
				OSName: oshelp.Ubuntu,
			},
			wantErr: true,
		},
		{
			name: "err on bad arch",
			platform: &types.Platform{
				Arch:   "bad",
				OS:     oshelp.OSLinux,
				OSName: oshelp.Ubuntu,
			},
			want: &types.Platform{
				Arch:   "bad",
				OS:     oshelp.OSLinux,
				OSName: oshelp.Ubuntu,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SetPlatformDefaults(tt.platform)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetPlatformDefaults() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetPlatformDefaults() = %v, want %v", got, tt.want)
			}
		})
	}
}
