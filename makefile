.PHONY: all consumer like queue start-all stop-all stop-consumer stop-like stop-queue


consumer:
	@echo "Starting consumer-service..."
	@cd consumer-service && nohup go run . > ../consumer.log 2>&1 & echo $$! > ../consumer.pid

like:
	@echo "Starting like-service..."
	@cd like-service && nohup go run . > ../like.log 2>&1 & echo $$! > ../like.pid

queue:
	@echo "Starting queue-service..."
	@cd queue-service && nohup go run . > ../queue.log 2>&1 & echo $$! > ../queue.pid


start-all: consumer like queue
	@echo "All services started."


stop-consumer:
	@echo "Stopping consumer-service..."
	@kill `cat consumer.pid` && rm consumer.pid || echo "consumer-service not running."

stop-like:
	@echo "Stopping like-service..."
	@kill `cat like.pid` && rm like.pid || echo "like-service not running."

stop-queue:
	@echo "Stopping queue-service..."
	@kill `cat queue.pid` && rm queue.pid || echo "queue-service not running."


stop-all: stop-consumer stop-like stop-queue
	@echo "All services stopped."
