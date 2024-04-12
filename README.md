# Test Server
## Simple test server for triggering specifig responses and behaviours.
- trigger specifig http response codes
- simulate slow responses
- view rate and count metrics per triggered endpoint

## Get Started

1. make sure you have go installed
2. run
   ```console
   go install github.com/devinmarder/test-server
   ```
3. run the testserver with
   ```console
   test-server -port=:8000
   ```
the port is optional and defaults to `8000`
