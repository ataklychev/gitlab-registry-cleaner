[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/traefik/traefik/blob/master/LICENSE.md)

Gitlab container registry cleaner
---

## Features
Periodic removes old images from GitLab container registry in all projects.

## Install

- Create access token in your gitlab service.
- Run from source ```go run main.go clean```
- Run from docker image ```docker run -e BASE_API_URL=https://gitlab.com/api/v4 -e ACCESS_TOKEN=XXX ataklychev/gitlab-registry-cleaner```

## Environment variables
- DEBUG=true ( configure log format )
- THRESHOLD=3 ( images over threshold will be deleted automatically )
- BASE_API_URL=https://gitlab.com/api/v4 ( gitlab api endpoint )
- ACCESS_TOKEN=XXX ( gitlab access token, see https://docs.gitlab.com/ee/user/profile/personal_access_tokens.html#create-a-personal-access-token )
- CRON_TIME=01:11 ( time of day to run clean )

## Run clean command in docker
Single operation, clean container registry and exit
```
docker run -e BASE_API_URL=https://gitlab.com/api/v4 -e ACCESS_TOKEN=XXX ataklychev/gitlab-registry-cleaner
```

## Run cron command in docker
Run clean container registry operation every day at specified time
```
docker run -e BASE_API_URL=https://gitlab.com/api/v4 -e ACCESS_TOKEN=XXX ataklychev/gitlab-registry-cleaner cron
```