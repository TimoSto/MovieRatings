# MovieRatings
Die Corona-Pandemie hat vieles verändert. Doch hat sie auch einen Einfluss auf die Gewohnheiten der Filmliebhaber? Schauen die Menschen jetzt mehr Komödien oder mehr Distopien-Filme? 
Um diese Frage zu beantworten, können mit dieser Dashboard-Anwendung die Anzahl der in jeder Woche auf TMDb geloggten Filme je nach Genre mit dem Verlauf der Corona-Zahlen verglichen werden.
## Systemvoraussetzungen

## Datenbank
Um die mySQL-Datenbank zu füllen, wird ein Kommandozeilen-Programm geschrieben in Golang verwendet. Dieses erledigt folgende Aufgaben:
- Die ersten 100 Filme aus den TMDb-Trends der aktuellen Woche ermitteln
- Ggf. die Einträge in der Movies-Tabelle ergänzen
- Ggf. werden die Einträge in der Genres-Tabelle ergänzt
- Ggf. die Einträge in der MovieGenre-Tabelle ergänzen
- GGf. die Einträge in der MovieWeekPopularity-Tabelle ergänzen
Die SQL-Befehle werden in die Datei `FILLDB.sql` geschrieben, sodass der selbe Datenbank-Zustand erreicht werden kann, wenn man diese Datei ausführt.
### Bauen der Anwendung
Die Anwendung kann entweder direkt gebaut/gestartet werden, wenn man lokal Golang installiert hat, oder in Docker. Dafür wird über die Batch-Datei `buildExeInDockerAndExtractIt.bat` ein Linux-Image gebaut, in welchem die EXE liegt. Dann wird ein Container mit diesem Image gestartet und die EXE wird auf den lokalen Rechner kopiert.