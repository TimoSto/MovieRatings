docker build -f saveTrendsToDB/Dockerfile -t gotmdbsql .
docker run -d -ti --name dummy gotmdbsql
docker cp dummy:/tmp/SaveTrendsOfWeek.exe out/SaveTrendsOfWeek.exe
docker cp dummy:/tmp/DBHandler.exe out/DBHandler.exe
docker stop dummy
docker rm dummy