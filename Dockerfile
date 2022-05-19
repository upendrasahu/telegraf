FROM telegraf:1.22.2
COPY telegraf /usr/bin/telegraf
RUN setcap 'cap_net_raw=-ep' /usr/bin/telegraf
CMD ["telegraf", "--config", "/etc/telegraf/telegraf.conf"]
