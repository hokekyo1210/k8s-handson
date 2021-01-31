FROM nginx:1.19

RUN apt-get update && apt-get install -y \
    curl \
    dnsutils \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

COPY ./nginx.conf /etc/nginx/nginx.conf
COPY ./server.conf /etc/nginx/templates/server.conf.template