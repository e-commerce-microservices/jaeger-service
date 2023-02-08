.PHONY: rebuild
rebuild:
	docker build -t ngoctd/test:latest . && \
	docker push ngoctd/test

.PHONY: redeploy
redeploy:
	kubectl rollout restart deployment depl-test

.PHONY: redeploy_jaeger
redeploy_jaeger:
	kubectl rollout restart deployment jaeger-all-in-one