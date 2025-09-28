## local startup, since not using Docker, even though I should
1. `nvm use 20` 
    to use the version of Node that Vite needs now. Note, most of my projects are with Node 18 so keeping that around.
2. ` brew services start mysql`


DB SCHEMA SKETCH

TABLE USER
id unique autoincrement
fname varchar 30
lname varchar 60
email varchar unique 60
pwd varchar +store using argon2 hash or go built in cryptography

TABLE ARTIST
id unique autoincrement
name fname

Relationshop
many user to many 'artist' i.e. user is parent and artist is child, we are building up their art

TABLE Artwork
id unique autoincrement
grade? varchar 20 // could be kindergaten or college senior or 2
School? varchar 20 // could be PS512
Image (1 to 1 to image relationship)
Description? varchar 500

TABLE IMAGE
id unique autoincrement
Artwork ID foreign key
url? varhar //if there is a link put that here
thumb smallBLOB
image MED BLOB


