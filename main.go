package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/hay-kot/gonv/gonv"
	"github.com/hay-kot/gonv/pkgs/ui"
	"github.com/hay-kot/yal"
	"github.com/urfave/cli/v2"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
	builtBy = "unknown"
)

func shortCommit() string {
	if len(commit) > 7 {
		return commit[:7]
	}
	return commit
}

func main() {
	yal.Log.Level = yal.LevelDebug

	app := &cli.App{
		Version: fmt.Sprintf("%s (commit: %s)", version, shortCommit()),
		Name:    "gonv",
		Usage:   "a command-line utility for managing local environment variables for your system",
		Flags: []cli.Flag{
			&cli.PathFlag{
				Name:    "file",
				Aliases: []string{"f"},
				Usage:   "the file to load environment variables from",
				Value:   "~/.env.local",
				EnvVars: []string{"GONV_FILE"},
			},
			&cli.PathFlag{
				Name:    "config",
				Aliases: []string{"t"},
				Usage:   "the config file used to load environment variables from",
				Value:   "~/.config/gonv-config.json",
				EnvVars: []string{"GONV_CONFIG"},
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "set",
				Aliases: []string{"s"},
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "clobber",
						Usage: "overwrite existing environment variables",
					},
				},
				Usage:  "set a local environment variable and append it to your file",
				Action: set,
			},
			{
				Name:    "remove",
				Aliases: []string{"r"},
				Usage:   "remove a local environment variable from your file",
				Action:  remove,
			},
			{
				Name:    "list",
				Aliases: []string{"l"},
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "show-secrets",
						Value: false,
						Usage: "also displays the secrets",
					},
				},
				Usage:  "list all local environment variable keys",
				Action: list,
			},
			{
				Name:   "encrypt",
				Usage:  "encrypt the env file with a user provided password",
				Action: encrypt,
			},
			{
				Name:   "decrypt",
				Usage:  "decrypt the env file with a user provided password",
				Action: decrypt,
			},
			{
				Name:   "shred",
				Usage:  "shreds the env file similar to shred -n 1",
				Action: shred,
			},
			{
				Name:   "load",
				Usage:  "load the env file into the environment",
				Action: load,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		os.Exit(1)
	}
}

func set(c *cli.Context) error {
	env, err := gonv.New(c.String("file"), c.String("config"))
	if err != nil {
		return err
	}

	args := c.Args().Slice()

	if len(args) != 2 {
		yal.Error("expected exactly two arguments: key value")
		return fmt.Errorf("expected two arguments: key and value")
	}

	err = env.Set(args[0], args[1], c.Bool("clobber"))

	if err != nil {
		if errors.Is(err, gonv.ErrKeyExists) {
			yal.Error("key already exists use --clobber to overwrite")
		}
		return err
	}

	return env.Save()
}

func remove(c *cli.Context) error {
	env, err := gonv.New(c.String("file"), c.String("config"))
	if err != nil {
		return err
	}

	args := c.Args().Slice()
	env.Remove(args...)

	return env.Save()
}

func list(c *cli.Context) error {
	env, err := gonv.New(c.String("file"), c.String("config"))
	if err != nil {
		return err
	}

	table := [][]string{
		{"Key", "Value"},
	}

	showSecrets := c.Bool("show-secrets")

	for k, v := range env.Vars {
		if showSecrets {
			table = append(table, []string{k, v})
			continue
		}
		table = append(table, []string{k, "********"})
	}

	fmt.Println("\n" + ui.Table(table))
	return nil
}

func encrypt(c *cli.Context) error {
	env, err := gonv.New(c.String("file"), c.String("config"))
	if err != nil {
		return err
	}

	passphrase := strings.Join(c.Args().Slice(), " ")
	if passphrase == "" {
		yal.Error("you must provide a passphrase")
		return errors.New("expected one argument: passphrase")
	}

	return gonv.EncryptFile(env.EnvFile, passphrase)
}

func decrypt(c *cli.Context) error {
	env, err := gonv.New(c.String("file"), c.String("config"))
	if err != nil {
		return err
	}

	passphrase := strings.Join(c.Args().Slice(), " ")
	if passphrase == "" {
		yal.Error("you must provide a passphrase")
		return errors.New("expected one argument: passphrase")
	}

	return gonv.DecryptFile(env.EnvFile, passphrase)
}

func load(c *cli.Context) error {
	yal.Error("loaders not yet implemented")
	return nil
}

func shred(c *cli.Context) error {
	yal.Error("shred not yet implemented")
	return nil
}
