package commands

import (
	"flag"
	"fmt"
	syslog "log"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/eventlog"

	"github.com/uhppoted/uhppote-core/uhppote"
	"github.com/uhppoted/uhppoted-lib/config"
	filelogger "github.com/uhppoted/uhppoted-lib/eventlog"
	"github.com/uhppoted/uhppoted-mqtt/log"
)

type service struct {
	name   string
	conf   *config.Config
	logger *syslog.Logger
	cmd    *Run
}

type EventLog struct {
	log *eventlog.Log
}

var RUN = Run{
	configuration: filepath.Join(workdir(), "uhppoted.conf"),
	dir:           workdir(),
	pidFile:       filepath.Join(workdir(), fmt.Sprintf("%s.pid", SERVICE)),
	logLevel:      "info",
	logFile:       filepath.Join(workdir(), "logs", fmt.Sprintf("%s.log", SERVICE)),
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
	flagset.BoolVar(&r.console, "console", r.console, "Run as command-line application")
	flagset.BoolVar(&r.debug, "debug", r.debug, "Displays internal information for diagnosing errors")

	return flagset
}

func (r *Run) Execute(args ...interface{}) error {
	log.Infof("", "%s service %s - %s (PID %d)\n", SERVICE, uhppote.VERSION, "Microsoft Windows", os.Getpid())

	f := func(c *config.Config) error {
		return r.start(c)
	}

	return r.execute(f)
}

func (cmd *Run) start(c *config.Config) error {
	var logger *syslog.Logger

	if cmd.console {
		logger = syslog.New(os.Stdout, "", syslog.LstdFlags)
	} else if eventlogger, err := eventlog.Open(SERVICE); err == nil {
		defer eventlogger.Close()

		events := EventLog{eventlogger}
		logger = syslog.New(&events, SERVICE, syslog.Ldate|syslog.Ltime|syslog.LUTC)
	} else {
		events := filelogger.Ticker{Filename: cmd.logFile, MaxSize: cmd.logFileSize}
		logger = syslog.New(&events, "", syslog.Ldate|syslog.Ltime|syslog.LUTC)
	}

	log.SetLogger(logger)
	log.SetLevel(cmd.logLevel)
	log.Infof("", "%s service - start\n", SERVICE)

	if cmd.console {
		interrupt := make(chan os.Signal, 1)

		signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
		cmd.run(c, logger, interrupt)
		return nil
	}

	uhppoted := service{
		name:   SERVICE,
		conf:   c,
		logger: logger,
		cmd:    cmd,
	}

	log.Infof("", "%s service - starting\n", SERVICE)

	if err := svc.Run(SERVICE, &uhppoted); err != nil {
		fmt.Printf("   Unable to execute ServiceManager.Run request (%v)\n", err)
		fmt.Println()
		fmt.Printf("   To run %s as a command line application, type:\n", SERVICE)
		fmt.Println()
		fmt.Printf("     > %s --console\n", SERVICE)
		fmt.Println()

		log.Fatalf("", "Error executing ServiceManager.Run request: %v", err)
		return err
	}

	log.Infof("", "%s daemon - started\n", SERVICE)

	return nil
}

func (s *service) Execute(args []string, r <-chan svc.ChangeRequest, status chan<- svc.Status) (ssec bool, errno uint32) {
	s.logger.Printf("%s service - Execute\n", SERVICE)

	const commands = svc.AcceptStop | svc.AcceptShutdown

	status <- svc.Status{State: svc.StartPending}

	interrupt := make(chan os.Signal, 1)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		s.cmd.run(s.conf, s.logger, interrupt)

		s.logger.Printf("exit\n")
	}()

	status <- svc.Status{State: svc.Running, Accepts: commands}

loop:
	for c := range r {
		s.logger.Printf("%s service - select: %v  %v\n", SERVICE, c.Cmd, c.CurrentStatus)
		switch c.Cmd {
		case svc.Interrogate:
			s.logger.Printf("%s service - svc.Interrogate %v\n", SERVICE, c.CurrentStatus)
			status <- c.CurrentStatus

		case svc.Stop:
			interrupt <- syscall.SIGINT
			s.logger.Printf("%s service- svc.Stop\n", SERVICE)
			break loop

		case svc.Shutdown:
			interrupt <- syscall.SIGTERM
			s.logger.Printf("%s service - svc.Shutdown\n", SERVICE)
			break loop

		default:
			s.logger.Printf("%s service - svc.????? (%v)\n", SERVICE, c.Cmd)
		}
	}

	s.logger.Printf("%s service - stopping\n", SERVICE)
	status <- svc.Status{State: svc.StopPending}
	wg.Wait()
	status <- svc.Status{State: svc.Stopped}
	s.logger.Printf("%s service - stopped\n", SERVICE)

	return false, 0
}

func (e *EventLog) Write(p []byte) (int, error) {
	err := e.log.Info(1, string(p))
	if err != nil {
		return 0, err
	}

	return len(p), nil
}
