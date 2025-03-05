package main

import (
    "bufio"
    "encoding/csv"
    "fmt"
    "os"
    "strconv"
    "strings"
    "time"
)

type MenuItem struct {
    Name        string
    Price       float64
    Description string
    Category    string
    Stock       int
}

type Order struct {
    ItemName  string
    Quantity  int
    TotalCost float64
    OrderTime time.Time
}

type User struct {
    Username string
    Password string
    Role     string
}

var menu []MenuItem
var orders []Order
var activityLog []string
var users = []User{
    {"admin", "admin123", "admin"},
}

func main() {
    reader := bufio.NewReader(os.Stdin)
    fmt.Println("Login System")

    var user *User
    for attempts := 0; attempts < 3; attempts++ {
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

        fmt.Println("Invalid credentials. Attempts left:", 2-attempts)
    }

    if user == nil {
        fmt.Println("Too many failed attempts. Exiting...")
        return
    }

    for {
        fmt.Println("\n=== Restaurant Management System ===")
        fmt.Println("1. Manage Menu")
        fmt.Println("2. View Orders")
        fmt.Println("3. Financial Report")
        fmt.Println("4. Export Orders to CSV")
        fmt.Println("5. View Activity Logs")
        fmt.Println("6. Exit")
        fmt.Print("Choose an option: ")

        input, _ := reader.ReadString('\n')
        option, err := strconv.Atoi(strings.TrimSpace(input))
        if err != nil {
            fmt.Println("Invalid input. Please enter a number.")
            continue
        }

        switch option {
        case 1:
            manageMenu(reader)
        case 2:
            viewOrders()
        case 3:
            viewFinancialReport()
        case 4:
            exportOrdersToCSV()
        case 5:
            viewActivityLogs()
        case 6:
            fmt.Println("Exiting program. Goodbye!")
            return
        default:
            fmt.Println("Invalid option.")
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
    fmt.Println("1. Add Item")
    fmt.Println("2. Update Stock")
    fmt.Println("3. View Menu")
    fmt.Print("Choose an option: ")
    input, _ := reader.ReadString('\n')
    option, err := strconv.Atoi(strings.TrimSpace(input))
    if err != nil {
        fmt.Println("Invalid input.")
        return
    }

    switch option {
    case 1:
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
            fmt.Println("Invalid price.")
            return
        }

        fmt.Print("Enter item description: ")
        description, _ := reader.ReadString('\n')
        description = strings.TrimSpace(description)

        fmt.Print("Enter initial stock: ")
        stockInput, _ := reader.ReadString('\n')
        stock, err := strconv.Atoi(strings.TrimSpace(stockInput))
        if err != nil || stock < 0 {
            fmt.Println("Invalid stock.")
            return
        }

        menu = append(menu, MenuItem{Name: name, Price: price, Description: description, Category: category, Stock: stock})
        logActivity(fmt.Sprintf("Added menu item: %s", name))
        fmt.Println("Menu item added successfully!")
    case 2:
        viewMenu()
        fmt.Print("Enter the number of the item to update stock: ")
        input, _ := reader.ReadString('\n')
        itemIndex, err := strconv.Atoi(strings.TrimSpace(input))
        if err != nil || itemIndex < 1 || itemIndex > len(menu) {
            fmt.Println("Invalid menu item.")
            return
        }

        fmt.Print("Enter new stock quantity: ")
        stockInput, _ := reader.ReadString('\n')
        stock, err := strconv.Atoi(strings.TrimSpace(stockInput))
        if err != nil || stock < 0 {
            fmt.Println("Invalid stock.")
            return
        }

        menu[itemIndex-1].Stock = stock
        logActivity(fmt.Sprintf("Updated stock for %s to %d", menu[itemIndex-1].Name, stock))
        fmt.Println("Stock updated successfully!")
    case 3:
        viewMenu()
    default:
        fmt.Println("Invalid option.")
    }
}

func viewMenu() {
    if len(menu) == 0 {
        fmt.Println("\nNo menu items available yet.")
        return
    }

    fmt.Println("\n--- Menu List ---")
    fmt.Println("No | Name          | Price      | Stock | Description      | Category")
    fmt.Println("---|---------------|------------|-------|------------------|---------")
    for i, item := range menu {
        fmt.Printf("%2d | %-13s | %-10.2f | %-5d | %-16s | %-8s\n",
            i+1, item.Name, item.Price, item.Stock, item.Description, item.Category)
    }
}

func logActivity(message string) {
    activityLog = append(activityLog, fmt.Sprintf("%s: %s", time.Now().Format("2006-01-02 15:04:05"), message))
}

func viewActivityLogs() {
    if len(activityLog) == 0 {
        fmt.Println("\nNo activity logs yet.")
        return
    }

    fmt.Println("\n--- Activity Logs ---")
    for _, log := range activityLog {
        fmt.Println(log)
    }
}

func exportOrdersToCSV() {
    file, err := os.Create("orders.csv")
    if err != nil {
        fmt.Println("Error creating file:", err)
        return
    }
    defer file.Close()

    writer := csv.NewWriter(file)
    defer writer.Flush()

    err = writer.Write([]string{"Item Name", "Quantity", "Total Cost"})
    if err != nil {
        fmt.Println("Error writing to file:", err)
        return
    }

    for _, order := range orders {
        err = writer.Write([]string{order.ItemName, strconv.Itoa(order.Quantity), fmt.Sprintf("%.2f", order.TotalCost)})
        if err != nil {
            fmt.Println("Error writing to file:", err)
            return
        }
    }

    fmt.Println("Orders exported successfully to 'orders.csv'")
}

func viewOrders() {
    if len(orders) == 0 {
        fmt.Println("\nNo orders have been placed yet.")
        return
    }

    fmt.Println("\n--- Order List ---")
    fmt.Println("No | Item Name     | Quantity | Total Cost | Order Time")
    fmt.Println("---|---------------|----------|------------|-------------------")
    for i, order := range orders {
        fmt.Printf("%2d | %-13s | %-8d | %-10.2f | %s\n",
            i+1, order.ItemName, order.Quantity, order.TotalCost, order.OrderTime.Format("2006-01-02 15:04:05"))
    }
}

func viewFinancialReport() {
    if len(orders) == 0 {
        fmt.Println("\nNo financial data available. No orders have been placed.")
        return
    }

    totalRevenue := 0.0
    categoryRevenue := make(map[string]float64)

    for _, order := range orders {
        totalRevenue += order.TotalCost
        for _, item := range menu {
            if item.Name == order.ItemName {
                categoryRevenue[item.Category] += order.TotalCost
                break
            }
        }
    }

    fmt.Println("\n--- Financial Report ---")
    fmt.Printf("Total Revenue: Rp%.2f\n", totalRevenue)
    fmt.Println("Revenue by Category:")
    for category, revenue := range categoryRevenue {
        fmt.Printf("  %s: Rp%.2f\n", category, revenue)
    }
}
