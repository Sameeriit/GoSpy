<h1 align="center">GoSpy</h1>

<p align="center">
    <img height=200 width=200 src="./icon.png"/>
    <br/>
    <a href="https://github.com/psidex/GoSpy/actions" >
        <img src="https://github.com/psidex/GoSpy/workflows/go%20build%20windows/badge.svg" />
    </a>
    <a href="https://github.com/psidex/GoSpy/actions" >
        <img src="https://github.com/psidex/GoSpy/workflows/go%20build%20ubuntu/badge.svg" />
    </a>
    <br/>
    <a href="https://goreportcard.com/report/github.com/psidex/GoSpy" >
        <img src="https://goreportcard.com/badge/github.com/psidex/GoSpy" />
    </a>
    <a href="./LICENSE" >
        <img src="https://img.shields.io/github/license/psidex/GoSpy" />
    </a>
    <a href="https://ko-fi.com/M4M18XB1" >
        <img src="https://img.shields.io/badge/support%20me-Ko--fi-orange.svg?style=flat&colorA=35383d" />
    </a>
</p>

<p align="center">A cross-platform remote access tool</p>

## Disclaimer

This project should be used for authorized testing or educational purposes only.

It is the final user's responsibility to obey all applicable local, state, and federal laws.

Authors assume no liability and are not responsible for any misuse or damage caused by this program.

## Usage

GoSpy consists of 2 binaries, the client is what you execute on your target machine and the server is what you run on
your machine to interact with the client.

## Features

These are almost all currently a WIP

- [x] Cross-platform with `CGO_ENABLED=0` (compiles to any target that Go supports)
- [x] Safe error handling so the client / server won't suddenly drop on error
- [x] Automatic reconnect for both client and server
- [x] Reverse shell
- [x] File grab (send a file from the client to the server)
- [x] File drop (send a file from the server to the client)
- [ ] Execute Lua scripts on target machine (using [gopher-lua](https://github.com/yuin/gopher-lua))
  - Useful if you have managed to execute the client on your target but (for whatever reason) the reverse shell can't
  execute things / isn't working
- [ ] More?
  - SSL/TLS?

## Screenshot

![](./demo.png)

## Why?

I wrote this project to learn more about both Go and penetration testing, as I recently completed an "Ethical Hacking"
unit for my university course and am interested in learning more.

## Architecture

The client maintains a main connection to the server, nicknamed `CmdCon`. This is only used to exchange commands and
arguments.

Any other time data needs to be transferred, a new connection is initiated (e.g. when sending a file). This means
that if anything goes wrong (e.g. a file read/write error) then the connection can just be closed instead of having
to deal with complicated communication logic (e.g. letting the client know an error ocurred when it's trying to send
file data).

## Credits

- [gopherize.me](https://gopherize.me/) for the icon
- [c-bata/go-prompt](https://github.com/c-bata/go-prompt/) for the interactive prompt on the server
- [vfedoroff/go-netcat](https://github.com/vfedoroff/go-netcat/blob/master/main.go) for some reverse shell net logic
