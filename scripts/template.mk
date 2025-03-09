project_structure:
	#create main.go
	mkdir -p cmd/app
	touch cmd/app/main.go

	#create app dir
	mkdir -p internal/app
	touch internal/app/app.go

	#create config file
	mkdir -p configs
	touch configs/config.toml

	# create internal dirs
	mkdir -p internal/ports/http
	mkdir -p internal/ports/grpc
	mkdir -p internal/services
	mkdir -p internal/clients
	mkdir -p internal/brokers
	mkdir -p internal/domains/entity
	mkdir -p internal/domains/enum
	mkdir -p internal/repositories
	mkdir -p internal/util

	#create dirs for proto files dependencies
	mkdir -p deps/api/proto
	mkdir -p pkg/deps/api/proto

	#create dir for service proto files
	mkdir -p api/proto
	mkdir -p pkg/api/proto


	touch .env
	echo "ENV=dev" > .env

	touch .gitignore
	echo ".env" > .gitignore
	echo ".idea" >> .gitignore

