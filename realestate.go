package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sort"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:    "realestate",
		Version: "0.1.0",
		Usage:   "CLI Helper for Real Estate app",
		Action: func(c *cli.Context) error {
			fmt.Println("Hello friend!")
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:    "start",
				Aliases: []string{"st"},
				Usage:   "Start docker",
				Action: func(c *cli.Context) error {
					cmd := exec.Command("docker-compose", "-f", "docker-compose.yml", "up", "-d")
					stdout, err := cmd.Output()
					if err != nil {
						return err
					}
					fmt.Print(string(stdout))
					return nil
				},
			},
			{
				Name:    "stop",
				Aliases: []string{"sp"},
				Usage:   "Stop docker",
				Action: func(c *cli.Context) error {
					cmd := exec.Command("docker-compose", "-f", "docker-compose.yml", "down")
					stdout, err := cmd.Output()
					if err != nil {
						return err
					}
					fmt.Print(string(stdout))
					return nil
				},
			},
			{
				Name:    "database",
				Aliases: []string{"db"},
				Usage:   "Access database (psql) CLI",
				Action: func(c *cli.Context) error {
					cmd := exec.Command("docker-compose", "-f", "docker-compose.yml", "exec", "database", "psql", "-U", "realestate", "-d", "realestate_db")
					var out bytes.Buffer
					var stderr bytes.Buffer
					cmd.Stdout = &out
					cmd.Stderr = &stderr
					err := cmd.Run()
					if err != nil {
						fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
						return err
					}
					fmt.Println("Result: " + out.String())
					return nil
				},
			},
			{
				Name:    "migrate-create",
				Aliases: []string{"mc"},
				Usage:   "Creates up and down migration files",
				Action: func(c *cli.Context) error {
					cmd := exec.Command("migrate", "create", "-ext", "sql", "-dir", "db/migrations", "-seq", c.Args().Get(0))
					var out bytes.Buffer
					var stderr bytes.Buffer
					cmd.Stdout = &out
					cmd.Stderr = &stderr
					err := cmd.Run()
					if err != nil {
						fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
						return err
					}
					fmt.Println("Result: " + out.String())
					return nil
				},
			},
			{
				Name:    "migrate-up",
				Aliases: []string{"mu"},
				Usage:   "Run up migrations",
				Action: func(c *cli.Context) error {
					cmd := exec.Command("/bin/sh", "-c", "migrate -database postgres://realestate:realestate_pass@localhost:5432/realestate_db?sslmode=disable -path db/migrations up")
					stdout, err := cmd.Output()
					if err != nil {
						return err
					}
					fmt.Print(string(stdout))
					return nil
				},
			},
			{
				Name:    "migrate-down",
				Aliases: []string{"md"},
				Usage:   "Run down migrations",
				Action: func(c *cli.Context) error {
					cmd := exec.Command("/bin/sh", "-c", "migrate -database postgres://realestate:realestate_pass@localhost:5432/realestate_db?sslmode=disable -path db/migrations down")
					stdout, err := cmd.Output()
					if err != nil {
						return err
					}
					fmt.Print(string(stdout))
					return nil
				},
			},
			{
				Name:    "cli-server",
				Aliases: []string{"cs"},
				Usage:   "Access go container cli",
				Action: func(c *cli.Context) error {
					cmd := exec.Command("/bin/sh", "-c", "docker exec -it realestate_server /bin/sh")
					stdout, err := cmd.Output()
					if err != nil {
						return err
					}
					fmt.Print(string(stdout))
					return nil
				},
			},
			{
				Name:    "cli-database",
				Aliases: []string{"cdb"},
				Usage:   "Access database container cli",
				Action: func(c *cli.Context) error {
					cmd := exec.Command("/bin/sh", "-c", "docker exec -it realestate_database /bin/bash")
					stdout, err := cmd.Output()
					if err != nil {
						return err
					}
					fmt.Print(string(stdout))
					return nil
				},
			},
			{
				Name:    "js-bundle",
				Aliases: []string{"jsb"},
				Usage:   "Bundle JavaScript assets and watch for changes",
				Action: func(c *cli.Context) error {
					cmd := exec.Command("/bin/sh", "-c", "docker exec -it realestate_server /bin/sh -c \"npm run bundle\"")
					stdout, err := cmd.Output()
					if err != nil {
						return err
					}
					fmt.Print(string(stdout))
					return nil
				},
			},
			{
				Name:    "js-install",
				Aliases: []string{"jsi"},
				Usage:   "Install NPM package",
				Action: func(c *cli.Context) error {
					cmd := exec.Command("docker", "exec", "realestate_server", "/bin/sh", "-c", fmt.Sprintf("npm install %v", c.Args().Get(0)))
					var out bytes.Buffer
					var stderr bytes.Buffer
					cmd.Stdout = &out
					cmd.Stderr = &stderr
					err := cmd.Run()
					if err != nil {
						fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
						return err
					}
					fmt.Println("Result: " + out.String())
					return nil
				},
			},
		},
	}

	sort.Sort(cli.CommandsByName(app.Commands))

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
