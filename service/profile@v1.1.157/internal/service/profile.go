package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	userAccountModel "code.emcdtech.com/emcd/service/accounting/model"
	accountingRepository "code.emcdtech.com/emcd/service/accounting/repository"

	userAccountModelEnum "code.emcdtech.com/emcd/service/accounting/model/enum"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"code.emcdtech.com/emcd/sdk/log"

	"code.emcdtech.com/emcd/service/profile/internal/jobs"
	"code.emcdtech.com/emcd/service/profile/internal/model"
	"code.emcdtech.com/emcd/service/profile/internal/notification"
	"code.emcdtech.com/emcd/service/profile/internal/repository"
)

const (
	kycApprovedStatus = 2
	defFee            = 0.015
)

var (
	ErrInconsistentUserInfo    = errors.New("user info inconsistent between emcd.users and profile.users")
	ErrAddressChangeNotAllowed = errors.New("address change not allowed")
	ErrBothAddressIsTheSame    = errors.New("both address is the same")
	ErrMinPayNotValid          = errors.New("min pay value is not valid")
	ErrCoinNotFound            = errors.New("coin not found")
	ErrTokenNotFound           = errors.New("token not found")
	ErrNotUniqueUsername       = errors.New("not unique username")
)

type Profile interface {
	Create(ctx context.Context, pr *model.Profile) (int32, error)
	SaveV3(ctx context.Context, pr *model.Profile) (int32, error)
	SaveV4(ctx context.Context, u *model.User) (int32, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*model.Profile, error)
	GetUserIDByOldID(ctx context.Context, oldID int32) (uuid.UUID, error)
	UpdatePassword(ctx context.Context, email, password string, whiteLabelID uuid.UUID) error
	GetUserByEmailAndWl(ctx context.Context, email string, whiteLabelID uuid.UUID) (*model.User, error)
	GetOldUserByEmailAndWl(ctx context.Context, email string, whiteLabelID uuid.UUID) (*model.User, error)
	GetAllUsers(ctx context.Context, skip, take int, sortField string, asc bool, searchEmail string) ([]*model.User, int, error)
	GetAllUsersByWlID(
		ctx context.Context,
		skip, take int,
		sortField string,
		asc bool,
		searchEmail string,
		whiteLabelID uuid.UUID,
	) ([]*model.User, int, error)
	GetUserByApiKey(ctx context.Context, apiKey string) (*model.User, error)
	GetOldByEmailAndWl(ctx context.Context, email string, whiteLabelID uuid.UUID) (*model.Profile, error)
	GetSuspended(ctx context.Context, userID uuid.UUID) (bool, error)
	SetSuspended(ctx context.Context, userIDs []uuid.UUID, suspended bool) error
	SendEmailToChangeAddress(ctx context.Context, userID int, username, ip, coinCode, address, domain string) error
	UpdateMinPay(ctx context.Context, userID int, coinCode string, value float32) error
	ChangeWalletAddressConfirm(ctx context.Context, userID int32, token string) (*model.ChangeWalletAddressConfirmResult, error)
	UpdateUserIsActive(ctx context.Context, email string, active bool) error
	GetUserIsActive(ctx context.Context, userID uuid.UUID) (bool, error)
	GetKycStatus(ctx context.Context, userID int) (*model.Kyc, error)
	SetKycStatus(ctx context.Context, userID, status int) error
	InsertKycHistory(ctx context.Context, userID int, data []byte) error
	CheckAppleAccount(ctx context.Context, appleID, email string) (bool, string, error)
	GetNotificationSettings(ctx context.Context, userID uuid.UUID) (*model.NotificationSettings, error)
	SaveNotificationSettings(ctx context.Context, settings *model.ChangeableNotificationSettings) error
	RelatedUsers(ctx context.Context, firstID, secondID uuid.UUID) (bool, error)
	GetAllSubUsers(ctx context.Context, userID uuid.UUID) ([]*model.User, error)
	GetAllUserIDsByUsername(ctx context.Context) (map[string]uuid.UUID, error)
	GetReferrals(ctx context.Context, userID uuid.UUID, skip, take int, sortField string, asc bool) ([]*model.User, int, error)
	GetUsernamesByIDs(ctx context.Context, userIDs []uuid.UUID) (map[uuid.UUID]string, error)
	GetEmailsByIDs(ctx context.Context, userIDs []uuid.UUID) (map[uuid.UUID]string, error)
	SafeDeleteByID(ctx context.Context, userID uuid.UUID) error
	GetByUsernames(ctx context.Context, usernames []string) ([]*model.Profile, error)
	UpdateRefID(ctx context.Context, oldUserID int32, userID uuid.UUID, refID int32, newRefID uuid.UUID) error
	WasReferralLinkGenerated(ctx context.Context, userID uuid.UUID) (bool, error)
	SetFlagReferralLinkGenerated(ctx context.Context, userID uuid.UUID, value bool) error
	UpdateProfile(ctx context.Context, u *model.User) error
	SetUserAttributes(ctx context.Context, userID uuid.UUID, attrs model.UserAttributes) error
	SetTimezone(ctx context.Context, userID uuid.UUID, timezone string) error
	SetLanguage(ctx context.Context, userID uuid.UUID, language string) error
	GetAddresses(ctx context.Context, userID uuid.UUID) ([]*model.Address, error)
	ChangeWalletAddress(ctx context.Context, userID int32, coinID int, address, coinName string) error
	GetCoinByCode(ctx context.Context, coinCode string) (*model.Coin, error)

	GetCountUserWithWL(ctx context.Context, wlUUID uuid.UUID, offset, limit int32) ([]model.UserShortInfo, int64, error)

	GetUsersByUUIDs(ctx context.Context, userUUIDs []uuid.UUID) ([]model.UserShortInfo, error)

	/* API key's section */

	// CreateNewAPIKey create new API key for user, parentID could be empty
	CreateNewAPIKey(ctx context.Context, req CreateNewAPIKeyReq) (string, error)
	// DeleteAPIKey deletes api key for user
	DeleteAPIKey(ctx context.Context, req DeleteAPIKeyReq) error
	// GetAPIKey by user id or pair user uuid and parent uuid
	GetAPIKey(ctx context.Context, req GetAPIKeyReq) (string, error)

	/* No pay section */

	// GetNoPayStatus get no pay status by user UUID
	GetNoPayStatus(ctx context.Context, userUUID uuid.UUID) (bool, time.Time, error)
	// UpdateNoPayToFalse set no pay to false
	UpdateNoPayToFalse(ctx context.Context, userUUID uuid.UUID) error
	// CancelJobOffNoPay cancel no pay job for user id
	CancelJobOffNoPay(ctx context.Context, userUUID uuid.UUID) error

	SoftDeleteSubUser(ctx context.Context, subUserID, newParentID uuid.UUID) error

	GetByUsernamesForReferrals(ctx context.Context, usernames []string) ([]*model.Profile, error)
	CreateSubUser(ctx context.Context, parentUUID uuid.UUID, username string, addresses []*model.CoinAndAddress) (uuid.UUID, int32, error)

	GetUserByTg(ctx context.Context, tgID string) (*model.User, error)

	CreateProfile(ctx context.Context, req CreateProfileRequest) (*CreateProfileResponse, error)

	GetUserBySegmentID(ctx context.Context, segmentID int32) (model.UserInfoBySegmentID, error)

	// CheckSignature checks signature and nonce
	CheckSignature(ctx context.Context, req CheckSignatureReq) (bool, error)

	// GetOrCreateSecretKey returns secret key for user
	GetOrCreateSecretKey(ctx context.Context, req GetOrCreateSecretKeyReq) (string, error)
}

