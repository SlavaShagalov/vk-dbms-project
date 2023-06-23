INSERT INTO users(nickname, fullname, about, email)
VALUES ('slava', 'Slava', 'Student', 'slava@vk.com'),
       ('kirill', 'Kirill', 'Student', 'kirill@vk.com'),
       ('petya', 'Petya', 'Student', 'petya@vk.com'),
       ('evgenii', 'Evgenii', 'Student', 'evgenii@vk.com');

INSERT INTO forums(title, slug, user_nickname)
VALUES ('gig', 'gigs', 'slava'),
       ('gig2', 'gigs2', 'slava'),
       ('gig3', 'gigs3', 'kirill');

INSERT INTO threads(title, author, forum, message, slug)
VALUES ('qwerty', 'slava', 'gigs', 'hahahaha', 'sdfa3212'),
       ('asdfgh', 'kirill', 'gigs', 'lalalalala', 'sdfa3212312'),
       ('zxcvb', 'petya', 'gigs3', 'rerererererere', 'sdfa3232233');

