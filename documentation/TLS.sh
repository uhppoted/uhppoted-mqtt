#/bin/bash

rm -f  localhost.jks
rm -rf docker
rm -rf etc

mkdir -p ./tmp

keytool -genkey -keyalg RSA -alias hivemq -keystore ./tmp/localhost.jks \
        -storepass hivemq -keypass hivemq \
        -validity 365 \
        -keysize 2048 \
        -ext "SAN=DNS:localhost,IP:192.168.1.100" \
        -dname "CN=localhost, OU=hivemq, O=uhppoted, L=docker, ST=docker, C=MQ"

keytool -exportcert -keystore ./tmp/localhost.jks -alias hivemq -keypass hivemq -storepass hivemq -rfc -file ./tmp/localhost.pem

openssl x509 -in ./tmp/localhost.pem -noout -text

mkdir -p ./docker
mkdir -p ./docker/hivemq
mkdir -p ./docker/uhppoted-mqtt
mkdir -p ./docker/integration-tests/hivemq
mkdir -p ./docker/integration-tests/mqttd

cp ./tmp/localhost.jks ./docker/hivemq/localhost.jks
cp ./tmp/localhost.pem ./docker/hivemq/localhost.pem
cp ./tmp/localhost.jks ./docker/integration-tests/hivemq/localhost.jks
cp ./tmp/localhost.pem ./docker/integration-tests/hivemq/localhost.pem
cp ./tmp/localhost.pem ./docker/integration-tests/mqttd/localhost.pem
cp ./tmp/localhost.pem ./docker/uhppoted-mqtt/hivemq.pem

mkdir -p ./etc/com.github.uhppoted
cp ./tmp/localhost.pem ./etc/com.github.uhppoted/hivemq.pem