type StoreInterface interface {
	CheckAndUpdateNonce(ctx context.Context, userID string, newNonce int64) (bool, error)
}

type profile struct {
	profileRepo              repository.Profile
	oldUsers                 repository.OldUsers
	miningCoins              []string
	userLogs                 repository.UserLogs
	kyc                      repository.Kyc
	notificationSettings     repository.NotificationSettings
	accessSecret             string
	analyticsNotifier        notification.Analytics
	changePhoneAllowed       int
	changePhoneAllowedPeriod time.Duration
	idenfyRetryDelay         time.Duration
	referral                 repository.Referral
	email                    repository.Email
	minGetAllReferralsTake   int
	jobsClient               *jobs.APIJobsClient
	userAccountRepo          UserAccountService
	walletClient             repository.Wallet
	coinRepository           repository.Coin
	whiteLabel               repository.Whitelabel
	nonceStore               StoreInterface
}

func NewProfile(
	profileRepo repository.Profile,
	oldUsers repository.OldUsers,
	userLogs repository.UserLogs,
	notificationSettings repository.NotificationSettings,
	kyc repository.Kyc,
	accessSecret string,
	analyticsNotifier notification.Analytics,
	idenfyRetryDelay time.Duration,
	referral repository.Referral,
	email repository.Email,
	minGetAllReferralsTake int,
	jobsClient *jobs.APIJobsClient,
	userAccountRepo UserAccountService,
	walletClient repository.Wallet,
	coinRepository repository.Coin,
	whiteLabel repository.Whitelabel,
	nonceStore StoreInterface,
) *profile {
	return &profile{
		profileRepo:            profileRepo,
		oldUsers:               oldUsers,
		userLogs:               userLogs,
		kyc:                    kyc,
		accessSecret:           accessSecret,
		analyticsNotifier:      analyticsNotifier,
		idenfyRetryDelay:       idenfyRetryDelay,
		notificationSettings:   notificationSettings,
		referral:               referral,
		email:                  email,
		minGetAllReferralsTake: minGetAllReferralsTake,
		jobsClient:             jobsClient,
		userAccountRepo:        userAccountRepo,
		coinRepository:         coinRepository,
		walletClient:           walletClient,
		whiteLabel:             whiteLabel,
		nonceStore:             nonceStore,
	}
}

func (p *profile) GetUserBySegmentID(ctx context.Context, segmentID int32) (model.UserInfoBySegmentID, error) {
	userID, err := p.oldUsers.GetUserIDBySegmentID(ctx, segmentID)
	if err != nil {
		return model.UserInfoBySegmentID{}, fmt.Errorf("GetUserBySegmentID, GetUserIDBySegmentID: %w", err)
	}

	user, err := p.oldUsers.GetUserByOldID(ctx, int(userID))
	if err != nil {
		return model.UserInfoBySegmentID{}, fmt.Errorf("GetUserBySegmentID, GetUserByOldID: %w", err)
	}

	return model.UserInfoBySegmentID{
		ID:       userID,
		UUID:     user.ID,
		UserName: user.Username,
		Email:    user.Email,
	}, nil
}

func (p *profile) Create(ctx context.Context, pr *model.Profile) (int32, error) {
	// при использовании альтернативных способов регистрации пользователь не вводит пароль,
	// Будем генерировать случайный. Потом пользователь сможет сбросить пароль при необходимости.
	if len(pr.User.Password) == 0 {
		err := pr.User.SetRandomPassword()
		if err != nil {
			return 0, fmt.Errorf("set random password:  %w", err)
		}
		log.Info(ctx, "set random password for user %s", pr.User.Email)
	}
	oldID, id, err := p.oldUsers.Create(ctx, pr.User)
	if err != nil {
		return 0, fmt.Errorf("profile: create: %w", err)
	}
	pr.User.ID = id
	err = p.profileRepo.Create(ctx, pr)
	if err != nil {
		return 0, fmt.Errorf("profile: create: %w", err)
	}

	err = p.setDefaultReferralSettings(ctx, pr.User.ID, pr.User.WhiteLabelID, pr.User.NewRefID)
	if err != nil {
		log.Error(ctx, "profile.Create: %v", err)
	}

	return oldID, nil
}

func (p *profile) SaveV3(ctx context.Context, pr *model.Profile) (int32, error) {
	// do this WithinTransaction
	trx := p.oldUsers.Begin()
	var userId int32
	err := trx.WithinTransactionWithOptions(
		ctx,
		func(txCtx context.Context) error {
			id, inTxErr := p.oldUsers.SaveV2(txCtx, pr.User)
			if inTxErr != nil {
				return fmt.Errorf("profile: save 1: %w", inTxErr)
			}
			inTxErr = p.profileRepo.SaveV3(ctx, pr)
			if inTxErr != nil {
				return fmt.Errorf("profile: save 2: %w", inTxErr)
			}
			userId = id
			return nil
		},
		pgx.TxOptions{IsoLevel: pgx.Serializable},
	)
	if err != nil {
		return 0, err
	}

	err = p.setDefaultReferralSettings(ctx, pr.User.ID, pr.User.WhiteLabelID, pr.User.NewRefID)
	if err != nil {
		log.Error(ctx, "profile.Create: %v", err)
	}

	return userId, nil
}

func (p *profile) UpdatePassword(ctx context.Context, email, password string, whiteLabelID uuid.UUID) error {
	err := p.profileRepo.UpdatePassword(ctx, email, password)
	if err != nil {
		return fmt.Errorf("profile: update password: %w", err)
	}
	return nil
}

