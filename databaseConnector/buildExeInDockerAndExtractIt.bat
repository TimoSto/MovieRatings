docker build -t gotmdbsql .
docker create -ti --name dummy gotmdbsql bash
docker cp dummy:/go/databaseConnector/GetTMDbDataOfWeek.exe ./GetTMDbDataOfWeek.exe
docker rm -f dummy