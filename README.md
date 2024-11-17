# Hot-Coffee Application

This is a coffee shop management system designed to provide a RESTful API for managing orders, menu items, and inventory. Built using **Go**, this application implements a **three-layered software architecture** focused on **clean code**, **maintainability**, and **scalability**.

## Architecture Overview

The application is structured into three main layers:

1. **Presentation Layer (Handlers)** - Manages HTTP requests and responses.
2. **Business Logic Layer (Services)** - Contains core functionality and business rules.
3. **Data Access Layer (Repositories)** - Handles data storage and retrieval from local JSON files.

Each layer interacts through well-defined interfaces to ensure decoupling, and to support future testing and extension.

## Key Features

- **RESTful API** to manage:
  - Orders
  - Menu items
  - Inventory
  - Aggregations (e.g., total sales, popular menu items)
- **Data Persistence**:
  - Data is stored locally in separate JSON files for each entity (`orders.json`, `menu_items.json`, `inventory.json`).
  - JSON files are located in a designated `data/` directory.
- **Logging**:
  - The application uses Go's `log/slog` package for logging significant events such as requests, errors, and business logic processing.
- **Error Handling**:
  - Errors are handled gracefully with appropriate HTTP status codes, and all actions are logged.

## Functional Requirements

### 1. Presentation Layer (Handlers)

- **Responsibilities**:
  - Handle HTTP requests and responses.
  - Parse JSON input and format JSON output.
  - Invoke appropriate methods from the Business Logic Layer.
- **Implementation**:
  - Handlers are organized by entities such as `order_handler.go`, `menu_handler.go`, `inventory_handler.go`.
  - Use Go's `net/http` package to set up routes and handle HTTP requests.
  - Validate input data and return meaningful error messages where necessary.

### 2. Business Logic Layer (Services)

- **Responsibilities**:
  - Implement core business logic (e.g., handling orders, managing menu items).
  - Perform computations such as aggregations (e.g., total sales, popular menu items).
  - Interact with repositories to read/write data.
- **Implementation**:
  - Define interfaces for each service (e.g., `OrderService`, `MenuService`, `InventoryService`).
  - Implement methods for core functionality (e.g., `GetTotalSales`, `GetPopularMenuItems`).
  - Ensure services are independent and can be tested in isolation.

### 3. Data Access Layer (Repositories)

- **Responsibilities**:
  - Manage data storage and retrieval in local JSON files.
  - Ensure data integrity and consistency.
  - Provide interfaces for repositories for flexibility and decoupling.
- **Implementation**:
  - Define repository interfaces for each entity (e.g., `OrderRepository`, `MenuRepository`, `InventoryRepository`).
  - Implement repositories to read from and write to JSON files in the `data/` directory.

## JSON File Structure

The data for orders, menu items, and inventory is stored in separate JSON files located in the `data/` directory:

- **orders.json**: Stores order details.
- **menu_items.json**: Stores information about menu items.
- **inventory.json**: Stores inventory levels and details.

```json
// orders.json
[
  {
    "id": 1,
    "menuItemID": 2,
    "quantity": 3,
    "totalPrice": 15.00,
    "status": "completed",
    "createdAt": "2024-11-10T12:30:00Z"
  }
]
```

```json
// menu_items.json
[
  {
    "id": 1,
    "name": "Espresso",
    "price": 5.00
  }
]
```

```json
// inventory.json
[
  {
    "id": 1,
    "name": "Espresso Beans",
    "quantity": 100
  }
]
```

## Endpoints

### Orders

- **POST /orders** - Create a new order.
- **GET /orders/{id}** - Retrieve an order by ID.
- **PUT /orders/{id}** - Update an existing order.
- **DELETE /orders/{id}** - Delete an order.
- **POST /orders/[id}close** - Closed the order.


### Menu Items

- **POST /menu-items** - Add a new menu item.
- **GET /menu-items/{id}** - Retrieve a menu item by ID.
- **PUT /menu-items/{id}** - Update a menu item.
- **DELETE /menu-items/{id}** - Delete a menu item.

### Inventory

- **POST /inventory** - Add an item to inventory.
- **GET /inventory/{id}** - Retrieve inventory information by item ID.
- **PUT /inventory/{id}** - Update inventory details.
- **DELETE /inventory/{id}** - Delete an inventory item.

### Aggregations

- **GET /aggregations/total-sales** - Get total sales based on all orders.
- **GET /aggregations/popular-menu-items** - Get a list of popular menu items based on order frequency.

## Usage

1. **Clone the repository:**

   ```bash
   git clone https://github.com/msuyundy/hot-coffee.git
   cd hot-coffee
   ```

2. **Install dependencies (Go modules):**

   ```bash
   go mod tidy
   ```

3. **Run the application:**

   The server will listen on a configurable port, defaulting to port 8080.

   ```bash
   go run main.go
   ```

4. **Test the API endpoints** using a tool like [Postman](https://www.postman.com/) or `curl`.

   Example request to create an order:

   ```bash
   curl -X POST -H "Content-Type: application/json" \
   -d '{"menuItemID": 1, "quantity": 2}' \
   http://localhost:8080/orders
   ```

## Configuration

You can configure the application to run on a different port by setting the `PORT` environment variable:

```bash
go run main.go --port 8080
```

## Logging

The application uses the `log/slog` package for logging. All requests and significant actions (e.g., creation, updates, deletions) will be logged to the console, and any errors encountered will be logged with appropriate error messages.

Example log entry:

```text
INFO: Order created successfully: {OrderID: 123, TotalPrice: 15.00}
ERROR: Failed to read menu item data: File not found
```

## Contributing

Feel free to fork this project and submit pull requests. Contributions are welcome, especially those that improve functionality, scalability, or error handling.

---

**Happy Brewing! ☕**