func (p *profile) GetUserByEmailAndWl(ctx context.Context, email string, whiteLabelID uuid.UUID) (*model.User, error) {
	u, err := p.profileRepo.GetUserByEmailAndWl(ctx, email, whiteLabelID)
	if err != nil {
		return nil, fmt.Errorf("GetUserByEmailAndWl: %w", err)
	}

	log.Info(ctx, "getted user %+v", u)

	if u == nil {
		return nil, nil
	}

	usr, err := p.oldUsers.GetUserByUUID(ctx, u.ID)
	if err != nil {
		return nil, fmt.Errorf("GetUserByUUID: %w", err)
	}

	if usr == nil {
		id, err := p.oldUsers.SaveID(ctx, u.ID, u.Email)
		if err != nil {
			return nil, fmt.Errorf("SaveID id: %v %v %v %w", u.ID, u.OldID, err, ErrInconsistentUserInfo)
		}
		log.Info(ctx, "fixed %v %u", id, u)
		u.OldID = id.Old
	} else {
		u.OldID = usr.OldID
	}
	return u, nil
}

func (p *profile) GetOldUserByEmailAndWl(ctx context.Context, email string, whiteLabelID uuid.UUID) (*model.User, error) {
	u, err := p.oldUsers.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("profile: get old user by email: %w", err)
	}
	if u == nil {
		return nil, nil
	}
	return u, nil
}

func (p *profile) GetAllUsers(
	ctx context.Context,
	skip, take int,
	sortField string,
	asc bool,
	searchEmail string,
) ([]*model.User, int, error) {
	// we left Phone field empty since
	// it seems that it don't required on receiver side
	users, totalCount, err := p.profileRepo.GetAllUsers(ctx, skip, take, sortField, asc, searchEmail)
	if err != nil {
		return nil, 0, fmt.Errorf("profile: get all users: %w", err)
	}
	newIDs := make([]uuid.UUID, len(users))
	for i := range newIDs {
		newIDs[i] = users[i].ID
	}
	ids, err := p.oldUsers.GetIDs(ctx, newIDs)
	if err != nil {
		return nil, 0, fmt.Errorf("profile: get all users by wl id: %w", err)
	}
	idsMap := make(map[uuid.UUID]int32)
	for i := range ids {
		idsMap[ids[i].New] = ids[i].Old
	}
	for i := range users {
		oldID, ok := idsMap[users[i].ID]
		if !ok {
			log.Error(ctx, "get all users: %v. email: %s", ErrInconsistentUserInfo, users[i].Email)
		}
		users[i].OldID = oldID
	}
	return users, totalCount, nil
}

func (p *profile) GetAllUsersByWlID(
	ctx context.Context,
	skip, take int,
	sortField string,
	asc bool,
	searchEmail string,
	whiteLabelID uuid.UUID,
) ([]*model.User, int, error) {
	// we left Phone field empty since
	// it seems that it don't required on receiver side
	users, totalCount, err := p.profileRepo.GetAllUsersByWlID(ctx, skip, take, sortField, asc, searchEmail, whiteLabelID)
	if err != nil {
		return nil, 0, fmt.Errorf("profile: get all users by wl id: %w", err)
	}
	newIDs := make([]uuid.UUID, len(users))
	for i := range newIDs {
		newIDs[i] = users[i].ID
	}
	ids, err := p.oldUsers.GetIDs(ctx, newIDs)
	if err != nil {
		return nil, 0, fmt.Errorf("profile: get all users by wl id: %w", err)
	}
	idsMap := make(map[uuid.UUID]int32)
	for i := range ids {
		idsMap[ids[i].New] = ids[i].Old
	}
	for i := range users {
		oldID, ok := idsMap[users[i].ID]
		if !ok {
			log.Error(ctx, "get all users by wl id: %v. email: %s", ErrInconsistentUserInfo, users[i].Email)
		}
		users[i].OldID = oldID
	}
	return users, totalCount, nil
}

func (p *profile) GetUserByApiKey(ctx context.Context, apiKey string) (*model.User, error) {
	const op = "service.profile.GetUserByApiKey"
	u, err := p.profileRepo.GetUserByApiKey(ctx, apiKey)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if u == nil {
		return nil, nil
	}
	usr, err := p.oldUsers.GetUserByUUID(ctx, u.ID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if usr == nil {
		return nil, fmt.Errorf("%s: %w", op, ErrInconsistentUserInfo)
	}
	u.OldID = usr.OldID
	return u, nil
}

func (p *profile) GetByUserID(ctx context.Context, userID uuid.UUID) (*model.Profile, error) {
	const op = "service.profile.GetByUserID"
	pr, err := p.profileRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if pr == nil {
		return nil, nil
	}
	u, err := p.oldUsers.GetUserByUUID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if u == nil {
		return nil, fmt.Errorf("%s: %w", op, ErrInconsistentUserInfo)
	}
	pr.User.OldID = u.OldID
	return pr, nil
}

func (p *profile) GetUserIDByOldID(ctx context.Context, oldID int32) (uuid.UUID, error) {
	user, err := p.oldUsers.GetUserByOldID(ctx, int(oldID))
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("profile: get by old user id: %w", err)
	}
	if user == nil {
		return uuid.UUID{}, fmt.Errorf("profile: get by old user id: %w", ErrInconsistentUserInfo)
	}
	return user.ID, nil
}

func (p *profile) GetByUsernames(ctx context.Context, usernames []string) ([]*model.Profile, error) {
	return p.profileRepo.GetByUsernames(ctx, usernames)
}

func (p *profile) GetByUsernamesForReferrals(ctx context.Context, usernames []string) ([]*model.Profile, error) {
	return p.profileRepo.GetByUsernamesForReferrals(ctx, usernames)
}

func (p *profile) GetOldByEmailAndWl(ctx context.Context, email string, whiteLabelID uuid.UUID) (*model.Profile, error) {
	const op = "service.profile.GetOldByEmailAndWl"
	u, err := p.oldUsers.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("%s: get by email: %w", op, err)
	}
	if u == nil {
		return nil, nil
	}
	return &model.Profile{
		User: u,
	}, nil
}

func (p *profile) GetCountUserWithWL(
	ctx context.Context,
	wlUUID uuid.UUID,
	offset, limit int32,
) ([]model.UserShortInfo, int64, error) {
	return p.profileRepo.GetCountUserWithWL(ctx, wlUUID, offset, limit)
}

func (p *profile) GetSuspended(ctx context.Context, userID uuid.UUID) (bool, error) {
	suspended, err := p.profileRepo.GetSuspended(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("getSuspended: %w", err)
	}
	return suspended, nil
}

func (p *profile) SetSuspended(ctx context.Context, userIDs []uuid.UUID, suspended bool) error {
	if err := p.profileRepo.SetSuspended(ctx, userIDs, suspended); err != nil {
		return fmt.Errorf("setSuspended: %w", err)
	}
	return nil
}

const (
	btcLowerCase = "btc"
	btcUpperCase = "BTC"
)

