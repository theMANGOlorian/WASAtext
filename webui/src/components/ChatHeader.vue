<template>
    <header class="chat-header" @click.stop="toggleSettingsMenu"> 
        <h1>{{ conversationTitle }}</h1>

        <!-- Dropdown menu (tendina) che appare quando clicchi sull'header -->
        <div v-if="settingsMenuVisible" class="settings-menu" ref="settingsMenu" @click.stop>
            <!-- Sezione sinistra: Foto e nome della chat -->
            <div class="left-section">
                <img :src="conversation.photoURL" alt="Group Photo" class="group-photo"/>
                <span v-if="!isEditingName" style="font-size: 1.9em;">{{ conversation.conversationName }}</span>

                <!-- Sezione per l'input per modificare il nome -->
                <div v-if="isEditingName">
                    <input 
                        v-model="newGroupName" 
                        type="text" 
                        class="name-input" 
                        placeholder="Enter new name"
                        @keyup.enter="updateGroupName"
                    />
                    <button @click="updateGroupName">Save</button>
                </div>
            </div>
            <div v-if="isGroup">
                <!-- Sezione destra: Bottoni per le azioni -->
                <div v-if="showingUserList" class="user-list">
                    <!-- Ciclo sugli utenti per creare un elenco -->
                    <div 
                        v-for="user in userList" 
                        :key="user" 
                        class="user-item" 
                        @click="addUserToGroup(user)" >
                        {{ user }}
                    </div>
                    <button @click="hideUserList" style="color: red;">Back</button>
                </div>
                <div v-else>
                    <div class="menu-item" @click="addPerson">Add Member</div>
                    <div class="menu-item" @click="changeName">Change Name</div>
                    <div class="menu-item">
                        <input type="file" accept="image/png" @change="handlePhotoChange" style="display: none;" ref="photoInput" />
                        <button @click="triggerFileInput">Change Photo</button>
                    </div>
                    <div class="menu-item" style="color: red;" @click="exitGroup">Exit Group</div>
                </div>
            </div>
        </div>
    </header>
</template>



<script>
import defaultProfileImage from '../assets/images/default-profile.jpg';
export default {
    props: {
        fatherConversation: {
            type: Object,
            required: true,
        },
    },
    data() {
        return {
            settingsMenuVisible: false, // Stato per la visibilità del menu delle impostazioni
            isGroup: false,
            isEditingName: false,
            newGroupName: '',
            userList: [],
            showingUserList: false,
            conversation: this.fatherConversation,
        };
    },

    watch: {
        fatherConversation: {
            handler(newVal) {
                // Quando la prop 'ObjConv' cambia, aggiorna 'conversation'
                this.conversation = { ...newVal };
            },
            deep: true,
        },
    },

    computed: {
        conversationTitle() {
            return this.conversation ? this.conversation.conversationName : "Chat";
        },
    },
    methods: {

        async addUserToGroup(user) {
            try {
                const response = await this.$axios.put(
                    `/conversations/${this.conversation.conversationId}/member`,
                    { Username: user }, // Corpo della richiesta
                    {
                        headers: {
                            Authorization: `Bearer ${sessionStorage.getItem('Auth')}`,
                            "Content-Type": "application/json",
                        },
                    }
                );
                alert(`Utente ${user} aggiunto con successo al gruppo!`);
            } catch (error) {
                console.error("Errore durante l'aggiunta dell'utente al gruppo:", error);
                alert(`Errore: impossibile aggiungere ${user} al gruppo.`);
            }
        },

        async addPerson() {
            try {
                const response = await this.$axios.get('/users', {
                    headers: {
                        Authorization: `Bearer ${sessionStorage.getItem('Auth')}`,
                    },
                });
                this.userList = response.data.users; // Salva i dati degli utenti
                this.showingUserList = true; // Mostra la lista utenti
            } catch (error) {
                console.error("Errore nel recuperare gli utenti:", error);
            }
        },

        // Nasconde la lista degli utenti
        hideUserList() {
            this.showingUserList = false;
        },

        // Toggle visibility della tendina
        toggleSettingsMenu() {
            if (this.conversation.conversationType === 'group' || this.conversation.conversationType === 'private'){

                if (this.conversation.conversationType === 'group'){
                    this.isGroup = true; 
                } else {
                    this.isGroup = false; 
                }

                this.settingsMenuVisible = !this.settingsMenuVisible;

                if (this.settingsMenuVisible) {
                    // Aggiungi il listener per il clic esterno quando il menu è visibile
                    document.addEventListener('click', this.handleClickOutside);
                    // Effettua la richiesta per ottenere l'immagine solo quando il menu è visibile
                    if (this.conversation.photoProfileCode) {
                        this.getPhotoConversation();
                    }
                } else {
                    // Rimuovi il listener quando il menu è chiuso
                    document.removeEventListener('click', this.handleClickOutside);
                }
            } else {
                this.settingsMenuVisible = false;
            }
        },

        // Controlla se il clic è fuori dal menu
        handleClickOutside(event) {
            if (this.settingsMenuVisible && !this.$refs.settingsMenu.contains(event.target)) {
                this.settingsMenuVisible = false;
            }
        },

        // Trigger per aprire il selettore di file
        triggerFileInput() {
            this.$refs.photoInput.click();
        },

        // Gestisce il cambiamento della foto
        async handlePhotoChange(event) {
            const file = event.target.files[0];

            if (!file) {
                return;
            }

            // Passa il file a SetPhotoGroup
            await this.SetPhotoGroup(file);
        },

        async getPhotoConversation() {
            try {
                if (this.conversation.photoProfileCode) {
                    const response = await this.$axios.get(`/images/${this.conversation.photoProfileCode}/photo`, {
                        headers: {
                            Authorization: `Bearer ${sessionStorage.getItem('Auth')}`, 
                        },
                        responseType: 'blob',
                    });
                    // Imposta la photoURL
                    this.conversation.photoURL = URL.createObjectURL(response.data);
                } else {
                    this.conversation.photoURL = defaultProfileImage;
                }
            } catch (error) {
                console.log("getPhotoConversation: " + error);
                // Se c'è un errore, imposta l'immagine di default
                this.conversation.photoURL = defaultProfileImage;
            }
        },

        // Metodo per cambiare la foto del gruppo
        async SetPhotoGroup(file) {
            try {
                if (file.type !== 'image/png') {
                    this.ErrorMessage = "Only PNG images are allowed.";
                    return;
                }
                let response = await this.$axios.put(`/conversations/${this.conversation.conversationId}/groupPhoto`, file, {
                    headers: {
                        Authorization: `Bearer ${sessionStorage.getItem('Auth')}`,
                        "Content-Type": "image/png",
                    },
                });

                this.conversation.photoURL = URL.createObjectURL(file);
                this.conversation.photoProfileCode = response.data.imageCode;
                console.log("Foto cambiata con successo.");
            } catch(error) {
                if (error.response) {
                    this.ErrorMessage = `Error: ${error.response.status} - ${error.response.data}`;
                } else if (error.request) {
                    this.ErrorMessage = "Network error, please try again.";
                } else {
                    this.ErrorMessage = "Unexpected error: " + error.message;
                }
            }
        },

        async updateGroupName() {
            if (this.newGroupName !== this.conversation.conversationName){
                try {
                let response = await this.$axios.put(`/conversations/${this.conversation.conversationId}/groupName`, { name: this.newGroupName }, {
                    headers: {
                        Authorization: `Bearer ${sessionStorage.getItem('Auth')}`,
                        "Content-Type": "application/json",
                    },
                });
                this.conversation.conversationName = this.newGroupName;
                this.isEditingName = false;
            } catch (error) {
                this.ErrorMessage = error.response ? `Error: ${error.response.status} - ${error.response.data}` : "Unexpected error: " + error.message;
            }
            } else {
                this.isEditingName = false;
            }
        },
        
        changeName() {
            this.isEditingName = true;
        },
        
        async exitGroup() {
            if (!confirm("Sei sicuro di voler uscire dal gruppo?")) {
                return; // Interrompi l'azione se l'utente annulla
            }

            try {
                const response = await this.$axios.delete(
                    `/conversations/${this.conversation.conversationId}`,
                    {
                        headers: {
                            Authorization: `Bearer ${sessionStorage.getItem('Auth')}`,
                        },
                    }
                );
                this.$emit("leftGroup", null);
                console.log("Sei uscito dal gruppo con successo!");
            } catch (error) {
                console.error("Errore durante l'uscita dal gruppo:", error);
                alert("Errore: impossibile uscire dal gruppo.");
            }
        },

    },

    beforeUnmount() {
        // Rimuovi il listener per evitare memory leaks quando il componente è distrutto
        document.removeEventListener('click', this.handleClickOutside);
    },
};

