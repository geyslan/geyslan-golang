#!/bin/bash

file=test.access.log
# There are 4 valid (3 GET, 1 POST) and 2 invalid (\x16\x03\x01) log lines
lines=$(cat <<EOT
209.17.97.82 - - [11/May/2020:02:19:28 +0000] "GET / HTTP/1.1" 301 - "-" "Mozilla/5.0 (compatible; Nimbostratus-Bot/v1.3.2; http://cloudsystemnetworks.com)"
44.224.22.196 - - [12/May/2020:01:52:33 +0000] "\x16\x03\x01" 400 226 "-" "-"
209.17.97.106 - - [13/May/2020:02:34:07 +0000] "GET / HTTP/1.0" 301 - "-" "Mozilla/5.0 (compatible; Nimbostratus-Bot/v1.3.2; http://cloudsystemnetworks.com)"
44.224.22.196 - - [14/May/2020:01:52:33 +0000] "\x16\x03\x01" 400 226 "-" "-"
172.104.108.109 - - [15/May/2020:03:15:31 +0000] "GET / HTTP/1.1" 301 - "-" "Mozilla/5.0"
201.184.225.146 - - [16/May/2020:03:24:35 +0000] "POST /doLogin HTTP/1.1" 404 196 "-" "-"
EOT
)

printf "\n* Only proceed if the Watch logic is already sleeping!\n"
printf "* Because I'll append some forged log lines to $file.\n\n"
read -p "* Press Enter to continue..."
printf "\n"

printf "$lines"
echo "$lines" >> "$file"

printf "\n\n* Now wait the Watch wake up again and go back consume the REST API.\n"
printf "* See you!\n\n"