# ETI_Assginment1
<h2>Architecture diagram</h2>

![Untitled Diagram (1)](https://user-images.githubusercontent.com/78250532/145850298-9ac8a5c2-ac57-485b-99e8-8a7e83db817a.jpg)

<h2>3 microservices</h2>

<h2>Driver</h2>

<ul>
  <li>create driver</li>
  <li>select driver</li>
  <li>select driver</li>
</ul>

<h2>Passenger</h2>

<ul>
  <li>create passenger</li>
  <li>select passenger</li>
  <li>select passenger</li>
</ul>

<h2>Trips</h2>

<ul>
  <li>create trip</li>
  <li>select trip</li>
  <li>select trip</li>
</ul>
  
<p>attempted to use monolith front end, but was'nt able to do much with it. Each microservice had its onw database table and was communicates with each other to retrieve trips and was planning on connecting it to the frontend<p>
 
<h2>Persistent storage of information using database with mySQL</h2>
<p>3 tables for each microservice<p>

<h2>Prerequisites</h2>
<p>Please ensure that GOLANG and MYSQL is installed on your system, and is fully operational</p>

<p>Please do also ensure that your SQL user login is as such:</p>

   <p>Username: root</p>
   <p>Password: 12N28c02</p>
   
<h2><b>Installation</b></h2>
<h3>Clone the repo</h3>
<p>git clone https://github.com/JSQJ/ETI_Assginment1.git</p>

<h3>Install necessary libraries</h3>
<ul>
  <li>go get -u github.com/go-sql-driver/mysql</li>
  <li>go get -u github.com/gorilla/mux</li>
  <li>go get -u github.com/gorilla/handlers</li>
</ul>
<p>Execute SQL script in mySQL database</p>
