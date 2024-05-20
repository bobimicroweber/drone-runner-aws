// Copyright 2020 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package setup

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/drone-runners/drone-runner-aws/command/config"
	"github.com/drone-runners/drone-runner-aws/internal/drivers"
	"github.com/drone-runners/drone-runner-aws/internal/lehelper"
	"github.com/drone-runners/drone-runner-aws/internal/poolfile"
	"github.com/drone-runners/drone-runner-aws/store/database"
	"github.com/drone/runner-go/client"
	"github.com/drone/runner-go/logger"
	"github.com/drone/signal"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
)

// empty context.
var nocontext = context.Background()

type setupCommand struct {
	envFile             string
	awsAccessKeyID      string
	awsAccessKeySecret  string
	digitalOceanPAT     string
	hetznerToken        string
	googleProjectID     string
	googleJSONPath      string
	ankaVMName          string
	azureClientID       string
	azureClientSecret   string
	azureTenantID       string
	azureSubscriptionID string
	ankabuildVMName     string
	ankabuildURL        string
	ankabuildToken      string
}

const (
	testPoolName    = "testpool"
	runnerName      = "setup"
	healthCheckWait = time.Minute * 10
)

func (c *setupCommand) run(*kingpin.ParseContext) error { //nolint
	// load environment variables from file.
	err := godotenv.Load(c.envFile)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	// load the configuration from the environment
	env, err := config.FromEnviron()
	if err != nil {
		return err
	}
	// check cli options and map arguments to configuration
	switch {
	case c.ankaVMName != "":
		env.Anka.VMName = c.ankaVMName
		logrus.Infoln("setup: using anka")
	case c.ankabuildVMName != "" && c.ankabuildURL != "":
		logrus.Infoln("setup: using anka build")
		env.AnkaBuild.Token = c.ankabuildToken
		env.AnkaBuild.URL = c.ankabuildURL
		env.AnkaBuild.VMName = c.ankabuildVMName
	case c.awsAccessKeyID != "" && c.awsAccessKeySecret != "":
		logrus.Infoln("setup: using amazon")
		env.AWS.AccessKeyID = c.awsAccessKeyID
		env.AWS.AccessKeySecret = c.awsAccessKeySecret
	case c.azureClientID != "" && c.azureClientSecret != "" && c.azureSubscriptionID != "" && c.azureTenantID != "":
		env.Azure.ClientID = c.azureClientID
		env.Azure.ClientSecret = c.azureClientSecret
		env.Azure.SubscriptionID = c.azureSubscriptionID
		env.Azure.TenantID = c.azureTenantID
		logrus.Infoln("setup: using azure")
	case c.digitalOceanPAT != "":
		env.DigitalOcean.PAT = c.digitalOceanPAT
		logrus.Infoln("setup: using digital ocean")
    case c.hetznerToken != "":
		env.Hetzner.Token = c.hetznerToken
		logrus.Infoln("setup: using hetzner token")
	case c.googleProjectID != "":
		env.Google.ProjectID = c.googleProjectID
		// use the default path if the user did not specify one
		if c.googleJSONPath != "" {
			env.Google.JSONPath = c.googleJSONPath
		}
		logrus.Infoln("setup: using google")
	default:
		logrus.
			Fatalln(`unsupported driver, please choose a driver setting the mandatory fields:
							for Amazon        --aws-access-key-id and --aws-access-key-secret
							for Anka          --anka-vm-name
							for AnkaBuild     --ankabuild-vm-name --ankabuild-url --ankabuild-token
							for Azure         --azure-client-id --azure-client-secret --azure-subscription-id --azure-tenant-id
							for Digital Ocean --digital-ocean-pat
							for Hetzner       --hetzner-token
							for Google        --google-project-id`)
	}

	// setup the global logrus logger.
	setupLogger(&env)

	ctx, cancel := context.WithCancel(nocontext)
	defer cancel()
	// listen for termination signals to gracefully shutdown.
	ctx = signal.WithContextFunc(ctx, func() {
		println("setup: received signal, terminating process")
		cancel()
	})

	cli := client.New(
		env.Client.Address,
		env.Client.Secret,
		env.Client.SkipVerify,
	)
	if env.Client.Dump {
		cli.Dumper = logger.StandardDumper(
			env.Client.DumpBody,
		)
	}
	cli.Logger = logger.Logrus(
		logrus.NewEntry(
			logrus.StandardLogger(),
		),
	)

	// use a single instance db, as we only need one machine
	store, _, err := database.ProvideStore(database.SingleInstance, "")
	if err != nil {
		logrus.WithError(err).Fatalln("Unable to start the database")
	}

	poolManager := drivers.New(ctx, store, &env)

	configPool, confErr := poolfile.ConfigPoolFile("", &env)
	if confErr != nil {
		logrus.WithError(confErr).
			Fatalln("Unable to load pool file, or use an in memory pool")
	}
	// process the pool file
	pools, processErr := poolfile.ProcessPool(configPool, runnerName)
	if processErr != nil {
		logrus.WithError(processErr).
			Fatalln("setup: unable to process pool file")
	}
	// there is only one instance
	addErr := poolManager.Add(pools[0])
	if addErr != nil {
		logrus.WithError(addErr).
			Fatalln("setup: unable to add pool")
	}
	// provision
	instance, provisionErr := poolManager.Provision(ctx, testPoolName, runnerName, runnerName, "drone", "", &env, nil)
	if provisionErr != nil {
		consoleLogs, consoleErr := poolManager.InstanceLogs(ctx, testPoolName, instance.ID)
		logrus.Infof("setup: instance logs for %s: %s", instance.ID, consoleLogs)
		logrus.WithError(provisionErr).
			WithError(consoleErr).
			Fatalln("setup: unable to provision instance")
	}
	// display the console logs
	consoleLogs, consoleLogsErr := poolManager.InstanceLogs(ctx, testPoolName, instance.ID)
	if consoleLogsErr != nil {
		logrus.WithError(consoleLogsErr).
			Errorln("setup: unable to get instance logs")
	}
	logrus.Infof("setup: instance logs for %s: %s", instance.ID, consoleLogs)
	// start the instance
	_, startErr := poolManager.StartInstance(ctx, testPoolName, instance.ID)
	if startErr != nil {
		cleanErr := poolManager.Destroy(ctx, testPoolName, instance.ID)
		consoleLogs, consoleErr := poolManager.InstanceLogs(ctx, testPoolName, instance.ID)
		logrus.Infof("setup: instance logs for %s: %s", instance.ID, consoleLogs)
		logrus.WithError(startErr).
			WithError(cleanErr).
			WithError(consoleErr).
			WithField("instance", instance.ID).
			Fatalln("setup: unable to start instance")
	}
	// create an LE client so we can test the instance
	leClient, leErr := lehelper.GetClient(instance, runnerName, instance.Port, env.LiteEngine.EnableMock, env.LiteEngine.MockStepTimeoutSecs)
	if leErr != nil {
		cleanErr := poolManager.Destroy(ctx, testPoolName, instance.ID)
		consoleLogs, consoleErr := poolManager.InstanceLogs(ctx, testPoolName, instance.ID)
		logrus.Infof("setup: instance logs for %s: %s", instance.ID, consoleLogs)
		logrus.WithError(leErr).
			WithError(cleanErr).
			WithError(consoleErr).
			WithField("instance", instance.ID).
			Fatalln("setup: unable to start lite engine")
	}

	// try the healthcheck api on the lite-engine until it responds ok
	logrus.Traceln("setup: running healthcheck and waiting for an ok response")
	performDNSLookup := drivers.ShouldPerformDNSLookup(ctx, instance.Platform.OS)

	healthResponse, healthErr := leClient.RetryHealth(ctx, healthCheckWait, performDNSLookup)
	if healthErr != nil {
		cleanErr := poolManager.Destroy(ctx, testPoolName, instance.ID)
		logrus.WithError(err).Errorln("failed health check with instance")
		consoleLogs, consoleErr := poolManager.InstanceLogs(ctx, testPoolName, instance.ID)
		logrus.Infof("setup: instance logs for %s: %s", instance.ID, consoleLogs)
		logrus.WithError(healthErr).
			WithError(cleanErr).
			WithError(consoleErr).
			WithField("instance", instance.ID).
			Fatalln("setup: health check failed")
	}
	logrus.WithField("response", fmt.Sprintf("%+v", healthResponse)).Traceln("LE.RetryHealth check complete")
	// print the pool file
	poolfile.PrintPoolFile(configPool)
	// finally clean the instance
	destroyErr := poolManager.Destroy(ctx, testPoolName, instance.ID)
	return destroyErr
}

