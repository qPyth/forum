PRAGMA foreign_keys = ON;

CREATE TABLE if not exists users (
           id INTEGER PRIMARY KEY AUTOINCREMENT,
           username VARCHAR(255) NOT NULL UNIQUE ,
           email VARCHAR(255) NOT NULL UNIQUE,
           password VARCHAR(255) NOT NULL
);

CREATE TABLE if not exists posts (
           id INTEGER PRIMARY KEY autoincrement,
           title VARCHAR(255) NOT NULL,
           content TEXT NOT NULL,
           image_path TEXT,
           user_id INTEGER NOT NULL,
           category_id INTEGER NOT NULL,
           created_date DATETIME DEFAULT (datetime(CURRENT_TIMESTAMP, 'localtime')),
            like_count INTEGER DEFAULT 0,
           dislike_count INTEGER DEFAULT 0,
           FOREIGN KEY (user_id) REFERENCES users(id),
           FOREIGN KEY (category_id) REFERENCES categories(id)
);

CREATE TABLE IF NOT EXISTS tags (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    category_id INTEGER NOT NULL,
    name VARCHAR(255) NOT NULL UNIQUE,
    FOREIGN KEY (category_id) REFERENCES categories(id)
);

CREATE TABLE IF NOT EXISTS posts_tags (
    post_id INTEGER,
    tags_id INTEGER,
    FOREIGN KEY (post_id) REFERENCES posts(id),
    FOREIGN KEY (tags_id) REFERENCES tags(id),
    PRIMARY KEY (post_id, tags_id)
);

CREATE TABLE if not exists comments (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            content TEXT NOT NULL,
            user_id INTEGER NOT NULL,
            post_id INTEGER NOT NULL,
            created_date DATETIME DEFAULT (datetime(CURRENT_TIMESTAMP, 'localtime')),
            like_count INTEGER DEFAULT 0,
            dislike_count INTEGER DEFAULT 0,
            FOREIGN KEY (user_id) REFERENCES users(id),
            FOREIGN KEY (post_id) REFERENCES posts(id)
);

CREATE TABLE if not exists sessions (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            user_id INTEGER NOT NULL,
            token TEXT NOT NULL,
            expired_date DATETIME NOT NULL,
            created_date DATETIME DEFAULT (datetime(CURRENT_TIMESTAMP, 'localtime')),
            FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS categories (
            id INT PRIMARY KEY,
            name VARCHAR(255) NOT NULL,
            description VARCHAR(255) NOT NULL,
            url VARCHAR(255) not null
);

CREATE TABLE IF NOT EXISTS votes (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            user_id INTEGER NOT NULL,
            action VARCHAR(255) NOT NULL,
            item_id INTEGER NOT NULL,
            item varchar(255) not null,
            FOREIGN KEY (user_id) REFERENCES users(id)
);

INSERT
    OR IGNORE INTO users (username, email, password)
values ('adikey', 'iskakov921@gmail.com', '$2a$10$Y.Xb0cLGkF1nkpH7g7gr0OcT23mv49LgVL9D43UAkFD4TrVbEB9Xu');

INSERT
    OR IGNORE INTO categories (id, name, description, url)
VALUES (1, 'General discussion', 'The place to discuss general game topics', 'general-discussion'),
       (2, 'Heroes', 'Builds, gameplay, combinations', 'heroes'),
       (3, 'Updates', 'Discussion of new and upcoming patches', 'updates'),
       (4,'LFT/LFP', 'Place to search for teammates or teams', 'lft-lfp'),
       (5,'eSport', 'Discussion of eSports matches', 'esport');


INSERT
    OR IGNORE INTO posts (id, title, content, image_path, user_id, category_id, like_count, dislike_count) VALUES (1,'How to up MMR','Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry''s standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum.','',  1, 1,10,5),
                                                          (2,'How to use quickstarts', 'Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry''s standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum.','', 1,1,11,1),
                                                          (3,'50% system', 'Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry''s standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum.','', 1,1,4,2),
                                                          (4,'Hidden pull is real','Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry''s standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum.','', 1,1,5,3);

INSERT
    OR IGNORE INTO comments (id, content, user_id, post_id,like_count, dislike_count) VALUES (1, 'content content content content content content content content content content', 1,1, 5,4),
                                                                                             (2, 'content content content content content content content content content content', 1,1, 10,2);

INSERT
    OR IGNORE INTO tags (category_id, name) VALUES (1, 'All pick'),
                                                   (1, 'Ranked'),
                                                   (1, 'Turbo'),
                                                   (1, 'Low priority'),
                                                   (1, 'Other'),
                                                   (2, 'Mid lane'),
                                                   (2, 'Offline'),
                                                   (2, 'Safe line'),
                                                   (2, 'Support'),
                                                   (2, 'Full support'),
                                                   (2, 'Other'),
                                                   (3, 'Map'),
                                                   (3, 'Economics'),
                                                   (3, 'Heroes'),
                                                   (3, 'Items'),
                                                   (3, 'Other'),
                                                   (4, 'Ranked'),
                                                   (4, 'Semi-pro'),
                                                   (4, 'Pro-level'),
                                                   (4, 'Other'),
                                                   (5, 'Games'),
                                                   (5, 'Bets'),
                                                   (5, 'Roster shuffle'),
                                                   (5, '322'),
                                                   (5, 'Other');
INSERT
    OR IGNORE INTO posts_tags (post_id, tags_id) VALUES (1,2),
                                                        (2, 1),
                                                        (2, 2),
                                                        (2, 3),
                                                        (3, 2),
                                                        (4, 5)

