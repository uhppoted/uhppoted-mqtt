package commands

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"

	"github.com/uhppoted/uhppote-core/uhppote"
	"github.com/uhppoted/uhppoted-lib/config"
	"github.com/uhppoted/uhppoted-lib/locales"
	"github.com/uhppoted/uhppoted-lib/monitoring"
	"github.com/uhppoted/uhppoted-mqtt/auth"
	"github.com/uhppoted/uhppoted-mqtt/mqtt"
)

type Run struct {
	configuration       string
	dir                 string
	pidFile             string
	logFile             string
	logFileSize         int
	console             bool
	debug               bool
	healthCheckInterval time.Duration
	watchdogInterval    time.Duration
}

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
		log.Printf("WARN  Could not load configuration (%v)", err)
	}

	if err := os.MkdirAll(cmd.dir, os.ModeDir|os.ModePerm); err != nil {
		return fmt.Errorf("Unable to create working directory '%v': %v", cmd.dir, err)
	}

	pid := fmt.Sprintf("%d\n", os.Getpid())

	_, err := os.Stat(cmd.pidFile)
	if err == nil {
		return fmt.Errorf("PID lockfile '%v' already in use", cmd.pidFile)
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("Error checking PID lockfile '%v' (%v)", cmd.pidFile, err)
	}

	if err := os.WriteFile(cmd.pidFile, []byte(pid), 0644); err != nil {
		return fmt.Errorf("Unable to create PID lockfile: %v", err)
	}

	defer func() {
		os.Remove(cmd.pidFile)
	}()

	return f(conf)
}

func (cmd *Run) run(c *config.Config, logger *log.Logger, interrupt chan os.Signal) {
	logger.Printf("START")

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
		c.MQTT.Permissions.Groups,
		logger)
	if err != nil {
		log.Printf("ERROR: %v", err)
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
			Events:   c.Topics.Resolve(c.Topics.Events),
			System:   c.Topics.Resolve(c.Topics.System),
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
		Protocol:       c.Protocol,

		Debug: cmd.debug,
	}

	if c.AWS.Credentials != "" {
		mqttd.AWS.Credentials = credentials.NewSharedCredentials(c.AWS.Credentials, c.AWS.Profile)
		mqttd.AWS.Region = c.AWS.Region
	}

	// ... TLS

	if strings.HasPrefix(mqttd.Connection.Broker, "tls:") {
		pem, err := os.ReadFile(c.Connection.BrokerCertificate)
		if err != nil {
			logger.Printf("ERROR: %v", err)
		} else {
			mqttd.TLS.InsecureSkipVerify = false
			mqttd.TLS.RootCAs = x509.NewCertPool()

			if ok := mqttd.TLS.RootCAs.AppendCertsFromPEM(pem); !ok {
				logger.Printf("ERROR: Could not initialise MQTTD CA certificates")
			}
		}

		certificate, err := tls.LoadX509KeyPair(c.Connection.ClientCertificate, c.Connection.ClientKey)
		if err != nil {
			logger.Printf("ERROR: %v", err)
		} else {
			mqttd.TLS.Certificates = []tls.Certificate{certificate}
		}
	}

	// ... authentication

	hmac, err := auth.NewHMAC(c.HMAC.Required, c.HMAC.Key)
	if c.HMAC.Required && err != nil {
		logger.Printf("ERROR: %v", err)
		return
	}

	hotp, err := auth.NewHOTP(c.MQTT.HOTP.Range, c.MQTT.HOTP.Secrets, c.MQTT.HOTP.Counters, logger)
	if mqttd.Authentication == "HOTP" && err != nil {
		logger.Printf("ERROR: %v", err)
		return
	}

	rsa, err := auth.NewRSA(c.RSA.KeyDir, logger)
	if mqttd.Authentication == "RSA" && err != nil {
		logger.Printf("ERROR: %v", err)
		return
	}

	nonce, err := auth.NewNonce(c.Nonce.Required, c.Nonce.Server, c.Nonce.Clients, logger)
	if err != nil {
		logger.Printf("ERROR: %v", err)
		return
	}

	mqttd.HMAC = *hmac
	mqttd.Encryption.HOTP = hotp
	mqttd.Encryption.RSA = rsa
	mqttd.Encryption.Nonce = *nonce

	// ... locales

	if c.Locale != "" {
		folder := filepath.Dir(c.Locale)
		file := filepath.Base(c.Locale)
		fs := os.DirFS(folder)
		if err := locales.Load(fs, file); err != nil {
			logger.Printf("WARN  %v", err)
		} else {
			logger.Printf("INFO  using translations from %v", c.Locale)
		}
	}

	// ... monitoring

	healthcheck := monitoring.NewHealthCheck(u, c.HealthCheckIdle, c.HealthCheckIgnore, logger)

	// ... authorized card list

	cards, err := authorized(c.MQTT.Cards)
	if err != nil {
		logger.Printf("WARN  %v", err)
	}

	// ... listen

	err = cmd.listen(u, &mqttd, devices, &healthcheck, cards, logger, interrupt)
	if err != nil {
		logger.Printf("ERROR %v", err)
	}

	logger.Printf("INFO  exit")
}

func (r *Run) listen(
	u uhppote.IUHPPOTE,
	mqttd *mqtt.MQTTD,
	devices []uhppote.Device,
	healthcheck *monitoring.HealthCheck,
	authorized []string,
	logger *log.Logger,
	interrupt chan os.Signal) error {

	// ... MQTT

	pid := fmt.Sprintf("%d\n", os.Getpid())
	workdir := filepath.Dir(r.pidFile)
	lockfile := filepath.Join(workdir, fmt.Sprintf("%s.lock", mqttd.Connection.ClientID))

	_, err := os.Stat(lockfile)
	if err == nil {
		return fmt.Errorf("MQTT client lockfile '%v' already in use", lockfile)
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("Error checking MQTT client lockfile '%v' (%v_)", lockfile, err)
	}

	if err := os.WriteFile(lockfile, []byte(pid), 0644); err != nil {
		return fmt.Errorf("Unable to create MQTT client lockfile: %v", err)
	}

	defer func() {
		os.Remove(lockfile)
	}()

	if err := mqttd.Run(u, devices, authorized, logger); err != nil {
		return err
	}

	defer mqttd.Close(logger)

	// ... monitoring

	monitor := mqtt.NewSystemMonitor(mqttd, logger)
	watchdog := monitoring.NewWatchdog(healthcheck, logger)
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
			logger.Printf("... interrupt")
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
