FROM --platform=${BUILDPLATFORM} golang:1.14.3-alpine AS build
LABEL maintainer="daan.seynaeve@openanalytics.eu"
WORKDIR /src
ENV CGO_ENABLED=0
COPY go.* .
RUN go mod download
COPY . .
ARG TARGETOS
ARG TARGETARCH
RUN mkdir -p /out
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /out/rdepot .

FROM scratch AS bin-unix
COPY --from=build /out/rdepot /

FROM bin-unix AS bin-linux
FROM bin-unix AS bin-darwin

FROM scratch AS bin-windows
COPY --from=build /out/rdepot /rdepot.exe

FROM bin-${TARGETOS} AS bin

FROM alpine AS image
COPY --from=build /out/rdepot /bin/rdepot
ENTRYPOINT ["/bin/rdepot"]
CMD ["--help"]

