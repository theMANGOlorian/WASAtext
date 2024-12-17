
<script>
export default {
    props: {
        conversation: {
            type: Object,
            required: false,
        },
    },
    data() {
        return {
            messages: [],
            newMessage: '',
            auth: '',
            pollingInterval: null,
            lastMessagesHash: null, // Salva l'hash della lista dei messaggi
            selectedMessage: null,
            emojiMenuVisible: false, // Per la seconda ContextMenu
            emojiList: ["😊", "😂", "😍", "😎", "😭"],
            replyTo: null,
            highlightedMessageId: null
        };
    },
    mounted() {
        this.auth = sessionStorage.getItem('Auth');
        if (this.conversation) {
            this.startPolling();
        }
        document.addEventListener('click', this.handleClickOutside);
    },

    watch: {
        conversation: {
            immediate: true,
            handler(newConversation) {
                if (newConversation) {
                    this.fetchMessages(newConversation.conversationId);
                    this.startPolling();
                } else {
                    this.stopPolling();
                    this.messages = [];
                    this.lastMessagesHash = null;
                }
            },
        },
    },

    beforeDestroy() {
        this.stopPolling(); // Ferma il polling quando il componente viene distrutto
        document.removeEventListener('click', this.handleClickOutside);
    },
    
    methods: {
        startPolling() {
            this.stopPolling(); // Evita duplicati
            this.pollingInterval = setInterval(() => {
                this.fetchMessages(this.conversation.conversationId);
            }, 5000); // Ogni 5 secondi
        },

        stopPolling() {
            if (this.pollingInterval) {
                clearInterval(this.pollingInterval);
                this.pollingInterval = null;
            }
        },

        async fetchMessages(conversationId) {
            try {
                const response = await this.$axios.get(`/conversations/${conversationId}/messages`, {
                    headers: {
                        Authorization: `Bearer ${this.auth}`,
                    },
                    params: {
                        limit: 50,
                        cursor: '',
                    },
                });

                let newMessages = response.data.messages;

                if (response.data.messages === null) {
                    this.updateMessages([]);
                    this.lastMessagesHash = null;
                }else {
                    const newHash = this.computeMessagesHash(newMessages);

                    // controlla se ci sono modifiche nei messaggi
                    if (newHash !== this.lastMessagesHash) {
                        this.lastMessagesHash = newHash;
                        this.updateMessages(newMessages);
                    }
                }
            } catch (error) {
                console.error('Errore nel caricamento dei messaggi:', error.message);
            }
        },

        updateMessages(newMessages) {
            this.messages = newMessages;
            for (let message of this.messages) {
                if (message.typeContent === 'photo') {
                    this.GetImages(message);
                }
            }
        },

        computeMessagesHash(messages) {
            return messages.map(msg => {
                const reactionsHash = msg.reactions ? msg.reactions.map(reaction => `${reaction.emoji}-${reaction.userId}`).join(',') : '';
                return `${msg.messageId}-${msg.text || ''}-${reactionsHash}`;
            }).join('|');
        },


        async sendMessage() {
            if (!this.newMessage.trim()) return;

            try {
                const response = await this.$axios.post(`/conversations/${this.conversation.conversationId}/text`, 
                {
                    bodyMessage: this.newMessage,
                    replyTo: this.replyTo.messageId,
                },
                {
                    headers: {
                        Authorization: `Bearer ${this.auth}`,
                    },
                });
                this.messages.push(response.data);
                this.newMessage = '';
                this.scrollToBottom();
                this.replyTo = null;
                await this.fetchMessages(this.conversation.conversationId);
            } catch (error) {
                console.error('Errore durante l\'invio del messaggio:', error.data);
            }
        },

        async sendPhoto(file) {
            try {
                if (file.type !== 'image/png') {
                    this.ErrorMessage = "Only PNG images are allowed.";
                    return;
                }

                let response = await this.$axios.post(`/conversations/${this.conversation.conversationId}/image`, file, 
                {
                    headers: {
                        Authorization: `Bearer ${this.auth}`,
                        'Content-Type': 'image/png',
                    },
                });

                await this.fetchMessages(this.conversation.conversationId);
            } catch (error) {
                console.error('Errore durante l\'invio della foto:', error.message);
            }
        },

        handleFileChange(event) {
            const file = event.target.files[0];
            if (file) {
                this.sendPhoto(file);
            }
        },

        async scrollToBottom() {
            this.$nextTick(() => {
                if (this.$refs.messageList) {
                    this.$refs.messageList.scrollTop = this.$refs.messageList.scrollHeight;
                }
            });
        },

        formatDate(timestamp) {
            const options = {
                year: 'numeric',
                month: 'short',
                day: 'numeric',
                hour: 'numeric',
                minute: 'numeric',
            };
            return new Date(timestamp).toLocaleString('it-IT', options);
        },

        async GetImages(message) {
            try {
                let response = await this.$axios.get(`/images/${message.image}/photo`, {
                    headers: {
                        Authorization: `Bearer ${this.auth}`,
                    },
                    responseType: 'blob',
                });

                const imageUrl = URL.createObjectURL(response.data);
                message.imageUrl = imageUrl;
            } catch (error) {
                console.error('Errore nel recupero dell\'immagine:', error.message);
            }
        },

        onImageClick() {
            this.$refs.fileInput.click();
        },

        handleClickOutside(event) {
            this.$refs.contextMenu.closeMenu();
        },

        onMessageClick(event, message) {
            event.preventDefault(); // previene il menu del browser
            let options = ['Reply','Forward','React'];
            
            if (message.reactions && message.reactions.length > 0) {
                const userReaction = message.reactions.find(reaction => reaction.owner === this.auth);
                if (userReaction) {
                    options.push('Unreact');
                }
            }
            if (this.auth === message.senderId) {
                options.push('Delete')
            }

            this.selectedMessage = message;
            this.$refs.contextMenu.openMenu(event.clientX, event.clientY, options);
        },

        onContextMenuOptionClick(option) {
            if (option === 'Delete') {
                this.deleteMessage(this.selectedMessage.messageId);
            } else if (option === 'Reply') {
                this.replyTo = this.selectedMessage
            } else if (option === 'Forward') {
                console.log('Inoltra il messaggio.');
            } else if (option === 'React') {
                this.showEmojiMenu();
            } else if (option === 'Unreact') {
                this.removeReaction();
            }
        },

        showEmojiMenu() {
            const { x, y } = this.$refs.contextMenu.getMenuPosition();
            this.$refs.emojiMenu.openMenu(x, y, this.emojiList);
        },

        onEmojiSelect(emoji) {
            this.reactMessage(emoji),
            // Chiudi la ContextMenu dopo la selezione
            this.$refs.emojiMenu.closeMenu();
        },

        async reactMessage(emoji) {
            try {
                const response = await this.$axios.put(`/messages/${this.selectedMessage.messageId}/comment`,
                {
                    "Reaction": emoji
                }, 
                {
                    headers: {
                        Authorization: `Bearer ${this.auth}`,
                    },
                });
                await this.fetchMessages(this.conversation.conversationId);
            } catch (error) {
                console.error('Errore durante l\'invio della reazione:', error.message);
            }
        },

        async removeReaction() {
            try {
                const response = await this.$axios.delete(`/messages/${this.selectedMessage.messageId}/comment`,
                {
                    headers: {
                        Authorization: `Bearer ${this.auth}`,
                    },
                });
                await this.fetchMessages(this.conversation.conversationId);
            } catch (error) {
                console.error('Errore durante la cancellazione della reazione', error.message);
            }
        },

        async deleteMessage(messageId) {
            try {
                const response = await this.$axios.delete(`/messages/${messageId}`, 
                {
                    headers: {
                        Authorization: `Bearer ${this.auth}`,
                    },
                });
                await this.fetchMessages(this.conversation.conversationId);
            } catch (error) {
                console.error('Errore durante la cancellazione del messaggio:', error.message);
            }
        },
        
        truncateMessage(message, max) {
            if (message.length > max) {
                return message.substring(0, max) + '...';
            }
            return message;
        },

        scrollToMessage(targetMessageId) {
            // Trova l'indice del messaggio nella lista
            const messageIndex = this.messages.findIndex(
                (message) => message.messageId === targetMessageId
            );

            if (messageIndex !== -1) {
                // Recupera il riferimento del messaggio
                const targetElement = this.$refs.messagesRefs[messageIndex];

                if (targetElement) {
                    // Scorri fino al messaggio selezionato
                    targetElement.scrollIntoView({ behavior: "smooth", block: "center" });

                    // Applica la classe di evidenziazione temporanea
                    this.highlightedMessageId = targetMessageId;

                    // Rimuovi la classe dopo 1 secondo
                    setTimeout(() => {
                        this.highlightedMessageId = null;
                    }, 1000);
                }
            }
        },

        getTextReplied(messageId) {
            const targetMessage = this.messages.find(message => message.messageId === messageId);
            if (targetMessage) {
                return this.truncateMessage(targetMessage.text,20) || "Image";
            } else {
                console.log("Messaggio non trovato: getTextReplied")
                return "";
            }
        }
    },
};
</script>


