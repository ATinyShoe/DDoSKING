FROM 98a2693c996c2294358552f48373498d
ARG DEBIAN_FRONTEND=noninteractive

# Update apt sources and install required packages (combined in one layer)
RUN apt-get update && \
    apt-get install -y --no-install-recommends iperf3 nginx && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Create directory and copy file
RUN mkdir -p /var/www/html/download/
COPY test.txt /root/
RUN cp /root/test.txt /var/www/html/download/

# Configure Nginx
RUN echo 'server { \
    listen 80; \
    server_name _; \
    location /download { \
        alias /var/www/html/download/; \
        default_type application/octet-stream; \
        add_header Content-Disposition "attachment; filename=test.txt"; \
    } \
}' > /etc/nginx/sites-available/default

# Integrate with Seedemu
COPY 082b96ec819c95ae773daebde675ef80 /start.sh
COPY d18858afc6bb66ec3a19d872077acfd2 /seedemu_sniffer
COPY d3d51fdf7f4bad30dc5db560a01ce629 /interface_setup
COPY 17ac2d812a99a91e7f747e1defb72a29 /seedemu_worker
COPY e01e36443f9f72c6204189260d0bd276 /ifinfo.txt
RUN chmod +x /start.sh
RUN chmod +x /seedemu_sniffer
RUN chmod +x /seedemu_worker

# Ensure Nginx starts with container
RUN echo "service nginx start" >> /start.sh

CMD ["/start.sh"]
