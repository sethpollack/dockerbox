package tools

applets: [Name=_]: {
  applet_name: Name
}

volumes: [Name=_]: {
  name: Name
}

networks: [Name=_]: {
  name: Name
}

baseEnvs: [
  for k, v in environ if k != "HOME" && k != "PATH" && k != "TMPDIR" {
  	"\(k)=\(v)"
  },
]

baseVolumes: ["\(environ.PWD):/src"]

#base: #Applet & {
  applet_name: string
  entrypoint:  string | *"\(applet_name)"
  work_dir:    string | *"/src"
  environment: [...string] | *baseEnvs
  volumes:     [...string] | *baseVolumes
}

#ruby: #base & {
  image:   "ruby"
  volumes: ["bundle:/usr/local/bundle"] + baseVolumes
}

#node: #base & {
  image:   "node"
  volumes: [
  	"yarn:/root/.yarn",
  	"node:/usr/local/lib/node_modules",
  ] + baseVolumes
  ports: ["3000:3000", "5173:5173"]
}

#rust: #base & {
  image: "rust"
  volumes: ["cargo:/root/.cargo"] + baseVolumes
  ports: ["8000:8000"]
}


applets: {
  rust:    #rust & {}
  cargo:   #rust & {}

  ruby:    #ruby & {}
  rspec:   #ruby & {}
  rubocop: #ruby & {}
  bundle:  #ruby & {}

  node:    #node & {}
  yarn:    #node & {}
  npx:     #node & {}
  npm:     #node & {}
  pnpm:    #node & {}
  tsc:     #node & {}

  ...
}

volumes: {
  yarn: {}
  bundle: {}
  node: {}
  cargo: {}

  ...
}


ignore: {
  pnpm: {}

  ...
}
