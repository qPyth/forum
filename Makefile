build:
	docker build -t forum .
run-img:
	docker run --network=host --name=forum -p 8080:8080 --rm -d forum
run:
	go run ./cmd
stop:
	docker stop forum