# Gonv

Gonv is a command line utility for manage local user environment variables.

**This project is a WIP and should be used with caution**

## Usage

```shell
NAME:
   gonv - a command-line utility for managing local environment variables for your system

USAGE:
   gonv [global options] command [command options] [arguments...]

VERSION:
   unknown (commit: )

COMMANDS:
   set, s     set a local environment variable and append it to your file
   remove, r  remove a local environment variable from your file
   list, l    list all local environment variable keys
   encrypt    encrypt the env file with a user provided password
   decrypt    decrypt the env file with a user provided password
   shred      shreds the env file similar to shred -n 1
   load       load the env file into the environment
   help, h    Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --config value, -t value  the config file used to load environment variables from (default: "~/.config/gonv-config.json") [$GONV_CONFIG]
   --file value, -f value    the file to load environment variables from (default: "~/.env.local") [$GONV_FILE]
   --help, -h                show help (default: false)
   --version, -v             print the version (default: false)
```

## Use Case

In most cases you'll have some subset of environmental variables that you'll want to manage for a specific machine. This project is designed to simply the process of managing those and cleaning up when you're finished. In your shell profile you should configure the `GONV_FILE` to the location of your local env file that you want to be managed by `gonv`. Additionally, you should set that file to be sources in your shell profile or configure a command to source that file.

**Setup by specifying a .env file to use in your dotfiles**

```shell
export GONV_FILE=~/.env.local
```

**Add A New Variable**

```shell
genv set FOO bar
```

**Remove a Variable**

```shell
genv remove FOO
```

**(wip): Load Env Variables from BitWarden**

```shell
genv load --bitwarden
```