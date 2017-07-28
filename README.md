# DSC SDEI Project

## Build
```
docker build -t dsc . # Disable the import of access log if using nginx service
```

## Check Stats
```
docker ps # Get container ID
docker exec -i -t ID /bin/tail -f /var/log/stats.log
```

You can append data to access to /var/log/nginx/access.log to watch it update
