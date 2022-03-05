.PHONY: swag

swag:
	@echo "!!! Running swag !!!"
	swag init --md ./ --pd -g ./pkg/adapters/http/*.go
