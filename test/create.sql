CREATE TABLE user ( 
    name VarChar(255)  PRIMARY KEY ,
    id INTEGER ,
    registered BOOLEAN , 
    acount Integer  FOREIGN KEY REFERENCES Account.id 
);