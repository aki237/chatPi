package user

import (
	"net"
	"fmt"
	"time"
	"errors"
	"math/rand"
	"crypto/sha1"
	"encoding/hex"
	b64 "encoding/base64"
)

type User struct {
	Cookie,Nick,IP string
	conn        net.Conn
}

type Users []User
var ul Users

var (
	ErrNickAlreadyTaken = errors.New("Error : Nick has already been taken")
	ErrMultipleInstance = errors.New("Error : Multiple clients per IP is not allowed in this server.")
	ErrUserNotFound = errors.New("Error : User not found for the given cookie")
	ErrUnAuthorised = errors.New("Error : The user is not authorised to do that")
)


func NewChat(ip string) {
	ul = make(Users,0)
	chatPi := User{
		Nick     : "*ChatPi*",
		IP       : ip,
		Cookie   : "chatpiultimatecookieat" + fmt.Sprint(time.Now().Unix()),
	}
	ul = append(ul,chatPi)
}

func NewUser (nick,password,ip string, conn net.Conn) (User,error) {
	_, err := GetUser(nick)
	if err == nil {
		return User{},ErrNickAlreadyTaken
	}
	// for _,val := range ul {
	// 	if val.IP == ip {
	// 		return User{},ErrMultipleInstance
	// 	}
	// }
	pad := fmt.Sprint(time.Now().Unix()) + fmt.Sprint(rand.Int63()) + nick + password
	newuser := User{
		Nick : nick,
		IP : ip,
		Cookie : hash(pad),
		conn : conn,
	}
	chatPi,_ := GetUser("*ChatPi*")
	for _,val := range ul {
		if val.Nick == chatPi.Nick {
			continue
		}
		chatPi.Broadcast(Ulistxml())
	}
	ul = append(ul,newuser)
	return newuser,nil
}

//
func RemoveUser (u User) (error) {
	_, ok := GetUser(u.Cookie)
	if ok != nil {
		return ok
	}
	// if bye.IP != u.IP {
	// 	return ErrUnAuthorised
	// }
	for index,_ := range ul {
		if ul[index] == u {
			ul = append(ul[:index],ul[index+1:]...)
			break
		}
	}
	return nil
}

//
func Userlist() ([]string) {
	nl := make([]string,0)
	for _,val := range ul {
		nl = append(nl,val.Nick)
	}
	return nl
}

//
func GetUser (username string) (User,error) {
	for _,val := range ul {
		if (val.Nick == username) {
			return val, nil
		}
	}
	return User{}, ErrUserNotFound
}

//
func (u *User) MessageTo (username string, msg string) (error) {
	reciever, err := GetUser(username)
	if err != nil {
		return err
	}
	reciever.conn.Write([]byte(formMessageXML(u.Nick, msg)))
	return nil
}

//
func (u *User) Broadcast (msg string) {
	for _,val := range ul {
		if val.Nick == "*ChatPi*" {
			continue
		}
		u.MessageTo(val.Nick, msg)
	}
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

func formMessageXML (from,message string) string {
	message = b64.StdEncoding.EncodeToString([]byte(message))
	msg := "<MESSAGE><FROM>"+from+"</FROM><CONTENT>"+message+"</CONTENT></MESSAGE>\n"
	return msg
}
