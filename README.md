[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/traefik/traefik/blob/master/LICENSE.md)

Gitlab container registry cleaner
---

## Features
Periodic removes old images from GitLab container registry in all projects.

## Install

- Create access token in your gitlab service.
- Set required params (BASE_API_URL, ACCESS_TOKEN) in configuration file ```cp ./.env.example ./.env```.
- 

## Environment variables
- PRODUCTION=false ( configure log format )
- THRESHOLD=3 ( images over threshold will be deleted automatically )
- BASE_API_URL=https://gitlab.com/api/v4 ( gitlab api endpoint )
- ACCESS_TOKEN=XXX ( gitlab access token, see https://docs.gitlab.com/ee/user/profile/personal_access_tokens.html#create-a-personal-access-token )
- CRON_TIME=01:11 ( time of day to run clean )

## Run clean command in docker
Single clean container registry operation
```
docker run \
--name gitlab-registry-cleaner \
--env-file=./.env \
--env COMMAND_NAME=clean \
ataklychev/gitlab-registry-cleaner
```

## Run cron command in docker
Run clean container registry operation every day at specified time
```
docker run \
--name gitlab-registry-cleaner-cron \
--env-file=./.env \
--env COMMAND_NAME=cron \
ataklychev/gitlab-registry-cleaner
```