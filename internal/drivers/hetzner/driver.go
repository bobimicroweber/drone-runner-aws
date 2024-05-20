package hetzner

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/drone-runners/drone-runner-aws/internal/drivers"
	"github.com/drone-runners/drone-runner-aws/internal/lehelper"
	"github.com/drone-runners/drone-runner-aws/types"
	"github.com/drone/runner-go/logger"

	"github.com/dchest/uniuri"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

// config is a struct that implements drivers.Pool interface
type config struct {
	token      string
	region     string
	size       string
	tags       []string
	FirewallID string
	SSHKeys    []string
	userData   string
	rootDir    string

	image string

	hibernate bool
}

func New(opts ...Option) (drivers.Driver, error) {
	p := new(config)
	for _, opt := range opts {
		opt(p)
	}
	return p, nil
}

func (p *config) DriverName() string {
	return string(types.Hetzner)
}

func (p *config) InstanceType() string {
	return p.image
}

func (p *config) RootDir() string {
	return p.rootDir
}

func (p *config) CanHibernate() bool {
	return p.hibernate
}

func (p *config) Ping(ctx context.Context) error {
	client := newClient(ctx, p.token)

	_, err := client.Server.All(context.Background())

	return err
}

// Create an AWS instance for the pool, it will not perform build specific setup.
func (p *config) Create(ctx context.Context, opts *types.InstanceCreateOpts) (instance *types.Instance, err error) {
	startTime := time.Now()
	logr := logger.FromContext(ctx).
		WithField("driver", types.Hetzner).
		WithField("pool", opts.PoolName).
		WithField("image", p.image).
		WithField("hibernate", p.CanHibernate())
	var name = fmt.Sprintf("%s-%s-%s", opts.RunnerName, opts.PoolName, uniuri.NewLen(8)) //nolint:gomnd
	logr.Infof("hetzner: creating instance %s", name)

	// create the instance
	client := newClient(ctx, p.token)

	req := hcloud.ServerCreateOpts{
		Name:       "my-ubuntu-server",
		Image:      &hcloud.Image{Name: "ubuntu-22.04"},
		ServerType: &hcloud.ServerType{Name: "cx11"},
		Location:   &hcloud.Location{Name: "nbg1"},
	}

	// set the ssh keys if they are provided
	if len(p.SSHKeys) > 0 {
		req.SSHKeys = createSSHKeys(p.SSHKeys)
	}

	createServer, _, err := client.Server.Create(context.Background(), req)
	if err != nil {
		logr.WithError(err).
			Errorln("cannot create instance")
		return nil, err
	}
	logr.Infof("hetzner: instance created %s", name)
	// get firewall id
	// 	if p.FirewallID == "" {
	// 		id, getFirewallErr := getFirewallID(ctx, client, len(p.SSHKeys) > 0)
	// 		if getFirewallErr != nil {
	// 			logr.WithError(getFirewallErr).
	// 				Errorln("cannot get firewall id")
	// 			return nil, getFirewallErr
	// 		}
	// 		p.FirewallID = id
	// 	}
	// 	// setup the firewall
	// 	_, firewallErr := client.Firewalls.AddDroplets(ctx, p.FirewallID, server.ID)
	// 	if firewallErr != nil {
	// 		logr.WithError(firewallErr).
	// 			Errorln("cannot assign instance to firewall")
	// 		return nil, firewallErr
	// 	}
	// 	logr.Infof("hetzner: firewall configured %s", name)
	// initialize the instance
	instance = &types.Instance{
		Name:         name,
		Provider:     types.Hetzner,
		State:        types.StateCreated,
		Pool:         opts.PoolName,
		Region:       p.region,
		Image:        p.image,
		Size:         p.size,
		Platform:     opts.Platform,
		CAKey:        opts.CAKey,
		CACert:       opts.CACert,
		TLSKey:       opts.TLSKey,
		TLSCert:      opts.TLSCert,
		Started:      startTime.Unix(),
		Updated:      startTime.Unix(),
		IsHibernated: false,
		Port:         lehelper.LiteEnginePort,
	}
	// poll the digitalocean endpoint for server updates and exit when a network address is allocated.
	interval := time.Duration(0)
poller:
	for {
		select {
		case <-ctx.Done():
			logr.WithField("name", instance.Name).
				Debugln("cannot ascertain network")

			return instance, ctx.Err()
		case <-time.After(interval):
			interval = time.Minute

			logr.WithField("name", instance.Name).
				Debugln("find instance network")

			server, _, err := client.Server.GetByID(context.Background(), createServer.Server.ID)
			if err != nil {
				logr.WithError(err).
					Errorln("cannot find instance")
				return instance, err
			}
			// 			instance.ID = fmt.Sprint(createServer.Server.ID)
			// 			for _, network := range server.Networks.V4 {
			// 				if network.Type == "public" {
			// 					instance.Address = network.IPAddress
			// 				}
			// 			}
			//
			// 			if instance.Address != "" {
			// 				break poller
			// 			}
		}
	}

	return instance, err
}

