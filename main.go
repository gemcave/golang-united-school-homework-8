package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

type Arguments map[string]string

type User struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

func parseArgs() Arguments {
	id := flag.String("id", "", "user ID")
	item := flag.String("item", "", "json")
	operation := flag.String("operation", "", "method")
	fileName := flag.String("fileName", "", "file name")

	flag.Parse()

	return Arguments{
		"id":        *id,
		"item":      *item,
		"operation": *operation,
		"fileName":  *fileName,
	}
}

const ITEM_FLAG = "item"
const ID_FLAG = "id"
const FILENAME_FLAG = "fileName"
const OPERATION_FLAG = "operation"
const LIST_FLAG = "list"
const ADD_OP = "add"
const LIST_OP = "list"
const FIND_BY_ID_OP = "findById"
const REMOVE_OP = "remove"

func getUsers(fileName string, w io.Writer) {
	f, _ := ioutil.ReadFile(fileName)

	w.Write(f)
}

func writeToFile(f *os.File, b *[]byte) {
	f.Seek(0, io.SeekStart)
	f.Truncate(0)
	f.Write(*b)
}

func addUser(fileName, item string, w io.Writer) {
	f, _ := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0755)
	defer f.Close()

	b, _ := ioutil.ReadAll(f)

	users := make([]User, 0, 1)
	if len(b) > 0 {
		_ = json.Unmarshal(b, &users)
	}

	var newUser User
	json.Unmarshal([]byte(item), &newUser)
	for _, user := range users {
		if newUser.Id == user.Id {
			w.Write([]byte("Item with id " + user.Id + " already exists"))
			return
		}
	}

	users = append(users, newUser)
	jsn, _ := json.Marshal(&users)

	writeToFile(f, &jsn)
}

func removeUser(fileName, id string, w io.Writer) {
	f, _ := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0755)
	defer f.Close()

	b, _ := ioutil.ReadAll(f)
	var users []User
	if err := json.Unmarshal(b, &users); err != nil {
		return
	}

	newUsers := make([]User, 0, len(users))
	for _, user := range users {
		if user.Id != id {
			newUsers = append(newUsers, user)
		}
	}

	if len(users) == len(newUsers) {
		w.Write([]byte("Item with id " + id + " not found"))
		return
	}

	jsn, _ := json.Marshal(newUsers)
	writeToFile(f, &jsn)
}

func findById(fileName, id string, w io.Writer) {
	var users []User
	b, _ := ioutil.ReadFile(fileName)
	json.Unmarshal(b, &users)

	for _, user := range users {
		if user.Id == id {
			b, _ = json.Marshal(user)
			w.Write(b)
		}
	}
}

func Perform(args Arguments, writer io.Writer) error {
	operation := args[OPERATION_FLAG]
	if operation == "" {
		return fmt.Errorf("-operation flag has to be specified")
	}

	fileName := args[FILENAME_FLAG]
	if fileName == "" {
		return fmt.Errorf("-fileName flag has to be specified")
	}

	switch operation {
	case ADD_OP:
		item := args[ITEM_FLAG]
		if item == "" {
			return fmt.Errorf("-item flag has to be specified")
		}

		addUser(fileName, item, writer)
	case LIST_OP:
		getUsers(fileName, writer)
	case FIND_BY_ID_OP:
		id := args[ID_FLAG]
		if id == "" {
			return fmt.Errorf("-id flag has to be specified")
		}

		findById(fileName, id, writer)
	case REMOVE_OP:
		id := args[ID_FLAG]
		if id == "" {
			return fmt.Errorf("-id flag has to be specified")
		}

		removeUser(fileName, id, writer)
	default:
		return fmt.Errorf("Operation %s not allowed!", operation)
	}

	return nil
}

func main() {
	err := Perform(parseArgs(), os.Stdout)
	if err != nil {
		panic(err)
	}
}
