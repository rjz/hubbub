package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	hubbubCli "github.com/rjz/hubbub/cli"
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

// exec loads a policyFile and a repoFile and applies the policy to each repo
func exec(policyFileName, reposFileName *string) {

	reposFile := fmt.Sprintf("./config/repos/%s.json", *reposFileName)
	repositories, err := hubbub.LoadRepositories(reposFile)
	if err != nil {
		fmt.Printf("Failed loading repositories '%s'\n", reposFile)
		fmt.Println(err)
		os.Exit(1)
	}

	policyFile := fmt.Sprintf("./config/policies/%s.json", *policyFileName)
	Policy, err := hubbub.LoadPolicy(policyFile)
	if err != nil {
		fmt.Printf("Failed loading policy '%s'\n", policyFile)
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
			Name:  "goals",
			Usage: "list available policy goals",
			Action: func(c *cli.Context) {
				hubbubCli.ListPolicyGoals()
			},
		},
		{
			Name:  "policies",
			Usage: "list available policies",
			Action: func(c *cli.Context) {
				hubbubCli.ListConfigFiles("./config/policies")
			},
		},
		{
			Name:  "repositories",
			Usage: "list available repositories",
			Action: func(c *cli.Context) {
				hubbubCli.ListConfigFiles("./config/repos")
			},
		},
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
