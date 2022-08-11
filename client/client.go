package main

import (

	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

var room_id = ""
func receiverMess(conn net.Conn) {

	for {

		readcon := bufio.NewReader(conn)
		mgs, err := readcon.ReadString('\n')
		if err != nil {

			break
		}
		var messages = strings.Split(mgs, ",")
		if len(messages) == 2 {

			var infos = strings.Split(messages[1] ,":")
			if len(infos) == 2 {

				m1 := strings.TrimSpace(infos[0])
				if m1 == "roomID" {

					room_id = strings.TrimSpace(infos[1])
				}
			}
			infos = strings.Split(messages[0] ,":")
			if len(infos) == 2 {

				m1 := strings.TrimSpace(infos[1])
				if m1 == "exit oke" {

					room_id = ""
				}
			}
		}
		fmt.Print(mgs)
	}
}

func command() string {

	for {

		fmt.Println("Enter your command : ")
		readCommand := bufio.NewReader(os.Stdin)
		mgs, err := readCommand.ReadString('\n')
		if err != nil {

			log.Fatal(err)
		}
		mgs = strings.TrimSpace(mgs)
		if mgs == "help" {

			fmt.Println("command 'exit' : leave the room")
			fmt.Println("command 'join' : join the room")
			fmt.Println("command 'create' : create the room")
			fmt.Println("command 'list_user' : show all the user")
			fmt.Println("command 'list_room' : show all the room")
			fmt.Println("command 'chat_with' : chat with a user")
			fmt.Println("command 'exit_command' : exit command")
			continue
		}
		if mgs == "exit_command" {

			return mgs
		}
		if mgs == "exit" {

			return "command : exit"
		}
		if mgs == "list_user" {

			return "command : list_user"
		}
		if mgs == "list_room" {

			return "command : list_room"
		}
		if mgs == "chat_with" {

			fmt.Println("Enter user_id : ")
			var user_id string
			for {

				if len(user_id) > 0 {

					break
				}
				read := bufio.NewReader(os.Stdin)
				name, err := read.ReadString('\n')
				if err != nil {

					log.Fatal(err)
				}
				user_id = strings.TrimSpace(name)
			}
			return "command : " + mgs + ", user ID : " + user_id
		}
		if mgs == "join" || mgs == "create" {

			fmt.Println("Enter your Room ID : ")
			var room_id string
			for {

				if len(room_id) > 0 {

					break
				}
				read := bufio.NewReader(os.Stdin)
				name, err := read.ReadString('\n')
				if err != nil {

					log.Fatal(err)
				}
				room_id = strings.TrimSpace(name)
			}
			return "command : " + mgs + ", room ID : " + room_id
		}
		fmt.Println("You entered the wrong command, please re-enter !!!")
	}
}

func main() {

	var (

		name string
		err error
	)
	for {

		if name != "" {

			break
		}
		fmt.Printf("Enter your name: ")
		reados := bufio.NewReader(os.Stdin)
		name, err = reados.ReadString('\n')
		if err != nil {

			log.Fatal(err)
		}
		name = strings.TrimSpace(name)
		if name == "" {

			fmt.Println("Your name must is not empty!!!")
		}
	}

	connection, err := net.Dial("tcp", "localhost:5055")
	if err != nil {

		log.Fatal(err)
	}
	connection.Write([]byte(name + "\n"))
	fmt.Printf("Welcome %s to the chatApp\n", name)
	fmt.Println("********** MESSAGES **********")
	go receiverMess(connection)
	for {

		readmgs := bufio.NewReader(os.Stdin)
		mgs, err := readmgs.ReadString('\n')
		mgs = strings.TrimSpace(mgs)
		var cmd = false
		if mgs == "command" {

			fmt.Println("!!!! Are you enter command? Yes or No?")
			read := bufio.NewReader(os.Stdin)
			oke, err := read.ReadString('\n')
			if err != nil {

				log.Fatal(err)
			}
			oke = strings.TrimSpace(oke)
			oke = strings.ToLower(oke)
			if oke == "yes" {

				mgs = command()
				cmd = true
			}
		}
		if err != nil {

			break
		}
		if len(room_id) == 0 {

			if len(mgs) < 7 || mgs[:7] != "command" {

				fmt.Println("You aren't in any room. Please enter command!!!")
				mgs = command()
				cmd = true
			}
		}
		if cmd && len(mgs) > 0 && mgs == "exit_command" {

			continue
		}
		if !cmd {

			mgs = "message : " + mgs
		}
		mgs = fmt.Sprintf("Room_id : %s, %s : %s\n", room_id, name, mgs)
		connection.Write([]byte(mgs))
	}

	connection.Close()
}