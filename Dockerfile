FROM scratch
COPY zswchain-go-demo /usr/bin/zswchain-go-demo
ENTRYPOINT ["/usr/bin/zswchain-go-demo"]
