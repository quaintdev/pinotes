# pinotes
Save notes from desktop/mobile browser address bar. 

It is highly recommended that you set this up on a Raspberry PI or on system that stays online 24x7. We want to keep notes accessible within LAN but not Internet. Following commands disable all incoming connections first and then enable local LAN connections.

```
# on raspbian/ubuntu
 sudo ufw default deny incoming       # disables all incoming connections
 sudo ufw allow from 192.168.0.0/16   # allows connections within local LAN
```

## Setup

### Raspberry PI/Desktop
1. Install using `go get github.com/quaintdev/pinotes`
2. Create a config file 
3. ./pinotes

### Browser
1. Create a search engine using this url http://raspberrypi.local:8008/add?q=
2. Assign a keyword such as `pin`

## Taking notes
#### Browser Address Bar
Start taking notes from your browser address bar with

```
  pin grocery!rice

  pin todo!pay electricity bill

  pin bmark!htttp://news.ycombinator.com

```
Above searches will create `grocery.md`, `todo.md` & `bmark.md` in directory specified by config file.  
**7-Jun-20 update: using ! as separator in url instead of # due to conflicts with url**

#### Browser Context Menu
You can also add a context menu using addons like [this](https://addons.mozilla.org/en-US/firefox/addon/context-search-we/) to send selected text as a note to your PI server. By default, any note without ! will directly be saved to `web_notes.md` as define in config.json

Note that both browser address bar and context menu option can be used with mobiles too!

## View Notes

Visit http://raspberrypi.local:8008/ to view the saved notes.
