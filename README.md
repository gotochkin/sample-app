# sample-app
Sample web app written on Go and working with Oracle and Postgresql backends
The Readme is still in progress. Please have a look the https://blog.gleb.ca/2022/09/23/sample-go-application-with-database-backend/ to get the basic instructions until the Readme is fully finished. 

# Description (TODO)

# Usage
## Preparations 
### Database backend
You need to create a database and provide a network connection from the application side. Depending on a database flavour the path can be different. Here are the basic steps to create an Oracle or Postgres backends.
1. Oracle Cloud Autonomous database backend. 
    *  Create a database. For Oracle Cloud Autonomous database register or login to your Oracle Cloud console,click on the “hamburger” (top left corner) and choose “Oracle Database” -> “Autonomous Database”. Then push the button   “Create Autonomous Database” and fill up the required fields. Pay attention to the “Choose network access” section. The simplest way is to choose access from everywhere but if you choose private endpoints or allowed IPs then make sure it is going to be accessible from your application host. 
    * Create the user. Open SqlDeveloper or SQL*Plus application, connect as admin and create a user which will be used for connection. Write down the password for the user. If you are not familiar with Oracle you can start from using the user admin and password used during creation of the database but keep in mind it is not the best practice.
    * Create the table and fill it with the sample data. You can use the statements from hr_cre.sql from the application’s “sql” directory. 
    * In the cloud console go to the database page, open “Database Connection” select “Instance Wallet” from the dropdown menu and download the wallet file.
2. PostreSQL backend. (To be added)

### Prepare the application host.
1.  Install the necessary packages to the application host. I will be showing using Oracle 8 on a VM as an example.
    ```
    $ sudo dnf update -y
    $ sudo dnf install golang git -y
    ```
2. Clone the application from GitHub.
    ```
    $ git clone https://github.com/gotochkin/sample-app.git
    $ cd sample-app
    ```
3. Create a directory and unzip the wallet downloaded from the Oracle cloud to the directory
    ```
    $ mkdir ssl
    $ unzip ~/Wallet_myatpdb.zip -d ssl/
    ```
4. Create a file with environment variables and execute it or export the variables using command line one by one.
    ```
    $ vi ~/sampleapp.env
    $ cat ~/sampleapp.env

    ###############The file with env variables################
    DBNAME=m5c5hcat3eqqydh_myatpdb_tp.adb.oraclecloud.com
    DBWALLET=~/sample-app/ssl
    DBHOST=adb.us-ashburn-1.oraclecloud.com
    DBPORT=1522
    DBPASS=MyExtremelyDifficultToRememberPassword
    DBUSER=sampleapp
    DBVERSION=ORACLE
    export DBNAME DBHOST DBPORT DBPASS DBUSER DBWALLET DBVERSION
    #####################################################

    $ source ~/sampleapp.env
    ```
4. Open port 8080 for web access to the application.
    ```
    $ sudo firewall-cmd --permanent --zone=public --add-port=8080/tcp
    $ sudo firewall-cmd --reload
    ```
## Run the application
### Run interactively in console
1. Change directory to the application’s source code folder import the variables and run the application. 
    ```
    $ cd ~/sample-app
    $ source ~/sampleapp.env
    $ go run cmd/app/main.go 

    #The expected output is 

    2022/09/23 15:06:47 Listening on port :8080
    ```
Try to access the application using http://<the host ip>:8080 You should see the last 10 rows from the table and a form to add a new employee information.

### Build and run it
1. Change directory to the application’s source code folder and build the application
    ```
    $ cd ~/sample-app
    $ go build -v -o sampleapp ./cmd/app/
    ```
2. Run the application
    ```
    $ ./sampleapp 

    #The expected output is 
    2022/09/23 15:06:47 Listening on port :8080
    ```

### Run from a docker container
Assuming you already have the docker installed. 
1.  Copy the wallet file to the ssl directory in the root folder of the application. 
    ```
    $ cp ~/Wallet_myatpdb.zip ssl/wallet.zip
    ```
2. Build the image
    ```
    $ cd ~/sample-app
    $ docker build -t otochkin/sampleapp:0.01 .
    ```
3. Prepare the  ~/sampleapp.env with variables 
    ```
    ###############The file with env variables################
    DBNAME=m5c5hcat3eqqydh_myatpdb_tp.adb.oraclecloud.com
    DBWALLET=~/sample-app/ssl
    DBHOST=adb.us-ashburn-1.oraclecloud.com
    DBPORT=1522
    DBPASS=MyExtremelyDifficultToRememberPassword
    DBUSER=sampleapp
    DBVERSION=ORACLE
    #####################################################

    ```
3. Create the container
    ```
    $ docker run --name sampleapp -p 8080:8080 --env-file ~/sampleapp.env sampleapp:0.01
    ```
    #The expected output :
    2022/09/22 19:22:18 Listening on port :8080

4. Access the application using url http://localhost:8080
