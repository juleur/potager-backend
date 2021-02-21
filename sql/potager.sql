SET TIME ZONE 'Europe/Paris';
CREATE EXTENSION pg_trgm;

DROP TABLE IF EXISTS users;
CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  username VARCHAR(25) NOT NULL,
  email VARCHAR(55) NOT NULL,
  password TEXT NOT NULL,
  refresh_token CHAR(32) CHECK(length(refresh_token) = 32),
  registered_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  verified_at TIMESTAMP NULL,
  last_logged_in_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE UNIQUE INDEX users_email_unique_idx ON users (email) WITH (deduplicate_items = off);
CREATE UNIQUE INDEX users_username_unique_idx ON users (username) WITH (deduplicate_items = off);
CREATE UNIQUE INDEX users_refresh_token_unique_idx ON users (refresh_token) WITH (deduplicate_items = off);

DROP TABLE IF EXISTS users_verification_codes;
CREATE TABLE users_verification_codes (
  pwd_reset_code INTEGER NOT NULL,
  verify_account_code INTEGER NOT NULL,
  pwd_reset_code_at TIMESTAMP NOT NULL,
  verify_account_code_at TIMESTAMP NOT NULL,
  user_id INTEGER NOT NULL CHECK(user_id > 0),
  CONSTRAINT chk_users_verification_codes_pwd_reset_code CHECK(pwd_reset_code >= 1000 AND pwd_reset_code <= 9999),
  CONSTRAINT chk_users_verification_codes_verify_account_code CHECK(verify_account_code >= 1000 AND verify_account_code <= 9999),
  CONSTRAINT fk_users_verification_codes_user_id FOREIGN KEY(user_id) REFERENCES users(id)
);
CREATE UNIQUE INDEX users_verification_codes_user_id_unique_idx ON users_verification_codes (user_id) WITH (deduplicate_items = off);

DROP TABLE IF EXISTS users_permission;
CREATE TABLE users_permission (
  failed_login_attempts SMALLINT DEFAULT 0 NOT NULL,
  locked_until TIMESTAMP NULL,
  status SMALLINT DEFAULT 2 NOT NULL, -- 0: banni, 1: inactif(supprimer), 2: actif
  user_id INTEGER NOT NULL CHECK(user_id > 0),
  CONSTRAINT chk_users_permission_failed_login_attempts CHECK(status >= 0 AND status <= 3),
  CONSTRAINT chk_users_permission_status CHECK(status >= 0 AND status <= 2),
  CONSTRAINT fk_users_permission_user_id FOREIGN KEY(user_id) REFERENCES users(id)
);
CREATE UNIQUE INDEX users_permission_user_id_unique_idx ON users_permission (user_id) WITH (deduplicate_items = off);
CREATE INDEX users_permission_status_idx ON users_permission (status) WHERE status = 2;

DROP TABLE IF EXISTS users_farmer;
CREATE TABLE users_farmer(
  id SERIAL PRIMARY KEY,
  img_url TEXT NOT NULL,
  description TEXT NULL,
  commune VARCHAR(45) NOT NULL,
  coordonnees GEOGRAPHY(Point,4326) NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NULL,
  temporary_disabled BOOLEAN DEFAULT FALSE NOT NULL,
  user_id INTEGER NOT NULL CHECK(user_id > 0),
  CONSTRAINT fk_users_farmer_user_id FOREIGN KEY(user_id) REFERENCES users(id)
);
CREATE UNIQUE INDEX users_farmer_user_id_unique_idx ON users_farmer (user_id) WITH (deduplicate_items = off);
CREATE INDEX users_farmer_gix ON users_farmer USING GIST (coordonnees);
CREATE INDEX users_farmer_user_id_tp_idx ON users_farmer(id, temporary_disabled) WHERE temporary_disabled = FALSE;


