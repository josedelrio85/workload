# Workload tool

## Start mysql docker db

```bash
docker-compose -f dev/mysql/docker-compose.yml up

# or

docker-compose -f dev/mysql/docker-compose.yml start # if the instance was already created

# connect to mysql
mysql -h 127.0.0.1 -u root -p
```

## Project

* Get all projects

```bash
curl -X GET http://localhost:9001/workload/projects
```

* Update a project

```bash
curl -X POST -H "Content-type: application/json" http://localhost:9001/workload/projects

curl -X POST -H "Content-type: application/json" --data '{"type":1, "project": [{"id": 99,"name": "Test"}]}' http://localhost:9001/workload/get
```