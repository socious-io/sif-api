# Socious Fund

Transforming crowdfunding through extended quadratic funding—with zero operational costs.

## Overview

Socious Fund revolutionizes crowdfunding with an effective quadratic voting system. By combining small contributions from socially impactful individuals with larger matching funds, we create an equitable funding platform. Our unique approach amplifies community voices through a robust impact scoring system and secure identity verification.

## Key Features

- **Democratic Decision-Making**: Small donations combine to create greater impact
- **Fair Resource Distribution**: Prevents concentration of funding power
- **Impact Optimization**: Funding allocated by community support and verified social impact
- **Transparent Tracking**: Blockchain-based funding tracking system
- **Milestone-based Distribution**: Funds released as goals are achieved

## User Stories

1. Log in and sign up
2. Complete the KYB and KYC process
3. Create, publish, and receive funds for a project
4. Vote for projects
5. Receive impact points and achievement badges (on Socious ID platform)

## Get Started

Visit https://socious.org to learn more about how you can participate in transforming crowdfunding for social impact.

## Documentation

For detailed information about our impact score calculation and methodology, visit our [Whitepaper](https://socious.gitbook.io/whitepaper/impact-score/how-is-impact-score-calculated).
For tutorials, please visit this link: https://socious.gitbook.io/fund

## Documentation

For detailed information about our impact score calculation and methodology, visit our [Whitepaper](https://socious.gitbook.io/whitepaper/impact-score/how-is-impact-score-calculated).
For tutorials, please visit this link: https://socious.gitbook.io/fund


## Database Query system
**Temp Note**: *would update db query related documents here*

## Migration system
**Temp Note**: *would update migration related documents here*

## Quick start
**should take care of matching config file to related connection such as pg and nats**
```
$ cd sif-api
$ cp .tmp.config.yml config.yml
$ sudo docker-compose up -d
$ go get
$ go run cmd/migration/main.go up
$ go run cmd/app/main.go
``` 
