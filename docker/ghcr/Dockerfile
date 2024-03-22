FROM alpine

EXPOSE 8080

RUN mkdir -p /var/uhppoted/etc
RUN mkdir -p /var/uhppoted/etc/mqtt
RUN mkdir -p /var/uhppoted/etc/mqtt/rsa
RUN mkdir -p /var/uhppoted/etc/mqtt/rsa/signing
RUN mkdir -p /var/uhppoted/etc/mqtt/rsa/encryption
RUN mkdir -p /var/uhppoted/var

COPY uhppoted.conf /var/uhppoted/etc/

WORKDIR /opt/uhppoted 
COPY uhppoted-mqtt .

ENTRYPOINT /opt/uhppoted/uhppoted-mqtt --debug --config /var/uhppoted/etc/uhppoted.conf --console
