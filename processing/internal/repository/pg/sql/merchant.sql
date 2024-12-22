-- name: GetMerchantWithTariff :one
select m.id, mt.upper_fee, mt.lower_fee, mt.min_pay, mt.max_pay
from merchant m
         join public.merchant_tariff mt on m.id = mt.merchant_id
where m.id = @id;

-- name: SaveMerchantID :exec
insert into merchant (id)
values (@id)
on conflict do nothing;

-- name: SaveMerchantTariff :exec
insert into merchant_tariff (merchant_id, upper_fee, lower_fee, min_pay, max_pay)
values (@merchant_id, @upper_fee, @lower_fee, @min_pay, @max_pay)
on conflict (merchant_id) do update set upper_fee = excluded.upper_fee,
                                        lower_fee = excluded.lower_fee,
                                        min_pay   = excluded.min_pay,
                                        max_pay   = excluded.max_pay;