func (p *profile) SendEmailToChangeAddress(ctx context.Context, userID int, username, ip, coinCode, address, domain string) error {
	if coinCode == "" {
		coinCode = btcLowerCase
	}

	var coinName string
	if coinCode != btcLowerCase {
		coinName = coinCode
	}

	coin, err := p.coinRepository.GetByCode(ctx, coinCode)
	if err != nil {
		return fmt.Errorf("p.coinsCli.GetByCode: %w", err)
	}
	if coin == nil {
		return fmt.Errorf("coinCode %s nil coin", coinCode)
	}

	user, err := p.oldUsers.GetUserByOldID(ctx, userID)
	if err != nil {
		return fmt.Errorf("p.oldUsers.GetUserByOldID: %w", err)
	}

	miningAccountID, err := p.getOrCreateMiningAccount(ctx, userID, user.ID, coin.ID)
	if err != nil {
		log.Error(ctx, "SendEmailToChangeAddress: getOrCreateMiningAccount: user - %s: coin id: %d: %v", user.ID, coin.ID, err)
		return fmt.Errorf("getOrCreateMiningAccount: %w", err)
	}

	oldAddress, err := p.oldUsers.GetDefaultAddress(ctx, miningAccountID)
	if err != nil {
		return fmt.Errorf("getDefaultAddress: %w", err)
	}
	if oldAddress == address {
		return ErrBothAddressIsTheSame
	}

	err = p.userLogs.DeactivateByType(ctx, userID, fmt.Sprintf("%s%s", coinName, "address"))
	if err != nil {
		return fmt.Errorf("deactivateUserLogByType: %w", err)
	}

	token, err := p.createLogToken(userID, userID, coin.ID)
	if err != nil {
		return err
	}

	err = p.userLogs.Create(ctx, &model.UserLog{
		UserID:        int32(userID),
		ChangeType:    fmt.Sprintf("%s%s", coinName, "address"),
		IP:            ip,
		Token:         token,
		OldValue:      oldAddress,
		Value:         address,
		Active:        true,
		IsSegmentSent: true,
	})
	if err != nil {
		return fmt.Errorf("createUserLog: %w", err)
	}

	err = p.email.SendWalletChangedAddress(
		ctx,
		user.ID,
		token,
		coinName,
	)
	if err != nil {
		return fmt.Errorf("email.SendWalletChangedAddress: %w", err)
	}

	return nil
}

func (p *profile) GetCoinByCode(ctx context.Context, coinCode string) (*model.Coin, error) {
	return p.coinRepository.GetByCode(ctx, coinCode)
}

func (p *profile) UpdateMinPay(ctx context.Context, userID int, coinCode string, value float32) error {
	if coinCode == "" {
		coinCode = btcLowerCase
	}
	coin, err := p.coinRepository.GetByCode(ctx, coinCode)
	if err != nil {
		return fmt.Errorf("getCoinIDByCode: %w", err)
	}

	if !model.IsMinPayValid(value, coin.ID) {
		return ErrMinPayNotValid
	}

	user, err := p.oldUsers.GetUserByOldID(ctx, userID)
	if err != nil {
		return fmt.Errorf("p.oldUsers.GetUserByOldID: %w", err)
	}

	miningAccountID, err := p.getOrCreateMiningAccount(ctx, userID, user.ID, coin.ID)
	if err != nil {
		return fmt.Errorf("getOrCreateMiningAccount: %w", err)
	}

	if err = p.oldUsers.UpdateUserMinPay(ctx, miningAccountID, value); err != nil {
		return fmt.Errorf("updateUserMinPay: %w", err)
	}

	return nil
}

const (
	idLogTokenParam     = "id"
	userIDLogTokenParam = "userid"
	subIDLogTokenParam  = "subid"
	coinIDLogTokenParam = "coinid"
	expireLogTokenParam = "exp"
)

func (p *profile) createLogToken(userID, subID, coinID int) (string, error) {
	atClaims := jwt.MapClaims{}
	atClaims[idLogTokenParam] = userID
	atClaims[userIDLogTokenParam] = userID
	atClaims[subIDLogTokenParam] = subID
	atClaims[coinIDLogTokenParam] = coinID
	atClaims[expireLogTokenParam] = time.Now().Add(time.Hour * 24 * 30).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(p.accessSecret))
	if err != nil {
		return "", fmt.Errorf("createLogToken: %w", err)
	}
	return token, nil
}

func (p *profile) ChangeWalletAddressConfirm(ctx context.Context, _ int32, token string) (*model.ChangeWalletAddressConfirmResult, error) {
	userID, subID, coinID, err := p.getLogTokenData(token)
	log.Info(ctx, "change wallet address confirm: user id: %v. coin id: %v. sub id: %v",
		userID, coinID, subID)
	if err != nil {
		return nil, fmt.Errorf("change wallet addres confirm: %w", err)
	}
	coin, err := p.coinRepository.GetByLegacyID(ctx, coinID)
	if err != nil {
		return nil, fmt.Errorf("change wallet address confirm: %w", err)
	}
	if coin == nil {
		return nil, fmt.Errorf("change wallet address confirm: %w", ErrCoinNotFound)
	}

	if coin.Name == btcUpperCase {
		coin.Name = ""
	}

	userLog, inTxErr := p.userLogs.Get(ctx, token, strings.ToLower(coin.Name)+model.UserLogChangeTypeChangeAddress, userID, true)
	if inTxErr != nil {
		return nil, fmt.Errorf("change wallet address confirm: %w", inTxErr)
	}
	if userLog == nil {
		return nil, fmt.Errorf("change wallet address confirm: %w", ErrTokenNotFound)
	}

	coinName := coin.Name
	if coinName == "" {
		coinName = btcLowerCase
	}

	err = p.ChangeWalletAddress(ctx, userID, coinID, userLog.Value, coinName)
	if err != nil {
		return nil, err
	}

	return &model.ChangeWalletAddressConfirmResult{
		Address: userLog.Value,
		UserID:  userID,
		CoinID:  coin.Code,
	}, nil
}

func (p *profile) ChangeWalletAddress(ctx context.Context, userID int32, coinID int, address, coinName string) error {
	trx := p.oldUsers.Begin()
	return trx.WithinTransactionWithOptions(
		ctx,
		func(txCtx context.Context) error {
			inTxErr := p.oldUsers.DeleteAddress(txCtx, userID, coinID)
			if inTxErr != nil {
				return fmt.Errorf("change wallet address confirm: %w", inTxErr)
			}
			inTxErr = p.oldUsers.CreateAddress(txCtx, userID, coinID, address)
			if inTxErr != nil {
				return fmt.Errorf("change wallet address confirm: %w", inTxErr)
			}

			inTxErr = p.oldUsers.UpdateUserAccountChangedAt(txCtx, userID, coinID)
			if inTxErr != nil {
				log.Error(txCtx, "change wallet address confirm: %v", inTxErr)
			}

			email, inTxErr := p.oldUsers.GetEmailWithParentID(txCtx, userID)
			if inTxErr != nil {
				log.Error(txCtx, "change wallet address confirm: %v", inTxErr)
				return nil
			}
			segmentID, inTxErr := p.oldUsers.GetSegmentID(txCtx, userID)
			if inTxErr != nil {
				log.Error(txCtx, "change wallet address confirm: %v", inTxErr)
				return nil
			}
			p.analyticsNotifier.WithdrawalAddressChanged(coinName, email, address, segmentID)
			return nil
		},
		pgx.TxOptions{IsoLevel: pgx.Serializable},
	)
}

