# Minio Synchronize CLI

Minio Sync is an command line tool for synchronize a list of buckets from source to destination MinIO Server


## Configuration

minio-sync uses yaml config file (default is $HOME/.minio-sync.yaml)

```yaml
src:
  endpoint: localhost:9000
  accessKey: minioadmin
  secretKey: minioadmin
  useSSL: true
dest:
  endpoint: minio.dev:9000
  accessKey: minioadmin
  secretKey: minioadmin
  useSSL: false
```

## Usage

```
minio-sync [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  download    A brief description of your command
  help        Help about any command
  remoteSync  A brief description of your command
  upload      A brief description of your command

Flags:
      --config string   config file (default is $HOME/.minio-sync.yaml)
  -h, --help            help for minio-sync
```

### Download 

Download objects from source MinIO Server to local directory

```
Usage:
  minio-sync download [flags]

Flags:
  -d, --directory string   Local directory to save objects (default "tmp")
  -f, --file string        List bucket file name (.csv) (default "test.csv")
  -h, --help               help for download

Global Flags:
      --config string   config file (default is $HOME/.minio-sync.yaml)
```

### Upload 

Upload files from local directory to destination MinIO Server.

```
Usage:
  minio-sync upload [flags]

Flags:
  -d, --directory string   Local directory to get files (default "tmp")
  -f, --file string        List bucket file name (.csv) (default "test.csv")
  -h, --help               help for upload

Global Flags:
      --config string   config file (default is $HOME/.minio-sync.yaml)
```

### Remote Sync 

Synchronize buckets directly from source to destination MinIO Server.

```
Usage:
  minio-sync remoteSync [flags]

Flags:
  -f, --file string   List bucket file name (.csv) (default "test.csv")
  -h, --help          help for remoteSync

Global Flags:
      --config string   config file (default is $HOME/.minio-sync.yaml)
```

