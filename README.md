# Private Theatre Management Application

This application provides an API for managing private theatres, focusing on two main aspects:

1. Users booking theatres
2. Admins managing theatres and other resources

## Features

### General Features (For All Users)

- Theatre browsing and details viewing
- Slot availability checking
- Addon browsing
- Order creation and management (booking theatres)
- User authentication and authorization
- Payment processing with Razorpay integration

### Admin-Specific Features

- Theatre management (creation and modification)
- Slot management (creating new time slots)
- Addon management (creating and modifying addons)
- User management (creating new user accounts)
- Access to all orders and bookings

## API Endpoints

### Health Check

- `GET /healthz`: Check the health status of the API

### Slots

- `POST /slots`: Create a new slot (Admin only)
- `GET /slots`: Retrieve all slots

### Theatres

- `POST /theatres`: Create a new theatre (Admin only)
- `GET /theatres`: Retrieve all theatres
- `GET /theatres/{id}`: Get details of a specific theatre

### Addons

- `POST /addons`: Create a new addon (Admin only)
- `GET /addons`: Retrieve all addons
- `GET /addons/categories`: Get addon categories

### Orders

- `POST /orders`: Create a new order
- `GET /orders`: Retrieve all orders
- `GET /orders/{orderId}`: Get details of a specific order

### Users

- `POST /users`: Create a new user (Admin only)
- `POST /login`: User login
- `POST /refresh-token`: Refresh authentication token

### Payments

- `POST /verify-payment`: Verify payment status

## Technologies Used

- Go (Golang)
- Chi router
- Razorpay for payment processing

## Setup and Installation

1. Clone the repository
2. Install dependencies: `go mod tidy`
3. Set up environment variables:
   - Copy the `.env.example` file to `.env`
   - Fill in the required values in the `.env` file
4. Run the application:
   - For development with auto-reload: `air`
   - For production: `go run main.go`

## Environment Variables

This project uses environment variables for configuration. Please refer to the `.env.example` file in the repository for the required variables. Make sure to set these up before running the application.

Key variables include database connection details, Razorpay API keys, JWT configuration, and server settings. Ensure all variables are properly set, especially sensitive information like database credentials and API keys.

## Development

This project uses [Air](https://github.com/cosmtrek/air) for live reloading during development. To use Air:

1. Install Air: `go install github.com/cosmtrek/air@latest`
2. Run the application with: `air`

Air will watch for file changes and automatically rebuild and restart the application.

## Authentication

The application uses token-based authentication. Some endpoints require admin authorization, which is handled by the `AdminAuthorization` middleware. This ensures that only authorized admin users can access and modify sensitive data or perform administrative actions.

## License

This project is licensed under the [MIT License](LICENSE).
