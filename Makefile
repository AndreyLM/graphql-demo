init:
	go run scripts/gqlgen.go init
update:
	go run scripts/gqlgen.go -v
run:
	go run server/server.go
dataloader:
	cd ./api/dataloaders && dataloaden VideoLoader int []*github.com/andreylm/graphql-demo/api.Video