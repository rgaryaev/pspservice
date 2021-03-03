# pspservice
This stand alone service can be used as an online check tool for russian internal passports. 
Verification is provided by using the official database of expired and non-valid passports (from Russian Ministry of Internal Affairs)

The service is developed as an academic project with purpose to learn Golang.

For more information
[http://сервисы.гувм.мвд.рф/info-service.htm?sid=2000](http://xn--b1afk4ade4e.xn--b1ab2a0a.xn--b1aew.xn--p1ai/info-service.htm?sid=2000).


## License
The MIT License (MIT). Please see [`LICENSE`](./LICENSE) for more information.

### Usage in Docker

```shell
git clone https://github.com/rgaryaev/pspservice.git
cd ./pspservice
docker build . -t pspservice
docker run -d -p 8080:8080 pspservice
```

### Usage from CLI
You have to have installed Golang development evironment  

```shell
git clone https://github.com/rgaryaev/pspservice.git
cd ./pspservice
mkdir .data
go build 
./pspservice
```
#### or on Windows you have to run
pspservice.exe 

## Example of using
You can see examples (scripts) in the testsrc/testweb
'powershell' folder contains script for Power Shell with examples how to use the service
'bash' folders contains scripts for curl and ab (ApacheBench ) with examples how to use and make a bechmark, accordingly 

##  Params in configuration file (config.json)
Default config has the following parameters
```
{
 	"storage": {
 		"passport_data": "./.data/list_of_expired_passports.csv",
 		"engine": "roaring_bitmap"
 	},
 	"listener": {
 		"address": "0.0.0.0",
 		"port": "8080",
 		"max_passport_per_request": 100
 	},
 	"loader": {
 		"source_url": "http://guvm.mvd.ru/upload/expired-passports/list_of_expired_passports.csv.bz2",
 		"every_x_day": 1,
 		"last_update": "2021-03-01"
 	}
}
```
- "passport_data" - path to file with passport data. When service is starting first time this file doesn't extst usually and will be downloaded automatically.

- "engine" - this parameter defines storage engine. Possible values are:  "roaring_bitmap"  or "sparse_bitmap".  
           **"roaring_bitmap"** option uses compressed bitmap and is most effiecient from point of view of memory consumption.
           So far the current passport data requires about 42 - 44 Mb in memory. 
           For more information about roaring bitmap you can vizit http://roaringbitmap.org/. 
           **"sparse_bitmap"** is a simple bitmap where passport series are rows and passport numbers are bitmap colums. This engine requires about 1.25 Gb in memory, 
           as we have 9999 rows * (999999 numbers / 64 bit) * 8 is about 1.25 Gb. This engine is expected to be faster than "roraring_bitmap" but actually there is no 
           big difference between them, so "sparse_bitmap" is not recomended to use.

- "address" and "port" - parameters for http listener 

- "max_passport_per_request" - default number of passports (series and number) in the Body of http request.  in other words it is 
                             size of JSON array (see testsrv/testweb examples for jason request). 
                             The array from the request will truncated if it contains number of elements more than value of this parameter.
- "source_url"  - full url path where original passport data file is stored

- "every_x_day" - this parameter defines regularity how often to do update (if need).For instance, 1 means every day,  2 one time per 2 days and etc.  
                Anyway before updating the service is checking Last Modification of file. if it was modified since last update then the service starts downloading and updating                   passport data. 

- "last_update" - last date of completed update of passport data. This parameter is overwriting by service
