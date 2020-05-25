# pinotes
Save notes from desktop/mobile browser to raspberry pi

## Setup

### Raspberry PI/Desktop
1. Install using `go get github.com/quaintdev/pinotes`
2. Create a config file 
3. ./pinotes

### Browser
1. Create a search engine using this url http://raspberrypi.local:8008/add?q=
2. Assign a keyword such as `ank`

Start taking notes from your browser with

```
  ank grocery#rice

  ank todo#pay electricity bill

  ank bmark#htttp://news.ycombinator.com
```

Above searches will create `grocery.md`, `todo.md` & `bmark.md` in directory specified by config file.
