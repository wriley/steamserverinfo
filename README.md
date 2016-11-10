# steamserverinfo

A server query program written in [golang](https://golang.org/) that uses the [Steam Server Query API](https://developer.valvesoftware.com/wiki/Server_queries)

## Example Usage

### DayZ server
```
$ ./steamserverinfo 184.172.24.19 2512 2>/dev/null | grep PLAYERS:
PLAYERS: 50
MAXPLAYERS: 50
```

### Miscreated server
```
$ ./steamserverinfo 91.198.152.238 64092 2>/dev/null | grep PLAYERS:
PLAYERS: 22
MAXPLAYERS: 36
```

### Debug output from a DayZ server
```
# Note:
# any 3rd argument enables verbose output
# if verbose then any 4th argument enables color for hex dumps
#
$ steamserverinfo 173.192.105.212 2512 v c
Opening UDP connection...
Sending A2S_INFO...
Writing...
00000000  ff ff ff ff 54 53 6f 75  72 63 65 20 45 6e 67 69  |....TSource Engi|
00000010  6e 65 20 51 75 65 72 79  00                       |ne Query.|
Wrote 25 bytes
Reading...
Read 127 bytes
00000000  ff ff ff ff 49 11 2f 72  2f 44 61 79 7a 55 6e 64  |....I./r/DayzUnd|
00000010  65 72 67 72 6f 75 6e 64  20 55 53 00 32 38 32 33  |erground US.2823|
00000020  34 31 00 64 61 79 7a 00  44 61 79 5a 00 00 00 09  |41.dayz.DayZ....|
00000030  32 00 64 77 00 01 30 2e  36 30 2e 31 33 33 39 31  |2.dw..0.60.13391|
00000040  33 00 b1 c6 09 09 f0 2b  11 2b 1e 40 01 62 61 74  |3......+.+.@.bat|
00000050  74 6c 65 79 65 2c 6e 6f  33 72 64 2c 70 72 69 76  |tleye,no3rd,priv|
00000060  48 69 76 65 2c 65 74 6d  34 2e 30 30 30 30 30 30  |Hive,etm4.000000|
00000070  2c 30 32 3a 35 35 00 ac  5f 03 00 00 00 00 00     |,02:55.._......|
HEADER: 0x49
PROTOCOL: 0x11
NAME: /r/DayzUnderground US
MAP: 282341
FOLDER: dayz
GAME: DayZ
ID: 0
PLAYERS: 9
MAXPLAYERS: 50
BOTS: 0
SERVERTYPE: d
ENVIRONMENT: w
VISIBILITY: 0
VAC: 1
VERSION: 0.60.133913
PORT: 2502
KEYWORDS: @battleye,no3rd,privHive,etm4.000000,02:55
Sending A2S_RULES...
Writing...
00000000  ff ff ff ff 56 ff ff ff  ff                       |....V....|
Wrote 9 bytes
Reading...
Read 9 bytes
00000000  ff ff ff ff 41 92 7d a9  06                       |....A.}..|
Challenge number: 111771026
Sending A2S_RULES...
Writing...
00000000  ff ff ff ff 56 92 7d a9  06                       |....V.}..|
Wrote 9 bytes
Reading...
Read 140 bytes
00000000  ff ff ff ff 45 08 00 61  6c 6c 6f 77 65 64 42 75  |....E..allowedBu|
00000010  69 6c 64 00 31 33 33 39  31 33 00 64 65 64 69 63  |ild.133913.dedic|
00000020  61 74 65 64 00 31 00 69  73 6c 61 6e 64 00 43 68  |ated.1.island.Ch|
00000030  65 72 6e 61 72 75 73 50  6c 75 73 00 6c 61 6e 67  |ernarusPlus.lang|
00000040  75 61 67 65 00 36 35 35  34 35 00 70 6c 61 74 66  |uage.65545.platf|
00000050  6f 72 6d 00 77 69 6e 00  72 65 71 75 69 72 65 64  |orm.win.required|
00000060  42 75 69 6c 64 00 31 33  33 39 31 33 00 72 65 71  |Build.133913.req|
00000070  75 69 72 65 64 56 65 72  73 69 6f 6e 00 36 30 00  |uiredVersion.60.|
00000080  74 69 6d 65 4c 65 66 74  00 31 35 00              |timeLeft.15.|
RULE LIST:
allowedBuild 133913
dedicated 1
island ChernarusPlus
language 65545
platform win
requiredBuild 133913
requiredVersion 60
timeLeft 15
Sending A2S_PLAYER...
Writing...
00000000  ff ff ff ff 55 ff ff ff  ff                       |....U....|
Wrote 9 bytes
Reading...
Read 9 bytes
00000000  ff ff ff ff 41 92 7d a9  06                       |....A.}..|
Challenge number: 111771026
Sending A2S_PLAYER...
Writing...
00000000  ff ff ff ff 55 92 7d a9  06                       |....U.}..|
Wrote 9 bytes
Reading...
Read 162 bytes
00000000  ff ff ff ff 44 09 00 45  6c 65 70 68 61 6e 74 00  |....D..Elephant.|
00000010  00 00 00 00 1a d7 b6 45  00 53 70 61 6e 6b 79 00  |.......E.Spanky.|
00000020  00 00 00 00 4b 10 a8 45  00 53 68 61 6f 6c 69 6e  |....K..E.Shaolin|
00000030  00 00 00 00 00 79 7a 7e  45 00 4c 61 63 65 73 54  |.....yz~E.LacesT|
00000040  6f 6f 4c 6f 6e 67 00 00  00 00 00 c1 72 6e 45 00  |ooLong......rnE.|
00000050  4e 61 67 61 76 00 00 00  00 00 97 43 44 45 00 67  |Nagav......CDE.g|
00000060  65 6e 61 72 00 00 00 00  00 c3 a0 60 44 00 53 70  |enar.......`D.Sp|
00000070  65 63 74 72 65 00 00 00  00 00 8b 0f 53 44 00 41  |ectre.......SD.A|
00000080  78 65 6c 20 41 67 75 73  74 69 6e 00 00 00 00 00  |xel Agustin.....|
00000090  e4 3a 75 43 00 42 75 73  68 00 00 00 00 00 c3 33  |.:uC.Bush......3|
000000a0  80 41                                             |.A|
PLAYER LIST:
Elephant 0 5851
Spanky 0 5378
Shaolin 0 4072
LacesTooLong 0 3815
Nagav 0 3140
genar 0 899
Spectre 0 844
Axel Agustin 0 245
Bush 0 16
```
