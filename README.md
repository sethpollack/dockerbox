# DockerBox

like busybox, just for docker.

## Install

`go get github.com/sethpollack/dockerbox`

Add `dockerbox` to your path

```
export PATH=$PATH:$HOME/.dockerbox/bin/
```

Add a `registry.yaml` file to `$HOME/.dockerbox/registry.yaml`

```yaml
configs:
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
    - $PWD:/app
    work_dir: /app
```

To update your local applet repo (`$HOME/.dockerbox/repo.yaml`) from all the configs in the registry run `dockerbox update`.

To see the list of available applets run `dockerbox list`

And to install an applet run `dockerbox install -i <applet name>` or `dockerbox install -a` to install all available applets. `dockerbox` installs applets by creating a symlink from `$HOME/.dockerbox/<applet name>` to the dockerbox binary  `$GOPATH/bin/dockerbox`.