DROP TABLE IF EXISTS favorite_potagers;
CREATE TABLE favorite_potagers (
  farmer_id INTEGER NOT NULL CHECK(farmer_id > 0),
  user_id INTEGER NOT NULL CHECK(user_id > 0),
  CONSTRAINT chk_favorite_potagers_not_himself CHECK(farmer_id <> user_id),
  CONSTRAINT fk_favorite_potagers_farmer_id FOREIGN KEY(farmer_id) REFERENCES users_farmer(id),
  CONSTRAINT fk_favorite_potagers_user_id FOREIGN KEY(user_id) REFERENCES users(id)
);
CREATE UNIQUE INDEX favorite_potagers_unique_idx ON favorite_potagers (user_id, farmer_id) WITH (deduplicate_items = off);

DROP TABLE IF EXISTS muted_potagers;
CREATE TABLE muted_potagers (
  farmer_id INTEGER NOT NULL CHECK(farmer_id > 0),
  user_id INTEGER NOT NULL CHECK(user_id > 0),
  CONSTRAINT chk_muted_potagers_not_himself CHECK(farmer_id <> user_id),
  CONSTRAINT fk_muted_potagers_farmer_id FOREIGN KEY(farmer_id) REFERENCES users_farmer(id),
  CONSTRAINT fk_muted_potagers_user_id FOREIGN KEY(user_id) REFERENCES users(id)
);
CREATE UNIQUE INDEX muted_potagers_unique_idx ON muted_potagers (user_id, farmer_id) WITH (deduplicate_items = off);

DROP TABLE IF EXISTS fruits;
CREATE TABLE fruits (
  id SERIAL PRIMARY KEY,
  img_url TEXT NULL,
  nom VARCHAR(20) NOT NULL,
  variete VARCHAR(30) NOT NULL,
  systeme_echange SMALLINT[] NOT NULL, -- Don: 0, Troc: 1, Vente: 2
  prix NUMERIC(4,2) NULL CHECK (prix > 0 AND prix < 99),
  unite_mesure SMALLINT NOT NULL, -- Botte: 0, Kg: 1, Piece: 2
  stock SMALLINT DEFAULT 2 NOT NULL, -- Épuisé: 0, Peu: 1, Disponible: 2
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NULL,
  CONSTRAINT chk_fruits_unite_mesure CHECK(unite_mesure >= 0 AND unite_mesure <= 2),
  CONSTRAINT chk_fruits_stock CHECK(stock >= 0 AND stock <= 2)
);

DROP TABLE IF EXISTS rel_fruits_farmers;
CREATE TABLE rel_fruits_farmers (
  farmer_id INTEGER NOT NULL CHECK (farmer_id > 0),
  fruit_id INTEGER NOT NULL CHECK (fruit_id > 0),
  CONSTRAINT fk_rel_fruits_users_farmer_farmer_id FOREIGN KEY(farmer_id) REFERENCES users_farmer(id),
  CONSTRAINT fk_rel_fruits_users_farmer_fruit_id FOREIGN KEY(fruit_id) REFERENCES fruits(id) ON DELETE CASCADE
);
CREATE INDEX rel_fruits_users_farmer_farmer_id_idx ON rel_fruits_farmers (farmer_id) WITH (deduplicate_items = off);
CREATE INDEX rel_fruits_users_farmer_fruit_id_idx ON rel_fruits_farmers (fruit_id) WITH (deduplicate_items = off);

DROP TABLE IF EXISTS graines;
CREATE TABLE graines (
  id SERIAL PRIMARY KEY,
  img_url TEXT NULL,
  nom VARCHAR(20) NOT NULL,
  variete VARCHAR(30) NOT NULL,
  systeme_echange SMALLINT[] NOT NULL, -- Don: 0, Troc: 1, Vente: 2
  prix NUMERIC(4,2) NULL CHECK (prix > 0 AND prix < 99),
  stock SMALLINT DEFAULT 2 NOT NULL, -- Épuisé: 0, Peu: 1, Disponible: 2
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NULL,
  CONSTRAINT chk_graines_stock CHECK(stock >= 0 AND stock <= 2)
);

