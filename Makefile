
target=quay.io/cw/kubernetes-entrypoint:v0.3.0c

default: runtest

container:
	rm -f kubernetes-entrypoint *~
	sudo docker build -t $(target) . | cat
	sudo docker images $(target)

runtest: container
	sudo docker run --rm -e ENTRYPOINT_DEBUG_LOAD=true $(target) 2>&1

push: container
	sudo docker push $(target) | cat


clean:
	rm -f *~