func (p *profile) UpdateUserIsActive(ctx context.Context, email string, active bool) error {
	err := p.profileRepo.UpdateUserIsActive(ctx, email, active)
	if err != nil {
		return fmt.Errorf("updateUserIsActive: %w", err)
	}
	return nil
}

func (p *profile) GetUserIsActive(ctx context.Context, userID uuid.UUID) (bool, error) {
	isActive, err := p.profileRepo.GetUserIsActive(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("getUserIsActive: %w", err)
	}
	nopay, err := p.oldUsers.GetNoPay(ctx, userID)
	if err != nil {
		return false, err
	}
	if isActive == true && nopay == false {
		return true, nil
	}
	return false, nil
}

func (p *profile) getLogTokenData(tokenString string) (int32, int, int, error) {
	var userID int32
	var subID, coinID int
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(p.accessSecret), nil
	})
	if err != nil {
		return userID, subID, coinID, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		if u, ok := claims[userIDLogTokenParam].(float64); ok {
			userID = int32(u)
		}
		if s, ok := claims[subIDLogTokenParam].(float64); ok {
			subID = int(s)
		}
		if c, ok := claims[coinIDLogTokenParam].(float64); ok {
			coinID = int(c)
		}

		return userID, subID, coinID, err
	}

	return userID, subID, coinID, err
}

func (p *profile) getOrCreateMiningAccount(ctx context.Context, userID int, userUuid uuid.UUID, coinID int) (int, error) {
	_, isMiningAccountMustExist := model.CoinsWithMiningAccounts[coinID]
	// miningAccountID, err := p.oldUsers.GetMiningAccountID(ctx, userID, coinID)
	miningAccountID, err := p.userAccountRepo.GetUserAccountIdByLegacyParams(
		ctx,
		int32(userID),
		int32(coinID),
		userAccountModelEnum.MiningAccountTypeID.ToInt32(),
	)
	if err != nil {
		if errors.Is(err, ErrUserAccountListIsEmpty) && !isMiningAccountMustExist {
			defaultMinPay := model.MinPayDefault[coinID]
			if defaultMinPay == 0 {
				return 0, fmt.Errorf("GetMiningAccountID: coin_id %d zero minpay", coinID)
			}

			// miningAccountID, err = p.oldUsers.CreateMiningAccount(ctx, userID, coinID, defaultMinPay)
			miningAccountID, err = p.userAccountRepo.CreateUserAccount(
				ctx,
				int32(userID),
				userUuid,
				int32(coinID),
				userAccountModelEnum.MiningAccountTypeID.ToInt32(),
				defaultMinPay,
			)
			if err != nil {
				return 0, fmt.Errorf("CreateMiningAccount: %w", err)
			}
		} else {
			return 0, fmt.Errorf("GetMiningAccountID: %w", err)
		}
	}

	if miningAccountID == 0 {
		return 0, fmt.Errorf("getOrCreateMiningAccount zero miningAccountID user_id:%d coin_id:%d", userID, coinID)
	}

	return miningAccountID, nil
}

func (p *profile) GetKycStatus(ctx context.Context, userID int) (*model.Kyc, error) {
	const op = "service.profile.GetKycStatus"
	kycStatus, lastTry, overall, docCheck, faceCheck, err := p.kyc.GetUserStatus(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &model.Kyc{
		RetryAfter:   lastTry.Add(p.idenfyRetryDelay),
		DelayMinutes: int(p.idenfyRetryDelay.Minutes()),
		Status:       model.KycStatus(kycStatus),
		IsAllowed:    kycStatus != kycApprovedStatus,
		Overall:      overall,
		DocCheck:     docCheck,
		FaceCheck:    faceCheck,
	}, nil
}

func (p *profile) SetKycStatus(ctx context.Context, userID, status int) error {
	err := p.kyc.SetUserStatus(ctx, userID, status)
	if err != nil {
		return fmt.Errorf("setUserStatus: %w", err)
	}
	return nil
}

func (p *profile) InsertKycHistory(ctx context.Context, userID int, data []byte) error {
	err := p.kyc.InsertHistory(ctx, userID, data)
	if err != nil {
		return fmt.Errorf("insertHistory: %w", err)
	}
	return nil
}

func (p *profile) CheckAppleAccount(ctx context.Context, appleID, email string) (bool, string, error) {
	exist, existedEmail, err := p.profileRepo.ExistAppleID(ctx, appleID)
	if err != nil {
		return false, "", fmt.Errorf("checkAppleAccount: %w", err)
	}
	if exist {
		return false, existedEmail, nil
	}

	bind, err := p.bindAppleID(ctx, email, appleID)
	if err != nil {
		return false, "", fmt.Errorf("bindAppleID: %w", err)
	}
	if bind {
		return false, email, nil
	} else {
		return true, "", nil
	}
}

func (p *profile) bindAppleID(ctx context.Context, email, appleID string) (bool, error) {
	updatedRows, err := p.profileRepo.UpdateAppleID(ctx, email, appleID)
	if err != nil {
		return false, fmt.Errorf("updateAppleID: %w", err)
	}

	if updatedRows == 0 {
		return false, nil
	} else {
		return true, nil
	}
}

func (p *profile) GetNotificationSettings(ctx context.Context, userID uuid.UUID) (*model.NotificationSettings, error) {
	notificationSettings, err := p.notificationSettings.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("notificationSettings.GetByUserID: %w", err)
	}
	return notificationSettings, nil
}

func (p *profile) SaveNotificationSettings(ctx context.Context, settings *model.ChangeableNotificationSettings) error {
	err := p.notificationSettings.Save(ctx, settings)
	if err != nil {
		return fmt.Errorf("notificationSettings.Save: %w", err)
	}
	return nil
}

func (p *profile) RelatedUsers(ctx context.Context, firstID, secondID uuid.UUID) (bool, error) {
	related, err := p.profileRepo.RelatedUsers(ctx, firstID, secondID)
	if err != nil {
		return false, fmt.Errorf("relatedAccounts: %w", err)
	}
	return related, nil
}

func (p *profile) GetAllSubUsers(ctx context.Context, userID uuid.UUID) ([]*model.User, error) {
	users, err := p.profileRepo.GetAllSubUsers(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("getAllSubUsers: %w", err)
	}
	return users, nil
}