DROP TABLE IF EXISTS rel_graines_farmers;
CREATE TABLE rel_graines_farmers (
  farmer_id INTEGER NOT NULL CHECK (farmer_id > 0),
  graine_id INTEGER NOT NULL CHECK (graine_id > 0),
  CONSTRAINT fk_rel_graines_users_farmer_farmer_id FOREIGN KEY(farmer_id) REFERENCES users_farmer(id),
  CONSTRAINT fk_rel_graines_users_farmer_graine_id FOREIGN KEY(graine_id) REFERENCES graines(id) ON DELETE CASCADE
);
CREATE INDEX rel_graines_users_farmer_farmer_id_idx ON rel_graines_farmers (farmer_id) WITH (deduplicate_items = off);
CREATE INDEX rel_graines_users_farmer_graine_id_idx ON rel_graines_farmers (graine_id) WITH (deduplicate_items = off);

-- le type de nom devrait être remplacer par jsonb type pour la traduction
DROP TABLE IF EXISTS legumes;
CREATE TABLE legumes (
  id SERIAL PRIMARY KEY,
  img_url TEXT NULL,
  nom VARCHAR(20) NOT NULL,
  variete VARCHAR(30) NOT NULL,
  systeme_echange SMALLINT[] NOT NULL, -- Don: 0, Troc: 1, Vente: 2
  prix NUMERIC(4,2) NULL CHECK (prix > 0 AND prix < 99),
  unite_mesure SMALLINT NOT NULL, -- Botte: 0, Kg: 1, Piece: 2
  stock SMALLINT DEFAULT 2 NOT NULL, -- 0: epuisé, 1: peu, 2: disponible
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NULL,
  CONSTRAINT chk_legumes_unite_mesure CHECK(unite_mesure >= 0 AND unite_mesure <= 2),
  CONSTRAINT chk_legumes_stock CHECK(stock >= 0 AND stock <= 2)
);

DROP TABLE IF EXISTS rel_legumes_farmers;
CREATE TABLE rel_legumes_farmers (
  farmer_id INTEGER NOT NULL CHECK (farmer_id > 0),
  legume_id INTEGER NOT NULL CHECK (legume_id > 0),
  CONSTRAINT fk_rel_legumes_users_farmer_farmer_id FOREIGN KEY(farmer_id) REFERENCES users_farmer(id),
  CONSTRAINT fk_rel_legumes_users_farmer_legume_id FOREIGN KEY(legume_id) REFERENCES legumes(id) ON DELETE CASCADE
);
CREATE INDEX rel_legumes_users_farmer_farmer_id_idx ON rel_legumes_farmers (farmer_id) WITH (deduplicate_items = off);
CREATE INDEX rel_legumes_users_farmer_legume_id_idx ON rel_legumes_farmers (legume_id) WITH (deduplicate_items = off);

-- -- mot de passe julien
-- -- user 1
INSERT INTO users (username, email, password, refresh_token) VALUES ('jules', 'jul@les.fr', '$argon2id$v=19$m=65536,t=1,p=2$L8cbcywf+q8USprZOBx4rg$TGrZ8aUCviK56IemYOvekGoPjiwGb7lg5unqJtSz2ys', '9j9kM4wOTZ38Bt4RMaLIYoSbCkrnLqOP');
INSERT INTO users_permission (user_id) VALUES(1);
INSERT INTO users_farmer (img_url, commune, coordonnees, user_id) VALUES ('https://lh3.googleusercontent.com/proxy/DK2ItEYNnaP5cD6kAq3-PLgdECN2TV-spHWYlhVCRJqFTOY4hzPg_4n_EICO-B_1-xxUKjheo1JkQrFSmdptNFFl81p7Y1Shz4VpQPAw1u28Si47cWW9Sy97MA8qpVV6m6BsKvSmyhsss38DgANpc99y5_X9J2Pa', 'Nonards', ST_SetSRID(ST_MakePoint(1.8095, 45.0174), 4326), 1);
-- -- -- ajout du premier légume
INSERT INTO legumes (img_url, nom, variete, systeme_echange, unite_mesure, stock) VALUES ('https://binette-et-cornichon.com/bundles/chouchieplant/images/plants/carotte/carotte.jpg', 'carotte', 'chantenay à coeur rouge', '{0,1}', 1, 2);
INSERT INTO rel_legumes_farmers (farmer_id, legume_id) VALUES (1, 1);
-- -- -- ajout du premier fruit
INSERT INTO fruits (img_url, nom, variete, systeme_echange, prix, unite_mesure, stock) VALUES ('https://c1dd285b9e085ddb1966-b5ece2cd3a8c2c0d8bc51df36519794c.ssl.cf1.rackcdn.com/boire_vignettes/raisin-grappe-violet-vigne.jpg', 'raisin', 'red globe', '{2}', 0.20, 2, 1);
INSERT INTO rel_fruits_farmers (farmer_id, fruit_id) VALUES (1, 1);
-- -- -- ajout du deuxième légume
INSERT INTO legumes (img_url, nom, variete, systeme_echange, prix, unite_mesure, stock) VALUES ('https://monjardinmamaison.maison-travaux.fr/wp-content/uploads/sites/8/2018/09/gettyimages-531203595-615x410.jpg', 'courge', 'butternut', '{1,2}', 1.90, 2, 0);
INSERT INTO rel_legumes_farmers (farmer_id, legume_id) VALUES (1, 2);



