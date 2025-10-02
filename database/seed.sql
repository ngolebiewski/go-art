-- curl -X POST http://localhost:8080/api/auth/register \
--   -H "Content-Type: application/json" \
--   -d '{
--     "fname": "Nick",
--     "lname": "Golebiewski",
--     "email": "email@email.email",
--     "password": "super-secret-pwd" 
-- }'


INSERT INTO artists (name, codename) VALUES ('Test Artist', 'Yaag');
INSERT INTO user_artists (user_id,artist_id) VALUES (2,1);
INSERT INTO artworks (artist_id, title, grade, school, description) VALUES (1, 'Bowling', 'Post School', 'N/A','Commision of Voelkers Bowling Alley in Buffalo, NY');