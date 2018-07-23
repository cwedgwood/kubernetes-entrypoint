# -*- conf -*-

FROM cwedgwood/bldr:0.03

# put the src into the container and build it
RUN mkdir -p /build/
WORKDIR /build/

# . (of the repo) -> github.com/stackanetes/kubernetes-entrypoint/
COPY ./ ./

ENV GOPATH=/gopath/
RUN mkdir -p /gopath/src/github.com/stackanetes/ ; ln -s /build/ /gopath/src/github.com/stackanetes/kubernetes-entrypoint
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o kubernetes-entrypoint

FROM busybox
COPY --from=0 /build/kubernetes-entrypoint /
ENV PATH=/
CMD [ "/kubernetes-entrypoint" ]
