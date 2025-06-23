# Roblox Account Age
Automatically checks accounts creation date using HTTP requests

- [Usage](#usage)
- [Accounts list format](#accounts-list-format)
- [Building a binary](#building-a-binary)

## Usage
1. **Download** a binary for your OS from [Releases](https://github.com/WhateverNick7/roblox-account-age/releases/latest)
2. **Run** downloaded binary
3. **Drag and Drop** accounts list file or manually type path to the file
4. **Press Enter** to confirm the path
5. **Wait** for output

Note: Do not put too many accounts, or else Roblox will rate limit your IP

## Accounts list format
Expected accounts file format
```
Player1:password
Player2:apple
Player3:banana
```
File formats that will still work
```
Player1
Player2
Player3
```
and
```
Player1
Player2:this text will be ignored, could be used as comment xd
--- 01.01.25
adcdefg blah-blah-blah
```
Only requirement for username is text that would be detected by Roblox's API

## Building a binary
1. Clone repository and open it
2. Download [Golang](https://go.dev/) (go1.24.4+) and install it
### Building on Windows:
```bat
go build -o rbxage.exe main.go
```
### Building on Linux / MacOS:
```sh
go build -o rbxage main.go
```
