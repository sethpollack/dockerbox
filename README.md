# DockerBox

`dockerbox` is a single executable that runs docker containers based on local and remote run configurations.

## Install

```
$ go get github.com/sethpollack/dockerbox
```

> By default dockerbox configuration files live in `$HOME/.dockerbox/` and binaries are installed at `$HOME/.dockerbox/bin/`. To override these defaults you can set the following environment variables `DOCKERBOX_ROOT_DIR` and `DOCKERBOX_INSTALL_DIR`.


Add `dockerbox` to your path:

```
$ export PATH=$HOME/.dockerbox/bin/:$PATH
```

Add a remote repo to your registry:

```
$ dockerbox registry add example https://raw.githubusercontent.com/sethpollack/dockerbox/master/example/example.yaml
```

Create a local repo:

```
$ cat <<'EOF' >$HOME/.dockerbox/k8s.yaml
applets:
  kubectl:
    name: kubectl
    image: sethpollack/kubectl
    image_tag: 1.8.4
    entrypoint: kubectl
    environment:
    - KUBECONFIG=/root/.kube/config
    volumes:
    - $HOME/.kube:/root/.kube
    - $PWD:/app
    work_dir: /app
EOF
```

Add it to your registry:

```
$ dockerbox registry add k8s $HOME/.dockerbox/k8s.yaml
```

Update your applet cache (`$HOME/.dockerbox/.cache.yaml`) from all the repos in the registry:

```
$ dockerbox update
```

Check which applets are available:

```
$ dockerbox list
kubectl:1.8.4
terraform:0.11.8
```

Install an applet with `dockerbox install -i <applet name>` or `dockerbox install -a` to install all available applets.

```
$ dockerbox install -i kubectl
$ which kubectl
/Users/seth/.dockerbox/bin/kubectl
```

`dockerbox` installs applets by creating a symlink from `$HOME/.dockerbox/bin/<applet name>` to the dockerbox binary.

```
$ ls -l $HOME/.dockerbox/bin
kubectl -> /Users/seth/go/bin/dockerbox
terraform -> /Users/seth/go/bin/dockerbox
```

Full applet spec:

- `name` string
- `work_dir` string
- `entrypoint` string
- `restart` string
- `network` string
- `hostname` string
- `rm` bool (defaults true)
- `tty` bool (defaults true)
- `interactive` bool (defaults true)
- `privileged` bool
- `detach` bool
- `kill` bool (kills running container with the same name before running)
- `pull` bool (pulls image before running)
- `dns` list
- `dns_search` list
- `dns_option` list
- `environment` list
- `volumes` list
- `ports` list
- `env_file` list
- `all_envs` bool (loads all local environment variables)
- `env_filter` string (regex to filter environment variables when using `all_envs`)
- `dependencies` (deprecated - use `before_hooks` instead)
- `before_hooks` list (list of applets to run first)
- `after_hooks` list (list of applets to run after)
- `links` list
- `image` string
- `image_tag` string
- `command` list

You can also overide an applets settings at runtime with flags followed by a seperator. The default seperator is `--` and can be configured with the `DOCKERBOX_SEPARATOR` environment variable.

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
  install     install docker applet
  list        list all available applets in the repo
  registry
  uninstall   uninstall docker applet
  update      update the repo from the registry configs
  version
```
```
Usage:
  dockerbox registry [command]

Available Commands:
  add         Add or update a repo in the registry.
  remove      Remove a repo from the registry
```