func (p *profile) GetAllUserIDsByUsername(ctx context.Context) (map[string]uuid.UUID, error) {
	ids, err := p.profileRepo.GetAllUserIDsByUsername(ctx)
	if err != nil {
		return nil, fmt.Errorf("getAllUserIDsByUsername: %w", err)
	}
	return ids, nil
}

func (p *profile) GetReferrals(
	ctx context.Context,
	userID uuid.UUID,
	skip, take int,
	sortField string,
	asc bool,
) ([]*model.User, int, error) {
	referrals, count, err := p.profileRepo.GetReferrals(ctx, userID, skip, take, sortField, asc)
	if err != nil {
		return nil, 0, fmt.Errorf("getReferrals: %w", err)
	}
	return referrals, count, nil
}

func (p *profile) GetUsernamesByIDs(ctx context.Context, userIDs []uuid.UUID) (map[uuid.UUID]string, error) {
	res, err := p.profileRepo.GetUsernamesByIDs(ctx, userIDs)
	if err != nil {
		return nil, fmt.Errorf("getUsernamesByIDs: %w", err)
	}
	return res, nil
}

func (p *profile) GetEmailsByIDs(ctx context.Context, userIDs []uuid.UUID) (map[uuid.UUID]string, error) {
	res, err := p.profileRepo.GetEmailsByIDs(ctx, userIDs)
	if err != nil {
		return nil, fmt.Errorf("getEmailsByIDs: %w", err)
	}
	return res, nil
}

func (p *profile) SafeDeleteByID(ctx context.Context, userID uuid.UUID) error {
	err := p.profileRepo.SafeDeleteByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("safeDeleteByID : %w", err)
	}
	return nil
}

func (p *profile) setDefaultReferralSettings(ctx context.Context, userID, wlID, refID uuid.UUID) error {
	if err := p.referral.Create(ctx, userID, wlID, refID); err != nil {
		return fmt.Errorf("setDefaultReferralSettings.Create: user_id: %s, wlID: %s, refID: %s, err: %v",
			userID.String(),
			wlID.String(),
			refID.String(),
			err,
		)
	}

	return nil
}

func (p *profile) setReferralSettingsForSubAccount(ctx context.Context, userUUID, parentUUID, wlID, refID uuid.UUID) error {
	if err := p.referral.CreateForSubAccount(ctx, userUUID, parentUUID, wlID, refID); err != nil {
		return fmt.Errorf("setReferralSettingsForSubAccount.Create: user_id: %s, wlID: %s, refID: %s, parentID %s, err: %w",
			userUUID.String(),
			wlID.String(),
			refID.String(),
			parentUUID.String(),
			err,
		)
	}
	return nil
}

func (p *profile) createUser(ctx context.Context, u *model.User) (int32, error) {
	var err error
	u.OldID, err = p.oldUsers.CreateUser(ctx, u)
	if err != nil {
		return 0, fmt.Errorf("profile: create: %w", err)
	}

	err = p.profileRepo.Create(ctx, &model.Profile{
		User: u,
	})
	if err != nil {
		return 0, fmt.Errorf("profile: create: %w", err)
	}

	return u.OldID, nil
}

// CreateSubAccount only for create, it doesn't update user like just Create func
func (p *profile) CreateSubAccount(ctx context.Context, u *model.User) (int32, error) {
	id, err := p.createUser(ctx, u)
	if err != nil {
		return 0, fmt.Errorf("profile: create: %w", err)
	}

	err = p.setReferralSettingsForSubAccount(ctx, u.ID, u.ParentID, u.WhiteLabelID, u.NewRefID)
	if err != nil {
		log.Error(ctx, "profile.CreateSubAccount: %v", err)
	}

	return id, nil
}

func (p *profile) SaveV4(ctx context.Context, u *model.User) (int32, error) {
	const op = "service.profile.SaveV4"

	// do this WithinTransaction
	trx := p.oldUsers.Begin()
	var userId int32
	err := trx.WithinTransactionWithOptions(
		ctx,
		func(txCtx context.Context) error {
			id, inTxErr := p.oldUsers.SaveV2(txCtx, u)
			if inTxErr != nil {
				return fmt.Errorf("%s: save to old users: %w", op, inTxErr)
			}
			inTxErr = p.profileRepo.SaveUser(ctx, u)
			if inTxErr != nil {
				return fmt.Errorf("%s: save to profile: %w", op, inTxErr)
			}
			userId = id
			return nil
		},
		pgx.TxOptions{IsoLevel: pgx.Serializable},
	)
	if err != nil {
		return 0, err
	}

	err = p.setDefaultReferralSettings(ctx, u.ID, u.WhiteLabelID, u.NewRefID)
	if err != nil {
		log.Error(ctx, "%s: %v", op, err)
	}

	return userId, nil
}

func (p *profile) UpdateRefID(
	ctx context.Context,
	oldUserID int32,
	userID uuid.UUID,
	refID int32,
	newRefID uuid.UUID,
) error {
	log.Info(ctx, "profile service: user_id: %s old_user_id: %d update ref_id: %s", userID, oldUserID, refID)
	// do this update WithinTransaction
	trx := p.oldUsers.Begin()

	err := trx.WithinTransactionWithOptions(
		ctx,
		func(txCtx context.Context) error {
			inTxErr := p.oldUsers.UpdateRefID(txCtx, oldUserID, refID)
			if inTxErr != nil {
				return fmt.Errorf("profile service: update ref_id in old_users: %w", inTxErr)
			}
			inTxErr = p.profileRepo.UpdateRefID(ctx, userID, refID, newRefID)
			if inTxErr != nil {
				return fmt.Errorf("profile service: update ref_id in profile: %w", inTxErr)
			}

			referralErr := p.referral.SetReferralUUID(ctx, userID, newRefID)
			if referralErr != nil {
				return fmt.Errorf("profile service: set referral uuid in referral: %w", referralErr)
			}

			return nil
		},
		pgx.TxOptions{IsoLevel: pgx.Serializable},
	)

	return err
}

func (p *profile) WasReferralLinkGenerated(ctx context.Context, userID uuid.UUID) (bool, error) {
	flg, err := p.profileRepo.GetFlagReferralLinkGenerated(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("get flag referral link generated: %w", err)
	}
	return flg, nil
}

func (p *profile) SetFlagReferralLinkGenerated(ctx context.Context, userID uuid.UUID, value bool) error {
	err := p.profileRepo.SetFlagReferralLinkGenerated(ctx, userID, value)
	if err != nil {
		return fmt.Errorf("set flag referral link generated: %w", err)
	}
	return nil
}

