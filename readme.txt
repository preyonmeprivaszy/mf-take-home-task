useful commands:
start the docker service: sudo systemctl start docker
wipe a pre existing db: sudo docker-compose down -v
generate the db: sudo docker-compose up -d mfoxdb
rebuild the whole project for docker: sudo docker-compose build
list all products: sudo docker-compose run --rm app list
buy something: sudo docker-compose run --rm app increase --id req-001 --sku M-PSH-IP16PR-BK-000 --qty 10 --reason "incrementos of produktos"
get a report: sudo docker-compose run --rm app report --top 3 --low-stock 20
