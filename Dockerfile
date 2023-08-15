FROM --platform=${BUILDPLATFORM} golang AS build
WORKDIR /go/src/github.com/sethpollack/dockerbox
ARG TARGETOS
ARG TARGETARCH
ARG VERSION
RUN --mount=target=. \
  --mount=type=cache,target=/root/.cache/go-build \
  --mount=type=cache,target=/go/pkg/mod \
  GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -ldflags \
  "-X github.com/sethpollack/dockerbox/version.Version=${VERSION} -X github.com/sethpollack/dockerbox/version.Commit=`git log -n 1 --pretty=format:"%h"`" \
  -o /build/dockerbox_v${VERSION}_${TARGETOS}_${TARGETARCH}

FROM scratch AS bin
COPY --from=build /build/ ./