func (p *profile) UpdateProfile(ctx context.Context, u *model.User) error {
	// do this WithinTransaction
	trx := p.oldUsers.Begin()
	err := trx.WithinTransactionWithOptions(
		ctx,
		func(txCtx context.Context) error {
			inTxErr := p.oldUsers.UpdateUser(txCtx, u)
			if inTxErr != nil {
				return fmt.Errorf("profile: old users: update user: %w", inTxErr)
			}
			inTxErr = p.profileRepo.UpdateUser(ctx, u.ID, func(user *model.User) error {
				*user = *u
				return nil
			})
			if inTxErr != nil {
				return fmt.Errorf("profile: update user: %w", inTxErr)
			}
			return nil
		},
		pgx.TxOptions{IsoLevel: pgx.Serializable},
	)
	return err
}

func (p *profile) SetUserAttributes(ctx context.Context, userID uuid.UUID, attrs model.UserAttributes) error {
	var u *model.User
	err := p.profileRepo.UpdateUser(ctx, userID,
		func(user *model.User) error {
			user.SetAttributes(attrs)
			u = user
			return nil
		})
	if err != nil {
		return fmt.Errorf("profile: update user: %w", err)
	}

	err = p.oldUsers.UpdateUser(ctx, u)
	if err != nil {
		return fmt.Errorf("profile: old users: update user: %w", err)
	}

	return nil
}

func (p *profile) SetTimezone(ctx context.Context, userID uuid.UUID, timezone string) error {
	u, err := p.oldUsers.GetUserByUUID(ctx, userID)
	if err != nil {
		return fmt.Errorf("getID: %w", err)
	}
	err = p.oldUsers.SetTimezone(ctx, int(u.OldID), timezone)
	if err != nil {
		return fmt.Errorf("setTimezone: %w", err)
	}
	return nil
}

func (p *profile) SetLanguage(ctx context.Context, userID uuid.UUID, language string) error {
	const op = "service.profile.SetLanguage"
	// do this WithinTransaction
	trx := p.oldUsers.Begin()
	err := trx.WithinTransactionWithOptions(
		ctx,
		func(txCtx context.Context) error {
			u, inTxErr := p.oldUsers.GetUserByUUID(txCtx, userID)
			if inTxErr != nil {
				log.Error(ctx, "%s: get old id. user id - %s: %w", op, userID, inTxErr)
				return fmt.Errorf("%s: failed to get old id", op)
			}
			inTxErr = p.oldUsers.SetLanguage(txCtx, int(u.OldID), language)
			if inTxErr != nil {
				log.Error(ctx, "%s: set language in old repo. user id - %s: %w", op, userID, inTxErr)
				return fmt.Errorf("%s: failed to set language", op)
			}
			inTxErr = p.profileRepo.SetLanguage(ctx, userID, language)
			if inTxErr != nil {
				log.Error(ctx, "%s: set language in profile. user id - %s: %w", op, userID, inTxErr)
				return fmt.Errorf("%s: failed to set language", op)
			}
			return nil
		},
		pgx.TxOptions{IsoLevel: pgx.Serializable},
	)

	return err
}

func (p *profile) GetAddresses(ctx context.Context, userID uuid.UUID) ([]*model.Address, error) {
	coinCodesMining, err := p.coinRepository.GetMiningCoins(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed get mining coin serivce.profile.GetAddresses: %w", err)
	}

	address, err := p.oldUsers.GetAddresses(ctx, userID, coinCodesMining)
	if err != nil {
		return nil, fmt.Errorf("failed get addresses serivce.profile.GetAddresses: %w", err)
	}

	return address, nil
}

func (p *profile) SoftDeleteSubUser(ctx context.Context, subUserID, newParentID uuid.UUID) error {
	const op = "service.profile.SoftDeleteSubUser"
	parent, err := p.oldUsers.GetUserByUUID(ctx, newParentID)
	if err != nil {
		log.Error(ctx, "%s: get ids: new parent - %s: %v", op, newParentID, err)
		return fmt.Errorf("%s: get ids", op)
	}
	subuser, err := p.oldUsers.GetUserByUUID(ctx, subUserID)
	if err != nil {
		log.Error(ctx, "%s: get ids: subuser - %s: %v", op, subUserID, err)
		return fmt.Errorf("%s: get ids", op)
	}
	if subuser == nil {
		return fmt.Errorf("%s: subuser is nil", op)
	}

	profile, err := p.profileRepo.GetByUserID(ctx, newParentID)
	if err != nil || profile == nil || profile.User == nil {
		return fmt.Errorf("%s: get by user id: %w", op, err)
	}
	newEmail := model.NewSubUserEmailByParentEmail(profile.User.Email)
	err = p.profileRepo.SoftDeleteSubUser(
		ctx,
		subUserID,
		newParentID,
		newEmail,
	)
	if err != nil {
		return fmt.Errorf("%s: soft delete in profile db: %w", op, err)
	}

	err = p.oldUsers.SoftDeleteSubUser(
		ctx,
		subuser.OldID,
		parent.OldID,
		newEmail,
	)
	if err != nil {
		// пропускаем ошибку, т.к. не поддерживаем соглсованность данных
		// profile остается точкой правды
		log.Error(ctx, "%s: soft delete in emcd db: %v", op, err)
	}
	return nil
}

func (p *profile) GetUsersByUUIDs(ctx context.Context, userUUIDs []uuid.UUID) ([]model.UserShortInfo, error) {
	res, err := p.profileRepo.GetUsersByUUIDs(ctx, userUUIDs)
	if err != nil {
		return nil, fmt.Errorf("GetUsersByUUIDs: %w", err)
	}
	return res, nil
}

