package commands

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	syslog "log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
	paho "github.com/eclipse/paho.mqtt.golang"

	"github.com/uhppoted/uhppote-core/uhppote"
	"github.com/uhppoted/uhppoted-lib/config"
	"github.com/uhppoted/uhppoted-lib/locales"
	"github.com/uhppoted/uhppoted-lib/lockfile"
	"github.com/uhppoted/uhppoted-lib/monitoring"

	"github.com/uhppoted/uhppoted-mqtt/acl"
	"github.com/uhppoted/uhppoted-mqtt/auth"
	"github.com/uhppoted/uhppoted-mqtt/log"
	"github.com/uhppoted/uhppoted-mqtt/mqtt"
)

type Run struct {
	configuration       string
	dir                 string
	pidFile             string
	logLevel            string
	logFile             string
	logFileSize         int
	console             bool
	debug               bool
	healthCheckInterval time.Duration
	watchdogInterval    time.Duration
	//lockfile            config.Lockfile
}

const LOG_TAG = ""

func (cmd *Run) Name() string {
	return "run"
}

func (cmd *Run) Description() string {
	return fmt.Sprintf("Runs the %s daemon/service until terminated by the system service manager", SERVICE)
}

func (cmd *Run) Usage() string {
	return fmt.Sprintf("%s [run] [--console] [--config <file>] [--dir <workdir>] [--pid <file>] [--logfile <file>] [--logfilesize <bytes>] [--debug]", SERVICE)
}

func (cmd *Run) Help() {
	fmt.Println()
	fmt.Printf("  Usage: %s  [--console] [--config <file>] [--dir <workdir>] [--pid <file>] [--logfile <file>] [--logfilesize <bytes>] [--debug]\n", SERVICE)
	fmt.Println()

	helpOptions(cmd.FlagSet())
}

func (cmd *Run) execute(f func(*config.Config) error) error {
	conf := config.NewConfig()
	if err := conf.Load(cmd.configuration); err != nil {
		log.Warnf(LOG_TAG, "Could not load configuration (%v)", err)
	}

	if err := os.MkdirAll(cmd.dir, os.ModeDir|os.ModePerm); err != nil {
		return fmt.Errorf("unable to create working directory '%v': %v", cmd.dir, err)
	}

	// ... create lockfile
	pidfile := config.Lockfile{
		File:   cmd.pidFile,
		Remove: conf.LockfileRemove,
	}

	if pidfile.File == "" {
		pidfile.File = filepath.Join(os.TempDir(), fmt.Sprintf("%s.pid", SERVICE))
	}

	if lock, err := lockfile.MakeLockFile(pidfile); err != nil {
		return err
	} else {
		defer func() {
			lock.Release()
		}()

		log.AddFatalHook(func() {
			lock.Release()
		})
	}

	return f(conf)
}

