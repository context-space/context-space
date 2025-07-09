#!/bin/bash

# Pass all arguments to docker compose logs command
COMPOSE_PROJECT_NAME=${COMPOSE_PROJECT_NAME:-context-space} docker compose -f docker/docker-compose.yml logs "$@"
