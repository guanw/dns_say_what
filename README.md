## host ui

```
// in client root
$ npm run build
// in server root
$ go run main.go
```

## setup local dev for tls/ssl cert

```
$ openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365 -nodes \
  -subj "/CN=localhost"
```

note: for production cert, have to get a public domain name + hosting via cloud provider with https://letsencrypt.org/getting-started/

## set up go environment with vscode

add following to ~/Library/Application Support/Code/User/settings.json

```
"go.alternateTools": {
  "go": "/opt/homebrew/opt/go@1.23/bin/go"
},
"go.useLanguageServer": true
```

cmd+shift+p -> Go: Install/Update tool -> add all

after completion, restart vscode
