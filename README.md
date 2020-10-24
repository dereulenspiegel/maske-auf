# maske-auf

Dies ist ein kleiner Service, der mit einer URL gefüttert werden kann, die eine GeoJSON Featurecollection zurück gibt.
Der Dienst nimmt an, dass in jedem Feature draußen Maskenpflicht herrscht.
Der Standort der Nutzer wird mittels der Javascript Geolocation API ermittelt und an das Backend übertragen.
Danach zeigt die Seite dem Nutzer ob auch draußen Masken verpflichtend sind oder nicht.

# Building

Der Dienst benötigt packr2 und go um gebaut zu werden.

Schritte

* go installieren (https://golang.org/dl/)
* packr2 installieren `go get -u github.com/gobuffalo/packr/v2/packr2`
* Make installieren, falls noch nicht vorhanden
* `make build` oder `make build-linux` aufrufen. `make build` baut für dein aktuelles Betriebssystem. `make build-linux` immer für Linux mit Arch AMD64

# Benutzung

Der Dienst kann einfach gestartet werden benötigt aber ein paar Parameter

| Parameter | Default | Verpflichtend? | Bedeutung |
|-----------|---------|----------------|-----------|
| `-url`    | keiner  | ja             | Die HTTPS-URL zu der GeoJSON Feature Collection 
| `-interval` | `5m`    | nein    | In welchen Intervallen die GeoJSON Feature Collection erneut heruntergeladen werden soll 
| `-addr`   | `:8080`  | nein          | Die Listen-Addr auf denen der Server lauscht. Per Default ist das alle Adressen auf Port 8080

# Hinweise

Dieser Dienst muss per HTTPs erreichbar sein. Man sollte ihn also hinter einem Reverse Proxy wie nginx, Traefik, Caddy 
etc. betreiben.

# ToDo

Die Seite für Benutzer ist aktuell sehr minimal. Wer gut mit CSS umgehen kann und Lust hat eine schöne, barrierefreie und 
responsive Seite zu bauen ist willkommen. Einfach eine PR erstellen.