func (cmd *Run) run(c *config.Config, logger *syslog.Logger, interrupt chan os.Signal) {
	log.SetLogger(logger)
	log.Infof(LOG_TAG, "START")

	cmd.healthCheckInterval = c.HealthCheckInterval
	cmd.watchdogInterval = c.WatchdogInterval

	// ... initialise MQTT

	bind, broadcast, listen := config.DefaultIpAddresses()

	if c.BindAddress != nil {
		bind = *c.BindAddress
	}

	if c.BroadcastAddress != nil {
		broadcast = *c.BroadcastAddress
	}

	if c.ListenAddress != nil {
		listen = *c.ListenAddress
	}

	devices := []uhppote.Device{}
	for id, d := range c.Devices {
		if device := uhppote.NewDevice(d.Name, id, d.Address, d.Doors); device != nil {
			devices = append(devices, *device)
		}
	}

	u := uhppote.NewUHPPOTE(bind, broadcast, listen, c.Timeout, devices, cmd.debug)

	permissions, err := auth.NewPermissions(
		c.MQTT.Permissions.Enabled,
		c.MQTT.Permissions.Users,
		c.MQTT.Permissions.Groups)
	if err != nil {
		log.Errorf(LOG_TAG, "%v", err)
		return
	}

	mqttd := mqtt.MQTTD{
		ServerID: c.ServerID,
		TLS:      &tls.Config{},
		Connection: mqtt.Connection{
			Broker:   fmt.Sprintf(c.Connection.Broker),
			ClientID: c.Connection.ClientID,
			UserName: c.Connection.Username,
			Password: c.Connection.Password,
		},
		Topics: mqtt.Topics{
			Requests: c.Topics.Resolve(c.Topics.Requests),
			Replies:  c.Topics.Resolve(c.Topics.Replies),
			Events: struct {
				Feed string
				Live string
			}{
				Feed: c.Topics.Resolve(c.Topics.EventsFeed),
				Live: c.Topics.Resolve(c.Topics.LiveEvents),
			},
			System: c.Topics.Resolve(c.Topics.System),
		},
		Alerts: mqtt.Alerts{
			QOS:      c.Alerts.QOS,
			Retained: c.Alerts.Retained,
		},
		Encryption: mqtt.Encryption{
			SignOutgoing:    c.SignOutgoing,
			EncryptOutgoing: c.EncryptOutgoing,
			EventsKeyID:     c.EventsKeyID,
			SystemKeyID:     c.SystemKeyID,
			HOTP:            nil,
		},
		Authentication: c.Authentication,
		Permissions:    *permissions,
		EventMap:       c.EventIDs,
		AWS:            mqtt.AWS{},
		ACL: mqtt.ACL{
			Verify: map[acl.Verification]bool{},
		},
		Protocol: c.MQTT.Protocol,

		Debug: cmd.debug,
	}

	if c.AWS.Credentials != "" {
		mqttd.AWS.Credentials = credentials.NewSharedCredentials(c.AWS.Credentials, c.AWS.Profile)
		mqttd.AWS.Region = c.AWS.Region
	}

	if strings.Contains(strings.ToLower(c.ACL.Verify), "none") {
		mqttd.ACL.Verify[acl.None] = true
	}

	if strings.Contains(strings.ToLower(c.ACL.Verify), "not-empty") {
		mqttd.ACL.Verify[acl.NotEmpty] = true
	}

	if strings.Contains(strings.ToLower(c.ACL.Verify), "rsa") {
		mqttd.ACL.Verify[acl.RSA] = true
	}

	// ... TLS
	allowInsecure := false
	if c.Connection.Verify == "allow-insecure" {
		allowInsecure = true
	}

	if strings.HasPrefix(mqttd.Connection.Broker, "tls:") {
		pem, err := os.ReadFile(c.Connection.BrokerCertificate)
		if err != nil {
			log.Errorf(LOG_TAG, "%v", err)
		} else {
			mqttd.TLS.InsecureSkipVerify = allowInsecure
			mqttd.TLS.RootCAs = x509.NewCertPool()

			if ok := mqttd.TLS.RootCAs.AppendCertsFromPEM(pem); !ok {
				log.Errorf(LOG_TAG, "could not initialise MQTTD CA certificates")
			}
		}

		certificate, err := tls.LoadX509KeyPair(c.Connection.ClientCertificate, c.Connection.ClientKey)
		if err != nil {
			log.Errorf(LOG_TAG, "%v", err)
		} else {
			mqttd.TLS.Certificates = []tls.Certificate{certificate}
		}
	}

	// ... authentication
	hmac, err := auth.NewHMAC(c.HMAC.Required, c.HMAC.Key)
	if c.HMAC.Required && err != nil {
		log.Errorf(LOG_TAG, "%v", err)
		return
	}

	hotp, err := auth.NewHOTP(c.MQTT.HOTP.Range, c.MQTT.HOTP.Secrets, c.MQTT.HOTP.Counters)
	if mqttd.Authentication == "HOTP" && err != nil {
		log.Errorf(LOG_TAG, "%v", err)
		return
	}

	rsa, err := auth.NewRSA(c.RSA.KeyDir)
	if mqttd.Authentication == "RSA" && err != nil {
		log.Errorf(LOG_TAG, "%v", err)
		return
	}

	nonce, err := auth.NewNonce(c.Nonce.Required, c.Nonce.Server, c.Nonce.Clients)
	if err != nil {
		log.Errorf(LOG_TAG, "%v", err)
		return
	}

	mqttd.HMAC = *hmac
	mqttd.Encryption.HOTP = hotp
	mqttd.Encryption.RSA = rsa
	mqttd.Encryption.Nonce = *nonce

	// ... locales
	if c.MQTT.Locale != "" {
		folder := filepath.Dir(c.MQTT.Locale)
		file := filepath.Base(c.MQTT.Locale)
		fs := os.DirFS(folder)
		if err := locales.Load(fs, file); err != nil {
			log.Warnf(LOG_TAG, "%v", err)
		} else {
			log.Infof(LOG_TAG, "using translations from %v", c.MQTT.Locale)
		}
	}

	// ... monitoring
	healthcheck := monitoring.NewHealthCheck(
		u,
		c.HealthCheckInterval,
		c.HealthCheckIdle,
		c.HealthCheckIgnore)

	mqtt.SetDisconnectsEnabled(c.MQTT.Disconnects.Enabled)
	mqtt.SetDisconnectsInterval(c.MQTT.Disconnects.Interval)
	mqtt.SetMaxDisconnects(c.MQTT.Disconnects.Max)

	// ... authorized card list
	cards, err := authorized(c.MQTT.Cards)
	if err != nil {
		log.Warnf(LOG_TAG, "%v", err)
	}

	// ... listen
	err = cmd.listen(u, &mqttd, devices, &healthcheck, cards, logger, interrupt)
	if err != nil {
		log.Errorf(LOG_TAG, "%v", err)
	}

	log.Infof(LOG_TAG, "exit")
}

