# urls.ac: Distributed URL Shortener
Personal project I made to get better at distributed system design and learn Kubernetes. This is a monorepo that
contains two backend services, k8s configs, a Docker Compose setup to test locally, and a Makefile to automate 
some commands.

## Tech Stack
- FastAPI
- Redis
- InfluxDB
- MongoDB
- Docker
- Kubernetes
- Contour


## System Design
TODO

## Project Structure
```
url-shortener/
├── urls_service/
├── frontend/
├── analytics_service/
├── k8s/
├── nginx/
├── frontend/
├── docker-compose.yml
├── Dockerfile
└── README.md
```

## Running Locally
To run locally, you need to set up the `.env.local` files in the `urls_service` and the `analytics_service`.
Then, run `make compose-up`.
