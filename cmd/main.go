package main

import (
	"fmt"
	"github.com/Foxcapades/Argonaut/v0"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"os/exec"
)

const (
	help = "Reads the configuration file 'sites.yml' and runs the test scripts " +
		"against the site/prefix combinations configured there."
	authHelp = "QA Authorization token, used for running against QA sites from" +
		" your local machine.  For further details see the help output from " +
		"param.sh or pub-strat.sh"
	vHelp = "enable debug level logging"
	testBlock = `
==================================================

    %s

==================================================

`
)

func main() {
	var auth string
	var verb bool

	cli.NewCommand().
		Description(help).
		Flag(cli.LFlag("auth", authHelp).Bind(&auth, true)).
		Flag(cli.SFlag('v', vHelp).Bind(&verb, true)).
		MustParse()

	raw, err := ioutil.ReadFile("./sites.yml")
	if err != nil {
		panic(err)
	}

	var config Config
	if err = yaml.Unmarshal(raw, &config); err != nil {
		panic(err)
	}

	for _, site := range config.Sites {
		for _, prefix := range config.Prefixes {
			var path string
			if prefix == "" {
				path = site
			} else {
				path = prefix + "." + site
			}
			fmt.Printf(testBlock, path)

			runCmd("./param.sh", path, auth, verb)
			runCmd("./pub-strat.sh", path, auth, verb)
		}
	}
}

type Config struct {
	Prefixes []string `yaml:"prefixes"`
	Sites    []string `yaml:"sites"`
}

func runCmd(com, site, auth string, verb bool) {
	cmd := exec.Command(com, "--summary=yaml")

	if len(auth) > 0 {
		cmd.Args = append(cmd.Args, "--auth="+auth)
	}

	if verb {
		cmd.Args = append(cmd.Args, "-v")
	}
	cmd.Args = append(cmd.Args, site)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		panic(err)
	}
}
