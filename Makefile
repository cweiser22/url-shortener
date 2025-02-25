# App Variables
APP_NAME = urls-service
ANALYTICS_NAME = analytics-service
K8S_DEPLOYMENT_DIR = k8s

# Docker Compose
compose-up:
	sudo docker compose up -d

compose-build:
	sudo docker compose build

compose-down:
	sudo docker compose down -v

# Kubernetes
k8s-up: docker-hub-push k8s-load k8s-configure k8s-deploy
k8s-down: k8s-clean k8s-clean-configure

k8s-build:
	sudo docker build -t urls-service:latest ./urls_service
	sudo docker build -t analytics-service:latest ./analytics_service

# Load images into local Kubernetes (for Minikube or Kind)
k8s-load:
	minikube image load urls-service:latest
	minikube image load analytics-service:latest

k8s-deploy:
	kubectl apply -f $(K8S_DEPLOYMENT_DIR)/urls-deployment.yml
	kubectl apply -f $(K8S_DEPLOYMENT_DIR)/analytics-deployment.yml
	kubectl apply -f $(K8S_DEPLOYMENT_DIR)/ingress.yml
	helm install url-kafka bitnami/kafka --namespace kafka --create-namespace


k8s-clean:
	kubectl delete -f $(K8S_DEPLOYMENT_DIR)/urls-deployment.yml
	kubectl delete -f $(K8S_DEPLOYMENT_DIR)/analytics-deployment.yml
	kubectl delete -f $(K8S_DEPLOYMENT_DIR)/ingress.yml

k8s-configure:
	kubectl create configmap urls-config --from-env-file=./urls_service/.env
	kubectl create configmap analytics-config --from-env-file=./analytics_service/.env
	kubectl create secret generic influx-secrets --from-env-file=./.env.influx
	kubectl create secret generic mongo-secrets --from-env-file=./.env.mongo
	kubectl create secret generic redis-secrets --from-env-file=./.env.redis
	kubectl create secret generic docker-secrets --from-env-file=./.env.dockerhub
	kubectl create secret generic kafka-secrets --from-env-file=./.env.kafka


k8s-clean-configure:
	kubectl delete configmap urls-config
	kubectl delete configmap analytics-config
	kubectl delete secret influx-secrets
	kubectl delete secret mongo-secrets
	kubectl delete secret redis-secrets
	kubectl delete secret docker-secrets

docker-hub-push:
	sudo docker build --platform=linux/amd64 -t cooperw22/urls-service:latest ./urls_service
	sudo docker build --platform=linux/amd64 -t cooperw22/analytics-consumer:latest -f new_analytics_service/consumer.Dockerfile ./new_analytics_service
	sudo docker build --platform=linux/amd64 -t cooperw22/analytics-producer:latest -f new_analytics_service/producer.Dockerfile ./new_analytics_service
	docker push cooperw22/urls-service:latest
	docker push cooperw22/analytics-service:latest

postgres-init:
	# use docker run to run init.sql on postgres container

k8s-install-kafka:
	helm install url-kafka bitnami/kafka --namespace kafka --create-namespace -f k8s/kafka-overrides.yml