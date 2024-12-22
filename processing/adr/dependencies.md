# Зависимости

В этом файле описаны наши зависимости, которые мы используем для получения данных, необходимый для работы сервиса

## Используемые методы других сервисов

### address.CreateProcessingAddress - создать адрес

https://code.emcdtech.com/emcd/blockchain/address/-/blob/master/protocol/proto/address.proto?ref_type=heads#L34

Input:

* user_uuid - ID emcd-юзера мерчанта, должны получать от мерчанта при создании инвойса
* network
* processing_uuid - ключ идемпотентности

Output:

* address_uuid - мы используем для привязки к платежу как ключ вместо пары (адрес, сеть)
* address
* ...

Заметки

* Предварительно, у себя надо будет хранить все выписанные адреса, и в какой они сети

### accounting.GetOrCreateUserAccountByArgs - получить account_id по user_id

https://code.emcdtech.com/emcd/service/accounting/-/blob/coinwatch/repository/user_account.go?ref_type=heads#L22

Input:

* userId - легаси инт user id, получить из profile сервиса
* userIdNew - uuid, который у нас уже есть
* coinIdNew - придет в сообщении от coinhold, также можем получить из coin service и сразу требовать от фронта
* accountTypeId = WalletAccountTypeID
* minPay = 0
* fee = 0

### coins.GetCoins - получить все существующие монеты

https://code.emcdtech.com/emcd/service/coin/-/blob/master/protocol/proto/coin.proto?ref_type=heads#L12

Мы можем использовать этот метод, чтобы получить все монеты, для каждой монеты ее сети и отдать это на фронт, чтоб
пользователь мог выбрать нужную. По получении монеты мы сохраняем ее вместе с данными инвойса.

Надо в запросе передать лимит 999999, чтоб получить все монеты. Обновлять регулярно (раз в минуту, например)

### profile.GetByUserIDV2 - получить юзера по id

https://code.emcdtech.com/emcd/service/profile/blob/382a6e4d9d04b4797c11c6acf09fe29adcad6d4f/protocol/proto/profile.proto#L14

Получает пользователя по id, там есть oldID

### accounting.ChangeBalance - перевести деньги

Используем, чтоб списать комиссию с мерчанта

Input

* type - надо завести нам новый тип под процессинг
* createdAt = time.Now()
* senderAccountID - мерчант
* receiverAccountID - наш счет для прибыли
* amount - размер комиссии
* actionID - ключ идемпотентности
* coinID - int, старый id монеты

## Coinwatch сообщения

Это сообщение кидается, когда мы видим транзакцию, и после того, как она подтверждена и мы начислили деньги на баланс
мерчанта. Надо делать только для процессинговых адресов.

```json
{
  "address": "some_addr",
  "amount": "123.456",
  "tx_hash": "...",
  "coin_code": "xxx",
  "is_confirmed": true,
  "user_uuid": "merchant_id",
  "processing_uuid": "processing_id"
}
```

`received_address` - адрес, который был получен при создании адреса в address service. На него пришли деньги
`is_confirmed` - если false, значит транзакция еще не подтверждена, если true, значит подтверждена и деньги зачислены на
баланс мерчанта
`processing_uuid` - ключ идемпотентности, использованный при создании адреса
