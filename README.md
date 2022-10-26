# Web3-blog

This is a DBlog.

- Sending email to create, update, or delete blog
- Login with wallet
- Store your data in ipfs and arweave
- View your articles like normal blog sites

## Debug

- Configure SES, SNS, SQS in AWS, make sure you can send email to SQS
- Update config/default.yaml, replace your property
- Run project with command `air`

## I am in developing

### Completed job

- Email receiving and sending over AWS SES, SNS, SQS
- Parse eml file ([RFC5322](https://www.rfc-editor.org/rfc/rfc5322))
- Login with [METAMASK](https://metamask.io/)

### To be completed

- Store in IPFS
- Store in arweave
- Sqlite file backup