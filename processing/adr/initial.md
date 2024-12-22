# ADR по сервису процессинга

## Мерчант создает платеж и отправляет юзера на оплату

```mermaid
sequenceDiagram
    autonumber
    actor M as Merchant
    participant P as Processing
    participant DB as DB
    participant As as Address Service
    actor B as Buyer
    B ->> M: Want to pay X$
    M ->> P: Create invoice for X$
    activate M
    note over M, P: external_id, amount,<br>currency, network, email, ...
    P ->> DB: Get deposit address from pool
    alt no available address in pool
        P -->> As: rpc CreateProcessingAddress()
        As -->> P: AddressResponse
        P -->> DB: Save received new address to pool
    end
    P ->> P: Save to DB
    P ->> M: Invoice info (id, link, ...)
    deactivate M
    M ->> B: Invoice link
```

> Часть данных может быть неизвестна при создании платежа. В таком случае фронт мерчанта или наш (наша форма оплаты)
> должна получить их от пользователя и только потом кинуть запрос на бэк для создания инвойса. Инвойс не может создаться
> при отсутствии хотя бы одного обязательного поля.

Если мерчанту хочет использовать нашу форму оплаты и предзаполнить в ней часть данных, то он может сделать это,
используя query-параметры ссылки на нашу форму, которую вставляет на наш сайт. Наш фронт при рендеринге возьмет данные и
заполнит нужные поля формы. После получения остальных полей от юзера он отправит запрос на создание инвойса и пойдет
дальше по флоу оплаты инвойса.

## Покупатель оплачивает инвойс

```mermaid
sequenceDiagram
    autonumber
    actor B as Buyer
    participant P as Processing
    participant Ac as Accounting
    participant R as RabbitMQ
    participant Cw as CoinWatch
    participant Bc as Blockchain

   Note over P: Invoice status WAITING_FOR_DEPOSIT
   loop until PAYMENT_ACCEPTED
        B -->> Bc: Transfer money
        Bc -->> Cw: CoinWatch detects the money
        Cw -->> R: Transaction acknowledged, not confirmed
        R -->> P: Transaction acknowledged, not confirmed
        Note over P: Invoice status PAYMENT_CONFIRMATION
        Cw -->> Ac: Money deposited on merchant account
        Cw -->> R: Transaction confirmed
        R -->> P: Transaction confirmed
        alt Not full amount
            Note over P: Invoice status PARTIALLY_PAID
        else Full amount
            Note over P: Invoice status PAYMENT_ACCEPTED
        end
    end
    P ->> P: Return deposit wallet into pool
    P ->> Ac: Transfer fee from merchant account to revenue account
    Note over P: Invoice status FINISHED
```

Address Pool, Processing - наши сервисы, мы их напишем с 0.

Accounting - используем существующий сервис.

Заметки:

1. CoinWatch отправляет данные через RabbitMQ, контракт согласуем перед началом разработки
2. Для отправки денег от мерчанта на счет с комиссиями используем метод ChangeBalance

## Схема движения средств

```mermaid
sequenceDiagram
    actor B as Buyer
    participant CW as Cold Wallet
    actor M as Merchant
    participant PA as Processing Account<br>(dedicated)
    B ->> CW: Transfers money
    CW ->> M: CoinWatch transfers money to merchant
    M ->> PA: Fee withdrawn from merchant
```

