# moddoc

This is a server that takes a GOPROXY url as an argument and gives you a UI documentation similar to godoc.org except all data comes from the GOPROXY and not from VCS. 

If your GOPROXY supports a /catalog endpoint, then you can see and search the list of existing modules on the home page. 

## Quick start

```bash
~ go get marwan.io/moddoc
# Assuming you have a GOPROXY server running at http://localhost:3000
~ GOPROXY=http://localhost:3000 moddoc
```

Visit http://localhost:3001 

You can also visit `http://localhost:3001/<module>/@v/<version>`  to see a documentation package directly. 
For example, http://localhost:3001/github.com/pkg/errors/@v/v0.8.1

## Development

There is a UI server as well as a backend server. You need to run both of them. 

#### UI Server: 
Make sure you have Node/NPM and the Vue CLI installed.

```bash
cd frontend && vui ui
```

The command will open up a dashboard where you can hit `tasks` and click `run task`. 


#### Development server:

```bash
go run .
```

The UI dev server automatically proxies 


## Status

WIP