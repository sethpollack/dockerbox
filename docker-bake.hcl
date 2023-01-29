variable "DESTDIR" {
  default = ".bin"
}

variable "VERSION" {
  default = "2.0.0"
}

group "default" {
  targets = ["binaries"]
}

target "binaries" {
  args = {
    VERSION = VERSION
  }
  output = ["${DESTDIR}/build"]
  platforms = ["local"]
}

target "binaries-cross" {
  inherits = ["binaries"]
  platforms = [
    "darwin/amd64",
    "linux/amd64",
  ]
}
