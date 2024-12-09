CREATE TRIGGER IF NOT EXISTS insert_own_status_messagge
AFTER INSERT ON messages
FOR EACH ROW
BEGIN
    INSERT INTO users_msg (userId, msgId)
    SELECT userId, NEW.id
    FROM members
    WHERE conversationId = NEW.conversation;
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



