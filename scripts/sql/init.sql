CREATE SCHEMA avito_features;

CREATE TABLE IF NOT EXISTS avito_features.users (
        id VARCHAR(255) UNIQUE NOT NULL PRIMARY KEY,
        created_at TIMESTAMP NOT NULL DEFAULT NOW(),
        deleted_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS avito_features.features (
        id SERIAL PRIMARY KEY,
        slug VARCHAR(255) NOT NULL,
        created_at TIMESTAMP NOT NULL DEFAULT NOW(),
        deleted_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS avito_features.user_feature (
        user_id VARCHAR(255) NOT NULL,
        feature_id BIGINT NOT NULL,
        expires_at TIMESTAMP,

        FOREIGN KEY(user_id) REFERENCES avito_features.users(id),
        FOREIGN KEY(feature_id) REFERENCES avito_features.features(id),
        UNIQUE(user_id, feature_id)
);

CREATE TABLE IF NOT EXISTS avito_features.history (
        user_id VARCHAR(255) NOT NULL,
        feature_id BIGINT NOT NULL,
        operation varchar(255),
        date TIMESTAMP NOT NULL DEFAULT NOW(),

        FOREIGN KEY(user_id) REFERENCES avito_features.users(id),
        FOREIGN KEY(feature_id) REFERENCES avito_features.features(id)
);

CREATE OR REPLACE FUNCTION create_history_row()
RETURNS TRIGGER AS $$
BEGIN
  INSERT INTO avito_features.history (user_id, feature_id, date, operation)
  VALUES (NEW.user_id, NEW.feature_id, NOW(), 'added');
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER insert_history_row
AFTER INSERT ON avito_features.user_feature
FOR EACH ROW
EXECUTE FUNCTION create_history_row();


CREATE OR REPLACE FUNCTION delete_history_row()
RETURNS TRIGGER AS $$
BEGIN
  IF OLD.expires_at < NOW() THEN
    INSERT INTO avito_features.history (user_id, feature_id, date, operation)
    VALUES (OLD.user_id, OLD.feature_id, OLD.expires_at, 'deleted');
  ELSE
    INSERT INTO avito_features.history (user_id, feature_id, date, operation)
    VALUES (OLD.user_id, OLD.feature_id, NOW(), 'deleted');
  END IF;
  RETURN OLD;
END;
$$ LANGUAGE plpgsql;


CREATE TRIGGER delete_history_row
BEFORE DELETE ON avito_features.user_feature
FOR EACH ROW
EXECUTE FUNCTION delete_history_row();