-- -- -- user 2
INSERT INTO users (username, email, password, refresh_token) VALUES ('max234', 'maxium@dupont.fr', '$argon2id$v=19$m=65536,t=1,p=2$rNavMPHif5zbz9FmkED6Kg$+P9RkQHYpBvtwn7gFdYPPvbPbvPO6Rz/E/IxBN+bjc0', 'XzwTH14LvJeTsoOf4pfFpmCTmg09jvQd');
INSERT INTO users_permission (user_id) VALUES(2);
-- -- INSERT INTO users_farmer (img_url, commune, coordonnees, temporary_disabled) VALUES ('https://www.jardiner-malin.fr/wp-content/uploads/2018/06/potager-facile-legumes.jpg', 'Bassignac-le-Bas', ST_SetSRID(ST_MakePoint(1.8713, 45.0272), 4326), 4326), 2);
-- -- -- -- ajout de la première graine
-- -- INSERT INTO graines (img_url, nom, variete, systeme_echange, stock) VALUES ('https://shop.babyplante.fr/721-home_default/25-graines-bambou-geant-moso-phyllostachys-edulis-heterocycla-ou-pubescens.jpg', 'bambusa', 'arundinacea', '{Don}', 0);
-- -- INSERT INTO rel_graines_users (user_id, graine_id) VALUES (2, 1);
-- -- -- -- ajout du deuxième fruit
-- -- INSERT INTO fruits (img_url, nom, variete, systeme_echange, prix, unite_mesure, stock) VALUES ('https://www.mangeonslocal-en-idf.com/sites/default/files/Pomme%20Faro2.JPG', 'pomme', 'pink cripps', '{Vente}', 2.30, 'Kg', 1);
-- -- INSERT INTO rel_fruits_users (user_id, fruit_id) VALUES (2, 2);
-- -- -- -- ajout du troisième légume
-- -- INSERT INTO legumes (img_url, nom, variete, systeme_echange, prix, unite_mesure, stock) VALUES ('https://static.passeportsante.net/230x185/i58858-radis.jpg', 'radis', 'read meat', '{Troc, Vente}', 0.95, 'Botte', 2);
-- -- INSERT INTO rel_legumes_users (user_id, legume_id) VALUES (2, 3);

