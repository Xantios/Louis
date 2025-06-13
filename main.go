package main

import (
	"fmt"
	"github.com/goccy/go-yaml"
	ClientDrivers "github.com/xantios/louis/clients"
	"github.com/xantios/louis/clients/OPNSense"
	"github.com/xantios/louis/clients/Proxmox"
	ClientHooks "github.com/xantios/louis/hooks"
	"github.com/xantios/louis/hooks/Slack"
	"github.com/xantios/louis/hooks/Webhook"
	"os"
	"time"
)

type CfgOpnSense struct {
	URL        string `yaml:"url"`
	Username   string `yaml:"username"`
	Password   string `yaml:"password"`
	BackupPath string `yaml:"backupPath"`
}

type CfgProxmox struct {
	URL      string `yaml:"url"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type Config struct {
	Debug    bool                   `yaml:"debug"`
	Interval int                    `yaml:"interval"`
	OPNSense map[string]CfgOpnSense `yaml:"OPNSense"`
	Proxmox  map[string]CfgProxmox  `yaml:"Proxmox"`
	Hooks    map[string]string      `yaml:"hooks"`
}

var clients = make(map[string]ClientDrivers.Client)
var hooks = make(map[string]ClientHooks.Hook)

func RegisterClient(name string, client ClientDrivers.Client) {
	clients[name] = client
}

func RegisterHook(name string, hook ClientHooks.Hook) {
	hooks[name] = hook
}

func main() {

	// Pull config from disk
	cfgBytes, err := os.ReadFile("./config.yaml")
	if err != nil {
		fmt.Println("Failed to read config.yaml\n")
		fmt.Println(err)
		os.Exit(1)
	}

	var cfg Config
	err = yaml.Unmarshal(cfgBytes, &cfg)
	if err != nil {
		fmt.Println("Failed to parse config.yaml\n")
		fmt.Println(err)
		os.Exit(1)
	}

	if cfg.Debug {
		fmt.Println("Debug mode enabled")
	}

	// Register OPNSense boxes
	for name, config := range cfg.OPNSense {
		c := OPNSense.New(config.URL, config.BackupPath, config.Username, config.Password, cfg.Debug)
		RegisterClient(fmt.Sprintf("OPNSense_%s", name), c)
	}

	// Register Proxmox boxes
	for name, config := range cfg.Proxmox {
		p := Proxmox.New(name, config.URL, config.Username, config.Password)
		RegisterClient(fmt.Sprintf("Proxmox_%s", name), p)
	}

	// Register hooks
	s := Slack.New(cfg.Hooks["slack"])
	RegisterHook("Slack", s)

	w := Webhook.New(cfg.Hooks["webhook"])
	RegisterHook("Webhook", w)

	ticker := time.NewTicker(time.Minute * time.Duration(cfg.Interval))
	for {
		Run()
		fmt.Printf("Sleeping for %d minute(s)...", cfg.Interval)
		<-ticker.C
	}
}

func Dispatch(msg string) {
	for name, hook := range hooks {
		e := hook.Send(msg)
		if e != nil {
			fmt.Printf("Failed to send message to %s hook", name)
			fmt.Println(e)
		}
	}
}

func Run() {
	fmt.Printf("Running %d clients\n", len(clients))
	for name, instance := range clients {
		fmt.Println("Checking updates for ", name)
		shouldUpdate, msg, err := instance.Update()
		if err != nil {
			Dispatch(err.Error())
		}

		if shouldUpdate {
			Dispatch(msg)
		}
	}
}
