docker build -f databaseConnector/Dockerfile -t gotmdbsql .
docker create -ti --name dummy gotmdbsql bash
mkdir out
docker cp dummy:/go/databaseConnector/SaveTrendsOfWeek.exe ./out/SaveTrendsOfWeek.exe
docker rm -f dummy