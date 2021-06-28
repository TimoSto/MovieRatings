***Abfrage für die Anzahl der Top100 Filme je Provider (die in Flat enthalten sind)***
SELECT pname, COUNT(movieId) AS AnzahlFilmeInFlat FROM `movieprovider` AS mp 
INNER JOIN providers AS p ON mp.provider = p.id 
WHERE service = 'flat' GROUP BY mp.provider

***Abfrage für die Anzahl der Top100 Serien je Provider (die in Flat enthalten sind)***
SELECT pname, COUNT(seriesId) AS AnzahlSerienInFlat FROM seriesprovider AS sp 
INNER JOIN providers AS p ON sp.provider = p.id 
WHERE service = 'flat' GROUP BY sp.provider

***Wo kann man den in der aktuellen KW beliebtesten Film anschauen?***
select mwp.weekNr, mwp.movieId, concat(mwp.voteAvg, "/10"), m.title, GROUP_CONCAT(p.pname SEPARATOR ', '), mwp.popularity as max from movieweekpopularity as mwp
inner join movies as m on mwp.movieId = m.id
inner join movieprovider as mp ON mp.movieId = m.id
INNER JOIN providers AS p ON mp.provider = p.id
where mwp.popularity = (Select max(mwp2.popularity) from movieweekpopularity as mwp2 where mwp2.weekNr = mwp.weekNr)
group by mwp.weekNr
---> !Problem: der jeweils beliebteste Film ist nicht in der Tabelle movieprovider zu finden

---Serien --> selbes Problem
select swp.weekNr, swp.seriesId, concat(swp.voteAvg,"/10"), s.title, GROUP_CONCAT(p.pname SEPARATOR ', '), swp.popularity as max from seriesweekpopularity as swp
inner join series as s on swp.seriesId = s.id
inner join seriesprovider as sp on sp.seriesId = s.id
inner join providers as p on p.id = sp.provider
where swp.popularity = (Select max(swp2.popularity) from seriesweekpopularity as swp2 where swp2.weekNr = swp.weekNr)
group by swp.weekNr
