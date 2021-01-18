# deSECupdater
Simple tool written in GO.

## About
The project's idea is to automatically update the A-entries for a specific domain at [deSEC.io](https://desec.io).

Main use case is in environments with changing IP addresses. This way a domain will always point to the right IPv4 address (as long as the execution of this tool is automated).
The IPv4 address is automatically pulled from the public API [ipify](https://www.ipify.org).

## Parameters
Parameter | Expected Input
------------ | -------------
-domain | This the domain you want to change (e.g. "domain.one")
-token | The access token you get from deSEC under "Token Management" (is 28 chars long)