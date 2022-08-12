package commands

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	syslog "log"
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
	"github.com/uhppoted/uhppoted-mqtt/log"
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
	softlock            softlock
}

type softlock struct {
	enabled  bool
	interval time.Duration
	wait     time.Duration
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

	log.SetFatalHook(func() {
		os.Remove(cmd.pidFile)
	})

	return f(conf)
}

func (cmd *Run) run(c *config.Config, logger *syslog.Logger, interrupt chan os.Signal) {
	log.SetLogger(logger)
	log.Infof(LOG_TAG, "START")

	cmd.healthCheckInterval = c.HealthCheckInterval
	cmd.watchdogInterval = c.WatchdogInterval

	cmd.softlock.enabled = c.MQTT.Softlock.Enabled
	cmd.softlock.interval = c.MQTT.Softlock.Interval
	cmd.softlock.wait = c.MQTT.Softlock.Wait

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
		Protocol:       c.MQTT.Protocol,

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
			log.Errorf(LOG_TAG, "%v", err)
		} else {
			mqttd.TLS.InsecureSkipVerify = false
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

	rsa, err := auth.NewRSA(c.RSA.KeyDir, logger)
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

	healthcheck := monitoring.NewHealthCheck(u, c.HealthCheckIdle, c.HealthCheckIgnore, logger)

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
	if lockfile, err := r.lock(mqttd.Connection.ClientID, interrupt); err != nil {
		return err
	} else if lockfile == "" {
		return fmt.Errorf("invalid MQTT client lockfile '%v'", lockfile)
	} else {
		defer func() {
			os.Remove(lockfile)
		}()
	}

	// ... MQTT

	if err := mqttd.Run(u, devices, authorized, logger); err != nil {
		return err
	}

	defer mqttd.Close()

	// ... monitoring

	monitor := mqtt.NewSystemMonitor(mqttd)
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

// Uses SHA-256 hash of lockfile contents because os.Stat updates the mtime
// of the lockfile and in any event, mtime has a resolution of 1 minute.
func (r *Run) lock(clientID string, interrupt chan os.Signal) (string, error) {
	workdir := filepath.Dir(r.pidFile)
	lockfile := filepath.Join(workdir, fmt.Sprintf("%s.lock", clientID))

	hash := func(file string) (string, error) {
		if bytes, err := os.ReadFile(file); err != nil {
			return "", err
		} else {
			sum := sha256.Sum256(bytes)

			return hex.EncodeToString(sum[:]), nil
		}
	}

	touch := func() error {
		pid := fmt.Sprintf("%d", os.Getpid())
		now := time.Now().Format("2006-01-02 15:04:05")
		v := fmt.Sprintf("%v\n%v\n", pid, now)

		return os.WriteFile(lockfile, []byte(v), 0644)
	}

	checksum, err := hash(lockfile)

	switch {
	case err != nil && !os.IsNotExist(err):
		return "", err

	case err != nil && os.IsNotExist(err):
		if err := touch(); err != nil {
			return "", err
		}

	case err == nil && !r.softlock.enabled:
		return "", fmt.Errorf("MQTT client lockfile '%v' already in use", lockfile)

	case err == nil && r.softlock.enabled && checksum == "":
		return "", fmt.Errorf("invalid MQTT client lockfile checksum")

	case err == nil && r.softlock.enabled && checksum != "":
		log.Warnf(LOG_TAG, "MQTT client lockfile '%v' exists, delaying for %v", lockfile, r.softlock.wait)

		wait := time.After(r.softlock.wait)

		select {
		case <-wait:

		case <-interrupt:
			return "", fmt.Errorf("interrupted")
		}

		h, err := hash(lockfile)
		switch {
		case err != nil && !os.IsNotExist(err):
			return "", err

		case err != nil && os.IsNotExist(err):
			if err := touch(); err != nil {
				return "", err
			}

		case h != checksum:
			return "", fmt.Errorf("MQTT client lockfile '%v' in use", lockfile)

		default:
			log.Warnf(LOG_TAG, "replacing MQTT client lockfile '%v'", lockfile)
			if err := touch(); err != nil {
				return "", err
			}
		}

	default:
		return "", fmt.Errorf("failed to acquire MQTT client lock")
	}

	if r.softlock.enabled {
		tick := time.Tick(r.softlock.interval)
		go func() {
			for {
				<-tick
				log.Infof(LOG_TAG, "touching MQTT client lockfile '%v'", lockfile)
				touch()
			}
		}()
	}

	return lockfile, nil
}
