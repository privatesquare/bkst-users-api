FROM alpine:3.15

RUN mkdir /workspace

ARG NAME="privatesquare/bkst-users-api"
ARG BINARY_NAME="bkst-users-api"
ARG DESCRIPTION="Book store users api"
ARG VERSION=${VERSION}
ARG RUN_AS_USER="apiuser"
ARG APP_PATH="/appl"

ENV PATH="$PATH:/appl"

LABEL \
label-schema.maintainer="Private Square <privatesquare99@gmail.com>" \
label-schema.name=${NAME} \
label-schema.description=${DESCRIPTION} \
label-schema.vendor="Private Square" \
label-schema.version=${VERSION} \
label-schema.schema-version="1.0"

RUN addgroup -S ${RUN_AS_USER} && adduser -S ${RUN_AS_USER} -G ${RUN_AS_USER} && mkdir ${APP_PATH}

WORKDIR ${APP_PATH}

COPY ${BINARY_NAME}-linux-amd64 ${BINARY_NAME}
COPY config.env config.env

RUN chmod +x bkst-users-api \
 && chown -R ${RUN_AS_USER}:${RUN_AS_USER} ${APP_PATH}

USER ${RUN_AS_USER}

CMD ["bkst-users-api"]
