# Referral

## Introduction.
> Service is an implementation of a new referral program. It provides a complete set of features for managing and monitoring the program. Exclusively utilizes a separate `referral` database, ensuring that only this service has access, thereby promoting data security and integrity. Other services should interact with this service through gRPC, ensuring efficient and safe inter-service communication.

## Functionality:
> This section outlines the primary functionality of the service. It describes the roles and responsibilities of various interfaces within the system

> - **DefaultSettings**: Management of the referrals default settings. Provides methods for creating, updating, deleting, and retrieving the default settings.
> - **DefaultWhitelabelSettings**: Management of the whitelabel referrals default settings. Allows creating, updating, deleting, and retrieving the settings data.
> - **Referral**: Management of referrals. It includes creating, updating, deleting, retrieving information about a specific referral, viewing the list of referrals, and viewing referral history.
> - **Reward**: Calculation of rewards for referrals based on settings.

## Usage

> This service is used in the following areas:
> - P2P: https://code.emcdtech.com/emcd/fiat/p2p
> - Mining: todo

## Additional Information

>The service uses a PostgreSQL trigger to maintain a log of changes to referral data. This trigger is activated upon updates to the `referral.referrals` table. It logs the previous state of any updated records, ensuring the preservation of historical data.