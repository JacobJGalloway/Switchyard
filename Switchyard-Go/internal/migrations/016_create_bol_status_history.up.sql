CREATE TABLE plan_bol_status_history (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    plan_bol_id UUID        NOT NULL REFERENCES plan_bol_record(id),
    from_status TEXT,
    to_status   TEXT        NOT NULL,
    changed_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_bol_status_history_bol ON plan_bol_status_history(plan_bol_id, changed_at);

CREATE OR REPLACE FUNCTION fn_record_bol_status_change()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO plan_bol_status_history(plan_bol_id, from_status, to_status, changed_at)
    VALUES (
        NEW.id,
        CASE WHEN TG_OP = 'INSERT' THEN NULL ELSE OLD.status END,
        NEW.status,
        now()
    );
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger fires on INSERT (captures initial status) and on UPDATE only when
-- status actually changes — no noise from other column updates.
CREATE TRIGGER trg_bol_status_history
AFTER INSERT OR UPDATE OF status ON plan_bol_record
FOR EACH ROW EXECUTE FUNCTION fn_record_bol_status_change();
