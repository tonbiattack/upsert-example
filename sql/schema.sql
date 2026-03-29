CREATE TABLE IF NOT EXISTS user_settings (
  user_id    INT NOT NULL,
  theme      VARCHAR(50) NOT NULL,
  language   VARCHAR(10) NOT NULL,
  updated_at DATETIME NOT NULL,
  PRIMARY KEY (user_id)
);

CREATE TABLE IF NOT EXISTS orders (
  id         INT NOT NULL AUTO_INCREMENT,
  user_id    INT NOT NULL,
  amount     DECIMAL(10,2) NOT NULL,
  status     VARCHAR(20) NOT NULL,
  ordered_at DATETIME NOT NULL,
  PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS monthly_order_summaries (
  user_id        INT NOT NULL,
  target_month   DATE NOT NULL,
  order_count    INT NOT NULL,
  total_amount   DECIMAL(12,2) NOT NULL,
  aggregated_at  DATETIME NOT NULL,
  PRIMARY KEY (user_id, target_month),
  UNIQUE KEY uq_user_month (user_id, target_month)
);

CREATE TABLE IF NOT EXISTS contracts (
  id         INT NOT NULL AUTO_INCREMENT,
  company_id INT NOT NULL,
  plan       VARCHAR(50) NOT NULL,
  signed_at  DATE NOT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY uq_company (company_id)
);

CREATE TABLE IF NOT EXISTS user_tags (
  user_id INT NOT NULL,
  tag_id  INT NOT NULL,
  PRIMARY KEY (user_id, tag_id)
);
