
target=quay.io/cw/kubernetes-entrypoint:v0.3.1d

default: runtest

container:
	sudo docker build -t $(target) .
	sudo docker images $(target)

runtest: container
	sudo docker run --rm -e ENTRYPOINT_DEBUG_LOAD=true $(target) 2>&1

push: container
	sudo docker push $(target)


clean:
	rm -f *~
