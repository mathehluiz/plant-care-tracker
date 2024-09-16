# Plant Care Tracker

The **Plant Care Tracker** API is a comprehensive solution designed to simplify plant care management. It empowers users to efficiently track their plants' needs, set up and manage care routines, and oversee their account details with ease. Whether you're a hobbyist gardener or a plant enthusiast, this API provides the necessary tools to automate and personalize your plant care schedules.

## Features

- User authentication, authorization and management
- Plant management (create, read, update, delete)
- Care routines management (create, read, update, delete)
- Integration with Redis for caching
- JWT-based secure authentication
- Role-based access control (RBAC)
- PostgreSQL as the primary database
- Docker for containerization

## API Endpoints

### Authentication

- `POST /api/v1/login`: Login user
- `POST /api/v1/verify-code`: Verify the code sent to the user
- `POST /api/v1/refresh-token`: Refresh authentication token
- `GET /api/v1/me`: Get current user details
- `POST /api/v1/register`: Register a new user
- `POST /api/v1/verify-email`: Verify user email
- `POST /api/v1/reset-password`: Request password reset
- `POST /api/v1/reset-password/:id`: Reset password using token
- `GET /api/v1/reset-password/:id`: Check password reset status
- `PATCH /api/v1/set-active`: Set user active status
- `DELETE /api/v1/delete-user/:id`: Delete user by ID
- `POST /api/v1/change-roles`: Change user roles

### Plant Management

- `POST /api/v1/plants`: Create a new plant
- `GET /api/v1/plants/:id`: Get plant details by ID
- `GET /api/v1/plants`: Get all plants for the user
- `PATCH /api/v1/plants/:id`: Update plant details by ID
- `DELETE /api/v1/plants/:id`: Delete plant by ID

### Care Management

- `POST /api/v1/cares`: Create a new care routine
- `GET /api/v1/cares/:id`: Get care routine details by ID
- `GET /api/v1/cares/plant/:id`: Get all care routines for a specific plant
- `PATCH /api/v1/cares/:id`: Update care routine by ID
- `DELETE /api/v1/cares/:id`: Delete care routine by ID

## Getting Started

To get a local copy of the project up and running for development and testing, follow these instructions.

### Prerequisites

- Docker
- Docker Compose

### Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/mathehluiz/plant-care-tracker.git

   ```

2. Navigate to project directory

   ```bash
   cd plant-care-tracker

   ```

3. Set up environment variables: Copy the .env.example file to .env and update the values as needed.

4. Build and run the Docker containers:

   ```bash
   docker-compose up --build

   ```

5. Run database migrations
   ```bash
   docker-compose exec app migrate -path ./migrations -database "postgres://your_user:your_pass@localhost:5432/plant-care-tracker?sslmode=disable" up
   ```

### Documentation

API documentation is available on insomnia_v4.json format at `/docs`.

### Contributing

We welcome contributions from the community! If you would like to contribute, please follow the guidelines below:

1. Fork the repository.
2. Create a new branch for your feature or bugfix.
3. Write clean and readable code.
4. Submit a pull request with a detailed description of your changes.
