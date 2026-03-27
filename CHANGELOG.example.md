# Change Log
All notable changes to this project will be documented in this file.
This project adheres to [Semantic Versioning](http://semver.org/).

## [3.81.0] - 2026-03-26 (86.8%)(Dev)
### Deployment
- Make sure to add new environment variable `KNOX_SERVICE_PLUGIN_KEY`
### Added
- Added setting the licence key in the knox service plugin's managed config
### Security
- Updated `google.golang.org/grpc` from `v1.72.2` to `v1.79.3`

## [3.80.0] - 2026-03-16 (86.8%)(Dev)(Prod)
### Added
- Added sonnet 4.6 model to scribe
- Added sonnet 4.5, haiku 4.5 and opus 4.5 to scribe
- Added error handling if dynamic prompt fails to be parsed
- Added an option on the enterprise to enable/disable the Samsung clock
- Added enabling and managing the knox service plugin
- Added setting Read Notification permission in the service plugin's managed config

## [3.79.1] - 2026-03-03 (86.8%)(Dev)(Prod)
### Changed
- Replaced Gorm's logging with a structured logging implementation.

## [3.79.0] - 2026-03-03 (86.8%)(Dev)(Prod)
### Added
- Added pending index for `webhook_events` table
- Added the agent app to be always included in the notification whitelisted apps
- Added logging message id and applications list md5 hash in the device application lists
- Added preventing the device application list being saved if there is no change from the previous device application list
- Added Knox Service Plugin managed configuration policy creator into `files/ksp` folder.
- Added disabling showing parameters in the SQL logs
### Changed
- Changed `scribe_jobs` table to mb4 collation for hte odd characters that can come through in the prompts and summaries.
- Removed a few places where `logprettyprint` and generic logs have been left in
- Changed the message updating process in the device details processing function in the queue service to use `Updates` instead of `Save` on the message
- Moved the device application list processing from queue service to the message service
- Refactored message processing from the cron handler to the message service
- After discussion with Geroge, added the Samsung Clock back in.
### Security
- Updated `github.com/cloudflare/circl` from `v1.6.1` to `v1.6.3`

## [3.78.2] - 2026-03-03 (86.8%)(Dev)(Prod)
### Added
- Added NotificationsEnabled flag to enterprise applications
- Added a managed config that passes a list of applications that have notifications enabled

## [3.78.1] - 2026-03-03 (86.8%)(Dev)(Prod)
### Changed
- Remove the Samsung clock from the homescreen due to music player prompt.

## [3.78.0] - 2026-02-25 (86.8%)(Dev)(Prod)
### Added
- Updated `ManagedAccountUpgrade` process to automatically sync device when updated. Added test to auto log owner in policy sync process.
- Added an option to set an app as a credential provider
- Added an option to set the credential provider policy default on an enterprise
- Added the ability to control the app blacklist an api key.
### Fixed
- Fixed the reset password endpoint not logging failed login when the passed token is invalid
- Role update endpoint wasn't validating provided permissions, added a check to see if they exist
### Security
- Update indirect package `filippo.io/edwards25519` from `v1.1.0` to `v1.1.1` for security fix
- Updated `github.com/go-git/go-git/v5` from `v5.13.1` to `v5.16.5` for dependabot issue.
- Replaced depreciated package `github.com/go-shiori/go-readability` used for BBC Feed processing with `codeberg.org/readeck/go-readability/v2`

## [3.77.0] - 2026-02-17 (86.8%)(Dev)(Prod)
### Added
- Added notificationServerUrl and notification tokens in to the system.
- Added a notification token to the agent application and to additional applications that have mdm auth enabled.
- Added notification token and server URL to oauth if hey have packages added ot their oauth client.
### Changed
- Allow the clock on the homescreen now we're using the knoxSDK

## [3.76.0] - 2026-02-10 (86.8%)(Dev)(Prod)
### Changed
- refactored the create owner token so the options are parsed in the function that uses it.
### Fixed
- Youtube code was trying to use items without checking if they were nil. Added protection so they can't panic.

## [3.75.0] - 2026-02-05 (86.8%)(Dev)(Prod)
### Added
- Added the ability for an enterprise to use Email 2FA for logging in
- Added `deviceUid` to app manged config if mdm auth enabled. Added to the agent app every time. Required for notifications.
- Added action filter for the owner activity endpoints.
### Changed
- Improved logging in scribe upload endpoint to include job ID
### Fixed
- Fixed device application overrides value in database being set to the string "null" rather than an empty string when the last override is deleted.
- Updated `ArticleUpdateRequest` to limit `body` and `excerpt` to `65000` as it was set to `500000` which is larger than the `TEXT` type in MySQL/MariaDB

## [3.74.0] - 2026-01-22 (86.8%)(Dev)(Prod)
### Added
- Added an option to enable managed google account upgrade on the device and optionally provide account email
- Added the ability to blacklist and whitelist PurplePlay channels by changing some options on the api key.
### Changed
- Changed how owner access keys are made, making it easier to add new values to the token in the future.
### Fixed
- Fixed purpledocs local domain not being a sub-domain of localmdm preventing the /files cookie from being used.
### Security
- Updated `google.golang.org/api` version `v0.230.0` to `v0.236.0` as it was needed to access the new policy setting `WorkAccountSetupConfig`

## [3.73.1] - 2026-01-12 (86.8%)(Dev)(Prod)
### Fixed
- Make sure the clear app cache feature works for the BT user and group UID override.
- Updated `Summary` general to log a warning using structured logging instead of an error.

## [3.73.0] - 2026-01-09 (86.8%)(Dev)(Prod)
### Added
- Added clear application cache feature, a table for tracking progress an endpoint to get the progress of this request.
### Changed
- Changed logging to use structured logging with slog. Also added some additional information to these logs, as well as some helper functions to create logs.
- Binding type for a new enterprise is now set to `MANAGED_GOOGLE_PLAY_ACCOUNTS_ENTERPRISE`, instead of setting it from the returned AM enterprise

## [3.72.0] - 2025-12-08 (86.8%)(Dev)(Prod)
### Added
- Added Binding Type to the Enterprises table and an option to upgrade an enterprise to the Managed Google Domain type

## [3.71.1] - 2025-12-08 (86.8%)(Dev)(Prod)
### Changed
- Made sure to exclude the agent app from application overrides
### Fixed
- Fixed enrollment groups not having the default keyguard_option 'DISABLE_ALL' preventing updating

## [3.71.0] - 2025-11-27 (86.8%)(Dev)(Prod)
### Added
- Added an option to set managed configuration on the device application override and an option to enable/disable it
- Added CIDR range validation to the enterprise proxy update endpoint to restrict allowed ranges to /28 or smaller
- Added setting device uid and enterprise uid on the message, when the message is of type Command
- Added an option to start the api with the mocked out YouTube service
- Added tests for bind failures in handlers increasing coverage
- Added the ability for keyguard_options to be configured on a group. Default is `DISABLE_ALL`
### Changed
- Renamed Track Override to Application Override
- Refactored handlers to return ErrorResponse instead of just the error for bind error handling
### Fixed
- Fixed TestOwnerUpdatePolicy (owner_test) by initialising the wait group before running, preventing race condition that occurred when ran individually.
- Fixed TestUpdateTask (task_test) by enabling the notification_enabled flag on the default enterprise allowing the notification database check to pass.
- Fixed scribe seeder by adding Unscoped() so it checks for deleted items before trying to create.
- Fixed webhook events always being processed in chronological order
### Security
- Updated `	golang.org/x/crypto` from `v0.37.0` to `v0.45.0`

## [3.70.1] - 2025-11-17 (86.1%)(Dev)(Prod)
### Deployment
- NOTE: Make sure to deploy the SSO server at the same time
### Added
- Added Purple Play dev endpoint to CORS config.
- Added skipping logging health checks from AWS in the SSO server
### Changed
- Refactored handlers in restricted routes so that the handler is next to the routes they handle

## [3.70.0] - 2025-11-10 (86.1%)(Dev)(Prod)
### Added
- Added order cancellation functionality - implements the ability for owners to cancel their orders
- Added logs when a device is removed from an owner using Delete in attachment handler.

## [3.69.0] - 2025-11-07 (86.1%)(Dev)(Prod)
### Added
- Added notifications system
- Added a cron job for sending notification alert emails to users
### Changed
- In production, postman `from` name updated to `Purple MDM` rather than `Madepurple Support` at the request of Rampton.

## [3.68.0] - 2025-11-04 (86.1%)(Dev)(Prod)
### Deployment
- Need to manually remove `system.log` from `owner_activities` table
### Added
- Added MOCK_SUMMARY_SERVICE environment variable to allow the mocking of the summary service when running locally.
- Added unit test for the main server
- Added index for the device activities table on the enterprise and device fields.
- Added an option to exclude global blacklist URLs from the kiosk sites
### Changed
- Updated device activities and device owner activities repository methods to have the enterprise uid field passed through (to make use of indexes on the table).
- Removed `system.log` filter from `owner_activities` repo as we no long log these items in this table. Handler for app now rejects them.
- Added error logging to google policy updates.
- Added an alternative function to the date_range package to just return dates, as their one returned nanoseconds which wasn't needed. This package needs removing eventually.

## [3.67.2] - 2025-10-27 (85.9%)(Dev)(Prod)
### Added
- Added SSO flag to app update and response
- Protected the forgot password endpoint against timing attacks by sending emails in a goroutine.
### Changed
- Added permission to devices/stats endpoint to check for Device.Read

## [3.67.1] - 2025-10-24 (85.9%)(Dev)(Prod)
### Added
- Added owner identifier and order ID for shopping into the task name for the tracker
- Added search by modelID to the tracker.

## [3.67.0] - 2025-10-23 (85.9%)(Dev)(Prod)
### Added
- Added lookup table for our model names to bedrock model names
- Added migration to update existing model names
### Changed
- Updated Go in the docker containers from `v1.23` to `v1.25.3` for security fixes.
- Updated the databases to match production from `mysql:5.7` to `mariadb:10.11.10`. Old image failed on the latest linux kernal and no updates.
- Fixed the `hashid` package as the number of allocations has reduced due to updates in the Go version. We now check for a maximum number of allocations instead.
- Updated env.dist to include a default cookie secret
- Validate prompt model on scribe prompt create & update
- Modified all hardcoded prompts to use new AIModel values
- Modified BasePrompt to translate between our AIModel values and the bedrock models
- Updated version of swag and swaggo used in the dev Dockerfile
- Changed collation of webhook_events table to support unicode emojis
### Removed
- Removed old General worker

## [3.66.2] - 2025-10-13 (85.9%)(Dev)(Prod)
### Fixed
- Allow the agent add to be added to the override list per device for track overrides

## [3.66.1] - 2025-10-13 (85.9%)(Dev)(Prod)
### Fixed
- Fixed override app package names not being included in the login response

## [3.66.0] - 2025-10-09 (85.9%)(Dev)(Prod)
### Added
- Added an option to override the video access on the individual purple play channels on enterprise

## [3.65.0] - 2025-10-02 (85.9%)(Dev)(Prod)
### Added
- Added kiosk site override on the applications
- Added internal name to the kiosk sites
- Added kiosk application flag to kiosk sites
### Changed
- The kiosk sites now also need the kiosk application flag to be set to true in order to be returned during the app login

## [3.64.0] - 2025-09-26 (85.9%)(Dev)(Prod)
### Added
- Added an option to disable battery optimization on an application

## [3.63.0] - 2025-09-24 (85.9%)(Dev)(Prod)
### Added
- Added `IsASCIIPrintableId` to the `str` package so we can check IDs passed through as valid printable characters as the string collation is latin
- Added `rot8` to the `osconfigs` table
- Added rot8 link to the os config response
- Added setting fake google credentials in the base test to fix failing tests that are using the real enterprise service and mocking http transport
### Changed
- Set custom IP extractor to extract only the trustable IP address from the last AWS proxy

## [3.62.1] - 2025-09-22 (85.9%)(Dev)(Prod)
### Fixed
Fixed a bug on the order page where the switch account enterprise ID was not being used

## [3.62.0] - 2025-09-22 (85.9%)(Dev)(Prod)
### Added
- Added UID search for enterprises into admin list endpoint
### Added
- Added `dev.scribe.purplemdm.com` and `scribe.purplemdm.com` to CORS

## [3.61.0] - 2025-09-11 (85.9%)(Dev)(Prod)
### Added
- Added custom responses to the http transport mock, so that we can check/validate data in individual requests
- Added a test to the update group tests that checks if the track override is applied correctly to the applications on the devices
- Added custom User Agent when fetching RSS feeds and Article Content
- Added onboard keyboard's width/height fields to the MPOS Policy
- Added ability for to admin users to clear user's failed logins .
- Added track override capability for agent applications
- Added russian support to the hub
### Changed
- Changed the Fetch Article Content tests so that they use the http transport mock
- Updated handlers so all errors when reading hte BBC website are now warnings
- Updated the `error` to `warning` in the file reader to remove the alerts for this like this in the scribe upload: `Warning: unable to read file info: error parsing multipart form: unexpected EOF`

## [3.60.4] - 2025-09-04 (85.9%)(Dev)(Prod)
### CHanged
- Exclude `bbc.co.uk/programmes` as these appear on some feeds, and then disappear in 5 minutes. This was causing lots of errors in the logs.

## [3.60.3] - 2025-09-04 (85.9%)(Dev)(Prod)
### Changed
- Updated Scribe upload error message `Error reading file info` to a warning `Warning: unable to read file info`, as it happens all the time and no one is reporting any issues. Chirs investigated twice and said he doesn't know what it could be, but most likely is just a wifi disconnect or network issue when uploading.
### Added
- Added `bbc.co.ukundefined` and `bbc.comundefined` to the list of blacklisted paths in the rss feed service, so that we can filter out the feed items with broken links

## [3.60.2] - 2025-09-03 (85.9%)(Dev)(Prod)
### Added
- Added a new OwnerActivity permissions and protested the existing `/owners/activities` and `/owners/:ownerUid/activities` endpoints with it.

## [3.60.1] - 2025-08-29 (85.9%)(Dev)(Prod)
### Fixed
- Fixed the `deleted_at` field in machine responses and document response, coming back with 0000-00-00 rather than null if it's not deleted.
- Fixed date range code that wouldn't work at the end of certain months due to implementation

## [3.60.0] - 2025-08-22 (85.9%)(Dev)(Prod)
### Added
- Added external endpoint to create scribe jobs
- Added external endpoint to upload files for scribe jobs

## [3.59.2] - 2025-08-15 (85.9%)(Dev)(Prod)
### Changed
- remove additional `error` in chain for news feeds to stop alerts.

## [3.59.1] - 2025-08-14 (85.9%)(Dev)(Prod)
### Added
- Added tests for the news feeds content endpoint
### Changed
- Changed error log printing in the get news feeds handler to a warning if the error contains `context deadline exceeded`
### Removed
- Removed error log printing in the rss feed service if rss feed items fetch fails

## [3.59.0] - 2025-08-09 (85.9%)(Dev)(Prod)
### Added
- Added an option to set/edit/delete track overrides on the device
- Added an option to sync policy on the device
### Changed
- Changed log message when there is a http error for the news feed so it just marks it as an issue and won't be alerted in Cloudwatch

## [3.58.8] - 2025-08-07 (85.9%)(Dev)(Prod)
### Changed
- For the top 10 application graph, pull the app names from the application on the enterprise after getting the top 10. This is so I don't have to put the description, where it was getting hte information from, into the index.  Added an index for this query too, to speed it up. NOTE: Should crunch number daily instead
### Fixed
- Fixed the top tags endpoint to also return data from the android devices. NOTE: This endpoint is bad and will need to be refactored in the future.

## [3.58.7] - 2025-08-01 (85.9%)(Dev)(Prod)
### Fixed
- Fixed the issue with global proxy exclude lists duplicating when saving proxy exclude lists by adding an ID check in uuidV7's `BeforeCreate`, so it doesn't create a new ID if ID is already set

## [3.58.6] - 2025-07-22 (85.9%)(Dev)(Prod)
### Changed
- Updated owner accessibility language column to support more than 4 characters
### Fixed
- Fixed the issue with machine logs and API key field. Was missing a foreignKey.

## [3.58.5] - 2025-07-20 (85.8%)(Dev)(Prod)
### Fixed
- Quick fix for Google pub sub messages not getting `ack` messages. We're reading them twice for some reason. If there is a duplicate DB error, then mark as acknowledged and return.

## [3.58.4] - 2025-07-17 (85.8%)(Dev)(Prod)
### Added
- If BT send though and IEP of 0, we now automatically default it to 2

## [3.58.3] - 2025-07-16 (85.8%)(Dev)(Prod)
### Fixed
- Removed the `order by id asc message` on message process cron jon, as it was stopping the index from being used!

## [3.58.2] - 2025-07-14 (85.8%)(Dev)(Prod)
### Fixed
- Updated SSO server to accept unique owner UIDs for BT, up to 32 characters long. All other SSO users will use their UUID V4 generated ID's.

## [3.58.1] - 2025-07-09 (85.8%)(Dev)(Prod)
### Fixed
- Updated a couple of endpoints that were passing the response directly to gorm rather than using the DTO pattern. They were returning null rather than an empty array `[]`
- Fixed a couple of responses that return `deleted_at` as it was using a gorm.DeletedAt type rather than a regular time.Time

## [3.58.0] - 2025-07-08 (85.8%)(Dev)(Prod)
### Deployment
- Coverage dropped as more errors are being handled as part of the DB change, but some of those aren't being triggered in tests.
### Changed
- Updated Gorm to `gorm.io/gorm v1.30.0` from the old `github.com/jinzhu/gorm` package. Updated code to reflect those changes.

## [3.57.1] - 2025-07-08 (85.9%)(Dev)(Prod)
### Fixed
- Fixed developer settings which were removed in the policy, set by the group as part of force unlock work.

## [3.57.0] - 2025-07-07 (85.9%)(Dev)(Prod)
### Added
- Device command logging
- Canceling non-executed device commands of the same type before issuing a new one
- Added endpoint for getting user activities by the uid.
### Changed
- Moved the `LockDevice`, `ResetPasswordDevice` and `RebootDevice` functions to the device service so that it's easier to extensively test the command logging and canceling
### Fixed
- Added `&timeTruncate=1s` to the DB connection string to make sure it removes the precision from the time.

## [3.56.1] - 2025-06-24 (85.9%)(Dev)(Prod)
### Fixed
- Fixed bug in proxy whitelist where it was including applications that had been deleted.
- Fixed tests that reduced coverage due to some functions only sometimes running (sort etc.)

## [3.56.0] - 2025-06-23 (85.8%)(Dev)(Prod)
### Added
- Added support for parsing new elements from BBC articles, including; publish time/date, author, inline image/videos, and captions.
### Changed
- Refactored helper method inside `ProcessNode` into its own top-level function (`AddBlockIfNotEmpty`) in the RSSProcessor.
### Fixed
- Fixed processing of BBC sport and "video-only" articles.
- Fixed the parsing of another type of BBC article that only show a video and a small amount of text, including new tests that work with a real HTML file.

## [3.55.2] - 2025-06-23 (85.7%)(Dev)(Prod)
### Added
- Added `com.sec.android.soagent` and `com.wssyncmldm` to policy if OTA upgrade is allowed.

## [3.55.1] - 2025-06-19 (85.7%)(Dev)(Prod)
### Fixed
- Fixed file upload limit for THe Hub generate from PDF feature. It was set to 4.7MB and not 4.5MB due to the differences between MiB and MB

## [3.55.0] - 2025-06-19 (85.7%)(Dev)(Prod)
### Added
- Added an option to hide the system roles on an enterprise
- Add a `hidden` flag onto a role so we cna manually hide a role from the DB

## [3.54.0] - 2025-06-18 (85.7%)(Dev)(Prod)
### Changed
- Refactored RSS news article content processing, enabling granular control over text styling and formatting within frontend.

## [3.53.0] - 2025-06-12 (85.7%)(Dev)(Prod)
### Added
- Added markdown summary type format support
- Added an option on the group to disable the sd card on the device
### Changed
- Video file max size for the hub

## [3.52.2] - 2025-06-09 (85.7%)(Dev)(Prod)
### Deployment
- Make sure to clear 'Successful Oauth2 Token Created' activities from the database.
### Added
- Added more indexes to `activites` and `owner_activites` to help with queries found in slow queries log in prod.
### Changed
- Removed the Activity log for when oauth tokens are created, and moved to AWS CloudWatch log instead.
- Updated `target` in `owner_activites` table from 1000 char to 60 char so it can be indexed. It should only have contained the app package name anyway.

## [3.52.1] - 2025-06-05 (85.7%)(Dev)(Prod)
### Added
- Added support for Claud Sonnet & Opus V4

## [3.52.0] - 2025-06-04 (85.7%)(Dev)(Prod)
### Added
- Added index to activities table for enterprise and user ID.
- Added scribe group information to scribe job responses
- Added endpoint to fetch all the scribe group memberships for a given user
### Changed
- Preload on the hub models to activities
- Updated messages from tags handler to use Topics (Which is what it is called on the frontend)
- Modified file service read file info function to return a struct with filesize on it
- Modified article and tag file upload endpoints to have a max file size
- Article delete endpoint now deletes all comments related to the article
### Fixed
- Fixed the webhook event service so that it can create events for scribe jobs that have been deleted

## [3.51.0] - 2025-06-03 (85.8%)(Dev)(Prod)
### Added
- Added endpoint to list activities for the current user
### Changed
- Updated `StandardClaims` in all JWT tokens to `RegisteredClaims` as `StandardClaims` was depreciated
- Updated name for  `rss.RSSFeedService` to `rss.FeedService` for Qodana

## [3.50.0] - 2025-06-02 (85.8%)(Dev)(Prod)
### Added
- Added image URL to news feed content response
- Added model ID filter to activities list endpoint
- Added language code to device owners' accessibility settings
### Changed
- Updated the activity messages for all TheHub resources


## [3.49.0] - 2025-05-23 (85.8%)(Dev)(Prod)
### Added
- Added support for AWS nova models in summary service
### Changed
- Refactored RSS service
### Fixed
- Fixed dynamic prompt summary model version
- Fixed support for special characters in scribe transcriptions and summaries

## [3.48.1] - 2025-05-13 (85.6%)(Dev)(Prod)
### Changed
- Changed collation of the owner_activities and device_activities as sometimes the apps and data we use will have unicode characters in them.

## [3.48.0] - 2025-05-13 (85.6%)(Dev)(Prod)
### Added
- Added an option on the group to disable the USB access or restrict it to specific USB types

## [3.47.1] - 2025-05-13 (85.6%)(Dev)(Prod)
### Fixed
- Fixed issues generating summaries of custom types for scribe jobs

## [3.47.0] - 2025-05-13 (85.6%)(Dev)(Prod)
### Added
- Added dynamic scribe prompts for scribe summary types
- Added admin management endpoints for scribe prompts
- Added scribe summary worker factory to handle both static and dynamic summary types

## [3.46.0] - 2025-05-07 (85.3%)(Dev)(Prod)
### Added
- Added ability to search for tags depending on featured status
- Support for fetching article images and content within the `RssFeedService`.
- New endpoint in the apps level `NewsFeedHandler` for fetching article content.
### Changed
- Reworked RSS service tests

## [3.45.0] - 2025-05-02 (85.4%)(Dev)(Prod)
### Added
- Added disabling of the google voice typing keyboard when the default keyboard is disabled
### Fixed
- Removed 'Automatic' EAP Inner from the tests and the network request example as it's no longer used

## [3.44.0] - 2025-05-02 (85.4%)(Dev)(Prod)
### Added
- Added admin endpoint to list all users filtered by enterprise
- Added article copy service to allow admins to copy articles from one enterprise to another
### Changed
- Removed ability to search articles by their title and body. It is now just by the title.
- Removed the sanitise service check from the update endpoint for TheHub article comments

## [3.43.0] - 2025-04-28 (85.4%)(Dev)(Prod)
### Fixed
- Fixed the `TestListTaskMessages` tests as it was done on created_at, when ti should be done on ID for consistent ordering. Tests have always been wrong but as they happen in the same second it would only occasionally show.
### Security
- Replaced `github.com/golang-jwt/jwt` version `v3.2.2+incompatible` with `github.com/golang-jwt/jwt/v4 v4.5.2` due to https://github.com/golang-jwt/jwt/security/advisories/GHSA-mh63-6h87-95cp
- Updated `golang.org/x/oauth2` from version `v0.21.0` to `v0.29.0` due to https://www.mend.io/vulnerability-database/CVE-2025-22868?utm_source=JetBrains
- Updated `golang.org/x/net` from `v0.37.0` to `v0.38.0` due to `golang.org/x/net vulnerable to Cross-site Scripting`
- Updated `google.golang.org/api` version `v0.186.0` to `v0.230.0` as it was targeting an old version of protobuf.
### Removed
- Removed YouTube video sync log messages as there are 4k messages each night just saying each channel has been synchronised. We just need to log errors or useful information.

## [3.42.0] - 2025-04-23 (85.4%)(Dev)(Prod)
### Added
- Added endpoints for admins to manage additional headers for webhooks
### Changed
- Updated the webhook service to fetch and set all additional headers before sending the request
### Removed
- Removed YouTube video sync log messages as there are 4k messages each night just saying each channel has been synchronised. We just need to log errors or useful information.

## [3.41.0] - 2025-04-17 (85.3%)(Dev)(Prod)
### Added
- Added arabic support to the translation and text to speech services
- Added MP3 support for TheHub article video files
- Added flag for allowing TheHub profile page to enterprises
### Changed
- Updated TheHub article list status filter to include a scheduled option

## [3.40.0] - 2025-04-04 (85.3%)(Dev)(Prod)
### Added
- Added an endpoint to retrieve a specific network
- Added name of the proxy exclude list onto the network response
- Added checks to ONC to amke sure the proxy exclude lists have been loaded onto the network if they are supposed to exist. Policies won't be changed if they are not. THis will ensure we don't tank a device in a secure environment.
### Fixed
- Added network exclude lists into groups when updating policies from device updates.
- Fixed a couple of tests that were failing due to using BST in mysql

## [3.39.2] - 2025-03-27 (85.3%)(Dev)(Prod)
### Changed
- Changed collation of articles, article translations and article comments

## [3.39.1] - 2025-03-27 (85.3%)(Dev)(Prod)
### Added
- Excluded `/ip` endpoint form the logger as it is used as a ping endpoint by some clients

## [3.39.0] - 2025-03-25 (85.3%)(Dev)(Prod)
### Deployment
- Make sure to add new environment variables `PKI_HOST`, `PKI_API_KEY` and `MOCK_PKI_API`
- Check the current proxy exclude lists and global proxy exclude lists and make sure that the domains within them are split by new line
### Added
- Optional origin field to scribe jobs
- Added global proxy exclude list endpoints
- Added ability to attach global proxy exclude lists to proxy exclude list
- Added `policyUpdatedAt` to the managed config and amapi policy updates table
- Added error responses to the enterprise whitelist endpoint if the data fetch from any of the repos fails
- Added support for provisioned certificates on the device
- Added an option to manually set the client certificate on the network
- Added an option to disable the network
- Added the PKI agent mock so that it can be used when running the API locally
- Added groups for Purple Scribe resources and service for determining access based on group
### Changed
- Order Calendar Events and Kiosk Announcements returned by their start dates.
- Changed PKI mock into http transport mock to test the PKI service
- Changed the domains within proxy exclude lists to be split by new line characters instead of commas
### Fixed
- Updated scribe job and summary types so that name uniqueness is checked on create/update
### Removed
- Removed the configuration of network settings from the Create endpoint
### Security
- Updated `golang.org/x/net` from `v0.33.0` to `v0.36.0` for HTTP Proxy bypass using IPv6 Zone IDs in golang.org/x/net
- Updated `github.com/corazawaf/coraza/v3` from  `v3.2.1` to `v3.3.3` for OWASP Coraza WAF has parser confusion which leads to wrong URI in `REQUEST_FILENAME`

## [3.38.0] - 2025-03-18 (85.5%)(Dev)(Prod)
### Changed
- Updated the length of TheHUB article titles from 128 characters to 1000 characters
- Allowed TheHub articles to have PDFs that exceed 4.5MB but don't allow summaries on them
### Added
- Added sorting to TheHub endpoint for getting articles by tag
- Added search filter for EAPI policies list endpoint
- Added Scribe Job Types tables
- Added Scribe Summary Types tables
- Added admin endpoints to list, get, create, update and delete scribe job types.
- Added admin endpoints to list, get, create, update and delete scribe summary types.
- Added list endpoint so system users can read all scribe job types that their enterprise has access to.
- Added list endpoint so system users can read all scribe summary types that their enterprise has access to.
- Added seeders for scribe job types, summary types and the relationships between them

## [3.27.2] - 2025-03-11 (85.4%)(Dev)(Prod)
### Added
- Additional USA -> UK word replacements for purple scribe
- Added error logging to the activity DB inserts, as we received an error message that the `message` was too long but can't trace where it came from.
### Changed
- Modified Serco Case Management Prompt

## [3.37.1] - 2025-02-27 (85.5%)(Dev)(Prod)
### Added
- Added error message in PDF is too long to summarise
### Changed
- Modified model used by PDF summary prompt to claud 3.5
- Modified PDF summary prompt to ignore instructions provided in the PDF

## [3.37.0] - 2025-02-26 (85.5%)(Dev)(Prod)
### Added
- Added summarise PDF method to summary service with a new PDF specific prompt
- Added an endpoint to get a summary of an article's PDF using the new summary service method
- Added WYSIWYG method to the sanitise service to retain tags permitted by the WYSIWYG

## [3.36.1] - 2025-02-25 (85.5%)(Dev)(Prod)
### Added
- Added in the youtube channel name to the issues email that gets sent.
- Added in `https://dev.app.purple-scribe.com` to CORs rules.
- Added ability to sort TheHub tags by name

## [3.36.0] - 2025-02-24 (85.5%)(Dev)(Prod)
### Added
- Added a check in the enterprise proxy update endpoint if any of the enabled enterprise proxy ip addresses already exists on another enabled enterprise proxy
- Added indexes for heavily used tables.
### Changed
- Changed sequential ID to UUIDv7 in the calendar events table
- Changed model id from int to string in the activities table to accommodate uuid entries
- Updated `Authentication` action for owner activities table to `AuthenticationSuccess` and `AuthenticationFailed` so we cna easily group it for reporting.

## [3.35.1] - 2025-02-19 (85.5%)(Dev)(Prod)
### Fixed
- Fixed issue with new summary types not being added to the list of valid summary types

## [3.35.0] - 2025-02-19 (85.5%)(Dev)(Prod)
### Added
- Allowing logs to be made for deleted devices
- Device Uid to the logger
- Storing device policy version in the device data
- New device activity when the device policy version is updated
- Added looking up saved youtube videos before creating new ones during the sync to avoid duplicate entry errors
- Added sending notification email if any youtube channels fail to be fetched during the sync
- Added new job and summary types for scribe: ACCTReview, ASCFParoleInterview, ASCFInitialInterview, ASCFParoleAssessmentReport
### Fixed
- Fixed an issue with Article comments as it's ordering by date desc, rather than the ID.
- Fixed issue with TestTheHubResidentArticleList tests as sporadically failing due to the ordering of the articles

## [3.34.1] - 2025-02-14 (85.5%)(Dev)(Prod)
### Deployment
- Make sure to set the new `COUNTRY` environment variable
### Fixed
- Fixed invite and password reset URLs not working for Australia
- Check a device exists before trying to resync it's policy in the delayed calls.

## [3.34.0] - 2025-02-14 (85.5%)(Dev)(Prod)
### Added
- Added PDF field to articles
- Added endpoints to upload and delete article PDF
- Added endpoint for owner to get files cookie
### Changed
- Delete article PDF when article is deleted

## [3.33.0] - 2025-02-07 (85.5%)(Dev)(Prod)
### Added
- Added S3 service interface
- Added S3 service to Deps
- Added Text To Speech migrations and models
- Added `TTS_BUCKET` ENV var
- Added Text To Speech Service, Tests & Integration Tests
- Added Polly Service + mock
- Installed AWS Polly SDK
- Added TTS handlers and tests
- Added `job_title` and `pronouns` fields to Users table + ability to update
- Added a Profile Picture for users + endpoints to update and delete
- Added endpoint to support uploading images from the article WYSIWYG
### Changed
- Added Additional methods to S3 service (Delete, DeleteFromURI, Get, GetFromURI)
- Modified Transcription service to use S3 Interface
- Article Update endpoint now deletes old TTS entries
- Article Delete endpoint now deletes tag associations, updates tag article count, and deletes TTS entries
- Grated Device owners access to profile pictures of system users (Used for author images)

## [3.32.1] - 2025-02-06 (85.5%)(Dev)(Prod)
### Added
- Added in some Indexes for the youtube API
### Changed
- Removed duplicate channel ID warnings from the errors list as these will happen periodically anyway.

## [3.32.0] - 2025-02-06 (85.5%)(Dev)(Prod)
### Added
- Added in the Feed ID and URL to any error logs for hte RSS feeds so we know which one is causing issues.
- Added a AMAPI policy update tracking table, and updated all places where policyies are updated so we can track the process origins for the update. Also track the update itself in json. NOTE: we may need to truncate this table periodically.
### Fixed
- Added migrations to fix issue with translations by setting article & tags collation to `utf8_general_ci`
- Fixed a bug where the owner wasn't being passed though when resyncing a policy after allocation.
### Security
- Updated package from `github.com/microcosm-cc/bluemonday` from `v1.0.25` to `v1.0.27` as recommended by the package owner.

## [3.31.2] - 2025-01-31 (85.6%)(Dev)(Prod)
### Fixed
- Fixed issue loading tag images for resident

## [3.31.1] - 2025-01-31 (85.6%)(Dev)(Prod)
### Fixed
- Fixed issue getting article thumbnail when loading by tag id

## [3.31.0] - 2025-01-31 (85.6%)(Dev)(Prod)
### Added
- Added tables for articles, tag and translations copied over from TheHub repository
- Added endpoints to list all articles and get or update specific articles information
- Added Owner Accessibility table and endpoints
- Added Translation service for Articles and Tags
- Added the hub urls to CORs rules
### Changed
- Granted owners read access to the hub folder in the files cookie
- Changed files cookie to include domain

## [3.30.0] - 2025-01-24 (85.6%)(Dev)(Prod)
### Added
- Device policy resync 2 minutes after the device has been attached to an owner as a preventative measure in case android doesn't start the application installation straight away.
### Fixed
- Fixed failing pdf forms Read and ReadAll tests

## [3.29.4] - 2025-01-15 (85.6%)(Dev)(Prod)
### Changed
- Updated `ip_address` fields to 200 characters in `failed_key_validations`, `failed_logins`, `activities` tables
### Fixed
- Fixed race condition in `pdf-forms` handler where the default ordering was based on created_at desc. Updated to id desc.

## [3.29.3] - 2025-01-15 (85.6%)(Dev)(Prod)
### Added
- Added indexes to the messages table.
### Fixed
- Fixed the calendar event update endpoint so that it's updating the event by ID rather than model uid
- Added in the device UID when deleting other devices in teh devices table. This should ensure that if two enrollment messages come in at exactly the same time, they won't end up deleting each other.

## [3.29.2] - 2025-01-14 (85.6%)(Dev)(Prod)
### Added
- The enrollment messages happened to have unique ID's so it's confirmed bug from android now. Added some code to only process the first enrollment message in the queue_service.

## [3.29.1] - 2025-01-13 (85.6%)(Dev)(Prod)
### Added
- Added in a `msg_id` to the pub/sub messages from Android Management API as we've recently received the same message twice, causing issues with the enrollment.

## [3.29.0] - 2025-01-07 (85.6%)(Dev)(Prod)
### Deployment
- Make sure the environment variables `COOKIE_SECRET` and `COOKIE_CRYPT` are set for mdm-api deployment
### Added
- Login error handling in the SSO test client
- Added "App SSO endpoint" so that certain apps can auto login the user to the SSO server
- Changed owner UIDs for bt from OwnerUid to EnterpriseUID-OwnerUid
- Automatically uppercase the prisoner ID if the BT integration is turned on as they can't handle lower case Nomis ID's. Their system throws a 500 error.
- Added Device application list and admin device application list endpoints
### Changed
- Updated docker build to go version `v1.23`
### Security
- Updated `golang.org/x/net` from `v0.26.0` to `v0.33.0` for `Non-linear parsing of case-insensitive content in golang.org/x/net/html`
### HOTFIX
- Added code for owner tokens which would have the wrong owner ID for BT integrations after we have updated the DB to enable unique identifiers across the BT enterprises.


## [3.28.1] - 2024-12-19 (85.5%)(Dev)(Prod)
### Fixed
- Updated the expiry for the SSO ID token to 10 minutes

## [3.28.0] - 2024-12-13 (85.5%)(Dev)(Prod)
### Added
- Added additional banned words for scribe summaries
- Added device activity and device owner activity endpoints
- Added an endpoint where external API users can see what IP address they are coming from.
### Changed
- Updated default policy for samsung keyboard application
### Security
- Updated `golang.org/x/crypto` from `v0.24.0` to  `v0.31.0`

## [3.27.0] - 2024-12-06 (85.4%)(Dev)(Prod)
### Added
- Added keyboard_installed flag on the device so that we know when we are free to disable the samsung keyboard
- Added a check in the queue service and device activity handler to verify if the simple keyboard has been installed, setting keyboard_installed to true on the device and triggering a device policy sync
- Added device_application_lists table to store the list of installed applications on the device from the status report
### Changed
- Enabled updating DefaultKeyboardEnabled setting on the group

## [3.26.1] - 2024-12-06 (85.4%)(Dev)(Prod)
### Added
- Added additional banned words to help with new prompts

## [3.26.0] - 2024-12-06 (85.4%)(Dev)(Prod)
### Added
- Added "webm" file extension and mime type to be allowed when uploading data for scribe jobs.
- Added presentation review prompt and job type
- Added prisoner phone call prompt and job type
- Added business call prompt and job type
### Changed
- Reduced max concurrent calls for the webhook service
- Removed experiments folder from coverage report
- Removed purple-scribe.com and www.purple-scribe.com from CORS. Added app.purple-scribe.com
- Updated prompt for Purple Visits summary

## [3.25.1] - 2024-11-27 (85.4%)(Dev)(Prod)
### Added
- Added deallocation policy sync for device after 60 seconds of the device being deallocated.
- Added cron job to scheduler and hander for checking for devices stuck in deallocation, and synchronising every 5 minutes up to a cap of 12 attempts.
- Added `SercoCaseManagement` Prompt and Summary
- Added new transformer for Assurance Alliance Action Plan
- Added `scribe_format` to Summaries
### Changed
- Moved cron handler into cron folder
- Added in a fromDevice flag on the device_activity table so we can start to add logs against the device for processed in the server.
- Updated Q&A Prompt
- Updated Assurance Alliance Action Plan Prompt
- Renamed current `scribe_type` to `scribe_format`

## [3.25.0] - 2024-11-26 (85.5%)(Dev)(Prod)
### Added
- Scribe: Additional USA -> UK spellings
- Scribe: Tests for scribe stats endpoints
- Added `Scribe.Stats` permissions
- Added device activities table and device activity handler
- Added owner uid to the device mdm token and device uid to the owner mdm token
- Added `mdmDeviceToken` to the device's managed config
### Changed
- Modified `date_range` to support multiple common date formats
- Modified Scribe stats endpoints to use new stats permissions
- Changed the path to the enrollment handler (still supporting the previous one)

## [3.24.6] - 2024-11-20 (85.4%)(Dev)(Prod)
### Fixed
- Fixed issue with purple scribe stats endpoints

## [3.24.5] - 2024-11-20 (85.4%)(Dev)(Prod)
### Changed
- Increased max characters for transcript and summary update request

## [3.24.4] - 2024-11-20 (85.4%)(Dev)(Prod)
### Changed
- Changed CPS Meeting and CPS Meeting PDF to use Sonnet 3.5

## [3.24.3] - 2024-11-20 (85.4%)(Dev)(Prod)
### Added
- Added route to update scribe transcriptions
- Added MeetingActionsAndTakeaways, CPSMeeting and CPSMeetingPDF summaries & job types
### Changed
- Ability to update scribe summaries
- Modified `url.go` to allow `/scribe-transcriptions` route through.

## [3.24.2] - 2024-11-18 (85.4%)(Dev)(Prod)
### Added
- Added secure and http flags to the SSO server cookie
- Added `FULLWINDOW` as an options for group updates. This allows a `WINDOW` which is the full day. This is mainly for forcing updates to be done at any time.
### Changed
- SSO redirect uris are now lowercased when they're compared
### Fixed
- Fixed issue with being able to generate summaries of the same type

## [3.24.1] - 2024-11-18 (85.4%)(Dev)(Prod)
### Added
- Added `purple-scribe.com` and `www.purple-scribe.com` to CORS
### Changed
- `GroupUID` on owner activties is no longer a point
- Set groupUID on owner activity creation in all places

## [3.24.0] - 2024-11-15 (85.4%)(Dev)(Prod)
### Added
- Added `group_uid` to owner activities
- Added group_uid and date filtering to owner activity list endpoint
- Added endpoint to retrieve the number of successful login activities sorted by date and filtered by group, owner and date
- Added endpoint to get top kiosk tags for owner activities, filtered by group, owner and date
- Added a date range validator to the `pkg` directory
### Changed
- Updated owner activity filtering to use `from` and `to` dates rather than a selected period

## [3.23.2] - 2024-11-14 (85.3%)(Dev)(Prod)
### Fixed
- Fixed issue creating webhook event when generating scribe summary on demand
- Fixed issue with summary mdm_summary_type and model being set to incorrect values

## [3.23.1] - 2024-11-13 (85.3%)(Dev)(Prod)
### Fixed
- Fixed fetching transcription when generating summary in external handler

## [3.23.0] - 2024-11-13 (85.3%)(Dev)(Prod)
### Added
- Added summary service for scribe jobs using AWS Bedrock
- Added endpoint to allow users to generate a summary for a scribe session on demand
### Changed
- Updated external endpoint to generate summaries for scribe jobs
- Updated DST prompt to a new version
- Updated Keyworker prompt to a new version

## [3.22.0] - 2024-10-22 (85.3%)(Dev)(Prod)
### Added
- Added disabling factory applications
- Handling of additional Bad Request response (400) from the BT authentication end point
- Handling of Bad Request response from the BT authentication end point
- Added device OTA Update enable and disable endpoints
- Added device and enterprise uid columns in the messages table and setting them in the queue service
- Added security patch level and android version to the device table and included it in the device response

## [3.21.0] - 2024-10-22 (85.3%)(Dev)(Prod)
### Added
- Added read permissions for the Purple Visits, Purple Post & theHub endpoints
- Added `DeleteDir` to File Store
- Added `audio_expired`, `transcripts_expired` and `summaries_expired` flags to scribe jobs
- Added endpoint to download the stored audio file for a scribe job
- Added a closer.Close() function so you can defer the closure of an io.Closer and handle the error
### Changed
- Updated cron job to delete folder for scribe job instead of just the audio file itself
- Updated login user response to include enterprise allow scribe audio download
- Updated BT logs so they are JSON only without the log date being printed.
### Fixed
- Now calling `defer cancel` for `WithTimeout` functions when using Bedrock

## [3.20.0] - 2024-10-22 (85.4%)(Dev)(Prod)
### Added
- Data retention limits to enterprises for scribe audio, transcripts and summaries
- Expiry dates for stored scribe audio, transcripts and summaries
- Cron jobs to delete expired audio, transcripts and summaries at midnight every day
### Changed
- Enterprise updates now include retention limits (in days) for scribe audio files, transcripts and summaries

## [3.19.0] - 2024-10-21 (85.4%)(Dev)(Prod)
### Added
- SSO - Added a `home` page for `/` route
- Added styling for documentation page, and wrote up full docs. Added version into the page
- Added new endpoint to let the user resend the password on the device that's in the enrollment group
### Changed
- Changed `OwnerLogin` function in the bt pin service to return the failed login message alongside the error
### Fixed
- SSO - Fixed scroll bar issue on Error and logout pages
- Owner activities count, now excluding activities with `system.log` target from the count

## [3.18.0] - 2024-10-16 (85.4%)(Dev)(Prod)
### Deployment
- Add the SSO related server variables need adding to both the new SSO Server and to the API, so the api can generate cookies for use in the app.
- Need to manually update the existing service accounts for client credentials to have a user_oauth row.
### Added
- Added SSO server, docker files and deployment code. This enables the Authentication Code Flow on a dedicated service.
- Added validation in oauth code credential type, and ensured that the flow is enabled onthe client.
### Changed
- Moved the `server` folder for the API into an `api` folder as we now also have a second server called `sso` (and `scheduler` etc.).

## [3.17.0] - 2024-10-16 (85.5%)(Dev)(Prod)
### Added
- Added proxy excluded lists table containing the list of domains to be excluded from using the proxy on the device
- Checksum of the returned data in the proxy whitelist response
- Added logging public activity when the owner's group is changed during the 'assign a device' process

## [3.16.0] - 2024-10-11 (85.5%)(Dev)(Prod)
### Added
- Use the application track ID when creating a policy.
### Changed
- Filter for searching scribe jobs now includes the unique ID as well as the title

## [3.15.0] - 2024-10-02 (85.5%)(Dev)(Prod)
### Deployment
- Make sure to set the new `PROXY_CODE_KEY` environment variable
### Added
- Enabling proxy on the enterprise
- Passing proxy details in the device policy
- Setting the proxy details on the network
- Added `global_proxy_domains` table for storing globally whitelisted urls for the proxy
### Fixed
- Fixed the `.env.dist` file as `MOCK_GOOGLE_API` had been overwritten.

## [3.14.0] - 2024-09-19 (85.5%)
### Added
- Added `%device_id%` and `%owner_identifier%` to the managed configuration
- Added validation of the app package names
- Added an endpoint to the device assignment handler that registers the device with BT
### Changed
- The `IsDeviceSecure` flag now also reports on `DevicePosture`. IsDeviceSecure was originally just showing if a pin/password was set, but we also now amke sure it reports on the DevicePosture
- Changed the Status Bar Options in the group settings so that System Info Only option is no longer considered insecure
- Changed group UIDs for bt from GroupNumber to EnterpriseUID-GroupNumber

## [3.13.1] - 2024-09-05 (85.5%)
### Changed
- Allow m4a files with the old `audio/x-m4a` mimetype to be uploaded to scribe jobs. Note: Go detects this mimetype as `application/octet-stream`

## [3.13.0] - 2024-09-05 (85.6%)
### Added
- Added webhooks model and migration
- Added super admin webhook endpoints
- Added webhook events model and service
- Added webhook service and tests
### Changed
- Modified scribe endpoints to call the webhook events service

## [3.12.1] - 2024-09-04 (85.7%)
### Added
- Added the following headers X-Permitted-Cross-Domain-Policies, X-Frame-Options , X-Content-Type-Options

## [3.12.0] - 2024-09-03 (85.7%)
### Added
- Common password list for the password security service
- Added `123` as common pattern to block for allowed passwords
- Added max character validation to all requests
- Added the duration flag onto commands, set as a week
- Monthly brute force limit for the login endpoint
### Fixed
- MIME Type checks in the scribe job handler

## [3.11.0] - 2024-08-20 (85.6%)
### Added
- New param to possession category to allow excluding from location max item checks.
- Added `LoginRateLimit` to config.
### Fixed
- Check deleted_at is null when listing items for volumetric controls
- Fixed issue with rate limit applied to login routes when running tests. Made it configurable for testing.

## [3.10.1] - 2024-08-02 (85.6%, 100%)
### Changed
- Updated the owner uid when bt integration is enabled so that it's an uppercased identifier instead of being random
- Changed group UIDs for bt from EnterpriseUID-IEP-GroupNumber to GroupNumber

## [3.10.0] - 2024-08-01 (85.6%, 100%)
### Added
- Added various new job types for scribe
### Changed
- Ability to upload MP4 files for scribe jobs

## [3.9.1] - 2024-07-27 (85.6%, 100%)
### Changed
- Updated rate limiter to 10/s from 20/s

## [3.9.0] - 2024-07-27 (85.6%, 100%)
### Added
- Made the length of the password for Owners configurable per enterprise. Minimum length is 6.
- Added an option on the enterprise for "Simple Password" protection for Users
- Added an option on the enterprise for "Simple Password" protection for Owners too, same as above
- Added in a basic breach list where the top 10k known passwords are checked against.
- Rate limit authentication endpoint.
- Ensure the 2FA OTP page has brute force protection built in. We should be adding something to the brute force protection check each time.
- Update the users password to a default of 12. Updated all tests to make sure they are all updated for the new password requirement.
### Changed
- Updated logging for the BT service to just provide JSON rather than it been mixed

## [3.8.2] - 2024-07-23 (85.6%, 100%)
### Fixed
- Fixed issue causing server to crash when updating a devices owners group. Checks for if the wait group were nil were missing so have been added in
### Changed
- Locations now come back in alphabetical order
- Categories now come back in alphabetical order

## [3.8.1] - 2024-07-23 (85.6%, 100%)
### Changed
- Removed required validation for `max_items` when updating a possession category. This is so max items can be set to 0.

## [3.8.0] - 2024-07-19 (85.6%, 100%)
### Added
- Added `BEDROCK_REGION` to ENV file
- Added owner possession category endpoints
- Added `is_prohibited` filter to owner possessions LIST endpoint
### Changed
- Users file token now uses the `RefreshTokenDuration` specified on their enterprise
- Ability to create PDF forms without specifying an owner identifier

## [3.7.3] - 2024-07-18 (85.5%, 100%)
### Added
- Added BT Login to the `Attach` handler when the BT integration is enabled, same with BT Provisioning in the queue service. Pure API calls to BT have integration test but not in this test coverage so dropped percentage

## [3.7.2] - 2024-07-18 (85.6%, 100%)
### Added
- Added an extra scope in Activities List to filter out 'Successful token created' activities, also added service_account to the activity list response
- Granular update control per device. Added `allowOTAUpgrade` to groups and to devices. Groups will be set first on a policy, and if the device exists it will be overridden by whatever is set on the device.
### Changed
- Updated new groups so that network escape hatch is disabled by default

## [3.7.1] - 2024-07-16 (85.5%, 100%)
### Added
- Added status change on device owner login
- Added `status` and `previous_status` to device data endpoint.

## [3.7.0] - 2024-07-16 (85.5%, 100%)
### Deployment
- Make sure to add new environment variables `BT_HOST` and `BT_API_KEY`
- Make sure once the migrations have run, the default status for devices is `active` or `enrolled`
### Added
- Added Status to device with all stages: "provisioning","enrollment","active","disabled","locked","deallocate"
- Added deallocated, and job to ensure applications are removed before going back into the enrollment group
### Changed
- Changed qodana code quality GitHub action from `v2023.2` to `v2024.1`
- Refactored queue service setup so the B64 credentials are only fetched when called specifically in production.
### Fixed
- When a device is deleted it sends though status report without a policy name on them. Fixed so it wasn't trying to change the policy on the device.
- Removed HOST form files cookie, as it will autogenerate it.

## [3.6.2] - 2024-07-16 (85.1%, 100%)
### Added
- Added `requires_possession_owner_id` filter to owners location handler
### Changed
- Updated location seeder to require possessions in stores to have an owner
### Fixed
- Fixed issue getting images for possessions from LIST endpoint

## [3.6.1] - 2024-07-12 (85.0%, 100%)
### Fixed
- Used correct ENV for file token url

## [3.6.0] - 2024-07-12 (85.0%, 100%)
### Added
- Updated the owner's password validation so that the owner's new/updated password is checked against the enterprise's password requirements
- Added Device Owner Login Enabled field to the enterprise settings
- Added marking groups as secure
### Fixed
- Added files token for users and owners. Added services to support this
- Endpoint to get files with the files token
- Added `max_items` to locations
- Added `possession_request_message` to enterprise
### Changed
- Moved files_service into pkg folder to fix issue with circular deps with new files token
- Modified users refresh endpoint to give files token
- Modified owners login endpoint to give files token
- Modified file routes to use middleware auth with files token. Left public folder as is to aid with backwards compatibility
- Modified volumetric controls to respect the `max_items` on a location
- Modified possessions requests to respect a locations `max_items` if volumetric controls and enabled
- Modified volumetric control response to include location info
- Removed Room from location seeder. Changed stores to be ID 1
- Ability to update `possession_request_message` on an enterprise
- Return `possession_request_message` to an owner when possession request is created
- Queue service should fetch the deleted devices too when making updates, in case the device is "stuck" and still sending reports even though it's been asked to be WIPED.
- Issue with being able to exceed volumetric controls on possession requests
### Security
- Updated `google.golang.org/grpc` from `v1.64.0` to `v1.64.1` due to `GHSA-xr7q-jx4m-x55m`

## [3.5.1] - 2024-07-08 (85.0%, 100%)
### Added
- Added in swagger validation step to local development Dockerfile. Note: require rebuild using `docker-compose build mdm_app`
- Added additional check to see if the owner has reached the limit of allowed devices when assigning a new device
- App enrollment endpoint
- Added search and sorting filters to machines list method.
### Changed
- Moved the `queue` service so that it is now included in the full test suite. This dropped the overall coverage by 0.2% as there are functions in there specifically for the Google Pub Sub Service that aren't mocked.
- Removed the `Permissions` from the enterprise token for the iframe for Google play

## [3.5.0] - 2024-07-03 (85.1%, 100%)
### Added
- Added Possessions, possession categories and possession requests
- Added Locations

## [3.4.1] - 2024-06-28 (84.7%, 100%)
### Changed
- Updated ScreenCaptureAllowed to be controlled by the group
- Added `screenCaptureDisabled` to managed config for the agent application
- Added applied group policy version info to devices

## [3.4.0] - 2024-06-28 (84.7%, 100%)
### Added
- Added TimeToLock to groups
### Changed
-`google.golang.org/api` updated from `v0.128.0` to `v0.186.0`

## [3.3.0] - 2024-06-27 (84.6%, 100%)
### Added
- Added speaker detect related models and migrations for purple scribe.
- Added ScribeSpeakerSegments and Speakers to ScribeJob response and preloaded them on
- Added endpoint to update Scribe Transcription Speakers
- Added external endpoint to create speaker aware transcripts
### Fixed
- In Scribe Jobs Get endpoint fixed the read all permission that was being checks. Changed from PDF forms to Scribe Job

## [3.2.0] - 2024-06-25 (84.7%)
### Added
- Added the MAC address randomization mode selection
- Added swagger validation library
- Added auto-generation of swagger docs back in
- Added disabling/enabling devices functionality
### Changed
- Updated Go to 1.21 from 1.19. Updated docker images so local testing image will need rebuilding for each dev.
### Fixed
- Fixed swagger documentation in handler comments for various endpoints.

## [3.1.0] - 2024-06-20 (84.7%)
### Added
- Added oauth2 token endpoint for authentication.
### Fixed
- Fixed all the imports for repositories2 and models2 to remove the incorrect aliasing

## [3.0.1] - 2024-06-20 (84.7%)
### Changed
- Default enrollment group will now include the Status update messages enabled.
- Start the Doc generation for development on docker-compose start again

## [3.0.0] - 2024-06-17 (84.7%)
### Deployment
- Make sure to set the `ENROLLMENT_TOTP_SECRET` environment variable
### Added
- Added in `track` to enterprise applications.
- Added search feature to enterprise list
- Added Owner Device Limit to enterprises
- Added Enrollment Group endpoints
- Enrollment Group: Added migration code to switch devices from being controlled by an Owner policy to its own device level policy
- Added group_uid filter to devices page. Special override for asking for `enrollment` which links to the enrollment group.
- Added endpoints to attach and detach devices from owners. Detaching a device will put it back into the enrollment group.
- Added an Admin level Device page
- Added a device TOTP endpoint for devices without an owner
### Changed
- Removed Nexus env variables
- Moved code for handling deletion of owner favourites in Group handler.
- `owner_policy_version` removed from the owner response
- Enrollment Group: System now works at device level policies, not owner level policies. Removed `UpdateOwnerPolicyFromGroup` and replaced with `UpdateDevicePolicyFromGroup`
- Enrollment Group: Updating a group now updates every device attached to that group using their device policy.
- Enrollment Group: Updating a group on an owner no updates every device owned by the owner using each device policy.
- `ExpectedCallBack` in the testing framework is now called before all other expected validations. this allows us to dynamically insert code to wait for wait groups, before validating other items in the test.
### Removed
- Owner polices are no longer used in the system. Where present on a device, they will automatically be switched whenever a group is updated or when an owner changes group.
- Code for deleting owner favourites when an owner changes group or when a group's app is disabled.

## [2.43.1] - 2024-05-24 (83.6%)
### Changed
- ModelId param to ModelUid on CalendarEvent model.
- Create/Update CalendarEvent requests to use time.Time objects instead of strings.
- Replaced scopes in CalendarEvent List handler for new repo method.

## [2.43.0] - 2024-05-22 (83.6%)
### Added
- Functionality for Calendar Events.

## [2.42.3] - 2024-05-13 (82.9%)
### Added
- Added track ID capability to the enterprise.
- Added `firmware_auto_update_allowed` option in managed configuration for Knox.

## [2.42.2] - 2024-05-10 (82.8%)
### Added
- Added in `track` to enterprise applications.
- Updated google play token to pull URL from environment.
### Changed
- Removed Nexus env variables
- Bumped all docker images to go v1.19
### Security
- Bumped `golang.org/x/net` from `v0.21.0` to `v0.23.0` for dependabot issue `close connections when receiving too many headers`

## [2.42.1] - 2024-04-18 (82.7%)
### Added
- Add check to see if a device has already been isolated, and if it has just update it in the database so that it's synchronised.

## [2.42.0] - 2024-04-17 (82.7%)
### Added
- FUnction added to android api to clone policies and isolate devices
- Endpoint added to isolate4 a device for superadmins. Override policy is returned on device is overridden.

## [2.41.1] - 2024-03-25 (82.4%)
### Fixed
- Fixed issue with getting channels that were assigned to deleted groups
### Changed
- Owner youtube channels now come back in alphabetical order

## [2.41.0] - 2024-03-20 (82.3%)
### Added
- Added `start_date` and `end_date` to kiosk announcements table
- Added ability to set `start_date` and `end_date` when creating/updating a kiosk announcement
- Added filter to list endpoint for purple play channels to get if it has been assigned to a group or not
- Favourite type to Application Favourite model.
### Changed
- Changed apps kiosk announcement handler to only return announcements within their start and end date window
- Renamed YoutubeChannel permission to PurplePlay
- Curating videos for a group now requires the PurplePlay Update permission.
- Purple Play channels now come back in alphabetical order
- Updated owner creation & update endpoint to remove whitespace from identifier

## [2.40.0] - 2024-03-12 (83.1%)
### Added
- Added `AWS_REGION` ENV variable
### Changed
- Changed Transcription Service, Transcription Queue Service & S3 Service to make use of the `AWS_REGION` Env variable.
- Allowed MP3 files to be uploaded for scribe v2

## [2.39.0] - 2024-03-07 (83.1%)
### Added
- Passed the enterpriseUID to the managed configuration

## [2.38.0] - 2024-03-07 (83.1%)
### Added
- Added purple scribe job type field with filtering
- Added Scribe Billing Record model & type
- Added `is_super_admin` flag to activities model. This denotes if an activity was created by another super admin. These activities can only been seen by other super admins.
### Changed
- When completing scribe job, meta and billable secs are now sent from the worker
- Billable secs now returned when fetching scribe job
- Updated scribe job and pdf form handlers to support Delete and Delete All permissions.
- Update profile handler to use current enterprise instead of user's enterprise
- Updated enterprise handler to clear all access records when access type is changed
- Created events when a scribe job is Updated, Viewed, Deleted or has a file uploaded
- Modified the activities endpoints to return only non-super admin activities to non-super admins.
- When creating activities set the `is_super_admin` flag depending upon the users super admin status.

## [2.37.0] - 2024-02-10 (82.4%)
### Added
- Added summary type column to scribe summaries table

## [2.36.0] - 2024-02-02 (82.4%)
### Added
- New endpoint and handler method for getting a group's news feeds.
- New tests for group news feeds Get method.
### Fixed
- List method in NewsFeedHandler.
- Admin update news feed method.
### Changed
- NewsFeedHandler List tests to match new fix.

## [2.35.1] - 2024-01-31 (83%)
### Fixed
- Added in staging app (`com.madepurple.agent.staging`) so logs update app name automatically for activity log.

## [2.35.0] - 2024-01-31 (82.4%)
### Added
- Functionality for RSS News Feeds.
- Added functionality for Owner Favourites.
- Functionality for Featured Apps on Group level.

## [2.34.2] - 2024-01-26 (82.1%)
### Fixed
- Fixed issue in owner and user middleware where the last_seen date was being updated for everyone

## [2.34.0] - 2024-01-24 (82.1%)
### Fixed
- Replaced package https://github.com/satori/go.uuid as there were security concerns with it.
### Changed
- Added migration to change type of scribe transcription from text to longtext.
- Update middleware to set last seen date on user and owner, instead of being set on activity creation.
### Added
- Added list all endpoint to pdf forms and scribe job handlers with permissions
- Added delete endpoints for scribe jobs and pdf forms

## [2.33.1] - 2023-11-15 (82.8%)
### Added
- Ability to order device owners by `identifier', `additional_info` and `last_seen`.
### Changed
- Updated endpoint that is called for hte application filters as it was using an admin only one.
- Changed application results in the api toi a slim version when underprivileged users are calling it.
### Fixed
- Replaced filter for `Opened application` on the top activities.
- Updated doc comments so new docs can build for Kiosk Announcements.

## [2.33.0] - 2023-11-14 (82.1%)
### Added
- Added `search` query param to users LIST endpoint
### Changed
- Modernised users handler. Switched to using server repos, renamed model2 and minor tidy up.
- Groups list endpoint now returns groups order by name ASC

## [2.32.0] - 2023-11-13 (82.1%)
### Added
- `last_seen` field to User model.
- Added `group_uid` query param to Owners LIST endpoint
### Changed
- Set `last_seen` on a user whenever a new activity is created for them.

## [2.31.0] - 2023-11-10 (82.1%)
### Added
- New module for Kiosk Announcements.

## [2.30.0] - 2023-10-18 (82.7%)
### Added
- Added kiosk global css table and endpoints
### Changed
- Logo and positioning has been updated.
- PurpleDocs mobile layout has been fixed.

## [2.29.2] - 2023-09-26 (82.0%)
### Added
- Added `docs.purplemdm.com` to CORS allowed origins

## [2.29.1] - 2023-09-26 (82.0%)
### Added
- Added in missing migrations for Purple Docs

## [2.29.0] - 2023-09-26 (82.0%)
### Added
- Tables for Purple Docs: Documents, Document Revisions, Templates and Template Categories.
- Admin endpoints for Purple Docs: Templates (CRUD & upload-image) and Template Categories (CRUD).
- Restricted endpoints for Purple Docs: Documents (list, get, delete, list-revisions).
- App endpoints for Purple Docs: Documents (CRUD) and Templates (list).
- BlueMonday library for HTML sanitization in Purple Docs WYSIWYG editor.

## [2.28.0] - 2023-09-22 (81.8%)
### Added
- Added forgot and reset password endpoints
- Added email when inviting user
- Added forgot password enabled flag to enterprise with custom email template
- Added IPv6 test cases for IP Package and ws_tests.
### Changed
- Create new user endpoint only takes name and email, not password

## [2.27.0] - 2023-08-29 (81.8%)
### Added
- Added option to add all YouTube channels to group via passing only 0 as the channel ID
- Added option to remove all YouTube channels from group via passing only 0 as the channel ID
- Added global css flag to kiosk sites
- Notes field to scribe job model
- Scribe job update handler with ability to update job notes and title
- Added `https://youtube.purplemdm.com` to CORS allowed origins
- Added last seen column to owners table
### Changed
- Changed update kiosk site endpoint to allow for global css flag to be set
- Changed prisoners get kiosk site endpoint to respect global css flag
- Changed owner response to include group name

## [2.26.1] - 2023-08-22  (82.5%)
### Fixed
- Add `apps/youtube` endpoint back in.

## [2.26.0] - 2023-08-03  (82.5%)
### Added
- Added ability to search by `additional_info` on owners list page
- Added `Repos` to server object as a quick way to access repositories anywhere that has access to the server.
- Added `youtube_video` table.
- Added GET endpoint for owner youtube channels
- Added App endpoints for youtube videos
- Added restricted endpoints for youtube videos
- Added YouTube service
- Added `YOUTUBE_API_KEY` env variable
- Added endpoint to sync videos for all youtube channels
- Added job to scheduler to call youtube video sync endpoint at midnight everyday
- Added `PurplePlayVideoAccess` to Enterprise model. This controls if the video access is a whitelist or blacklist.
- Added Youtube Video Access handler
### Changed
- Moved all Repos from the Dependency service to the `Repos` struct.
- Moved YouTube endpoints into their own handler
- Renamed `/youtube-channels` endpoint to `/youtube/channels`
- Changed GORM default connection charset from `utf8` to `utf8mb4` to support emojis in YouTube video titles/descriptions
- Changed route from `apps/youtube/channels/{id}` to `apps/youtube/channels/{channel_id}`. This instead uses the channels youtube id, not our id in the database.
- Youtube Video App handlers - Will not return any videos blocked/not included by the access list
- Minor modification to TestWalletTransactionsList. There was a chance the number `9` could appear in the owner_uid and cause the tests to fail. I've fixed by testing for `"amount":9`.
### Removed
- Nexus API Service

## [2.25.0] - 2023-08-03  (82.4%)
### Added
- Added `Repos` to server object as a quick way to access repositories anywhere that has access to the server.
- Added code to generate Open Network configuration properly from structs, and include basic EAP without certificates
### Changed
- Moved YouTube endpoints into their own handler
- Moved jwt_service into it's own package.
- Moved `enterprise_service` into its own pkg called `amapi`.

## [2.24.1] - 2023-07-18  (82.4%)
### Added
- Log issues when creating policies for users in Android Management API
### Fixed
- Issue with AWS SQS de-duplication ID containing invalid characters, such as a whitespace

## [2.24.0] - 2023-07-18  (82.4%)
### Added
- Added scribe summaries table & model
- Scribe summary update endpoint
- External endpoints to create scribe summaries
- Added not_in_group filter to YouTube channels LIST endpoint.
- Added endpoints for owners to change password from MDM-Electron
### Changed
- When creating YouTube channels syncing with the nexus API is not done in a go routine
- Update the issue email response to include the enterprise name.
- Including CanChangePassword field in Enterprises
- Including CanChangePassword and PasswordQuality in AppLoginResponse
- Renamed Transcription model & table to Scribe Transcriptions
- Getting a scribe job now preloads on any summaries

## [2.23.0] - 2023-07-13  (82.5%)
### Added
- Added admin upload YouTube channel image endpoint.
### Fixed
- Purple scribe handlers not respecting impersonated enterprise
- Purple form handlers not respecting impersonated enterprise
- Updated inline functions to work with go v1.20.6
### Changed
- Scribe jobs now return newest first

## [2.22.0] - 2023-07-03  (81.8%)
### Added
- External scribe job handlers for marking a job as complete for failed
- Transcription queue service for sending jobs to be processed
- V2 Upload endpoint for whisper based scribe jobs
- AWS SQS package
### Changed
- Modified 'group - remove YouTube channels' req method: DELETE to POST
- Turned external handler into a folder, and split it into purple post & scribe jobs
- Create scribe job endpoint now takes an option version. Default/V1 = AWS, V2 = Whisper.
### Fixed
- List of applications for electron login required package name of `com.purplevisits.mdm` to be able to store electron activities.

## [2.21.1] - 2023-06-15  (81.9%)
### Added
- Added owner ID into managed config for agent

## [2.21.0] - 2023-06-15  (81.9%)
### Added
- Added PDF Forms table
- Added PDF forms permissions
- Added basic PDF form endpoints - List, Get, GetFile & Upload
- On user create, also check for deleted users by email and show warning
### Fixed
- Fixed issue with AWS Transcribe job names conflicting between dev and prod. Added unix timestamp to the end of job name.

## [2.20.0] - 2023-06-08  (81.8%)
### Added
- scribe_jobs, transcriptions & files tables
- Added in Read & Create permissions for scribe
- Added `TRANSCRIPTION_BUCKET` env variable
- Add cron job to update scribe jobs that are processing
- Added cron endpoint to update processing jobs
- Added scribe jobs endpoints
- Added transcription service for AWS + mock
- Added `github.com/aws/aws-sdk-go/service/transcribeservice` package
### Changed
- Removed the need to pass the db to the files service as it did not need it.
- Added `Upload` to the existing `S3Service`

## [2.19.1] - 2023-06-06  (82.4%)
### Added
- If a user is a sup[er admin then add in ]
### Fixed
- Fixed task handler API to use hte enterpriseUID from the context and not from the user.

## [2.19.0] - 2023-06-06  (82.4%)
### Added
- Added an enterprise seeder.
- Added get issue endpoint.
### Changed
- updated local development URL to localmdm.com. This is to avoid issues with the domain cookie set for the WS server as you can't override a SECURE cookie.
- Updated create issue endpoint to handle general and machine issues.
### Fixed
- Failed owner logins were adding activities saying "Successful Login" rather that "Failed Login". Updated text to fix.
- Updated WS cookie to use enterpriseUID from token for superadmin switch users.

## [2.18.1] - 2023-05-18  (82.4%)
### Changed
- updated WS users endpoint from user auth to cookie auth.

## [2.18.0] - 2023-05-18  (82.4%)
### Added
- Added machine chassis_serial to machine eapi endpoints and to results returned for machine.
- Added websocket server
- Added commands table for the machines
- Added machine chassis_serial to machine eapi endpoints and to results returned for machine.
- Added `MDM_NEXUS_API_KEY` env variable and`MDM_NEXUS_API_URL` env variable.
- Added Nexus API Service and mock.
- Added `youtube_channels` feature and endpoints
- Added list, get and create endpoints for machine commands
- Added WS server Client. Commands added to machines will be sent across the WS server straight away.
### Changed
- Updated machine stats endpoint to add the latest stats and last seen timestamp to machine.
- Implemented UTC as the default timezone.
- Moved `middleware` and `validation` to `pkg` folder so that they can be used by the API and WS.
- Add `youtube_enabled` to groups

## [2.17.3] - 2023-05-02  (83.0%)
### Added
- Added ability to switch enterprises when registering pre-exisitng machine. Old machine is deleted with machine log. New machine is created
### Changed
- Updated EAPI machine middleware to include a filter for Enterprise UID from KEY for machine. When machine try to record stats this will kick back a not found error forcing re-registration

## [2.17.2] - 2023-04-26  (83.1%)
### Added
- Added machine stats endpoints
- Added HL key info into machines list and get endpoints

## [2.17.1] - 2023-04-25  (83.0%)
### Added
- Added endpoints for update and delete for Machines.

## [2.17.0] - 2023-04-24  (82.8%)
### Added
- Added endpoints for mdm-portal users to fetch list of machines, and get individual machines
- Added a library for diffing two objects and storing the results in the activities table.
- Added diff package for saving changes to models.

## [2.16.0] - 2023-04-20  (82.8%)
### Added
- Add filter for retrieval of soft deleted devices
- Deleted existing devices if the same serial is passed though.
- Added EAPI Machine, machine logs and machine stats resources.
- Added registration process for machines
### Changed
- Changed tasks description column from varChar(500) to a text field
- Tracker: Attached links to the end of security tasks description for the Owner and Device
- Tracker: Removed validation for due date when updating task
### Fixed
- Changed filtering for getting top owner activities for portal dashboard

## [2.15.2] - 2023-03-22  (82.8%)
### Added
- Added skipper to nocache tto skip the files endpoint from the middleware
- Added filescache middleware which will add  a years cache onto the image files.
### Fixed
- Preload the file service to ensure the storage location can be found when running in production.
- Added logs to file system for when issues are found. Update error messages to include `error` so that are picked up in warning system.
- Added fatal when starting kiosk app if `API_URL` is not set.

## [2.15.1] - 2023-03-21  (82.8%)
### Fixed
- Set user UID on transaction types for refund.

## [2.15.0] - 2023-03-21  (82.8%)
### Added
- Added delete endpoint for network
- Added ErrorResponseWithMeta to responses
- Shopping: Added authorise and reject endpoints for orders
- Added global css pkg and prepend css to kiosk_css when retrieved

## [2.14.0] - 2023-03-17  (82.7%)
### Added
- Added new feature to groups to allowed for "NoPin" to be set. This will block the setting of the pin on the device. Password setting needs to be set to NONE on the enterprise for this to work.
- Added refresh token duration onto enterprise table and set refresh token expiry time based on that value
- Added file system
- Added Electron application fields to kiosk table with appropriate update endpoint
- Added kiosk tags table and end points
- Cron task now processes all device report types: Active, Provisioning, Deleted
- Added software update check to messages cron job
- Tracker: Added CreateShoppingTask to task service
- Tracker: Added task user handler for returning users that can be assigned to a task
### Changed
- Changed app refresh response to return a slim list of kiosk sites for the owner
- kiosk site create request now requires a name
- Tracker: ask due date now set on creation
- Tracker: Tasks returned in ascending order by due date
### Fixed
- Fixed task Id not being set in 'processing_info' when creating a security task

## [2.13.0] - 2023-03-16  (82.3%)
### Added
- Added in osconfig object to the EAPI policy for controlling the MadepurpleOS applications
### Changed
- Converted magzter endpoint from GET to POST.
- Replaced `go-migrate` with our own version in `pkg`

## [2.12.0] - 2023-01-24  (82.3%)
### Added
- Added CRON Scheduler to process message reports every 5 minutes
- Added cron endpoint for processing device reports
- Added tracker_processing_enabled field to enterprise
- Added tracker task service
- Added device reports application checks
### Changed
- Moved Docker Files into docker folder
- Refactored main.go -> api.go and moved into cmd/api folder
### Fixed
- Added 'Status' to task service 'CreateSecurityTask'

## [2.11.0] - 2023-01-18  (82.3%)
### Added
- List/Get Order endpoints and List endpoint for Order-rows
- Creates Shopping task when owner places an order
- Added Meta model and a Type to Tasks table
- Added Status to Orders table
### Changes
- Changed ModelID (uint) -> ModelUid (string) on task table
- Renamed Meta -> ResponseMeta in responses folder

## [2.10.0] - 2023-01-13  (82.0%)
### Added
- Added custom logger using zerolog which includes the users UID and enterprise ID if set.
### Changed
- Bumped Echo framework from v4.1.16 to v4.9.1 so that we can use custom logging using zero log

## [2.9.0] - 2023-01-13  (77.8%)
### Added
- Added Tasks and Task Messages tables/endpoints/test for the MDM Tracker
- Added 'Task' filter into activities repo
### Changed
- Device enrollment endpoints return error if group has Purple Launcher disabled
### Fixed
- Fixed issue with kiosk being unable to set shipping cost to 0

## [2.8.2] - 2023-01-12  (77.2%)
### Added
- Added options in to policy for keyboard being turned on or off.

## [2.8.1] - 2023-01-12  (77.2%)
### Added
- Added default keyboard enabled flag on the group.

## [2.8.0] - 2023-01-11  (77.2%)
### Added
- Added custom timestamp option for generating owner TOTP (admin-only)
- Added new simple keyboard to enrollment, and disabled the Samsung keyboard for MDM V2

## [2.7.0] - 2022-12-13  (77.4%)
### Added
- Added user ip ranges and owner ip ranges to enterprise model
- Added Ip checks based on enterprise ip settings on owner and user logins
- Added middleware that check users ip with enterprise ip settings
- Added Ip package to contain all ip related functions

## [2.6.0] - 2022-12-13  (77.2%)
### Added
- Added Email to be sent to enterprise shopping email when owner makes shopping order
- Added hashIds into packages
- model, model uid and description added to wallet transaction table
- Added security package and middleware to check user agents and query parameters
### Changed
- Changed orderUids to be generated with the hashIds package

## [2.5.1] - 2022-12-01  (76%)
### Added
- Added purple wallet enabled to app user login response

## [2.5.0] - 2022-11-29  (76%)
### Added
- Added `purple_wallet_enabled` to enterprise model and table (defaults to false)
- Added Endpoints for listing wallet transactions, adding funds and removing funds to owners wallet
- Added Bills table with crud endpoints
- Added endpoint for listing unbilled enterprises
- Added Order amd Order_row table to store owners shopping orders with endpoints
- Added Shopping settings to kiosk model with updated endpoints
### Changed
- Changed refresh endpoint to return `enterprise_wallet_enabled` as part of the user data

## [2.4.0] - 2022-10-20  (75.9%)
### Added
- Return purple account details on owner login as base64'd string.

## [2.3.0] - 2022-09-12 (75.9%)
### Added
- Added `allow_firmware_recovery` option to DB and group. Defaults to False. This option only works with KNOX.

## [2.2.0] - 2022-09-01 (75.8%)
### Added
- Added super-admin CRUD endpoints for EAPI Policies
- Added bulk update endpoint to change all MPOS policies of same kind
- Added `onboard` keyboard settings to the eapi policy. Can create and update with other policy settings
### Changed
- Including Policy ID in EAPI Request and Response

## [2.1.0] - 2022-09-01 (75.4%)
### Added
- Added Super-Admin CRUD endpoints for Purple Accounts
- Added EnterpriseUid filtering to List Purple Accounts
- Added ability to Remove owner from Purple Account
### Changed
- super-admin routes now hanging from /admin, i.e. /enterprises -> /admin/enterprises
- Removed Unique constraint from mdm_name in kiosk_sites table
- Removed DeletedAt column from kiosk_blacklists table

## [2.0.0] - 2022-08-22 (75.1%)
### Added
- Added KNOX_KEY to .env file.
- Send Knox key, factory reset and settings items through managed configuration to the mdm agent application.
- Added TOTP secret key to owner on create and update. Added secret in to policy in managed config too.
- Added TOTP endpoint to fetch for owner.
- Added mdm Version in to each enterprise and a way to update from enterprise update.
- Enterprise version now returned in refresh token

## [1.10.5] - 2022-07-01 (75.3%)
### Added
- Disable Bluetooth in Enterprise Policy

## [1.10.4] - 2022-07-01 (75.3%)
### Added
- Added Factory Reset Protection using email address in each policy.

## [1.10.3] - 2022-07-01 (75.3%)
### Added
- WifiConfigEnabled added to groups and update endpoint.
- Added `isIn` and `isNotIn` for tests to replace contains so that the contains item is shown first in the error report.

## [1.10.2] - 2022-06-20 (75.3%)
### Added
- ResourceWhitelist added to kiosk sites. This whitelist lists what resource URLs the wrapped sites are allowed to load.

## [1.10.1] - 2022-06-20 (75.3%)
### Fixed
- Removed unnecessary middleware on /kiosk/:mdmName route

## [1.10.0] - 2022-06-17 (75.3%)
### Added
- Enabled the Application Logs from Devices within Groups.
- Added a `GC_TOPIC_ID` environment variable to separate topics from subs.
- Kiosk Blacklist and Kiosk Sites with Owner and Admin endpoints have been added.

## [1.9.4] - 2022-06-14  (74.2%)
### Added
- Added secure and encrypted json flags to individual device endpoint.
- Ability to ban some default apps.
### Changed
- Ordered owners desc by created_at

## [1.9.3] - 2022-05-26 (74.6%)
### Added
- Added `create_windows_disabled` flag to group. Added to update endpoint and to enterprise service.

## [1.9.2] - 2022-05-24 (74.6%)
### Changed
- Updated the enterprise policy setting to disable modification of accounts and restrict some other parts.

## [1.9.1] - 2022-05-10 (74.7%)
### Added
- Added code to amke sure a group is empty and has no owners so that it can be deleted.

## [1.9.0] - 2022-05-10
### Added
- Added endpoint to delete a group, as long as there are no users attached.
- Added code to run coverage for the whole project in the makefile.
### Changed
- Moved integration tests into an integration test package so that they can be excluded from coverage

## [1.8.0] - 2022-04-20
### Added
- Shopping Email added to enterprise. Can be set through update endpoint
- Shopping endpoint added for OWners to submit shopping lists
- Threads added as a way of storing shopping lists. Can be changed into a service desk eventually as threads will have comments.
### Fixed
- Made sure the activity count included the latest days activities.

## [1.7.1] - 2022-03-24
### Added
- Added the ability to filter by specific applications for owners activities (9now and 7plus)
- Added Count by date for owner activities so that we can generate count graph.

## [1.7.0] - 2022-03-15
### Added
- Added in Top Owner Activities endpoint based on launched applications.
- Added Launcher Override so that we can set a custom package as the launcher for a policy is set on the enterprise.

## [1.6.1] - 2022-03-08
### Added
- app blacklist added to owner api keys so we cna filter allowed application by API key.
- Stats endpoints added for devices.

## [1.6.0] - 2022-02-16
### Added
- endpoint for fetching all activities for an enterprise with Search and filter by application.
### Changed
- refactored owner activity repository so that we can pass through scopes for use in GET and LIST

## [1.5.0] - 2022-02-11
### Added
- Added confirmation message a logged the command name when a password has been set on a device. Also set it straight away and 45 seconds later.
- Added a Create Issue endpoint `/issues` for portal users.
- Added Mailer and Email Templates
### Changed
- Updated `ResetPasswordDevice` so we can change the password on the device without it instantly locking
- Updated go to golang:1.17.3-alpine builder for production to enable embedding

## [1.4.5] - 2021-12-20
### Added
- Allow soft deleted applications to be restored when creating a new app on an enterprise.
- Added `madepurpleos` and `mp.` as prefixes to application which will be ignored by google enterprise AMAPI.

## [1.4.4] - 2021-12-13
### Added
- Added a `label` to the api keys
### Changed
- Updated page titles and links. Users -> System Users. Owners -> Device Owners. Activities -> Staff Activities.

## [1.4.3] - 2021-12-09
### Added
- Added signing capabilities for CloudFront files. Added new environment variable to pass private key in.
### Changed
- EAPI policies now provide cloudfront signed links under our own custom URL: `files.purplemdm.com`

## [1.4.2] - 2021-12-02
### Changed
- Shortened the API keys to 16 characters. We have brute force protection so this *should be* more than enough

## [1.4.1] - 2021-12-01
### Added
- Added an Enterprise AMI (eami) router group and middleware based on API key in header.
- Added /eami/policy endpoint to pull a policy related to an API key.
### Changed
- Updated golang builder to version golang:1.17.3
- Refactored routes and handlers into separate files based on router groups
- Refactored API middleware and added an API service to group those functions
- Move all the settings that lived on the server into their own singleton config struct which can be imported anywhere.

## [1.3.4] - 2021-11-23
### Added
- EXTERNAL_API_KEY env variable
- External handler & tests
- External handler route to fetch if owners at an enterprise have access to purple post in cell
- API key middleware
- IP range middleware
- Added `GetGroupWithApplications` to the groups repo

## [1.3.3] - 2021-11-05
### Added
- Added Permission grant feature to Applications

## [1.3.2] - 2021-11-03
### Added
- Added `AdvancedSecurityOverrides` with developer options as we were unable to go into debugging mode after google changes.
### Changed
- Removed `AppAutoUpdatePolicy` from group as it's now depreciated.
- Removed `InstallUnknownSourcesAllowed` from policy and instead made it control `UntrustedAppsPolicy`
- Removed `DebuggingFeaturesAllowed` from policy and instead made it control `DeveloperSettings`

## [1.3.1] - 2021-10-29
### Changed
- Removed disable window popups as it might be stopping me debug on the tablet. Also disabled other hard coded fields.

## [1.3.0] - 2021-10-29
### Added
- Added EmmCreateEnterprise process to create enterprises using the new process flow.
- Added managed configuration for Applications

## [1.2.2] - 2021-10-13
### Added
- Set default launcher app to update straight away `appLauncherPolicy.AutoUpdateMode = "AUTO_UPDATE_HIGH_PRIORITY"`

## [1.2.1] - 2021-10-07
### Changed
- Allow owner_uid to also be NULL when inserted to the database. It would only find it previously if it was ""

## [1.2.0] - 2021-09-22
### Added
- Options added to group: MountPhysicalMediaEnabled BluetoothConfigEnabled VpnConfigDisabled ShareLocationEnabled UsbFileTransferEnabled StatusBar SystemNavigation
- Added domain whitelist generator and a list of the domains per application in a json file.

## [1.1.9] - 2021-09-15
### Added
- Added Purple Accounts for enterprises. Accounts are automatically be assigned when an owner is created.
- Added endpoint to manually assign a purple account.
- Added structured logger that prints out JSON for use in Cloudwatch monitoring.
### Fixed
- Fixed issue where NONE for Enterprise is not allowed for a password policy. Had to exclude completely for that option.

## [1.1.8] - 2021-09-14
### Added
- Options for install_unknown_sources_allowed for groups
- Now able to disable launcher app on a group through the API.
- Managed configuration now includes Settings Allow option for showing settings link in app.
- Password requirements added to Enterprise
- Added Update endpoint for Enterprises.
- Added Data protection officer and contact detail to enterprises ready for Auto Signup and New way of creating enterprises to be built in next.
- Added an update_policy for each of the applications added to an enterprise.
- Added option on applications to set the Default runtime permissions (allow, deny, prompt)
### Changed
- Enabled network escape hatch by default, disable safe boot by default, disable debugging features by default.

## [1.1.7] - 2021-09-07
### Changed
- Only send the password set message when we receive the enrollment notification, and not on every message

## [1.1.6] - 2021-09-07
### Added
- Set the `SetConnMaxLifetime` in the DB when created as some threads are being timed out by the server. AWS Mariadb Server default is 8 hours, setting ours to 1 hour.
- Wrapped the startListener in a loop so if it drops for any reason it will restart itself. Printed message out for logging too.
### Changed
- Replaced the system activity messages when owners are logging in with Owner Device activity messages instead as it makes more sense for them to appear there.

## [1.1.5] - 2021-08-23
### Changed
- Nothing. This is just to test the ECR push with Australia

## [1.1.4] - 2021-08-04
### Added
- Added `Update` and `Delete` endpoints for the API Keys

## [1.1.3] - 2021-08-02
### Added
- Added Validation endpoint to check a key exists and the current IP address works for it.
- Added brute force protection for the key validation endpoint
- Added security headers to API.

## [1.1.2] - 2021-07-29
### Added
- Added and Owners Login endpoint
- Added API key facilities for each enterprise to protect the owner's login endpoints.

## [1.1.1] - 2021-07-19
### Added
- Added in no-cache middleware
- Added same site=strict for cookie as well as limiting the path of the refresh cookie to the refresh endpoint.

## [1.1.0] - 2021-07-14
### Added
- Added 2FA to enterprises
- Added password complexity checks.
### Changed
- Removed `newPolicy.WifiConfigDisabled = true` from policy as it overrides network escape hatch.

## [1.0.9] - 2021-05-24
### Added
- Added options for network escape hatch and requiring password on boot.

## [1.0.8] - 2021-05-21
### Added
- Added in all remaining restrictive policies for adding users, changing accounts etc.
### Changed
- Set min password length for owner to 4 as a hack for the Lenovo 4 digit pin issue
### Fixed
- Made sure password is pulled through for Wifi configuration.

## [1.0.7] - 2021-05-20
### Added
- Return the enterprise information in the auth login endpoints
- Added first repo (Enterprise) to the dependencies service. Will migrate some others over as a new way of working.
- Added an environment variable for mocking the Google API rather than it being driven from the environment of dev or prod.
- Added logout endpoint
- Added enterprise Switch endpoint
- Added Enterprise Get endpoint
### Changed
- Now return the refresh token as a HttpOnly cookie in auth endpoints.
- Now accept the refresh token in a cookie for the refresh endpoint. Refresh endpoint is now `GET` rather than `POST`
- For Refresh and Access tokens, when they are generated we now use the enterprise passed through so that we can enable switching.
### Fixed
- Fixed refresh endpoint not returning all permissions for the user

## [1.0.6] - 2021-05-17
### Added
- When changing an owner's password, also loop through their devices syncing the Lock Screen password.
- Set a `PasswordPolicy` to ensure `SOMETHING` is set for a password. Note, `SOMETHING` is a real policy.
- Now added `PolicyEnforcementRules` to each device which forces a password to be set, or the device ot be blocked
- Added Lock screen message which asks for the password and prints out the owner's identifier. Will help to identify tablets.
- Added encryption policy
- On ENROLLMENT or first message from device through PubSub, we now Sync the owner's password to the device.
### Changed
- QueueServices now accept dependency service rather than DB so that it can now access other services.

## [1.0.5] - 2021-05-12
### Added
- /apps/me endpoint for Owners to get personal info.

## [1.0.4] - 2021-05-11
### Added
- Add the app name to the owner activity when the logging comes in.
- Added Magzter endpoints and tests.

## [1.0.3] - 2021-05-06
### Added
- Search on the owner list endpoint by identifier *
- Profile endpoints (update name and email, and update password) *
- Added Endpoint to change the policy on a device *
- Owner passwords, update owner password endpoint
- Refactored to set a policy for the owner instead of for the group. Still based on the group policies though.
- Added authentication for Owners using their own token JWTs which are passed through managed configuration
- Added endpoint to log activities from owners
### Changed
- Made the QR code bigger by default (400px)
- For the tests set the bcrypt library to a cost of 4 to speed up encryption.

## [1.0.2] - 2021-04-12
### Added
- Settings for groups including kiosk security options
- Auto-Run for PubSub Messages
- Added certificates to docker scratch build
- Added switch for enterprise iFrame url

## [1.0.1] - 2021-03-31
### Added
- Initial version built.



# Notes
[Deployment] Notes for deployment
[Added] for new features.
[Changed] for changes in existing functionality.
[Deprecated] for once-stable features removed in upcoming releases.
[Removed] for deprecated features removed in this release.
[Fixed] for any bug fixes.
[Security] to invite users to upgrade in case of vulnerabilities.
[YANKED] Note the emphasis, used for Hotfixes
