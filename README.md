# pinotes
Save notes from desktop/mobile browser address bar.   

It is highly recommended that you set this up on a Raspberry PI or on system that stays online 24x7. Do configure your firewall so that these notes are always served within your LAN and not Internet.

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

## View Notes
You can always visit your server url http://raspberrypi.local:8008/ to list and view saved notes. Note that these are viewed in plain text and not markdown.
