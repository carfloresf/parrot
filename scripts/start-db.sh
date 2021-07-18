docker build -t parrot-db init/.
docker run -e POSTGRES_PASSWORD=docker -d --name parrot-db --publish 5433:5432 parrot-db
