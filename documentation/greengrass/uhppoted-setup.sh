#!/bin/bash -x

# Tim Irvin
# (c) 2022 NetTempo, Inc.
# MIT License

# The following are the parameters of the deployment passed in using the IoT recipe
WORKDIR=$1; shift
CONFIGDIR=$1; shift
ARTIFACTDIR=$1; shift
THING=$1; shift

# If the THING name wasn't passed in with the parameters, use the name in the effectiveConfig.yaml file.
# If we are using a "virtual" THING name to get around the limitation that Core devices can't have their own
# credentials, onlyu THINGS that aren't Core devices, that name is passed in as a parameter.
if [ -z "${THING}" ]; then
  THING=$(grep thingName: ${CONFIGDIR}/effectiveConfig.yaml | awk -F': ' '{print $2}' | sed 's/"//g')
fi

REGION=$(grep awsRegion: ${CONFIGDIR}/effectiveConfig.yaml | awk -F': ' '{print $2}' | sed 's/"//g')

# Our recipe brings in a copy of uhppote-cli, and this will install it in the PATH if it doesn't already exist there.
if ! which uhppote-cli 2>&1 > /dev/null; then
  cp ${ARTIFACTDIR}/uhppote-cli /usr/local/bin/uhppote-cli && sudo chmod +x /usr/local/bin/uhppote-cli
fi

# Install any missing needed tools
if ! which curl 2>&1 > /dev/null; then
  apt-get -y update && \
  apt-get -i install curl
fi
if ! which jq 2>&1 > /dev/null; then
  apt-get -y update && \
  apt-get -i install jq
fi
if ! which aws 2>&1 > /dev/null; then
  curl "https://awscli.amazonaws.com/awscli-exe-linux-$(uname -m).zip" -o "awscliv2.zip" && \
  unzip awscliv2.zip && \
  ./aws/install
  rm -rf awscliv2.zip aws
fi

# Get the certificates for this THING using IoT Discovery
curl --cert ${ARTIFACTDIR}/client-certificate.pem.crt \
     --key  ${ARTIFACTDIR}/client-private.pem.key \
      https://greengrass-ats.iot.${REGION}.amazonaws.com:8443/greengrass/discover/thing/${THING} > ${WORKDIR}/discovery || exit 1
HOST=$(jq -r '.GGGroups[0].Cores[0].Connectivity[0].HostAddress' < ${WORKDIR}/discovery)
PORT=$(jq -r '.GGGroups[0].Cores[0].Connectivity[0].PortNumber' < ${WORKDIR}/discovery)
cat ${WORKDIR}/discovery | jq -r '.GGGroups[0].CAs[0]' > ${WORKDIR}/broker-certificate.pem
rm -f ${WORKDIR}/discovery

# Take the base uhppoted.conf file that was provided by the IoT recipe and add the discovered certs, host and port.
cp ${ARTIFACTDIR}/uhppoted.conf ${WORKDIR}/uhppoted.conf
cat >> ${WORKDIR}/uhppoted.conf <<EOF
mqtt.connection.broker = tls://${HOST}:${PORT}
mqtt.connection.client.ID = ${THING}
EOF

# Copy the client certs to the working directory to make them easier to find
cp ${ARTIFACTDIR}/client-certificate.pem.crt ${WORKDIR}/client-certificate.pem.crt
chown ggc_user:ggc_group ${WORKDIR}/client-certificate.pem.crt
cp ${ARTIFACTDIR}/client-private.pem.key ${WORKDIR}/client-private.pem.key
chown ggc_user:ggc_group ${WORKDIR}/client-private.pem.key

# Install other files provided by the recipe to the working directory
cp ${ARTIFACTDIR}/cards ${WORKDIR}/cards
chown ggc_user:ggc_group ${WORKDIR}/cards

# Use the uhppote-cli devices discovery to find all the local UHPPOTE devices and add them to the conf file
# and configure those devices to set this THING as the listener (using the closest interface to that device)
IFS=$'\n'
for dev in $(uhppote-cli get-devices); do
  IFS=" " read -r devserial devip devmask devmac devvers devverdate <<< $dev
  myip=$(ip -j route get $devip | jq -r '.[0].prefsrc')
  printf "UT0311-L0x.%s.address = %s:60000" $devserial $devip >> ${WORKDIR}/uhppoted.conf
  uhppote-cli set-listener $devserial $myip:60001
done
IFS=$' '

# Clean up any orphaned *.pid and *.lock files
if [ -f ${WORKDIR}/uhppoted-mqtt.pid ]; then
  PID=$(cat ${WORKDIR}/uhppoted-mqtt.pid)
  if ! kill -0 ${PID} 2> /dev/null; then
    rm -f ${WORKDIR}/uhppoted-mqtt.pid
    rm -f ${WORKDIR}/${THING}.lock
  fi
else
  rm -f ${WORKDIR}/${THING}.lock
fi

exit 0
