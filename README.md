# Usage

```
# install - or download from https://github.com/korylprince/fleetapns/releases
GOBIN=`pwd` go install github.com/korylprince/fleetapns@v0.1.0
./fleetapns -email "<your email>" -org "<your org>" -password pass
CSR submitted to Fleet's server. Check inbox at <your email> for email from Fleet.
# upload CSR from Fleet email to identity.apple.com to create/renew APNS cert
# upload cert/key to MicroMDM
mdmctl mdmcert upload -cert <cert downloaded from identity.apple.com> -private-key mdm-certificates/PushCertificatePrivateKey.key -password=pass
```
