# pinotes

Self-hosted notes solution. Primarily targeting Raspberry PI but should work on
any system that can stay online 24x7. 

**Breaking change**  
5th Sept 21:
The application now uses sqlite instead of plain markdown files. If you had like to migrate 
your existing notes to sqlite database you can use the command line `-migrate` option. All the 
markdown files in current directory will be migrated to sqlite database.

## Setup

### Configure Firewall
We want to make sure that the notes are served within LAN and not on 
the Internet.

```
# on raspbian/ubuntu
 sudo ufw default deny incoming       # disables all incoming connections
 sudo ufw allow from 192.168.0.0/16   # allows connections within local LAN
```

### Install
The setup below is for raspbian. You may have to modify some steps as per your 
distribution.
1. Clone the repository to your desktop 
   `git clone https://github.com/quaintdev/pinotes.git`  

2. Create a systemd service file as below and move it to 
   `/etc/systemd/system/pinotes.service` on your Raspberry pi. 

    ```shell
    [Unit]
    Description=A self hosted notes service
    After=network.target
    
    [Service]
    User=pi
    WorkingDirectory=/home/user/pinotes
    LimitNOFILE=4096
    ExecStart=/home/user/pinotes/pinotes.bin
    Restart=always
    RestartSec=10
    StartLimitIntervalSec=0
    
    [Install]
    WantedBy=multi-user.target
    ```
   
2. Create a deployment script as below. You will have to modify it 
   for your env and Pi version. This is for Raspberry Pi 2 B.
   
    ```shell
    cd build
    export CGO_ENABLED=1
    export CC=arm-linux-gnueabi-gcc
    GOOS=linux GOARCH=arm GOARM=7 go build github.com/quaintdev/pinotes/cmd/
    mv cmd pinotes.bin
    scp pinotes.bin user@piaddress:/home/user/pinotes/
    ssh user@piaddress <<'ENDSSH'
    cd ~/user/pinotes
    sudo systemctl stop pinotes
    rm pinotes.bin
    sudo systemctl start pinotes
    ENDSSH
    ```
4. Create a config file as per your requirement. You can use config.json
   in this repository.
5. Verify your setup by visiting http://piaddress:8008/. You will see an
   empty list of topics `[]` if this is your first time.

### Browser Config
1. Create a search engine in your browser using url http://piaddress:8008/add?q=%s
2. Assign a keyword such as `pi`. You should now be able to add notes like below
   ```shell
   pi todo - buy groceries
   pi readlater - http://wikipedia.com/
   ```
3. All your notes will be saved in 'defaultTopic' defined in config.json
4. You can view any topic using http://piaddress:8008/topic/topicname
