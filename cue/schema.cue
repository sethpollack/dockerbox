#Applet: {
  applet_name: string

  name?: string
  entrypoint?: string
  image: string
  image_tag: string | *"latest"
  work_dir?: string

  restart?: string | "no" | "always" | "on-failure" | "unless-stopped"
  hostname?: string

  interactive: bool | *true
  tty: bool | *true
  rm: bool | *true

  kill?: bool
  pull?: bool
  detach?: bool
  privileged?: bool

  after_hooks?: [...#Applet]
  before_hooks?: [...#Applet]
  command?: [...string]
  dns?: [...string]
  dns_option?: [...string]
  dns_search?: [...string]
  environment?: [...string]
  env_file?: [...string]
  links?: [...string]
  ports?: [...string]
  volumes?: [...string]
  networks?: [...string]
}

#Network: {
  name: string
  driver?: string
}

#Volume: {
  name: string
  driver?: string
}

environ: [string]: string
applets: [string]: #Applet
networks: [string]: #Network
volumes: [string]: #Volume
ignore: [string]: #Applet