func setupLogger(c *config.EnvConfig) {
	logger.Default = logger.Logrus(
		logrus.NewEntry(
			logrus.StandardLogger(),
		),
	)
	if c.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
	if c.Trace {
		logrus.SetLevel(logrus.TraceLevel)
	}
}

// Register the setup command.
func Register(app *kingpin.Application) {
	c := new(setupCommand)

	cmd := app.Command("setup", "sets up the runner").
		Action(c.run)
	cmd.Flag("envfile", "load the environment variable file").
		Default(".env").
		StringVar(&c.envFile)
	// Amazon specific flags
	cmd.Flag("aws-access-key-id", "aws access key ID").
		Default("").
		StringVar(&c.awsAccessKeyID)
	cmd.Flag("aws-access-key-secret", "aws access key secret").
		Default("").
		StringVar(&c.awsAccessKeySecret)
	// Anka specific flags
	cmd.Flag("anka-vm-name", "Anka VM name").
		Default("").
		StringVar(&c.ankaVMName)
	// Anka Build specific flags
	cmd.Flag("ankabuild-vm-name", "Anka Build VM name or ID").
		Default("").
		StringVar(&c.ankabuildVMName)
	cmd.Flag("ankabuild-url", "Anka Build Token").
		Default("").
		StringVar(&c.ankabuildURL)
	cmd.Flag("ankabuild-token", "Anka Build URL").
		Default("").
		StringVar(&c.ankabuildToken)
	// Azure specific flags
	cmd.Flag("azure-client-id", "Azure client ID").
		Default("").
		StringVar(&c.azureClientID)
	cmd.Flag("azure-client-secret", "Azure client secret").
		Default("").
		StringVar(&c.azureClientSecret)
	cmd.Flag("azure-subscription-id", "Azure subscription ID").
		Default("").
		StringVar(&c.azureSubscriptionID)
	cmd.Flag("azure-tenant-id", "Azure tenant ID").
		Default("").
		StringVar(&c.azureTenantID)
	// Digital Ocean specific flags
	cmd.Flag("digital-ocean-pat", "digital ocean token").
		Default("").
		StringVar(&c.digitalOceanPAT)
    // Hetzner specific flags
	cmd.Flag("hetzner-token", "hetzner token").
		Default("").
		StringVar(&c.hetznerToken)
	// Google specific flags
	cmd.Flag("google-project-id", "Google project ID").
		Default("").
		StringVar(&c.googleProjectID)
	cmd.Flag("google-json-path", "Google JSON path").
		StringVar(&c.googleJSONPath)
}
