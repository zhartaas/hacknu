-- Create users table
CREATE TABLE IF NOT EXISTS chats (
                                     id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                     title      TEXT,
                                     model      TEXT NOT NULL,                  -- e.g. "gpt-4o-mini"
                                     created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

INSERT INTO chats (id, title, model) VALUES
                                         ('1719e433-4215-4450-9a72-ae2ec5956224', 'Sample chat', 'gpt-4o-mini');


CREATE TABLE IF NOT EXISTS messages (
                                        id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                        chat_id    UUID NOT NULL REFERENCES chats(id) ON DELETE CASCADE,
                                        role       VARCHAR(32) NOT NULL,           -- 'user' | 'assistant' | 'system' (if you ever need it)
                                        content    TEXT NOT NULL,
                                        created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

INSERT INTO messages (chat_id, role, content) VALUES
                                                         ('1719e433-4215-4450-9a72-ae2ec5956224', 'user',
                                                          'Hello, how are you?'),
                                                         ('1719e433-4215-4450-9a72-ae2ec5956224', 'assistant',
                                                          'Im just a computer program, but Im here and ready to help you! How can I assist you today?');
