build:
	go build -ldflags="-s -w" -trimpath

deploy: build
	rsync -avzL docs privtracker cert ubuntu@tracker.whereami.com.cn:

