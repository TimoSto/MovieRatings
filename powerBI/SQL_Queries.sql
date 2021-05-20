/*
 Hier können wir die SQL-Queries, die wir dann für PowerBI nutzen, sammeln
 */

/* Covid Differenzen zum vorherigen Monat) */
select weekNr, cases-LAG(cases) OVER (ORDER BY weekNr ASC), deaths-LAG(deaths) OVER (ORDER BY weekNr ASC), recovered-LAG(recovered) OVER (ORDER BY weekNr ASC) from covid

/*Fälle-Tote-Genesene*/
select weekNr, cases-recovered-deaths from covid
/*Film-Trends aufgeschlüsselt nach Genre und Woche*/
select mwp.weekNr, sum(mwp.popularity), count(*), mg.genreId, g.genre from movieweekpopularity as mwp
inner join movies as m on m.id = mwp.movieId
inner join moviegenre as mg on mg.movieId = m.id
inner join genres as g on g.id = mg.genreId
group by mwp.weekNr, g.id
order by mwp.weekNr
