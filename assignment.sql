CREATE database assignment_db;

USE assignment_db;

CREATE TABLE Passenger (
 PassengerID int NOT NULL AUTO_INCREMENT PRIMARY KEY,
 FirstName VARCHAR(30), 
 LastName VARCHAR(30), 
 MobileNumber varchar(8), 
 EmailAddress varchar(60)
 );
 
CREATE TABLE Driver (
DriverID varchar(50) NOT NULL UNIQUE PRIMARY KEY, 
FirstName VARCHAR(30), 
LastName VARCHAR(30), 
MobileNumber varchar(8), 
EmailAddress varchar(50), 
LicenseNumber varchar(7)
); 

CREATE TABLE Trips(
  TripID int NOT NULL AUTO_INCREMENT,
  DriverID varchar(50) NOT NULL,
  PassengerID int NOT NULL,
  PickUpLocation varchar(255) NOT NULL,
  DropOffLocation varchar(255) NOT NULL,
  PickUpTime varchar(255) NOT NULL,
  DropOffTime varchar(255) NOT NULL,
  Status varchar(255) NOT NULL CHECK (Status IN ('Pending', 'On The Way', 'In Transit', 'Completed', 'Failed')),
  PRIMARY KEY (TripID) 
);

DROP TABLE IF EXISTS Passenger;
DROP TABLE IF EXISTS Driver;
DROP TABLE IF EXISTS Trips;


INSERT INTO Driver (DriverID, FirstName, LastName, MobileNumber, EmailAddress, IdNumber, CarLicenseNumber) VALUES ("0001", "John", "Jones", "93123412", "driver@gmail.com", 1, "SBH3158E");

INSERT INTO Passenger (FirstName, LastName, MobileNumber, EmailAddress)
VALUES("Jake","Lee","93123412","passenger@gmail.com");
INSERT INTO Driver (DriverID, FirstName, LastName, MobileNumber, EmailAddress, LicenseNumber) 
VALUES("S1234567D", "John","Jones","93123412","driver@gmail.com","S126LO");
INSERT INTO Trips (DriverID, PassengerID, PickupLocation, DropoffLocation, PickUpTime, DropOffTime, Status)
VALUES("S1234567D",3,"123456", "123457", "12:30", "13:00", "Completed");
INSERT INTO Trips (DriverID, PassengerID, PickupLocation, DropoffLocation, PickUpTime, DropOffTime, Status)
VALUES("S1234567D",3,"123456", "123457", "12:30", " ", "Pending");

