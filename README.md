# Auction System - Go Routines

## Overview
This project is a Go-based auction system leveraging goroutines to automatically close auctions after a configurable time period. It integrates with MongoDB for persistent storage and provides a RESTful API for managing auctions, bids, and users.

## Features
- Automatic auction closure using Go routines
- RESTful API for auction and bid management
- MongoDB integration for data persistence
- Dockerized setup for easy deployment

## Getting Started

### Prerequisites
- Docker and Docker Compose installed
- Git (optional, for cloning the repository)

### Installation
1. Clone the repository (if not already downloaded):
   ```sh
   git clone https://github.com/BrunoGuimaraesSilva/labs-auction-goexpert.git
   cd labs-auction-goexpert
   ```
2. Start the application and MongoDB using Docker Compose:
   ```sh
   docker-compose up --build
   ```
   - The `--build` flag ensures the application image is rebuilt if changes are made.

### Accessing the Application
Once the containers are running, access the API at:
- **Base URL:** [http://localhost:8080](http://localhost:8080)

## API Endpoints

### Auction Management
- **Create an Auction**
  - **Endpoint:** `POST /auction`
  - **Request Body:**
    ```json
    {
        "product_name": "Smartphone",
        "category": "Electronics",
        "description": "iPhone 15 Pro Max",
        "condition": 1
    }
    ```
    - `condition`: Integer (e.g., 0 = Used, 1 = New)

- **List All Auctions**
  - **Endpoint:** `GET /auction`

- **Get Auction by ID**
  - **Endpoint:** `GET /auction/:auctionId`

- **Get Winning Bid**
  - **Endpoint:** `GET /auction/winner/:auctionId`

### Bid Management
- **Place a Bid**
  - **Endpoint:** `POST /bid`
  - **Request Body:**
    ```json
    {
        "user_id": "user456",
        "auction_id": "auction789",
        "amount": 150.75
    }
    ```
    - `user_id`: ID of the bidding user
    - `auction_id`: ID of the auction being bid on
    - `amount`: Bid amount (floating-point value)

- **List Bids for an Auction**
  - **Endpoint:** `GET /bid/:auctionId`

### User Management
- **Get User by ID**
  - **Endpoint:** `GET /user/:userId`

## Configuration
The application uses environment variables defined in `cmd/auction/.env`:
```
MONGO_INITDB_ROOT_USERNAME=admin
MONGO_INITDB_ROOT_PASSWORD=admin
MONGODB_URL=mongodb://admin:admin@mongodb:27017/auctions?authSource=admin
MONGODB_DB=auctions
AUCTION_INTERVAL=20s
```
- `AUCTION_INTERVAL`: Duration after which auctions close (e.g., `20s`).

## Running Tests
To execute automated tests:
1. Ensure Docker containers are running:
   ```sh
   docker-compose up -d
   ```
2. Run tests locally (assuming MongoDB is accessible):
   ```sh
   go test -v ./internal/infra/database/auction
   ```
   - Tests assume MongoDB is at `localhost:27017` with `admin:admin` credentials if run outside Docker.

### Debugging
Check container logs for troubleshooting:
```sh
docker-compose logs
```

## Repository
- **GitHub:** [labs-auction-goexpert](https://github.com/BrunoGuimaraesSilva/labs-auction-goexpert)

