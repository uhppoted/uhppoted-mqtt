package commands

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/uhppoted/uhppoted-api/config"
	xpath "github.com/uhppoted/uhppoted-api/encoding/plist"
)

type info struct {
	Label      string
	Executable string
	WorkDir    string
	StdLogFile string
	ErrLogFile string
}

type plist struct {
	Label             string
	Program           string
	WorkingDirectory  string
	ProgramArguments  []string
	KeepAlive         bool
	RunAtLoad         bool
	StandardOutPath   string
	StandardErrorPath string
}

const newsyslog = `#logfilename                                       [owner:group]  mode  count  size   when  flags [/pid_file]  [sig_num]
{{range .}}{{.LogFile}}  :              644   30     10000  @T00  J     {{.PID}}
{{end}}`

var DAEMONIZE = Daemonize{
	plist:   fmt.Sprintf("com.github.uhppoted.%s.plist", SERVICE),
	workdir: "/usr/local/var/com.github.uhppoted",
	logdir:  "/usr/local/var/com.github.uhppoted/logs",
	config:  "/usr/local/etc/com.github.uhppoted/uhppoted.conf",
	hotp:    "/usr/local/etc/com.github.uhppoted/mqtt.hotp.secrets",
}

type Daemonize struct {
	plist   string
	workdir string
	logdir  string
	config  string
	hotp    string
}

func (cmd *Daemonize) Name() string {
	return "daemonize"
}

func (cmd *Daemonize) FlagSet() *flag.FlagSet {
	return flag.NewFlagSet("daemonize", flag.ExitOnError)
}

func (cmd *Daemonize) Description() string {
	return fmt.Sprintf("Daemonizes %s as a service/daemon", SERVICE)
}

func (cmd *Daemonize) Usage() string {
	return ""
}

func (cmd *Daemonize) Help() {
	fmt.Println()
	fmt.Printf("  Usage: %s daemonize\n", SERVICE)
	fmt.Println()
	fmt.Printf("    Daemonizes %s as a service/daemon that runs on startup\n", SERVICE)
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

	i := info{
		Label:      fmt.Sprintf("com.github.uhppoted.%s", SERVICE),
		Executable: executable,
		WorkDir:    cmd.workdir,
		StdLogFile: filepath.Join(cmd.logdir, fmt.Sprintf("%s.log", SERVICE)),
		ErrLogFile: filepath.Join(cmd.logdir, fmt.Sprintf("%s.err", SERVICE)),
	}

	if err := cmd.launchd(&i); err != nil {
		return err
	}

	if err := cmd.mkdirs(); err != nil {
		return err
	}

	if err := cmd.logrotate(&i); err != nil {
		return err
	}

	if err := cmd.firewall(&i); err != nil {
		return err
	}

	if err := cmd.conf(&i); err != nil {
		return err
	}

	if err := cmd.genkeys(&i); err != nil {
		return err
	}

	fmt.Printf("   ... %s registered as a LaunchDaemon\n", i.Label)
	fmt.Println()
	fmt.Println("   The daemon will start automatically on the next system restart - to start it manually, execute the following command:")
	fmt.Println()
	fmt.Printf("   sudo launchctl load /Library/LaunchDaemons/com.github.uhppoted.%s.plist\n", SERVICE)
	fmt.Println()
	fmt.Println("   Please replace the default RSA keys for event and system messages:")
	fmt.Printf("     - %s\n", filepath.Join(filepath.Dir(cmd.config), "mqtt", "rsa", "encryption", "event.pub"))
	fmt.Printf("     - %s\n", filepath.Join(filepath.Dir(cmd.config), "mqtt", "rsa", "encryption", "system.pub"))
	fmt.Println()
	fmt.Println()

	return nil
}

