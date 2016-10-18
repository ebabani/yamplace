package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"

	"github.com/concourse/fly/template"
	flags "github.com/jessevdk/go-flags"
	"github.com/pkg/errors"
)

var opts struct {
	YamlFile string   `short:"y" long:"yaml" description:"YAML file with the placeholders" required:"true"`
	VarsFrom []string `short:"f" long:"vars-from" description:"Load placeholders from this file"`
}

func init() {
	args, err := flags.Parse(&opts)
	if err != nil {
		log.Printf("&v\n", args)
		os.Exit(1)
	}
}

func main() {
	configFile, err := ioutil.ReadFile(opts.YamlFile)
	if err != nil {
		log.Println(errors.Wrap(err, "Unable to read yaml file").Error())
		os.Exit(1)
	}

	var resultVars template.Variables
	for _, path := range opts.VarsFrom {
		fileVars, err := template.LoadVariablesFromFile(path)
		if err != nil {
			log.Println(errors.Wrap(err, "Unable to read plaheholder from "+path).Error())
			os.Exit(1)
		}

		resultVars = resultVars.Merge(fileVars)
	}

	configFile, err = template.Evaluate(configFile, resultVars)
	if err != nil {
		log.Println(errors.Wrap(err, "Unable to replace placeholders").Error())
		os.Exit(1)
	}

	var configStructure interface{}

	err = yaml.Unmarshal(configFile, &configStructure)
	if err != nil {
		fmt.Println(errors.Wrap(err, "Unable to replace placeholders").Error())
		os.Exit(1)
	}

	parsedYaml, _ := yaml.Marshal(configStructure)
	fmt.Println("---")
	fmt.Println(string(parsedYaml))
}
