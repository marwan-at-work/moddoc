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

## Demo

[<img width="717" alt="Screen Shot 2019-03-22 at 1 32 36 AM" src="https://user-images.githubusercontent.com/16294261/54802943-d3b6c080-4c43-11e9-8886-a294e8ed8daa.png">](https://vimeo.com/325806835)

## Status

WIP
