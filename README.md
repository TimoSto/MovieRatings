# MovieRatings
Die Corona-Pandemie hat vieles verändert. Doch hat sie auch einen Einfluss auf die Gewohnheiten der Filmliebhaber? Schauen die Menschen jetzt mehr Komödien oder mehr Distopien-Filme? 
Um diese Frage zu beantworten, können mit dieser Dashboard-Anwendung die Anzahl der in jeder Woche auf TMDb geloggten Filme je nach Genre mit dem Verlauf der Corona-Zahlen verglichen werden.
## Informationen, die dargestellt werden (können)
- Verlauf der Coronazahlen
- Trends der jeweiligen Genres (Film oder TV) in jeder Woche (Als Balkendiagramm für Genre oder Kuchendiagramm für Woche)
- Streaming-Anbieter, welche die Trends zur Verfügung stellen
- Produktionsländer der Trends
- Networks der Serien
- Cast und Crew der Filme/Serien
- Beliebte Personen in der jeweiligen Woche
- Filme/Serien, in denen die Person mitgewirkt hat
## Systemvoraussetzungen
### Datenbank
- mySQL-Server
- Golang oder Docker (wenn die Anwendung neu gebaut werden soll, sonst kann die EXE verwendet werden)
## Datenbank
### Setup
Um die Datenbank aufzusetzen muss zunächst das SQL-Script `database/STARTUP_DB.sql` z.B. in der Workbench eines MySQL-Servers ausgeführt werden. Dies erzeugt die DB `movieratings`. Dann werden mit dem Skript `database/CREATE_DB.sql` die notwendigen Tabellen angelegt. Um die bisher abgerufenen Daten in diese Tabellen zu schreiben, kann das Skript `database/FILLDB.sql` verwendet werden.
### Transfer der Daten von der API in MySQL
Um die mySQL-Datenbank zu füllen, wird ein Kommandozeilen-Programm (`out/SaveTrendsOfWeek.exe`) geschrieben in Golang (`saveTrendsToDB/SaveTrendsOfWeek.go`) verwendet. Dieses erledigt folgende Aufgaben:
- Die ersten 100 Filme aus den TMDb-Trends der aktuellen Woche ermitteln
- Ggf. die Einträge in der Movies-Tabelle ergänzen
- Ggf. werden die Einträge in der Genres-Tabelle ergänzt
- Ggf. die Einträge in der MovieGenre-Tabelle ergänzen
- Ggf. werden die Einträge in der Personen-Tabelle ergänzt
- Ggf. werden die Einträge in der MovieCredits-Tabelle ergänzt
- Ggf. Countries-Tabelle ergänzen
- Ggf. MovieCountry-Tabelle ergänzen
- Ggf. Provider-Einträge ergänzen
- Ggf. Movie-Provider-Einträge ergänzen
- GGf. die Einträge in der MovieWeekPopularity-Tabelle ergänzen
- Die ersten 100 Serien aus den TMDb-Trends der aktuellen Woche ermitteln
- Ggf. die Einträge in der Series-Tabelle ergänzen
- Ggf. werden die Einträge in der Genres-Tabelle ergänzt
- Ggf. die Einträge in der SeriesGenre-Tabelle ergänzen
- Ggf. werden die Einträge in der Personen-Tabelle ergänzt
- Ggf. werden die Einträge in der SeriesCredits-Tabelle ergänzt
- Ggf. Countries-Tabelle ergänzen
- Ggf. SeriesCountry-Tabelle ergänzen
- Ggf. Provider-Einträge ergänzen
- Ggf. Series-Provider-Einträge ergänzen 
- Ggf. Network-Einträge ergänzen
- Ggf. Series-Network-Einträge ergänzen  
- GGf. die Einträge in der SeriesWeekPopularity-Tabelle ergänzen
- Die beliebten Personen werden in der Personen-Tabelle ergänz
- Ggf. wird die PersonWeek-Tabelle ergänzt
Die SQL-Befehle werden in die Datei `FILLDB.sql` geschrieben, sodass der selbe Datenbank-Zustand erreicht werden kann, wenn man diese Datei im Kontext der Datenbank ausführt.
### Bauen der DB-Anwendung
Die Anwendung kann entweder direkt gebaut/gestartet werden, wenn man lokal Golang installiert hat, oder in Docker. Dafür wird über die Batch-Datei `build.bat` ein Linux-Image gebaut, in welchem die EXE liegt. Dann wird ein Container mit diesem Image gestartet und die EXE wird auf den lokalen Rechner kopiert. Die EXE liegt dann im Out-Verzeichnis und kann von dort gestartet werden.