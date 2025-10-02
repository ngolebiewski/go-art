-- ------------------------
-- Table: users
-- ------------------------
CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    fname VARCHAR(30) NOT NULL,
    lname VARCHAR(60) NOT NULL,
    email VARCHAR(60) NOT NULL UNIQUE,
    pwd CHAR(60) NOT NULL -- store Argon2 or bcrypt hash
);

-- ------------------------
-- Table: artists
-- ------------------------
CREATE TABLE artists (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(60) NOT NULL, -- in the api, make sure to limit this field to the user 'admin' associated with it. 
    codename VARCHAR(60); -- optional to store alias for artist name, i.e. political artist, grafitti artist, student, for privacy
);

-- ------------------------
-- Many-to-many relationship: user <-> artist
-- ------------------------
CREATE TABLE user_artists (
    user_id INT NOT NULL,
    artist_id INT NOT NULL,
    PRIMARY KEY(user_id, artist_id),
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY(artist_id) REFERENCES artists(id) ON DELETE CASCADE
);

-- ------------------------
-- Table: artworks
-- ------------------------
CREATE TABLE artworks (
    id INT AUTO_INCREMENT PRIMARY KEY,
    artist_id INT NOT NULL,               -- Artwork belongs to an artist
    grade VARCHAR(20),                    -- optional
    school VARCHAR(30),                   -- optional
    title VARCHAR(100),                    -- optional
    description VARCHAR(500),             -- optional
    FOREIGN KEY(artist_id) REFERENCES artists(id) ON DELETE CASCADE
);

-- ------------------------
-- Table: images
-- ------------------------
CREATE TABLE images (
    id INT AUTO_INCREMENT PRIMARY KEY,
    artwork_id INT NOT NULL UNIQUE,       -- 1:1 with artwork
    url VARCHAR(255),                     -- optional link
    thumb BLOB,                           -- thumbnail image <64KB
    image MEDIUMBLOB NOT NULL,            -- full image
    FOREIGN KEY(artwork_id) REFERENCES artworks(id) ON DELETE CASCADE
);

-- ------------------------
-- Table: mediums
-- ------------------------
CREATE TABLE mediums (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(60) NOT NULL UNIQUE
);

-- ------------------------
-- Many-to-many relationship: artworks <-> mediums
-- ------------------------
CREATE TABLE artworks_mediums (
    artwork_id INT NOT NULL,
    medium_id INT NOT NULL,
    PRIMARY KEY(artwork_id, medium_id),
    FOREIGN KEY(artwork_id) REFERENCES artworks(id) ON DELETE CASCADE,
    FOREIGN KEY(medium_id) REFERENCES mediums(id) ON DELETE CASCADE
);

DROP VIEW IF EXISTS all_artwork_data;

-- ------------------------
-- Collates together all the artist/artwork/medium data + a thumbnail and a link. May still b eto BIG with BLOBS
-- Note: Coalesce is cool, If codename (artist alias) is set, it shows up in the view, otherwise if null the name will show. 
-- ------------------------

If codename is NULL, the real name is shown.
CREATE VIEW all_artwork_data AS
SELECT 
    a.id AS artwork_id,
    a.grade,
    a.school,
    a.description,
    COALESCE(ar.codename, ar.name) AS artist_name,
    i.url,
    i.thumb, -- include thumbnail but not full image
    GROUP_CONCAT(m.name ORDER BY m.name SEPARATOR ', ') AS mediums
FROM artworks a
JOIN artists ar ON a.artist_id = ar.id
LEFT JOIN images i ON a.id = i.artwork_id
LEFT JOIN artworks_mediums am ON a.id = am.artwork_id
LEFT JOIN mediums m ON am.medium_id = m.id
GROUP BY a.id, a.grade, a.school, a.description, ar.name, ar.codename, i.url, i.thumb;


-- For fast artist -> artworks lookups
CREATE INDEX idx_artworks_artist_id ON artworks(artist_id);

-- For images lookups by artwork
CREATE UNIQUE INDEX idx_images_artwork_id ON images(artwork_id);

-- For join table lookups
CREATE INDEX idx_user_artists_user_id ON user_artists(user_id);
CREATE INDEX idx_user_artists_artist_id ON user_artists(artist_id);

CREATE INDEX idx_artworks_mediums_artwork_id ON artworks_mediums(artwork_id);
CREATE INDEX idx_artworks_mediums_medium_id ON artworks_mediums(medium_id);

