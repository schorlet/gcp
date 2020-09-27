FROM golang:1.14-buster as compiler

ARG SHORT_SHA=unset
WORKDIR /go/src/app

ENV CGO_ENABLED 0
ENV GO111MODULE on

RUN set -eux && \
	addgroup --gid 3000 plop && \
	adduser --gid 3000 --uid 3000 --no-create-home --disabled-login --gecos "" plop

COPY go.mod go.sum /go/src/app/
RUN go mod download

COPY . /go/src/app/
RUN go build \
	-v -trimpath -mod=readonly \
	-ldflags "-X 'main.Version=$SHORT_SHA'" \
	-o /go/bin/webapp

# -------------------

#FROM gcr.io/distroless/base:debug as final
FROM gcr.io/distroless/static
USER plop:plop
EXPOSE 8011

COPY --from=compiler /etc/passwd /etc/group /etc/
COPY --from=compiler --chown=plop:plop /go/bin/webapp /webapp

ENTRYPOINT ["/webapp"]
