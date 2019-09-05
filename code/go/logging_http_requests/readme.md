This is a sample application that implements a version of Go HTTP server logging as described in https://onepage.nopub.io/p/Logging-HTTP-requests-in-Go-233de7fe59a747078b35b82a1b035d36

License: public domain, you can use it for free, without attribution.

You can just use it your own HTTP servers, use it selected parts of the implementation or use it as a template and tweak the implementation for your own needs,

Here's what happens when run this app with `go run .`:

- an HTTP server starts up
- it starts logging to a file in `${HOME}/data/myserver/http_log/<log-file>.txt`
- we make a few http requests to demonstrate logging
- we print the logs file to stdout
- we exit
