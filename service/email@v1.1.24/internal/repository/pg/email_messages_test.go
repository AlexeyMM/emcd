package pg

//nolint
//
//func TestEmailMessages_Create(t *testing.T) {
//	ctx := context.Background()
//	defer truncateEmailMessages(ctx)
//	expected := &model.EmailMessageEvent{
//		ID:        uuid.New(),
//		Email:     "my_email",
//		Type:      model.WorkerChangedState,
//		CreatedAt: time.Now().Truncate(time.Millisecond).UTC(),
//	}
//	err := emailMessagesRepo.Create(ctx, expected)
//	require.NoError(t, err)
//	actual := selectEmailMessages(ctx)
//	require.Len(t, actual, 1)
//	require.Equal(t, expected, actual[0])
//}
//
//func selectEmailMessages(ctx context.Context) []*model.EmailMessageEvent {
//	rows, err := dbPool.Query(ctx, "SELECT id,email,type,created_at FROM sent_email_messages")
//	if err != nil {
//		log.Fatal().Msg(err.Error())
//	}
//	res := make([]*model.EmailMessageEvent, 0)
//	for rows.Next() {
//		var e model.EmailMessageEvent
//		err = rows.Scan(&e.ID, &e.Email, &e.Type, &e.CreatedAt)
//		if err != nil {
//			log.Fatal().Msg(err.Error())
//		}
//		res = append(res, &e)
//	}
//	rows.Close()
//	return res
//}
//
//func truncateEmailMessages(ctx context.Context) {
//	_, err := dbPool.Exec(ctx, "TRUNCATE TABLE sent_email_messages")
//	if err != nil {
//		log.Fatal().Msg(err.Error())
//	}
//}
