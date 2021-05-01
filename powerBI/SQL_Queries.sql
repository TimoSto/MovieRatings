/*
 Hier können wir die SQL-Queries, die wir dann für PowerBI nutzen, sammeln
 */

/* Covid Differenzen zum vorherigen Monat) */
select weekNr, cases-LAG(cases) OVER (ORDER BY weekNr ASC), deaths-LAG(deaths) OVER (ORDER BY weekNr ASC), recovered-LAG(recovered) OVER (ORDER BY weekNr ASC) from covid

/*Fälle-Tote-Genesene*/
select weekNr, cases-recovered-deaths from covid