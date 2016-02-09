package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	hubbub "github.com/rjz/hubbub/common"
	_ "github.com/rjz/hubbub/services"
	"os"
	"sync"
)

// environmentalFacts provides defaults from the environment
func environmentalFacts() map[string]interface{} {
	envFacts := map[string]interface{}{}

	// Github-related defaults
	envFacts["github.access_token"] = os.Getenv("HUBBUB_GITHUB_ACCESS_TOKEN")

	// Travis-related defaults
	envFacts["travis.org_token"] = os.Getenv("HUBBUB_TRAVIS_ORG_TOKEN")
	envFacts["travis.pro_token"] = os.Getenv("HUBBUB_TRAVIS_PRO_TOKEN")

	return envFacts
}

// printGoals describes all globally-registered goals
func printGoals() {
	fmt.Println("Goals:")
	serviceFactories := hubbub.ServiceFactories()
	for _, name := range serviceFactories.Goals() {
		fmt.Println("  *", name)
	}
	os.Exit(0)
}

// exec loads a policyFile and a repoFile and applies the policy to each repo
func exec(policyFile, reposFile *string) {

	repositories, err := hubbub.LoadRepositories(*reposFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	Policy, err := hubbub.LoadPolicy(*policyFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var wg sync.WaitGroup
	for _, repo := range *repositories {

		facts := hubbub.NewFacts(environmentalFacts())
		facts.SetRepository(&repo)

		wg.Add(1)
		go func(sess *hubbub.Session) {
			sess.Run()
			wg.Done()
		}(hubbub.NewSession(&Policy, facts))
	}
	wg.Wait()
}

func main() {
	app := cli.NewApp()
	app.Name = "hubbub"
	app.Usage = "apply a policy to a repository list"

	app.Commands = []cli.Command{
		{
			Name:  "apply",
			Usage: "apply policy",
			Action: func(c *cli.Context) {
				reposFile := c.String("repositories")
				policyFile := c.String("policy")
				exec(&policyFile, &reposFile)
			},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "repositories",
					Usage: "name of repository list",
				},
				cli.StringFlag{
					Name:  "policy",
					Usage: "name of policy",
				},
			},
		},
	}

	app.Run(os.Args)
}
