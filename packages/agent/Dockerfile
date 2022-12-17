FROM rust:1.66.0-slim-buster as builder
WORKDIR /app

RUN rustup target add x86_64-unknown-linux-musl

COPY . .
RUN cargo build --locked --release --target=x86_64-unknown-linux-musl


FROM alpine:3.17.0 as runner
WORKDIR /workspace

RUN apk update && apk add python3 nodejs bash openrc util-linux curl sudo gcc libc-dev ca-certificates

COPY --from=builder /app/target/x86_64-unknown-linux-musl/release/agent /usr/local/bin/agent
RUN chmod +x /usr/local/bin/agent

COPY openrc-agent.sh /etc/init.d/agent
RUN chmod +x /etc/init.d/agent

RUN rc-update add agent default
RUN rc-update add devfs boot
RUN rc-update add procfs boot
RUN rc-update add sysfs 
RUN echo "root:root" | chpasswd
