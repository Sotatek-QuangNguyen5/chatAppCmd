package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strconv"
	"strings"
)

type Client struct {
	
	con      net.Conn
	userName string
	userId   int
}

type Room struct {
	room_id   string
	isPrivate bool
	conns     []Client
}

type Message struct {
	mgs       string
	isCommand bool
	command   string
	room_id   string
	Client    Client
}

var (
	rooms       []Room
	clients     []Client
	clientCh    = make(chan Client)
	mgsCh       = make(chan Message)
	closeClient = make(chan Client)
)

func checkIDClient(id int, clients []Client) bool {

	for _, client := range clients {

		if id == client.userId {

			return true
		}
	}
	return false
}

func checkIDRoom(id string, rooms []Room) bool {

	for _, room := range rooms {

		if id == room.room_id {

			return true
		}
	}
	return false
}

func makeClientId(clients []Client) int {

	var newId = rand.Intn(10000)
	for {

		if checkIDClient(newId, clients) {

			newId = rand.Intn(10000)
		} else {

			return newId
		}
	}
}

func makeRoomId() string {

	var newId = strconv.Itoa(rand.Intn(100000))
	for {

		if checkIDRoom(newId, rooms) {

			newId = strconv.Itoa(rand.Intn(100000))
		} else {

			break
		}
	}
	return newId
}

func findRoomIndex(room_id string, rooms []Room) int {

	for index, room := range rooms {

		if room.room_id == room_id {

			return index
		}
	}
	return -1
}

func findClient(id int, clients []Client) int {

	for index := range clients {

		if clients[index].userId == id {

			return index
		}
	}
	return -1
}

func AddClient(conn net.Conn, clients *([]Client)) Client {

	readcon := bufio.NewReader(conn)
	mgs, err := readcon.ReadString('\n')
	if err != nil {

		log.Fatal(err)
	}
	var newClient = Client{

		userName: mgs[:len(mgs)-1],
		userId:   makeClientId(*clients),
		con:      conn,
	}
	fmt.Printf("Welcome %s to the room!!!\n", mgs[:len(mgs)-1])
	*clients = append(*clients, newClient)
	return newClient
}

func decryption(mgs string, client Client) Message {

	infos := strings.Split(mgs, ",")
	var (
		isCom   = false
		mess    string
		command string
	)
	if strings.TrimSpace(strings.Split(infos[1], ":")[1]) == "command" {

		isCom = true
		command = strings.TrimSpace(strings.Split(infos[1], ":")[2])
		if len(infos) > 2 {

			mess = strings.TrimSpace(strings.Split(infos[2], ":")[1])
		}
	} else {

		mess = strings.TrimSpace(strings.Split(infos[1], ":")[2])
	}
	var message = Message{

		mgs:       mess,
		isCommand: isCom,
		command:   command,
		room_id:   strings.TrimSpace(strings.Split(infos[0], ":")[1]),
		Client:    client,
	}
	return message
}

func leaveRoom(client Client, room_id string) {

	index := findRoomIndex(room_id, rooms)
	if index < 0 {

		return
	}
	for i, cl := range rooms[index].conns {

		if cl.userId == client.userId {

			rooms[index].conns = append(rooms[index].conns[:i], rooms[index].conns[i+1:]...)
			break
		}
	}
	if len(rooms[index].conns) == 0 {

		//rooms = append(rooms[:index], rooms[index+1:]...)
	}
}

func exitRoom(mess Message) {

	if mess.mgs == "" {

		mess.Client.con.Write([]byte("You are not in any room!!!\n"))
		return
	}
	leaveRoom(mess.Client, mess.room_id)
	index := findRoomIndex(mess.room_id, rooms)
	if index < 0 {

		return
	}
	for _, cl := range rooms[index].conns {

		mgs := fmt.Sprintf("%s left the room :((\n", mess.Client.userName)
		cl.con.Write([]byte(mgs))
	}
}

func joinRoom(mess Message) {

	if mess.mgs == mess.room_id {

		mess.Client.con.Write([]byte("You are in this room!!!\n"))
		return
	}
	indexRoom := findRoomIndex(mess.mgs, rooms)
	if indexRoom < 0 {

		mess.Client.con.Write([]byte("This Room has ID : " + mess.mgs + "is not survive !!!\n"))
		return
	}
	if rooms[indexRoom].isPrivate {

		mess.Client.con.Write([]byte("This room is private, cannot join!!!\n"))
	} else {

		exitRoom(mess)
		for _, cl := range rooms[indexRoom].conns {

			mgs := fmt.Sprintf("%s joined the room <3\n", mess.Client.userName)
			cl.con.Write([]byte(mgs))
		}
		rooms[indexRoom].conns = append(rooms[indexRoom].conns, mess.Client)
		mess.Client.con.Write([]byte("You joined the room, roomID : " + mess.mgs + "\n"))
	}
}