-- -- -- user 3
INSERT INTO users (username, email, password, refresh_token) VALUES ('soph19', 'sophie@manilo.fr', '$argon2id$v=19$m=65536,t=1,p=2$vFUjB8Bgaei8A38a6VOBlA$SBxuyA/lkB3RskK2mjc5EWY1VtWqoZOPQWCGw0WpQTs', 'yxvil0A6aIKdL1mKgtyQi9jh0aThPfrD');
INSERT INTO users_permission (user_id) VALUES(3);
INSERT INTO users_farmer (img_url, commune, coordonnees, temporary_disabled, user_id) VALUES ('https://www.aquaportail.com/pictures1802/potager-fleuri.jpg', 'Queyssac-les-Vignes', ST_SetSRID(ST_MakePoint(1.771, 44.963), 4326), TRUE, 3);
-- -- -- ajout du troisième fruit
INSERT INTO fruits (img_url, nom, variete, systeme_echange, prix, unite_mesure, stock) VALUES ('https://img-3.journaldesfemmes.fr/BFQk3P_rKysA2PJeVl4HbfAcNbE=/1240x/smart/264a7ff437d54df18588f2899c74c462/ccmcms-jdf/10662444.jpg', 'noix', 'périgord', '{2}', 1.10, 1, 0);
INSERT INTO rel_fruits_farmers (farmer_id, fruit_id) VALUES (2, 2);
-- -- -- ajout du quatrième légume
INSERT INTO legumes (img_url, nom, variete, systeme_echange, prix, unite_mesure, stock) VALUES ('https://fr.rc-cdn.community.thermomix.com/recipeimage/2nc1g3lr-77de7-564089-cfcd2-b67ljhi0/0d0d5043-f22b-4b08-8d9c-2578e0802b37/main/salade-verte.jpg', 'salade', 'romaine', '{1,2}', 0.60, 2, 1);
INSERT INTO rel_legumes_farmers (farmer_id, legume_id) VALUES (2, 3);
-- -- --
INSERT INTO graines (img_url, nom, variete, systeme_echange, stock) VALUES ('https://shop.babyplante.fr/721-home_default/25-graines-bambou-geant-moso-phyllostachys-edulis-heterocycla-ou-pubescens.jpg', 'bambusa', 'arundinacea', {'0'}, 0);
INSERT INTO rel_graines_farmers (farmer_id, graine_id) VALUES (2, 1);


-- -- -- user 4
INSERT INTO users (username, email, password, refresh_token) VALUES ('jere323', 'jeremy@christian.fr', '$argon2id$v=19$m=65536,t=1,p=2$bQtg8LT41c8jPjYTzjgyVg$Yx1P/pH3taGdY0plwe8PQ+WaZdlvZn7zXKY6u448r+4', 'S0FxHRyP8zASgQTPXP13wExZ4W2GmGAc');
INSERT INTO users_permission (user_id) VALUES(4);
INSERT INTO users_farmer (img_url, commune, coordonnees, temporary_disabled, user_id) VALUES ('https://www.gammvert.fr/conseils/sites/default/files/styles/main_image/public/2019-02/AdobeStock_113911976-potager-carre.jpg?itok=KjLxllqD', 'Goulles', ST_SetSRID(ST_MakePoint(2.0691, 45.0067), 4326), TRUE, 4);
-- -- -- ajout du 4eme fruit
INSERT INTO fruits (img_url, nom, variete, systeme_echange, prix, unite_mesure, stock) VALUES ('https://i-reg.unimedias.fr/sites/art-de-vivre/files/styles/large/public/r69-poire-williams_fotolia.jpg?auto=compress%2Cformat&crop=faces%2Cedges&cs=srgb&fit=crop', 'poire', 'nashi', '{0,1}', 0.25, 2, 1);
INSERT INTO rel_fruits_farmers (farmer_id, fruit_id) VALUES (3, 3);
-- -- -- ajout du 5eme fruit
INSERT INTO fruits (img_url, nom, variete, systeme_echange, unite_mesure, stock) VALUES ('https://i-dja.unimedias.fr/sites/art-de-vivre/files/styles/large/public/dj_pomme_pressoir_cueillette.jpg?auto=compress%2Cformat&crop=faces%2Cedges&cs=srgb&fit=crop&h=600&w=900', 'pomme', 'boskoop', '{0}', 2, 2);
INSERT INTO rel_fruits_farmers (farmer_id, fruit_id) VALUES (3, 4);
-- -- -- ajout du 6eme fruit
INSERT INTO fruits (img_url, nom, variete, systeme_echange, prix, unite_mesure, stock) VALUES ('https://www.agrimaroc.ma/wp-content/uploads/clementines-espagne-1.jpg', 'clémentine', 'corse', '{2}', 2.49, 1, 1);
INSERT INTO rel_fruits_farmers (farmer_id, fruit_id) VALUES (3, 5);


