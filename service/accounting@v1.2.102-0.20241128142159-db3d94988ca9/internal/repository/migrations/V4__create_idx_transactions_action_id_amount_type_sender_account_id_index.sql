-- На проде юзается CONCURRENTLY в этих двух инструкциях.
-- Из-за особенностей мигратора здесь не можем прописать CONCURRENTLY,
-- так как мигратор под капотом создает транзакцию, в которой нельзя
-- использовать CONCURRENTLY.
DROP INDEX emcd.transactions_action_id_amount_type_index
;

CREATE UNIQUE INDEX transactions_action_id_amount_type_sender_account_id_index
    on emcd.transactions (action_id, amount, type, sender_account_id)
;
