FROM openresty/openresty:alpine

# Install additional packages
RUN apk add --no-cache \
    curl \
    git

# Install LuaRocks packages
# Install lua-resty-http directly (without LuaRocks)
RUN cd /tmp && \
    git clone https://github.com/ledgetech/lua-resty-http.git && \
    cd lua-resty-http && \
    mkdir -p /usr/local/openresty/lualib/resty && \
    cp lib/resty/http*.lua /usr/local/openresty/lualib/resty/


EXPOSE 80 443

CMD ["/usr/local/openresty/bin/openresty", "-g", "daemon off;"]
