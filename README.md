# DSC SDEI Project

## Build
```
docker build -t dsc . # Disable the import of access log if using nginx service
```

## Run
```
docker run --name dsc-container
```

## Check Stats
```
docker ps # Get container ID
docker exec -i -t ID /usr/bin/tail -f /var/log/stats.log
```

You can append data to access to /var/log/nginx/access.log to watch it update
