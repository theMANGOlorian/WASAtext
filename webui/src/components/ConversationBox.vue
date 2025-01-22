<template>
  <div>
    <div class="new-conversation-container">
        <div class="conversation-details">
            <h3 class="name-chat">New conversation</h3>
            <!-- Bottoni per Private e Group -->
            <div class="button-conversation">
                <button 
                    :class="{ active: conversationSelected === 'private' }" 
                    @click="conversationSelected = 'private'">
                    Private
                </button>
                <button 
                    :class="{ active: conversationSelected === 'group' }" 
                    @click="conversationSelected = 'group'">
                    Group
                </button>
            </div>
            
            <div class="add-conversation-container">
                <input 
                    v-model="newConversationName"
                    type="text"
                    :placeholder="placeholderText"
                    @input="fetchUsersList"
                    @keyup.enter="createConversation"
                    @blur="closeDropdown"
                    @focus="showAllUsers"
                    class="conversation-input"
                />
                <button @click="createConversation" class="send-button">Start</button>

                <!-- Dropdown con la lista utenti (solo per conversazioni private)-->
                <ul v-if="filteredUsers.length && this.conversationSelected === 'private'" class="user-dropdown">
                    <li 
                        v-for="user in filteredUsers" 
                        :key="user.id" 
                        @click="selectUser(user)"
                        class="user-item"
                    >
                        {{ user }}
                    </li>
                </ul>
            </div>
        </div>
    </div>
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
          conversations: [],  // Array delle conversazioni
          ErrorMessage: '',
          conversationSelected: 'private',
          newConversationName: '',
          conversationsInterval: null,  // Intervallo per l'aggiornamento delle conversazioni
          previousHash: null, // Salva l'hash della lista di conversazioni
          filteredUsers: [], // Lista utenti filtrata per il dropdown
          debounceTimeout: null, // Timeout per debouncing
      };
  },
  mounted() {
      this.fetchConversations();
      this.conversationsInterval = setInterval(() => {
          this.fetchConversations();
      }, 5000);
  },
  beforeDestroy() {
      // Pulisce l'intervallo quando il componente viene distrutto
      clearInterval(this.conversationsInterval);
  },
  computed: {
      placeholderText() {
          if (this.conversationSelected === 'group') {
            return `Enter group name`; 
          } else {
            return `Enter user name`;
          }
      },
  },
  methods: {

      async showAllUsers() {
          try {
              const response = await this.$axios.get('/users', {
                  headers: { Authorization: `Bearer ${sessionStorage.getItem('Auth')}` },
              });

              // L'API restituisce un array di stringhe
              const allUsers = response.data.users || [];
              this.filteredUsers = allUsers; // Mostra tutti gli utenti
          } catch (error) {
              console.error("Error fetching users:", error.message);
          }
      },

      closeDropdown() {
          setTimeout(() => {
              this.filteredUsers = [];
          }, 100); // Ritardo leggero per permettere il click sugli elementi
      },

      async fetchUsersList() {
        // Blocca la chiamata se l'input è vuoto o troppo corto
        if (this.newConversationName.trim().length < 2) {
          this.filteredUsers = [];
          return;
        }

        // Debounce per evitare troppe chiamate
        clearTimeout(this.debounceTimeout);
        this.debounceTimeout = setTimeout(async () => {
          try {
            const response = await this.$axios.get('/users', {
              headers: { Authorization: `Bearer ${sessionStorage.getItem('Auth')}` },
            });

            // L'API restituisce un array di stringhe
            const allUsers = response.data.users || [];
            this.filteredUsers = allUsers.filter(user =>
              user.toLowerCase().includes(this.newConversationName.trim().toLowerCase())
            );
          } catch (error) {
            console.error("Error fetching users:", error.message);
          }
        }, 300); // Debounce di 300ms
      },

      // Seleziona l'utente dal dropdown
      selectUser(user) {
        this.newConversationName = user;
        this.filteredUsers = []; // Chiude il dropdown
      },

      async createConversation() {
          if (this.newConversationName.trim()) {
            try {
              const response = await this.$axios.post(`/users/${sessionStorage.getItem('Auth')}/conversations`, 
              { 
                name: this.newConversationName,
                typeConversation: this.conversationSelected,
              },
              {
                headers: { Authorization: `Bearer ${sessionStorage.getItem('Auth')}`, },
              });
              this.fetchConversations(); // Aggiorna le conversazioni dopo aver creato
              this.newConversationName = ''; // Pulisce il campo input
              this.filteredUsers = [];
            } catch (error) {
              this.ErrorMessage = error.response ? `Error: ${error.response.status} - ${error.response.data}` : "Unexpected error: " + error.message;
            }
          } else {
            this.ErrorMessage = "Conversation name cannot be empty.";
          }
      },

      async fetchConversations() {
          try {
              const response = await this.$axios.get(`/users/${sessionStorage.getItem('Auth')}/conversations`, {
                  headers: {
                      Authorization: `Bearer ${sessionStorage.getItem('Auth')}`, 
                  },
              });

              const newConversations = response.data.conversations || [];
              const newHash = this.computeHash(newConversations);

              if (this.previousHash !== newHash) {
                  this.previousHash = newHash;
                  this.conversations = newConversations;
                  this.conversations.forEach(conversation => {
                      this.getPhotoConversation(conversation);
                  });
              }
          } catch (error) {
              console.log("fetchConversation: " + error);
          }
      },

      async getPhotoConversation(conversation) {
          try {
              if (conversation.photoProfileCode) {
                  const response = await this.$axios.get(`/images/${conversation.photoProfileCode}/photo`, {
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
              console.log("getPhotoConversation: " + error);
          }
      },

      computeHash(conversations) {
          const jsonString = JSON.stringify(conversations);
          let hash = 0, i, chr;
          for (i = 0; i < jsonString.length; i++) {
              chr = jsonString.charCodeAt(i);
              hash = ((hash << 5) - hash) + chr;
              hash |= 0; // Convert to 32bit integer
          }
          return hash;
      },

      formatDate(dateString) {
          if (dateString === "0000-00-00 00:00:00") {
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
.new-conversation-container {
  display: flex; /* Usa flexbox per il layout orizzontale */
  align-items: center; /* Allinea verticalmente al centro */
  border-bottom: 1px solid #d0ccd1;
  padding: 10px;
  height: 110px;
}

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

.button-conversation {
  display: flex;
  gap: 10px;
  margin-bottom: 10px;
  margin-top: 10px;
}

.button-conversation button {
  width: 70px;
  border: 1px solid #ccc;
  border-radius: 5px;
  background-color: #f1f1f1;
  color: #333;
  cursor: pointer;
  transition: all 0.3s ease;
}

.button-conversation button:hover {
  background-color: #e0e0e0;
}

.button-conversation button.active {
  background-color: #4caf50;
  color: white;
  border-color: #4caf50;
}


.add-conversation-container {
  display: flex;
  gap: 10px;
  margin-top: 10px;
  align-items: center;
}

.conversation-input {
  flex: 1;
  padding: 5px 8px;
  border: 1px solid #ccc;
  border-radius: 5px;
  font-size: 1em;
}

.send-button {
  padding: 5px 10px;
  background-color: #4caf50;
  color: white;
  border: none;
  border-radius: 5px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.send-button:hover {
  background-color: #45a049;
}

.button-conversation button {
  width: 70px;
  border: 1px solid #ccc;
  border-radius: 5px;
  background-color: #f1f1f1;
  color: #333;
  cursor: pointer;
  transition: all 0.3s ease;
}

.button-conversation button:hover {
  background-color: #e0e0e0;
}

.button-conversation button.active {
  background-color: #4caf50;
  color: white;
  border-color: #4caf50;
}

.add-conversation-container {
  display: flex;
  gap: 10px;
  margin-top: 10px;
  align-items: center;
  position: relative; /* Aggiunto per gestire la posizione assoluta del dropdown */
}

.user-dropdown {
  list-style: none;
  background-color: #fff;
  border: 1px solid #ccc;
  border-radius: 5px;
  max-height: 150px;
  overflow-y: auto;
  position: absolute;
  width: 100%; /* Prende tutta la larghezza del contenitore padre */
  z-index: 10;
  top: 100%; /* Posiziona la lista appena sotto l'input */
  left: 0; /* Allinea la lista con il lato sinistro dell'input */
  margin: 0;
  padding: 5px 0;
}

.user-item {
  padding: 8px 12px;
  cursor: pointer;
}

.user-item:hover {
  background-color: #f1f1f1;
}
</style>