<template>
    <div class="chat-box">
        <div v-if="!conversation" class="no-conversation">
            <p>Select a conversation and start chatting. Your friends are waiting you!</p>
        </div>
        <div v-else class="conversation">
            <div ref="messageList" class="message-list">
                <div
                    v-for="message in messages"
                    :key="message.messageId"
                    ref="messagesRefs"
                    class="message-item"                >
                    <div :class="{
                            'sent': message.senderId === this.auth, 
                            'received': message.senderId !== this.auth, 
                            'highlighted': highlightedMessageId === message.messageId
                        }" @contextmenu="onMessageClick($event, message)">
                        <p v-if="message.replyTo != ''" @click="scrollToMessage(message.replyTo)" class="reply-message"> {{ getTextReplied(message.replyTo) }}</p>
                        <p class="sender-name">{{ message.username }}</p>
                        <p v-if="message.typeContent === 'text'" class="message-text">{{ message.text }}</p>
                        <img v-else-if="message.typeContent === 'photo'" :src="message.imageUrl" alt="Image" class="message-image" />
                        <div class="reactions-container">
                            <span v-for="reaction in message.reactions" class="reactions">{{ reaction.emoji }}</span>
                        </div>
                        <span class="message-time">{{ formatDate(message.timestamp) }}</span>
                    </div>
                </div>
            </div>
            <div class="reply-box" v-if="this.replyTo !== null"><span>Reply to: {{ this.truncateMessage(this.replyTo.text,100) }}</span></div>
            <footer class="message-input">
                <input
                    v-model="newMessage"
                    type="text"
                    placeholder="Scrivi un messaggio..."
                    @keyup.enter="sendMessage"
                />
                <button class="send-message-button" @click="sendMessage">Send</button>
                <button class="load-image-button" @click="onImageClick">+</button>
                <input type="file" ref="fileInput" style="display: none" @change="handleFileChange" accept="image/png"/>
            </footer>
        </div>
        <ContextMenu ref="contextMenu" @option-click="onContextMenuOptionClick" />
        <ContextMenu ref="emojiMenu" @option-click="onEmojiSelect" />
    </div>
