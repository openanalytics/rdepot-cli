FROM --platform=${BUILDPLATFORM} golang:1.14.3-alpine AS build
LABEL maintainer="daan.seynaeve@openanalytics.eu"
WORKDIR /src
ENV CGO_ENABLED=0
COPY go.* .
RUN go mod download
COPY . .
ARG TARGETOS
ARG TARGETARCH
RUN mkdir -p /src/out
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /src/out/rdepot .

FROM alpine AS bin-unix
COPY --from=build /src/out/rdepot /bin/rdepot
ENTRYPOINT ["/bin/rdepot"]
CMD ["--help"]

