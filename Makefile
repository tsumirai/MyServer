all: install

build: clean gen_output build_without_glide
	mkdir -p output/bin
	mv $(GOPATH)/src/MyServer/MyServer/MyServer output/bin/

clean:
	rm -rf output/bin
	rm -rf output/conf

gen_output:
	mkdir -p output/log
	mkdir -p output/conf
	cp -r conf output/

build_without_glide:
	go build

run:
	cd ./output/ && ./bin/MyServer # 前台运行,方便退出

nohup:
	cd ./output/ && nohup ./bin/MyServer &