# Accounting

Это центральный микросервис для работы с балансами. Любая логика работы с балансами должна быть прописана тут.
Микросервис работает по протоколу gRPC. На данный момент реализованы следующие функции:

Работа с балансами:
 - ViewBalance
 - ChangeBalance
 - ChangeBalanceWithBlock
 - ChangeBalanceWithUnblock
 - Unblock
 - ChangeMultipleBalance

Выборка данных:
 - FindOperations
 - FindBatchOperations
 - FindOperationsWithBlocks
 - FindTransactionsWithBlocks
 - SearchOperations

Аналитика:
 - FindLastBlockTimeBalances
 - FindBalancesDiffMining
 - FindBalancesDiffWallet

История:
 - GetHistory