-- -- -- user 5
INSERT INTO users (username, email, password, refresh_token) VALUES ('bernard23', 'bernard@martel.fr', '$argon2id$v=19$m=65536,t=1,p=2$vGWS1IO7FhZAebXwq5Dn1Q$ugGEmlsES3F6oyFSKyRSpsCZrw0BDV9fBP/MszbQguQ', 'nmevzd7wZXJqJ5pqIvnVGxLYoxmilz49');
INSERT INTO users_permission (user_id) VALUES(5);
-- -- -- ajout du 7eme fruit
-- -- INSERT INTO fruits (img_url, nom, variete, systeme_echange, prix, unite_mesure, stock) VALUES ('https://static.aujardin.info/cache/th/img9/punica-granatum-grenade-600x450.jpg', 'grenade', 'lambda', '{Troc, Vente}', 0.55, 'Piece', 2);
-- -- INSERT INTO rel_fruits_users (user_id, fruit_id) VALUES (5, 7);
-- -- -- ajout de la 2eme graine
-- -- INSERT INTO graines (img_url, nom, variete, systeme_echange, prix, stock) VALUES ('https://www.manutan.fr/img/S/GRP/ST/AIG2132442.jpg', 'ficus', 'benjamina', '{Vente}', 2.30, 2);
-- -- INSERT INTO rel_graines_users (user_id, graine_id) VALUES (5, 2);
-- -- -- ajout de la 3eme graine
-- -- INSERT INTO graines (img_url, nom, variete, systeme_echange, prix, stock) VALUES ('https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRPAfmB2YF9BxmWCfg_b9eRzJClY1XLrfJYOg&usqp=CAU', 'sarracenia', 'purpurea', '{Troc, Vente}', 3.10, 1);
-- -- INSERT INTO rel_graines_users (user_id, graine_id) VALUES (5, 3);

-- -- -- user 6
INSERT INTO users (username, email, password, refresh_token) VALUES ('jeAnne78', 'jeanne@bineaux.fr', '$argon2id$v=19$m=65536,t=1,p=2$MTJFgvIUstpyCSaImFqH7A$0q1KUvaFQCBOk/gHEqzFF+0Fs4MNKAduEAyioRrlgF', 'CK0MJ4NF9vq4VCSjMsSxsZxIkAMqvrjo');
INSERT INTO users_permission (user_id) VALUES(6);
INSERT INTO users_farmer (img_url, commune, coordonnees, user_id) VALUES ('https://www.jardindeco.com/data/img/contentjardin_en_carres_marge%20gche.jpg', 'Nonards', ST_SetSRID(ST_MakePoint(1.8316, 45.0034), 4326), 6);
-- -- -- ajout du 8ème fruit
INSERT INTO fruits (img_url, nom, variete, systeme_echange, prix, unite_mesure, stock) VALUES ('https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcT11SXm2FkaQBmc6JSF3FEML0U9NTM3IUNjcw&usqp=CAU', 'pomme', 'corrèze', '{2}', 1.80, 1, 1);
INSERT INTO rel_fruits_farmers (farmer_id, fruit_id) VALUES (4, 6);
-- -- -- ajout du 5ème légume
INSERT INTO legumes (img_url, nom, variete, systeme_echange, unite_mesure, stock) VALUES ('https://wordpress.potagercity.fr/wp-content/uploads/2019/02/fiche-produit-radis-roses.jpg', 'radis', 'japonais', '{1}', 0, 2);
INSERT INTO rel_legumes_farmers (farmer_id, legume_id) VALUES (4, 4);
-- -- --
INSERT INTO graines (img_url, nom, variete, systeme_echange, prix, stock) VALUES ('https://www.manutan.fr/img/S/GRP/ST/AIG2132442.jpg', 'ficus', 'benjamina', '{1}', 2.30, 2);
INSERT INTO rel_graines_farmers (farmer_id, graine_id) VALUES (4, 2);

