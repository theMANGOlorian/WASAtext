<template>
    <div>
      <!-- ciclo sulle conversazioni -->
      <div v-if="conversations.length">
        <div v-for="conversation in conversations" :key="conversation.conversationId" @click="selectConversation(conversation)" class="conversation-item">

          <img :src="conversation.photoURL" alt="Profile Picture" class="profile-image" />
          
          <div class="conversation-details">
            <h3 class="name-chat">{{ conversation.conversationName }}</h3>
            <p>{{ truncateMessage(conversation.lastMessagePreview) }}</p>
            <span>{{ formatDate(conversation.lastMessageTimeStamp) }}</span>
          </div>
        </div>
      </div>
      <div v-else>
        <p>Nessuna conversazione disponibile.</p>
        <ErrorBox v-if="ErrorMessage" :msg="ErrorMessage"></ErrorBox>
      </div>
    </div>
  </template>
  

<script>
import defaultProfileImage from '../assets/images/default-profile.jpg';

export default {
    data() {
        return {
            conversations: [],  // Array vuoto che conterrà le conversazioni
            ErrorMessage: '',
        };
    },
    mounted() {
        this.fetchConversations();
    },
    methods: {
        async fetchConversations() {
            try {
                const response = await this.$axios.get("/users/"+sessionStorage.getItem('Auth')+"/conversations",{
                    headers: {
                                Authorization: `Bearer ${sessionStorage.getItem('Auth')}`, 
                            },
                });
                this.conversations = response.data.conversations;
                this.conversations.forEach(conversation => {
                    this.getPhotoConversation(conversation);
                });
            } catch (error) {
                this.ErrorMessage = error.response ? `Error: ${error.response.status} - ${error.response.data}` : "Unexpected error: " + error.message;
            }
        },

        async getPhotoConversation(conversation){
            try {
                if (conversation.photoProfileCode != ""){
                    let response = await this.$axios.get(`/images/${conversation.photoProfileCode}/photo`, {
                        headers: {
                            Authorization: `Bearer ${sessionStorage.getItem('Auth')}`, 
                        },
                        responseType: 'blob',
                    });
                    conversation.photoURL = URL.createObjectURL(response.data);
                } else {
                    conversation.photoURL = defaultProfileImage;
                }
                

            } catch (error) {
                this.ErrorMessage = error.response ? `Error: ${error.response.status} - ${error.response.data}` : "Unexpected error: " + error.message;
            }
        },

        formatDate(dateString) {
            if (dateString == "0000-00-00 00:00:00") {
                return "";
            }
            const options = { year: 'numeric', month: 'short', day: 'numeric', hour: 'numeric', minute: 'numeric' };
            return new Date(dateString).toLocaleString('it-IT', options);
        },
        
        truncateMessage(message) {
            if (message.length > 50) {
                return message.substring(0, 50) + '...';
            }
            return message;
        },

        selectConversation(conversation) {
          this.$emit('conversation-selected', conversation);
      },

    }
};
</script>

<style scoped>
.conversation-item {
  display: flex; /* Usa flexbox per il layout orizzontale */
  align-items: center; /* Allinea verticalmente al centro */
  border-bottom: 1px solid #d0ccd1;
  padding: 10px;
  height: 110px;
}

.conversation-item:hover {
    background-color: #d6d4d4;
}

.profile-image {
  width: 70px; 
  height: 70px; 
  border-radius: 100%; 
  margin-right: 15px;
}

.conversation-details {
  flex-grow: 1;
}

.name-chat {
  margin: 0;
  font-size: 1.4em;
}

p {
  color: #555;
  margin: 5px 0;
}

span {
  font-size: 0.9em;
  color: #888;
}
</style>

  