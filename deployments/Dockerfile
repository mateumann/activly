FROM golang:1.21.6-alpine3.18
WORKDIR /app
RUN adduser -h /app -s /dev/null -G users -D -u 5149 activly  # FIXME: this is for Tilt, change to > 10000
COPY --chown=activly:users build/activly ./activly
#ADD web web

ENTRYPOINT ["activly"]