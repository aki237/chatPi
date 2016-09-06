package user

import (
	"crypto/sha1"
	b64 "encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"time"
)

type User struct {
	Cookie, Nick, IP string
	conn             net.Conn
}

type Users []User

var ul Users

var (
	ErrNickAlreadyTaken = errors.New("Error : Nick has already been taken")
	ErrMultipleInstance = errors.New("Error : Multiple clients per IP is not allowed in this server.")
	ErrUserNotFound     = errors.New("Error : User not found for the given cookie")
	ErrUnAuthorised     = errors.New("Error : The user is not authorised to do that")
)

func NewChat(ip string) {
	ul = make(Users, 0)
	chatPi := User{
		Nick:   "*ChatPi*",
		IP:     ip,
		Cookie: "chatpiultimatecookieat" + fmt.Sprint(time.Now().Unix()),
	}
	ul = append(ul, chatPi)
}

func NewUser(nick, password, ip string, conn net.Conn) (User, error) {
	_, err := GetUser(nick)
	if err == nil {
		return User{}, ErrNickAlreadyTaken
	}
	// for _,val := range ul {
	// 	if val.IP == ip {
	// 		return User{},ErrMultipleInstance
	// 	}
	// }
	pad := fmt.Sprint(time.Now().Unix()) + fmt.Sprint(rand.Int63()) + nick + password
	newuser := User{
		Nick:   nick,
		IP:     ip,
		Cookie: hash(pad),
		conn:   conn,
	}
	ul = append(ul, newuser)
	chatPi, _ := GetUser("*ChatPi*")
	for _, val := range ul {
		if val.Nick == chatPi.Nick || val.Nick == nick {
			continue
		}
		//chatPi.MessageTo(val.Nick,Ulistxml())
		typeMessageTo(val.Nick, Ulistxml(), "userlist")
	}
	return newuser, nil
}

//
func RemoveUser(u User) error {
	_, ok := GetUser(u.Nick)
	if ok != nil {
		return ok
	}
	// if bye.IP != u.IP {
	// 	return ErrUnAuthorised
	// }
	un := make(Users, 0)
	for _, val := range ul {
		if val.Cookie != u.Cookie {
			un = append(un, val)
		}
	}
	ul = un
	chatPi, _ := GetUser("*ChatPi*")
	for _, val := range ul {
		if val.Nick == chatPi.Nick {
			continue
		}
		typeMessageTo(val.Nick, Ulistxml(), "userlist")
	}
	return nil
}

//
func Userlist() []string {
	nl := make([]string, 0)
	for _, val := range ul {
		nl = append(nl, val.Nick)
	}
	return nl
}

//
func GetUser(username string) (User, error) {
	for _, val := range ul {
		if val.Nick == username {
			return val, nil
		}
	}
	return User{}, ErrUserNotFound
}

//
func (u *User) MessageTo(username string, msg string) error {
	if username == "*ChatPi*" {
		u.conn.Write([]byte(FormMessageXML("private", "*ChatPi*", "I'm a the world's suckiest bot. I haven't been given a brain yet.", "message")))
		return nil
	}
	reciever, err := GetUser(username)
	if err != nil {
		return err
	}
	reciever.conn.Write([]byte(FormMessageXML("private", u.Nick, msg, "message")))
	return nil
}

//
func typeMessageTo(username string, msg string, msgtype string) error {
	reciever, err := GetUser(username)
	if err != nil {
		return err
	}
	reciever.conn.Write([]byte(FormMessageXML("private", "*ChatPi*", msg, msgtype)))
	return nil
}

//
func (u *User) Broadcast(msg string) {
	for _, val := range ul {
		if val.Nick == "*ChatPi*" {
			continue
		}
		u.Shout(val.Nick, msg)
	}
}

//
func (u *User) Shout(username string, msg string) error {
	if username == "*ChatPi*" {
		u.conn.Write([]byte(FormMessageXML("private", "*ChatPi*", "I'm a the world's suckiest bot. I haven't been given a brain yet.", "message")))
		return nil
	}
	reciever, err := GetUser(username)
	if err != nil {
		return err
	}
	reciever.conn.Write([]byte(FormMessageXML("broadcast", u.Nick, msg, "message")))
	return nil
}

func hash(str string) string {
	a := sha1.Sum([]byte(str))
	return hex.EncodeToString(a[:])
}

func Ulistxml() string {
	ul := Userlist()
	ret := "<MEMBERS>"
	for _, val := range ul {
		ret += "<MEMBER>" + val + "</MEMBER>"
	}
	return ret + "</MEMBERS>"
}

func FormMessageXML(channel, from, message, tpe string) string {
	message = b64.StdEncoding.EncodeToString([]byte(message))
	msg := "<MESSAGE CHANNEL=\"" + channel + "\"><FROM>" + from + "</FROM><CONTENT TYPE=\"" + tpe + "\">" + message + "</CONTENT></MESSAGE>\n"
	return msg
}
