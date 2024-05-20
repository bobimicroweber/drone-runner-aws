package hetzner

import (
	"fmt"
	"os"

	"github.com/drone-runners/drone-runner-aws/internal/oshelp"
	"github.com/drone-runners/drone-runner-aws/types"
	"github.com/sirupsen/logrus"
)

type Option func(*config)

func SetPlatformDefaults(platform *types.Platform) (*types.Platform, error) {
	if platform.Arch == "" {
		platform.Arch = oshelp.ArchAMD64
	}
	if platform.Arch != oshelp.ArchAMD64 {
		return platform, fmt.Errorf("invalid arch %s, has to be '%s'", platform.Arch, oshelp.ArchAMD64)
	}
	// verify that we are using sane values for OS
	if platform.OS == "" {
		platform.OS = oshelp.OSLinux
	}
	if platform.OS != oshelp.OSLinux {
		return platform, fmt.Errorf("hetzner - invalid OS %s, has to be '%s'", platform.OS, oshelp.OSLinux)
	}
	// set osname
	if platform.OS == oshelp.OSLinux {
		if platform.OSName == "" {
			platform.OSName = oshelp.Ubuntu
		}
	}
	return platform, nil
}

func WithToken(token string) Option {
	return func(p *config) {
		p.token = token
	}
}

func WithRegion(region string) Option {
	return func(p *config) {
		if region == "" {
			p.region = "nyc1"
		} else {
			p.region = region
		}
	}
}

func WithSize(size string) Option {
	return func(p *config) {
		if size == "" {
			p.size = "s-2vcpu-4gb"
		} else {
			p.size = size
		}
	}
}

func WithImage(image string) Option {
	return func(p *config) {
		if image == "" {
			p.image = "docker-18-04"
		} else {
			p.image = image
		}
	}
}

func WithFirewallID(firewallID string) Option {
	return func(p *config) {
		p.FirewallID = firewallID
	}
}

func WithTags(tags []string) Option {
	return func(p *config) {
		p.tags = tags
	}
}

func WithSSHKeys(sshKeys []string) Option {
	return func(p *config) {
		p.SSHKeys = sshKeys
	}
}

func WithUserData(text, path string) Option {
	if text != "" {
		return func(p *config) {
			p.userData = text
		}
	}
	return func(p *config) {
		if path != "" {
			data, err := os.ReadFile(path)
			if err != nil {
				logrus.WithError(err).
					Fatalln("failed to read user_data file")
				return
			}
			p.userData = string(data)
		}
	}
}

// WithRootDirectory sets the root directory for the virtual machine.
func WithRootDirectory(dir string) Option {
	return func(p *config) {
		p.rootDir = oshelp.JoinPaths(oshelp.OSLinux, "/tmp", "hetzner")
	}
}
