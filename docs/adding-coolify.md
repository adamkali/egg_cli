---
title: adding-coolify
description: Adding coolify
author: adamkali
categories: 
  - Egg
  - Golang
  - CLI
  - Add-ons
  - Wiki
created: 2025-08-16T18:34:51-0500
updated: 2025-09-03T18:24:43-0500
version: 1.1.1
---


# Adding Coolify with a webhook. 

If you want to be able to deploy to coolify with a webhook, after building the docker image via github container registry, there are some steps that you have to work through. 
1. Create a new coolify project
2. Create a new webhook for the project
3. Get the webhook url for the project in the webhook settings
4. Add the coolify webhook url to the github repo as a secret 
5. Create the coolify token in coolify's api tab
6. Add the coolify token to the github repo as a secret
7. Add the following yaml to the github worflows directory (`.github/workflows`)
8. Finally make sure that the coolify instance has premissions to interact with the [github container registry](https://coolify.io/docs/knowledge-base/docker/registry)

```yaml
on:
  push:
    branches:
      - main
permissions:
  contents: read
  packages: write
jobs:
  build-image:
    runs-on: ubuntu-latest
    needs:
      - build-and-test
    steps:
      - name: Checkout repository
        uses: actions/checkout@v3
      - name: Log in to the Container registry
        uses: docker/login-action@f054a8b539a109f9f41c372932f1ae047eff08c9
        with:
          registry: https://ghcr.io
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}
      - name: Build and push Docker image
        uses: docker/build-push-action@ad44023a93711e3deb337508980b4b5e9bcdc5dc
        with:
          context: .
          push: true
          tags: ghcr.io/adamkali/fullstack_app:stable

  deploy:
    runs-on: ubuntu-latest
    needs:
      - build-and-test
      - build-image
    steps:
      - name: Checkout repository
        uses: actions/checkout@v3
      - name: Deploy to coolify
        run: |
  curl --request GET "${{ secrets.COOLIFY_WEBHOOK }}" --header "Authorization: Bearer ${{ secrets.COOLIFY_TOKEN }}"
```

