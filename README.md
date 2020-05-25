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

## Taking notes
#### Browser Address Bar
Start taking notes from your browser address bar with

```
  ank grocery#rice

  ank todo#pay electricity bill

  ank bmark#htttp://news.ycombinator.com
```

Above searches will create `grocery.md`, `todo.md` & `bmark.md` in directory specified by config file.

#### Browser Context Menu
You can also add a context menu using addons like [this](https://addons.mozilla.org/en-US/firefox/addon/context-search-we/) to send selected text as a note to your PI server. By default, any note without # will directly be saved to `web_notes.md` as define in config.json

Note that both browser address bar and context menu option can be used with mobiles too!