func (p *profile) CreateSubUser(ctx context.Context, parentUUID uuid.UUID, username string, addresses []*model.CoinAndAddress) (uuid.UUID, int32, error) {
	uniqueUsername, err := p.profileRepo.IsUniqueUsername(ctx, username)
	if err != nil {
		return uuid.Nil, 0, fmt.Errorf("isUniqueUsername: %w", err)
	}
	if !uniqueUsername {
		return uuid.Nil, 0, ErrNotUniqueUsername
	}

	parent, err := p.profileRepo.GetByUserID(ctx, parentUUID)
	if err != nil {
		return uuid.Nil, 0, fmt.Errorf("getByUserID: %w", err)
	}

	newSubUserEmail := model.NewSubUserEmailByParentEmail(parent.User.Email)

	fee := p.getSpecRefFee(ctx, int32(parent.User.RefID))

	oldParentIDs, err := p.oldUsers.GetIDs(ctx, []uuid.UUID{parentUUID})
	if err != nil {
		return uuid.Nil, 0, fmt.Errorf("oldUsers.GetIDs: %w", err)
	}

	if len(oldParentIDs) == 0 {
		return uuid.Nil, 0, fmt.Errorf("oldUsers.GetIDs: no parents users")
	}

	newSubID := uuid.New()
	usr := model.User{
		ID:           newSubID,
		ParentID:     parentUUID,
		OldParentID:  &oldParentIDs[0].Old,
		Email:        newSubUserEmail.String(),
		Username:     username,
		Password:     "subaccount",
		RefID:        parent.User.RefID,
		NewRefID:     parent.User.NewRefID,
		IsActive:     true,
		WhiteLabelID: uuid.Nil,
		CreatedAt:    time.Now().UTC(),
	}
	subAccountID, err := p.CreateSubAccount(ctx, &usr)
	if err != nil {
		return uuid.Nil, 0, fmt.Errorf("createSubUser: %w %v email: %s", err, usr, newSubUserEmail)
	}

	newSubAcc, err := p.profileRepo.GetByUserID(ctx, newSubID)
	if err != nil {
		return uuid.Nil, 0, fmt.Errorf("getByUserID: %w", err)
	}

	if newSubAcc == nil || newSubAcc.User == nil {
		return uuid.Nil, 0, fmt.Errorf("newSubAcc == nil || newSubAcc.User == nil: %w", ErrNotFound)
	}

	err = p.oldUsers.ApplyParentPromoCodes(ctx, int(parent.User.OldID), subAccountID)
	if err != nil {
		return uuid.Nil, 0, fmt.Errorf("applyParentPromoCodes: %w", err)
	}

	txAccounts, err := p.oldUsers.Begin().Runner(ctx).Begin(ctx)
	if err != nil {
		return uuid.Nil, 0, fmt.Errorf("accounts.Begin: %w", err)
	}
	_, err = p.createAccounts(ctx, txAccounts, subAccountID, newSubID, fee, 0, addresses)
	if err != nil {
		return uuid.Nil, 0, fmt.Errorf("createUserAccounts error: %w", err)
	}
	err = txAccounts.Commit(ctx)
	if err != nil {
		log.Error(ctx, "createSubUser: commit error: %w", err)
		return uuid.Nil, 0, err
	}

	return newSubID, subAccountID, nil
}

func (p *profile) getSpecRefFee(ctx context.Context, refID int32) float64 {
	fee := defFee
	specRefFee, err := p.oldUsers.GetSpecRefFee(ctx, refID)
	if err != nil {
		increasedFee, ok := model.IncreasedFeeRefMap[refID]
		if ok {
			fee = increasedFee
		}
	} else {
		fee = specRefFee
	}
	return fee
}

func (p *profile) createAccounts(ctx context.Context, tx pgx.Tx, userID int32,
	userIdNew uuid.UUID, fee float64, refID int32, addresses []*model.CoinAndAddress,
) ([]int32, error) {
	createdAccs, err := p.createUserAccounts(ctx, userID, userIdNew, fee)
	if err != nil {
		return nil, fmt.Errorf("accounts.CreateUserAccounts: %w", err)
	}

	miningAccounts := make([]int32, 0)
	for _, acc := range createdAccs {
		switch acc.AccountTypeID.ToInt32() {
		case userAccountModelEnum.ReferralAccountTypeID.ToInt32():
			// TODO: remove it because all commission save to referral service
			if acc.CoinID == 0 {
				log.Error(ctx, "attempt to create ref account without coin id")
				continue
			}
			err = p.oldUsers.CreateReferralAccounts(ctx, tx, acc.ID, acc.CoinID)
			if err != nil {
				return nil, fmt.Errorf("accounts.CreateReferralAccounts: %w", err)
			}
		case userAccountModelEnum.MiningAccountTypeID.ToInt32():
			miningAccounts = append(miningAccounts, acc.ID)
			if refID != model.GarantexUserID {
				err = p.oldUsers.CreateAccountsPool(ctx, tx, acc.ID)
				if err != nil {
					return nil, fmt.Errorf("accounts.CreateAccountsPool: %w", err)
				}
			}

			autoPayAddress := ""
			if acc.Address.Valid && acc.Address.String != "" {
				autoPayAddress = acc.Address.String
			}

			for _, addr := range addresses {
				if strings.EqualFold(acc.CoinNew.String, addr.Coin) && addr.Address != "" {
					autoPayAddress = addr.Address
					break
				}
			}

			if autoPayAddress != "" {
				errInner := p.oldUsers.CreateAutoPayAddress(ctx, tx, acc.ID, acc.CoinID, autoPayAddress)
				if errInner != nil {
					return nil, fmt.Errorf("accounts.CreateAutoPayAddress: %w", errInner)
				}
			}

		default:
			log.Error(ctx, "not handled account type %d", acc.AccountTypeID)
		}
	}
	return miningAccounts, nil
}

func (p *profile) createUserAccounts(ctx context.Context, userID int32, userIdNew uuid.UUID, fee float64) (userAccountModel.UserAccounts, error) {
	// pass address wallets
	userAccountsRequest := accountingRepository.GenerateDefaultUsersAccounts(userID, userIdNew, 0.0, fee)

	for _, acc := range userAccountsRequest {
		if acc.AccountTypeID.ToInt32() != userAccountModelEnum.MiningAccountTypeID.ToInt32() || len(acc.Address.String) > 0 {
			continue
		}

		addrResponse, err := p.walletClient.GetAddress(ctx, userID, acc.CoinNew.String)
		if err != nil {
			log.Error(ctx, "walletClient.GetAddress: %w", err)
			continue
		}

		addrs := addrResponse.GetAddresses()
		if len(addrs) == 0 {
			log.Error(ctx, "walletClient.GetAddress: no addresses for coin %s", acc.CoinNew.String)
			continue
		}

		acc.Address = sql.NullString{
			String: addrs[0].Address,
			Valid:  true,
		}
	}

	userAccounts, err := p.userAccountRepo.CreateUserAccounts(ctx, userID, userIdNew, userAccountsRequest)
	if err != nil {
		log.Error(ctx, "error creating user accounts %s, %d, %s", err.Error(), userID, userIdNew)
		return nil, err
	}

	return userAccounts, nil
}

func (p *profile) GetUserByTg(ctx context.Context, tgID string) (*model.User, error) {
	log.Info(ctx, "tg id: %s", tgID)
	u, err := p.profileRepo.GetUserByTg(ctx, tgID)
	if err != nil {
		return nil, fmt.Errorf("GetUserByTg: %w", err)
	}

	log.Info(ctx, "user from database: %+v", u)

	if u == nil {
		return nil, nil
	}

	usr, err := p.oldUsers.GetUserByUUID(ctx, u.ID)
	if err != nil {
		return nil, fmt.Errorf("GetUserByUUID: %w", err)
	}

	if usr == nil {
		id, err := p.oldUsers.SaveID(ctx, u.ID, u.Email)
		if err != nil {
			return nil, fmt.Errorf("SaveID id: %v %v %v %w", u.ID, u.OldID, err, ErrInconsistentUserInfo)
		}
		log.Info(ctx, "fixed %v %u", id, u)
		u.OldID = id.Old
		return u, nil
	}
	u.OldID = usr.OldID
	return u, nil
}