func (cmd *Daemonize) launchd(i *info) error {
	path := filepath.Join("/Library/LaunchDaemons", cmd.plist)
	_, err := os.Stat(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	pl := plist{
		Label:             i.Label,
		Program:           i.Executable,
		WorkingDirectory:  i.WorkDir,
		ProgramArguments:  []string{},
		KeepAlive:         true,
		RunAtLoad:         true,
		StandardOutPath:   i.StdLogFile,
		StandardErrorPath: i.ErrLogFile,
	}

	if !os.IsNotExist(err) {
		current, err := cmd.parse(path)
		if err != nil {
			return err
		}

		pl.WorkingDirectory = current.WorkingDirectory
		pl.ProgramArguments = current.ProgramArguments
		pl.KeepAlive = current.KeepAlive
		pl.RunAtLoad = current.RunAtLoad
		pl.StandardOutPath = current.StandardOutPath
		pl.StandardErrorPath = current.StandardErrorPath
	}

	return cmd.daemonize(path, pl)
}

func (cmd *Daemonize) parse(path string) (*plist, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	p := plist{}
	decoder := xpath.NewDecoder(f)
	err = decoder.Decode(&p)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (cmd *Daemonize) daemonize(path string, p interface{}) error {
	fmt.Printf("   ... creating '%s'\n", path)
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	defer f.Close()

	encoder := xpath.NewEncoder(f)
	if err = encoder.Encode(p); err != nil {
		return err
	}

	return nil
}

func (cmd *Daemonize) mkdirs() error {
	directories := []string{
		cmd.workdir,
		cmd.logdir,
	}

	for _, dir := range directories {
		fmt.Printf("   ... creating '%s'\n", dir)

		if err := os.MkdirAll(dir, 0644); err != nil {
			return err
		}
	}

	return nil
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

func (cmd *Daemonize) logrotate(i *info) error {
	pid := filepath.Join(cmd.workdir, fmt.Sprintf("%s.pid", SERVICE))
	logfiles := []struct {
		LogFile string
		PID     string
	}{
		{
			LogFile: i.StdLogFile,
			PID:     pid,
		},
		{
			LogFile: i.ErrLogFile,
			PID:     pid,
		},
	}

	t := template.Must(template.New("logrotate.conf").Parse(newsyslog))
	path := filepath.Join("/etc/newsyslog.d", fmt.Sprintf("%s.conf", SERVICE))

	fmt.Printf("   ... creating '%s'\n", path)

	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	defer f.Close()

	return t.Execute(f, logfiles)
}

func (cmd *Daemonize) firewall(i *info) error {
	fmt.Println()
	fmt.Println("   ***")
	fmt.Printf("   *** WARNING: adding '%s' to the application firewall and unblocking incoming connections\n", SERVICE)
	fmt.Println("   ***")
	fmt.Println()

	path := i.Executable

	command := exec.Command("/usr/libexec/ApplicationFirewall/socketfilterfw", "--getglobalstate")
	out, err := command.CombinedOutput()
	fmt.Printf("   > %s", out)
	if err != nil {
		return fmt.Errorf("Failed to retrieve application firewall global state (%v)", err)
	}

	if strings.Contains(string(out), "State = 1") {
		command = exec.Command("/usr/libexec/ApplicationFirewall/socketfilterfw", "--setglobalstate", "off")
		out, err = command.CombinedOutput()
		fmt.Printf("   > %s", out)
		if err != nil {
			return fmt.Errorf("Failed to disable the application firewall (%v)", err)
		}

		command = exec.Command("/usr/libexec/ApplicationFirewall/socketfilterfw", "--add", path)
		out, err = command.CombinedOutput()
		fmt.Printf("   > %s", out)
		if err != nil {
			return fmt.Errorf("Failed to add 'uhppoted-rest' to the application firewall (%v)", err)
		}

		command = exec.Command("/usr/libexec/ApplicationFirewall/socketfilterfw", "--unblockapp", path)
		out, err = command.CombinedOutput()
		fmt.Printf("   > %s", out)
		if err != nil {
			return fmt.Errorf("Failed to unblock 'uhppoted-rest' on the application firewall (%v)", err)
		}

		command = exec.Command("/usr/libexec/ApplicationFirewall/socketfilterfw", "--setglobalstate", "on")
		out, err = command.CombinedOutput()
		fmt.Printf("   > %s", out)
		if err != nil {
			return fmt.Errorf("Failed to re-enable the application firewall (%v)", err)
		}

		fmt.Println()
	}

	return nil
}

func (cmd *Daemonize) genkeys(i *info) error {
	return genkeys(filepath.Dir(cmd.config), cmd.hotp)
}
