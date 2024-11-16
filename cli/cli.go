package cli

import (
	"bufio"
	"fmt"
	"go-learn/db"
	"go-learn/db/models"
	"os"
	"strconv"
)

func Loop(productRepo db.Repository[models.Product]) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		var cmd string
		fmt.Scan(&cmd)

		switch cmd {
		case "getall":
			products, err := productRepo.GetAll()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				continue
			}
			for _, p := range products {
				fmt.Println(p)
			}
		case "get":
			fmt.Print("id: ")
			scanner.Scan()
			id, err := strconv.ParseUint(scanner.Text(), 10, 64)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Invalid id")
				continue
			}

			product, err := productRepo.GetOne("id", id)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				continue
			}

			fmt.Println(product)
		case "exit":
			fmt.Println("Bye!")
			return
		}
	}
}
