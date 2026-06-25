# bbs-go Docker Offline Package

This directory is copied into `dist/bbs-go-docker-offline` by `scripts/package-docker-offline.sh`.

On the target machine:

```sh
./deploy.sh
```

The package contains:

- `images/bbs-go-offline.tar`
- `images/mysql-8.4.tar`
- `docker-compose.yaml`
- `load-images.sh`
- `deploy.sh`

Runtime files are stored in `docker-data/`, including MySQL data.
