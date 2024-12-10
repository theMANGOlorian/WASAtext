--- questo trigger non servirà più
CREATE TRIGGER update_message_status
AFTER UPDATE ON members
FOR EACH ROW
BEGIN
    -- Aggiorna lo stato dei messaggi a 'recv' se tsLastRecv è stato modificato
    UPDATE messages
    SET status = 'recv'
    WHERE conversation = NEW.conversationId
      AND timestamp <= NEW.tsLastRecv
      AND status = 'none'
      AND NOT EXISTS (
        SELECT 1
        FROM members m
        WHERE m.conversationId = NEW.conversationId
      AND m.userId != messages.sender
      AND m.tsLastRecv <= messages.timestamp
        )
      AND NEW.tsLastRecv IS NOT NULL
      AND NEW.tsLastRecv != OLD.tsLastRecv; -- Verifica che tsLastRecv sia cambiato

    -- Aggiorna lo stato dei messaggi a 'read' se tsLastRead è stato modificato
    UPDATE messages
    SET status = 'read'
    WHERE conversation = NEW.conversationId
      AND timestamp <= NEW.tsLastRead
      AND status IN ('none', 'recv')  -- Solo messaggi con stato 'none' o 'recv'
      AND NOT EXISTS (
        SELECT 1
        FROM members m
        WHERE m.conversationId = NEW.conversationId
      AND m.userId != messages.sender
      AND m.tsLastRead <= messages.timestamp
        )
      AND NEW.tsLastRead IS NOT NULL
      AND NEW.tsLastRead != OLD.tsLastRead; -- Verifica che tsLastRead sia cambiato
END;


CREATE TRIGGER IF NOT EXISTS check_reaction_is_allowed
BEFORE INSERT ON reactions
FOR EACH ROW
BEGIN
SELECT
    CASE
        WHEN NOT EXISTS (
            SELECT 1 FROM members m JOIN conversations c ON m.conversationId = c.id JOIN messages msg ON msg.conversation = c.id
            WHERE m.userId = NEW.owner
              AND NEW.messageId = msg.id
        ) THEN
            RAISE (ABORT, 'NOT ALLOWED: user is not a member')
        END;
END;



