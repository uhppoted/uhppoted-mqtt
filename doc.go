// Copyright 2023 uhppoted@twyst.co.za. All rights reserved.
// Use of this source code is governed by an MIT-style license
// that can be found in the LICENSE file.

/*
Package uhppoted-mqtt implements an MQTT client for the UHPPOTE TCP/IP Wiegand-26 access controllers.

The MQTT client wraps the low level UDP API implemented by uhppote-core in an set of MQTT messages that
add support for authentication and authorisation as well as adding functionality to manage access control
lists and events.

The MQTT client is based on Eclipse Paho and at this point in time supports only MQTT v3.1.
*/
package mqtt
