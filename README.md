# MovieRatings
Die Corona-Pandemie hat vieles verändert. Doch hat sie auch einen Einfluss auf die Gewohnheiten der Filmliebhaber? Schauen die Menschen jetzt mehr Komödien oder mehr Distopien-Filme? 
Um diese Frage zu beantworten, können mit dieser Dashboard-Anwendung die Anzahl der in jeder Woche auf TMDb geloggten Filme und Serien sowie die Infos über die beliebten Personen in verschiedenen Verknüpfungen dargestellt werden, während immer die Veränderungen der (deutschen) Coronazahlen zum Vergleich ersichtlich sind.  
Eine detailiertere Dokumentation der technischen und der organisatorischen Aspekte dieser Anwendung ist im [Wiki](https://github.com/TimoSto/MovieRatings/wiki) zu finden.
## Trends in DB speichern
1. ggf. `config.json` aktualisieren
2. im out-Verzeichnis `SaveTrendsOfWeek.exe` über CMD aufrufen
## Coronazahlen eintragen
Im out-Verzeichnis:  
`DBHandler --covid week=26 cases=12345678 deaths=12345 recovered=546327`
## Anwendungen bauen
Nur notwendig, wenn der Quellcode verändert wurde:
1. Docker installieren und auf Linux-Container einstellen
2. `build.bat`