</script>


<style scoped>
.chat-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 10px;
    background-color: #ffffff;
    color: rgb(0, 0, 0);
    position: relative;
    cursor: pointer;
}

.chat-header h1 {
    margin: 0;
    font-size: 1.5em;
}

.chat-header button {
    background: none;
    border: none;
    color: rgb(0, 0, 0);
    cursor: pointer;
    font-size: 1.2em;
}

.buttonSettings:hover {
    background-color: whitesmoke;
}

/* Stile per il menu delle impostazioni */
.settings-menu {
    position: absolute;
    top: 40px; /* Spostato sotto la header */
    left: 50%; /* Centra la tendina orizzontalmente */
    transform: translateX(-50%); /* Centra perfettamente la tendina */
    background-color: #ffffff;
    border: 1px solid #ddd;
    border-radius: 8px;
    box-shadow: 0px 4px 6px rgba(0, 0, 0, 0.1);
    width: 600px; /* Aumenta la larghezza della tendina */
    padding: 80px;
    z-index: 100;
    display: flex; /* Dispone gli elementi in due colonne */
    justify-content: space-between; /* Distanza fra le due sezioni */
}

/* Sezione sinistra (foto e nome del gruppo) */
.left-section {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    text-align: center;
}

.group-photo {
    width: 200px;
    height: 200px;
    border-radius: 50%;
    margin-bottom: 10px;
}

.left-section span {
    font-size: 1.2em;
}

/* Sezione destra (bottoni) */
.right-section {
    display: flex;
    flex-direction: column;
    justify-content: flex-start;
    align-items: flex-start;
}

.menu-item {
    padding: 10px;
    font-size: 1.1em;
    cursor: pointer;
    width: 200px; /* Imposta la larghezza fissa per i bottoni */
    text-align: center;
}
.menu-item button {
    font-size: 1.0em;
}


.menu-item:hover {
    background-color: #f0f0f0;
}


.user-item {
    cursor: pointer;
    font-size: 1.1em;
    width: 100px;
    text-align: center;
}

.user-item:hover{
    color: green;
    background-color: #e9e9e9;
}

</style>
