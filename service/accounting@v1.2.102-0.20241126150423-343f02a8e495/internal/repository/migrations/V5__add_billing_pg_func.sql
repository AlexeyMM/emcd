-- Создание проводки по биллингу
CREATE OR REPLACE FUNCTION emcd.accountingCreateTransactions(
    BalanceWhiteList text[],
    ProcessingType varchar,
    ActionId text,
    ToAccountId integer,
    FromAccountId integer,
    AmountSum numeric,
    TrTypeId integer,
    CommentData varchar,
    HashData varchar,
    HashrateData bigint,
    ReceiverAddress varchar,
    TokenID integer,
    BlockedTill timestamp,
    UnblockToAccount integer,
    CreatedAt timestamp,
    FromReferralId integer
) RETURNS INT AS
$$
DECLARE
    TrID  INT;
    CoinID INT;
    CoinIdCheck BOOLEAN;
    Balance NUMERIC;
    WhiteListCheck INT;
    --BalanceWhiteList TEXT ARRAY;
    BlockID INT;
    BlockSenderID INT;
    BlockReceiverID INT;
    BlockCoinID INT;
    BlockAmount NUMERIC;
    UnBlockTrID INT;
    SenderUsername varchar;
BEGIN
    -- ProcessingType: default block unblock

    --Белый лист, не проверяем балансы тут для тестов pg
    --SELECT ARRAY ['admin', 'payouts_node', 'emcd_lost_users', 'coinholdref'] INTO BalanceWhiteList;

    --если тип проводки не указан
    IF ProcessingType IS NULL OR ProcessingType = '' THEN
        RAISE EXCEPTION 'processing type is empty, sender_account_id: %, receiver_account_id: %', FromAccountId, ToAccountId;
    END IF;
    --если ActionId не указан
    IF ActionId IS NULL OR ActionId = '' THEN
        RAISE EXCEPTION '%: action_id is empty, sender_account_id: %, receiver_account_id: %', ProcessingType, FromAccountId, ToAccountId;
    END IF;

    /* Убрали для обратной совместимости при внедрении
    IF EXISTS (SELECT FROM emcd.accounting_actions WHERE action_id = ActionId::uuid)
    THEN
        --RAISE NOTICE '%', ActionId;
    ELSE
        RAISE EXCEPTION 'action is not found, action_id: %', ActionId;
    END IF;
    */


    --Получаем id монеты и проверяем что оба аккаунта в одной монете
    SELECT max(coin_id),  (sum(coin_id)/count(coin_id) = max(coin_id)) FROM emcd.users_accounts WHERE id IN(FromAccountId,ToAccountId,UnblockToAccount) INTO CoinID, CoinIdCheck;
    IF CoinIdCheck != true THEN
        RAISE EXCEPTION '%: coin check is fail, sender_account_id: %, receiver_account_id: %, to_account_id: %', ProcessingType, FromAccountId, ToAccountId, UnblockToAccount;
    END IF;

    IF ProcessingType = 'default' OR ProcessingType = 'block' THEN
        --проверяем пользователя на сервисный акк из белого листа
        SELECT COALESCE(array_position(BalanceWhiteList, LOWER(u.username)::text), 0), LOWER(u.username)
        FROM emcd.users_accounts ua
                 INNER JOIN emcd.users u on u.id = ua.user_id
        WHERE ua.id=FromAccountId INTO WhiteListCheck, SenderUsername;

        --Проверяем баланс отправителя на наличие средств, если он не в белом списке
        IF WhiteListCheck = 0 THEN
        SELECT COALESCE(ROUND(sum(amount), 8), 0) FROM emcd.operations WHERE account_id=FromAccountId INTO Balance;
        IF AmountSum > Balance THEN
            RAISE EXCEPTION '%: balance is less than the transfer amount, sender_account_id: %, balance: %, amount: %, username: %', ProcessingType, FromAccountId, Balance, AmountSum, SenderUsername;
        END IF;
        END IF;
    END IF;

    --собираем все необходимые данные из транзакции блокировки и записи блокировки, даже если это запрос на блокировку(проверяем, а вдруг она уже есть)
    IF ProcessingType = 'block' THEN
        SELECT tb.id, tb.unblock_transaction_id, t.sender_account_id, t.receiver_account_id, t.coin_id, t.amount FROM emcd.transactions_blocks tb
                   LEFT JOIN emcd.transactions t on t.id = tb.block_transaction_id
        WHERE t.action_id = ActionId::uuid and t.sender_account_id=FromAccountId and t.receiver_account_id=ToAccountId and amount=AmountSum and t.type=TrTypeId
                   LIMIT 1 INTO BlockID, UnblockTrID, BlockSenderID, BlockReceiverID, BlockCoinID, BlockAmount;
        ELSEIF ProcessingType = 'unblock' OR ProcessingType = 'reject' THEN
            SELECT tb.id, tb.unblock_transaction_id, t.sender_account_id, t.receiver_account_id, t.coin_id, t.amount FROM emcd.transactions_blocks tb
                   LEFT JOIN emcd.transactions t on t.id = tb.block_transaction_id
            WHERE t.action_id = ActionId::uuid
              and t.receiver_account_id=FromAccountId
              and amount=AmountSum
              and tb.unblock_transaction_id IS NULL
                   LIMIT 1 INTO BlockID, UnblockTrID, BlockSenderID, BlockReceiverID, BlockCoinID, BlockAmount;
    END IF;

    --отклоняем запрос на блокировку, если она существует
    IF ProcessingType = 'block' AND BlockID IS NOT NULL THEN
        RAISE EXCEPTION '%: balance block is already exists, action_id: %, block_id: %', ProcessingType, ActionId, BlockID;
    END IF;

    --отклоняем запрос на разблокировку, если уже разблокировано или не существует
    IF ProcessingType = 'unblock' OR ProcessingType = 'reject' THEN
        IF UnblockTrID IS NOT NULL THEN
            RAISE EXCEPTION '%: balance unblock is already exists, action_id: %, block_id: %, Amount: %', ProcessingType, ActionId, BlockID, AmountSum;
        END IF;
        IF BlockID IS NULL THEN
            RAISE EXCEPTION '%: balance block is not found, action_id: %, type: %, amount: %', ProcessingType, ActionId, TrTypeId, AmountSum;
        END IF;
    END IF;

    IF CreatedAt IS NULL THEN
        SELECT now() INTO CreatedAt;
    END IF;

    INSERT INTO emcd.transactions (action_id, type, sender_account_id, receiver_account_id, coin_id, amount, comment, receiver_address, hashrate, token_id, hash, created_at, from_referral_id)
    VALUES (ActionId::uuid, TrTypeId, FromAccountId, ToAccountId, CoinID, AmountSum, CommentData, ReceiverAddress, HashrateData, TokenID, HashData, CreatedAt, FromReferralId)
    RETURNING id INTO TrID;

    INSERT INTO emcd.operations (type, transaction_id, account_id, coin_id, amount, hash, created_at)
    VALUES (TrTypeId, TrID, FromAccountId, CoinID, -AmountSum, HashData, CreatedAt),
           (TrTypeId, TrID, ToAccountId, CoinID, AmountSum, HashData, CreatedAt);

    --если это блокировка, создаем запись блокировки
    IF ProcessingType = 'block' THEN
        --если время жизни блокировки не указано - ставим на год
        IF BlockedTill is NULL THEN
            SELECT now() + '1 YEAR' INTO BlockedTill;
        END IF;
        --если время жизни блокировки истекло - вернем ошибку
        IF BlockedTill < NOW()  THEN
            RAISE EXCEPTION '%: blocking time has passed, blocking_till: %, now: %, action_id: %', ProcessingType, BlockedTill, NOW(), ActionId;
        END IF;
        INSERT INTO emcd.transactions_blocks (block_transaction_id, unblock_transaction_id, unblock_to_account_id, blocked_till)
        VALUES (TrID, NULL, UnblockToAccount, BlockedTill)  RETURNING id INTO BlockID;
    END IF;

    --если это разблокировка - обновляем запись блокировки
    IF ProcessingType = 'unblock' THEN
        UPDATE emcd.transactions_blocks SET unblock_transaction_id = TrID WHERE id = BlockID;
        IF HashData IS NOT NULL AND length(HashData) > 0 THEN
            UPDATE emcd.transactions SET hash=HashData WHERE id=(SELECT block_transaction_id FROM emcd.transactions_blocks WHERE id=BlockID) and hash IS NULL;
        END IF;
    END IF;

    RETURN TrID;
END
$$
    LANGUAGE plpgsql;
