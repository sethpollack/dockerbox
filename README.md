# DockerBox

dockerbox is a single executable that runs docker containers based on local and remote run configurations.

## Install

`go get github.com/sethpollack/dockerbox`

Add `dockerbox` to your path

```
export PATH=$PATH:$HOME/.dockerbox/bin/
```

Add a `registry.yaml` file to `$HOME/.dockerbox/registry.yaml`

```yaml
repos:
- name: local
  path: $HOME/.dockerbox/local.yaml
  type: file
- name: example
  path: https://raw.githubusercontent.com/sethpollack/dockerbox/master/example/example.yaml
  type: url
```

add a local config to `$HOME/.dockerbox/local.yaml`

```yaml
applets:
  kubectl:
    name: kubectl
    image: beenverifiedinc/kubectl
    image_tag: 1.7.11
    entrypoint: kubectl
    environment:
    - KUBECONFIG=/root/.kube/config
    volumes:
    - $HOME/.kube:/root/.kube
    - $PWD:/app
    work_dir: /app
```

To update your applet cache (`$HOME/.dockerbox/.cache.yaml`) from all the repos in the registry run `dockerbox update`.

To see the list of available applets run `dockerbox list`

And to install an applet run `dockerbox install -i <applet name>` or `dockerbox install -a` to install all available applets. `dockerbox` installs applets by creating a symlink from `$HOME/.dockerbox/<applet name>` to the dockerbox binary `$GOPATH/bin/dockerbox`.

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
  uninstall   uninstall docker applet
  update      update the repo from the registry configs
```
