variable "DESTDIR" {
  default = ".bin"
}

variable "VERSION" {
  default = "2.0.4"
}

group "default" {
  targets = ["binaries-cross"]
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
    "darwin/arm64",
    "linux/amd64",
    "linux/arm64",
  ]
}
