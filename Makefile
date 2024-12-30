generate-mock-all:
	mockgen -source=internal/repositories/interfaces/balance_history.go -destination=internal/repositories/mocks/mock_balance_history.go -package=mocks
	mockgen -source=internal/repositories/interfaces/wallet_balance.go -destination=internal/repositories/mocks/mock_wallet_balance.go -package=mocks
