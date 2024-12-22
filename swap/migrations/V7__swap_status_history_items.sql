CREATE TABLE swap.swap_status_history_items
(
    swap_id UUID REFERENCES swap.swaps(id) NOT NULL,
    status  SMALLINT REFERENCES swap.swap_statuses(id) NOT NULL,
    set_at  TIMESTAMP NOT NULL,
    CONSTRAINT PK_swap_id_set_at PRIMARY KEY (swap_id, set_at)
);
