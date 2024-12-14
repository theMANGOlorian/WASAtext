
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
        };
    },
    mounted() {
        this.auth = sessionStorage.getItem('Auth');
    },

    watch: {
        conversation: {
            immediate: true,
            handler(newConversation) {
                if (newConversation) {
                    this.fetchMessages(newConversation.conversationId);
                } else {
                    this.messages = [];
                }
            },
        },
    },
    methods: {
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
                this.messages = response.data.messages;
                for (let message of this.messages) {
                    if (message.typeContent === 'photo') {
                        await this.GetImages(message);
                    }
                }
            } catch (error) {
                console.error('Errore nel caricamento dei messaggi:', error.message);
            }
        },
        async sendMessage() {
            if (!this.newMessage.trim()) return;

            const message = {
                text: this.newMessage,
                timestamp: new Date().toISOString(),
                senderId: this.auth,
                typeContent: 'text',
            };

            try {
                const response = await this.$axios.post(`/conversations/${this.conversation.conversationId}/text`, message, {
                    headers: {
                        Authorization: `Bearer ${sessionStorage.getItem('Auth')}`,
                    },
                });
                this.messages.push(response.data.message);
                this.newMessage = '';
                this.scrollToBottom();
                if (response.data.message.typeContent === 'image') {
                    await this.GetImages(response.data.message);
                }
            } catch (error) {
                console.error('Errore durante l\'invio del messaggio:', error.message);
            }
        },
        scrollToBottom() {
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
                // Recupera l'immagine dal server
                let response = await this.$axios.get(`/images/${message.image}/photo`, {
                    headers: {
                        Authorization: `Bearer ${this.auth}`,
                    },
                    responseType: 'blob',  // Imposta il tipo di risposta per immagini
                });
                
                // Crea un URL per l'immagine
                const imageUrl = URL.createObjectURL(response.data);
                
                // Aggiorna il messaggio con l'URL dell'immagine
                message.imageUrl = imageUrl;


            } catch (error) {
                console.error('Errore nel recupero dell\'immagine:', error.message);
            }
        },
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
                    class="message-item"
                >
                    <div :class="{'sent': message.senderId === this.auth, 'received': message.senderId !== this.auth}">
                        <p class="sender-name">{{ message.username }}</p>
                        <p v-if="message.typeContent === 'text'" class="message-text">{{ message.text }}</p>
                        <img v-else-if="message.typeContent === 'photo'" :src="message.imageUrl" alt="Image" class="message-image" />
                        <span class="message-time">{{ formatDate(message.timestamp) }}</span>
                    </div>
                </div>
            </div>

            <footer class="message-input">
                <input
                    v-model="newMessage"
                    type="text"
                    placeholder="Scrivi un messaggio..."
                    @keyup.enter="sendMessage"
                />
                <button class="send-message-button" @click="sendMessage">Send</button>
                <button class="load-image-button" @click="loadImage">+</button>
            </footer>
        </div>
    </div>
</template>

<style scoped>
.chat-box {
    display: flex;
    flex-direction: column;
    height: 100%;
    background-color: #f9f9f9;
}

.no-conversation {
    display: flex;
    justify-content: center;
    align-items: center;
    height: 100%;
    color: #aaa;
    font-size: 1.2em;
}


.message-list {
    flex: 1;
    overflow-y: auto;
    padding: 10px;
}

.message-item {
    margin: 10px 0;
    display: flex; 
    flex-direction: column;  /* Allinea verticalmente il contenuto */
}

.sent,
.received {
    padding: 10px;
    border-radius: 10px;
    max-width: 70%;
    display: flex;
    flex-direction: column;  /* Allinea verticalmente il contenuto */
}

.sent {
    align-self: flex-end;
    background-color: #1edb0548;
}

.received {
    align-self: flex-start;
    background-color: #dddddd7e;
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
    margin-top: 5px; /* Spazio tra l'immagine/testo e l'orario */
    display: block;  /* Forza l'orario su una nuova riga */
}


.message-input {
    display: flex;
    padding: 10px;
    border-top: 1px solid #ddd;
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

</style>
