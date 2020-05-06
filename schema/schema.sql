CREATE TABLE books (
   id BIGSERIAL PRIMARY KEY,
   name TEXT NOT NULL,
   pages INT NOT NULL CHECK ( pages > 0 ),
   -- don't remove any data from db
   removed BOOLEAN DEFAULT FALSE,
   file_name TEXT NOT NULL,
   category TEXT,
   book_in_hand TEXT ,
    -- amount_of_likes INT CHECK ( amount_of_likes > 0 )
    amount_of_likes INT
);


drop table books;


CREATE TABLE comments (
    id BIGSERIAL PRIMARY KEY ,
    comment TEXT NOT NULL ,
    book_id  BIGSERIAL REFERENCES books (id),
    commentator_name TEXT NOT NULL,
    removed BOOLEAN DEFAULT FALSE
);

drop table comments


CREATE TABLE users(
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL ,
    login TEXT NOT NULL UNIQUE ,
    password TEXT NOT NULL ,
    role BOOLEAN DEFAULT TRUE ,
    ban BOOLEAN DEFAULT FALSE ,
    removed BOOLEAN DEFAULT FALSE
)