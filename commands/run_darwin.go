package commands

import (
	"flag"
	"fmt"
	syslog "log"
	"os"
	"os/signal"
	"syscall"

	"github.com/uhppoted/uhppote-core/uhppote"
	"github.com/uhppoted/uhppoted-lib/config"
	"github.com/uhppoted/uhppoted-lib/eventlog"
	"github.com/uhppoted/uhppoted-mqtt/log"
)

var RUN = Run{
	configuration: "/usr/local/etc/com.github.uhppoted/uhppoted.conf",
	dir:           "/usr/local/var/com.github.uhppoted",
	logLevel:      "info",
	pidFile:       fmt.Sprintf("/usr/local/var/com.github.uhppoted/%s.pid", SERVICE),
	logFile:       fmt.Sprintf("/usr/local/var/com.github.uhppoted/logs/%s.log", SERVICE),
	logFileSize:   10,
	console:       false,
	debug:         false,
}

func (r *Run) FlagSet() *flag.FlagSet {
	flagset := flag.NewFlagSet("run", flag.ExitOnError)

	flagset.StringVar(&r.configuration, "config", r.configuration, "Sets the configuration file path")
	flagset.StringVar(&r.dir, "dir", r.dir, "Work directory")
	flagset.StringVar(&r.pidFile, "pid", r.pidFile, "Sets the service PID file path")
	flagset.StringVar(&r.logLevel, "log-level", r.logFile, "Sets the logging level (debug, info, warning or error)")
	flagset.StringVar(&r.logFile, "logfile", r.logFile, "Sets the log file path")
	flagset.IntVar(&r.logFileSize, "logfilesize", r.logFileSize, "Sets the log file size before forcing a log rotate")
	flagset.BoolVar(&r.console, "console", r.console, "Writes log entries to stdout")
	flagset.BoolVar(&r.debug, "debug", r.debug, "Displays internal information for diagnosing errors")

	return flagset
}

func (r *Run) Execute(args ...interface{}) error {
	log.Infof("", "%s service %s - %s (PID %d)\n", SERVICE, uhppote.VERSION, "MacOS", os.Getpid())

	f := func(c *config.Config) error {
		return r.exec(c)
	}

	return r.execute(f)
}

func (cmd *Run) exec(c *config.Config) error {
	logger := syslog.New(os.Stdout, "", syslog.LstdFlags)
	interrupt := make(chan os.Signal, 1)

	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	if !cmd.console {
		events := eventlog.Ticker{Filename: cmd.logFile, MaxSize: cmd.logFileSize}
		logger = syslog.New(&events, "", syslog.Ldate|syslog.Ltime|syslog.LUTC)
		rotate := make(chan os.Signal, 1)

		signal.Notify(rotate, syscall.SIGHUP)

		go func() {
			for {
				<-rotate
				log.Infof("", "Rotating %s log file '%s'\n", SERVICE, cmd.logFile)
				events.Rotate()
			}
		}()
	}

	log.SetLogger(logger)
	log.SetLevel(cmd.logLevel)

	cmd.run(c, logger, interrupt)

	return nil
}
