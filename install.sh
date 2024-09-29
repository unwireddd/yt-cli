#!/usr/bin/env bash

echo "Starting installation..."
sudo mkdir /etc/yt-cli
sudo cp channels.go /etc/yt-cli
sudo cp framework.go /etc/yt-cli
sudo cp main.go /etc/yt-cli
sudo cp go.mod /etc/yt-cli
sudo cp go.sum /etc/yt-cli
echo Directory yt-cli created successfuly
sudo chmod +x yt-cli
sudo cp yt-cli /usr/local/bin
echo Program installed succesfully, you can now edit your subscribed channels in /etc/yt-cli/channels.go
echo Before you run the program for the first time make sure you have the mpv, yt-dlp and youtube-dl packages installed
