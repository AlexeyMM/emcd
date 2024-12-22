# Coin Service

## Introduction.

>
>  A service responsible for saving, receiving, updating static data for all crypto coins in the system.
> ```
> Gitlab repository: https://code.emcdtech.com/emcd/service/coin
> ```
> This is an internal service designed to interact with our other services through the gRPC protocol. There should not be a direct appeal of clients to this service. The service has a separate database of coin.

## Integration In Project.
#### To integrate the service with our other services via gRPC, you need to:
>```go
>   replace code.emcdtech.com/emcd/service/coin => code.emcdtech.com/emcd/service/coin.git v.(последная версиа)
>```
>

## Error handling.
>
> On the service side, when an error occurs, a status code and an error body are returned, consisting of:
> ```json
>   {
>       "code":   "string", – unique error identifier.
>       "message": "string", – a more detailed explanation of the error.
>   }
>```
> Error codes:
>  ```json
>   C000001 - Internal server error.
>   C000002 - Forbidden request for user.
>```

## Functional.
> ### Available for admin only:
> #### Creating a new coin in wallets:
> - Request:
> ```json
> rpc CreateWalletCoin(CreateWalletCoinRequest) returns (CreateWalletCoinResponse) {}
>
> {
>    "coin": {
>        "coin": "string",
>        "title": "string",
>        "description": "string",
>        "media_id": "string",
>        "decimals": double,
>        "withdrawal_fee": double,
>        "withdrawal_min_limit": double,
>        "is_active": bool
>    },
>    "token": "string"
> }
>```
> - Response:
> ```json
> {
>    "success": bool
> }
>```
> #### Creating a new network for a coin:
>
> - Request:
> ```json
> rpc CreateWalletCoinNetwork(CreateWalletCoinNetworkRequest) returns (CreateWalletCoinNetworkResponse) {}
>
> {
>    "coin": {
>        "coin": "string",
>        "code": "string",
>        "title": "string",
>        "description": "string",
>        "is_active": bool
>    },
>    "token": "string"
> }
>```
> - Response:
> ```json
> {
>    "success": bool
> }
>```
> #### Creating a new coin in mining:
>
> - Request:
> ```json
> rpc CreateMiningCoin(CreateMiningCoinRequest) returns (CreateMiningCoinResponse) {}
>
> {
>    "coin": {
>        "coin": "string",
>        "code": "string",
>        "title": "string",
>        "description": "string",
>        "fee": double,
>        "marge_coin": "string",
>        "is_active": bool
>    },
>    "token": "string"
> }
>```
> - Response:
> ```json
> {
>    "success": bool
> }
>```
> #### Updating coins in wallets:
> - Request:
> ```json
> rpc UpdateWalletCoin(UpdateWalletCoinRequest) returns (UpdateWalletCoinResponse) {}
>
> {
>    "coin": {
>        "coin": "string",
>        "title": "string",
>        "description": "string",
>        "media_id": "string",
>        "decimals": double,
>        "withdrawal_fee": double,
>        "withdrawal_min_limit": double,
>        "is_active": bool
>    },
>    "token": "string"
> }
>```
> - Response:
> ```json
> {
>    "success": bool
> }
>```
> #### Coin network update:
>
> - Request:
> ```json
> rpc UpdateWalletCoinNetwork(UpdateWalletCoinNetworkRequest) returns (UpdateWalletCoinNetworkResponse) {}
>
> {
>    "coin": {
>        "coin": "string",
>        "code": "string",
>        "title": "string",
>        "description": "string",
>        "is_active": bool
>    },
>    "token": "string"
> }
>```
> - Response:
> ```json
> {
>    "success": bool
> }
>```
> #### Mining coin update:
>
> - Request:
> ```json
> rpc UpdateMiningCoin(UpdateMiningCoinRequest) returns (UpdateMiningCoinResponse) {}
>
> {
>    "coin": {
>        "coin": "string",
>        "code": "string",
>        "title": "string",
>        "description": "string",
>        "fee": double,
>        "marge_coin": "string",
>        "is_active": bool
>    },
>    "token": "string"
> }
>```
> - Response:
> ```json
> {
>    "success": bool
> }
>```
> ### Only available for all services:
> #### Getting a list of coins in wallets::
>
> - Request:
> ```json
> rpc GetWalletCoins(GetCoinsRequest) returns (GetWalletCoinsResponse) {}
>
> {
>     "sort_type": "string",
>     "limit": int,
>     "offset" int
> }
>```
> - Response:
> ```json
> {
>    "coins": [] {
>        "coin": "string",
>        "title": "string",
>        "description": "string",
>        "media_id": "string",
>        "decimals": double,
>        "withdrawal_fee": double,
>        "withdrawal_min_limit": double,
>        "is_active": bool
>    },
>    "total_count": int
> }
>```
> #### Getting a list of coins in wallets::
>
> - Request:
> ```json
> rpc GetMiningCoins(GetCoinsRequest) returns (GetMiningCoinsResponse) {}
>
> {
>     "sort_type": "string",
>     "limit": int,
>     "offset" int
> }
>```
> - Response:
> ```json
> {
>    "coins": [] {
>        "coin": "string",
>        "code": "string",
>        "title": "string",
>        "description": "string",
>        "fee": double,
>        "marge_coin": "string",
>        "is_active": bool
>    },
>    "total_count": int
> }
>```

## Confluence link:
> https://emcd-io.atlassian.net/wiki/spaces/PD/pages/484213177/Coin+Service