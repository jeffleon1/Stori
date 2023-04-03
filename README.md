<p align="center">
    <img src="https://stori-email-images.s3.amazonaws.com/StoriLogo.jpg" width="600px" height="150">
</p>


### Architecture
<p align="center">
    <img src="https://stori-email-images.s3.amazonaws.com/readme/diagrama_de_componentes.png" width="600px" height="250">
</p>


### Run API

Before you run the API need a .env file, check the folder of each one project email-ms and transaction-ms check the .env-sample for the variables that you need and create yours or change the name for use it

* Go to email-ms 
```shell
    cd ./email-ms
```

* Go to transaction-ms 
```shell
    cd ./transaction-ms
```

ones you are done with your env files you can run the app

* First of all go to the project folder
if you are in email-ms or transaction-ms use the command below
```shell
    cd ../project
```
* If you are in the root folder of the project type the command below

```shell
    cd ./project
```

* then run the app with the make file command
```shell
    make up_all
```
__Note:__
###### please check first if your docker is running you need it for run the app

---
### Testing the app
ones your application is app and running please go to documentation [http://localhost:8080/api/stori/v1/public/docs/index.html](http://localhost:8080/api/stori/v1/public/docs/index.html)

* Then you can see something like this
<p align="center">
    <img src="https://stori-email-images.s3.amazonaws.com/readme/swaggerinit.png" height="300">
</p>

* Then click in account balance proccesor
<p align="center">
    <img src="https://stori-email-images.s3.amazonaws.com/readme/swagger2.png" height="300">
</p>

* Click in try it out then choose a csv file and the email that you want, in the root folder you can see a example.csv this file contain the format that the app need for run in good way 
<p align="center">
    <img src="https://stori-email-images.s3.amazonaws.com/readme/swaggertryout.png" height="300">
</p>

Then Check the account resume was create in mongo db via mongo Express [http://localhost:8081](http://localhost:8081)

* Enter the credentials that you use in DB_PASSWORD and DB_USER in .env file

<p align="center">
    <img src="https://stori-email-images.s3.amazonaws.com/readme/mongoexpress.png" width="600px" height="250">
</p>

* Click in green button View in row stori
<p align="center">
    <img src="https://stori-email-images.s3.amazonaws.com/readme/mongoexpressstori.png" width="600px" height="250">
</p>

* Click in green button View
<p align="center">
    <img src="https://stori-email-images.s3.amazonaws.com/readme/mongoexpressview.png" width="600px" height="250">
</p>

* Finally find your account resume
<p align="center">
    <img src="https://stori-email-images.s3.amazonaws.com/readme/mongoexpressviewcollection.png" width="600px" height="250">
</p>

Then check in mailhog the mail that you send [http://localhost:8025](http://localhost:8025)

* Then you can see something like this
<p align="center">
    <img src="https://stori-email-images.s3.amazonaws.com/readme/mailhoginit.png" height="300">
</p>

* The box highlighted in blue is the inbox in this part you can see the entrance mails, the first one that you see is the last one that you send

<p align="center">
    <img src="https://stori-email-images.s3.amazonaws.com/readme/mailhoginbox.png" height="300">
</p>

* Last but not least see the email click in accounts.balance@stori.com

<p align="center">
    <img src="https://stori-email-images.s3.amazonaws.com/readme/email.png" height="300">
</p>

----
### MongoExpress

* check [Mongo Express](http://localhost:8081)

__Note:__
###### the app need to be up and running
----
### Email

* check [Mail Hog](http://localhost:8025)

__Note:__
###### the app need to be up and running

----
### Docs
* See the documentation ones the API is up in the following URL

[Swagger documentation](http://localhost:8080/api/stori/v1/public/docs/index.html)

__Note:__
###### the app need to be up and running