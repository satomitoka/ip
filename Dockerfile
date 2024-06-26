FROM wjqserver/caddy:latest

RUN mkdir -p /data/www
RUN mkdir -p /data/ipinfo/db
RUN mkdir -p /data/ipinfo/log
RUN wget -O /data/www/index.html https://raw.githubusercontent.com/WJQSERVER-STUDIO/ip/main/index.html
RUN wget -O /data/caddy/Caddyfile https://raw.githubusercontent.com/WJQSERVER-STUDIO/ip/main/Caddyfile
RUN wget -O /data/ipinfo/ip https://raw.githubusercontent.com/WJQSERVER-STUDIO/ip/main/ip
RUN wget -O /usr/local/bin/init.sh https://raw.githubusercontent.com/WJQSERVER-STUDIO/ip/main/init.sh
RUN chmod +x /data/ipinfo/ip
RUN chmod +x /usr/local/bin/init.sh

CMD ["/usr/local/bin/init.sh"]
