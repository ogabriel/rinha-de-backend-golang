CREATE INDEX pessoas_busca_gist_trgm_ops_index ON pessoas USING gist (busca gist_trgm_ops (siglen = '1024'));
