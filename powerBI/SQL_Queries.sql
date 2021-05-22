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
/*Movie-Week-Popularity*/
select mwp.weekNr, mwp.movieId, m.title, mwp.voteAvg, mwp.popularity, mwp.voteCount from movieweekpopularity as mwp
inner join movies as m on mwp.movieId = m.id
order by mwp.weekNr asc, mwp.popularity desc
/*Series-Week-Popularity*/
select swp.weekNr, swp.seriesId, s.title, swp.voteAvg, swp.popularity, swp.voteCount from seriesweekpopularity as swp
inner join series as s on swp.seriesId = s.id
order by swp.weekNr asc, swp.popularity desc
/*Top-Movie jeder Woche*/
select mwp.weekNr, mwp.movieId, concat(mwp.voteAvg, "/10"), m.title, GROUP_CONCAT(g.genre SEPARATOR ', '), mwp.popularity as max from movieweekpopularity as mwp
inner join movies as m on mwp.movieId = m.id
inner join moviegenre as mg on m.id = mg.movieId
inner join genres as g on g.id = mg.genreId
where mwp.popularity = (Select max(mwp2.popularity) from movieweekpopularity as mwp2 where mwp2.weekNr = mwp.weekNr)
group by mwp.weekNr
/*Top-Serie jeder Woche*/
select swp.weekNr, swp.seriesId, concat(swp.voteAvg,"/10"), s.title, GROUP_CONCAT(g.genre SEPARATOR ', '), swp.popularity as max from seriesweekpopularity as swp
inner join series as s on swp.seriesId = s.id
inner join seriesgenre as sg on s.id = sg.seriesId
inner join genres as g on g.id = sg.genreId
where swp.popularity = (Select max(swp2.popularity) from seriesweekpopularity as swp2 where swp2.weekNr = swp.weekNr)
group by swp.weekNr