#Base: #Applet & {
	name:       string
	entrypoint: "\(name)"
	image:      "repo/\(name)"
	work_dir:   string | *"\(environ.PWD)"
	volumes:    [...string] | *[
		"\(environ.HOME):\(environ.HOME)",
	]
	environment: [
		for k, v in environ if k != "HOME" && k != "PWD" && k != "TMPDIR" {
			"\(k)=\(v)"
		},
	]
}

applets: [Name=_]: #Base & {
	name: Name
}

applets: {
	karpenter: {}
	kubectl: {}
	kustomize: {}
}

ignore: [applets.karpenter]
