CREATE SCHEMA avito_features;

CREATE TABLE IF NOT EXISTS avito_features.users
(
    id         VARCHAR(255) UNIQUE NOT NULL PRIMARY KEY,
    created_at TIMESTAMP           NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS avito_features.features
(
    slug       VARCHAR(255) NOT NULL PRIMARY KEY,
    created_at TIMESTAMP    NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS avito_features.user_feature
(
    user_id      VARCHAR(255) NOT NULL,
    feature_slug VARCHAR(255) NOT NULL,
    expires_at   TIMESTAMP,

    FOREIGN KEY (user_id) REFERENCES avito_features.users (id) ON DELETE CASCADE,
    FOREIGN KEY (feature_slug) REFERENCES avito_features.features (slug) ON DELETE CASCADE,
    UNIQUE (user_id, feature_slug)
);

CREATE TABLE IF NOT EXISTS avito_features.history
(
    user_id      VARCHAR(255) NOT NULL,
    feature_slug VARCHAR(255) NOT NULL,
    operation    varchar(255),
    date         TIMESTAMP    NOT NULL DEFAULT NOW()
);

CREATE OR REPLACE FUNCTION avito_features.feature_insert_trigger()
    RETURNS TRIGGER AS
$$
BEGIN
    INSERT INTO avito_features.history(user_id, feature_slug, operation, date)
    VALUES (NEW.user_id, NEW.feature_slug, 'add', NOW());
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER user_feature_insert_trigger
    AFTER INSERT
    ON avito_features.user_feature
    FOR EACH ROW
EXECUTE PROCEDURE avito_features.feature_insert_trigger();


CREATE OR REPLACE FUNCTION avito_features.feature_delete_trigger()
    RETURNS TRIGGER AS
$$
BEGIN
    INSERT INTO avito_features.history(user_id, feature_slug, operation, date)
    VALUES (OLD.user_id, OLD.feature_slug, 'delete', NOW());
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER user_feature_delete_trigger
    BEFORE DELETE
    ON avito_features.user_feature
    FOR EACH ROW
EXECUTE PROCEDURE avito_features.feature_delete_trigger();
