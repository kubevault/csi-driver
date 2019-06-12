FROM {ARG_FROM}

RUN set -x \
  && apk add --update --no-cache ca-certificates tzdata e2fsprogs findmnt

ADD bin/{ARG_OS}_{ARG_ARCH}/{ARG_BIN} /{ARG_BIN}

ENTRYPOINT ["/{ARG_BIN}"]
