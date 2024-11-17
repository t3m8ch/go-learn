package cli

import (
	"bufio"
	"fmt"
	"go-learn/db"
	"go-learn/db/models"
	"os"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"
)

func Loop(productRepo db.Repository[models.Product]) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		var cmd string
		fmt.Scan(&cmd)

		switch cmd {
		case "add":
			fmt.Print("n = ")
			scanner.Scan()
			n, err := strconv.Atoi(scanner.Text())

			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				continue
			}

			products := make([]models.Product, 0, n)
			for i := 0; i < n; i++ {
				fmt.Printf("n: %d\n", i+1)

				var product models.Product

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

				products = append(products, product)
			}

			ps, err := productRepo.Add(products...)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				continue
			}

			for _, p := range ps {
				fmt.Println(p)
			}
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
		case "getfirst":
			product, err := productRepo.GetOneSql(func(cols []string) string {
				return fmt.Sprintf("SELECT %s FROM products ORDER BY title ASC", strings.Join(cols, ","))
			})
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				continue
			}

			fmt.Println(product)
		case "getlast":
			product, err := productRepo.GetOneSql(func(cols []string) string {
				return fmt.Sprintf("SELECT %s FROM products ORDER BY title DESC", strings.Join(cols, ","))
			})
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				continue
			}

			fmt.Println(product)
		case "getodd":
			products, err := productRepo.GetManySql(func(cols []string) string {
				return fmt.Sprintf("SELECT %s FROM products WHERE id %% 2 = 1", strings.Join(cols, ","))
			})
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				continue
			}
			for _, p := range products {
				fmt.Println(p)
			}
		case "exit":
			fmt.Println("Bye!")
			return
		}
	}
}
