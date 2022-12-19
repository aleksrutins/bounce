# Architecture
The Bounce code execution engine is designed to execute code at scale in a secure and fault-tolerant manner. It consists of four main components: the agent, the worker, the backend, and the messaging system.

# Agent

The agent is a Rust program that runs inside a Firecracker micro VM. Firecracker is a virtualization technology that allows for the creation of secure and lightweight VMs. It is responsible for executing the code, collecting the output, and returning the results to the worker.

# Worker

The worker is a Go program that consumes code submissions from the messaging system (in this case, RabbitMQ) and forwards them to the agent for execution. It also handles the results returned by the agent and makes them available to the backend.

# Backend

The backend is a TypeScript-based web server that uses the Koa.js framework. It offers a REST API for submitting code for execution and retrieving the results. It pushes new tasks to the messaging system for the worker to consume and stores code submissions and their outputs in a PostgreSQL database.

# Messaging System (RabbitMQ)

The messaging system (in this case, RabbitMQ) is used to communicate between the worker and the backend. The backend pushes new tasks to the messaging system for the worker to consume, and the worker sends the results of the code execution back to the backend through the messaging system.
