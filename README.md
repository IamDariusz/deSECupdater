# deSECupdater
Simple tool written in GO.

## About
The project's idea is to automatically update the A-entries for a specific domain at [deSEC.io](https://desec.io).

Main use case is in environments with changing IP addresses. This way a domain will always point to the right IPv4 address (as long as the execution of this tool is automated).
The IPv4 address is automatically pulled from the public API [ipify](https://www.ipify.org).

Developed according to [deSEC's API documentation](https://desec.readthedocs.io/).

## Parameters
Parameter | Expected Input
------------ | -------------
-domain | The domain you want to change (e.g. "domain.one")
-subdomain | The subdomain under the main domain you want to change (e.g. "mysubdomain")
-token | The access token you get from deSEC under "Token Management" (is 28 chars long)
-ip | The IPv4 address to use, if not used the script will determine the WAN's IPv4 address automatically via [ipify](https://www.ipify.org).