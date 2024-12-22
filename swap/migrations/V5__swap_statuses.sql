CREATE TABLE swap.swap_statuses
(
    id          int PRIMARY KEY,
    name        varchar(30)  NOT NULL,
    description varchar(255) NOT NULL
);

INSERT INTO swap.swap_statuses (id, name, description)
VALUES (0, 'UNKNOWN', 'Неизвестный статус'),
       (1, 'WaitDeposit', 'Ожидаем поступления депозита'),
       (2, 'CheckDeposit', 'Проверяем поступивший депозит. ANL проверки и т.д.'),
       (3, 'DepositError', 'deprecated'),
       (4, 'TransferToUnified', 'Переводим с FUND счёта на торговый, для дальнейшего размещения ордеров'),
       (5, 'CreateOrder', 'Создаём ордер, сохраняем в БД'),
       (6, 'PlaceOrder', 'Размещаем ордер на биржи'),
       (7, 'CheckOrder', 'Проверяем статус ордера'),
       (8, 'PlaceAdditionalOrder', 'Если не прямой обмен, размещаем дополнительный ордер'),
       (9, 'CheckAdditionalOrder', 'Проверяем статус дополнительного ордера'),
       (10, 'TransferFromSubToMaster', 'Переводим монеты с суб аккаунта на мастер для последующего вывода'),
       (11, 'CheckTransferFromSubToMaster', 'Проверяем состояние перевода'),
       (12, 'PrepareWithdraw', 'Подготавливаем вывод, сохраняем в БД'),
       (13, 'WithdrawSwapStatus', 'Выводим монеты пользователю'),
       (14, 'WaitWithdraw', 'Ждём пока вывод попадёт в блокчейн'),
       (15, 'Completed', 'Своп завершён успешно'),
       (16, 'Cancel', 'Своп отменён'),
       (17, 'Error', 'Ошибка во время выполнения свопа на любом из этапов'),
       (18, 'ManualCompleted', 'Довели своп вручную, после ошибки, монеты выведены пользователю')