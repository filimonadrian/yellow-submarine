# Tema 2: Yellow Submarine

Student: Filimon Adrian  
Grupa: 342C2

## Implementare

### Serverul de REST API si Interfata grafica

- scrise in limbajul Go
- Cele doua aplicatii sunt legate prin intermediul unei conexiuni pe socketi:
    - aplicatia `server` porneste prima si initializeaza serverul de http si asteapta conexiuni pe socketul deschis(portul 8000)
    - cand aplicatia `gui` porneste, aceasta se conecteaza la server
    - din acest moment, comenzile primite de serverul de http vor fi trimise mai departe, la interfata grafica
        - Serverul de http primeste pe rute diferite comenzile pentru submarin, pesti sau artefacte
        - Pentru a usura munca facuta de gui si transmisia datelor, structura primita de server este reimpachetata intr-o alta structura care contine si numele resursei. Ex:
            ```
            {
                Type: submarine,
                x:10,
                y:15
            }
            ```
    - Un aspect care merita mentionat aici este ca daca se intrerupe conexiunea intre cele 2 si apoi isi revine, pot comunica in continuare(exista un thread separat pentru pornirea si oprirea conexiunilor)

#### Interfata grafica

- foloseste modulul [tcell](github.com/gdamore/tcell)
- Pentru redirectarea datelor la consola `/dev/ttyAMA0`, am folosit o functie oferita nativ de acest modul: `NewDevTtyFromDev`
- Redirectarea iesirilor si intrarilor standard la consola nu a mers(nu am gasit problema)
- O problema cu utilizarea acestei functii este ca desi pe masina locala functioneaza, cand rulez in qemu considera ca am o consola de doar 80x40(marimea standard pentru terminal) si nu am putut modifica acest lucru

### Imaginea linux de baza

- Am creat layer-ul meta-tema-si, care contine mai multe retete:
    - recipes-core, ce contine imaginea de baza
    - recipes-gui, ce contine daemon-ul de display
    - recipes-server, ce contine serverul de http

1. Recipes-core
    - Am facut append la imaginea de baza, unde am mai adaugat cateva pachete: `dhcp-client`, `go`, `avahi-daemon`, `ssh` si cele 2 retete custom, `server` si `gui`
    - De asemenea, tot aici am adaugat si parametrii pentru creearea de user
    - Pentru schimbarea hostname-ului, am facut append la reteta `base-files`

2. Recipes-gui si recipes-server
    - Acestea sunt retetele create de mine pentru server si interfata grafica
    - Desi sunt 99% asemanatoare, am preferat sa fac 2 retete diferite din motive de design(si sa le pot builda separat)
    - Fisierele .bb contin functia de compilare si functia de install, unde se incarca in sistem(in /usr/bin)
    - Pentru a realiza rularea task-urilor la boot, am creat cate un serviciu pentru fiecare aplicatie
        - Scriptul de init(`gui-init`/`server-init`) este adaugat in /etc/init.d(unde se afla toate serverviciile) si este rulat la boot

# Linkuri utile in dezvoltarea temei

- Serverul de http si comunicarea pe socketi:
    - https://golangr.com/golang-http-server/
    - https://itnext.io/plain-socket-communication-between-two-go-programs-the-easy-way-bd5ac5819eb6

- Crearea imaginii de linux:
    - https://ocw.cs.pub.ro/courses/si/laboratoare/yocto-extra1
    - https://lynxbee.com/bitbake-yocto-recipes-for-cross-compiling-golang-program/
    - https://docs.google.com/presentation/d/1V8X7M14CuTXLThrJPnLgTML3BWMP3appLTsm1GFjX3E/edit#slide=id.g10640917eb2_0_24
    - https://wiki.yoctoproject.org/wiki/Cookbook:Appliance:Startup_Scripts
    - https://gist.github.com/naholyr/4275302
    - https://git.yoctoproject.org/poky/plain/meta/recipes-extended/go-examples/go-helloworld_0.1.bb
    - https://www.jamescoyle.net/cheat-sheets/791-update-rc-d-cheat-sheet
    - https://stackoverflow.com/questions/63723563/run-application-at-start-up-with-yocto-dunfell
    - http://docs.yoctoproject.org/ref-manual/variables.html#term-CORE_IMAGE_EXTRA_INSTALL