func createRoom(mess Message) {

	if checkIDRoom(mess.mgs, rooms) {

		mess.Client.con.Write([]byte("This room ID already exist!!!\n"))
	} else {

		exitRoom(mess)
		rooms = append(rooms, Room{

			room_id:   mess.mgs,
			isPrivate: false,
			conns:     []Client{mess.Client},
		})
		mgs := fmt.Sprintf("Create room success, roomID : %s\n", mess.mgs)
		mess.Client.con.Write([]byte(mgs))
	}
}

func showListUser(mess Message) {

	mgs := "List user is active :\t"
	for index := 0; index < len(clients); index++ {

		fmt.Println(clients[index], index, mess.Client.userId)
		if clients[index].userId != mess.Client.userId {

			mgs = mgs + fmt.Sprintf("UserName : %s, UserId : %d\t\t\t", clients[index].userName, clients[index].userId)
			//mess.Client.con.Write([]byte(mgs))
		}
	}
	mess.Client.con.Write([]byte(mgs + "\n"))
}

func showListRoom(mess Message) {

	mgs := "List Room :\t"
	//mess.Client.con.Write([]byte(mgs))
	for _, room := range rooms {

		if room.room_id != mess.room_id {

			mgs = fmt.Sprintf("RoomId %s, IsPrivate :%t\t\t\t", room.room_id, room.isPrivate)
			//mess.Client.con.Write([]byte(mgs))
		}
	}
	mess.Client.con.Write([]byte(mgs + "\n"))
}

func createPrivateRoom(mess Message, newClient Client) {

	exitRoom(mess)
	room_id := makeRoomId()
	rooms = append(rooms, Room{

		room_id:   room_id,
		isPrivate: true,
		conns:     []Client{mess.Client, newClient},
	})
	mgs := fmt.Sprintf("Chat with %s, roomID : %s\n", newClient.userName, room_id)
	mess.Client.con.Write([]byte(mgs))
	mgs = fmt.Sprintf("Chat with %s, roomID : %s\n", mess.Client.userName, room_id)
	newClient.con.Write([]byte(mgs))
}

func isFree(client Client) bool {

	for _, room := range rooms {

		for _, cl := range room.conns {

			if cl.userId == client.userId {

				return false
			}
		}
	}
	return true
}
func chatWith(mess Message) {

	user_id, err := strconv.Atoi(mess.mgs)
	if err != nil {

		mess.Client.con.Write([]byte("UserID is invalid!!!\n"))
		return
	}
	index := findClient(user_id, clients)
	if index < 0 {

		mess.Client.con.Write([]byte("UserID is invalid!!!\n"))
		return
	}
	if isFree(clients[index]) {

		createPrivateRoom(mess, clients[index])
	} else {

		mess.Client.con.Write([]byte("User " + clients[index].userName + " is in another room!!!\n"))
	}
}

func send(mess Message) {

	if mess.isCommand {

		if mess.command == "exit" {

			exitRoom(mess)
		} else if mess.command == "join" {

			joinRoom(mess)
		} else if mess.command == "create" {

			createRoom(mess)
		} else if mess.command == "list_user" {

			showListUser(mess)
		} else if mess.command == "list_room" {

			showListRoom(mess)
		} else if mess.command == "chat_with" {

			chatWith(mess)
		}
	} else {

		index := findRoomIndex(mess.room_id, rooms)
		if index < 0 {

			return
		}
		for _, client := range rooms[index].conns {

			if client.userId != mess.Client.userId {

				mgs := fmt.Sprintf("%s : %s\n", mess.Client.userName, mess.mgs)
				client.con.Write([]byte(mgs))
			}
		}
	}
}

func onMessage(client Client) {

	for {

		readcon := bufio.NewReader(client.con)
		mgs, err := readcon.ReadString('\n')
		if err != nil {

			break
		}
		mess := decryption(mgs[:len(mgs)-1], client)
		mgsCh <- mess
	}
	closeClient <- client
}

func deleteClient(client Client) {

	for index, cl := range clients {

		if client.userId == cl.userId {

			clients = append(clients[:index], clients[index+1:]...)
			break
		}
	}
	for i := range rooms {

		for j, cl := range rooms[i].conns {

			if cl.userId == client.userId {

				rooms[i].conns = append(rooms[i].conns[:j], rooms[i].conns[j+1:]...)
				break
			}
		}
	}
}

func main() {

	server, err := net.Listen("tcp", "localhost:5055")
	if err != nil {

		log.Fatal(err)
	}
	fmt.Println("Server is running on port 5055!!!")
	go func() {

		for {

			conn, err := server.Accept()
			if err != nil {

				log.Fatal(err)
			}

			clientCh <- AddClient(conn, &clients)
			fmt.Println(clients)
		}
	}()

	for {

		select {

		case conn := <-clientCh:
			go onMessage(conn)
		case mgs := <-mgsCh:
			go send(mgs)
		case client := <-closeClient:
			deleteClient(client)

		}
	}

}
