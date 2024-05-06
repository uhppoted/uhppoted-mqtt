package commands

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"text/template"

	"github.com/uhppoted/uhppoted-lib/config"
)

var DAEMONIZE = Daemonize{
	usergroup: "uhppoted:uhppoted",
	workdir:   "/var/uhppoted",
	logdir:    "/var/log/uhppoted",
	config:    "/etc/uhppoted/uhppoted.conf",
	hotp:      "/etc/uhppoted/mqtt.hotp.secrets",
}

type usergroup string

type info struct {
	Description   string
	Documentation string
	Executable    string
	PID           string
	User          string
	Group         string
	Uid           int
	Gid           int
	LogFiles      []string
}

const serviceTemplate = `[Unit]
Description={{.Description}}
Documentation={{.Documentation}}
After=syslog.target network-online.target
Wants=syslog.target network-online.target

[Service]
Type=simple
ExecStart={{.Executable}}
PIDFile={{.PID}}
User={{.User}}
Group={{.Group}}

[Install]
WantedBy=multi-user.target
`

const logRotateTemplate = `{{range .LogFiles}}{{.}} {{end}}{
    daily
    rotate 30
    compress
        compresscmd /bin/bzip2
        compressext .bz2
        dateext
    missingok
    notifempty
    su uhppoted uhppoted
    postrotate
       /usr/bin/killall -HUP uhppoted
    endscript
}
`

type Daemonize struct {
	usergroup usergroup
	workdir   string
	logdir    string
	config    string
	hotp      string
}

func (cmd *Daemonize) Name() string {
	return "daemonize"
}

func (cmd *Daemonize) FlagSet() *flag.FlagSet {
	flagset := flag.NewFlagSet("daemonize", flag.ExitOnError)
	flagset.Var(&cmd.usergroup, "user", "user:group for uhppoted service")

	return flagset
}

func (cmd *Daemonize) Description() string {
	return fmt.Sprintf("Daemonizes %s as a service/daemon", SERVICE)
}

func (cmd *Daemonize) Usage() string {
	return "daemonize [--user <user:group>]"
}

func (cmd *Daemonize) Help() {
	fmt.Println()
	fmt.Println("  Usage: uhppoted daemonize [--user <user:group>]")
	fmt.Println()
	fmt.Printf("    Registers %s as a systemd service/daemon that runs on startup.\n", SERVICE)
	fmt.Println("      Defaults to the user:group uhppoted:uhppoted unless otherwise specified")
	fmt.Println("      with the --user option")
	fmt.Println()

	helpOptions(cmd.FlagSet())
}

func (cmd *Daemonize) Execute(args ...interface{}) error {
	dir := filepath.Dir(cmd.config)
	r := bufio.NewReader(os.Stdin)

	fmt.Println()
	fmt.Printf("     **** PLEASE MAKE SURE YOU HAVE A BACKUP COPY OF THE CONFIGURATION INFORMATION AND KEYS IN %s ***\n", dir)
	fmt.Println()
	fmt.Printf("     Enter 'yes' to continue with the installation: ")

	text, err := r.ReadString('\n')
	if err != nil || strings.TrimSpace(text) != "yes" {
		fmt.Println()
		fmt.Printf("     -- installation cancelled --")
		fmt.Println()
		return nil
	}

	return cmd.execute()
}

