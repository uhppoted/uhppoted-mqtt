FROM alpine

EXPOSE 8080

RUN mkdir -p /etc/uhppoted
RUN mkdir -p /etc/uhppoted/mqtt
RUN mkdir -p /etc/uhppoted/mqtt/rsa
RUN mkdir -p /etc/uhppoted/mqtt/rsa/signing
RUN mkdir -p /etc/uhppoted/mqtt/rsa/encryption
RUN mkdir -p /var/uhppoted

COPY uhppoted.conf /etc/uhppoted/

WORKDIR /opt/uhppoted 
COPY uhppoted-mqtt .

ENTRYPOINT /opt/uhppoted/uhppoted-mqtt --debug --console
