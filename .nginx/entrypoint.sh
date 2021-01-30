envsubst '$$NODE_IP $$POD_IP'< /etc/nginx/conf.d/server.conf.template > /etc/nginx/conf.d/server.conf
nginx -g 'daemon off;'