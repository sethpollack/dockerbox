# DockerBox

`dockerbox` is a single executable that runs docker containers based on local and remote run configurations.

## Install

```
$ go get github.com/sethpollack/dockerbox
```

> By default dockerbox configuration files live in `$HOME/.dockerbox/` and binaries are installed at `$HOME/.dockerbox/bin/`. To override these defaults you can set the following environment variables `DOCKERBOX_ROOT_DIR` and `DOCKERBOX_INSTALL_DIR`.


Add `dockerbox` to your path:

```
export PATH=$HOME/.dockerbox/bin/:$PATH
```

Add a remote repo to your registry:

```
$ dockerbox registry add example https://raw.githubusercontent.com/sethpollack/dockerbox/master/example/example.yaml
```

Create a local repo:

```
cat <<'EOF' >$HOME/.dockerbox/k8s.yaml
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
dockerbox registry add k8s $HOME/.dockerbox/k8s.yaml
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
- `rm` bool (defaults true)
- `tty` bool (defaults true)
- `interactive` bool (defaults true)
- `privileged` bool
- `detach` bool
- `kill` bool (kills running container with the same name before running)
- `environment` list
- `volumes` list
- `ports` list
- `env_file` list
- `dependencies` list (list of applets to run first)
- `links` list
- `image` string
- `image_tag` string
- `command` list

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
