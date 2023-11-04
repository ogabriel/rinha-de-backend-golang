CREATE TABLE IF NOT EXISTS pessoas (
    id uuid PRIMARY KEY,
    apelido character varying(32),
    nome character varying(100),
    nascimento character varying(10),
    stack character varying(32)[],
    busca text,
);
