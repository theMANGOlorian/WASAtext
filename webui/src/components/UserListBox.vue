<template>
    <div class="user-list-overlay" v-on:click.self="close">
        <div class="user-list">
            <h2>Inoltra a ...</h2>
            <ul>
                <li 
                    v-for="recipient in recipientsList" 
                    :key="recipient.id"  
                    v-on:click="onRecipientClick(recipient)"
                >
                    <div class="contact-item">
                        <span class="contact-name">{{ recipient.name }}</span>
                    </div>
                </li>
            </ul>
        </div>
    </div>
</template>

<script>
import defaultPicture from '../assets/images/default-profile.jpg';
export default {
    props: {
        messageId: {
            type: String,
            required: true
        }
    },

    data() {
        return {
            auth: '',
            conversationsList: [],
            usersList: [],
            recipientsList: [], // Unisce utenti e gruppi
            defaultProfilePicture: defaultPicture,
        };
    },

    async beforeMount() {
        this.auth = sessionStorage.getItem('Auth');
        await this.fetchConversationsList();
        await this.fetchUsersList();
        this.combineRecipients();
    },

    methods: {
        async fetchUsersList() {
            try {
                const response = await this.$axios.get(`/users`, {
                    headers: { Authorization: `Bearer ${this.auth}` },
                });
                this.usersList = response.data.users?.map(username => ({
                    username,
                    photoURL: this.defaultProfilePicture 
                })) || [];
            } catch (error) {
                console.error("Error fetching users:", error.message);
                alert("Error Message: " + error.message);
            }
        },

        async fetchConversationsList() {
            try {
                const response = await this.$axios.get(`/users/${this.auth}/conversations`, {
                    headers: { Authorization: `Bearer ${this.auth}` },
                });
                this.conversationsList = response.data.conversations;
            } catch (error) {
                console.error("Error fetching conversations:", error.message);
                alert("Error Message: " + error.message);
            }
        },

        combineRecipients() {
            // Creiamo un Set con i nomi delle conversazioni esistenti
            const existingConversationNames = new Set(this.conversationsList.map(convo => convo.conversationName));

            this.recipientsList = [
                // Aggiungiamo i gruppi dalla conversationsList
                ...this.conversationsList.map(conversation => ({
                    type: 'group',
                    name: conversation.conversationName,
                    id: conversation.conversationId,
                })),
                
                // Aggiungiamo solo gli utenti che NON hanno già una conversazione
                ...this.usersList
                    .filter(user => !existingConversationNames.has(user.username))
                    .map(user => ({
                        type: 'private',
                        name: user.username,
                        id: user.username,
                    }))
            ];
        },


        async createConversation(username) {
            try {
                const response = await this.$axios.post(`/users/${this.auth}/conversations`, {
                    name: username,
                    typeConversation: "private",
                }, {
                    headers: { Authorization: `Bearer ${this.auth}` },
                });
                return response.data.identifier;
            } catch (error) {
                console.error("Error creating conversation:", error.message);
                alert("Error Message: " + error.message);
            }
        },

        async forwardMessage(conversationId) {
            try {
                await this.$axios.post(`messages/${this.messageId}/forward`, {
                    conversationId: conversationId,
                }, {
                    headers: { Authorization: `Bearer ${this.auth}` },
                });
            } catch (error) {
                console.error("Error forwarding message:", error.message);
                alert("Error Message: " + error.message);
            }
        },

        async onRecipientClick(recipient) {
            if (recipient.type === 'group') {
                await this.forwardMessage(recipient.id);
            } else {
                const existingConversation = this.conversationsList.find(convo => convo.conversationName === recipient.name);
                if (existingConversation) {
                    await this.forwardMessage(existingConversation.conversationId);
                } else {
                    const conversationId = await this.createConversation(recipient.name);
                    await this.forwardMessage(conversationId);
                }
            }
            this.close();
        },

        close() {
            this.$emit('close');
        }
    }
};
</script>

<style scoped>
.user-list-overlay {
    position: fixed;
    top: 0;
    left: 0;
    width: 100vw;
    height: 100vh;
    background: rgba(0, 0, 0, 0.5);
    display: flex;
    justify-content: center;
    align-items: center;
    z-index: 1000;
}

.user-list {
    background: #fff;
    border-radius: 8px;
    box-shadow: 0 4px 10px rgba(0, 0, 0, 0.3);
    width: 350px;
    padding: 20px;
    text-align: center;
}

.user-list h2 {
    margin-bottom: 20px;
    font-size: 1.5rem;
    color: #333;
}

.user-list ul {
    list-style: none;
    padding: 0;
    margin: 0;
}

.user-list li {
    margin: 10px 0;
    background: #f9f9f9;
    border: 1px solid #e0e0e0;
    border-radius: 8px;
    overflow: hidden;
    cursor: pointer;
    transition: background 0.3s, transform 0.2s;
}

.user-list li:hover {
    background: #65ee61;
    transform: scale(1.02);
}

.contact-item {
    display: flex;
    align-items: center;
    padding: 10px;
}

.contact-name {
    font-size: 1rem;
    font-weight: 500;
    color: #333;
    text-align: left;
    flex-grow: 1;
}
</style>