func (r *Run) listen(
	u uhppote.IUHPPOTE,
	mqttd *mqtt.MQTTD,
	devices []uhppote.Device,
	healthcheck *monitoring.HealthCheck,
	authorized []string,
	logger *syslog.Logger,
	interrupt chan os.Signal) error {

	// ... acquire MQTT client lock
	dir := os.TempDir()
	if r.pidFile != "" {
		dir = filepath.Dir(r.pidFile)
	}

	clientlock := config.Lockfile{
		File:   filepath.Join(dir, fmt.Sprintf("%s.lock", mqttd.Connection.ClientID)),
		Remove: lockfile.RemoveLockfile,
	}

	if kraken, err := lockfile.MakeLockFile(clientlock); err != nil {
		return err
	} else {
		defer func() {
			kraken.Release()
		}()

		log.AddFatalHook(func() {
			kraken.Release()
		})
	}

	// ... MQTT
	paho.CRITICAL = logger
	paho.ERROR = logger
	paho.WARN = logger

	if mqttd.Debug {
		paho.DEBUG = logger
	}

	if err := mqttd.Run(u, devices, authorized); err != nil {
		return err
	}

	defer mqttd.Close()

	// ... monitoring
	monitor := mqtt.NewSystemMonitor(mqttd)
	watchdog := monitoring.NewWatchdog(healthcheck)
	k := time.NewTicker(r.healthCheckInterval)

	defer k.Stop()

	go func() {
		for {
			<-k.C
			healthcheck.Exec(monitor)
		}
	}()

	// ... wait until interrupted
	w := time.NewTicker(r.watchdogInterval)

	defer w.Stop()

	for {
		select {
		case <-w.C:
			if err := watchdog.Exec(monitor); err != nil {
				return err
			}

		case <-interrupt:
			log.Infof(LOG_TAG, "interrupt")
			return nil
		}
	}
}

func authorized(file string) ([]string, error) {
	b, err := os.ReadFile(file)
	if err != nil {
		return []string{}, err
	}

	lines := []string{}
	scanner := bufio.NewScanner(bytes.NewBuffer(b))
	for scanner.Scan() {
		lines = append(lines, strings.TrimSpace(scanner.Text()))
	}

	return lines, scanner.Err()
}