// Destroy destroys the server AWS EC2 instances.
func (p *config) Destroy(ctx context.Context, instances []*types.Instance) (err error) {
	var instanceIDs []string
	for _, instance := range instances {
		instanceIDs = append(instanceIDs, instance.ID)
	}
	if len(instanceIDs) == 0 {
		return fmt.Errorf("no instance ids provided")
	}

	logr := logger.FromContext(ctx).
		WithField("id", instanceIDs).
		WithField("driver", types.DigitalOcean)

	client := newClient(ctx, p.token)
	for _, instanceID := range instanceIDs {
		id, err := strconv.Atoi(instanceID)
		if err != nil {
			return err
		}

		_, res, err := client.Server.Get(ctx, id)
		if err != nil && res.StatusCode == 404 {
			logr.WithError(err).
				Warnln("droplet does not exist")
			return fmt.Errorf("droplet does not exist '%s'", err)
		} else if err != nil {
			logr.WithError(err).
				Errorln("cannot find droplet")
			return err
		}
		logr.Debugln("deleting droplet")

		_, err = client.Droplets.Delete(ctx, id)
		if err != nil {
			logr.WithError(err).
				Errorln("deleting droplet failed")
			return err
		}
		logr.Debugln("droplet deleted")
	}
	logr.Traceln("digitalocean: VM terminated")
	return
}

func (p *config) Logs(ctx context.Context, instanceID string) (string, error) {
	return "no logs here", nil
}

func (p *config) Hibernate(ctx context.Context, instanceID, poolName string) error {
	return nil
}

func (p *config) Start(ctx context.Context, instanceID, poolName string) (string, error) {
	return "", nil
}

func (p *config) SetTags(ctx context.Context, instance *types.Instance,
	tags map[string]string) error {
	return nil
}

// helper function returns a new digitalocean client.
func newClient(ctx context.Context, token string) *hcloud.Client {

	return hcloud.NewClient(hcloud.WithToken(token))

}

// take a slice of ssh keys and return a slice of hcloud.SSHKey.Create
func createSSHKeys(sshKeys []string) {

	var keys []string
	for _, key := range sshKeys {
		opts := hcloud.SSHKeyCreateOpts{
			Name:      "drone-runner-ssh-key",
			PublicKey: key,
		}
		clinet := hcloud.SSHKeyClient()
		sshKey, _, err := clinet.create(context.Background(), opts)
		if err != nil {
			panic(err)
		}
	}

	return keys
}

// retrieve the runner firewall id or create a new one.
func getFirewallID(ctx context.Context, client *hcloud.Client, sshException bool) (string, error) {
	firewalls, _, listErr := client.Firewalls.List(ctx, &hcloud.ListOptions{})
	if listErr != nil {
		return "", listErr
	}
	// if the firewall already exists, return the id. NB we do not update any new firewall rules.
	for i := range firewalls {
		if firewalls[i].Name == "harness-runner" {
			return firewalls[i].ID, nil
		}
	}

	inboundRules := []hcloud.InboundRule{
		{
			Protocol:  "tcp",
			PortRange: "9079",
			Sources: &hcloud.Sources{
				Addresses: []string{"0.0.0.0/0", "::/0"},
			},
		},
	}
	if sshException {
		inboundRules = append(inboundRules, hcloud.InboundRule{
			Protocol:  "tcp",
			PortRange: "22",
			Sources: &hcloud.Sources{
				Addresses: []string{"0.0.0.0/0", "::/0"},
			},
		})
	}
	// firewall does not exist, create one.
	firewall, _, createErr := client.Firewalls.Create(ctx, &hcloud.FirewallRequest{
		Name:         "harness-runner",
		InboundRules: inboundRules,
		OutboundRules: []hcloud.OutboundRule{
			{
				Protocol:  "icmp",
				PortRange: "0",
				Destinations: &hcloud.Destinations{
					Addresses: []string{"0.0.0.0/0", "::/0"},
				},
			},
			{
				Protocol:  "tcp",
				PortRange: "0",
				Destinations: &hcloud.Destinations{
					Addresses: []string{"0.0.0.0/0", "::/0"},
				},
			},
			{
				Protocol:  "udp",
				PortRange: "0",
				Destinations: &hcloud.Destinations{
					Addresses: []string{"0.0.0.0/0", "::/0"},
				},
			},
		},
	})

	if createErr != nil {
		return "", createErr
	}
	return firewall.ID, nil
}
