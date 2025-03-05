package main

import (
    "bufio"
    "encoding/csv"
    "fmt"
    "os"
    "strconv"
    "strings"
)

type MenuItem struct {
    Name        string
    Price       float64
    Description string
    Category    string
}

type Order struct {
    ItemName  string
    Quantity  int
    TotalCost float64
}

type User struct {
    Username string
    Password string
    Role     string // "admin" or "customer"
}

var menu []MenuItem
var orders []Order
var users = []User{
    {"admin", "admin123", "admin"},
    {"customer", "cust123", "customer"},
}

func main() {
    reader := bufio.NewReader(os.Stdin)
    fmt.Println("Login System")
    
    attempts := 0
    var user *User
    for attempts < 3 {
        fmt.Print("Username: ")
        username, _ := reader.ReadString('\n')
        username = strings.TrimSpace(username)
        
        fmt.Print("Password: ")
        password, _ := reader.ReadString('\n')
        password = strings.TrimSpace(password)
        
        user = authenticate(username, password)
        if user != nil {
            break
        }
        
        attempts++
        fmt.Println("Invalid credentials. Attempts left:", 3-attempts)
    }
    
    if user == nil {
        fmt.Println("Too many failed attempts. Exiting...")
        return
    }
    
    for {
        fmt.Println("\n=== Restaurant Management System ===")
        if user.Role == "admin" {
            fmt.Println("1. Manage Menu")
            fmt.Println("2. View Orders")
        } else {
            fmt.Println("1. View Menu")
            fmt.Println("2. Place Order")
        }
        fmt.Println("3. Export Orders to CSV")
        fmt.Println("4. Exit")
        fmt.Print("Choose an option: ")

        input, _ := reader.ReadString('\n')
        option, err := strconv.Atoi(strings.TrimSpace(input))
        if err != nil {
            fmt.Println("Invalid input. Please enter a number.")
            continue
        }

        if user.Role == "admin" {
            switch option {
            case 1:
                manageMenu(reader)
            case 2:
                viewOrders()
            case 3:
                exportOrdersToCSV()
            case 4:
                fmt.Println("Exiting program. Goodbye!")
                return
            default:
                fmt.Println("Invalid option.")
            }
        } else {
            switch option {
            case 1:
                viewMenu()
            case 2:
                placeOrder(reader)
            case 3:
                exportOrdersToCSV()
            case 4:
                fmt.Println("Exiting program. Goodbye!")
                return
            default:
                fmt.Println("Invalid option.")
            }
        }
    }
}

func authenticate(username, password string) *User {
    for _, u := range users {
        if u.Username == username && u.Password == password {
            return &u
        }
    }
    return nil
}

func manageMenu(reader *bufio.Reader) {
    fmt.Println("\n--- Manage Menu ---")
    fmt.Print("Enter item name: ")
    name, _ := reader.ReadString('\n')
    name = strings.TrimSpace(name)

    fmt.Print("Enter category (Food/Drink/Dessert): ")
    category, _ := reader.ReadString('\n')
    category = strings.TrimSpace(category)

    fmt.Print("Enter item price: ")
    priceInput, _ := reader.ReadString('\n')
    price, err := strconv.ParseFloat(strings.TrimSpace(priceInput), 64)
    if err != nil {
        fmt.Println("Invalid price. Please enter a valid number.")
        return
    }

    fmt.Print("Enter item description: ")
    description, _ := reader.ReadString('\n')
    description = strings.TrimSpace(description)

    menu = append(menu, MenuItem{Name: name, Price: price, Description: description, Category: category})
    fmt.Println("Menu item added successfully!")
}

func viewMenu() {
    if len(menu) == 0 {
        fmt.Println("No menu items available.")
        return
    }
    fmt.Println("\n--- Menu List ---")
    for i, item := range menu {
        fmt.Printf("%d. %s - %.2f (%s) [%s]\n", i+1, item.Name, item.Price, item.Description, item.Category)
    }
}

func placeOrder(reader *bufio.Reader) {
    viewMenu()
    fmt.Print("Enter the number of the item to order: ")
    input, _ := reader.ReadString('\n')
    itemIndex, err := strconv.Atoi(strings.TrimSpace(input))
    if err != nil || itemIndex < 1 || itemIndex > len(menu) {
        fmt.Println("Invalid selection. Please choose a valid menu item.")
        return
    }

    fmt.Print("Enter quantity: ")
    quantityInput, _ := reader.ReadString('\n')
    quantity, err := strconv.Atoi(strings.TrimSpace(quantityInput))
    if err != nil || quantity <= 0 {
        fmt.Println("Invalid quantity. Please enter a positive number.")
        return
    }

    selectedItem := menu[itemIndex-1]
    totalCost := selectedItem.Price * float64(quantity)
    if selectedItem.Category == "Drink" && totalCost > 20000 {
        fmt.Println("You get a 10% discount on drinks over Rp20,000!")
        totalCost *= 0.9
    }

    orders = append(orders, Order{ItemName: selectedItem.Name, Quantity: quantity, TotalCost: totalCost})
    fmt.Printf("Order for %s (x%d) has been placed successfully! Total: %.2f\n", selectedItem.Name, quantity, totalCost)
}

func viewOrders() {
    if len(orders) == 0 {
        fmt.Println("No orders placed.")
        return
    }
    fmt.Println("\n--- Order List ---")
    for i, order := range orders {
        fmt.Printf("%d. %s x %d - Total: %.2f\n", i+1, order.ItemName, order.Quantity, order.TotalCost)
    }
}

func exportOrdersToCSV() {
    file, err := os.Create("orders.csv")
    if err != nil {
        fmt.Println("Error creating file.")
        return
    }
    defer file.Close()

    writer := csv.NewWriter(file)
    defer writer.Flush()

    writer.Write([]string{"Item Name", "Quantity", "Total Cost"})
    for _, order := range orders {
        writer.Write([]string{order.ItemName, strconv.Itoa(order.Quantity), fmt.Sprintf("%.2f", order.TotalCost)})
    }
    fmt.Println("Orders exported to orders.csv")
}
