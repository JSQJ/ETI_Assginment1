# ETI_Assginment1
<h2>Architecture diagram</h2>

![Untitled Diagram (1)](https://user-images.githubusercontent.com/78250532/145850298-9ac8a5c2-ac57-485b-99e8-8a7e83db817a.jpg)

<h2>3 microservices</h2>
-Driver
  -create driver
  -select driver
  -update driver
 -Passenger
  -create driver
  -select driver
  -update driver
-Trips
  -select trip
  -create trip
  -update trip
  
attempted to use monolith front end, but was'nt able to do much with it
 
Persistent storage of information using database with mySQL
-3 tables for each microservice

<h2>Prerequisites</h2>
Please ensure that GOLANG and MYSQL is installed on your system, and is fully operational

Please do also ensure that your SQL user login is as such:

   Username: root
   Password: 12N28c02
   
<h2><b>Installation</b></h2>
Clone the repo
git clone https://github.com/JSQJ/ETI_Assginment1.git

Install necessary libraries
<a>go get -u github.com/go-sql-driver/mysql</a>
<a>go get -u github.com/gorilla/mux</a>
<a>go get -u github.com/gorilla/handlers</a>

Execute SQL script in mySQL database
