FROM alpine

EXPOSE 60001/udp

RUN mkdir -p /usr/local/etc/uhppoted/mqtt
RUN mkdir -p /usr/local/etc/uhppoted/mqtt/rsa
RUN mkdir -p /usr/local/etc/uhppoted/mqtt/rsa/signing
RUN mkdir -p /usr/local/etc/uhppoted/mqtt/rsa/encryption

COPY uhppoted.conf /usr/local/etc/uhppoted

WORKDIR /opt/uhppoted 
COPY uhppoted-mqtt .

ENTRYPOINT /opt/uhppoted/uhppoted-mqtt --debug --config /usr/local/etc/uhppoted/uhppoted.conf --console
