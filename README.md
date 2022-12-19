# Bounce

Bounce is a scalable and fault-tolerant code execution engine that uses [Firecracker](https://github.com/firecracker-microvm/firecracker) for running code in isolated microVMs and [RabbitMQ](https://rabbitmq.com) for scheduling them.

## Getting Started
1. Install the following dependencies:
    - [Docker](https://docs.docker.com/install/)
    - [Go](https://golang.org/doc/install)
    - [Node.js](https://nodejs.org/en/download/)

2. Clone the repository:
```bash
git clone https://github.com/wyzlle/bounce.git .
```
3. Start RabbitMQ and Postgres containers:
```bash
docker-compose up -d
```
4. Start the worker (and agent):
```bash
./packages/worker/start.sh
```
5. Start the backend server:
```bash
cd packages/backend
yarn install
yarn build
yarn start
```

## API

Bounce exposes a simple API for submitting code and retrieving the results.

### POST /exec

Submits code for execution. Accepts the following parameters in the request body:
- code: the code to be executed (required)
- language: the language the code is written in (required)
- input: input to be provided to the code during execution (optional)

### GET /exec/:id

Retrieves the results of a code execution by its submission ID. Returns a JSON object with the following properties:

- stdout: the standard output produced by the code
- stderr: the standard error produced by the code
- exitCode: the exit code of the code

### Examples
```bash
# submit code for execution
curl -X POST -H "Content-Type: application/json" -d '{"code": "print(input())", "input": "Hello World", "language": "py"}' http://localhost:3004/exec

# get execution results
curl http://localhost:3004/exec/<submission_id>
```

### Architecture

## License

Bounce is licensed under the [MIT License](https://github.com/wyzlle/bounce/blob/main/LICENSE).