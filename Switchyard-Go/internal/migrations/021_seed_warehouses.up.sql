INSERT INTO warehouse (id) VALUES
    ('WH001'), ('WH002'), ('WH003'), ('WH004'), ('WH005'),
    ('WH006'), ('WH007'), ('WH008'), ('WH009')
ON CONFLICT (id) DO NOTHING;