</template>

<style scoped>

.chat-box {
    display: flex;
    flex-direction: column;
    height: 100%;
    background-color: #f9f9f9;
    background-image: url("../assets/images/background.png");
}

.no-conversation {
    display: flex;
    justify-content: center;
    align-items: center;
    height: 100%;
    color: #000000;
    font-size: 1.4em;
}

.conversation {
    display: flex;
    flex-direction: column;
    height: 93.7%; 
    overflow: hidden;
}


.message-list {
    flex: 1;
    overflow-y: auto; 
    padding: 10px;
}

.message-item {
    margin: 10px 0;
    display: flex; 
    flex-direction: column;
}

.sent,
.received {
    padding: 10px;
    border-radius: 10px;
    max-width: 70%;
    display: flex;
    flex-direction: column; 
}

.sent {
    align-self: flex-end;
    background-color: #46e631c7;
}

.received {
    align-self: flex-start;
    background-color: #fffafaea;
}

.message-text {
    margin: 0;
    font-size: 1em;
}

.message-image {
    max-width: 100%;
    border-radius: 10px;
}

.message-time {
    font-size: 0.8em;
    color: #888;
    text-align: right;
    margin-top: 5px; 
    display: block;  
}


.message-input {
    display: flex;
    align-items: center; 
    padding: 10px;
    border-top: 1px solid #ddd;
    background-color: #ffffff;
    width: 100%;
    box-sizing: border-box; 
}


.message-input input {
    flex: 1;
    padding: 10px;
    border: 1px solid #ccc;
    border-radius: 5px;
}

.message-input button {
    border: none;
    background-color: #057722;
    color: white;
    cursor: pointer;
    transition: 0.3s ease;
}
.message-input button:hover {
    background-color: #6fc909;
}

.sender-name {
    font-weight: bolder;
    color: rgb(3, 87, 3)
}

.load-image-button {
    border-radius: 100%;
    margin-left: 10px;
    text-align: center;
    width: 50px;
    height: 50px;
    font-size: x-large;
}
.send-message-button {
    border-radius: 5px;
    margin-left: 10px;
    padding: 10px 20px;
}

.reactions-container {
    width: fit-content;
    background-color: rgba(235, 235, 235, 0.692);
    padding: 1px; 
    border-radius: 5px;
    margin-top: 10px;
}

.reactions {
    margin: 0 1px;
    font-size: 1.3em;
}

.reply-box {
    background-color: rgba(206, 206, 206, 0.801);
    width: fit-content;
    padding: 10px 100px 10px 20px; /*padding: top dx bottom sx */
    border-radius: 10px 10px 0 0;
}

.reply-message {
    background-color: rgba(202, 226, 202, 0.658);
    border-radius: 5px 5px 0 0;
    text-align: center;
    cursor: pointer;
}

.highlighted {
    background-color: #14e4ffb4 !important; 
}
</style>
