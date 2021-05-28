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

/*Popularität der Film-Genres jede Woche*/
select mwp.weekNr, sum(mwp.popularity) as Popularität, count(*) as AnzahlTrends, mg.genreId, g.genre as Genre from movieweekpopularity as mwp
inner join movies as m on m.id = mwp.movieId
inner join moviegenre as mg on mg.movieId = m.id
inner join genres as g on g.id = mg.genreId
group by mwp.weekNr, g.id
order by mwp.weekNr asc, sum(mwp.popularity) desc

/*Popularität der Serien-Genres jede Woche*/
select swp.weekNr, sum(swp.popularity), count(*), mg.genreId, g.genre from seriesweekpopularity as swp
inner join series as s on s.id = swp.seriesId
inner join seriesgenre as mg on mg.seriesId = s.id
inner join genres as g on g.id = mg.genreId
group by swp.weekNr, g.id
order by swp.weekNr asc, sum(swp.popularity) desc

/*Gender-Verteilung in Film-Trends*/
select mwp.weekNr, count(*), p.gender from movieweekpopularity as mwp
inner join moviecredits as mc on mc.movieId = mwp.movieId
inner join personen as p on p.id = mc.personId
group by mwp.weekNr, p.gender
order by mwp.weekNr asc

/*Gender-Verteilung in Serien-Trends*/
select swp.weekNr, count(*), p.gender from seriesweekpopularity as swp
inner join seriescredits as mc on mc.seriesId = swp.seriesId
inner join personen as p on p.id = mc.personId
group by swp.weekNr, p.gender
order by swp.weekNr asc

/*Verteilung der Film-Trends auf die Länder*/
select mwp.weekNr, count(*) as AnzahlTrends, c.cname as Land, c.id as ISO_ID from movieweekpopularity as mwp
inner join moviecountry as mc on mc.movieId = mwp.movieId
inner join countries as c on c.id = mc.countryId
group by mwp.weekNr, c.id
order by mwp.weekNr asc

/*Verteilung der Seiren-Trends auf die Länder*/
select swp.weekNr, count(*) as AnzahlTrends, c.cname as Land, c.id as ISO_ID from seriesweekpopularity as swp
inner join seriescountry as sc on sc.seriesId = swp.seriesId
inner join countries as c on c.id = sc.countryId
group by swp.weekNr, c.id
order by swp.weekNr asc

/*Covid-Zahlen als String in Tsd.-Einheiten*/
select weekNr, 
concat(round((cases-LAG(cases) OVER (ORDER BY weekNr ASC))/1000, 1), " Tsd.") as cases, 
convert(deaths-LAG(deaths) OVER (ORDER BY weekNr ASC), char) as deaths, 
concat(round((recovered-LAG(recovered) OVER (ORDER BY weekNr ASC))/1000, 1), " Tsd.") as recovered 
from covid

/*Trends pro Woche und Land*/
select mwp.weekNr, mwp.popularity, m.title, c.cname as Land, c.id as ISO_ID from movieweekpopularity as mwp
inner join movies as m on m.id = mwp.movieId
inner join moviecountry as mc on mc.movieId = mwp.movieId
inner join countries as c on c.id = mc.countryId
order by mwp.weekNr asc, ISO_ID

/*Trends-Country-Week*/
SELECT movies.title, moviecountry.countryId, countries.cname, movieweekpopularity.weekNr, movieweekpopularity.movieId, movieweekpopularity.popularity, sub.AnzahlTrends FROM movies INNER JOIN moviecountry ON movies.id = moviecountry.movieId INNER JOIN countries ON countries.id = moviecountry.countryId 
INNER JOIN movieweekpopularity ON movieweekpopularity.movieId = movies.id 
left join (select mwp.weekNr, count(*) as AnzahlTrends, c.cname as Land, c.id as ISO_ID, mwp.movieId from movieweekpopularity as mwp
inner join moviecountry as mc on mc.movieId = mwp.movieId
inner join countries as c on c.id = mc.countryId
group by mwp.weekNr, c.id
order by mwp.weekNr asc) as sub on moviecountry.countryId = sub.ISO_ID and sub.weekNr = movieweekpopularity.weekNr and sub.movieId = movieweekpopularity.movieId
ORDER BY cname, movieweekpopularity.weekNr, movieweekpopularity.popularity

/*Alter in Movie-Trends*/
select mwp.weekNr, count(*),
case when timestampdiff(YEAR,STR_TO_DATE( p.birthday, '%Y-%m-%d'), STR_TO_DATE( m.releaseDate, '%Y-%m-%d')) >= 80 then '>80'
when timestampdiff(YEAR,STR_TO_DATE( p.birthday, '%Y-%m-%d'), STR_TO_DATE( m.releaseDate, '%Y-%m-%d')) >= 70 then '70-79'
when timestampdiff(YEAR,STR_TO_DATE( p.birthday, '%Y-%m-%d'), STR_TO_DATE( m.releaseDate, '%Y-%m-%d')) >= 60 then '60-69'
when timestampdiff(YEAR,STR_TO_DATE( p.birthday, '%Y-%m-%d'), STR_TO_DATE( m.releaseDate, '%Y-%m-%d')) >= 50 then '50-59'
when timestampdiff(YEAR,STR_TO_DATE( p.birthday, '%Y-%m-%d'), STR_TO_DATE( m.releaseDate, '%Y-%m-%d')) >= 40 then '40-49'
when timestampdiff(YEAR,STR_TO_DATE( p.birthday, '%Y-%m-%d'), STR_TO_DATE( m.releaseDate, '%Y-%m-%d')) >= 30 then '30-39'
when timestampdiff(YEAR,STR_TO_DATE( p.birthday, '%Y-%m-%d'), STR_TO_DATE( m.releaseDate, '%Y-%m-%d')) >= 20 then '20-29'
else '0-19'
 end as agegroups
from movieweekpopularity as mwp
inner join movies as m on m.id = mwp.movieId
inner join moviecredits as mc on mc.movieId = m.id
inner join personen as p on p.id = mc.personId
where p.birthday != ''
group by mwp.weekNr, agegroups
order by mwp.weekNr
/*Beliebteste Personen der woche nach Gender*/
select weekNr, p.gender, p.name from personweek as pw0
inner join personen as p on p.id = pw0.personId
where personId in (select p.id from personweek as pw
inner join personen as p on p.id = pw.personId
where pw.popularity = (Select max(pw2.popularity) from personweek as pw2 inner join personen as p2 on p2.id = pw2.personId where pw2.weekNr = pw.weekNr and gender = 1 and pw.weekNr = pw0.weekNr)
order by pw.weekNr asc, pw.popularity desc) or personId in (select p.id from personweek as pw
inner join personen as p on p.id = pw.personId
where pw.popularity = (Select max(pw2.popularity) from personweek as pw2 inner join personen as p2 on p2.id = pw2.personId where pw2.weekNr = pw.weekNr and gender = 2 and pw.weekNr = pw0.weekNr)
order by pw.weekNr asc, pw.popularity desc)
order by pw0.weekNr asc, p.gender desc