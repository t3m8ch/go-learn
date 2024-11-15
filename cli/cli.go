package cli

import (
	"bufio"
	"fmt"
	"go-learn/db"
	"os"
	"strconv"

	"github.com/shopspring/decimal"
)

func Loop(productRepo db.ProductRepository) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		var cmd string
		fmt.Scan(&cmd)

		switch cmd {
		case "add":
			var product db.Product

			fmt.Print("title: ")
			scanner.Scan()
			product.Title = scanner.Text()

			fmt.Print("description: ")
			scanner.Scan()
			product.Description = scanner.Text()

			fmt.Print("price: ")
			scanner.Scan()
			priceStr := scanner.Text()

			price, err := decimal.NewFromString(priceStr)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Invalid price")
				continue
			}
			product.Price = price

			p, err := productRepo.Save(product)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				continue
			}

			fmt.Println(p)
		case "getall":
			products, err := productRepo.GetAll()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				continue
			}
			for _, p := range *products {
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

			product, err := productRepo.GetById(id)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				continue
			}

			fmt.Println(product)
		}
	}
}