func (cmd *Daemonize) execute() error {
	fmt.Println()
	fmt.Println("   ... daemonizing")

	executable, err := os.Executable()
	if err != nil {
		return err
	}

	uid, gid, err := getUserGroup(string(cmd.usergroup))
	if err != nil {
		return err
	}

	bind, _, _ := config.DefaultIpAddresses()

	i := info{
		Description:   "UHPPOTE UTO311-L0x access card controllers MQTT service/daemon ",
		Documentation: "https://github.com/uhppoted/uhppoted-mqtt",
		Executable:    executable,
		PID:           fmt.Sprintf("/var/uhppoted/%s.pid", SERVICE),
		User:          "uhppoted",
		Group:         "uhppoted",
		Uid:           uid,
		Gid:           gid,
		LogFiles:      []string{fmt.Sprintf("/var/log/uhppoted/%s.log", SERVICE)},
	}

	if err := cmd.systemd(&i); err != nil {
		return err
	}

	if err := cmd.mkdirs(&i); err != nil {
		return err
	}

	if err := cmd.logrotate(&i); err != nil {
		return err
	}

	if err := cmd.conf(&i); err != nil {
		return err
	}

	if err := cmd.genkeys(&i); err != nil {
		return err
	}

	fmt.Printf("   ... %s registered as a systemd service\n", SERVICE)
	fmt.Println()
	fmt.Println("   The daemon will start automatically on the next system restart - to start it manually, execute the following command:")
	fmt.Println()
	fmt.Printf("     > sudo systemctl start  %v\n", SERVICE)
	fmt.Printf("     > sudo systemctl status %v\n", SERVICE)
	fmt.Println()
	fmt.Println("   For some system configurations it may be necessary to also enable the service:")
	fmt.Println()
	fmt.Printf("     > sudo systemctl enable %v\n", SERVICE)
	fmt.Println()
	fmt.Println()
	fmt.Println("   The firewall may need additional rules to allow UDP broadcast e.g. for UFW:")
	fmt.Println()
	fmt.Printf("     > sudo ufw allow from %s to any port 60000 proto udp\n", bind.Addr())
	fmt.Println()
	fmt.Println("   Please replace the default RSA keys for event and system messages:")
	fmt.Printf("     - %s\n", filepath.Join(filepath.Dir(cmd.config), "mqtt", "rsa", "encryption", "event.pub"))
	fmt.Printf("     - %s\n", filepath.Join(filepath.Dir(cmd.config), "mqtt", "rsa", "encryption", "system.pub"))
	fmt.Println()

	return nil
}

func (cmd *Daemonize) systemd(i *info) error {
	service := fmt.Sprintf("%s.service", SERVICE)
	path := filepath.Join("/etc/systemd/system", service)
	t := template.Must(template.New(service).Parse(serviceTemplate))

	fmt.Printf("   ... creating '%s'\n", path)
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	defer f.Close()

	return t.Execute(f, i)
}

func (cmd *Daemonize) mkdirs(i *info) error {
	directories := []string{
		"/var/uhppoted",
		"/var/log/uhppoted",
		"/etc/uhppoted",
		"/etc/uhppoted/mqtt",
	}

	for _, dir := range directories {
		fmt.Printf("   ... creating '%s'\n", dir)

		if err := os.MkdirAll(dir, 0770); err != nil {
			return err
		}

		if err := os.Chown(dir, i.Uid, i.Gid); err != nil {
			return err
		}
	}

	return nil
}

func (cmd *Daemonize) logrotate(i *info) error {
	path := filepath.Join("/etc/logrotate.d", SERVICE)
	t := template.Must(template.New(fmt.Sprintf("%s.logrotate", SERVICE)).Parse(logRotateTemplate))

	fmt.Printf("   ... creating '%s'\n", path)
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	defer f.Close()

	return t.Execute(f, i)
}

func (cmd *Daemonize) conf(i *info) error {
	path := cmd.config

	fmt.Printf("   ... creating '%s'\n", path)

	// initialise config from existing uhppoted.conf
	cfg := config.NewConfig()
	if f, err := os.Open(path); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else {
		err := cfg.Read(f)
		f.Close()
		if err != nil {
			return err
		}
	}

	// generate HMAC and RSA keys
	if cfg.MQTT.HMAC.Key == "" {
		hmac, err := hmac()
		if err != nil {
			return err
		}

		cfg.MQTT.HMAC.Key = hmac
	}

	// write back config with any updated information
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	defer f.Close()

	return cfg.Write(f)
}

func (cmd *Daemonize) genkeys(i *info) error {
	return genkeys(filepath.Dir(cmd.config), cmd.hotp)
}

// usergroup::flag.Value

func getUserGroup(s string) (int, int, error) {
	match := regexp.MustCompile(`(\w+?):(\w+)`).FindStringSubmatch(s)
	if match == nil {
		return 0, 0, fmt.Errorf("invalid user:group '%s'", s)
	}

	u, err := user.Lookup(match[1])
	if err != nil {
		return 0, 0, err
	}

	g, err := user.LookupGroup(match[2])
	if err != nil {
		return 0, 0, err
	}

	uid, err := strconv.Atoi(u.Uid)
	if err != nil {
		return 0, 0, err
	}

	gid, err := strconv.Atoi(g.Gid)
	if err != nil {
		return 0, 0, err
	}

	return uid, gid, nil
}

func (f *usergroup) String() string {
	if f == nil {
		return "uhppoted:uhppoted"
	}

	return string(*f)
}

func (f *usergroup) Set(s string) error {
	_, _, err := getUserGroup(s)
	if err != nil {
		return err
	}

	*f = usergroup(strings.TrimSpace(s))

	return nil
}