-- -- -- user 7
INSERT INTO users (username, email, password, refresh_token) VALUES ('l3432ea', 'lea@moreau.fr', '$argon2id$v=19$m=65536,t=1,p=2$t9cdmsT3KMLvYe6givmLTw$qN0c8LTzSvviob9mU7BIfGRv1gzlUVxoyVuGMUCtnpY', '9bGXDIIdYm3AJyWS44a2EDQbsUZ7iNSm');
INSERT INTO users_permission (user_id) VALUES(7);
INSERT INTO users_farmer (img_url, commune, coordonnees, temporary_disabled, user_id) VALUES ('https://www.jardindeco.com/data/img/contentjardin_en_carres_marge%20gche.jpg', 'Bretenoux', ST_SetSRID(ST_MakePoint(1.845,44.9145), 4326), TRUE, 7);
-- -- -- ajout du 9ème fruit
INSERT INTO fruits (img_url, nom, variete, systeme_echange, prix, unite_mesure, stock) VALUES ('https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcS5rCUWhyAruXs4b1E4wwTGJkPi-hIuPKLB0w&usqp=CAU', 'pomme', 'golden jaune', '{1,2}', 0.50, 2, 2);
INSERT INTO rel_fruits_farmers (farmer_id, fruit_id) VALUES (5, 7);
-- -- -- ajout du 6ème légume
INSERT INTO legumes (img_url, nom, variete, systeme_echange, unite_mesure, stock) VALUES ('https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQocN5goGOd4MDJu5J35CcFvC1NEC8lqYMb7w&usqp=CAU', 'raisin', 'aladin', '{0,1}', 0, 2);
INSERT INTO rel_legumes_farmers (farmer_id, legume_id) VALUES (5, 5);
-- -- --
INSERT INTO graines (img_url, nom, variete, systeme_echange, prix, stock) VALUES ('https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRPAfmB2YF9BxmWCfg_b9eRzJClY1XLrfJYOg&usqp=CAU', 'sarracenia', 'purpurea', '{1,2}', 3.10, 1);
INSERT INTO rel_graines_farmers (farmer_id, graine_id) VALUES (5, 3);

-- -- -- user 8
INSERT INTO users (username, email, password, refresh_token) VALUES ('camille192', 'camille@fauduit.fr', '$argon2id$v=19$m=65536,t=1,p=2$NMJ/Xy6PGg47WXKaXV6Bmg$AqMqdnQhCqr78v0XzQ8Hmu75aeaOb5oQwxU4YHtdYZk', 'yzbnPCIMPIlobUkH4iULZYtazS5f34z0');
INSERT INTO users_permission (user_id) VALUES(8);
INSERT INTO users_farmer (img_url, commune, coordonnees, temporary_disabled, user_id) VALUES ('https://upload.wikimedia.org/wikipedia/commons/0/06/Jardin_potager_001.JPG', 'Sérilhac', ST_SetSRID(ST_MakePoint(1.7445569470490274, 45.10304302614448), 4326), TRUE, 8);
-- -- -- ajout du 9ème fruit
INSERT INTO graines (img_url, nom, variete, systeme_echange, prix, stock) VALUES ('https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRPAfmB2YF9BxmWCfg_b9eRzJClY1XLrfJYOg&usqp=CAU', 'sésame', 'peuimporte', '{2}', 2.40, 1);
INSERT INTO rel_graines_farmers (farmer_id, graine_id) VALUES (6, 4);
-- -- -- ajout du 6ème légume
INSERT INTO legumes (img_url, nom, variete, systeme_echange, unite_mesure, stock) VALUES ('https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQocN5goGOd4MDJu5J35CcFvC1NEC8lqYMb7w&usqp=CAU', 'raisin', 'bordeau', '{1}', 0, 2);
INSERT INTO rel_legumes_farmers (farmer_id, legume_id) VALUES (6, 6);
-- -- --
INSERT INTO fruits (img_url, nom, variete, systeme_echange, prix, unite_mesure, stock) VALUES ('https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcS5rCUWhyAruXs4b1E4wwTGJkPi-hIuPKLB0w&usqp=CAU', 'pomme', 'golden jaune', '{1,2}', 0.80, 2, 0);
INSERT INTO rel_fruits_farmers (farmer_id, fruit_id) VALUES (6, 8);