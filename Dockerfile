FROM --platform=${BUILDPLATFORM} golang:1.17-alpine AS base
WORKDIR /miniproject2
ENV CGO_ENABLED=0
COPY . .
ARG TARGETOS
ARG TARGETARCH
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /out/example .

FROM scratch AS bin
COPY --from=base /out/example /