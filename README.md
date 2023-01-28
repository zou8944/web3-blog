# Web3-blog

This is a DBlog

- Sending email to create, update, or delete blog
- Login with wallet to admin page to manage your article
- Store your data in IPFS and ArWeave
- View your articles as normal blog site

## Build

```shell
make build-docker
```

## Run

### Prepare

- Configure SES, SNS, SQS in AWS, make sure you can send email to SQS
- Register a wallet account, get key file (keyfile.json)

### CLI

- Copy config/default.yaml to config/{your-env}.yaml, replace your configuration
- Copy keyfile.json to project root directory
- Run `ENV={your-env} go run main.go`

### Docker

- Prepare a directory A to mount to container, such as `~/web3-blog/config`
- Copy config/default.yaml to A/{your-env}.yaml, replace your configuration field
- Copy keyfile.json to A
- Run `sudo docker run --name web3-blog -v A:/config -e ENV={your-env} {image-name}:latest`


## Job to do

- [x] Email receiving and sending over AWS SES, SNS, SQS
- [x] Parse eml file ([RFC5322](https://www.rfc-editor.org/rfc/rfc5322))
- [x] Login with [METAMASK](https://metamask.io/)
- [x] Store in IPFS
- [x] Store in arweave
- [x] front end page
- [x] dockerize
- [ ] blog migration from hexo or other platform
- [ ] admin page
- [ ] Sqlite file backup
