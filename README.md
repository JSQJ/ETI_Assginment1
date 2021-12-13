# ETI_Assginment1
Architecture diagram

![Untitled Diagram (1)](https://user-images.githubusercontent.com/78250532/145850298-9ac8a5c2-ac57-485b-99e8-8a7e83db817a.jpg)

3 microservices
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

Prerequisites
Please ensure that GOLANG and MYSQL is installed on your system, and is fully operational

Please do also ensure that your SQL user login is as such:

   Username: root
   Password: 12N28c02
   
Installation
Clone the repo
git clone https://github.com/JSQJ/ETI_Assginment1.git

Install necessary libraries
go get -u github.com/go-sql-driver/mysql
go get -u github.com/gorilla/mux
go get -u github.com/gorilla/handlers

Execute SQL script in mySQL database
