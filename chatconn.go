package main

import (
	"net"
	"strings"
	"bufio"
	"github.com/aki237/chatPi/user"
)

type chatConn struct {
	net.Conn
}

const (
	ON_CONNECT  string = `...Welcome to chatPi...
        JOIN [nick] [password]                  -  join the chat room. This returns a "cookie" string.
        LIST USERS                              -  list of all users connected at that time.
        MSG WITH [cookie] TO [nick] [message]   -  message a specific user a message. Cookie must be specified each time.
        BROADCAST WITH [cookie] [message]       -  broadcast a message to all the users in a chat room.
        OUT WITH [cookie]                       -  disconnect from the chat server.
`
	SYNERR_JOIN string = "Wrong Syntax. Try :\n\tJOIN [nick] [password]\n"
	SYNERR_LIST string = "Wrong Syntax. Try :\n\tLIST USERS\n"
	SYNERR_MSG  string = "Wrong Syntax. Try :\n\tMSG WITH [cookie] TO [nick] [message]\n"
	SYNERR_OUT  string = "Wrong Syntax. Try :\n\tOUT [cookie]\n"
	SYNERR_BCST string = "Wrong Syntax. Try :\n\tBROADCAST WITH [cookie] [message]\n"
)

//
func (c *chatConn) Serve () {
	ips := strings.Split(c.RemoteAddr().String(),":")
	if len(ips) < 2 {
		c.Write([]byte("You are an alien\n"))
		c.Close()
		return
	}
	ip := ips[0]
	u := user.User{}
	authed := false
	//c.Write([]byte(ON_CONNECT))
	
	CONNLOOP :
	for {
		issued, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			log("REMOVING USER",err)
			log(user.RemoveUser(u))
			break			
		}
		message := strings.Trim(issued, "\n\r ")
		command := strings.Split(message, " ")
		switch command[0] {

		// Join functions
		case "JOIN" :
			if len(command) != 3 {
				c.Write([]byte(user.FormMessageXML("*ChatPi*",SYNERR_JOIN,"error")))
				continue
			}
			if !authed {
				nick := command[1]
				password := command[2]
				var erro error
				u, erro = user.NewUser(nick, password, ip, c.Conn)
				switch erro {
				// case user.ErrMultipleInstance :
				// 	c.Write([]byte("Error : Multiple connections per PC is not allowed\n"))
				// 	break
				case user.ErrNickAlreadyTaken :
					c.Write([]byte(user.FormMessageXML("*ChatPi*","Username already Taken","error")))
					continue
				case nil :
					c.Write([]byte(user.FormMessageXML("*ChatPi*",u.Cookie,"cookie")))
					authed = true
					continue
				}
			} else {
				c.Write([]byte(user.FormMessageXML("*ChatPi*","You have already been registered","error")))
				continue
			}
		case "LIST" :
			if (len(command) != 2) {
				log(command ," : Wrong SYNTAX")
				c.Write([]byte(user.FormMessageXML("*ChatPi*",SYNERR_LIST,"error")))
				continue
			}
			if (command[1] != "USERS") {
				c.Write([]byte(user.FormMessageXML("*ChatPi*",SYNERR_LIST,"error")))
				continue
			}
			c.Write([]byte(user.Ulistxml()+"\n"))
			continue
		case "MSG" :
			if !authed {
				c.Write([]byte(user.FormMessageXML("*ChatPi*","You have to register first","error")))
				continue
			}
			if len(command) < 6 ||
				command[1] != "WITH" ||
				command[3] != "TO" {
				c.Write([]byte(user.FormMessageXML("*ChatPi*",SYNERR_MSG,"error")))
				continue
			}
			if u.Cookie != command[2] {
				log(u.Cookie,command[2])
				c.Write([]byte(user.FormMessageXML("*ChatPi*","Error : Cookie wrong","error")))
				continue
			}
			reciever, err := user.GetUser(command[4])
			if err != nil {
				c.Write([]byte(user.FormMessageXML("*ChatPi*","User not found. Use 'LIST USERS' to view the user list","error")))
				continue
			}
			msg := concatenate(command[5:])
			u.MessageTo(reciever.Nick, msg)
		case "BROADCAST":
			if !authed {
				c.Write([]byte(user.FormMessageXML("*ChatPi*","You have to register first.","error")))
				continue
			}
			if len(command) < 4 || command[1] != "WITH" {
				log("ERROR WITH \"WITH\"")
				c.Write([]byte(user.FormMessageXML("*ChatPi*",SYNERR_BCST,"error")))
				continue
			}
			if u.Cookie != command[2] {
				log(u.Cookie,command[2])
				c.Write([]byte(user.FormMessageXML("*ChatPi*","Error : Cookie wrong","error")))
				continue
			}
			msg := concatenate(command[3:])
			u.Broadcast(msg)
		case "OUT":
			if authed {
				c.Write([]byte(user.FormMessageXML("*ChatPi*","Bye have a great time....","message")))
				user.RemoveUser(u)
			}
			break CONNLOOP
		default:
			c.Write([]byte(user.FormMessageXML("*ChatPi*","Command not found","error")))
		}
	}
	c.Close()
}

func concatenate(msg []string) string {
	ret := ""
	for index, val := range msg {
		if index == 0 {
			ret += val
		} else {
			ret += " " + val
		}
	}
	return ret
}
