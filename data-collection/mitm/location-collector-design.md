# Location Collector VPN Design

The goal of the location collector VPN is to capture and decrypt requests to gs-loc.apple.com and \*.ls.apple.com which make up the location services of IOS and MacOS.

This system replaces the combination of `mitmproxy` and WireGuard for the purpose of privacy. By processing the traffic on-device, we do not provide the means for the survey to accidentally collect private data.

The implementation should generate a root certificate which can be installed onto IOS as a profile, allowing us to decrypt traffic. Then, when enabled, the VPN simply forwards traffic as is from the current network interface without a 3rd party server. Using SNI sniffing, we can detect requests to our target hosts and conduct MITM, saving the relevant data to a datastore.

This VPN applies only to IOS for which we are doing a study on location services.
