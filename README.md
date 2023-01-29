
# DockerBox

`dockerbox` is a tool that allows you to run command line tools using Docker instead of native binaries. It works by symlinking commands to the Dockerbox binary, which then looks up the `docker run` configuration and executes the `docker run command` for you.


Dockerbox uses configuration files written in [Cuelang](https://cuelang.org/) for specifying the Docker run configuration. 

The root schema for the configuration files contains the following fields:

`environ: [string]: string` - This field is automatically populated with the host's environment variables at runtime.

`applets: [string]: #Applet` - This field is used to configure the applets.

`ignore:  [string]: #Applet` - This field is used to instruct dockerbox to skip certain applets when running the dockerbox install command.


`dockerbox` will look for configuration files by walking the path of your current directory and unifying all of the files.

> File names are arbitrary, but must end in `.dbx.cue`.


## Install

```
$ go get github.com/sethpollack/dockerbox
```

> By default dockerbox binaries are installed at `$HOME/.dockerbox/bin/`. To override you can set the following environment variable `DOCKERBOX_INSTALL_DIR`.

Then add `dockerbox` to your path:

```
$ export PATH=$DOCKERBOX_INSTALL_DIR:$PATH
```

To prevent the installation of an applet, add it to the ignore list in any of your configuration files.

```
$ cat <<'EOF' >$HOME/.ignore.dbx.cue
ignore: [applets.kubectl]
EOF
```

Install applets with `dockerbox install`.

```
$ dockerbox install
$ which kubectl
/Users/seth/.dockerbox/bin/kubectl
```

`dockerbox` installs applets by creating a symlink from `$HOME/.dockerbox/bin/<applet name>` to the dockerbox binary.

```
$ ls -l $HOME/.dockerbox/bin
kubectl -> /Users/seth/go/bin/dockerbox
terraform -> /Users/seth/go/bin/dockerbox
```

Full schema can be found [here](cue/schema.cue).

You can also override an applets settings at runtime with flags followed by a separator. The default separator is `--` and can be configured with the `DOCKERBOX_SEPARATOR` environment variable.

For example:

```
applets: {
  kubectl: {
	  name: "kubectl"
	  image: "kubectl"
  }
}
```

If the user wants to run `kubectl proxy` they can either add the ports section to the config, or they can just pass an additional runtime flag which will be passed along to the docker run command.

```
kubectl -p 8080:8080 -- proxy --port=8080
```

The full flag spec can be found below:

```
      --after-hook strings    Run container after
      --all-envs              Pass all envars to container
      --before-hook strings   Run container before
      --command strings       Command to run in container
      --dependency strings    Run container before
  -d, --detach                Run container in background and print container ID
      --dns strings           Set custom DNS servers
      --dns-option strings    Set DNS options
      --dns-search strings    Set custom DNS search domains
      --entrypoint string     Overwrite the default ENTRYPOINT of the image
      --env-file strings      Read in a file of environment variables
      --env-filter string     Filter env vars passed to container from --all-envs
  -e, --environment strings   Set environment variables
      --hostname string       Container host name
      --image string          Container image
  -i, --interactive           Keep STDIN open even if not attached
      --inverse               Inverse env-filter
      --kill                  Kill previous run on container with same name
      --link strings          Add link to another container
      --name string           Assign a name to the container
      --network string        Connect a container to a network
      --privileged            Give extended privileges to this container
  -p, --publish strings       Publish a container's port(s) to the host
      --pull                  Pull image before running it
      --restart string        Restart policy to apply when a container exits
      --rm                    Automatically remove the container when it exits
      --tag string            Container image tag
  -t, --tty                   Allocate a pseudo-TTY
  -v, --volume strings        Bind mount a volume
  -w, --workdir string        Working directory inside the container
```

## Usage
```
Usage:
  dockerbox [command]

Available Commands:
  help        Help about any command
  install     Install docker applet
  uninstall   Uninstall docker applet
  version
```
