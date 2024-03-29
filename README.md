# Wombat Task Queue

Wombat Task Queue, an open-source task queue service built with Go and leveraging MongoDB for storage, is engineered to simplify task management and processing in distributed applications. 


## Planned Features

- **Task Prioritization**: Assign priority levels to tasks for processing.
- **Stale Task Handling**: Automatically remove tasks that have been pending for too long.
- **Task Metrics**: Track task processing times and success rates.



## Usage


### Installation

1. Clone the repository:
  ```bash
  git clone https://github.com/xis/wombat.git
  ```
2. Navigate to the project directory:
  ```bash
  cd wombat
  ```
  3. Build the project:
  ```bash
  go build -o wombat ./cmd
  ```

## Usage

### Setting Up Environment Variables

Wombat uses environment variables for configuration. Before starting the service, ensure you have set the following variables:

- `MONGO_URI`: The connection string to your MongoDB instance.
- `HTTP_ADDR`: The HTTP address and port the service will listen on (default: `:8080`).

### Starting the Service

To start the Wombat service, run:
  
  ```bash 
  ./wombat
  ```
### API Endpoints

Wombat provides the following RESTful API endpoints for managing tasks within queues:

- **GET `/queues/:queueID/tasks`**: Retrieve pending tasks from a specific queue.
- **POST `/queues/:queueID/tasks`**: Create a new task in a specific queue.
- **PUT `/queues/:queueID/tasks/:taskID`**: Update the status of a specific task.

## Development

For developing on Wombat, you will need a Go development environment set up and an accessible MongoDB instance for testing.

## Contributing

Contributions are welcome! If you'd like to contribute, please fork the repository and create a pull request with your changes. For major changes, please open an issue first to discuss what you would like to change.

## License

Distributed under the MIT License. See `LICENSE` for more information.
