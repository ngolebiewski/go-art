# Go ART

*A React/Typescript + Go/MySql project to privately archive artworks in a family setting (i.e. parents wanting to save their kids art, and have it searchable, and shareable with a link). Kind of like a museum collections database meets Google Photos. Some of the experiment is pushing the limit of saving image data as a BLOB in the DB and finding that flux point when the database becomes overwhelmed, and fitting something within that space as an alternative to bucket storage (which is arguably a much better solution, but we all know that already).*


## About

Why MySQL?: I have an 'unlimited' instance of MySQL databases from a legacy host that I still have, so although I would prefer PostgreSQL, and see this migrating there, there is no additional cost and no cold-starts with this option. Perfect for experimenting and building an MVP -- minimal viable product.

Why Argon2: Again, MVP, small base, so ok with the extra cost of this extra secure cryptographic method (v. bcrypt). Also, its fun resarching hashing/cryptography. To be used for auth/logins.

## local startup
1. `nvm use 20` 
    to use the version of Node that Vite needs now. Note, most of my projects are with Node 18 so keeping that around.


### DB SCHEMA SKETCH
See database/schema.sql for the schema.

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


