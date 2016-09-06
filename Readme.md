# chatPi
*A simpe tcp chat server in golang*

## Installation

```shell
$ go get -u github.com/aki237/chatPi
```

## Prepare
You might have to change the IP and te port in the [chatPi.go](https://github.com/aki237/chatPi/blob/master/chatPi.go#L16) file.
Then inside the folder :
```shell
$ go install
```
or
```shell
$ go install github.com/aki237/chatPi
```

# Running
Just fire up the terminal and run the command `chatPi`. Say you are running the server at `192.168.0.100:6672`, `telnet` to that server and you'll be able to interact.
OR
Try the [chatterbox](https://github.com/aki237/chatterbox) GUI client.
