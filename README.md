# pidia == PiDiaShow Monitor
## Benutzung
### Pidiarc
#### Format

    <Ordner rel. zu Datei> [Anzeigezeit in Sekunden*]
    <Ordner absolut>       [Anzeigezeit in Sekunden*]
    <Datei rel. zu Datein> [Anzeigezeit in Sekunden]
    <Datei absolut>        [Anzeigezeit in Sekunden]

\* Die Anzeigezeit gilt für jedes Bild in dem Ordner und allen Unterordnern. Wird ein Unterordner zusätzlich eingefügt, werden die Bilder doppelt angezeigt.

### Verhalten
Mit den Standarteinstellungen liest ```pidia``` ...

- ... alle Ordner die in der ```.pidiarc``` Datei angegeben sind.
    - (Dabei werden Unterordner rekursiv gelesen)
- ... alle ```.pidiarc``` Dateien die von der ersten/standart ```.pidiarc``` Datei entweder direkt als Datei und indirekt durch den Ordner angegeben werden.
- ... alle Dateien die in einer der ```.pidiarc``` Dateien angegeben werden.
- ... niemals eine ```.pidiarc``` Datei mehrmals


## Serverbetreuung
### Samba Dateisystem
Das Samba Dateisytem wird über ```mount``` eingehängt.

```bash
mount -t cifs -o username=$USERNAME //$TARGETIP/pool $HOME/Fass/
```

### Compilieren

```bash
go get -u github.com/auvii/pidia/...
pidia -c="$HOME/Fass/Optik/PiDiaShow/.pidiarc" -p 8080
```

### Environment
```pidia``` greift auf die Environmentvariablen ```$HOME``` und ```$GOPATH``` zu.

0. Die Templates müssen unter ```$GOPATH```/src/github.com/auvii/pidia/diaweb/ zu finden sein
0. Die Standart Konfigurationsdatei liegt unter ```$HOME/.pidiarc```
