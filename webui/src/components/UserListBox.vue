<template>
    <div class="user-list-overlay" v-on:click.self="close">
        <div class="user-list">
            <h2>Inoltra a ...</h2>
            <ul>
                <li 
                    v-for="user in usersList" 
                    :key="user.username"  
                    v-on:click="onUserClick(user)" 
                >
                    <div class="contact-item">
                        <span class="contact-name">{{ user.username }}</span>
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
            conversationSelected: '',
            defaultProfilePicture: defaultPicture,
            usersList: [],
        };
    },


    beforeMount() {
        this.auth = sessionStorage.getItem('Auth');
        this.fetchConversationsList()
        this.fetchUsersList();
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
            }
        },
        
        async createConversation(username) {
            try {
              const response = await this.$axios.post(`/users/${sessionStorage.getItem('Auth')}/conversations`, 
              { 
                name: username,
                typeConversation: "private",
              },
              {
                headers: { Authorization: `Bearer ${sessionStorage.getItem('Auth')}`, },
              });
              return response.data.identifier;
            } catch (error) {
                console.error("Error starting a new conversation while forwaring: ", error.message);
            }
        },

        async fetchConversationsList() {

            try {
                const response = await this.$axios.get(`/users/${this.auth}/conversations`, {
                    headers: { Authorization: `Bearer ${this.auth}` },
                });
                this.conversationsList = response.data.conversations;

                this.conversationsList = await Promise.all(
                    response.data.conversations.map(async (conversation) => {
                        const photoURL = await this.getProfilePicture(conversation.photoProfileCode);
                        return {
                            ...conversation,
                            photoURL: photoURL || this.defaultProfilePicture, // URL immagine o default
                        };
                    })
                );

            } 
            catch (error) {
                console.error("Error fetching users:", error.message);
            }
        },

        async forwardMessage(conversationId) {
            try {
                const response = await this.$axios.post(`messages/${this.messageId}/forward`,
                    {
                        conversationId: conversationId,
                    },
                    { 
                        headers: { Authorization: `Bearer ${this.auth}` },
                    }
                );
            }
            catch (error) {
                console.error("Error Forward: ", error.message)
            }
        },

        async getProfilePicture(photoCode) {

            if (photoCode === ""){
                return this.defaultProfilePicture;
            }
            try {
                const response = await this.$axios.get(`images/${photoCode}/photo`,
                    { 
                        headers: { Authorization: `Bearer ${this.auth}` },
                        responseType: 'blob',
                    }
                    
                );
                return URL.createObjectURL(response.data);
            }
            catch (error) {
                console.error("Error nella richiesta per la foto", error.message);
                return this.defaultProfilePicture;
            }
        },

        close() {
            this.$emit('close'); // chiude la lista
        },

        async onUserClick(user) {
            const username = user.username;
            // Verifica se l'utente è già presente nelle conversazioni
            const existingConversation = this.conversationsList.find(conversation => conversation.conversationName === username);

            if (existingConversation) {
                // Se l'utente è già nelle conversazioni, forwarda il messaggio
                this.forwardMessage(existingConversation.conversationId);
            } else {
                // Se l'utente non è nelle conversazioni, crea una nuova conversazione
                const conversationId = await this.createConversation(username);
                // Dopo aver creato la conversazione, inoltra il messaggio alla nuova conversazione
                this.forwardMessage(conversationId);
            }
            this.close(); // Chiudi la lista dopo l'azione
        },
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
/*
.contact-photo {
    width: 40px;
    height: 40px;
    border-radius: 50%;
    margin-right: 15px;
    object-fit: cover;
    background: #ddd;
}
*/

.contact-name {
    font-size: 1rem;
    font-weight: 500;
    color: #333;
    text-align: left;
    flex-grow: 1;
}
</style>
