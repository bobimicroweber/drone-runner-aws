package types

import (
	"database/sql/driver"
)

type InstanceState string
type DriverType string

type ContextKey string

const (
	Hosted = ContextKey("hosted")
)

// Value converts the value to a sql string.
func (s InstanceState) Value() (driver.Value, error) {
	return string(s), nil
}

func (s DriverType) Value() (driver.Value, error) {
	return string(s), nil
}

const (
	Amazon       = DriverType("amazon")
	Anka         = DriverType("anka")
	AnkaBuild    = DriverType("ankabuild")
	Azure        = DriverType("azure")
	DigitalOcean = DriverType("digitalocean")
	Hetzner      = DriverType("hetzner")
	Google       = DriverType("google")
	VMFusion     = DriverType("vmfusion")
	Noop         = DriverType("noop")
	Nomad        = DriverType("nomad")
)

// InstanceState type enumeration.
const (
	StateCreated     = InstanceState("created")
	StateInUse       = InstanceState("inuse")
	StateHibernating = InstanceState("hibernating")
)

type Instance struct {
	ID           string        `db:"instance_id" json:"id"`
	NodeID       string        `db:"instance_node_id" json:"node_id"`
	Name         string        `db:"instance_name" json:"name"`
	Address      string        `db:"instance_address" json:"address"`
	Provider     DriverType    `db:"instance_provider" json:"provider"` // this is driver, though its the old legacy name of provider
	State        InstanceState `db:"instance_state" json:"state"`
	Pool         string        `db:"instance_pool" json:"pool"`
	Image        string        `db:"instance_image" json:"image"`
	Region       string        `db:"instance_region" json:"region"`
	Zone         string        `db:"instance_zone" json:"zone"`
	Size         string        `db:"instance_size" json:"size"`
	OwnerID      string        `db:"instance_owner_id" json:"owner_id"`
	Platform     `json:"platform"`
	CAKey        []byte `db:"instance_ca_key" json:"ca_key"`
	CACert       []byte `db:"instance_ca_cert" json:"ca_cert"`
	TLSKey       []byte `db:"instance_tls_key" json:"tls_key"`
	TLSCert      []byte `db:"instance_tls_cert" json:"tls_cert"`
	Stage        string `db:"instance_stage" json:"stage"`
	Updated      int64  `db:"instance_updated" json:"updated"`
	Started      int64  `db:"instance_started" json:"started"`
	IsHibernated bool   `db:"is_hibernated" json:"is_hibernated"`
	Port         int64  `db:"instance_port" json:"port"`
	RunnerName   string `db:"runner_name" json:"runner_name"`
}

type Tmate struct {
	Enabled bool
	Image   string
	Server  string
	Port    string
	RSA     string
	ED25519 string
}

type InstanceCreateOpts struct {
	CAKey          []byte
	CACert         []byte
	TLSKey         []byte
	TLSCert        []byte
	LiteEnginePath string
	Platform
	PoolName             string
	RunnerName           string
	Limit                int
	Pool                 int
	HarnessTestBinaryURI string
	PluginBinaryURI      string
	Tmate                Tmate
	AccountID            string
	IsHosted             bool
	ResourceClass        string
}

// Platform defines the target platform.
type Platform struct {
	OS      string `json:"os,omitempty" db:"instance_os" default:"linux"`
	Arch    string `json:"arch,omitempty" db:"instance_arch" default:"amd64"`
	Variant string `json:"variant,omitempty" yaml:"variant,omitempty" db:"instance_variant"`
	Version string `json:"version,omitempty" yaml:"version,omitempty" db:"instance_version"`
	OSName  string `json:"os_name,omitempty" yaml:"os_name,omitempty" db:"instance_os_name"`
}

type QueryParams struct {
	Status     InstanceState
	Stage      string
	Platform   *Platform
	RunnerName string
}

type StageOwner struct {
	StageID  string `db:"stage_id" json:"stage_id"`
	PoolName string `db:"pool_name" json:"pool_name"`
}
