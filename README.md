# steamserverinfo

A server query program written in [golang](https://golang.org/) that uses the [Steam Server Query API](https://developer.valvesoftware.com/wiki/Server_queries)

#### Example Usage

DayZ server
```
$ ./steamserverinfo 184.172.24.19 2512 2>/dev/null | grep PLAYERS:
PLAYERS: 50
MAXPLAYERS: 50
```

Miscreated server
```
$ ./steamserverinfo 91.198.152.238 64092 2>/dev/null | grep PLAYERS:
PLAYERS: 22
MAXPLAYERS: 36
```

Full output from a DayZ server
```
$ ./steamserverinfo 184.172.24.19 2512
Opening UDP connection...
Sending A2S_INFO...
Writing...
Wrote 25 bytes
Reading...
HEADER: 0x49
PROTOCOL: 0x11
NAME: /r/DayzUnderground
MAP: 282341
FOLDER: dayz
GAME: DayZ
ID: 0
PLAYERS: 50
MAXPLAYERS: 50
BOTS: 0
SERVERTYPE: d
ENVIRONMENT: w
VISIBILITY: 0
VAC: 1
VERSION: 0.60.133617
PORT: 2502
KEYWORDS: battleye,no3rd,privHive,etm4.000000,12:47
Sending A2S_RULES...
Writing...
Wrote 9 bytes
Reading...
Challenge number: 215285510
Sending A2S_RULES...
Writing...
Wrote 9 bytes
Reading...
RULE LIST:
allowedBuild 133617
dedicated 1
island ChernarusPlus
language 65545
platform win
requiredBuild 133617
requiredVersion 60
timeLeft 15
Sending A2S_PLAYER...
Writing...
Wrote 9 bytes
Reading...
Challenge number: 215285510
Sending A2S_PLAYER...
Writing...
Wrote 9 bytes
Reading...
PLAYER LIST:
WARDEN WILL 0 800
Oposum 0 800
Bandit Slayer 0 800
Hayden 0 800
AoF Mudd 0 794
Milo Windbane 0 794
[FMC] Altra 0 793
Gonzalo Darin 0 792
Soap Colton 0 789
Cletus Colton 0 789
BlinG* 0 785
Vlad Jiggers 0 784
Frank Pirnie 0 784
A Gwopstop 0 782
DrDirtyDan 0 779
CPRL. Muzzak 0 779
BlackBeard 0 779
Courtney Colton 0 777
[A3F] Shadow 0 773
John 0 772
Sobieski12 0 769
boykie 0 769
Ryan 0 768
Kleanuppguy 0 765
CuriousCole 0 763
[V/\] Jake 0 762
Ryan (2) 0 741
Mills 0 740
[V/\]Mr.bear 0 731
Algalica 0 728
[A3F] Stiles 0 728
Bravo-Skie 0 718
Tommie Dantonio 0 691
Redmoon 0 690
b0bth3k1ll3r 0 663
Mr.Nobody 0 662
M12 0 642
Sgt. Ganja 0 587
Jeff Miller 0 563
Rueby 0 552
Bestbuds 0 542
Wyatt Mann 0 538
Sally Carfano 0 484
Isaac\ 0 462
Jarrett 0 414
Jari 0 284
quickfinger 0 274
Ducktape 0 126
sniper DEW 0 126
 0 24
```
