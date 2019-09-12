CREATE DATABASE IF NOT EXISTS userDB;

USE userDB;

create table if not exists users (
    id int not null auto_increment primary key,
    email varchar(128) not null,
    first_name varchar(64) not null,
    last_name varchar(128) not null,
    user_name varchar(64) not null,
    pass_hash binary(60) not null,
    photo_url varchar(128) not null
);

create table if not exists sign_ins (
    id int not null auto_increment primary key,
    user_id int not null,
    date_time datetime not null,
    ip_address varchar(128) not null
);

-- create new for user leaderBoard (id, user1, user2, numguess right)
create table if not exists leaderboards (
    id int not null auto_increment primary key,
    actorID int not null,
    guesserID int not null,
    numGuessRight int not null,
    numGuessPlayed int not null
);
